package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

// ItemService handles item-related operations
type ItemService struct {
	client *api.Client
}

// NewItemService creates a new item service
func NewItemService(client *api.Client) *ItemService {
	return &ItemService{
		client: client,
	}
}

// ItemInfo represents simplified item information for display
type ItemInfo struct {
	ID         string
	Title      string
	Number     *int
	URL        *string
	Type       string
	State      string
	Repository *string
	Author     *string
	CreatedAt  string
	UpdatedAt  string
	Labels     []string
	Assignees  []string
}

// GetIssue gets a specific issue by repository and number
func (s *ItemService) GetIssue(ctx context.Context, owner, repo string, number int) (*graphql.Issue, error) {
	variables := graphql.BuildGetIssueVariables(owner, repo, number)

	var query graphql.GetIssueQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return &query.Repository.Issue, nil
}

// GetPullRequest gets a specific pull request by repository and number
func (s *ItemService) GetPullRequest(ctx context.Context, owner, repo string, number int) (*graphql.PullRequest, error) {
	variables := graphql.BuildGetPullRequestVariables(owner, repo, number)

	var query graphql.GetPullRequestQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return &query.Repository.PullRequest, nil
}

// SearchIssues searches for issues using GitHub search syntax
func (s *ItemService) SearchIssues(ctx context.Context, searchQuery string, limit int) ([]ItemInfo, error) {
	if limit <= 0 {
		limit = 10
	}

	opts := graphql.SearchOptions{
		Query: searchQuery,
		First: limit,
	}

	variables := graphql.BuildSearchIssuesVariables(opts)

	var query graphql.SearchIssuesQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %w", err)
	}

	items := make([]ItemInfo, 0, len(query.Search.Nodes))
	for i := range query.Search.Nodes {
		node := &query.Search.Nodes[i]
		issue := node.Issue

		labels := make([]string, len(issue.Labels.Nodes))
		for i, label := range issue.Labels.Nodes {
			labels[i] = label.Name
		}

		assignees := make([]string, len(issue.Assignees.Nodes))
		for i, assignee := range issue.Assignees.Nodes {
			assignees[i] = assignee.Login
		}

		items = append(items, ItemInfo{
			ID:         issue.ID,
			Title:      issue.Title,
			Number:     &issue.Number,
			URL:        &issue.URL,
			Type:       "Issue",
			State:      issue.State,
			Repository: &issue.Repository.NameWithOwner,
			Author:     &issue.Author.Login,
			CreatedAt:  issue.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  issue.UpdatedAt.Format("2006-01-02 15:04:05"),
			Labels:     labels,
			Assignees:  assignees,
		})
	}

	return items, nil
}

// SearchPullRequests searches for pull requests using GitHub search syntax
func (s *ItemService) SearchPullRequests(ctx context.Context, searchQuery string, limit int) ([]ItemInfo, error) {
	if limit <= 0 {
		limit = 10
	}

	opts := graphql.SearchOptions{
		Query: searchQuery,
		First: limit,
	}

	variables := graphql.BuildSearchPullRequestsVariables(opts)

	var query graphql.SearchPullRequestsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to search pull requests: %w", err)
	}

	items := make([]ItemInfo, 0, len(query.Search.Nodes))
	for i := range query.Search.Nodes {
		node := &query.Search.Nodes[i]
		pr := node.PullRequest

		labels := make([]string, len(pr.Labels.Nodes))
		for i, label := range pr.Labels.Nodes {
			labels[i] = label.Name
		}

		assignees := make([]string, len(pr.Assignees.Nodes))
		for i, assignee := range pr.Assignees.Nodes {
			assignees[i] = assignee.Login
		}

		state := pr.State
		if pr.Merged {
			state = "MERGED"
		}

		items = append(items, ItemInfo{
			ID:         pr.ID,
			Title:      pr.Title,
			Number:     &pr.Number,
			URL:        &pr.URL,
			Type:       "PullRequest",
			State:      state,
			Repository: &pr.Repository.NameWithOwner,
			Author:     &pr.Author.Login,
			CreatedAt:  pr.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  pr.UpdatedAt.Format("2006-01-02 15:04:05"),
			Labels:     labels,
			Assignees:  assignees,
		})
	}

	return items, nil
}

