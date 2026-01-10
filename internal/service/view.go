package service

import (
	"context"
	"fmt"
	"strings"

	gql "github.com/shurcooL/graphql"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

// ViewService handles view-related operations
type ViewService struct {
	client *api.Client
}

// NewViewService creates a new view service
func NewViewService(client *api.Client) *ViewService {
	return &ViewService{
		client: client,
	}
}

// ViewInfo represents simplified view information for display
type ViewInfo struct {
	Filter      *string
	ID          string
	Name        string
	Layout      graphql.ProjectV2ViewLayout
	ProjectID   string
	ProjectName string
	GroupBy     []ViewGroupByInfo
	SortBy      []ViewSortByInfo
	Number      int
}

// ViewGroupByInfo represents group by configuration information
type ViewGroupByInfo struct {
	FieldID   string
	FieldName string
	Direction graphql.ProjectV2ViewSortDirection
}

// ViewSortByInfo represents sort by configuration information
type ViewSortByInfo struct {
	FieldID   string
	FieldName string
	Direction graphql.ProjectV2ViewSortDirection
}

// CreateViewInput represents input for creating a view
type CreateViewInput struct {
	ProjectID string
	Name      string
	Layout    graphql.ProjectV2ViewLayout
}

// UpdateViewInput represents input for updating a view
type UpdateViewInput struct {
	Name   *string
	Filter *string
	ViewID string
}

// DeleteViewInput represents input for deleting a view
type DeleteViewInput struct {
	ViewID string
}

// CopyViewInput represents input for copying a view
type CopyViewInput struct {
	ProjectID string
	ViewID    string
	Name      string
}

// UpdateViewSortInput represents input for updating view sort configuration
type UpdateViewSortInput struct {
	ViewID    string
	SortByID  *string
	Direction graphql.ProjectV2ViewSortDirection
}

// UpdateViewGroupInput represents input for updating view group configuration
type UpdateViewGroupInput struct {
	ViewID    string
	GroupByID *string
	Direction graphql.ProjectV2ViewSortDirection
}

// CreateView creates a new project view
func (s *ViewService) CreateView(ctx context.Context, input CreateViewInput) (*graphql.ProjectV2View, error) {
	variables := graphql.BuildCreateViewVariables(&graphql.CreateViewInput{
		ProjectID: gql.ID(input.ProjectID),
		Name:      gql.String(input.Name),
		Layout:    input.Layout,
	})

	var mutation graphql.CreateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create view: %w", err)
	}

	return &mutation.CreateProjectV2View.ProjectV2View, nil
}

// UpdateView updates an existing project view
//
//nolint:dupl // Similar structure to UpdateProject but operates on different types
func (s *ViewService) UpdateView(ctx context.Context, input UpdateViewInput) (*graphql.ProjectV2View, error) {
	gqlInput := &graphql.UpdateViewInput{
		ViewID: gql.ID(input.ViewID),
	}
	if input.Name != nil {
		name := gql.String(*input.Name)
		gqlInput.Name = &name
	}
	if input.Filter != nil {
		filter := gql.String(*input.Filter)
		gqlInput.Filter = &filter
	}

	variables := graphql.BuildUpdateViewVariables(gqlInput)

	var mutation graphql.UpdateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update view: %w", err)
	}

	return &mutation.UpdateProjectV2View.ProjectV2View, nil
}

// DeleteView deletes a project view
func (s *ViewService) DeleteView(ctx context.Context, input DeleteViewInput) error {
	variables := graphql.BuildDeleteViewVariables(&graphql.DeleteViewInput{
		ViewID: gql.ID(input.ViewID),
	})

	var mutation graphql.DeleteProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete view: %w", err)
	}

	return nil
}

// CopyView creates a copy of an existing view
func (s *ViewService) CopyView(ctx context.Context, input CopyViewInput) (*graphql.ProjectV2View, error) {
	variables := graphql.BuildCopyViewVariables(&graphql.CopyViewInput{
		ProjectID: gql.ID(input.ProjectID),
		ViewID:    gql.ID(input.ViewID),
		Name:      gql.String(input.Name),
	})

	var mutation graphql.CopyProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to copy view: %w", err)
	}

	return &mutation.CopyProjectV2View.ProjectV2View, nil
}

// UpdateViewSort updates the sort configuration for a view
//
//nolint:dupl // Similar structure to UpdateViewGroup but operates on different field
func (s *ViewService) UpdateViewSort(ctx context.Context, input UpdateViewSortInput) error {
	gqlInput := &graphql.UpdateViewSortByInput{
		ViewID:    gql.ID(input.ViewID),
		Direction: input.Direction,
	}
	if input.SortByID != nil {
		id := gql.ID(*input.SortByID)
		gqlInput.SortByID = &id
	}

	variables := graphql.BuildUpdateViewSortByVariables(gqlInput)

	var mutation graphql.UpdateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to update view sort: %w", err)
	}

	return nil
}

// UpdateViewGroup updates the group configuration for a view
//
//nolint:dupl // Similar structure to UpdateViewSort but operates on different field
func (s *ViewService) UpdateViewGroup(ctx context.Context, input UpdateViewGroupInput) error {
	gqlInput := &graphql.UpdateViewGroupByInput{
		ViewID:    gql.ID(input.ViewID),
		Direction: input.Direction,
	}
	if input.GroupByID != nil {
		id := gql.ID(*input.GroupByID)
		gqlInput.GroupByID = &id
	}

	variables := graphql.BuildUpdateViewGroupByVariables(gqlInput)

	var mutation graphql.UpdateProjectViewMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to update view group: %w", err)
	}

	return nil
}

