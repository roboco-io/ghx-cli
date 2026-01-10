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

// AddOptionOptions holds options for the add-option command
type AddOptionOptions struct {
	FieldID     string
	Name        string
	Color       string
	Description string
	Format      string
}

// NewAddOptionCmd creates the add-option command
func NewAddOptionCmd() *cobra.Command {
	opts := &AddOptionOptions{}

	cmd := &cobra.Command{
		Use:   "add-option <field-id> <name>",
		Short: "Add option to single select field",
		Long: `Add a new option to a single select field.

This command only works with single select fields. You can specify the color
and an optional description for the new option.

Available colors: gray, red, orange, yellow, green, blue, purple, pink

Examples:
  ghx field add-option field-id "Critical"
  ghx field add-option field-id "High" --color red
  ghx field add-option field-id "Urgent" --color red --description "Requires immediate attention"`,

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FieldID = args[0]
			opts.Name = args[1]
			opts.Format = cmd.Flag("format").Value.String()
			return runAddOption(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Color, "color", "gray", "Color for the option (gray, red, orange, yellow, green, blue, purple, pink)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Optional description for the option")

	return cmd
}

func runAddOption(ctx context.Context, opts *AddOptionOptions) error {
	// Validate color
	if err := service.ValidateColor(opts.Color); err != nil {
		return err
	}

	// Normalize color
	normalizedColor := service.NormalizeColor(opts.Color)

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	fieldService := service.NewFieldService(client)

	// Create field option
	var description *string
	if opts.Description != "" {
		description = &opts.Description
	}

	input := service.CreateFieldOptionInput{
		FieldID:     opts.FieldID,
		Name:        opts.Name,
		Color:       normalizedColor,
		Description: description,
	}

	option, err := fieldService.CreateFieldOption(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create field option: %w", err)
	}

	// Output created option
	return outputCreatedOption(option, opts.Format)
}

func outputCreatedOption(option *graphql.ProjectV2SingleSelectFieldOption, format string) error {
	switch format {
	case formatJSON:
		return outputCreatedOptionJSON(option)
	case formatTable:
		return outputCreatedOptionTable(option)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputCreatedOptionTable(option *graphql.ProjectV2SingleSelectFieldOption) error {
	fmt.Printf("✅ Option '%s' added successfully\n\n", option.Name)

	fmt.Printf("Option Details:\n")
	fmt.Printf("  ID: %s\n", option.ID)
	fmt.Printf("  Name: %s\n", option.Name)
	fmt.Printf("  Color: %s\n", service.FormatColor(option.Color))

	if option.Description != nil && *option.Description != "" {
		fmt.Printf("  Description: %s\n", *option.Description)
	}

	return nil
}

func outputCreatedOptionJSON(option *graphql.ProjectV2SingleSelectFieldOption) error {
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
