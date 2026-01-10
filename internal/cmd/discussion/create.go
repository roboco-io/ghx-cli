package discussion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// CreateOptions holds options for the create command
type CreateOptions struct {
	Repo     string
	Category string
	Title    string
	Body     string
	Format   string
}

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create <owner/repo>",
		Short: "Create a new discussion",
		Long: `Create a new discussion in a repository.

You must specify a category (by slug) for the discussion.
Use 'ghp discussion category list' to see available categories.`,
		Example: `  ghx discussion create owner/repo --category ideas --title "New feature" --body "Description"
  ghx discussion create owner/repo -c general -t "Question" -b "How do I...?"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			return runCreate(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Category, "category", "c", "", "Category slug (required)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Discussion title (required)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Discussion body (required)")
	cmd.Flags().StringVar(&opts.Format, "format", formatDetails, "Output format: details, json")

	_ = cmd.MarkFlagRequired("category")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("body")

	return cmd
}

func runCreate(ctx context.Context, opts *CreateOptions) error {
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

	// Create discussion
	createOpts := service.CreateDiscussionOptions{
		Owner:    owner,
		Repo:     repo,
		Category: opts.Category,
		Title:    opts.Title,
		Body:     opts.Body,
	}

	discussion, err := discussionService.CreateDiscussion(ctx, createOpts)
	if err != nil {
		return fmt.Errorf("failed to create discussion: %w", err)
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
		fmt.Printf("Created discussion #%d: %s\n", discussion.Number, discussion.Title)
		fmt.Printf("URL: %s\n", discussion.URL)
	}

	return nil
}
