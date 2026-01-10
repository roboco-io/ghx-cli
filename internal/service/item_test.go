package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/roboco-io/ghx-cli/internal/api"
)

func TestItemService(t *testing.T) {
	t.Run("NewItemService creates new service", func(t *testing.T) {
		client := api.NewClient("test-token")
		service := NewItemService(client)

		assert.NotNil(t, service)
		assert.IsType(t, &ItemService{}, service)
	})

	t.Run("GetIssue with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewItemService(client)

		ctx := context.Background()
		issue, err := service.GetIssue(ctx, "owner", "repo", 123)

		assert.Error(t, err)
		assert.Nil(t, issue)
		assert.Contains(t, err.Error(), "failed to get issue")
	})

	t.Run("GetPullRequest with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewItemService(client)

		ctx := context.Background()
		pr, err := service.GetPullRequest(ctx, "owner", "repo", 123)

		assert.Error(t, err)
		assert.Nil(t, pr)
		assert.Contains(t, err.Error(), "failed to get pull request")
	})

	t.Run("SearchIssues with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewItemService(client)

		ctx := context.Background()
		issues, err := service.SearchIssues(ctx, "is:issue is:open", 10)

		assert.Error(t, err)
		assert.Nil(t, issues)
		assert.Contains(t, err.Error(), "failed to search issues")
	})

	t.Run("SearchPullRequests with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewItemService(client)

		ctx := context.Background()
		prs, err := service.SearchPullRequests(ctx, "is:pr is:open", 10)

		assert.Error(t, err)
		assert.Nil(t, prs)
		assert.Contains(t, err.Error(), "failed to search pull requests")
	})
}

func TestParseItemReference(t *testing.T) {
	t.Run("Parse owner/repo#123 format", func(t *testing.T) {
		owner, repo, number, err := ParseItemReference("octocat/Hello-World#123")

		assert.NoError(t, err)
		assert.Equal(t, "octocat", owner)
		assert.Equal(t, "Hello-World", repo)
		assert.Equal(t, 123, number)
	})

	t.Run("Parse GitHub issue URL", func(t *testing.T) {
		owner, repo, number, err := ParseItemReference("https://github.com/octocat/Hello-World/issues/123")

		assert.NoError(t, err)
		assert.Equal(t, "octocat", owner)
		assert.Equal(t, "Hello-World", repo)
		assert.Equal(t, 123, number)
	})

	t.Run("Parse GitHub PR URL", func(t *testing.T) {
		owner, repo, number, err := ParseItemReference("https://github.com/octocat/Hello-World/pull/456")

		assert.NoError(t, err)
		assert.Equal(t, "octocat", owner)
		assert.Equal(t, "Hello-World", repo)
		assert.Equal(t, 456, number)
	})

	t.Run("Invalid format returns error", func(t *testing.T) {
		_, _, _, err := ParseItemReference("invalid-format")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unrecognized item reference format")
	})

	t.Run("Invalid repository format returns error", func(t *testing.T) {
		_, _, _, err := ParseItemReference("invalid#123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid repository format")
	})

	t.Run("Invalid number returns error", func(t *testing.T) {
		_, _, _, err := ParseItemReference("owner/repo#abc")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid item number")
	})

	t.Run("Missing repository context returns error", func(t *testing.T) {
		_, _, _, err := ParseItemReference("#123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository context required")
	})
}

func TestFormatItemReference(t *testing.T) {
	t.Run("Format item reference correctly", func(t *testing.T) {
		ref := FormatItemReference("octocat", "Hello-World", 123)
		assert.Equal(t, "octocat/Hello-World#123", ref)
	})
}

func TestBuildSearchQuery(t *testing.T) {
	t.Run("Build basic issue search query", func(t *testing.T) {
		filters := SearchFilters{
			Type:  "issue",
			State: "open",
		}

		query := BuildSearchQuery(&filters)
		assert.Contains(t, query, "is:issue")
		assert.Contains(t, query, "is:open")
	})

	t.Run("Build PR search query with repository", func(t *testing.T) {
		filters := SearchFilters{
			Type:       "pr",
			State:      "closed",
			Repository: "octocat/Hello-World",
		}

		query := BuildSearchQuery(&filters)
		assert.Contains(t, query, "is:pr")
		assert.Contains(t, query, "is:closed")
		assert.Contains(t, query, "repo:octocat/Hello-World")
	})

	t.Run("Build complex search query", func(t *testing.T) {
		filters := SearchFilters{
			Type:       "issue",
			State:      "open",
			Repository: "octocat/Hello-World",
			Author:     "octocat",
			Assignee:   "hubot",
			Labels:     []string{"bug", "help-wanted"},
			Query:      "authentication error",
		}

		query := BuildSearchQuery(&filters)
		assert.Contains(t, query, "is:issue")
		assert.Contains(t, query, "is:open")
		assert.Contains(t, query, "repo:octocat/Hello-World")
		assert.Contains(t, query, "author:octocat")
		assert.Contains(t, query, "assignee:hubot")
		assert.Contains(t, query, "label:bug")
		assert.Contains(t, query, "label:help-wanted")
		assert.Contains(t, query, "authentication error")
	})

	t.Run("Build query with pullrequest alias", func(t *testing.T) {
		filters := SearchFilters{
			Type: "pullrequest",
		}

		query := BuildSearchQuery(&filters)
		assert.Contains(t, query, "is:pr")
	})

	t.Run("Build empty query", func(t *testing.T) {
		filters := SearchFilters{}

		query := BuildSearchQuery(&filters)
		assert.Equal(t, "", query)
	})
}
