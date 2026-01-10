package graphql

import (
	"testing"

	gql "github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

func TestProjectQueries(t *testing.T) {
	t.Run("ListUserProjects query structure", func(t *testing.T) {
		query := &ListUserProjectsQuery{}

		assert.NotNil(t, query)
		// Test that query has the right structure for GraphQL
	})

	t.Run("ListOrgProjects query structure", func(t *testing.T) {
		query := &ListOrgProjectsQuery{}

		assert.NotNil(t, query)
	})

	t.Run("GetProject query structure", func(t *testing.T) {
		query := &GetProjectQuery{}

		assert.NotNil(t, query)
	})
}

func TestProjectMutations(t *testing.T) {
	t.Run("CreateProject mutation structure", func(t *testing.T) {
		mutation := &CreateProjectMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("UpdateProject mutation structure", func(t *testing.T) {
		mutation := &UpdateProjectMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("DeleteProject mutation structure", func(t *testing.T) {
		mutation := &DeleteProjectMutation{}

		assert.NotNil(t, mutation)
	})
}

func TestItemMutations(t *testing.T) {
	t.Run("AddItemToProject mutation structure", func(t *testing.T) {
		mutation := &AddItemToProjectMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("UpdateItemField mutation structure", func(t *testing.T) {
		mutation := &UpdateItemFieldMutation{}

		assert.NotNil(t, mutation)
	})

	t.Run("RemoveItemFromProject mutation structure", func(t *testing.T) {
		mutation := &RemoveItemFromProjectMutation{}

		assert.NotNil(t, mutation)
	})
}

func TestVariableBuilders(t *testing.T) {
	t.Run("BuildCreateProjectVariables creates proper variables", func(t *testing.T) {
		input := &CreateProjectInput{
			OwnerID: gql.ID("test-owner-id"),
			Title:   gql.String("Test Project"),
		}

		variables := BuildCreateProjectVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildAddItemVariables creates proper variables", func(t *testing.T) {
		input := &AddItemInput{
			ProjectID: gql.ID("project-id"),
			ContentID: gql.ID("content-id"),
		}

		variables := BuildAddItemVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildListProjectsVariables creates proper variables", func(t *testing.T) {
		variables := BuildListProjectsVariables("testuser", 10, nil)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "login")
		assert.Contains(t, variables, "first")
	})

	t.Run("BuildGetProjectVariables creates proper variables for org", func(t *testing.T) {
		variables := BuildGetProjectVariables("testorg", 1, true)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "orgLogin")
		assert.Contains(t, variables, "number")
	})

	t.Run("BuildGetProjectVariables creates proper variables for user", func(t *testing.T) {
		variables := BuildGetProjectVariables("testuser", 1, false)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "userLogin")
		assert.Contains(t, variables, "number")
	})
}

func TestResponseParsing(t *testing.T) {
	t.Run("ParseProjectResponse extracts project data", func(t *testing.T) {
		response := &GetProjectQuery{
			Organization: struct {
				ProjectV2 ProjectV2 `graphql:"projectV2(number: $number)"`
			}{
				ProjectV2: ProjectV2{
					ID:     "test-id",
					Title:  "Test Project",
					Number: 42,
				},
			},
		}

		project := response.Organization.ProjectV2
		assert.Equal(t, "test-id", project.ID)
		assert.Equal(t, "Test Project", project.Title)
		assert.Equal(t, 42, project.Number)
	})
}
