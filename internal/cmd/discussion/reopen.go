package discussion

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// ReopenOptions holds options for the reopen command
type ReopenOptions struct {
	Repo   string
	Format string
	Number int
}

// NewReopenCmd creates the reopen command
func NewReopenCmd() *cobra.Command {
	opts := &ReopenOptions{}

	cmd := &cobra.Command{
		Use:   "reopen <owner/repo> <number>",
		Short: "Reopen a closed discussion",
		Long:  `Reopen a previously closed discussion.`,
		Example: `  ghx discussion reopen owner/repo 123
  ghx discussion reopen owner/repo 123 --format json`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number
			return runReopen(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", formatDetails, "Output format: details, json")

	return cmd
}

func runReopen(ctx context.Context, opts *ReopenOptions) error {
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

	// Reopen discussion
	discussion, err := discussionService.ReopenDiscussion(ctx, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to reopen discussion: %w", err)
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
		fmt.Printf("Reopened discussion #%d: %s\n", discussion.Number, discussion.Title)
	}

	return nil
}
