package graphql

import "time"

// =============================================================================
// DATA TYPES
// =============================================================================

// Discussion represents a GitHub Discussion
type Discussion struct {
	CreatedAt      time.Time          `graphql:"createdAt"`
	UpdatedAt      time.Time          `graphql:"updatedAt"`
	ClosedAt       *time.Time         `graphql:"closedAt"`
	AnswerChosenAt *time.Time         `graphql:"answerChosenAt"`
	Author         DiscussionActor    `graphql:"author"`
	Category       DiscussionCategory `graphql:"category"`
	Answer         *DiscussionComment `graphql:"answer"`
	Comments       struct {
		TotalCount int                 `graphql:"totalCount"`
		Nodes      []DiscussionComment `graphql:"nodes"`
	} `graphql:"comments(first: $commentFirst)"`
	Labels struct {
		Nodes []DiscussionLabel `graphql:"nodes"`
	} `graphql:"labels(first: 10)"`
	ID          string `graphql:"id"`
	Title       string `graphql:"title"`
	Body        string `graphql:"body"`
	BodyHTML    string `graphql:"bodyHTML"`
	URL         string `graphql:"url"`
	Number      int    `graphql:"number"`
	UpvoteCount int    `graphql:"upvoteCount"`
	Locked      bool   `graphql:"locked"`
	Closed      bool   `graphql:"closed"`
}

// DiscussionCategory represents a discussion category
type DiscussionCategory struct {
	CreatedAt    time.Time `graphql:"createdAt"`
	ID           string    `graphql:"id"`
	Name         string    `graphql:"name"`
	Slug         string    `graphql:"slug"`
	Description  string    `graphql:"description"`
	Emoji        string    `graphql:"emoji"`
	IsAnswerable bool      `graphql:"isAnswerable"`
}

// DiscussionComment represents a comment on a discussion
type DiscussionComment struct {
	CreatedAt   time.Time       `graphql:"createdAt"`
	UpdatedAt   time.Time       `graphql:"updatedAt"`
	Author      DiscussionActor `graphql:"author"`
	ID          string          `graphql:"id"`
	Body        string          `graphql:"body"`
	BodyHTML    string          `graphql:"bodyHTML"`
	UpvoteCount int             `graphql:"upvoteCount"`
	IsAnswer    bool            `graphql:"isAnswer"`
}

// DiscussionActor represents the author (User or Bot)
type DiscussionActor struct {
	Login     string `graphql:"login"`
	AvatarURL string `graphql:"avatarUrl"`
	Type      string `graphql:"__typename"`
}

// DiscussionLabel represents a GitHub label
type DiscussionLabel struct {
	ID          string `graphql:"id"`
	Name        string `graphql:"name"`
	Color       string `graphql:"color"`
	Description string `graphql:"description"`
}

// =============================================================================
// ENUM TYPES
// =============================================================================

// DiscussionCloseReason represents the reason for closing a discussion
type DiscussionCloseReason string

const (
	DiscussionCloseReasonResolved  DiscussionCloseReason = "RESOLVED"
	DiscussionCloseReasonOutdated  DiscussionCloseReason = "OUTDATED"
	DiscussionCloseReasonDuplicate DiscussionCloseReason = "DUPLICATE"
)

// DiscussionLockReason represents the reason for locking
type DiscussionLockReason string

const (
	DiscussionLockReasonOffTopic  DiscussionLockReason = "OFF_TOPIC"
	DiscussionLockReasonResolved  DiscussionLockReason = "RESOLVED"
	DiscussionLockReasonSpam      DiscussionLockReason = "SPAM"
	DiscussionLockReasonTooHeated DiscussionLockReason = "TOO_HEATED"
)

// DiscussionOrderField represents the field to order discussions by
type DiscussionOrderField string

const (
	DiscussionOrderFieldCreatedAt DiscussionOrderField = "CREATED_AT"
	DiscussionOrderFieldUpdatedAt DiscussionOrderField = "UPDATED_AT"
)

// DiscussionState represents the state of a discussion
type DiscussionState string

const (
	DiscussionStateOpen   DiscussionState = "OPEN"
	DiscussionStateClosed DiscussionState = "CLOSED"
)

// =============================================================================
// QUERIES
// =============================================================================

