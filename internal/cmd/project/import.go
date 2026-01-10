package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// ImportOptions holds options for the import command
type ImportOptions struct {
	File       string
	Owner      string
	DryRun     bool
	SkipItems  bool
	SkipFields bool
}

// NewImportCmd creates the import command
func NewImportCmd() *cobra.Command {
	opts := &ImportOptions{}

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import project data from a file",
		Long: `Import project data from a previously exported backup file.

This command recreates a project with all its configuration, items, fields, and workflows
from an exported backup file.

Examples:
  ghx project import --file project-backup.json --owner myorg
  ghx project import --file backup.json --owner myuser --dry-run
  ghx project import --file export.json --owner myorg --skip-items`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runImport(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.File, "file", "", "Import file path (required)")
	cmd.Flags().StringVar(&opts.Owner, "owner", "", "Project owner (user or organization)")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview import without making changes")
	cmd.Flags().BoolVar(&opts.SkipItems, "skip-items", false, "Skip importing project items")
	cmd.Flags().BoolVar(&opts.SkipFields, "skip-fields", false, "Skip importing custom fields")

	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("owner")

	return cmd
}

func runImport(ctx context.Context, opts *ImportOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)

	// Import project
	importOptions := &service.ProjectImportOptions{
		File:       opts.File,
		Owner:      opts.Owner,
		DryRun:     opts.DryRun,
		SkipItems:  opts.SkipItems,
		SkipFields: opts.SkipFields,
	}

	result, err := projectService.ImportProject(ctx, importOptions)
	if err != nil {
		return fmt.Errorf("failed to import project: %w", err)
	}

	if opts.DryRun {
		fmt.Printf("🔍 Dry run completed\n\n")
		fmt.Printf("Would create project: %s\n", result.ProjectTitle)
		fmt.Printf("Items to import: %d\n", result.ItemCount)
		fmt.Printf("Fields to import: %d\n", result.FieldCount)
		fmt.Printf("Views to import: %d\n", result.ViewCount)
	} else {
		fmt.Printf("✅ Successfully imported project\n\n")
		fmt.Printf("Project ID: %s\n", result.ProjectID)
		fmt.Printf("Project URL: %s\n", result.ProjectURL)
		fmt.Printf("Items imported: %d\n", result.ItemCount)
		fmt.Printf("Fields imported: %d\n", result.FieldCount)
		fmt.Printf("Views imported: %d\n", result.ViewCount)
	}

	return nil
}
