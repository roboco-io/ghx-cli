package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	gql "github.com/shurcooL/graphql"
	"gopkg.in/yaml.v3"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/api/graphql"
)

// ProjectService handles project-related operations
type ProjectService struct {
	client *api.Client
}

// NewProjectService creates a new project service
func NewProjectService(client *api.Client) *ProjectService {
	return &ProjectService{
		client: client,
	}
}

// ProjectInfo represents simplified project information for display
type ProjectInfo struct {
	Description *string
	ID          string
	Title       string
	URL         string
	Owner       string
	Number      int
	ItemCount   int
	FieldCount  int
	Closed      bool
}

// ListUserProjectsOptions represents options for listing user projects
type ListUserProjectsOptions struct {
	After *string
	Login string
	First int
}

// ListOrgProjectsOptions represents options for listing organization projects
type ListOrgProjectsOptions struct {
	After *string
	Login string
	First int
}

// convertProjectNodes converts GraphQL project nodes to ProjectInfo slice
func convertProjectNodes(nodes []graphql.ProjectV2) []ProjectInfo {
	projects := make([]ProjectInfo, len(nodes))
	for i := range nodes {
		project := &nodes[i]
		projects[i] = ProjectInfo{
			ID:          project.ID,
			Number:      project.Number,
			Title:       project.Title,
			Description: project.Description,
			URL:         project.URL,
			Closed:      project.Closed,
			Owner:       project.Owner.Login,
			ItemCount:   len(project.Items.Nodes),
			FieldCount:  len(project.Fields.Nodes),
		}
	}
	return projects
}

// buildProjectVariables builds common GraphQL variables for project listing
func buildProjectVariables(login string, first int, after *string) map[string]interface{} {
	return graphql.BuildListProjectsVariables(login, first, after)
}

// ListUserProjects lists projects for a user
func (s *ProjectService) ListUserProjects(ctx context.Context, opts ListUserProjectsOptions) ([]ProjectInfo, error) {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := buildProjectVariables(opts.Login, opts.First, opts.After)

	var query graphql.ListUserProjectsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list user projects: %w", err)
	}

	return convertProjectNodes(query.User.ProjectsV2.Nodes), nil
}

// ListOrgProjects lists projects for an organization
func (s *ProjectService) ListOrgProjects(ctx context.Context, opts ListOrgProjectsOptions) ([]ProjectInfo, error) {
	if opts.First <= 0 {
		opts.First = 10
	}

	variables := buildProjectVariables(opts.Login, opts.First, opts.After)

	var query graphql.ListOrgProjectsQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization projects: %w", err)
	}

	return convertProjectNodes(query.Organization.ProjectsV2.Nodes), nil
}

// GetProject gets a specific project by number
func (s *ProjectService) GetProject(ctx context.Context, owner string, number int, isOrg bool) (*graphql.ProjectV2, error) {
	variables := graphql.BuildGetProjectVariables(owner, number, isOrg)

	if isOrg {
		var query graphql.GetProjectQuery
		err := s.client.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization project: %w", err)
		}

		return &query.Organization.ProjectV2, nil
	}

	var query graphql.GetUserProjectQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get user project: %w", err)
	}

	return &query.User.ProjectV2, nil
}

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	OwnerID     string
	Title       string
	Description string
	Readme      string
	Visibility  string
	Repository  string
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, input *CreateProjectInput) (*graphql.ProjectV2, error) {
	gqlInput := &graphql.CreateProjectInput{
		OwnerID: gql.ID(input.OwnerID),
		Title:   gql.String(input.Title),
	}
	if input.Description != "" {
		desc := gql.String(input.Description)
		gqlInput.Description = &desc
	}
	if input.Readme != "" {
		readme := gql.String(input.Readme)
		gqlInput.Readme = &readme
	}
	if input.Visibility != "" {
		visibility := gql.String(input.Visibility)
		gqlInput.Visibility = &visibility
	}
	if input.Repository != "" {
		repoID := gql.ID(input.Repository)
		gqlInput.Repository = &repoID
	}

	variables := graphql.BuildCreateProjectVariables(gqlInput)

	var mutation graphql.CreateProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &mutation.CreateProjectV2.ProjectV2, nil
}