// ListDiscussionsQuery lists discussions for a repository
type ListDiscussionsQuery struct {
	Repository struct {
		Discussions struct {
			PageInfo   PageInfo     `graphql:"pageInfo"`
			Nodes      []Discussion `graphql:"nodes"`
			TotalCount int          `graphql:"totalCount"`
		} `graphql:"discussions(first: $first, after: $after, categoryId: $categoryId, answered: $answered)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// GetDiscussionQuery gets a specific discussion by number
type GetDiscussionQuery struct {
	Repository struct {
		Discussion Discussion `graphql:"discussion(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// ListDiscussionCategoriesQuery lists discussion categories for a repository
type ListDiscussionCategoriesQuery struct {
	Repository struct {
		DiscussionCategories struct {
			Nodes []DiscussionCategory `graphql:"nodes"`
		} `graphql:"discussionCategories(first: 50)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// GetDiscussionCategoryQuery gets a specific category by slug
type GetDiscussionCategoryQuery struct {
	Repository struct {
		DiscussionCategory *DiscussionCategory `graphql:"discussionCategory(slug: $slug)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// =============================================================================
// MUTATIONS
// =============================================================================

// CreateDiscussionMutation creates a new discussion
type CreateDiscussionMutation struct {
	CreateDiscussion struct {
		Discussion Discussion `graphql:"discussion"`
	} `graphql:"createDiscussion(input: $input)"`
}

// UpdateDiscussionMutation updates a discussion
type UpdateDiscussionMutation struct {
	UpdateDiscussion struct {
		Discussion Discussion `graphql:"discussion"`
	} `graphql:"updateDiscussion(input: $input)"`
}

// DeleteDiscussionMutation deletes a discussion
type DeleteDiscussionMutation struct {
	DeleteDiscussion struct {
		Discussion struct {
			ID string `graphql:"id"`
		} `graphql:"discussion"`
	} `graphql:"deleteDiscussion(input: $input)"`
}

// CloseDiscussionMutation closes a discussion
type CloseDiscussionMutation struct {
	CloseDiscussion struct {
		Discussion Discussion `graphql:"discussion"`
	} `graphql:"closeDiscussion(input: $input)"`
}

// ReopenDiscussionMutation reopens a discussion
type ReopenDiscussionMutation struct {
	ReopenDiscussion struct {
		Discussion Discussion `graphql:"discussion"`
	} `graphql:"reopenDiscussion(input: $input)"`
}

// LockDiscussionMutation locks a discussion (uses generic LockLockable)
type LockDiscussionMutation struct {
	LockLockable struct {
		LockedRecord struct {
			Locked bool `graphql:"locked"`
		} `graphql:"lockedRecord"`
	} `graphql:"lockLockable(input: $input)"`
}

// UnlockDiscussionMutation unlocks a discussion (uses generic UnlockLockable)
type UnlockDiscussionMutation struct {
	UnlockLockable struct {
		UnlockedRecord struct {
			Locked bool `graphql:"locked"`
		} `graphql:"unlockedRecord"`
	} `graphql:"unlockLockable(input: $input)"`
}

// AddDiscussionCommentMutation adds a comment to a discussion
type AddDiscussionCommentMutation struct {
	AddDiscussionComment struct {
		Comment DiscussionComment `graphql:"comment"`
	} `graphql:"addDiscussionComment(input: $input)"`
}

// UpdateDiscussionCommentMutation updates a discussion comment
type UpdateDiscussionCommentMutation struct {
	UpdateDiscussionComment struct {
		Comment DiscussionComment `graphql:"comment"`
	} `graphql:"updateDiscussionComment(input: $input)"`
}

// DeleteDiscussionCommentMutation deletes a discussion comment
type DeleteDiscussionCommentMutation struct {
	DeleteDiscussionComment struct {
		Comment struct {
			ID string `graphql:"id"`
		} `graphql:"comment"`
	} `graphql:"deleteDiscussionComment(input: $input)"`
}

// MarkDiscussionCommentAsAnswerMutation marks a comment as the answer
type MarkDiscussionCommentAsAnswerMutation struct {
	MarkDiscussionCommentAsAnswer struct {
		Discussion Discussion `graphql:"discussion"`
	} `graphql:"markDiscussionCommentAsAnswer(input: $input)"`
}

// UnmarkDiscussionCommentAsAnswerMutation unmarks a comment as the answer
type UnmarkDiscussionCommentAsAnswerMutation struct {
	UnmarkDiscussionCommentAsAnswer struct {
		Discussion Discussion `graphql:"discussion"`
	} `graphql:"unmarkDiscussionCommentAsAnswer(input: $input)"`
}

// =============================================================================
// INPUT TYPES
// =============================================================================

// CreateDiscussionInput represents input for creating a discussion
type CreateDiscussionInput struct {
	RepositoryID string `json:"repositoryId"`
	CategoryID   string `json:"categoryId"`
	Title        string `json:"title"`
	Body         string `json:"body"`
}

// UpdateDiscussionInput represents input for updating a discussion
type UpdateDiscussionInput struct {
	Title        *string `json:"title,omitempty"`
	Body         *string `json:"body,omitempty"`
	CategoryID   *string `json:"categoryId,omitempty"`
	DiscussionID string  `json:"discussionId"`
}

// DeleteDiscussionInput represents input for deleting a discussion
type DeleteDiscussionInput struct {
	ID string `json:"id"`
}

// CloseDiscussionInput represents input for closing a discussion
type CloseDiscussionInput struct {
	Reason       DiscussionCloseReason `json:"reason"`
	DiscussionID string                `json:"discussionId"`
}

// ReopenDiscussionInput represents input for reopening a discussion
type ReopenDiscussionInput struct {
	DiscussionID string `json:"discussionId"`
}

// LockDiscussionInput represents input for locking a discussion
type LockDiscussionInput struct {
	LockReason *DiscussionLockReason `json:"lockReason,omitempty"`
	LockableID string                `json:"lockableId"`
}

// UnlockDiscussionInput represents input for unlocking a discussion
type UnlockDiscussionInput struct {
	LockableID string `json:"lockableId"`
}

// AddDiscussionCommentInput represents input for adding a comment
type AddDiscussionCommentInput struct {
	ReplyToID    *string `json:"replyToId,omitempty"`
	DiscussionID string  `json:"discussionId"`
	Body         string  `json:"body"`
}

// UpdateDiscussionCommentInput represents input for updating a comment
type UpdateDiscussionCommentInput struct {
	CommentID string `json:"commentId"`
	Body      string `json:"body"`
}

// DeleteDiscussionCommentInput represents input for deleting a comment
type DeleteDiscussionCommentInput struct {
	ID string `json:"id"`
}

// MarkDiscussionCommentAsAnswerInput represents input for marking answer
type MarkDiscussionCommentAsAnswerInput struct {
	ID string `json:"id"`
}

// UnmarkDiscussionCommentAsAnswerInput represents input for unmarking answer
type UnmarkDiscussionCommentAsAnswerInput struct {
	ID string `json:"id"`
}

// =============================================================================
// VARIABLE BUILDERS
// =============================================================================

// BuildListDiscussionsVariables builds variables for listing discussions
func BuildListDiscussionsVariables(owner, name string, first int, after, categoryID *string, answered *bool) map[string]interface{} {
	vars := map[string]interface{}{
		"owner":        owner,
		"name":         name,
		"first":        first,
		"commentFirst": 0,
	}
	if after != nil {
		vars["after"] = *after
	}
	if categoryID != nil {
		vars["categoryId"] = *categoryID
	}
	if answered != nil {
		vars["answered"] = *answered
	}
	return vars
}

// BuildGetDiscussionVariables builds variables for getting a discussion
func BuildGetDiscussionVariables(owner, name string, number, commentFirst int) map[string]interface{} {
	return map[string]interface{}{
		"owner":        owner,
		"name":         name,
		"number":       number,
		"commentFirst": commentFirst,
	}
}

// BuildListDiscussionCategoriesVariables builds variables for listing categories
func BuildListDiscussionCategoriesVariables(owner, name string) map[string]interface{} {
	return map[string]interface{}{
		"owner": owner,
		"name":  name,
	}
}

// BuildGetDiscussionCategoryVariables builds variables for getting a category
func BuildGetDiscussionCategoryVariables(owner, name, slug string) map[string]interface{} {
	return map[string]interface{}{
		"owner": owner,
		"name":  name,
		"slug":  slug,
	}
}

// BuildCreateDiscussionVariables builds variables for creating a discussion
func BuildCreateDiscussionVariables(input *CreateDiscussionInput) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"repositoryId": input.RepositoryID,
			"categoryId":   input.CategoryID,
			"title":        input.Title,
			"body":         input.Body,
		},
		"commentFirst": 0,
	}
}

// BuildUpdateDiscussionVariables builds variables for updating a discussion
func BuildUpdateDiscussionVariables(input *UpdateDiscussionInput) map[string]interface{} {
	inputMap := map[string]interface{}{
		"discussionId": input.DiscussionID,
	}
	if input.Title != nil {
		inputMap["title"] = *input.Title
	}
	if input.Body != nil {
		inputMap["body"] = *input.Body
	}
	if input.CategoryID != nil {
		inputMap["categoryId"] = *input.CategoryID
	}
	return map[string]interface{}{
		"input":        inputMap,
		"commentFirst": 0,
	}
}

// BuildDeleteDiscussionVariables builds variables for deleting a discussion
func BuildDeleteDiscussionVariables(id string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"id": id,
		},
	}
}

