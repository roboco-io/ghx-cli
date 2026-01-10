package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	gql "github.com/shurcooL/graphql"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

// DiscussionService handles discussion-related operations
type DiscussionService struct {
	client *api.Client
}

// NewDiscussionService creates a new discussion service
func NewDiscussionService(client *api.Client) *DiscussionService {
	return &DiscussionService{
		client: client,
	}
}

// =============================================================================
// OUTPUT TYPES
// =============================================================================

// DiscussionInfo represents simplified discussion information
type DiscussionInfo struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ID           string
	Title        string
	Body         string
	URL          string
	State        string
	Category     string
	CategorySlug string
	Author       string
	Labels       []string
	Number       int
	CommentCount int
	UpvoteCount  int
	Locked       bool
	HasAnswer    bool
}

// DiscussionDetails represents detailed discussion information
type DiscussionDetails struct {
	DiscussionInfo
	ClosedAt     *time.Time
	Answer       *CommentInfo
	BodyHTML     string
	Comments     []CommentInfo
	CategoryInfo CategoryInfo
}

// CommentInfo represents comment information
type CommentInfo struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ID          string
	Body        string
	BodyHTML    string
	Author      string
	UpvoteCount int
	IsAnswer    bool
}

// CategoryInfo represents category information
type CategoryInfo struct {
	ID           string
	Name         string
	Slug         string
	Description  string
	Emoji        string
	IsAnswerable bool
}

// =============================================================================
// INPUT TYPES
// =============================================================================

// ListDiscussionsOptions represents options for listing discussions
type ListDiscussionsOptions struct {
	After    *string
	Answered *bool
	Owner    string
	Repo     string
	Category string
	State    string
	First    int
}

// CreateDiscussionOptions represents options for creating a discussion
type CreateDiscussionOptions struct {
	Owner    string
	Repo     string
	Category string
	Title    string
	Body     string
}

// UpdateDiscussionOptions represents options for updating a discussion
type UpdateDiscussionOptions struct {
	Title    *string
	Body     *string
	Category *string
	Owner    string
	Repo     string
	Number   int
}

// CloseDiscussionOptions represents options for closing a discussion
type CloseDiscussionOptions struct {
	Owner  string
	Repo   string
	Reason string
	Number int
}

// LockDiscussionOptions represents options for locking a discussion
type LockDiscussionOptions struct {
	Owner  string
	Repo   string
	Reason string
	Number int
}

// AddCommentOptions represents options for adding a comment
type AddCommentOptions struct {
	ReplyToID *string
	Owner     string
	Repo      string
	Body      string
	Number    int
}

// AnswerOptions represents options for marking/unmarking an answer
type AnswerOptions struct {
	Owner     string
	Repo      string
	CommentID string
	Number    int
	Unmark    bool
}

// =============================================================================
// QUERY METHODS
// =============================================================================

// ListDiscussions lists discussions for a repository
func (s *DiscussionService) ListDiscussions(ctx context.Context, opts ListDiscussionsOptions) ([]DiscussionInfo, error) {
	if opts.First <= 0 {
		opts.First = defaultDiscussionListLimit
	}

	// Get category ID if category slug provided
	var categoryID *string
	if opts.Category != "" {
		cat, err := s.GetCategoryBySlug(ctx, opts.Owner, opts.Repo, opts.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
		if cat == nil {
			return nil, fmt.Errorf("category not found: %s", opts.Category)
		}
		categoryID = &cat.ID
	}

	variables := graphql.BuildListDiscussionsVariables(
		opts.Owner, opts.Repo, opts.First, opts.After, categoryID, opts.Answered,
	)

	var query graphql.ListDiscussionsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list discussions: %w", err)
	}

	discussions := convertDiscussionNodes(query.Repository.Discussions.Nodes)

	// Filter by state if specified
	if opts.State != "" && opts.State != "all" {
		discussions = filterDiscussionsByState(discussions, opts.State)
	}

	return discussions, nil
}

// GetDiscussion gets a specific discussion by number
func (s *DiscussionService) GetDiscussion(ctx context.Context, owner, repo string, number, commentLimit int) (*DiscussionDetails, error) {
	if commentLimit <= 0 {
		commentLimit = defaultDiscussionCommentLimit
	}

	variables := graphql.BuildGetDiscussionVariables(owner, repo, number, commentLimit)

	var query graphql.GetDiscussionQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	return convertDiscussionToDetails(&query.Repository.Discussion), nil
}