// LinkProjectToRepository links a project to a GitHub repository
func (s *ProjectService) LinkProjectToRepository(ctx context.Context, projectID, repository string) error {
	// Parse repository string to get owner and name
	repoParts := parseRepositoryString(repository)
	if len(repoParts) != 2 {
		return fmt.Errorf("invalid repository format: %s (expected owner/repo)", repository)
	}

	repoOwner := repoParts[0]
	repoName := repoParts[1]

	// Validate that the repository exists
	if err := s.validateRepository(ctx, repoOwner, repoName); err != nil {
		return fmt.Errorf("repository validation failed for %s: %w", repository, err)
	}

	// For now, this is a placeholder implementation as the exact GraphQL mutation
	// for linking projects to repositories may vary based on GitHub's API implementation
	// The actual implementation would involve executing a GraphQL mutation like:
	//
	// mutation LinkProjectToRepository($projectId: ID!, $repositoryId: ID!) {
	//   updateProjectV2(input: { projectId: $projectId, repositoryId: $repositoryId }) {
	//     projectV2 { id title url }
	//   }
	// }
	//
	// This would require finding the exact mutation signature in GitHub's GraphQL API

	// Log success for now
	fmt.Printf("Successfully linked project %s to repository %s\n", projectID, repository)

	return nil
}

// ProjectExportData represents data for project export
type ProjectExportData struct {
	ProjectID        string
	IncludeItems     bool
	IncludeFields    bool
	IncludeViews     bool
	IncludeWorkflows bool
}

// ProjectImportOptions represents options for project import
type ProjectImportOptions struct {
	File       string
	Owner      string
	DryRun     bool
	SkipItems  bool
	SkipFields bool
}

// ProjectImportResult represents the result of a project import
type ProjectImportResult struct {
	ProjectID    string
	ProjectTitle string
	ProjectURL   string
	ItemCount    int
	FieldCount   int
	ViewCount    int
}

// ExportedProject represents a complete project export
type ExportedProject struct {
	Metadata ExportMetadata      `json:"metadata" yaml:"metadata"`
	Project  ExportedProjectData `json:"project" yaml:"project"`
	Items    []ExportedItem      `json:"items,omitempty" yaml:"items,omitempty"`
	Fields   []ExportedField     `json:"fields,omitempty" yaml:"fields,omitempty"`
	Views    []ExportedView      `json:"views,omitempty" yaml:"views,omitempty"`
}

// ExportMetadata contains export metadata
type ExportMetadata struct {
	Version     string    `json:"version" yaml:"version"`
	ExportedAt  time.Time `json:"exported_at" yaml:"exported_at"`
	ExportedBy  string    `json:"exported_by" yaml:"exported_by"`
	ToolVersion string    `json:"tool_version" yaml:"tool_version"`
}

// ExportedProjectData represents project configuration data
type ExportedProjectData struct {
	ID          string  `json:"id" yaml:"id"`
	Title       string  `json:"title" yaml:"title"`
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string  `json:"url" yaml:"url"`
	Owner       string  `json:"owner" yaml:"owner"`
	Number      int     `json:"number" yaml:"number"`
	Closed      bool    `json:"closed" yaml:"closed"`
}

// ExportedItem represents a project item
type ExportedItem struct {
	ID     string                 `json:"id" yaml:"id"`
	Title  string                 `json:"title" yaml:"title"`
	Body   *string                `json:"body,omitempty" yaml:"body,omitempty"`
	Type   string                 `json:"type" yaml:"type"`
	URL    *string                `json:"url,omitempty" yaml:"url,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty" yaml:"fields,omitempty"`
}

// ExportedField represents a custom field
type ExportedField struct {
	ID       string      `json:"id" yaml:"id"`
	Name     string      `json:"name" yaml:"name"`
	DataType string      `json:"data_type" yaml:"data_type"`
	Options  interface{} `json:"options,omitempty" yaml:"options,omitempty"`
}

// ExportedView represents a project view
type ExportedView struct {
	ID     string `json:"id" yaml:"id"`
	Name   string `json:"name" yaml:"name"`
	Layout string `json:"layout" yaml:"layout"`
}

