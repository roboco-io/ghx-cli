package discussion

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// ViewOptions holds options for the view command
type ViewOptions struct {
	Repo         string
	Format       string
	Number       int
	CommentLimit int
}

// NewViewCmd creates the view command
func NewViewCmd() *cobra.Command {
	opts := &ViewOptions{}

	cmd := &cobra.Command{
		Use:   "view <owner/repo> <number>",
		Short: "View a discussion",
		Long: `View the details of a specific discussion by its number.

This command displays the full discussion content including title, body,
category, state, and optionally comments.`,
		Example: `  ghp discussion view owner/repo 123              # View discussion #123
  ghp discussion view owner/repo 123 --comments 10 # Show 10 comments
  ghp discussion view owner/repo 123 --format json # Output as JSON`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			number, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid discussion number: %s", args[1])
			}
			opts.Number = number
			return runView(cmd.Context(), opts)
		},
	}

	cmd.Flags().IntVar(&opts.CommentLimit, "comments", defaultCommentLimit, "Number of comments to show")
	cmd.Flags().StringVar(&opts.Format, "format", formatDetails, "Output format: details, json")

	return cmd
}

func runView(ctx context.Context, opts *ViewOptions) error {
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

	// Get discussion
	discussion, err := discussionService.GetDiscussion(ctx, owner, repo, opts.Number, opts.CommentLimit)
	if err != nil {
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	return outputDiscussionDetails(discussion, opts.Format)
}

func outputDiscussionDetails(d *service.DiscussionDetails, format string) error {
	switch format {
	case formatJSON:
		return outputDiscussionDetailsJSON(d)
	case formatDetails, formatTable:
		return outputDiscussionDetailsText(d)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputDiscussionDetailsText(d *service.DiscussionDetails) error {
	// Header
	fmt.Printf("Discussion #%d: %s\n", d.Number, d.Title)
	fmt.Println(strings.Repeat("=", viewSeparatorWidth))

	// Metadata
	fmt.Printf("State:      %s", d.State)
	if d.Locked {
		fmt.Print(" (Locked)")
	}
	fmt.Println()
	fmt.Printf("Category:   %s %s\n", d.CategoryInfo.Emoji, d.CategoryInfo.Name)
	fmt.Printf("Author:     %s\n", d.Author)
	fmt.Printf("Created:    %s\n", d.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:    %s\n", d.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Comments:   %d\n", d.CommentCount)
	fmt.Printf("Upvotes:    %d\n", d.UpvoteCount)
	fmt.Printf("URL:        %s\n", d.URL)

	if len(d.Labels) > 0 {
		fmt.Printf("Labels:     %s\n", strings.Join(d.Labels, ", "))
	}

	if d.HasAnswer {
		fmt.Println("Answered:   Yes")
	}

	// Body
	fmt.Println()
	fmt.Println("Body:")
	fmt.Println(strings.Repeat("-", viewSeparatorWidth))
	fmt.Println(d.Body)
	fmt.Println(strings.Repeat("-", viewSeparatorWidth))

	// Answer
	if d.Answer != nil {
		fmt.Println()
		fmt.Println("ANSWER:")
		fmt.Println(strings.Repeat("-", viewSeparatorWidth))
		fmt.Printf("By %s on %s:\n", d.Answer.Author, d.Answer.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println(d.Answer.Body)
		fmt.Println(strings.Repeat("-", viewSeparatorWidth))
	}

	// Comments
	if len(d.Comments) > 0 {
		fmt.Println()
		fmt.Printf("Comments (%d):\n", len(d.Comments))
		fmt.Println(strings.Repeat("-", viewSeparatorWidth))

		for i, c := range d.Comments {
			if c.IsAnswer {
				continue // Skip the answer, already shown above
			}
			fmt.Printf("\n[%d] %s - %s", i+1, c.Author, c.CreatedAt.Format("2006-01-02 15:04:05"))
			if c.UpvoteCount > 0 {
				fmt.Printf(" (+%d)", c.UpvoteCount)
			}
			fmt.Println()
			fmt.Println(c.Body)
		}
	}

	return nil
}

func outputDiscussionDetailsJSON(d *service.DiscussionDetails) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
