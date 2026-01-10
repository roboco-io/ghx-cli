package field

import (
	"github.com/spf13/cobra"
)

// NewFieldCmd creates the field command group
func NewFieldCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "field",
		Short: "Manage project fields",
		Long: `Manage custom fields in GitHub Projects.

Fields allow you to track additional metadata for your project items.
GitHub Projects supports different field types including text, number,
date, single select, and iteration fields.

This command group provides comprehensive field management capabilities:

• Create new custom fields with various data types
• List and view existing fields in projects
• Update field names and properties
• Delete fields from projects
• Manage single select field options (add, update, delete)

Field Types:
  text         - Text field for arbitrary text input
  number       - Numeric field for numbers and calculations  
  date         - Date field for deadlines and milestones
  single_select - Single select field with predefined options
  iteration    - Iteration field for sprint/cycle planning

For more information about GitHub Projects fields, visit:
https://docs.github.com/en/issues/planning-and-tracking-with-projects`,

		Example: `  ghx field list octocat/123                    # List fields in project
  ghx field create octocat/123 "Priority" text     # Create text field
  ghx field create octocat/123 "Status" single_select --options "Todo,In Progress,Done"
  ghx field update field-id --name "New Priority"  # Rename field
  ghx field delete field-id --force                # Delete field
  ghx field add-option field-id "Critical" --color red  # Add select option`,
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewAddOptionCmd())
	cmd.AddCommand(NewUpdateOptionCmd())
	cmd.AddCommand(NewDeleteOptionCmd())

	return cmd
}
