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

// AnswerOptions holds options for the answer command
type AnswerOptions struct {
	Repo      string
	CommentID string
	Number    int
	Unmark    bool
}

// NewAnswerCmd creates the answer command
func NewAnswerCmd() *cobra.Command {
	opts := &AnswerOptions{}

	cmd := &cobra.Command{
		Use:   "answer <owner/repo> <number>",
		Short: "Mark or unmark a comment as the answer",
		Long: `Mark or unmark a comment as the answer to a discussion.

This only works for discussions in answerable categories (like Q&A).
Use --unmark to remove the answer designation from a comment.`,
		Example: `  ghp discussion answer owner/repo 123 --comment-id DC_xxx
  ghp discussion answer owner/repo 123 --comment-id DC_xxx --unmark`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number
			return runAnswer(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.CommentID, "comment-id", "", "Comment ID to mark as answer (required)")
	cmd.Flags().BoolVar(&opts.Unmark, "unmark", false, "Unmark the comment as answer")

	_ = cmd.MarkFlagRequired("comment-id")

	return cmd
}

func runAnswer(ctx context.Context, opts *AnswerOptions) error {
	// Parse repository reference (for validation)
	_, _, err := service.ParseRepositoryReference(opts.Repo)
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

	if opts.Unmark {
		// Unmark answer
		err = discussionService.UnmarkAnswer(ctx, opts.CommentID)
		if err != nil {
			return fmt.Errorf("failed to unmark answer: %w", err)
		}
		fmt.Printf("Unmarked comment %s as answer for discussion #%d\n", opts.CommentID, opts.Number)
	} else {
		// Mark answer
		err = discussionService.MarkAnswer(ctx, opts.CommentID)
		if err != nil {
			return fmt.Errorf("failed to mark answer: %w", err)
		}
		fmt.Printf("Marked comment %s as answer for discussion #%d\n", opts.CommentID, opts.Number)
	}

	return nil
}
