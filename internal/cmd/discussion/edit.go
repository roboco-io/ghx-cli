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

// EditOptions holds options for the edit command
type EditOptions struct {
	Category *string
	Title    *string
	Body     *string
	Repo     string
	Format   string
	Number   int
}

// NewEditCmd creates the edit command
func NewEditCmd() *cobra.Command {
	opts := &EditOptions{}
	var title, body, category string

	cmd := &cobra.Command{
		Use:   "edit <owner/repo> <number>",
		Short: "Edit a discussion",
		Long: `Edit an existing discussion's title, body, or category.

At least one of --title, --body, or --category must be specified.`,
		Example: `  ghx discussion edit owner/repo 123 --title "New title"
  ghx discussion edit owner/repo 123 --body "Updated description"
  ghx discussion edit owner/repo 123 --category ideas`,
		Args: cobra.ExactArgs(2),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if title == "" && body == "" && category == "" {
				return fmt.Errorf("at least one of --title, --body, or --category must be specified")
			}
			if title != "" {
				opts.Title = &title
			}
			if body != "" {
				opts.Body = &body
			}
			if category != "" {
				opts.Category = &category
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
			return runEdit(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "New discussion title")
	cmd.Flags().StringVarP(&body, "body", "b", "", "New discussion body")
	cmd.Flags().StringVarP(&category, "category", "c", "", "New category slug")
	cmd.Flags().StringVar(&opts.Format, "format", formatDetails, "Output format: details, json")

	return cmd
}

func runEdit(ctx context.Context, opts *EditOptions) error {
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

	// Update discussion
	updateOpts := service.UpdateDiscussionOptions{
		Owner:    owner,
		Repo:     repo,
		Number:   opts.Number,
		Title:    opts.Title,
		Body:     opts.Body,
		Category: opts.Category,
	}

	discussion, err := discussionService.UpdateDiscussion(ctx, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update discussion: %w", err)
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
		fmt.Printf("Updated discussion #%d: %s\n", discussion.Number, discussion.Title)
		fmt.Printf("URL: %s\n", discussion.URL)
	}

	return nil
}