// ListCategories lists discussion categories for a repository
func (s *DiscussionService) ListCategories(ctx context.Context, owner, repo string) ([]CategoryInfo, error) {
	variables := graphql.BuildListDiscussionCategoriesVariables(owner, repo)

	var query graphql.ListDiscussionCategoriesQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return convertCategoryNodes(query.Repository.DiscussionCategories.Nodes), nil
}

// GetCategoryBySlug gets a category by its slug
func (s *DiscussionService) GetCategoryBySlug(ctx context.Context, owner, repo, slug string) (*CategoryInfo, error) {
	variables := graphql.BuildGetDiscussionCategoryVariables(owner, repo, slug)

	var query graphql.GetDiscussionCategoryQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if query.Repository.DiscussionCategory == nil {
		return nil, nil
	}

	return convertCategory(query.Repository.DiscussionCategory), nil
}

// =============================================================================
// MUTATION METHODS
// =============================================================================

// CreateDiscussion creates a new discussion
func (s *DiscussionService) CreateDiscussion(ctx context.Context, opts CreateDiscussionOptions) (*DiscussionDetails, error) {
	// Get repository ID
	repoInfo, err := s.getRepositoryInfo(ctx, opts.Owner, opts.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	// Get category ID
	cat, err := s.GetCategoryBySlug(ctx, opts.Owner, opts.Repo, opts.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if cat == nil {
		return nil, fmt.Errorf("category not found: %s", opts.Category)
	}

	variables := graphql.BuildCreateDiscussionVariables(&graphql.CreateDiscussionInput{
		RepositoryID: gql.ID(repoInfo.ID),
		CategoryID:   gql.ID(cat.ID),
		Title:        gql.String(opts.Title),
		Body:         gql.String(opts.Body),
	})

	var mutation graphql.CreateDiscussionMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create discussion: %w", err)
	}

	return convertDiscussionToDetails(&mutation.CreateDiscussion.Discussion), nil
}

// UpdateDiscussion updates a discussion
func (s *DiscussionService) UpdateDiscussion(ctx context.Context, opts UpdateDiscussionOptions) (*DiscussionDetails, error) {
	// Get discussion ID first
	discussion, err := s.GetDiscussion(ctx, opts.Owner, opts.Repo, opts.Number, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	input := &graphql.UpdateDiscussionInput{
		DiscussionID: gql.ID(discussion.ID),
	}

	// Convert optional fields
	if opts.Title != nil {
		title := gql.String(*opts.Title)
		input.Title = &title
	}
	if opts.Body != nil {
		body := gql.String(*opts.Body)
		input.Body = &body
	}

	// Get category ID if changing category
	if opts.Category != nil {
		cat, catErr := s.GetCategoryBySlug(ctx, opts.Owner, opts.Repo, *opts.Category)
		if catErr != nil {
			return nil, fmt.Errorf("failed to get category: %w", catErr)
		}
		if cat == nil {
			return nil, fmt.Errorf("category not found: %s", *opts.Category)
		}
		catID := gql.ID(cat.ID)
		input.CategoryID = &catID
	}

	variables := graphql.BuildUpdateDiscussionVariables(input)

	var mutation graphql.UpdateDiscussionMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update discussion: %w", err)
	}

	return convertDiscussionToDetails(&mutation.UpdateDiscussion.Discussion), nil
}

// DeleteDiscussion deletes a discussion
func (s *DiscussionService) DeleteDiscussion(ctx context.Context, owner, repo string, number int) error {
	// Get discussion ID first
	discussion, err := s.GetDiscussion(ctx, owner, repo, number, 0)
	if err != nil {
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	variables := graphql.BuildDeleteDiscussionVariables(discussion.ID)

	var mutation graphql.DeleteDiscussionMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete discussion: %w", err)
	}

	return nil
}

// CloseDiscussion closes a discussion
func (s *DiscussionService) CloseDiscussion(ctx context.Context, opts CloseDiscussionOptions) (*DiscussionDetails, error) {
	discussion, err := s.GetDiscussion(ctx, opts.Owner, opts.Repo, opts.Number, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	reason := mapCloseReason(opts.Reason)
	variables := graphql.BuildCloseDiscussionVariables(discussion.ID, reason)

	var mutation graphql.CloseDiscussionMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to close discussion: %w", err)
	}

	return convertDiscussionToDetails(&mutation.CloseDiscussion.Discussion), nil
}

// ReopenDiscussion reopens a discussion
func (s *DiscussionService) ReopenDiscussion(ctx context.Context, owner, repo string, number int) (*DiscussionDetails, error) {
	discussion, err := s.GetDiscussion(ctx, owner, repo, number, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	variables := graphql.BuildReopenDiscussionVariables(discussion.ID)

	var mutation graphql.ReopenDiscussionMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to reopen discussion: %w", err)
	}

	return convertDiscussionToDetails(&mutation.ReopenDiscussion.Discussion), nil
}

// LockDiscussion locks a discussion
func (s *DiscussionService) LockDiscussion(ctx context.Context, opts LockDiscussionOptions) error {
	discussion, err := s.GetDiscussion(ctx, opts.Owner, opts.Repo, opts.Number, 0)
	if err != nil {
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	var reason *graphql.DiscussionLockReason
	if opts.Reason != "" {
		r := mapLockReason(opts.Reason)
		reason = &r
	}
	variables := graphql.BuildLockDiscussionVariables(discussion.ID, reason)

	var mutation graphql.LockDiscussionMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to lock discussion: %w", err)
	}

	return nil
}

// UnlockDiscussion unlocks a discussion
func (s *DiscussionService) UnlockDiscussion(ctx context.Context, owner, repo string, number int) error {
	discussion, err := s.GetDiscussion(ctx, owner, repo, number, 0)
	if err != nil {
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	variables := graphql.BuildUnlockDiscussionVariables(discussion.ID)

	var mutation graphql.UnlockDiscussionMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to unlock discussion: %w", err)
	}

	return nil
}

// AddComment adds a comment to a discussion
func (s *DiscussionService) AddComment(ctx context.Context, opts AddCommentOptions) (*CommentInfo, error) {
	discussion, err := s.GetDiscussion(ctx, opts.Owner, opts.Repo, opts.Number, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	variables := graphql.BuildAddDiscussionCommentVariables(discussion.ID, opts.Body, opts.ReplyToID)

	var mutation graphql.AddDiscussionCommentMutation
	err = s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	return convertComment(&mutation.AddDiscussionComment.Comment), nil
}

// MarkAnswer marks a comment as the answer
func (s *DiscussionService) MarkAnswer(ctx context.Context, commentID string) error {
	variables := graphql.BuildMarkAnswerVariables(commentID)

	var mutation graphql.MarkDiscussionCommentAsAnswerMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to mark answer: %w", err)
	}

	return nil
}

// UnmarkAnswer unmarks a comment as the answer
func (s *DiscussionService) UnmarkAnswer(ctx context.Context, commentID string) error {
	variables := graphql.BuildUnmarkAnswerVariables(commentID)

	var mutation graphql.UnmarkDiscussionCommentAsAnswerMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to unmark answer: %w", err)
	}

	return nil
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func (s *DiscussionService) getRepositoryInfo(ctx context.Context, owner, name string) (*graphql.RepositoryInfo, error) {
	var query graphql.RepositoryQuery
	variables := graphql.BuildRepositoryVariables(owner, name)

	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	return graphql.ParseRepositoryResponse(&query)
}

func convertDiscussionNodes(nodes []graphql.DiscussionSummary) []DiscussionInfo {
	discussions := make([]DiscussionInfo, len(nodes))
	for i := range nodes {
		discussions[i] = convertDiscussionSummary(&nodes[i])
	}
	return discussions
}

//nolint:dupl // DiscussionSummary and Discussion have same structure but different types
func convertDiscussionSummary(d *graphql.DiscussionSummary) DiscussionInfo {
	state := "OPEN"
	if d.Closed {
		state = "CLOSED"
	}

	labels := make([]string, 0, len(d.Labels.Nodes))
	for _, l := range d.Labels.Nodes {
		labels = append(labels, l.Name)
	}

	return DiscussionInfo{
		ID:           d.ID,
		Number:       d.Number,
		Title:        d.Title,
		Body:         d.Body,
		URL:          d.URL,
		State:        state,
		Locked:       d.Locked,
		Category:     d.Category.Name,
		CategorySlug: d.Category.Slug,
		Author:       d.Author.Login,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		CommentCount: d.Comments.TotalCount,
		UpvoteCount:  d.UpvoteCount,
		HasAnswer:    d.Answer != nil,
		Labels:       labels,
	}
}

//nolint:dupl // DiscussionSummary and Discussion have same structure but different types
func convertDiscussion(d *graphql.Discussion) DiscussionInfo {
	state := "OPEN"
	if d.Closed {
		state = "CLOSED"
	}

	labels := make([]string, 0, len(d.Labels.Nodes))
	for _, l := range d.Labels.Nodes {
		labels = append(labels, l.Name)
	}

	return DiscussionInfo{
		ID:           d.ID,
		Number:       d.Number,
		Title:        d.Title,
		Body:         d.Body,
		URL:          d.URL,
		State:        state,
		Locked:       d.Locked,
		Category:     d.Category.Name,
		CategorySlug: d.Category.Slug,
		Author:       d.Author.Login,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		CommentCount: d.Comments.TotalCount,
		UpvoteCount:  d.UpvoteCount,
		HasAnswer:    d.Answer != nil,
		Labels:       labels,
	}
}

func convertDiscussionToDetails(d *graphql.Discussion) *DiscussionDetails {
	info := convertDiscussion(d)

	details := &DiscussionDetails{
		DiscussionInfo: info,
		BodyHTML:       d.BodyHTML,
		ClosedAt:       d.ClosedAt,
		CategoryInfo: CategoryInfo{
			ID:           d.Category.ID,
			Name:         d.Category.Name,
			Slug:         d.Category.Slug,
			Description:  d.Category.Description,
			Emoji:        d.Category.Emoji,
			IsAnswerable: d.Category.IsAnswerable,
		},
	}

	if d.Answer != nil {
		details.Answer = convertComment(d.Answer)
	}

	for i := range d.Comments.Nodes {
		details.Comments = append(details.Comments, *convertComment(&d.Comments.Nodes[i]))
	}

	return details
}

func convertComment(c *graphql.DiscussionComment) *CommentInfo {
	return &CommentInfo{
		ID:          c.ID,
		Body:        c.Body,
		BodyHTML:    c.BodyHTML,
		Author:      c.Author.Login,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		IsAnswer:    c.IsAnswer,
		UpvoteCount: c.UpvoteCount,
	}
}

func convertCategoryNodes(nodes []graphql.DiscussionCategory) []CategoryInfo {
	categories := make([]CategoryInfo, len(nodes))
	for i := range nodes {
		categories[i] = *convertCategory(&nodes[i])
	}
	return categories
}

func convertCategory(c *graphql.DiscussionCategory) *CategoryInfo {
	return &CategoryInfo{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		Emoji:        c.Emoji,
		IsAnswerable: c.IsAnswerable,
	}
}

func filterDiscussionsByState(discussions []DiscussionInfo, state string) []DiscussionInfo {
	filtered := make([]DiscussionInfo, 0, len(discussions))
	targetState := strings.ToUpper(state)

	for _, d := range discussions {
		if d.State == targetState {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

func mapCloseReason(reason string) graphql.DiscussionCloseReason {
	switch strings.ToLower(reason) {
	case "outdated":
		return graphql.DiscussionCloseReasonOutdated
	case "duplicate":
		return graphql.DiscussionCloseReasonDuplicate
	default:
		return graphql.DiscussionCloseReasonResolved
	}
}

func mapLockReason(reason string) graphql.DiscussionLockReason {
	switch strings.ToLower(reason) {
	case "off_topic", "off-topic", "offtopic":
		return graphql.DiscussionLockReasonOffTopic
	case "spam":
		return graphql.DiscussionLockReasonSpam
	case "too_heated", "too-heated", "tooheated":
		return graphql.DiscussionLockReasonTooHeated
	default:
		return graphql.DiscussionLockReasonResolved
	}
}

// ParseRepositoryReference parses owner/repo format
func ParseRepositoryReference(ref string) (owner, repo string, err error) {
	parts := strings.Split(ref, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s (expected owner/repo)", ref)
	}
	return parts[0], parts[1], nil
}

// ValidateCloseReason validates a close reason string
func ValidateCloseReason(reason string) error {
	switch strings.ToLower(reason) {
	case "resolved", "outdated", "duplicate":
		return nil
	default:
		return fmt.Errorf("invalid close reason: %s (valid: resolved, outdated, duplicate)", reason)
	}
}

// ValidateLockReason validates a lock reason string
func ValidateLockReason(reason string) error {
	switch strings.ToLower(reason) {
	case "", "off_topic", "off-topic", "resolved", "spam", "too_heated", "too-heated":
		return nil
	default:
		return fmt.Errorf("invalid lock reason: %s (valid: off_topic, resolved, spam, too_heated)", reason)
	}
}

// Discussion service constants
const (
	defaultDiscussionListLimit    = 20
	defaultDiscussionCommentLimit = 50
)
