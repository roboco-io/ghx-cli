package view

import (
	"github.com/spf13/cobra"
)

// NewViewCmd creates the view command
func NewViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "Manage project views",
		Long: `Manage views in GitHub Projects.

Views provide different perspectives on your project data, allowing you to 
organize and visualize items in ways that best suit your workflow. GitHub 
Projects supports multiple view types:

• Table views - Traditional table layout with customizable columns
• Board views - Kanban-style boards with swimlanes
• Roadmap views - Timeline-based planning views

This command group provides comprehensive view management capabilities:

• List existing views in projects
• Create new views with different layouts
• Update view names and filters
• Copy views to create variations
• Delete views when no longer needed
• Configure view sorting and grouping

View Layouts:
  table       - Table view with customizable columns
  board       - Kanban board view with card layout  
  roadmap     - Timeline roadmap for planning

View Operations:
  list        - List all views in a project
  create      - Create a new project view
  update      - Update view name or filter
  copy        - Create a copy of an existing view
  delete      - Delete a project view
  sort        - Configure view sorting options
  group       - Configure view grouping options`,

		Example: `  # List all views in a project
  ghx view list octocat/123

  # Create a new table view
  ghx view create octocat/123 "Sprint Dashboard" table

  # Create a board view with filter
  ghx view create octocat/123 "Bug Tracking" board --filter "label:bug"

  # Copy an existing view
  ghx view copy view-id "New Sprint Board"

  # Update view name
  ghx view update view-id --name "Updated Dashboard"

  # Configure view sorting
  ghx view sort view-id --field priority --direction desc`,
	}

	// Add format flag to all subcommands
	cmd.PersistentFlags().String("format", "table", "Output format (table, json)")

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewUpdateCmd())
	cmd.AddCommand(NewCopyCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewSortCmd())
	cmd.AddCommand(NewGroupCmd())

	return cmd
}
