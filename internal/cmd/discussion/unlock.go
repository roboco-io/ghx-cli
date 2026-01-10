package discussion

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// UnlockOptions holds options for the unlock command
type UnlockOptions struct {
	Repo   string
	Number int
}

// NewUnlockCmd creates the unlock command
func NewUnlockCmd() *cobra.Command {
	opts := &UnlockOptions{}

	cmd := &cobra.Command{
		Use:     "unlock <owner/repo> <number>",
		Short:   "Unlock a discussion",
		Long:    `Unlock a previously locked discussion to allow new comments.`,
		Example: `  ghp discussion unlock owner/repo 123`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number
			return runUnlock(cmd.Context(), opts)
		},
	}

	return cmd
}

func runUnlock(ctx context.Context, opts *UnlockOptions) error {
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

	// Unlock discussion
	err = discussionService.UnlockDiscussion(ctx, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to unlock discussion: %w", err)
	}

	fmt.Printf("Unlocked discussion #%d\n", opts.Number)
	return nil
}
