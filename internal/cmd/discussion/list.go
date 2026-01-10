package discussion

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// ListOptions holds options for the list command
type ListOptions struct {
	Answered *bool
	Repo     string
	Category string
	State    string
	Format   string
	Limit    int
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	opts := &ListOptions{}
	var answered, unanswered bool

	cmd := &cobra.Command{
		Use:   "list <owner/repo>",
		Short: "List discussions in a repository",
		Long: `List discussions for a repository with optional filters.

You can filter discussions by category, state, and whether they have been answered.
The --answered and --unanswered flags are mutually exclusive.`,
		Example: `  ghp discussion list owner/repo                    # List all discussions
  ghp discussion list owner/repo --category ideas   # Filter by category
  ghp discussion list owner/repo --state open       # Filter by state
  ghp discussion list owner/repo --answered         # Show only answered
  ghp discussion list owner/repo --unanswered       # Show only unanswered
  ghp discussion list owner/repo --limit 50         # Limit results`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if answered && unanswered {
				return fmt.Errorf("--answered and --unanswered are mutually exclusive")
			}
			if answered {
				t := true
				opts.Answered = &t
			} else if unanswered {
				f := false
				opts.Answered = &f
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			return runList(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Category, "category", "", "Filter by category slug")
	cmd.Flags().StringVar(&opts.State, "state", stateAll, "Filter by state: open, closed, all")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", defaultListLimit, "Maximum number of discussions")
	cmd.Flags().StringVar(&opts.Format, "format", formatTable, "Output format: table, json")
	cmd.Flags().BoolVar(&answered, "answered", false, "Show only answered discussions")
	cmd.Flags().BoolVar(&unanswered, "unanswered", false, "Show only unanswered discussions")

	return cmd
}

func runList(ctx context.Context, opts *ListOptions) error {
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

	// List discussions
	listOpts := service.ListDiscussionsOptions{
		Owner:    owner,
		Repo:     repo,
		Category: opts.Category,
		Answered: opts.Answered,
		State:    opts.State,
		First:    opts.Limit,
	}

	discussions, err := discussionService.ListDiscussions(ctx, listOpts)
	if err != nil {
		return fmt.Errorf("failed to list discussions: %w", err)
	}

	return outputDiscussions(discussions, opts.Format)
}

func outputDiscussions(discussions []service.DiscussionInfo, format string) error {
	if len(discussions) == 0 {
		fmt.Println("No discussions found")
		return nil
	}

	switch format {
	case formatJSON:
		return outputDiscussionsJSON(discussions)
	case formatTable:
		return outputDiscussionsTable(discussions)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputDiscussionsTable(discussions []service.DiscussionInfo) error {
	fmt.Printf("%-6s %-42s %-15s %-8s %-15s %-8s %-8s\n",
		"NUM", "TITLE", "CATEGORY", "STATE", "AUTHOR", "COMMENTS", "ANSWERED")
	fmt.Println(strings.Repeat("-", tableSeparatorWidth))

	for _, d := range discussions {
		title := truncateString(d.Title, titleMaxLength, titleTruncateLength)
		category := truncateString(d.Category, categoryMaxLength, categoryTruncateLength)
		author := truncateString(d.Author, authorMaxLength, authorTruncateLength)

		answered := "No"
		if d.HasAnswer {
			answered = "Yes"
		}

		fmt.Printf("%-6d %-42s %-15s %-8s %-15s %-8d %-8s\n",
			d.Number, title, category, d.State, author, d.CommentCount, answered)
	}

	return nil
}

func outputDiscussionsJSON(discussions []service.DiscussionInfo) error {
	data, err := json.MarshalIndent(discussions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func truncateString(s string, maxLen, truncateLen int) string {
	if len(s) > maxLen {
		return s[:truncateLen] + "..."
	}
	return s
}
