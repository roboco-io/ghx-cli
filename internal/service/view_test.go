package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

func TestViewService(t *testing.T) {
	t.Run("NewViewService creates new service", func(t *testing.T) {
		client := api.NewClient("test-token")
		service := NewViewService(client)

		assert.NotNil(t, service)
		assert.IsType(t, &ViewService{}, service)
	})

	t.Run("CreateView with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		input := CreateViewInput{
			ProjectID: "test-project-id",
			Name:      "Test View",
			Layout:    graphql.ProjectV2ViewLayoutTable,
		}

		view, err := service.CreateView(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, view)
		assert.Contains(t, err.Error(), "failed to create view")
	})

	t.Run("UpdateView with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		newName := "Updated Test View"
		input := UpdateViewInput{
			ViewID: "test-view-id",
			Name:   &newName,
		}

		view, err := service.UpdateView(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, view)
		assert.Contains(t, err.Error(), "failed to update view")
	})

	t.Run("DeleteView with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		input := DeleteViewInput{
			ViewID: "test-view-id",
		}

		err := service.DeleteView(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete view")
	})

	t.Run("CopyView with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		input := CopyViewInput{
			ProjectID: "test-project-id",
			ViewID:    "test-view-id",
			Name:      "Copied View",
		}

		view, err := service.CopyView(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, view)
		assert.Contains(t, err.Error(), "failed to copy view")
	})

	t.Run("UpdateViewSort with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		sortByID := "test-field-id"
		input := UpdateViewSortInput{
			ViewID:    "test-view-id",
			SortByID:  &sortByID,
			Direction: graphql.ProjectV2ViewSortDirectionASC,
		}

		err := service.UpdateViewSort(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update view sort")
	})

	t.Run("UpdateViewGroup with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		groupByID := "test-field-id"
		input := UpdateViewGroupInput{
			ViewID:    "test-view-id",
			GroupByID: &groupByID,
			Direction: graphql.ProjectV2ViewSortDirectionASC,
		}

		err := service.UpdateViewGroup(ctx, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update view group")
	})

	t.Run("GetProjectViews with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		views, err := service.GetProjectViews(ctx, "test-project-id")

		assert.Error(t, err)
		assert.Nil(t, views)
		assert.Contains(t, err.Error(), "failed to get project views")
	})

	t.Run("GetView with invalid token returns error", func(t *testing.T) {
		client := api.NewClient("invalid-token")
		service := NewViewService(client)

		ctx := context.Background()
		view, err := service.GetView(ctx, "test-view-id")

		assert.Error(t, err)
		assert.Nil(t, view)
		assert.Contains(t, err.Error(), "failed to get view")
	})
}

func TestViewValidation(t *testing.T) {
	t.Run("ValidateViewName accepts valid names", func(t *testing.T) {
		err := ValidateViewName("Test View")
		assert.NoError(t, err)

		err = ValidateViewName("Dashboard View")
		assert.NoError(t, err)
	})

	t.Run("ValidateViewName rejects empty names", func(t *testing.T) {
		err := ValidateViewName("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "view name cannot be empty")

		err = ValidateViewName("   ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "view name cannot be empty")
	})

	t.Run("ValidateViewName rejects long names", func(t *testing.T) {
		longName := strings.Repeat("a", 101)
		err := ValidateViewName(longName)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot exceed 100 characters")
	})

	t.Run("ValidateViewLayout accepts valid layouts", func(t *testing.T) {
		layout, err := ValidateViewLayout("table")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewLayoutTable, layout)

		layout, err = ValidateViewLayout("BOARD")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewLayoutBoard, layout)

		layout, err = ValidateViewLayout("Roadmap")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewLayoutRoadmap, layout)

		layout, err = ValidateViewLayout("TABLE_VIEW")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewLayoutTable, layout)
	})

	t.Run("ValidateViewLayout rejects invalid layouts", func(t *testing.T) {
		_, err := ValidateViewLayout("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid view layout: invalid")
	})

	t.Run("ValidateSortDirection accepts valid directions", func(t *testing.T) {
		direction, err := ValidateSortDirection("asc")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewSortDirectionASC, direction)

		direction, err = ValidateSortDirection("DESC")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewSortDirectionDESC, direction)

		direction, err = ValidateSortDirection("Ascending")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewSortDirectionASC, direction)

		direction, err = ValidateSortDirection("DESCENDING")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ProjectV2ViewSortDirectionDESC, direction)
	})

	t.Run("ValidateSortDirection rejects invalid directions", func(t *testing.T) {
		_, err := ValidateSortDirection("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid sort direction: invalid")
	})

	t.Run("NormalizeSortDirection converts to proper format", func(t *testing.T) {
		assert.Equal(t, "ASC", NormalizeSortDirection("asc"))
		assert.Equal(t, "DESC", NormalizeSortDirection("DESC"))
		assert.Equal(t, "ASC", NormalizeSortDirection("ascending"))
		assert.Equal(t, "DESC", NormalizeSortDirection("DESCENDING"))
	})
}

func TestViewFormatting(t *testing.T) {
	t.Run("FormatViewLayout formats correctly", func(t *testing.T) {
		assert.Equal(t, "Table", FormatViewLayout(graphql.ProjectV2ViewLayoutTable))
		assert.Equal(t, "Board", FormatViewLayout(graphql.ProjectV2ViewLayoutBoard))
		assert.Equal(t, "Roadmap", FormatViewLayout(graphql.ProjectV2ViewLayoutRoadmap))
	})

	t.Run("FormatSortDirection formats correctly", func(t *testing.T) {
		assert.Equal(t, "Ascending", FormatSortDirection(graphql.ProjectV2ViewSortDirectionASC))
		assert.Equal(t, "Descending", FormatSortDirection(graphql.ProjectV2ViewSortDirectionDESC))
	})
}

func TestViewInfo(t *testing.T) {
	t.Run("ViewInfo structure", func(t *testing.T) {
		info := ViewInfo{
			ID:        "view-id",
			Name:      "Test View",
			Layout:    graphql.ProjectV2ViewLayoutTable,
			Number:    1,
			ProjectID: "project-id",
		}

		assert.Equal(t, "view-id", info.ID)
		assert.Equal(t, "Test View", info.Name)
		assert.Equal(t, graphql.ProjectV2ViewLayoutTable, info.Layout)
		assert.Equal(t, 1, info.Number)
		assert.Equal(t, "project-id", info.ProjectID)
	})

	t.Run("ViewGroupByInfo structure", func(t *testing.T) {
		info := ViewGroupByInfo{
			FieldID:   "field-id",
			FieldName: "Status",
			Direction: graphql.ProjectV2ViewSortDirectionASC,
		}

		assert.Equal(t, "field-id", info.FieldID)
		assert.Equal(t, "Status", info.FieldName)
		assert.Equal(t, graphql.ProjectV2ViewSortDirectionASC, info.Direction)
	})

	t.Run("ViewSortByInfo structure", func(t *testing.T) {
		info := ViewSortByInfo{
			FieldID:   "field-id",
			FieldName: "Priority",
			Direction: graphql.ProjectV2ViewSortDirectionDESC,
		}

		assert.Equal(t, "field-id", info.FieldID)
		assert.Equal(t, "Priority", info.FieldName)
		assert.Equal(t, graphql.ProjectV2ViewSortDirectionDESC, info.Direction)
	})
}