// BuildCloseDiscussionVariables builds variables for closing a discussion
func BuildCloseDiscussionVariables(discussionID string, reason DiscussionCloseReason) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"discussionId": discussionID,
			"reason":       reason,
		},
		"commentFirst": 0,
	}
}

// BuildReopenDiscussionVariables builds variables for reopening a discussion
func BuildReopenDiscussionVariables(discussionID string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"discussionId": discussionID,
		},
		"commentFirst": 0,
	}
}

// BuildLockDiscussionVariables builds variables for locking a discussion
func BuildLockDiscussionVariables(lockableID string, reason *DiscussionLockReason) map[string]interface{} {
	inputMap := map[string]interface{}{
		"lockableId": lockableID,
	}
	if reason != nil {
		inputMap["lockReason"] = *reason
	}
	return map[string]interface{}{
		"input": inputMap,
	}
}

// BuildUnlockDiscussionVariables builds variables for unlocking a discussion
func BuildUnlockDiscussionVariables(lockableID string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"lockableId": lockableID,
		},
	}
}

// BuildAddDiscussionCommentVariables builds variables for adding a comment
func BuildAddDiscussionCommentVariables(discussionID, body string, replyToID *string) map[string]interface{} {
	inputMap := map[string]interface{}{
		"discussionId": discussionID,
		"body":         body,
	}
	if replyToID != nil {
		inputMap["replyToId"] = *replyToID
	}
	return map[string]interface{}{
		"input": inputMap,
	}
}

