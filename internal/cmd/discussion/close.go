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

// CloseOptions holds options for the close command
type CloseOptions struct {
	Repo   string
	Reason string
	Format string
	Number int
}

// NewCloseCmd creates the close command
func NewCloseCmd() *cobra.Command {
	opts := &CloseOptions{}

	cmd := &cobra.Command{
		Use:   "close <owner/repo> <number>",
		Short: "Close a discussion",
		Long: `Close a discussion with an optional reason.

Valid reasons:
  - resolved (default): The discussion has been resolved
  - outdated: The discussion is no longer relevant
  - duplicate: The discussion is a duplicate of another`,
		Example: `  ghp discussion close owner/repo 123
  ghp discussion close owner/repo 123 --reason resolved
  ghp discussion close owner/repo 123 --reason outdated`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number

			// Validate reason
			if err := service.ValidateCloseReason(opts.Reason); err != nil {
				return err
			}

			return runClose(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Reason, "reason", closeReasonResolved, "Close reason: resolved, outdated, duplicate")
	cmd.Flags().StringVar(&opts.Format, "format", formatDetails, "Output format: details, json")

	return cmd
}

func runClose(ctx context.Context, opts *CloseOptions) error {
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

	// Close discussion
	closeOpts := service.CloseDiscussionOptions{
		Owner:  owner,
		Repo:   repo,
		Number: opts.Number,
		Reason: opts.Reason,
	}

	discussion, err := discussionService.CloseDiscussion(ctx, closeOpts)
	if err != nil {
		return fmt.Errorf("failed to close discussion: %w", err)
	}

	// Output result
	switch opts.Format {
	case formatJSON:
		data, jsonErr := json.MarshalIndent(discussion, "", "  ")
		if jsonErr != nil {
			return fmt.Errorf("failed to marshal JSON: %w", jsonErr)
		}
		fmt.Println(string(data))
	default:
		fmt.Printf("Closed discussion #%d: %s (reason: %s)\n", discussion.Number, discussion.Title, opts.Reason)
	}

	return nil
}