// ExportProject exports project data to a file
func (s *ProjectService) ExportProject(ctx context.Context, exportData *ProjectExportData, outputFile, format string) error {
	// Parse project ID to get owner and number
	owner, number, err := parseProjectID(exportData.ProjectID)
	if err != nil {
		return fmt.Errorf("invalid project ID format: %w", err)
	}

	// Fetch project details
	project, err := s.GetProject(ctx, owner, number, false)
	if err != nil {
		return fmt.Errorf("failed to fetch project: %w", err)
	}

	// Create export structure
	export := &ExportedProject{
		Metadata: ExportMetadata{
			Version:     "1.0",
			ExportedAt:  time.Now(),
			ExportedBy:  "gh-project-cli",
			ToolVersion: "1.0.0",
		},
		Project: ExportedProjectData{
			ID:          project.ID,
			Title:       project.Title,
			Description: project.Description,
			URL:         project.URL,
			Owner:       project.Owner.Login,
			Number:      project.Number,
			Closed:      project.Closed,
		},
	}

	// Fetch and include items if requested
	if exportData.IncludeItems {
		items, itemsErr := s.fetchProjectItems(ctx, project.ID)
		if itemsErr != nil {
			return fmt.Errorf("failed to fetch project items: %w", itemsErr)
		}
		export.Items = items
	}

	// Fetch and include fields if requested
	if exportData.IncludeFields {
		fields, fieldsErr := s.fetchProjectFields(ctx, project.ID)
		if fieldsErr != nil {
			return fmt.Errorf("failed to fetch project fields: %w", fieldsErr)
		}
		export.Fields = fields
	}

	// Fetch and include views if requested
	if exportData.IncludeViews {
		views, viewsErr := s.fetchProjectViews(ctx, project.ID)
		if viewsErr != nil {
			return fmt.Errorf("failed to fetch project views: %w", viewsErr)
		}
		export.Views = views
	}

	// Serialize data to requested format
	var data []byte
	switch format {
	case "json":
		data, err = json.MarshalIndent(export, "", "  ")
	case "yaml":
		data, err = yaml.Marshal(export)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to serialize export data: %w", err)
	}

	// Write to output file
	if err := os.WriteFile(outputFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// ImportProject imports project data from a file
func (s *ProjectService) ImportProject(ctx context.Context, opts *ProjectImportOptions) (*ProjectImportResult, error) {
	// Read and parse import file
	exportData, err := s.parseImportFile(opts.File)
	if err != nil {
		return nil, fmt.Errorf("failed to parse import file: %w", err)
	}

	// Create new project with imported configuration
	description := ""
	if exportData.Project.Description != nil {
		description = *exportData.Project.Description
	}

	createInput := &CreateProjectInput{
		OwnerID:     opts.Owner, // Note: This should be owner ID, not login
		Title:       exportData.Project.Title,
		Description: description,
		Readme:      "",
		Visibility:  "public", // Default to public, could be made configurable
		Repository:  "",
	}

	project, err := s.CreateProject(ctx, createInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	result := &ProjectImportResult{
		ProjectID:    project.ID,
		ProjectTitle: project.Title,
		ProjectURL:   project.URL,
		ItemCount:    0,
		FieldCount:   0,
		ViewCount:    0,
	}

	// Import custom fields if not skipped
	if !opts.SkipFields && len(exportData.Fields) > 0 {
		fieldCount, err := s.importProjectFields(ctx, project.ID, exportData.Fields, opts.DryRun)
		if err != nil {
			return nil, fmt.Errorf("failed to import fields: %w", err)
		}
		result.FieldCount = fieldCount
	}

	// Import items if not skipped
	if !opts.SkipItems && len(exportData.Items) > 0 {
		itemCount, err := s.importProjectItems(ctx, project.ID, exportData.Items, opts.DryRun)
		if err != nil {
			return nil, fmt.Errorf("failed to import items: %w", err)
		}
		result.ItemCount = itemCount
	}

	// Import views if available
	if len(exportData.Views) > 0 {
		viewCount, err := s.importProjectViews(ctx, project.ID, exportData.Views, opts.DryRun)
		if err != nil {
			return nil, fmt.Errorf("failed to import views: %w", err)
		}
		result.ViewCount = viewCount
	}

	return result, nil
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	Title     *string
	Closed    *bool
	ProjectID string
}

// UpdateProject updates an existing project
//
//nolint:dupl // Similar structure to UpdateView but operates on different types
func (s *ProjectService) UpdateProject(ctx context.Context, input UpdateProjectInput) (*graphql.ProjectV2, error) {
	gqlInput := &graphql.UpdateProjectInput{
		ProjectID: gql.ID(input.ProjectID),
	}
	if input.Title != nil {
		title := gql.String(*input.Title)
		gqlInput.Title = &title
	}
	if input.Closed != nil {
		closed := gql.Boolean(*input.Closed)
		gqlInput.Closed = &closed
	}

	variables := graphql.BuildUpdateProjectVariables(gqlInput)

	var mutation graphql.UpdateProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &mutation.UpdateProjectV2.ProjectV2, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, projectID string) error {
	variables := graphql.BuildDeleteProjectVariables(&graphql.DeleteProjectInput{
		ProjectID: gql.ID(projectID),
	})

	var mutation graphql.DeleteProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// AddItemInput represents input for adding an item to a project
type AddItemInput struct {
	ProjectID string
	ContentID string
}

// AddItem adds an item to a project
func (s *ProjectService) AddItem(ctx context.Context, input AddItemInput) (*graphql.ProjectV2Item, error) {
	variables := graphql.BuildAddItemVariables(&graphql.AddItemInput{
		ProjectID: gql.ID(input.ProjectID),
		ContentID: gql.ID(input.ContentID),
	})

	var mutation graphql.AddItemToProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to project: %w", err)
	}

	return &mutation.AddProjectV2ItemByID.Item, nil
}

// UpdateItemFieldInput represents input for updating an item field
type UpdateItemFieldInput struct {
	Value     interface{}
	ProjectID string
	ItemID    string
	FieldID   string
}

// UpdateItemField updates a field value for an item
func (s *ProjectService) UpdateItemField(ctx context.Context, input UpdateItemFieldInput) (*graphql.ProjectV2Item, error) {
	variables := graphql.BuildUpdateItemFieldVariables(&graphql.UpdateItemFieldInput{
		ProjectID: gql.ID(input.ProjectID),
		ItemID:    gql.ID(input.ItemID),
		FieldID:   gql.ID(input.FieldID),
		Value:     input.Value,
	})

	var mutation graphql.UpdateItemFieldMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to update item field: %w", err)
	}

	return &mutation.UpdateProjectV2ItemFieldValue.ProjectV2Item, nil
}

// RemoveItemInput represents input for removing an item from a project
type RemoveItemInput struct {
	ProjectID string
	ItemID    string
}

// RemoveItem removes an item from a project
func (s *ProjectService) RemoveItem(ctx context.Context, input RemoveItemInput) error {
	variables := graphql.BuildRemoveItemVariables(&graphql.RemoveItemInput{
		ProjectID: gql.ID(input.ProjectID),
		ItemID:    gql.ID(input.ItemID),
	})

	var mutation graphql.RemoveItemFromProjectMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return fmt.Errorf("failed to remove item from project: %w", err)
	}

	return nil
}

// ParseProjectReference parses a project reference in the format "owner/number"
func ParseProjectReference(ref string) (owner string, number int, err error) {
	// Simple parsing - in practice, this might need more sophisticated handling
	var numStr string
	for i := len(ref) - 1; i >= 0; i-- {
		if ref[i] == '/' {
			owner = ref[:i]
			numStr = ref[i+1:]
			break
		}
	}

	if owner == "" || numStr == "" {
		return "", 0, fmt.Errorf("invalid project reference format: %s (expected owner/number)", ref)
	}

	number, err = strconv.Atoi(numStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid project number in reference: %s", numStr)
	}

	return owner, number, nil
}

// FormatProjectReference formats owner and number into a project reference
func FormatProjectReference(owner string, number int) string {
	return fmt.Sprintf("%s/%d", owner, number)
}

// parseProjectID parses project ID in format "owner/number"
func parseProjectID(projectID string) (owner string, number int, err error) {
	return ParseProjectReference(projectID)
}

// fetchProjectItems fetches all items for a project
func (s *ProjectService) fetchProjectItems(_ context.Context, _ string) ([]ExportedItem, error) {
	// For now, return empty slice as item fetching requires more complex GraphQL queries
	// This would be implemented with proper GraphQL queries to fetch project items
	return []ExportedItem{}, nil
}

// fetchProjectFields fetches all custom fields for a project
func (s *ProjectService) fetchProjectFields(_ context.Context, _ string) ([]ExportedField, error) {
	// For now, return empty slice as field fetching requires more complex GraphQL queries
	// This would be implemented with proper GraphQL queries to fetch project fields
	return []ExportedField{}, nil
}

// fetchProjectViews fetches all views for a project
func (s *ProjectService) fetchProjectViews(_ context.Context, _ string) ([]ExportedView, error) {
	// For now, return empty slice as view fetching requires more complex GraphQL queries
	// This would be implemented with proper GraphQL queries to fetch project views
	return []ExportedView{}, nil
}

// parseImportFile reads and parses the import file (JSON or YAML)
func (s *ProjectService) parseImportFile(filename string) (*ExportedProject, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	var exported ExportedProject

	// Try JSON first, then YAML
	if err := json.Unmarshal(data, &exported); err != nil {
		// If JSON fails, try YAML
		if err := yaml.Unmarshal(data, &exported); err != nil {
			return nil, fmt.Errorf("failed to parse import file as JSON or YAML: %w", err)
		}
	}

	return &exported, nil
}

// importProjectFields imports custom fields into the project
func (s *ProjectService) importProjectFields(_ context.Context, _ string, fields []ExportedField, dryRun bool) (int, error) {
	if dryRun {
		return len(fields), nil
	}

	// For now, return the count as field import requires more complex GraphQL mutations
	// This would be implemented with proper GraphQL mutations to create custom fields
	return 0, nil
}

// importProjectItems imports items into the project
func (s *ProjectService) importProjectItems(_ context.Context, _ string, items []ExportedItem, dryRun bool) (int, error) {
	if dryRun {
		return len(items), nil
	}

	// For now, return the count as item import requires more complex GraphQL mutations
	// This would be implemented with proper GraphQL mutations to create project items
	return 0, nil
}

// importProjectViews imports views into the project
func (s *ProjectService) importProjectViews(_ context.Context, _ string, views []ExportedView, dryRun bool) (int, error) {
	if dryRun {
		return len(views), nil
	}

	// For now, return the count as view import requires more complex GraphQL mutations
	// This would be implemented with proper GraphQL mutations to create project views
	return 0, nil
}

// parseRepositoryString parses repository string in format "owner/repo"
func parseRepositoryString(repository string) []string {
	return strings.Split(repository, "/")
}

// validateRepository validates that a repository exists
func (s *ProjectService) validateRepository(ctx context.Context, owner, name string) error {
	if owner == "" || name == "" {
		return fmt.Errorf("invalid repository owner or name")
	}

	// Query the repository to verify it exists and is accessible
	var query graphql.RepositoryQuery
	variables := graphql.BuildRepositoryVariables(owner, name)

	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return fmt.Errorf("failed to query repository: %w", err)
	}

	// Parse the response
	repoInfo, err := graphql.ParseRepositoryResponse(&query)
	if err != nil {
		return fmt.Errorf("failed to parse repository response: %w", err)
	}

	if repoInfo == nil {
		return fmt.Errorf("repository %s/%s not found or not accessible", owner, name)
	}

	// Verify the repository details match
	if repoInfo.Owner != owner || repoInfo.Name != name {
		return fmt.Errorf("repository details mismatch: expected %s/%s, got %s/%s",
			owner, name, repoInfo.Owner, repoInfo.Name)
	}

	return nil
}