// BuildUpdateDiscussionCommentVariables builds variables for updating a comment
func BuildUpdateDiscussionCommentVariables(commentID, body string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"commentId": commentID,
			"body":      body,
		},
	}
}

// BuildDeleteDiscussionCommentVariables builds variables for deleting a comment
func BuildDeleteDiscussionCommentVariables(id string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"id": id,
		},
	}
}

// BuildMarkAnswerVariables builds variables for marking an answer
func BuildMarkAnswerVariables(commentID string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"id": commentID,
		},
		"commentFirst": 0,
	}
}

// BuildUnmarkAnswerVariables builds variables for unmarking an answer
func BuildUnmarkAnswerVariables(commentID string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"id": commentID,
		},
		"commentFirst": 0,
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// ValidCloseReasons returns all valid close reasons
func ValidCloseReasons() []string {
	return []string{
		string(DiscussionCloseReasonResolved),
		string(DiscussionCloseReasonOutdated),
		string(DiscussionCloseReasonDuplicate),
	}
}

// ValidLockReasons returns all valid lock reasons
func ValidLockReasons() []string {
	return []string{
		string(DiscussionLockReasonOffTopic),
		string(DiscussionLockReasonResolved),
		string(DiscussionLockReasonSpam),
		string(DiscussionLockReasonTooHeated),
	}
}

// FormatDiscussionState formats discussion state for display
func FormatDiscussionState(closed bool) string {
	if closed {
		return "Closed"
	}
	return "Open"
}

// FormatCloseReason formats close reason for display
func FormatCloseReason(reason DiscussionCloseReason) string {
	switch reason {
	case DiscussionCloseReasonResolved:
		return "Resolved"
	case DiscussionCloseReasonOutdated:
		return "Outdated"
	case DiscussionCloseReasonDuplicate:
		return "Duplicate"
	default:
		return string(reason)
	}
}

// FormatLockReason formats lock reason for display
func FormatLockReason(reason DiscussionLockReason) string {
	switch reason {
	case DiscussionLockReasonOffTopic:
		return "Off-topic"
	case DiscussionLockReasonResolved:
		return "Resolved"
	case DiscussionLockReasonSpam:
		return "Spam"
	case DiscussionLockReasonTooHeated:
		return "Too heated"
	default:
		return string(reason)
	}
}
