package view

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// CopyOptions holds options for the copy command
type CopyOptions struct {
	ViewID     string
	ProjectRef string
	Name       string
	Format     string
}

// NewCopyCmd creates the copy command
func NewCopyCmd() *cobra.Command {
	opts := &CopyOptions{}

	cmd := &cobra.Command{
		Use:   "copy <view-id> <new-name> [project-ref]",
		Short: "Copy a project view",
		Long: `Create a copy of an existing project view.

The copied view will have the same layout, filter, sorting, and grouping
configuration as the original view. You can optionally copy the view to
a different project by specifying the target project reference.

Examples:
  ghx view copy view-id "Sprint 2 Board"
  ghx view copy view-id "Bug Dashboard" octocat/456
  ghx view copy view-id "Roadmap Copy" --format json`,

		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ViewID = args[0]
			opts.Name = args[1]
			if len(args) > 2 {
				opts.ProjectRef = args[2]
			}
			opts.Format = cmd.Flag("format").Value.String()
			return runCopy(cmd.Context(), opts)
		},
	}

	cmd.Flags().Bool("org", false, "Copy to organization project")

	return cmd
}

func runCopy(ctx context.Context, opts *CopyOptions) error {
	// Validate view name
	if err := service.ValidateViewName(opts.Name); err != nil {
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
	viewService := service.NewViewService(client)

	var projectID string

	if opts.ProjectRef != "" {
		// Copy to different project
		projectService := service.NewProjectService(client)

		// Parse project reference
		parts := strings.Split(opts.ProjectRef, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid project reference format. Use: owner/project-number")
		}

		owner := parts[0]
		projectNumber, parseErr := strconv.Atoi(parts[1])
		if parseErr != nil {
			return fmt.Errorf("invalid project number: %s", parts[1])
		}

		// Get target project (with automatic owner detection)
		project, getErr := projectService.GetProjectWithOwnerDetection(ctx, owner, projectNumber)
		if getErr != nil {
			return fmt.Errorf("failed to get target project: %w", getErr)
		}

		projectID = project.ID
	} else {
		// Copy within same project - get source view to determine project
		sourceView, viewErr := viewService.GetView(ctx, opts.ViewID)
		if viewErr != nil {
			return fmt.Errorf("failed to get source view: %w", viewErr)
		}
		projectID = sourceView.ProjectID
	}

	// Copy view
	input := service.CopyViewInput{
		ProjectID: projectID,
		ViewID:    opts.ViewID,
		Name:      opts.Name,
	}

	view, err := viewService.CopyView(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to copy view: %w", err)
	}

	// Output copied view
	return outputCopiedView(view, opts.Format)
}

func outputCopiedView(view *graphql.ProjectV2View, format string) error {
	switch format {
	case formatJSON:
		return outputCopiedViewJSON(view)
	case formatTable:
		return outputCopiedViewTable(view)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputCopiedViewTable(view *graphql.ProjectV2View) error {
	fmt.Printf("✅ View '%s' copied successfully\n\n", view.Name)
	fmt.Printf("New View Details:\n")
	outputViewDetailsTable(view)
	return nil
}

func outputCopiedViewJSON(view *graphql.ProjectV2View) error {
	fmt.Printf("{\n")
	outputViewDetailsJSON(view)
	fmt.Printf("}\n")
	return nil
}
