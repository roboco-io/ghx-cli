package item

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

const (
	defaultListLimit          = 20
	tableHeaderSeparatorWidth = 100
	maxItemTypeLength         = 10
	itemTypeTruncateLength    = 7
	maxStateLength            = 8
	stateTruncateLength       = 5
	maxNumberLength           = 6
	numberTruncateLength      = 3
	maxTitleLength            = 28
	listTitleTruncateLength   = 25
	maxRepositoryLength       = 18
	repositoryTruncateLength  = 15
	maxAuthorLength           = 13
	authorTruncateLength      = 10
	dateOnlyLength            = 10
)

// ListOptions holds options for the list command
type ListOptions struct {
	Repository string
	Search     string
	Type       string
	State      string
	Author     string
	Assignee   string
	Format     string
	Labels     []string
	Limit      int
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list [repository]",
		Short: "List issues and pull requests",
		Long: `List issues and pull requests from repositories or search across GitHub.

You can list items from a specific repository or search across all of GitHub
using various filters.

Examples:
  ghx item list octocat/Hello-World                    # List items from repository
  ghx item list octocat/Hello-World --type issue       # List only issues
  ghx item list --search "is:issue is:open bug"       # Search across GitHub
  ghx item list --author octocat --state open          # Find items by author
  ghx item list --assignee @me --type pr               # Find PRs assigned to you`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Repository = args[0]
			}
			return runList(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Search, "search", "", "Search query (GitHub search syntax)")
	cmd.Flags().StringVar(&opts.Type, "type", "", "Item type: issue, pr, pullrequest")
	cmd.Flags().StringVar(&opts.State, "state", "", "Item state: open, closed, merged")
	cmd.Flags().StringVar(&opts.Author, "author", "", "Filter by author username")
	cmd.Flags().StringVar(&opts.Assignee, "assignee", "", "Filter by assignee username")
	cmd.Flags().StringSliceVar(&opts.Labels, "label", nil, "Filter by labels (can be used multiple times)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", defaultListLimit, "Maximum number of items to list")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

