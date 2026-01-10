package field

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// UpdateOptionOptions holds options for the update-option command
type UpdateOptionOptions struct {
	OptionID    string
	Name        string
	Color       string
	Description string
	Format      string
}

// NewUpdateOptionCmd creates the update-option command
func NewUpdateOptionCmd() *cobra.Command {
	opts := &UpdateOptionOptions{}

	cmd := &cobra.Command{
		Use:   "update-option <option-id>",
		Short: "Update single select field option",
		Long: `Update properties of a single select field option.

You can update the name, color, and description of existing options.
At least one property must be specified.

Available colors: gray, red, orange, yellow, green, blue, purple, pink

Examples:
  ghx field update-option option-id --name "Very High"
  ghx field update-option option-id --color red
  ghx field update-option option-id --name "Critical" --color red --description "Highest priority"`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.OptionID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runUpdateOption(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the option")
	cmd.Flags().StringVar(&opts.Color, "color", "", "New color for the option")
	cmd.Flags().StringVar(&opts.Description, "description", "", "New description for the option")

	return cmd
}

func runUpdateOption(ctx context.Context, opts *UpdateOptionOptions) error {
	// Validate at least one field is provided
	if opts.Name == "" && opts.Color == "" && opts.Description == "" {
		return fmt.Errorf("at least one of --name, --color, or --description must be provided")
	}

	// Validate color if provided
	var normalizedColor *string
	if opts.Color != "" {
		if err := service.ValidateColor(opts.Color); err != nil {
			return err
		}
		color := service.NormalizeColor(opts.Color)
		normalizedColor = &color
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	fieldService := service.NewFieldService(client)

	// Prepare input
	input := service.UpdateFieldOptionInput{
		OptionID: opts.OptionID,
	}

	if opts.Name != "" {
		input.Name = &opts.Name
	}
	if normalizedColor != nil {
		input.Color = normalizedColor
	}
	if opts.Description != "" {
		input.Description = &opts.Description
	}

	// Update field option
	option, err := fieldService.UpdateFieldOption(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update field option: %w", err)
	}

	// Output updated option
	return outputUpdatedOption(option, opts.Format)
}

func outputUpdatedOption(option *graphql.ProjectV2SingleSelectFieldOption, format string) error {
	switch format {
	case "json":
		return outputUpdatedOptionJSON(option)
	case "table":
		return outputUpdatedOptionTable(option)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputUpdatedOptionTable(option *graphql.ProjectV2SingleSelectFieldOption) error {
	fmt.Printf("✅ Option '%s' updated successfully\n\n", option.Name)

	fmt.Printf("Option Details:\n")
	fmt.Printf("  ID: %s\n", option.ID)
	fmt.Printf("  Name: %s\n", option.Name)
	fmt.Printf("  Color: %s\n", service.FormatColor(option.Color))

	if option.Description != nil && *option.Description != "" {
		fmt.Printf("  Description: %s\n", *option.Description)
	}

	return nil
}

func outputUpdatedOptionJSON(option *graphql.ProjectV2SingleSelectFieldOption) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", option.ID)
	fmt.Printf("  \"name\": \"%s\",\n", option.Name)
	fmt.Printf("  \"color\": \"%s\"", option.Color)

	if option.Description != nil {
		fmt.Printf(",\n  \"description\": \"%s\"", *option.Description)
	}

	fmt.Printf("\n}\n")

	return nil
}
