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

// LockOptions holds options for the lock command
type LockOptions struct {
	Repo   string
	Reason string
	Number int
}

// NewLockCmd creates the lock command
func NewLockCmd() *cobra.Command {
	opts := &LockOptions{}

	cmd := &cobra.Command{
		Use:   "lock <owner/repo> <number>",
		Short: "Lock a discussion",
		Long: `Lock a discussion to prevent new comments.

Valid reasons:
  - off_topic: The discussion is off-topic
  - resolved: The discussion has been resolved
  - spam: The discussion contains spam
  - too_heated: The discussion has become too heated`,
		Example: `  ghp discussion lock owner/repo 123
  ghp discussion lock owner/repo 123 --reason resolved
  ghp discussion lock owner/repo 123 --reason spam`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number

			// Validate reason if provided
			if opts.Reason != "" {
				if err := service.ValidateLockReason(opts.Reason); err != nil {
					return err
				}
			}

			return runLock(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Reason, "reason", "", "Lock reason: off_topic, resolved, spam, too_heated")

	return cmd
}

func runLock(ctx context.Context, opts *LockOptions) error {
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

	// Lock discussion
	lockOpts := service.LockDiscussionOptions{
		Owner:  owner,
		Repo:   repo,
		Number: opts.Number,
		Reason: opts.Reason,
	}

	err = discussionService.LockDiscussion(ctx, lockOpts)
	if err != nil {
		return fmt.Errorf("failed to lock discussion: %w", err)
	}

	if opts.Reason != "" {
		fmt.Printf("Locked discussion #%d (reason: %s)\n", opts.Number, opts.Reason)
	} else {
		fmt.Printf("Locked discussion #%d\n", opts.Number)
	}

	return nil
}