func runList(ctx context.Context, opts *ListOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	itemService := service.NewItemService(client)

	var items []service.ItemInfo

	if opts.Repository != "" {
		// List items from specific repository
		items, err = listRepositoryItems(ctx, itemService, opts)
	} else {
		// Search items across GitHub
		items, err = searchItems(ctx, itemService, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to list items: %w", err)
	}

	return outputItems(items, opts.Format)
}

func listRepositoryItems(ctx context.Context, itemService *service.ItemService, opts *ListOptions) ([]service.ItemInfo, error) {
	// Parse repository reference
	parts := strings.Split(opts.Repository, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format: %s (expected owner/repo)", opts.Repository)
	}
	owner, repo := parts[0], parts[1]

	var states []string
	if opts.State != "" {
		states = []string{strings.ToUpper(opts.State)}
	}

	var allItems []service.ItemInfo

	// Get issues if not specifically requesting PRs
	if opts.Type == "" || opts.Type == "issue" {
		issues, err := itemService.ListRepositoryIssues(ctx, owner, repo, states, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}
		allItems = append(allItems, issues...)
	}

	// Get PRs if not specifically requesting issues
	if opts.Type == "" || opts.Type == "pr" || opts.Type == "pullrequest" {
		prs, err := itemService.ListRepositoryPullRequests(ctx, owner, repo, states, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}
		allItems = append(allItems, prs...)
	}

	// Apply additional filters
	return applyFilters(allItems, opts), nil
}

func searchItems(ctx context.Context, itemService *service.ItemService, opts *ListOptions) ([]service.ItemInfo, error) {
	// Build search query
	filters := service.SearchFilters{
		Type:       opts.Type,
		State:      opts.State,
		Repository: opts.Repository,
		Author:     opts.Author,
		Assignee:   opts.Assignee,
		Labels:     opts.Labels,
		Query:      opts.Search,
	}

	searchQuery := service.BuildSearchQuery(&filters)
	if searchQuery == "" {
		return nil, fmt.Errorf("no search criteria specified")
	}

	var allItems []service.ItemInfo

	// Search issues if not specifically requesting PRs
	if opts.Type == "" || opts.Type == "issue" {
		issues, err := itemService.SearchIssues(ctx, searchQuery, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search issues: %w", err)
		}
		allItems = append(allItems, issues...)
	}

	// Search PRs if not specifically requesting issues
	if opts.Type == "" || opts.Type == "pr" || opts.Type == "pullrequest" {
		prs, err := itemService.SearchPullRequests(ctx, searchQuery, opts.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search pull requests: %w", err)
		}
		allItems = append(allItems, prs...)
	}

	return allItems, nil
}

func applyFilters(items []service.ItemInfo, opts *ListOptions) []service.ItemInfo {
	filtered := make([]service.ItemInfo, 0, len(items))

	for i := range items {
		if passesAuthorFilter(&items[i], opts.Author) &&
			passesAssigneeFilter(&items[i], opts.Assignee) &&
			passesLabelsFilter(&items[i], opts.Labels) {
			filtered = append(filtered, items[i])
		}
	}

	return applyLimit(filtered, opts.Limit)
}

func passesAuthorFilter(item *service.ItemInfo, authorFilter string) bool {
	if authorFilter == "" {
		return true
	}
	return item.Author != nil && *item.Author == authorFilter
}

func passesAssigneeFilter(item *service.ItemInfo, assigneeFilter string) bool {
	if assigneeFilter == "" {
		return true
	}

	for _, assignee := range item.Assignees {
		if assignee == assigneeFilter || (assigneeFilter == "@me" && assignee == "current-user") {
			return true
		}
	}
	return false
}

func passesLabelsFilter(item *service.ItemInfo, labelFilters []string) bool {
	if len(labelFilters) == 0 {
		return true
	}

	for _, requiredLabel := range labelFilters {
		if !itemHasLabel(item, requiredLabel) {
			return false
		}
	}
	return true
}

func itemHasLabel(item *service.ItemInfo, requiredLabel string) bool {
	for _, itemLabel := range item.Labels {
		if itemLabel == requiredLabel {
			return true
		}
	}
	return false
}

func applyLimit(items []service.ItemInfo, limit int) []service.ItemInfo {
	if limit > 0 && len(items) > limit {
		return items[:limit]
	}
	return items
}

func outputItems(items []service.ItemInfo, format string) error {
	if len(items) == 0 {
		fmt.Println("No items found")
		return nil
	}

	switch format {
	case formatJSON:
		return outputItemsJSON(items)
	case "table":
		return outputItemsTable(items)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputItemsTable(items []service.ItemInfo) error {
	printItemTableHeader()

	for i := range items {
		printItemTableRow(&items[i])
	}

	return nil
}

func printItemTableHeader() {
	fmt.Printf("%-10s %-8s %-6s %-30s %-20s %-15s %-12s\n",
		"TYPE", "STATE", "NUMBER", "TITLE", "REPOSITORY", "AUTHOR", "UPDATED")
	fmt.Println(strings.Repeat("-", tableHeaderSeparatorWidth))
}

func printItemTableRow(item *service.ItemInfo) {
	itemType := truncateString(item.Type, maxItemTypeLength, itemTypeTruncateLength)
	state := truncateString(item.State, maxStateLength, stateTruncateLength)
	number := formatItemNumber(item.Number)
	title := truncateString(item.Title, maxTitleLength, listTitleTruncateLength)
	repository := formatItemRepository(item.Repository)
	author := formatItemAuthor(item.Author)
	updated := formatItemDate(item.UpdatedAt)

	fmt.Printf("%-10s %-8s %-6s %-30s %-20s %-15s %-12s\n",
		itemType, state, number, title, repository, author, updated)
}

func truncateString(s string, maxLen, truncateLen int) string {
	if len(s) > maxLen {
		return s[:truncateLen] + "..."
	}
	return s
}

func formatItemNumber(number *int) string {
	if number == nil {
		return ""
	}
	result := fmt.Sprintf("#%d", *number)
	return truncateString(result, maxNumberLength, numberTruncateLength)
}

func formatItemRepository(repository *string) string {
	if repository == nil {
		return ""
	}
	return truncateString(*repository, maxRepositoryLength, repositoryTruncateLength)
}

func formatItemAuthor(author *string) string {
	if author == nil {
		return ""
	}
	return truncateString(*author, maxAuthorLength, authorTruncateLength)
}

func formatItemDate(updated string) string {
	if len(updated) > dateOnlyLength {
		return updated[:dateOnlyLength] // Keep only date part
	}
	return updated
}

func outputItemsJSON(items []service.ItemInfo) error {
	// Simplified JSON output
	fmt.Println("[")
	for i := range items {
		item := &items[i]
		fmt.Printf("  {\n")
		fmt.Printf("    \"type\": \"%s\",\n", item.Type)
		fmt.Printf("    \"title\": \"%s\",\n", item.Title)
		if item.Number != nil {
			fmt.Printf("    \"number\": %d,\n", *item.Number)
		}
		fmt.Printf("    \"state\": \"%s\",\n", item.State)
		if item.Repository != nil {
			fmt.Printf("    \"repository\": \"%s\",\n", *item.Repository)
		}
		if item.Author != nil {
			fmt.Printf("    \"author\": \"%s\",\n", *item.Author)
		}
		fmt.Printf("    \"updated_at\": \"%s\"\n", item.UpdatedAt)

		if i < len(items)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")

	return nil
}
