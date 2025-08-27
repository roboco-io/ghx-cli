package item

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/api/graphql"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

type bulkAddOptions struct {
	issues   string
	label    string
	fromFile string
}

// NewAddBulkCmd creates the add-bulk command
func NewAddBulkCmd() *cobra.Command {
	var opts bulkAddOptions

	cmd := &cobra.Command{
		Use:   "add-bulk PROJECT_ID",
		Short: "Add multiple issues to a project in bulk",
		Long: `Add multiple issues or pull requests to a GitHub Project in bulk.

This command allows you to add multiple items at once using various methods:
• Number range (e.g., 34-46)
• By label
• From a file containing issue URLs or numbers

Examples:
  # Add issues by number range
  ghp item add-bulk myorg/123 --issues 34-46
  
  # Add all issues with a specific label
  ghp item add-bulk myorg/123 --label epic
  
  # Add issues from a file (one per line)
  ghp item add-bulk myorg/123 --from-file issue-list.txt`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeBulkAdd(cmd.Context(), args[0], opts)
		},
	}

	cmd.Flags().StringVar(&opts.issues, "issues", "", "Issue number range (e.g., 34-46)")
	cmd.Flags().StringVar(&opts.label, "label", "", "Add all issues with this label")
	cmd.Flags().StringVar(&opts.fromFile, "from-file", "", "File containing issue URLs or numbers (one per line)")

	return cmd
}

// executeBulkAdd executes the bulk add operation
func executeBulkAdd(ctx context.Context, projectID string, opts bulkAddOptions) error {
	// Validate input options
	if err := validateBulkAddOptions(opts); err != nil {
		return err
	}

	// Collect items from all sources
	itemsToAdd, err := collectItemsToAdd(ctx, opts)
	if err != nil {
		return err
	}

	// Remove duplicates
	itemsToAdd = removeDuplicates(itemsToAdd)
	fmt.Printf("Adding %d items to project %s...\n", len(itemsToAdd), projectID)

	// Initialize services
	itemService, projectService, err := initializeServices()
	if err != nil {
		return err
	}

	// Get project
	project, err := getProjectFromID(ctx, projectService, projectID)
	if err != nil {
		return err
	}

	// Execute bulk add
	return executeBulkAddToProject(ctx, itemService, project.ID, itemsToAdd)
}

// validateBulkAddOptions validates the bulk add options
func validateBulkAddOptions(opts bulkAddOptions) error {
	if opts.issues == "" && opts.label == "" && opts.fromFile == "" {
		return fmt.Errorf("at least one of --issues, --label, or --from-file must be specified")
	}
	return nil
}

// collectItemsToAdd collects items from all specified sources
func collectItemsToAdd(ctx context.Context, opts bulkAddOptions) ([]string, error) {
	var itemsToAdd []string

	// Handle number range
	if opts.issues != "" {
		items, err := parseNumberRange(opts.issues)
		if err != nil {
			return nil, fmt.Errorf("invalid issue range: %w", err)
		}
		itemsToAdd = append(itemsToAdd, items...)
	}

	// Handle label
	if opts.label != "" {
		items, err := getIssuesByLabel(ctx, opts.label)
		if err != nil {
			return nil, fmt.Errorf("failed to get issues by label: %w", err)
		}
		itemsToAdd = append(itemsToAdd, items...)
	}

	// Handle file input
	if opts.fromFile != "" {
		items, err := readIssuesFromFile(opts.fromFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read issues from file: %w", err)
		}
		itemsToAdd = append(itemsToAdd, items...)
	}

	return itemsToAdd, nil
}

// initializeServices initializes authentication and services
func initializeServices() (*service.ItemService, *service.ProjectService, error) {
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return nil, nil, fmt.Errorf("authentication failed: %w", err)
	}

	client := api.NewClient(token)
	return service.NewItemService(client), service.NewProjectService(client), nil
}

// getProjectFromID parses project ID and retrieves project
func getProjectFromID(ctx context.Context, projectService *service.ProjectService, projectID string) (*graphql.ProjectV2, error) {
	parts := strings.Split(projectID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid project format: %s (expected owner/number)", projectID)
	}

	owner := parts[0]
	projectNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid project number: %s", parts[1])
	}

	project, err := projectService.GetProjectWithOwnerDetection(ctx, owner, projectNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// executeBulkAddToProject executes the bulk add operation to the project
func executeBulkAddToProject(ctx context.Context, itemService *service.ItemService, projectID string, itemsToAdd []string) error {
	bulkInput := service.BulkAddInput{
		ProjectID: projectID,
		Items:     make([]service.CreateItemInput, len(itemsToAdd)),
	}

	for i, item := range itemsToAdd {
		bulkInput.Items[i] = service.CreateItemInput{
			Title:       fmt.Sprintf("Item %s", item),
			Body:        fmt.Sprintf("Item added from bulk operation: %s", item),
			ContentType: "issue",
			ContentID:   &item,
		}
	}

	result, err := itemService.BulkAddItems(ctx, bulkInput)
	if err != nil {
		return fmt.Errorf("failed to add items in bulk: %w", err)
	}

	fmt.Printf("\n✓ Successfully added %d items to project\n", result.Added)
	return nil
}

// parseNumberRange parses a number range like "34-46" into a slice of strings
func parseNumberRange(rangeStr string) ([]string, error) {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range format, expected 'start-end'")
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid start number: %w", err)
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid end number: %w", err)
	}

	if start > end {
		return nil, fmt.Errorf("start number must be less than or equal to end number")
	}

	var result []string
	for i := start; i <= end; i++ {
		result = append(result, fmt.Sprintf("#%d", i))
	}

	return result, nil
}

// getIssuesByLabel retrieves issues with a specific label
func getIssuesByLabel(ctx context.Context, label string) ([]string, error) {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	itemService := service.NewItemService(client)

	// Use the service method to get items by label
	return itemService.GetItemsByLabel(ctx, label)
}

// readIssuesFromFile reads issue URLs or numbers from a file
func readIssuesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var issues []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") { // Skip comments
			issues = append(issues, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}

// removeDuplicates removes duplicate items from a slice
func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
