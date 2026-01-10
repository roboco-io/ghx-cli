package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/roboco-io/ghx-cli/internal/api"
)

func TestProjectService(t *testing.T) {
	t.Run("NewProjectService creates new service", func(t *testing.T) {
		client := api.NewClient("test-token")
		service := NewProjectService(client)

		assert.NotNil(t, service)
		assert.IsType(t, &ProjectService{}, service)
	})

	t.Run("ListUserProjects with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewProjectService(client)

		opts := ListUserProjectsOptions{
			Login: "testuser",
			First: 10,
		}

		ctx := context.Background()
		projects, err := service.ListUserProjects(ctx, opts)

		// With invalid token, this should return an error
		assert.Error(t, err)
		assert.Nil(t, projects)
		assert.Contains(t, err.Error(), "failed to list user projects")
	})

	t.Run("ListOrgProjects with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewProjectService(client)

		opts := ListOrgProjectsOptions{
			Login: "testorg",
			First: 10,
		}

		ctx := context.Background()
		projects, err := service.ListOrgProjects(ctx, opts)

		// With invalid token, this should return an error
		assert.Error(t, err)
		assert.Nil(t, projects)
		assert.Contains(t, err.Error(), "failed to list organization projects")
	})

	t.Run("GetProject with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewProjectService(client)

		ctx := context.Background()
		project, err := service.GetProject(ctx, "testowner", 1, false)

		// With invalid token, this should return an error
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "failed to get")
	})
}

func TestParseProjectReference(t *testing.T) {
	t.Run("Valid project reference", func(t *testing.T) {
		owner, number, err := ParseProjectReference("owner/123")

		assert.NoError(t, err)
		assert.Equal(t, "owner", owner)
		assert.Equal(t, 123, number)
	})

	t.Run("Valid project reference with org", func(t *testing.T) {
		owner, number, err := ParseProjectReference("my-org/456")

		assert.NoError(t, err)
		assert.Equal(t, "my-org", owner)
		assert.Equal(t, 456, number)
	})

	t.Run("Invalid project reference - no slash", func(t *testing.T) {
		owner, number, err := ParseProjectReference("invalid")

		assert.Error(t, err)
		assert.Empty(t, owner)
		assert.Zero(t, number)
		assert.Contains(t, err.Error(), "invalid project reference format")
	})

	t.Run("Invalid project reference - no number", func(t *testing.T) {
		owner, number, err := ParseProjectReference("owner/")

		assert.Error(t, err)
		assert.Empty(t, owner)
		assert.Zero(t, number)
		assert.Contains(t, err.Error(), "invalid project reference format")
	})

	t.Run("Invalid project reference - non-numeric", func(t *testing.T) {
		owner, number, err := ParseProjectReference("owner/abc")

		assert.Error(t, err)
		assert.Empty(t, owner)
		assert.Zero(t, number)
		assert.Contains(t, err.Error(), "invalid project number")
	})
}

func TestFormatProjectReference(t *testing.T) {
	t.Run("Format valid reference", func(t *testing.T) {
		ref := FormatProjectReference("owner", 123)
		assert.Equal(t, "owner/123", ref)
	})

	t.Run("Format org reference", func(t *testing.T) {
		ref := FormatProjectReference("my-org", 456)
		assert.Equal(t, "my-org/456", ref)
	})
}

func TestProjectServiceMethods(t *testing.T) {
	client := api.NewClient("test-token")
	service := NewProjectService(client)
	ctx := context.Background()

	t.Run("CreateProject with invalid token fails", func(t *testing.T) {
		input := CreateProjectInput{
			OwnerID: "test-owner-id",
			Title:   "Test Project",
		}

		project, err := service.CreateProject(ctx, &input)
		assert.Error(t, err)
		assert.Nil(t, project)
	})

	t.Run("UpdateProject with invalid token fails", func(t *testing.T) {
		title := "Updated Title"
		input := UpdateProjectInput{
			ProjectID: "test-project-id",
			Title:     &title,
		}

		project, err := service.UpdateProject(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, project)
	})

	t.Run("DeleteProject with invalid token fails", func(t *testing.T) {
		err := service.DeleteProject(ctx, "test-project-id")
		assert.Error(t, err)
	})

	t.Run("AddItem with invalid token fails", func(t *testing.T) {
		input := AddItemInput{
			ProjectID: "test-project-id",
			ContentID: "test-content-id",
		}

		item, err := service.AddItem(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, item)
	})

	t.Run("UpdateItemField with invalid token fails", func(t *testing.T) {
		input := UpdateItemFieldInput{
			ProjectID: "test-project-id",
			ItemID:    "test-item-id",
			FieldID:   "test-field-id",
			Value:     "test-value",
		}

		item, err := service.UpdateItemField(ctx, input)
		assert.Error(t, err)
		assert.Nil(t, item)
	})

	t.Run("RemoveItem with invalid token fails", func(t *testing.T) {
		input := RemoveItemInput{
			ProjectID: "test-project-id",
			ItemID:    "test-item-id",
		}

		err := service.RemoveItem(ctx, input)
		assert.Error(t, err)
	})
}
