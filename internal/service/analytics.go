package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

// AnalyticsService handles analytics and reporting operations
type AnalyticsService struct {
	client *api.Client
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(client *api.Client) *AnalyticsService {
	return &AnalyticsService{
		client: client,
	}
}

// Analytics Input Types

// ExportProjectInput represents input for project export
type ExportProjectInput struct {
	Filter           *string
	ProjectID        string
	Format           graphql.ProjectV2ExportFormat
	IncludeItems     bool
	IncludeFields    bool
	IncludeViews     bool
	IncludeWorkflows bool
}

// ImportProjectInput represents input for project import
type ImportProjectInput struct {
	ProjectID     string
	Format        graphql.ProjectV2ExportFormat
	Data          string
	MergeStrategy string
}

// BulkUpdateItemsInput represents input for bulk item update
type BulkUpdateItemsInput struct {
	Updates   map[string]interface{}
	ProjectID string
	ItemIDs   []string
}

// BulkDeleteItemsInput represents input for bulk item delete
type BulkDeleteItemsInput struct {
	ProjectID string
	ItemIDs   []string
}

// BulkArchiveItemsInput represents input for bulk item archive
type BulkArchiveItemsInput struct {
	ProjectID string
	ItemIDs   []string
}

// Service Methods

// GetProjectAnalytics gets analytics data for a project
func (s *AnalyticsService) GetProjectAnalytics(ctx context.Context, projectID string) (*graphql.ProjectV2Analytics, error) {
	var query graphql.GetProjectAnalyticsQuery
	variables := map[string]interface{}{
		"projectId": projectID,
	}

	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get project analytics: %w", err)
	}

	return &query.Node.ProjectV2, nil
}

// ExportProject exports a project with specified options
func (s *AnalyticsService) ExportProject(ctx context.Context, input ExportProjectInput) (*graphql.ProjectV2Export, error) {
	var mutation graphql.ExportProjectMutation
	variables := graphql.BuildExportProjectVariables(graphql.ExportProjectInput{
		ProjectID:        input.ProjectID,
		Format:           input.Format,
		IncludeItems:     input.IncludeItems,
		IncludeFields:    input.IncludeFields,
		IncludeViews:     input.IncludeViews,
		IncludeWorkflows: input.IncludeWorkflows,
		Filter:           input.Filter,
	})

	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to export project: %w", err)
	}

	return &mutation.ExportProjectV2.Export, nil
}

// ImportProject imports project data
func (s *AnalyticsService) ImportProject(ctx context.Context, input ImportProjectInput) (*graphql.BulkOperation, error) {
	var mutation graphql.ImportProjectMutation
	variables := graphql.BuildImportProjectVariables(graphql.ImportProjectInput{
		ProjectID:     input.ProjectID,
		Format:        input.Format,
		Data:          input.Data,
		MergeStrategy: input.MergeStrategy,
	})

	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to import project: %w", err)
	}

	return &mutation.ImportProjectV2.BulkOperation, nil
}

// BulkUpdateItems performs bulk update on items
func (s *AnalyticsService) BulkUpdateItems(ctx context.Context, input BulkUpdateItemsInput) (*graphql.BulkOperation, error) {
	var mutation graphql.BulkUpdateItemsMutation
	variables := graphql.BuildBulkUpdateItemsVariables(graphql.BulkUpdateItemsInput{
		ProjectID: input.ProjectID,
		ItemIDs:   input.ItemIDs,
		Updates:   input.Updates,
	})

	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk update items: %w", err)
	}

	return &mutation.BulkUpdateProjectV2Items.BulkOperation, nil
}

// BulkDeleteItems performs bulk delete on items
func (s *AnalyticsService) BulkDeleteItems(ctx context.Context, input BulkDeleteItemsInput) (*graphql.BulkOperation, error) {
	var mutation graphql.BulkDeleteItemsMutation
	variables := graphql.BuildBulkDeleteItemsVariables(graphql.BulkDeleteItemsInput{
		ProjectID: input.ProjectID,
		ItemIDs:   input.ItemIDs,
	})

	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk delete items: %w", err)
	}

	return &mutation.BulkDeleteProjectV2Items.BulkOperation, nil
}

// BulkArchiveItems performs bulk archive on items
func (s *AnalyticsService) BulkArchiveItems(ctx context.Context, input BulkArchiveItemsInput) (*graphql.BulkOperation, error) {
	var mutation graphql.BulkArchiveItemsMutation
	variables := graphql.BuildBulkArchiveItemsVariables(graphql.BulkArchiveItemsInput{
		ProjectID: input.ProjectID,
		ItemIDs:   input.ItemIDs,
	})

	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk archive items: %w", err)
	}

	return &mutation.BulkArchiveProjectV2Items.BulkOperation, nil
}

// GetBulkOperation gets bulk operation status
func (s *AnalyticsService) GetBulkOperation(ctx context.Context, operationID string) (*graphql.BulkOperation, error) {
	var query graphql.GetBulkOperationQuery
	variables := map[string]interface{}{
		"operationId": operationID,
	}

	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get bulk operation: %w", err)
	}

	return &query.Node.BulkOperation, nil
}

// Validation Functions

// ValidateExportFormat validates export format
func ValidateExportFormat(format string) (graphql.ProjectV2ExportFormat, error) {
	normalizedFormat := strings.ToUpper(strings.TrimSpace(format))

	switch normalizedFormat {
	case "JSON":
		return graphql.ProjectV2ExportFormatJSON, nil
	case "CSV":
		return graphql.ProjectV2ExportFormatCSV, nil
	case "XML":
		return graphql.ProjectV2ExportFormatXML, nil
	default:
		return "", fmt.Errorf("invalid export format '%s', must be one of: %v", format, graphql.ValidExportFormats())
	}
}

