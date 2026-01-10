package item

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// RemoveOptions holds options for the remove command
type RemoveOptions struct {
	ProjectRef string
	ItemID     string
	Force      bool
}

// NewRemoveCmd creates the remove command
func NewRemoveCmd() *cobra.Command {
	opts := &RemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove <project> <item-id>",
		Short: "Remove an item from a project",
		Long: `Remove an item from a project by its project item ID.

The item ID should be the project-specific item ID, not the issue or PR ID.
You can find item IDs by listing project items or using the GitHub web interface.

⚠️  WARNING: This action cannot be undone. The item will be removed from the project
but the underlying issue or PR will remain unchanged.

Examples:
  ghx item remove octocat/1 PVTI_lADOANN5s84ACbL0zgBZrOY    # Remove item from project
  ghx item remove myorg/2 item-123 --force                  # Skip confirmation`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.ItemID = args[1]
			return runRemove(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runRemove(ctx context.Context, opts *RemoveOptions) error {
	// Parse project reference
	projectOwner, projectNumber, err := service.ParseProjectReference(opts.ProjectRef)
	if err != nil {
		return fmt.Errorf("invalid project reference: %w", err)
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	itemService := service.NewItemService(client)
	projectService := service.NewProjectService(client)

	// Get project details
	project, err := projectService.GetProjectWithOwnerDetection(ctx, projectOwner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Show confirmation unless --force is used
	if !opts.Force {
		fmt.Printf("⚠️  You are about to remove item %s from project:\n\n", opts.ItemID)
		fmt.Printf("Project: %s (#%d)\n", project.Title, project.Number)
		fmt.Printf("Owner: %s\n", project.Owner.Login)
		fmt.Printf("\n⚠️  This action cannot be undone. The item will be removed from the project.\n")
		fmt.Printf("Type 'REMOVE' to confirm: ")

		var confirmation string
		_, scanErr := fmt.Scanln(&confirmation)
		if scanErr != nil {
			fmt.Println("❌ Failed to read confirmation.")
			return scanErr
		}

		if confirmation != "REMOVE" {
			fmt.Println("❌ Removal canceled.")
			return nil
		}
	}

	// Remove item from project
	err = itemService.RemoveItemFromProject(ctx, project.ID, opts.ItemID)
	if err != nil {
		return fmt.Errorf("failed to remove item from project: %w", err)
	}

	fmt.Printf("✅ Item %s removed from project successfully.\n", opts.ItemID)
	return nil
}
