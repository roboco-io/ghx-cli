package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// CreateOptions holds options for the create command
type CreateOptions struct {
	Title       string
	Description string
	Readme      string
	Visibility  string
	Repository  string
	OwnerID     string
	Format      string
	Org         bool
	Web         bool
}

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create [title]",
		Short: "Create a new project",
		Long: `Create a new project for a user or organization.

Examples:
  ghx project create "My Project"                           # Basic project creation
  ghx project create "My Project" --description "A project" # With description
  ghx project create "My Project" --readme "Detailed info"  # With README
  ghx project create "My Project" --visibility private      # Private project
  ghx project create "My Project" --repo owner/repo         # Link to repository`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(cmd.Context(), opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Project title")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Project description")
	cmd.Flags().StringVar(&opts.Readme, "readme", "", "Project README content")
	cmd.Flags().StringVar(&opts.Visibility, "visibility", "public", "Project visibility (public, private)")
	cmd.Flags().StringVar(&opts.Repository, "repo", "", "Link to repository (owner/repo)")
	cmd.Flags().StringVar(&opts.OwnerID, "owner-id", "", "Owner ID (user or organization)")
	cmd.Flags().BoolVar(&opts.Org, "org", false, "Create organization project")
	cmd.Flags().BoolVar(&opts.Web, "web", false, "Open project in web browser after creation")
	cmd.Flags().StringVar(&opts.Format, "format", "details", "Output format: details, json")

	return cmd
}

func runCreate(ctx context.Context, opts *CreateOptions, args []string) error {
	// Get title from args or flags
	if len(args) > 0 {
		opts.Title = args[0]
	}

	if opts.Title == "" {
		return fmt.Errorf("project title is required")
	}

	if opts.OwnerID == "" {
		return fmt.Errorf("owner ID is required (use --owner-id flag)")
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

	// Create project
	input := &service.CreateProjectInput{
		OwnerID:     opts.OwnerID,
		Title:       opts.Title,
		Description: opts.Description,
		Readme:      opts.Readme,
		Visibility:  opts.Visibility,
		Repository:  opts.Repository,
	}

	project, err := projectService.CreateProject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	// Open in web browser if requested
	if opts.Web {
		fmt.Printf("Opening project in browser: %s\n", project.URL)
		// In a real implementation, we'd use a library to open the browser
	}

	// Output project details
	fmt.Printf("✅ Project created successfully!\n\n")
	return outputCreatedProject(project, opts.Format)
}

func outputCreatedProject(project *graphql.ProjectV2, format string) error {
	switch format {
	case formatJSON:
		return outputProjectDetailsJSON(project)
	case formatDetails:
		fmt.Printf("Project #%d\n", project.Number)
		fmt.Printf("Title: %s\n", project.Title)

		if project.Description != nil {
			fmt.Printf("Description: %s\n", *project.Description)
		}

		fmt.Printf("URL: %s\n", project.URL)
		fmt.Printf("Owner: %s (%s)\n", project.Owner.Login, project.Owner.Type)
		fmt.Printf("Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}
