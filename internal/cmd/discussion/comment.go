package discussion

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// CommentOptions holds options for the comment command
type CommentOptions struct {
	ReplyToID *string
	Repo      string
	Body      string
	Format    string
	Number    int
}

// NewCommentCmd creates the comment command
func NewCommentCmd() *cobra.Command {
	opts := &CommentOptions{}
	var replyTo string

	cmd := &cobra.Command{
		Use:   "comment <owner/repo> <number>",
		Short: "Add a comment to a discussion",
		Long: `Add a comment to an existing discussion.

You can reply to a specific comment by using the --reply-to flag with the comment ID.`,
		Example: `  ghp discussion comment owner/repo 123 --body "This is my comment"
  ghp discussion comment owner/repo 123 -b "Reply to comment" --reply-to DC_xxx`,
		Args: cobra.ExactArgs(2),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if replyTo != "" {
				opts.ReplyToID = &replyTo
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number
			return runComment(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Comment body (required)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "Comment ID to reply to")
	cmd.Flags().StringVar(&opts.Format, "format", formatDetails, "Output format: details, json")

	_ = cmd.MarkFlagRequired("body")

	return cmd
}

func runComment(ctx context.Context, opts *CommentOptions) error {
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

	// Add comment
	commentOpts := service.AddCommentOptions{
		Owner:     owner,
		Repo:      repo,
		Number:    opts.Number,
		Body:      opts.Body,
		ReplyToID: opts.ReplyToID,
	}

	comment, err := discussionService.AddComment(ctx, commentOpts)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	// Output result
	switch opts.Format {
	case formatJSON:
		data, jsonErr := json.MarshalIndent(comment, "", "  ")
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal JSON: %w", jsonErr)
		}
		fmt.Println(string(data))
	default:
		fmt.Printf("Added comment to discussion #%d\n", opts.Number)
		fmt.Printf("Comment ID: %s\n", comment.ID)
		fmt.Printf("Author: %s\n", comment.Author)
		fmt.Printf("Created: %s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}
