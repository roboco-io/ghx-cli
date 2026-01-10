package analytics

import (
	"github.com/spf13/cobra"
)

// NewAnalyticsCmd creates the analytics command
func NewAnalyticsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analytics",
		Short: "Generate analytics and reports for GitHub Projects",
		Long: `Generate analytics and reports for GitHub Projects.

The analytics command provides comprehensive reporting and analysis capabilities
for GitHub Projects v2. You can generate various types of reports including:

• Project overview statistics and metrics
• Item distribution by status, assignee, labels, and milestones  
• Velocity and performance metrics
• Timeline analysis and milestone tracking
• Export project data in multiple formats (JSON, CSV, XML)
• Import project data with merge strategies
• Bulk operations on project items

Analytics Types:
  overview     - Project overview with item counts and basic statistics
  velocity     - Team velocity and performance metrics over time
  timeline     - Project timeline with milestones and activity analysis
  distribution - Item distribution across statuses, assignees, and labels

Export Formats:
  json         - JSON format for programmatic access
  csv          - CSV format for spreadsheet analysis
  xml          - XML format for structured data exchange

Bulk Operations:
  update       - Update multiple items at once
  delete       - Delete multiple items at once
  archive      - Archive multiple items at once

Import Strategies:
  merge        - Merge imported data with existing items
  replace      - Replace existing items with imported data
  append       - Add imported items without modifying existing ones
  skip_conflicts - Skip items that would cause conflicts

Examples:
  ghx analytics overview octocat/123
  ghx analytics velocity octocat/123 --period monthly
  ghx analytics export octocat/123 --format json --include-all
  ghx analytics import octocat/123 --file data.json --strategy merge
  ghx analytics bulk-update octocat/123 --items item1,item2 --field status --value Done`,

		Args: cobra.NoArgs,
	}

	// Add subcommands
	cmd.AddCommand(NewOverviewCmd())
	cmd.AddCommand(NewVelocityCmd())
	cmd.AddCommand(NewTimelineCmd())
	cmd.AddCommand(NewDistributionCmd())
	cmd.AddCommand(NewExportCmd())
	cmd.AddCommand(NewImportCmd())
	cmd.AddCommand(NewBulkUpdateCmd())
	cmd.AddCommand(NewBulkDeleteCmd())
	cmd.AddCommand(NewBulkArchiveCmd())
	cmd.AddCommand(NewOperationStatusCmd())

	return cmd
}