// ListRepositoryIssues lists issues in a repository
func (s *ItemService) ListRepositoryIssues(ctx context.Context, owner, repo string, states []string, limit int) ([]ItemInfo, error) {
	if limit <= 0 {
		limit = 10
	}

	opts := graphql.ListIssueOptions{
		Owner:  owner,
		Repo:   repo,
		States: states,
		First:  limit,
	}

	variables := graphql.BuildListIssuesVariables(opts)

	var query graphql.ListRepositoryIssuesQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list repository issues: %w", err)
	}

	items := make([]ItemInfo, len(query.Repository.Issues.Nodes))
	for i := range query.Repository.Issues.Nodes {
		issue := &query.Repository.Issues.Nodes[i]
		labels := make([]string, len(issue.Labels.Nodes))
		for j, label := range issue.Labels.Nodes {
			labels[j] = label.Name
		}

		assignees := make([]string, len(issue.Assignees.Nodes))
		for j, assignee := range issue.Assignees.Nodes {
			assignees[j] = assignee.Login
		}

		items[i] = ItemInfo{
			ID:         issue.ID,
			Title:      issue.Title,
			Number:     &issue.Number,
			URL:        &issue.URL,
			Type:       "Issue",
			State:      issue.State,
			Repository: &issue.Repository.NameWithOwner,
			Author:     &issue.Author.Login,
			CreatedAt:  issue.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  issue.UpdatedAt.Format("2006-01-02 15:04:05"),
			Labels:     labels,
			Assignees:  assignees,
		}
	}

	return items, nil
}

// ListRepositoryPullRequests lists pull requests in a repository
func (s *ItemService) ListRepositoryPullRequests(ctx context.Context, owner, repo string, states []string, limit int) ([]ItemInfo, error) {
	if limit <= 0 {
		limit = 10
	}

	opts := graphql.ListPullRequestOptions{
		Owner:  owner,
		Repo:   repo,
		States: states,
		First:  limit,
	}

	variables := graphql.BuildListPullRequestsVariables(opts)

	var query graphql.ListRepositoryPullRequestsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list repository pull requests: %w", err)
	}

	items := make([]ItemInfo, len(query.Repository.PullRequests.Nodes))
	for i := range query.Repository.PullRequests.Nodes {
		pr := &query.Repository.PullRequests.Nodes[i]
		labels := make([]string, len(pr.Labels.Nodes))
		for j, label := range pr.Labels.Nodes {
			labels[j] = label.Name
		}

		assignees := make([]string, len(pr.Assignees.Nodes))
		for j, assignee := range pr.Assignees.Nodes {
			assignees[j] = assignee.Login
		}

		state := pr.State
		if pr.Merged {
			state = "MERGED"
		}

		items[i] = ItemInfo{
			ID:         pr.ID,
			Title:      pr.Title,
			Number:     &pr.Number,
			URL:        &pr.URL,
			Type:       "PullRequest",
			State:      state,
			Repository: &pr.Repository.NameWithOwner,
			Author:     &pr.Author.Login,
			CreatedAt:  pr.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  pr.UpdatedAt.Format("2006-01-02 15:04:05"),
			Labels:     labels,
			Assignees:  assignees,
		}
	}

	return items, nil
}

// AddItemToProject adds an existing issue or PR to a project
func (s *ItemService) AddItemToProject(ctx context.Context, projectID, contentID string) (*graphql.ProjectV2Item, error) {
	input := AddItemInput{
		ProjectID: projectID,
		ContentID: contentID,
	}

	projectService := NewProjectService(s.client)
	return projectService.AddItem(ctx, input)
}

// CreateDraftIssue creates a draft issue in a project
func (s *ItemService) CreateDraftIssue(ctx context.Context, projectID, title string, body *string) (*graphql.ProjectV2Item, error) {
	input := graphql.CreateDraftIssueInput{
		ProjectID: projectID,
		Title:     title,
		Body:      body,
	}

	variables := graphql.BuildCreateDraftIssueVariables(input)

	var mutation graphql.CreateDraftIssueMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create draft issue: %w", err)
	}

	return &mutation.AddProjectV2DraftIssue.ProjectItem, nil
}

// UpdateDraftIssue updates a draft issue
func (s *ItemService) UpdateDraftIssue(ctx context.Context, draftIssueID string, title, body *string) (*graphql.DraftIssue, error) {
	input := graphql.UpdateDraftIssueInput{
		DraftIssueID: draftIssueID,
		Title:        title,
		Body:         body,
	}

	variables := graphql.BuildUpdateDraftIssueVariables(input)

	var mutation graphql.UpdateDraftIssueMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update draft issue: %w", err)
	}

	return &mutation.UpdateProjectV2DraftIssue.DraftIssue, nil
}

// RemoveItemFromProject removes an item from a project
func (s *ItemService) RemoveItemFromProject(ctx context.Context, projectID, itemID string) error {
	input := RemoveItemInput{
		ProjectID: projectID,
		ItemID:    itemID,
	}

	projectService := NewProjectService(s.client)
	return projectService.RemoveItem(ctx, input)
}

