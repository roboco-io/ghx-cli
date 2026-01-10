package graphql

import (
	"testing"

	gql "github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
)

const (
	testFieldID = "test-field-id"
	fieldID     = "field-id"
)

func TestProjectV2ViewTypes(t *testing.T) {
	t.Run("ProjectV2ViewLayout constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2ViewLayout("TABLE_VIEW"), ProjectV2ViewLayoutTable)
		assert.Equal(t, ProjectV2ViewLayout("BOARD_VIEW"), ProjectV2ViewLayoutBoard)
		assert.Equal(t, ProjectV2ViewLayout("ROADMAP_VIEW"), ProjectV2ViewLayoutRoadmap)
	})

	t.Run("ProjectV2ViewSortDirection constants", func(t *testing.T) {
		assert.Equal(t, ProjectV2ViewSortDirection("ASC"), ProjectV2ViewSortDirectionASC)
		assert.Equal(t, ProjectV2ViewSortDirection("DESC"), ProjectV2ViewSortDirectionDESC)
	})
}

func TestViewVariableBuilders(t *testing.T) {
	t.Run("BuildCreateViewVariables creates proper variables", func(t *testing.T) {
		input := &CreateViewInput{
			ProjectID: gql.ID("test-project-id"),
			Name:      gql.String("Test View"),
			Layout:    ProjectV2ViewLayoutTable,
		}

		variables := BuildCreateViewVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateViewVariables creates proper variables", func(t *testing.T) {
		name := gql.String("Updated View")
		filter := gql.String("status:todo")
		input := &UpdateViewInput{
			ViewID: gql.ID("test-view-id"),
			Name:   &name,
			Filter: &filter,
		}

		variables := BuildUpdateViewVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateViewVariables with minimal input", func(t *testing.T) {
		input := &UpdateViewInput{
			ViewID: gql.ID("test-view-id"),
		}

		variables := BuildUpdateViewVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildDeleteViewVariables creates proper variables", func(t *testing.T) {
		input := &DeleteViewInput{
			ViewID: gql.ID("test-view-id"),
		}

		variables := BuildDeleteViewVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildCopyViewVariables creates proper variables", func(t *testing.T) {
		input := &CopyViewInput{
			ProjectID: gql.ID("test-project-id"),
			ViewID:    gql.ID("test-view-id"),
			Name:      gql.String("Copied View"),
		}

		variables := BuildCopyViewVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateViewSortByVariables creates proper variables", func(t *testing.T) {
		sortByID := gql.ID(testFieldID)
		input := &UpdateViewSortByInput{
			ViewID:    gql.ID("test-view-id"),
			SortByID:  &sortByID,
			Direction: ProjectV2ViewSortDirectionASC,
		}

		variables := BuildUpdateViewSortByVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateViewSortByVariables without sortById", func(t *testing.T) {
		input := &UpdateViewSortByInput{
			ViewID:    gql.ID("test-view-id"),
			Direction: ProjectV2ViewSortDirectionDESC,
		}

		variables := BuildUpdateViewSortByVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateViewGroupByVariables creates proper variables", func(t *testing.T) {
		groupByID := gql.ID(testFieldID)
		input := &UpdateViewGroupByInput{
			ViewID:    gql.ID("test-view-id"),
			GroupByID: &groupByID,
			Direction: ProjectV2ViewSortDirectionASC,
		}

		variables := BuildUpdateViewGroupByVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildUpdateViewGroupByVariables without groupById", func(t *testing.T) {
		input := &UpdateViewGroupByInput{
			ViewID:    gql.ID("test-view-id"),
			Direction: ProjectV2ViewSortDirectionDESC,
		}

		variables := BuildUpdateViewGroupByVariables(input)

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "input")
	})

	t.Run("BuildGetProjectViewsVariables creates proper variables", func(t *testing.T) {
		variables := BuildGetProjectViewsVariables("test-project-id")

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "projectId")
	})

	t.Run("BuildGetViewVariables creates proper variables", func(t *testing.T) {
		variables := BuildGetViewVariables("test-view-id")

		assert.NotNil(t, variables)
		assert.Contains(t, variables, "viewId")
	})
}

