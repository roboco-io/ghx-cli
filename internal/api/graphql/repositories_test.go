package graphql

import (
	"testing"

	gql "github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryQueries(t *testing.T) {
	t.Run("RepositoryQuery structure", func(t *testing.T) {
		query := &RepositoryQuery{}
		assert.NotNil(t, query)
	})
}

func TestRepositoryVariableBuilders(t *testing.T) {
	t.Run("BuildRepositoryVariables creates proper variables", func(t *testing.T) {
		variables := BuildRepositoryVariables("octocat", "Hello-World")

		assert.NotNil(t, variables)
		assert.Equal(t, gql.String("octocat"), variables["owner"])
		assert.Equal(t, gql.String("Hello-World"), variables["name"])
	})
}

func TestParseRepositoryResponse(t *testing.T) {
	t.Run("ParseRepositoryResponse handles nil repository", func(t *testing.T) {
		resp := &RepositoryQuery{
			Repository: nil,
		}

		result, err := ParseRepositoryResponse(resp)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("ParseRepositoryResponse handles valid repository", func(t *testing.T) {
		description := "My first repository on GitHub!"
		resp := &RepositoryQuery{
			Repository: &RepositoryDetails{
				ID:          "R_123456789",
				Name:        "Hello-World",
				Owner:       RepositoryOwner{Login: "octocat", Type: "User"},
				IsPrivate:   false,
				Visibility:  "public",
				Description: &description,
				URL:         "https://github.com/octocat/Hello-World",
			},
		}

		result, err := ParseRepositoryResponse(resp)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "R_123456789", result.ID)
		assert.Equal(t, "Hello-World", result.Name)
		assert.Equal(t, "octocat", result.Owner)
		assert.Equal(t, "User", result.OwnerType)
		assert.False(t, result.IsPrivate)
		assert.Equal(t, "public", result.Visibility)
		assert.Equal(t, "My first repository on GitHub!", *result.Description)
		assert.Equal(t, "https://github.com/octocat/Hello-World", result.URL)
	})

	t.Run("ParseRepositoryResponse handles repository with nil description", func(t *testing.T) {
		resp := &RepositoryQuery{
			Repository: &RepositoryDetails{
				ID:          "R_123456789",
				Name:        "Hello-World",
				Owner:       RepositoryOwner{Login: "octocat", Type: "Organization"},
				IsPrivate:   true,
				Visibility:  "private",
				Description: nil,
				URL:         "https://github.com/octocat/Hello-World",
			},
		}

		result, err := ParseRepositoryResponse(resp)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "R_123456789", result.ID)
		assert.Equal(t, "Hello-World", result.Name)
		assert.Equal(t, "octocat", result.Owner)
		assert.Equal(t, "Organization", result.OwnerType)
		assert.True(t, result.IsPrivate)
		assert.Equal(t, "private", result.Visibility)
		assert.Nil(t, result.Description)
		assert.Equal(t, "https://github.com/octocat/Hello-World", result.URL)
	})
}

func TestRepositoryStructures(t *testing.T) {
	t.Run("RepositoryDetails structure validation", func(t *testing.T) {
		repo := &RepositoryDetails{}
		assert.NotNil(t, repo)
	})

	t.Run("RepositoryOwner structure validation", func(t *testing.T) {
		owner := &RepositoryOwner{}
		assert.NotNil(t, owner)
	})

	t.Run("RepositoryInfo structure validation", func(t *testing.T) {
		info := &RepositoryInfo{}
		assert.NotNil(t, info)
	})
}
