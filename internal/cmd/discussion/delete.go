package discussion

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// DeleteOptions holds options for the delete command
type DeleteOptions struct {
	Repo   string
	Number int
	Force  bool
}

// NewDeleteCmd creates the delete command
func NewDeleteCmd() *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <owner/repo> <number>",
		Short: "Delete a discussion",
		Long: `Delete a discussion from a repository.

This action is irreversible. All comments and replies will also be deleted.
You will be prompted to confirm unless --force is used.`,
		Example: `  ghp discussion delete owner/repo 123
  ghp discussion delete owner/repo 123 --force`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number
			return runDelete(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *DeleteOptions) error {
	// Parse repository reference
	owner, repo, err := service.ParseRepositoryReference(opts.Repo)
	if err != nil {
		return err
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	discussionService := service.NewDiscussionService(client)

	// Get discussion to confirm
	discussion, err := discussionService.GetDiscussion(ctx, owner, repo, opts.Number, 0)
	if err != nil {
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	// Confirm deletion
	if !opts.Force {
		fmt.Printf("You are about to delete the following discussion:\n\n")
		fmt.Printf("  #%d: %s\n", discussion.Number, discussion.Title)
		fmt.Printf("  Category: %s\n", discussion.Category)
		fmt.Printf("  Comments: %d\n", discussion.CommentCount)
		fmt.Printf("\n")
		fmt.Printf("This action cannot be undone. Type 'DELETE' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		confirmation, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("failed to read confirmation: %w", readErr)
		}

		if strings.TrimSpace(confirmation) != "DELETE" {
			fmt.Println("Deletion canceled")
			return nil
		}
	}

	// Delete discussion
	err = discussionService.DeleteDiscussion(ctx, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to delete discussion: %w", err)
	}

	fmt.Printf("Deleted discussion #%d: %s\n", discussion.Number, discussion.Title)
	return nil
}
