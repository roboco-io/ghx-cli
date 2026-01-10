package analytics

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// ExportOptions holds options for the export command
type ExportOptions struct {
	ProjectRef       string
	Format           string
	OutputFormat     string
	Filter           string
	IncludeItems     bool
	IncludeFields    bool
	IncludeViews     bool
	IncludeWorkflows bool
}

// NewExportCmd creates the export command
func NewExportCmd() *cobra.Command {
	opts := &ExportOptions{}

	cmd := &cobra.Command{
		Use:   "export <owner/project-number>",
		Short: "Export project data",
		Long: `Export GitHub Project data in various formats.

Export project data including items, fields, views, and workflows to different
formats for backup, analysis, or migration purposes.

Export Formats:
  json         - JSON format for programmatic access and API integration
  csv          - CSV format for spreadsheet analysis and reporting
  xml          - XML format for structured data exchange and integration

Include Options:
  --include-items      Include project items (issues, pull requests, draft items)
  --include-fields     Include custom fields and field configurations
  --include-views      Include project views and their configurations
  --include-workflows  Include automation workflows and their rules
  --include-all        Include all available data (items, fields, views, workflows)

Filter Options:
  --filter             Apply filter to limit exported items (e.g., "status:open", "assignee:octocat")

Examples:
  ghx analytics export octocat/123 --format json --include-all
  ghx analytics export octocat/123 --format csv --include-items --include-fields
  ghx analytics export octocat/123 --format xml --filter "status:open" --output json
  ghx analytics export --org myorg/456 --format json --include-workflows`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.OutputFormat = cmd.Flag("format").Value.String()

			// Handle include-all flag
			if includeAll, _ := cmd.Flags().GetBool("include-all"); includeAll {
				opts.IncludeItems = true
				opts.IncludeFields = true
				opts.IncludeViews = true
				opts.IncludeWorkflows = true
			}

			return runExport(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "json", "Export format (json, csv, xml)")
	cmd.Flags().BoolVar(&opts.IncludeItems, "include-items", false, "Include project items")
	cmd.Flags().BoolVar(&opts.IncludeFields, "include-fields", false, "Include custom fields")
	cmd.Flags().BoolVar(&opts.IncludeViews, "include-views", false, "Include project views")
	cmd.Flags().BoolVar(&opts.IncludeWorkflows, "include-workflows", false, "Include workflows")
	cmd.Flags().Bool("include-all", false, "Include all available data")
	cmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter for exported items")
	cmd.Flags().Bool("org", false, "Target organization project")

	return cmd
}

func runExport(ctx context.Context, opts *ExportOptions) error {
	// Parse project reference
	parts := strings.Split(opts.ProjectRef, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid project reference format. Use: owner/project-number")
	}

	owner := parts[0]
	projectNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid project number: %s", parts[1])
	}

	// Validate export format
	exportFormat, err := service.ValidateExportFormat(opts.Format)
	if err != nil {
		return err
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
	analyticsService := service.NewAnalyticsService(client)

	// Get project to validate access and get project ID (with automatic owner detection)
	project, err := projectService.GetProjectWithOwnerDetection(ctx, owner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Prepare export input
	input := service.ExportProjectInput{
		ProjectID:        project.ID,
		Format:           exportFormat,
		IncludeItems:     opts.IncludeItems,
		IncludeFields:    opts.IncludeFields,
		IncludeViews:     opts.IncludeViews,
		IncludeWorkflows: opts.IncludeWorkflows,
	}

	if opts.Filter != "" {
		input.Filter = &opts.Filter
	}

	// Export project
	export, err := analyticsService.ExportProject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to export project: %w", err)
	}

	// Output export result
	return outputExport(export, opts.OutputFormat)
}

func outputExport(export *service.ProjectV2Export, format string) error {
	switch format {
	case FormatJSON:
		return outputExportJSON(export)
	case FormatTable:
		return outputExportTable(export)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}

func outputExportTable(export *service.ProjectV2Export) error {
	fmt.Printf("✅ Project export completed successfully\n\n")

	fmt.Printf("Export Details:\n")
	fmt.Printf("  Project ID: %s\n", export.ProjectID)
	fmt.Printf("  Title: %s\n", export.Title)
	if export.Description != nil {
		fmt.Printf("  Description: %s\n", *export.Description)
	}
	fmt.Printf("  Export Date: %s\n", export.ExportDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Format: %s\n", service.FormatExportFormat(export.Format))

	fmt.Printf("\nExported Data:\n")
	fmt.Printf("  Items: %d\n", len(export.Items))
	fmt.Printf("  Fields: %d\n", len(export.Fields))
	fmt.Printf("  Views: %d\n", len(export.Views))
	fmt.Printf("  Workflows: %d\n", len(export.Workflows))

	// Show sample items if available
	if len(export.Items) > 0 {
		fmt.Printf("\nSample Items:\n")
		maxItems := 3
		if len(export.Items) < maxItems {
			maxItems = len(export.Items)
		}
		for i := 0; i < maxItems; i++ {
			item := export.Items[i]
			fmt.Printf("  • %s (%s) - %s\n", item.Title, item.Type, item.State)
		}
		if len(export.Items) > maxItems {
			fmt.Printf("  ... and %d more items\n", len(export.Items)-maxItems)
		}
	}

	return nil
}

func outputExportJSON(export *service.ProjectV2Export) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"success\": true,\n")
	fmt.Printf("  \"projectId\": \"%s\",\n", export.ProjectID)
	fmt.Printf("  \"title\": \"%s\",\n", export.Title)

	if export.Description != nil {
		fmt.Printf("  \"description\": \"%s\",\n", *export.Description)
	}

	fmt.Printf("  \"exportDate\": \"%s\",\n", export.ExportDate.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("  \"format\": \"%s\",\n", export.Format)
	fmt.Printf("  \"itemCount\": %d,\n", len(export.Items))
	fmt.Printf("  \"fieldCount\": %d,\n", len(export.Fields))
	fmt.Printf("  \"viewCount\": %d,\n", len(export.Views))
	fmt.Printf("  \"workflowCount\": %d\n", len(export.Workflows))
	fmt.Printf("}\n")

	return nil
}