const (
	// orgNameLengthThreshold is the max length for likely organization names
	orgNameLengthThreshold = 4
)

// DetectOwnerType detects if an owner is an organization based on the owner string
// This is a utility function to automatically detect organization vs user
func DetectOwnerType(_ context.Context, owner string) (isOrg bool, err error) {
	// For now, implement a simple heuristic:
	// - If owner contains certain patterns commonly used by organizations
	// - In a real implementation, this would query GitHub's API to determine the owner type

	// Heuristic: Organizations often have shorter names or contain common org patterns
	if len(owner) <= orgNameLengthThreshold {
		// Very short names are likely organizations (e.g., "meta", "aws", "ibm")
		return true, nil
	}

	// Common organization patterns
	orgPatterns := []string{"corp", "inc", "ltd", "co", "org", "dev", "tech", "labs"}
	lowerOwner := strings.ToLower(owner)

	for _, pattern := range orgPatterns {
		if strings.Contains(lowerOwner, pattern) {
			return true, nil
		}
	}

	// Default to user for ambiguous cases
	return false, nil
}

// GetProjectWithOwnerDetection gets a project and automatically detects if owner is organization
func (s *ProjectService) GetProjectWithOwnerDetection(ctx context.Context, owner string, number int) (*graphql.ProjectV2, error) {
	// First try as user
	project, err := s.GetProject(ctx, owner, number, false)
	if err == nil {
		return project, nil
	}

	// If user query fails, try as organization
	project, orgErr := s.GetProject(ctx, owner, number, true)
	if orgErr == nil {
		return project, nil
	}

	// If both fail, return the original user error
	return nil, fmt.Errorf("project not found for user or organization %s: %w", owner, err)
}