// ValidateBulkOperationType validates bulk operation type
func ValidateBulkOperationType(opType string) (graphql.BulkOperationType, error) {
	normalizedType := strings.ToUpper(strings.TrimSpace(opType))
	normalizedType = strings.ReplaceAll(normalizedType, "-", "_")

	switch normalizedType {
	case "UPDATE":
		return graphql.BulkOperationTypeUpdate, nil
	case "DELETE":
		return graphql.BulkOperationTypeDelete, nil
	case "IMPORT":
		return graphql.BulkOperationTypeImport, nil
	case "EXPORT":
		return graphql.BulkOperationTypeExport, nil
	case "ARCHIVE":
		return graphql.BulkOperationTypeArchive, nil
	case "MOVE":
		return graphql.BulkOperationTypeMove, nil
	default:
		return "", fmt.Errorf("invalid bulk operation type '%s', must be one of: %v", opType, graphql.ValidBulkOperationTypes())
	}
}

// ValidateMergeStrategy validates import merge strategy
func ValidateMergeStrategy(strategy string) error {
	normalizedStrategy := strings.ToLower(strings.TrimSpace(strategy))

	validStrategies := []string{"merge", "replace", "append", "skip_conflicts"}

	for _, valid := range validStrategies {
		if normalizedStrategy == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid merge strategy '%s', must be one of: %v", strategy, validStrategies)
}

// Helper Functions

// FormatExportFormat formats export format for display
func FormatExportFormat(format graphql.ProjectV2ExportFormat) string {
	return graphql.FormatExportFormat(format)
}

// FormatBulkOperationType formats bulk operation type for display
func FormatBulkOperationType(opType graphql.BulkOperationType) string {
	return graphql.FormatBulkOperationType(opType)
}

// FormatBulkOperationStatus formats bulk operation status for display
func FormatBulkOperationStatus(status graphql.BulkOperationStatus) string {
	return graphql.FormatBulkOperationStatus(status)
}

// Analytics Display Types

// AnalyticsInfo represents simplified analytics data for display
type AnalyticsInfo struct {
	VelocityData  *VelocityInfo
	TimelineData  *TimelineInfo
	ProjectID     string
	Title         string
	StatusStats   []StatusStat
	AssigneeStats []AssigneeStat
	ItemCount     int
	FieldCount    int
	ViewCount     int
}

// StatusStat represents status statistics
type StatusStat struct {
	Status string
	Count  int
}

// AssigneeStat represents assignee statistics
type AssigneeStat struct {
	Assignee string
	Count    int
}

// VelocityInfo represents velocity information
type VelocityInfo struct {
	Period         string
	LeadTime       string
	CycleTime      string
	CompletedItems int
	AddedItems     int
	ClosureRate    float64
}

// TimelineInfo represents timeline information
type TimelineInfo struct {
	StartDate      *string
	EndDate        *string
	Duration       int
	MilestoneCount int
	ActivityCount  int
}

// Add type aliases for GraphQL types used in service
type ProjectV2Export = graphql.ProjectV2Export
type BulkOperation = graphql.BulkOperation

// FormatAnalytics converts GraphQL analytics to display format
func FormatAnalytics(analytics *graphql.ProjectV2Analytics) *AnalyticsInfo {
	info := &AnalyticsInfo{
		ProjectID:  analytics.ProjectID,
		Title:      analytics.Title,
		ItemCount:  analytics.ItemCount,
		FieldCount: analytics.FieldCount,
		ViewCount:  analytics.ViewCount,
	}

	// Format status statistics
	for _, status := range analytics.ItemsByStatus {
		info.StatusStats = append(info.StatusStats, StatusStat{
			Status: status.Status,
			Count:  status.Count,
		})
	}

	// Format assignee statistics
	for _, assignee := range analytics.ItemsByAssignee {
		info.AssigneeStats = append(info.AssigneeStats, AssigneeStat{
			Assignee: assignee.Assignee,
			Count:    assignee.Count,
		})
	}

	// Format velocity data
	if analytics.Velocity.Period != "" {
		info.VelocityData = &VelocityInfo{
			Period:         analytics.Velocity.Period,
			CompletedItems: analytics.Velocity.CompletedItems,
			AddedItems:     analytics.Velocity.AddedItems,
			ClosureRate:    analytics.Velocity.ClosureRate,
			LeadTime:       fmt.Sprintf("%.1f %s", analytics.Velocity.LeadTime.Average, analytics.Velocity.LeadTime.Unit),
			CycleTime:      fmt.Sprintf("%.1f %s", analytics.Velocity.CycleTime.Average, analytics.Velocity.CycleTime.Unit),
		}
	}

	// Format timeline data
	var startDate, endDate *string
	if analytics.Timeline.StartDate != nil {
		start := analytics.Timeline.StartDate.Format("2006-01-02")
		startDate = &start
	}
	if analytics.Timeline.EndDate != nil {
		end := analytics.Timeline.EndDate.Format("2006-01-02")
		endDate = &end
	}

	info.TimelineData = &TimelineInfo{
		StartDate:      startDate,
		EndDate:        endDate,
		Duration:       analytics.Timeline.Duration,
		MilestoneCount: len(analytics.Timeline.Milestones),
		ActivityCount:  len(analytics.Timeline.Activities),
	}

	return info
}