// ParseItemReference parses an item reference in various formats
func ParseItemReference(ref string) (owner, repo string, number int, err error) {
	// Handle different formats:
	// - owner/repo#123
	// - https://github.com/owner/repo/issues/123
	// - https://github.com/owner/repo/pull/123
	// - #123 (requires current repo context)

	if strings.HasPrefix(ref, "https://github.com/") {
		return parseGitHubURL(ref)
	}

	if strings.Contains(ref, "#") {
		parts := strings.Split(ref, "#")
		if len(parts) != 2 {
			return "", "", 0, fmt.Errorf("invalid item reference format: %s", ref)
		}

		repoPath := parts[0]
		numberStr := parts[1]

		if repoPath == "" {
			return "", "", 0, fmt.Errorf("repository context required for reference: %s", ref)
		}

		repoParts := strings.Split(repoPath, "/")
		if len(repoParts) != 2 {
			return "", "", 0, fmt.Errorf("invalid repository format in reference: %s", ref)
		}

		owner = repoParts[0]
		repo = repoParts[1]

		number, err = strconv.Atoi(numberStr)
		if err != nil {
			return "", "", 0, fmt.Errorf("invalid item number in reference: %s", numberStr)
		}

		return owner, repo, number, nil
	}

	return "", "", 0, fmt.Errorf("unrecognized item reference format: %s", ref)
}

// parseGitHubURL parses a GitHub URL to extract owner, repo, and number
func parseGitHubURL(url string) (owner, repo string, number int, err error) {
	// Remove https://github.com/ prefix
	path := strings.TrimPrefix(url, "https://github.com/")

	parts := strings.Split(path, "/")
	if len(parts) < minSearchPartsLength {
		return "", "", 0, fmt.Errorf("invalid GitHub URL format: %s", url)
	}

	owner = parts[0]
	repo = parts[1]
	// parts[2] is "issues" or "pull"
	numberStr := parts[3]

	number, err = strconv.Atoi(numberStr)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid item number in URL: %s", numberStr)
	}

	return owner, repo, number, nil
}

// FormatItemReference formats owner, repo, and number into an item reference
func FormatItemReference(owner, repo string, number int) string {
	return fmt.Sprintf("%s/%s#%d", owner, repo, number)
}

// BuildSearchQuery builds a search query for items based on filters
func BuildSearchQuery(filters *SearchFilters) string {
	parts := make([]string, 0, defaultSearchPartsSize)

	// Base type filter
	if filters.Type == "issue" {
		parts = append(parts, "is:issue")
	} else if filters.Type == "pr" || filters.Type == "pullrequest" {
		parts = append(parts, "is:pr")
	}

	// State filter
	if filters.State != "" {
		parts = append(parts, fmt.Sprintf("is:%s", filters.State))
	}

	// Repository filter
	if filters.Repository != "" {
		parts = append(parts, fmt.Sprintf("repo:%s", filters.Repository))
	}

	// Author filter
	if filters.Author != "" {
		parts = append(parts, fmt.Sprintf("author:%s", filters.Author))
	}

	// Assignee filter
	if filters.Assignee != "" {
		parts = append(parts, fmt.Sprintf("assignee:%s", filters.Assignee))
	}

	// Label filters
	for _, label := range filters.Labels {
		parts = append(parts, fmt.Sprintf("label:%s", label))
	}

	// Text search
	if filters.Query != "" {
		parts = append(parts, filters.Query)
	}

	return strings.Join(parts, " ")
}

// SearchFilters represents filters for searching items
type SearchFilters struct {
	Type       string
	State      string
	Repository string
	Author     string
	Assignee   string
	Query      string
	Labels     []string
}

// CreateItemInput represents input for creating an item
type CreateItemInput struct {
	ProjectID   string
	Title       string
	Body        string
	ContentType string  // "issue", "pull_request", "draft_issue"
	ContentID   *string // GitHub issue/PR ID if linking existing content
}

// BulkUpdateInput represents input for bulk update operations
type BulkUpdateInput struct {
	ProjectID string
	ItemIDs   []string
	FieldName string
	Value     interface{}
}

// BulkAddInput represents input for bulk add operations
type BulkAddInput struct {
	ProjectID string
	Items     []CreateItemInput
}

// BulkUpdateResult represents result of bulk update operation
type BulkUpdateResult struct {
	Updated int
	Failed  int
	Errors  []string
}

// BulkAddResult represents result of bulk add operation
type BulkAddResult struct {
	Added  int
	Failed int
	Errors []string
}

