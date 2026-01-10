package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// LinkOptions holds options for the link command
type LinkOptions struct {
	Repository string
	Format     string
}

// NewLinkCmd creates the link command
func NewLinkCmd() *cobra.Command {
	opts := &LinkOptions{}

	cmd := &cobra.Command{
		Use:   "link PROJECT_ID",
		Short: "Link a project to a repository",
		Long: `Link an existing project to a GitHub repository.

This allows the project to automatically track issues and pull requests from the repository.

Examples:
  ghx project link myorg/123 --repo owner/repo    # Link project to repository
  ghx project link user/456 --repo myuser/myrepo  # Link to personal repository`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLink(cmd.Context(), opts, args)
		},
	}

	cmd.Flags().StringVar(&opts.Repository, "repo", "", "Repository to link (owner/repo)")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json, yaml")

	_ = cmd.MarkFlagRequired("repo")

	return cmd
}

func runLink(ctx context.Context, opts *LinkOptions, args []string) error {
	projectID := args[0]

	if opts.Repository == "" {
		return fmt.Errorf("repository is required")
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

	// Link project to repository
	err = projectService.LinkProjectToRepository(ctx, projectID, opts.Repository)
	if err != nil {
		return fmt.Errorf("failed to link project to repository: %w", err)
	}

	fmt.Printf("✅ Successfully linked project %s to repository %s\n", projectID, opts.Repository)
	return nil
}
