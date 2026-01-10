package item

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// EditOptions holds options for the edit command
type EditOptions struct {
	ProjectRef string
	ItemID     string
	FieldName  string
	Value      string
	Format     string
}

// NewEditCmd creates the edit command
func NewEditCmd() *cobra.Command {
	opts := &EditOptions{}

	cmd := &cobra.Command{
		Use:   "edit <project> <item-id> --field <field-name> --value <value>",
		Short: "Edit item field values",
		Long: `Edit field values for items in a project.

This command allows you to update custom field values for project items.
You need to specify the project-specific item ID (not the issue/PR ID).

Field values can be:
• Text values for text fields
• Numbers for number fields
• Option names for single-select fields
• Dates in YYYY-MM-DD format for date fields

Examples:
  ghx item edit octocat/1 PVTI_123 --field "Status" --value "In Progress"
  ghx item edit myorg/2 item-456 --field "Priority" --value "High"
  ghx item edit octocat/1 PVTI_789 --field "Due Date" --value "2024-12-31"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.ItemID = args[1]
			return runEdit(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.FieldName, "field", "", "Field name to update (required)")
	cmd.Flags().StringVar(&opts.Value, "value", "", "New field value (required)")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	_ = cmd.MarkFlagRequired("field")
	_ = cmd.MarkFlagRequired("value")

	return cmd
}

func runEdit(ctx context.Context, opts *EditOptions) error {
	// Parse project reference
	projectOwner, projectNumber, err := service.ParseProjectReference(opts.ProjectRef)
	if err != nil {
		return fmt.Errorf("invalid project reference: %w", err)
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)

	// Get project details to find field ID
	project, err := projectService.GetProjectWithOwnerDetection(ctx, projectOwner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Find the field by name
	var fieldID string
	for _, field := range project.Fields.Nodes {
		if field.Name == opts.FieldName {
			fieldID = field.ID
			break
		}
	}

	if fieldID == "" {
		return fmt.Errorf("field '%s' not found in project. Available fields", opts.FieldName)
	}

	// Prepare field value based on field type
	var fieldValue interface{}
	// For now, we'll treat all values as strings
	// In a full implementation, we'd parse based on field type
	fieldValue = opts.Value

	// If value looks like a number, try to convert it
	if numValue, parseErr := strconv.ParseFloat(opts.Value, 64); parseErr == nil {
		fieldValue = numValue
	}

	// Update item field
	input := service.UpdateItemFieldInput{
		ProjectID: project.ID,
		ItemID:    opts.ItemID,
		FieldID:   fieldID,
		Value:     fieldValue,
	}

	item, err := projectService.UpdateItemField(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item field: %w", err)
	}

	fmt.Printf("✅ Field '%s' updated successfully!\n\n", opts.FieldName)
	return outputUpdatedItem(item, opts.Format, opts.FieldName, opts.Value)
}

func outputUpdatedItem(_ interface{}, format, fieldName, value string) error {
	switch format {
	case formatJSON:
		return outputUpdatedItemJSON(fieldName, value)
	case "table":
		return outputUpdatedItemTable(fieldName, value)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputUpdatedItemTable(fieldName, value string) error {
	fmt.Printf("Field: %s\n", fieldName)
	fmt.Printf("New Value: %s\n", value)
	return nil
}

func outputUpdatedItemJSON(fieldName, value string) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"status\": \"updated\",\n")
	fmt.Printf("  \"field\": \"%s\",\n", fieldName)
	fmt.Printf("  \"value\": \"%s\"\n", value)
	fmt.Printf("}\n")
	return nil
}
