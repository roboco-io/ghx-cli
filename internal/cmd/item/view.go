package item

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

const (
	maxBodyDisplayLength = 500
)

// ViewOptions holds options for the view command
type ViewOptions struct {
	ItemRef string
	Format  string
	Web     bool
}

// NewViewCmd creates the view command
func NewViewCmd() *cobra.Command {
	opts := &ViewOptions{}

	cmd := &cobra.Command{
		Use:   "view <item>",
		Short: "View details of an issue or pull request",
		Long: `View detailed information about a specific issue or pull request.

Item references can be in the following formats:
• owner/repo#123 (issue or PR reference)
• https://github.com/owner/repo/issues/123 (GitHub issue URL)
• https://github.com/owner/repo/pull/456 (GitHub PR URL)

Examples:
  ghx item view octocat/Hello-World#123              # View issue details
  ghx item view https://github.com/cli/cli/pull/456  # View PR from URL
  ghx item view myorg/repo#789 --format json         # View in JSON format
  ghx item view octocat/Hello-World#123 --web        # Open in browser`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ItemRef = args[0]
			return runView(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "details", "Output format: details, json")
	cmd.Flags().BoolVar(&opts.Web, "web", false, "Open item in web browser")

	return cmd
}

func runView(ctx context.Context, opts *ViewOptions) error {
	// Parse item reference
	owner, repo, number, err := service.ParseItemReference(opts.ItemRef)
	if err != nil {
		return fmt.Errorf("invalid item reference: %w", err)
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	itemService := service.NewItemService(client)

	// Try to get as issue first
	issue, err := itemService.GetIssue(ctx, owner, repo, number)
	if err == nil {
		if opts.Web {
			fmt.Printf("Opening issue in browser: %s\n", issue.URL)
			return nil
		}
		return outputIssueDetails(issue, opts.Format)
	}

	// Try as pull request
	pr, err := itemService.GetPullRequest(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to find issue or pull request: %w", err)
	}

	if opts.Web {
		fmt.Printf("Opening pull request in browser: %s\n", pr.URL)
		return nil
	}

	return outputPullRequestDetails(pr, opts.Format)
}

func outputIssueDetails(issue *graphql.Issue, format string) error {
	switch format {
	case "json":
		return outputIssueDetailsJSON(issue)
	case "details":
		return outputIssueDetailsTable(issue)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputPullRequestDetails(pr *graphql.PullRequest, format string) error {
	switch format {
	case "json":
		return outputPullRequestDetailsJSON(pr)
	case "details":
		return outputPullRequestDetailsTable(pr)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputIssueDetailsTable(issue *graphql.Issue) error {
	fmt.Printf("Issue #%d\n", issue.Number)
	fmt.Printf("Title: %s\n", issue.Title)
	fmt.Printf("Repository: %s\n", issue.Repository.NameWithOwner)
	fmt.Printf("Author: %s\n", issue.Author.Login)
	fmt.Printf("State: %s\n", issue.State)
	if issue.Closed {
		fmt.Printf("Status: Closed\n")
	} else {
		fmt.Printf("Status: Open\n")
	}
	fmt.Printf("URL: %s\n", issue.URL)
	fmt.Printf("Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Labels
	if len(issue.Labels.Nodes) > 0 {
		fmt.Printf("\nLabels:\n")
		for _, label := range issue.Labels.Nodes {
			fmt.Printf("  • %s\n", label.Name)
		}
	}

	// Assignees
	if len(issue.Assignees.Nodes) > 0 {
		fmt.Printf("\nAssignees:\n")
		for _, assignee := range issue.Assignees.Nodes {
			fmt.Printf("  • %s\n", assignee.Login)
		}
	}

	// Body
	if issue.Body != "" {
		fmt.Printf("\nDescription:\n")
		fmt.Printf("────────────\n")

		// Truncate long body for readability
		body := issue.Body
		if len(body) > maxBodyDisplayLength {
			body = body[:497] + "..."
		}

		// Simple formatting
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			fmt.Printf("%s\n", line)
		}
	}

	return nil
}

func outputPullRequestDetailsTable(pr *graphql.PullRequest) error {
	fmt.Printf("Pull Request #%d\n", pr.Number)
	fmt.Printf("Title: %s\n", pr.Title)
	fmt.Printf("Repository: %s\n", pr.Repository.NameWithOwner)
	fmt.Printf("Author: %s\n", pr.Author.Login)

	state := pr.State
	if pr.Merged {
		state = "MERGED"
	}
	fmt.Printf("State: %s\n", state)

	fmt.Printf("URL: %s\n", pr.URL)
	fmt.Printf("Created: %s\n", pr.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", pr.UpdatedAt.Format("2006-01-02 15:04:05"))

	// Labels
	if len(pr.Labels.Nodes) > 0 {
		fmt.Printf("\nLabels:\n")
		for _, label := range pr.Labels.Nodes {
			fmt.Printf("  • %s\n", label.Name)
		}
	}

	// Assignees
	if len(pr.Assignees.Nodes) > 0 {
		fmt.Printf("\nAssignees:\n")
		for _, assignee := range pr.Assignees.Nodes {
			fmt.Printf("  • %s\n", assignee.Login)
		}
	}

	// Review requests
	if len(pr.ReviewRequests.Nodes) > 0 {
		fmt.Printf("\nReview Requests:\n")
		for _, request := range pr.ReviewRequests.Nodes {
			fmt.Printf("  • %s\n", request.RequestedReviewer.User.Login)
		}
	}

	// Body
	if pr.Body != "" {
		fmt.Printf("\nDescription:\n")
		fmt.Printf("────────────\n")

		// Truncate long body for readability
		body := pr.Body
		if len(body) > maxBodyDisplayLength {
			body = body[:497] + "..."
		}

		// Simple formatting
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			fmt.Printf("%s\n", line)
		}
	}

	return nil
}

func outputIssueDetailsJSON(issue *graphql.Issue) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"type\": \"Issue\",\n")
	fmt.Printf("  \"number\": %d,\n", issue.Number)
	fmt.Printf("  \"title\": \"%s\",\n", issue.Title)
	fmt.Printf("  \"repository\": \"%s\",\n", issue.Repository.NameWithOwner)
	fmt.Printf("  \"author\": \"%s\",\n", issue.Author.Login)
	fmt.Printf("  \"state\": \"%s\",\n", issue.State)
	fmt.Printf("  \"closed\": %t,\n", issue.Closed)
	fmt.Printf("  \"url\": \"%s\",\n", issue.URL)
	fmt.Printf("  \"created_at\": \"%s\",\n", issue.CreatedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("  \"updated_at\": \"%s\"\n", issue.UpdatedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("}\n")
	return nil
}

func outputPullRequestDetailsJSON(pr *graphql.PullRequest) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"type\": \"PullRequest\",\n")
	fmt.Printf("  \"number\": %d,\n", pr.Number)
	fmt.Printf("  \"title\": \"%s\",\n", pr.Title)
	fmt.Printf("  \"repository\": \"%s\",\n", pr.Repository.NameWithOwner)
	fmt.Printf("  \"author\": \"%s\",\n", pr.Author.Login)
	fmt.Printf("  \"state\": \"%s\",\n", pr.State)
	fmt.Printf("  \"closed\": %t,\n", pr.Closed)
	fmt.Printf("  \"merged\": %t,\n", pr.Merged)
	fmt.Printf("  \"url\": \"%s\",\n", pr.URL)
	fmt.Printf("  \"created_at\": \"%s\",\n", pr.CreatedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("  \"updated_at\": \"%s\"\n", pr.UpdatedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("}\n")
	return nil
}
