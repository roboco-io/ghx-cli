package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// ExportOptions holds options for the export command
type ExportOptions struct {
	Output           string
	Format           string
	IncludeItems     bool
	IncludeFields    bool
	IncludeViews     bool
	IncludeWorkflows bool
}

// NewExportCmd creates the export command
func NewExportCmd() *cobra.Command {
	opts := &ExportOptions{}

	cmd := &cobra.Command{
		Use:   "export PROJECT_ID",
		Short: "Export project data to a file",
		Long: `Export GitHub Project data including configuration, items, fields, and workflows.

This creates a backup file that can be used to restore the project or migrate to another location.

Examples:
  ghx project export myorg/123 --output project-backup.json
  ghx project export user/456 --output backup.json --format yaml
  ghx project export myorg/123 --output full-backup.json --include-all`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(cmd.Context(), opts, args)
		},
	}

	cmd.Flags().StringVar(&opts.Output, "output", "", "Output file path (required)")
	cmd.Flags().StringVar(&opts.Format, "format", "json", "Export format: json, yaml")
	cmd.Flags().BoolVar(&opts.IncludeItems, "include-items", true, "Include project items")
	cmd.Flags().BoolVar(&opts.IncludeFields, "include-fields", true, "Include custom fields")
	cmd.Flags().BoolVar(&opts.IncludeViews, "include-views", true, "Include project views")
	cmd.Flags().BoolVar(&opts.IncludeWorkflows, "include-workflows", true, "Include automation workflows")

	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func runExport(ctx context.Context, opts *ExportOptions, args []string) error {
	projectID := args[0]

	// Validate format
	if opts.Format != formatJSON && opts.Format != formatYAML {
		return fmt.Errorf("unsupported format: %s (supported: %s, %s)", opts.Format, formatJSON, formatYAML)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(opts.Output), dirPerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)

	// Export project
	exportData := &service.ProjectExportData{
		ProjectID:        projectID,
		IncludeItems:     opts.IncludeItems,
		IncludeFields:    opts.IncludeFields,
		IncludeViews:     opts.IncludeViews,
		IncludeWorkflows: opts.IncludeWorkflows,
	}

	err = projectService.ExportProject(ctx, exportData, opts.Output, opts.Format)
	if err != nil {
		return fmt.Errorf("failed to export project: %w", err)
	}

	// Calculate file size
	fileInfo, err := os.Stat(opts.Output)
	if err != nil {
		return fmt.Errorf("failed to get export file info: %w", err)
	}

	fmt.Printf("✅ Successfully exported project %s\n", projectID)
	fmt.Printf("   Output: %s\n", opts.Output)
	fmt.Printf("   Format: %s\n", strings.ToUpper(opts.Format))
	fmt.Printf("   Size: %d bytes\n", fileInfo.Size())

	return nil
}
