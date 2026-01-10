package project

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

// EditOptions holds options for the edit command
type EditOptions struct {
	Owner  string
	Title  string
	Format string
	Number int
	Org    bool
	Close  bool
	Reopen bool
}

// NewEditCmd creates the edit command
func NewEditCmd() *cobra.Command {
	opts := &EditOptions{}

	cmd := &cobra.Command{
		Use:   "edit {<owner>/<number> | <number>}",
		Short: "Edit a project",
		Long: `Edit an existing project.

Examples:
  ghx project edit 123 --title "New Title"      # Edit project title
  ghx project edit octocat/123 --close          # Close project
  ghx project edit myorg/456 --reopen --org     # Reopen org project`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(cmd.Context(), opts, args)
		},
	}

	cmd.Flags().BoolVar(&opts.Org, "org", false, "Project belongs to an organization")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "New project title")
	cmd.Flags().BoolVar(&opts.Close, "close", false, "Close the project")
	cmd.Flags().BoolVar(&opts.Reopen, "reopen", false, "Reopen the project")
	cmd.Flags().StringVar(&opts.Format, "format", "details", "Output format: details, json")

	return cmd
}

func runEdit(ctx context.Context, opts *EditOptions, args []string) error {
	// Parse project reference
	projectRef := args[0]
	var err error

	if strings.Contains(projectRef, "/") {
		opts.Owner, opts.Number, err = service.ParseProjectReference(projectRef)
		if err != nil {
			return fmt.Errorf("invalid project reference: %w", err)
		}
	} else {
		// Just a number, need to determine owner from context
		opts.Number, err = strconv.Atoi(projectRef)
		if err != nil {
			return fmt.Errorf("invalid project number: %s", projectRef)
		}

		// For now, require owner to be specified
		return fmt.Errorf("owner must be specified in format owner/number")
	}

	// Validate conflicting flags
	if opts.Close && opts.Reopen {
		return fmt.Errorf("cannot specify both --close and --reopen")
	}

	// Check if any changes are specified
	if opts.Title == "" && !opts.Close && !opts.Reopen {
		return fmt.Errorf("no changes specified (use --title, --close, or --reopen)")
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

	// First, get the current project to obtain its ID
	currentProject, err := projectService.GetProject(ctx, opts.Owner, opts.Number, opts.Org)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Prepare update input
	updateInput := service.UpdateProjectInput{
		ProjectID: currentProject.ID,
	}

	if opts.Title != "" {
		updateInput.Title = &opts.Title
	}

	if opts.Close {
		closed := true
		updateInput.Closed = &closed
	} else if opts.Reopen {
		closed := false
		updateInput.Closed = &closed
	}

	// Update project
	updatedProject, err := projectService.UpdateProject(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	// Output updated project
	fmt.Printf("✅ Project updated successfully!\n\n")
	return outputUpdatedProject(updatedProject, opts.Format)
}

func outputUpdatedProject(project *graphql.ProjectV2, format string) error {
	switch format {
	case formatJSON:
		return outputProjectDetailsJSON(project)
	case "details":
		fmt.Printf("Project #%d\n", project.Number)
		fmt.Printf("Title: %s\n", project.Title)

		if project.Description != nil {
			fmt.Printf("Description: %s\n", *project.Description)
		}

		fmt.Printf("URL: %s\n", project.URL)
		fmt.Printf("Owner: %s (%s)\n", project.Owner.Login, project.Owner.Type)

		state := "Open"
		if project.Closed {
			state = "Closed"
		}
		fmt.Printf("State: %s\n", state)

		fmt.Printf("Updated: %s\n", project.UpdatedAt.Format("2006-01-02 15:04:05"))

		return nil
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}