func TestViewHelperFunctions(t *testing.T) {
	t.Run("ValidViewLayouts returns all valid layouts", func(t *testing.T) {
		layouts := ValidViewLayouts()
		expected := []string{
			string(ProjectV2ViewLayoutTable),
			string(ProjectV2ViewLayoutBoard),
			string(ProjectV2ViewLayoutRoadmap),
		}

		assert.Equal(t, expected, layouts)
		assert.Len(t, layouts, 3)
	})

	t.Run("ValidSortDirections returns all valid directions", func(t *testing.T) {
		directions := ValidSortDirections()
		expected := []string{
			string(ProjectV2ViewSortDirectionASC),
			string(ProjectV2ViewSortDirectionDESC),
		}

		assert.Equal(t, expected, directions)
		assert.Len(t, directions, 2)
	})

	t.Run("FormatViewLayout formats correctly", func(t *testing.T) {
		assert.Equal(t, "Table", FormatViewLayout(ProjectV2ViewLayoutTable))
		assert.Equal(t, "Board", FormatViewLayout(ProjectV2ViewLayoutBoard))
		assert.Equal(t, "Roadmap", FormatViewLayout(ProjectV2ViewLayoutRoadmap))
		assert.Equal(t, "UNKNOWN_VIEW", FormatViewLayout(ProjectV2ViewLayout("UNKNOWN_VIEW")))
	})

	t.Run("FormatSortDirection formats correctly", func(t *testing.T) {
		assert.Equal(t, "Ascending", FormatSortDirection(ProjectV2ViewSortDirectionASC))
		assert.Equal(t, "Descending", FormatSortDirection(ProjectV2ViewSortDirectionDESC))
		assert.Equal(t, "UNKNOWN", FormatSortDirection(ProjectV2ViewSortDirection("UNKNOWN")))
	})
}

func TestViewStructures(t *testing.T) {
	t.Run("ProjectV2View structure validation", func(t *testing.T) {
		view := ProjectV2View{
			ID:     "view-id",
			Name:   "Test View",
			Layout: ProjectV2ViewLayoutTable,
			Number: 1,
		}

		assert.Equal(t, "view-id", view.ID)
		assert.Equal(t, "Test View", view.Name)
		assert.Equal(t, ProjectV2ViewLayoutTable, view.Layout)
		assert.Equal(t, 1, view.Number)
	})

	t.Run("ProjectV2ViewGroupBy structure validation", func(t *testing.T) {
		groupBy := ProjectV2ViewGroupBy{
			Direction: ProjectV2ViewSortDirectionASC,
		}
		groupBy.Field.ID = fieldID
		groupBy.Field.Name = "Status"

		assert.Equal(t, fieldID, groupBy.Field.ID)
		assert.Equal(t, "Status", groupBy.Field.Name)
		assert.Equal(t, ProjectV2ViewSortDirectionASC, groupBy.Direction)
	})

	t.Run("ProjectV2ViewSortBy structure validation", func(t *testing.T) {
		sortBy := ProjectV2ViewSortBy{
			Direction: ProjectV2ViewSortDirectionDESC,
		}
		sortBy.Field.ID = fieldID
		sortBy.Field.Name = "Priority"

		assert.Equal(t, fieldID, sortBy.Field.ID)
		assert.Equal(t, "Priority", sortBy.Field.Name)
		assert.Equal(t, ProjectV2ViewSortDirectionDESC, sortBy.Direction)
	})

	t.Run("ProjectV2ViewColumn structure validation", func(t *testing.T) {
		column := ProjectV2ViewColumn{
			ID:       "column-id",
			Name:     "Column Name",
			Width:    200,
			IsHidden: false,
		}
		column.Field.ID = fieldID
		column.Field.Name = "Field Name"

		assert.Equal(t, "column-id", column.ID)
		assert.Equal(t, "Column Name", column.Name)
		assert.Equal(t, 200, column.Width)
		assert.False(t, column.IsHidden)
		assert.Equal(t, fieldID, column.Field.ID)
		assert.Equal(t, "Field Name", column.Field.Name)
	})
}