// GetProjectViews gets all views for a project
func (s *ViewService) GetProjectViews(ctx context.Context, projectID string) ([]ViewInfo, error) {
	variables := graphql.BuildGetProjectViewsVariables(projectID)

	var query graphql.GetProjectViewsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get project views: %w", err)
	}

	views := make([]ViewInfo, len(query.Node.ProjectV2.Views.Nodes))
	for i := range query.Node.ProjectV2.Views.Nodes {
		view := &query.Node.ProjectV2.Views.Nodes[i]
		groupBy := make([]ViewGroupByInfo, len(view.GroupBy))
		for j, gb := range view.GroupBy {
			groupBy[j] = ViewGroupByInfo{
				FieldID:   gb.Field.ID,
				FieldName: gb.Field.Name,
				Direction: gb.Direction,
			}
		}

		sortBy := make([]ViewSortByInfo, len(view.SortBy))
		for j, sb := range view.SortBy {
			sortBy[j] = ViewSortByInfo{
				FieldID:   sb.Field.ID,
				FieldName: sb.Field.Name,
				Direction: sb.Direction,
			}
		}

		views[i] = ViewInfo{
			ID:        view.ID,
			Name:      view.Name,
			Layout:    view.Layout,
			Number:    view.Number,
			Filter:    view.Filter,
			ProjectID: projectID,
			GroupBy:   groupBy,
			SortBy:    sortBy,
		}
	}

	return views, nil
}

// GetView gets a specific view by ID
func (s *ViewService) GetView(ctx context.Context, viewID string) (*ViewInfo, error) {
	variables := graphql.BuildGetViewVariables(viewID)

	var query graphql.GetProjectViewQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get view: %w", err)
	}

	view := query.Node.ProjectV2View

	groupBy := make([]ViewGroupByInfo, len(view.GroupBy))
	for i, gb := range view.GroupBy {
		groupBy[i] = ViewGroupByInfo{
			FieldID:   gb.Field.ID,
			FieldName: gb.Field.Name,
			Direction: gb.Direction,
		}
	}

	sortBy := make([]ViewSortByInfo, len(view.SortBy))
	for i, sb := range view.SortBy {
		sortBy[i] = ViewSortByInfo{
			FieldID:   sb.Field.ID,
			FieldName: sb.Field.Name,
			Direction: sb.Direction,
		}
	}

	viewInfo := &ViewInfo{
		ID:      view.ID,
		Name:    view.Name,
		Layout:  view.Layout,
		Number:  view.Number,
		Filter:  view.Filter,
		GroupBy: groupBy,
		SortBy:  sortBy,
	}

	return viewInfo, nil
}

// ValidateViewName validates a view name
func ValidateViewName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("view name cannot be empty")
	}
	if len(name) > maxViewNameLength {
		return fmt.Errorf("view name cannot exceed %d characters", maxViewNameLength)
	}
	return nil
}

// ValidateViewLayout validates a view layout
func ValidateViewLayout(layout string) (graphql.ProjectV2ViewLayout, error) {
	switch strings.ToUpper(layout) {
	case "TABLE", "TABLE_VIEW":
		return graphql.ProjectV2ViewLayoutTable, nil
	case "BOARD", "BOARD_VIEW":
		return graphql.ProjectV2ViewLayoutBoard, nil
	case "ROADMAP", "ROADMAP_VIEW":
		return graphql.ProjectV2ViewLayoutRoadmap, nil
	default:
		validLayouts := graphql.ValidViewLayouts()
		return "", fmt.Errorf("invalid view layout: %s (valid layouts: %s)", layout, strings.ToLower(strings.Join(validLayouts, ", ")))
	}
}

// ValidateSortDirection validates a sort direction
func ValidateSortDirection(direction string) (graphql.ProjectV2ViewSortDirection, error) {
	switch strings.ToUpper(direction) {
	case "ASC", "ASCENDING":
		return graphql.ProjectV2ViewSortDirectionASC, nil
	case "DESC", "DESCENDING":
		return graphql.ProjectV2ViewSortDirectionDESC, nil
	default:
		validDirections := graphql.ValidSortDirections()
		validStr := strings.ToLower(strings.Join(validDirections, ", "))
		return "", fmt.Errorf("invalid sort direction: %s (valid directions: %s)", direction, validStr)
	}
}

// NormalizeSortDirection normalizes a sort direction string to the proper format
func NormalizeSortDirection(direction string) string {
	switch strings.ToUpper(direction) {
	case "ASC", "ASCENDING":
		return string(graphql.ProjectV2ViewSortDirectionASC)
	case "DESC", "DESCENDING":
		return string(graphql.ProjectV2ViewSortDirectionDESC)
	default:
		return strings.ToUpper(direction)
	}
}

// FormatViewLayout formats view layout for display
func FormatViewLayout(layout graphql.ProjectV2ViewLayout) string {
	return graphql.FormatViewLayout(layout)
}

// FormatSortDirection formats sort direction for display
func FormatSortDirection(direction graphql.ProjectV2ViewSortDirection) string {
	return graphql.FormatSortDirection(direction)
}