// BulkUpdateItems updates multiple items with same field value
func (s *ItemService) BulkUpdateItems(ctx context.Context, input BulkUpdateInput) (*BulkUpdateResult, error) {
	result := &BulkUpdateResult{}

	// For each item, update the field value using GraphQL mutation
	for _, itemID := range input.ItemIDs {
		// Use the existing UpdateItemField GraphQL mutation
		var mutation graphql.UpdateItemFieldMutation
		variables := map[string]interface{}{
			"projectId": input.ProjectID,
			"itemId":    itemID,
			"fieldId":   input.FieldName, // This should be field ID, not name
			"value": map[string]interface{}{
				"text": input.Value, // Assuming text field for now
			},
		}

		err := s.client.Mutate(ctx, &mutation, variables)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("failed to update item %s: %v", itemID, err))
		} else {
			result.Updated++
		}
	}

	return result, nil
}

// BulkAddItems adds multiple items to a project
func (s *ItemService) BulkAddItems(ctx context.Context, input BulkAddInput) (*BulkAddResult, error) {
	result := &BulkAddResult{}

	for _, item := range input.Items {
		// Use the existing AddItemToProject GraphQL mutation
		var mutation graphql.AddItemToProjectMutation
		variables := map[string]interface{}{
			"projectId": input.ProjectID,
			"contentId": item.ContentID, // This would be the issue/PR ID
		}

		err := s.client.Mutate(ctx, &mutation, variables)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("failed to add item %s: %v", item.Title, err))
		} else {
			result.Added++
		}
	}

	return result, nil
}

// GetItemsByFilter retrieves items based on filter criteria
func (s *ItemService) GetItemsByFilter(ctx context.Context, projectID, filter string) ([]string, error) {
	// Parse filter string (e.g., "label:epic", "assignee:@me")
	parts := strings.Split(filter, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid filter format, expected 'key:value'")
	}

	filterType := strings.TrimSpace(parts[0])
	filterValue := strings.TrimSpace(parts[1])

	switch filterType {
	case "label":
		return s.GetItemsByLabel(ctx, filterValue)
	case "assignee":
		return s.getItemsByAssignee(ctx, filterValue)
	case "state":
		return s.getItemsByState(ctx, filterValue)
	default:
		return nil, fmt.Errorf("unsupported filter type: %s", filterType)
	}
}

// GetItemsByLabel retrieves items with specific label
func (s *ItemService) GetItemsByLabel(ctx context.Context, label string) ([]string, error) {
	return s.searchIssuesByQuery(ctx, fmt.Sprintf("label:%s", label), "failed to search issues by label")
}

// getItemsByAssignee retrieves items by assignee
func (s *ItemService) getItemsByAssignee(ctx context.Context, assignee string) ([]string, error) {
	searchAssignee := assignee
	if assignee == "@me" {
		searchAssignee = "@me"
	}

	return s.searchIssuesByQuery(ctx, fmt.Sprintf("assignee:%s", searchAssignee), "failed to search issues by assignee")
}

// getItemsByState retrieves items by state
func (s *ItemService) getItemsByState(ctx context.Context, state string) ([]string, error) {
	return s.searchIssuesByQuery(ctx, fmt.Sprintf("state:%s", state), "failed to search issues by state")
}

// searchIssuesByQuery is a helper function to search issues with a query string
func (s *ItemService) searchIssuesByQuery(ctx context.Context, query, errorPrefix string) ([]string, error) {
	var gqlQuery graphql.SearchIssuesQuery
	variables := graphql.BuildSearchIssuesVariables(graphql.SearchOptions{
		Query: query,
		First: 100,
	})

	err := s.client.Query(ctx, &gqlQuery, variables)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorPrefix, err)
	}

	var itemIDs []string
	for _, node := range gqlQuery.Search.Nodes {
		if node.Issue.ID != "" {
			itemIDs = append(itemIDs, node.Issue.ID)
		}
	}

	return itemIDs, nil
}

// ParseNumberRange parses number range string (e.g., "34-46") into item IDs
func ParseNumberRange(rangeStr string) ([]string, error) {
	if strings.Contains(rangeStr, "-") {
		parts := strings.Split(rangeStr, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format: %s", rangeStr)
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start number: %s", parts[0])
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end number: %s", parts[1])
		}

		if start > end {
			return nil, fmt.Errorf("start number cannot be greater than end number")
		}

		var result []string
		for i := start; i <= end; i++ {
			result = append(result, fmt.Sprintf("item_%d", i))
		}
		return result, nil
	}

	// Single number
	num, err := strconv.Atoi(strings.TrimSpace(rangeStr))
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", rangeStr)
	}

	return []string{fmt.Sprintf("item_%d", num)}, nil
}

// RemoveDuplicates removes duplicate strings from slice
func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
