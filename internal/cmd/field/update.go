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

// UpdateOptions holds options for the update command
type UpdateOptions struct {
	FieldID string
	Name    string
	Format  string
}

// NewUpdateCmd creates the update command
func NewUpdateCmd() *cobra.Command {
	opts := &UpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <field-id>",
		Short: "Update a project field",
		Long: `Update properties of an existing project field.

Currently, you can update the field name. Other field properties like
data type cannot be changed after creation.

Examples:
  ghx field update field-id --name "New Priority"
  ghx field update field-id --name "Status Category" --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FieldID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runUpdate(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the field")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func runUpdate(ctx context.Context, opts *UpdateOptions) error {
	// Validate field name
	if err := service.ValidateFieldName(opts.Name); err != nil {
		return err
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

	// Update field
	input := service.UpdateFieldInput{
		FieldID: opts.FieldID,
		Name:    &opts.Name,
	}

	field, err := fieldService.UpdateField(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update field: %w", err)
	}

	// Output updated field
	return outputUpdatedField(field, opts.Format)
}

func outputUpdatedField(field *graphql.ProjectV2Field, format string) error {
	switch format {
	case "json":
		return outputUpdatedFieldJSON(field)
	case "table":
		return outputUpdatedFieldTable(field)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputUpdatedFieldTable(field *graphql.ProjectV2Field) error {
	fmt.Printf("✅ Field updated successfully\n\n")

	fmt.Printf("Field Details:\n")
	fmt.Printf("  ID: %s\n", field.ID)
	fmt.Printf("  Name: %s\n", field.Name)
	fmt.Printf("  Type: %s\n", service.FormatFieldDataType(field.DataType))

	return nil
}

func outputUpdatedFieldJSON(field *graphql.ProjectV2Field) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", field.ID)
	fmt.Printf("  \"name\": \"%s\",\n", field.Name)
	fmt.Printf("  \"dataType\": \"%s\"\n", field.DataType)
	fmt.Printf("}\n")

	return nil
}
