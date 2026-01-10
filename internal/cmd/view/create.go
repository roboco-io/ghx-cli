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

// CreateOptions holds options for the create command
type CreateOptions struct {
	ProjectRef string
	Name       string
	Layout     string
	Filter     string
	Format     string
}

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <owner/project-number> <name> <layout>",
		Short: "Create a new project view",
		Long: `Create a new view in a GitHub Project.

Views provide different ways to visualize and organize your project data.
You can create table, board, or roadmap views depending on your needs.

View Layouts:
  table       - Table view with customizable columns and sorting
  board       - Kanban board view with swimlanes and cards  
  roadmap     - Timeline roadmap for milestone planning

Examples:
  ghx view create octocat/123 "Sprint Dashboard" table
  ghx view create octocat/123 "Bug Board" board --filter "label:bug"
  ghx view create --org myorg/456 "Release Roadmap" roadmap
  ghx view create octocat/123 "High Priority" table --filter "priority:high" --format json`,

		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Name = args[1]
			opts.Layout = args[2]
			opts.Format = cmd.Flag("format").Value.String()
			return runCreate(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter expression for the view")
	cmd.Flags().Bool("org", false, "Create view in organization project")

	return cmd
}

func runCreate(ctx context.Context, opts *CreateOptions) error {
	// Validate view name
	if err := service.ValidateViewName(opts.Name); err != nil {
		return err
	}

	// Validate and normalize layout
	layout, err := service.ValidateViewLayout(opts.Layout)
	if err != nil {
		return err
	}

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

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)
	viewService := service.NewViewService(client)

	// Get project to validate access and get project ID (with automatic owner detection)
	project, err := projectService.GetProjectWithOwnerDetection(ctx, owner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Create view
	input := service.CreateViewInput{
		ProjectID: project.ID,
		Name:      opts.Name,
		Layout:    layout,
	}

	view, err := viewService.CreateView(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create view: %w", err)
	}

	// Update view with filter if provided
	if opts.Filter != "" {
		updateInput := service.UpdateViewInput{
			ViewID: view.ID,
			Filter: &opts.Filter,
		}

		updatedView, err := viewService.UpdateView(ctx, updateInput)
		if err != nil {
			// Log warning but don't fail the entire operation
			fmt.Printf("Warning: View created but failed to set filter: %v\n", err)
		} else {
			view = updatedView
		}
	}

	// Output created view
	return outputCreatedView(view, opts.Format)
}

func outputCreatedView(view *graphql.ProjectV2View, format string) error {
	switch format {
	case formatJSON:
		return outputCreatedViewJSON(view)
	case formatTable:
		return outputCreatedViewTable(view)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputCreatedViewTable(view *graphql.ProjectV2View) error {
	fmt.Printf("✅ View '%s' created successfully\n\n", view.Name)
	fmt.Printf("View Details:\n")
	outputViewDetailsTable(view)
	return nil
}

func outputCreatedViewJSON(view *graphql.ProjectV2View) error {
	fmt.Printf("{\n")
	outputViewDetailsJSON(view)
	fmt.Printf("}\n")
	return nil
}
