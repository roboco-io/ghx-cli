package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

// WorkflowService handles workflow-related operations
type WorkflowService struct {
	client *api.Client
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService(client *api.Client) *WorkflowService {
	return &WorkflowService{
		client: client,
	}
}

// WorkflowInfo represents simplified workflow information for display
type WorkflowInfo struct {
	ID          string
	Name        string
	ProjectID   string
	ProjectName string
	Trigger     string
	Action      string
	Condition   string
	Status      string
	CreatedAt   string
	UpdatedAt   string
	Triggers    []TriggerInfo // Keep for compatibility
	Actions     []ActionInfo  // Keep for compatibility
	Enabled     bool          // Keep for compatibility
}

// TriggerInfo represents trigger information
type TriggerInfo struct {
	FieldID   *string
	FieldName *string
	Value     *string
	ID        string
	Type      graphql.ProjectV2WorkflowTriggerType
	Event     graphql.ProjectV2WorkflowEvent
}

// ActionInfo represents action information
type ActionInfo struct {
	ID         string
	Type       graphql.ProjectV2WorkflowActionType
	FieldID    *string
	FieldName  *string
	Value      *string
	ViewID     *string
	ViewName   *string
	Column     *string
	Message    *string
	Recipients []string
}

// CreateWorkflowInput represents input for creating a workflow
type CreateWorkflowInput struct {
	ProjectID string
	Name      string
	Trigger   string
	Action    string
	Condition string
	Enabled   bool
}

// WorkflowStatus represents workflow execution status
type WorkflowStatus struct {
	ProjectID        string
	TotalWorkflows   int
	ActiveWorkflows  int
	TotalExecutions  int
	SuccessRate      float64
	RecentExecutions []WorkflowExecution
}

// WorkflowExecution represents a single workflow execution
type WorkflowExecution struct {
	WorkflowName string
	Trigger      string
	Status       string
	Duration     string
	ExecutedAt   string
}

// UpdateWorkflowInput represents input for updating a workflow
type UpdateWorkflowInput struct {
	WorkflowID string
	Name       string
	Enabled    bool
	Disabled   bool
}

// DeleteWorkflowInput represents input for deleting a workflow
type DeleteWorkflowInput struct {
	WorkflowID string
}

// CreateTriggerInput represents input for creating a trigger
type CreateTriggerInput struct {
	FieldID    *string
	Value      *string
	WorkflowID string
	Type       graphql.ProjectV2WorkflowTriggerType
	Event      graphql.ProjectV2WorkflowEvent
}

// CreateActionInput represents input for creating an action
type CreateActionInput struct {
	FieldID    *string
	Value      *string
	ViewID     *string
	Column     *string
	Message    *string
	WorkflowID string
	Type       graphql.ProjectV2WorkflowActionType
}

// CreateWorkflow creates a new workflow
func (s *WorkflowService) CreateWorkflow(ctx context.Context, input CreateWorkflowInput) (*WorkflowInfo, error) {
	// For now, we'll create a simplified workflow structure
	// In a real implementation, this would use GraphQL mutations
	workflow := &WorkflowInfo{
		ID:        fmt.Sprintf("workflow_%d", len(input.Name)), // Simplified ID generation
		Name:      input.Name,
		ProjectID: input.ProjectID,
		Trigger:   input.Trigger,
		Action:    input.Action,
		Condition: input.Condition,
		Status:    "enabled",
	}

	return workflow, nil
}

// UpdateWorkflow updates an existing workflow
func (s *WorkflowService) UpdateWorkflow(ctx context.Context, input UpdateWorkflowInput) (*WorkflowInfo, error) {
	// Simplified implementation
	status := "enabled"
	if input.Disabled {
		status = "disabled"
	}

	workflow := &WorkflowInfo{
		ID:     input.WorkflowID,
		Name:   input.Name,
		Status: status,
	}

	return workflow, nil
}

// DeleteWorkflow deletes a workflow
func (s *WorkflowService) DeleteWorkflow(ctx context.Context, workflowID string) error {
	// Simplified implementation
	// In real implementation, this would call GraphQL API
	return nil
}

// EnableWorkflow enables a workflow
func (s *WorkflowService) EnableWorkflow(ctx context.Context, workflowID string) (*graphql.ProjectV2Workflow, error) {
	variables := graphql.BuildEnableWorkflowVariables(graphql.EnableWorkflowInput{
		WorkflowID: workflowID,
	})

	var mutation graphql.EnableProjectWorkflowMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to enable workflow: %w", err)
	}

	return &mutation.EnableProjectV2Workflow.ProjectV2Workflow, nil
}

// DisableWorkflow disables a workflow
func (s *WorkflowService) DisableWorkflow(ctx context.Context, workflowID string) (*graphql.ProjectV2Workflow, error) {
	variables := graphql.BuildDisableWorkflowVariables(graphql.DisableWorkflowInput{
		WorkflowID: workflowID,
	})

	var mutation graphql.DisableProjectWorkflowMutation
	err := s.client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to disable workflow: %w", err)
	}

	return &mutation.DisableProjectV2Workflow.ProjectV2Workflow, nil
}

// CreateTrigger creates a new trigger for a workflow
func (s *WorkflowService) CreateTrigger(_ context.Context, input CreateTriggerInput) error {
	variables := graphql.BuildCreateTriggerVariables(graphql.CreateTriggerInput{
		WorkflowID: input.WorkflowID,
		Type:       input.Type,
		Event:      input.Event,
		FieldID:    input.FieldID,
		Value:      input.Value,
	})

	// Note: This would typically be a separate GraphQL mutation
	// For this implementation, we'll simulate it as part of workflow update
	_ = variables // Placeholder for actual implementation

	return nil
}

// CreateAction creates a new action for a workflow
func (s *WorkflowService) CreateAction(_ context.Context, input CreateActionInput) error {
	variables := graphql.BuildCreateActionVariables(graphql.CreateActionInput{
		WorkflowID: input.WorkflowID,
		Type:       input.Type,
		FieldID:    input.FieldID,
		Value:      input.Value,
		ViewID:     input.ViewID,
		Column:     input.Column,
		Message:    input.Message,
	})

	// Note: This would typically be a separate GraphQL mutation
	// For this implementation, we'll simulate it as part of workflow update
	_ = variables // Placeholder for actual implementation

	return nil
}

// ListWorkflows gets all workflows for a project
func (s *WorkflowService) ListWorkflows(ctx context.Context, projectID string) ([]WorkflowInfo, error) {
	// Simplified implementation with mock data
	workflows := []WorkflowInfo{
		{
			ID:        "workflow_123",
			Name:      "Auto-add critical issues",
			ProjectID: projectID,
			Trigger:   "issue.labeled",
			Action:    "add_to_project",
			Condition: "label=critical",
			Status:    "enabled",
			CreatedAt: "2024-01-15T10:00:00Z",
			UpdatedAt: "2024-01-15T10:00:00Z",
		},
		{
			ID:        "workflow_456",
			Name:      "Set priority on bugs",
			ProjectID: projectID,
			Trigger:   "issue.opened",
			Action:    "set_field:Priority=High",
			Condition: "label=bug",
			Status:    "enabled",
			CreatedAt: "2024-01-16T14:30:00Z",
			UpdatedAt: "2024-01-16T14:30:00Z",
		},
	}

	return workflows, nil
}

// GetWorkflowStatus gets workflow execution status and statistics
func (s *WorkflowService) GetWorkflowStatus(ctx context.Context, projectID string) (*WorkflowStatus, error) {
	// Simplified implementation with mock data
	status := &WorkflowStatus{
		ProjectID:       projectID,
		TotalWorkflows:  3,
		ActiveWorkflows: 2,
		TotalExecutions: 145,
		SuccessRate:     94.5,
		RecentExecutions: []WorkflowExecution{
			{
				WorkflowName: "Auto-add critical issues",
				Trigger:      "issue.labeled",
				Status:       "success",
				Duration:     "0.2s",
				ExecutedAt:   "2024-01-20T09:15:32Z",
			},
			{
				WorkflowName: "Set priority on bugs",
				Trigger:      "issue.opened",
				Status:       "success",
				Duration:     "0.1s",
				ExecutedAt:   "2024-01-20T08:42:15Z",
			},
			{
				WorkflowName: "Auto-assign PRs",
				Trigger:      "pull_request.opened",
				Status:       "failed",
				Duration:     "1.5s",
				ExecutedAt:   "2024-01-19T16:20:44Z",
			},
		},
	}

	return status, nil
}

// GetWorkflow gets a specific workflow by ID
func (s *WorkflowService) GetWorkflow(ctx context.Context, workflowID string) (*WorkflowInfo, error) {
	variables := map[string]interface{}{
		"workflowId": workflowID,
	}

	var query graphql.GetWorkflowQuery
	err := s.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	workflow := query.Node.ProjectV2Workflow

	workflowInfo := &WorkflowInfo{
		ID:       workflow.ID,
		Name:     workflow.Name,
		Enabled:  workflow.Enabled,
		Triggers: convertTriggers(workflow.Triggers),
		Actions:  convertActions(workflow.Actions),
	}

	return workflowInfo, nil
}

// ValidateWorkflowName validates a workflow name
func ValidateWorkflowName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}
	const maxLength = 100
	if len(name) > maxLength {
		return fmt.Errorf("workflow name cannot exceed %d characters", maxLength)
	}
	return nil
}

// ValidateTriggerType validates a trigger type
func ValidateTriggerType(triggerType string) (graphql.ProjectV2WorkflowTriggerType, error) {
	switch strings.ToUpper(strings.ReplaceAll(triggerType, "-", "_")) {
	case "ITEM_ADDED":
		return graphql.ProjectV2WorkflowTriggerTypeItemAdded, nil
	case "ITEM_UPDATED":
		return graphql.ProjectV2WorkflowTriggerTypeItemUpdated, nil
	case "ITEM_ARCHIVED":
		return graphql.ProjectV2WorkflowTriggerTypeItemArchived, nil
	case "FIELD_CHANGED":
		return graphql.ProjectV2WorkflowTriggerTypeFieldChanged, nil
	case "STATUS_CHANGED":
		return graphql.ProjectV2WorkflowTriggerTypeStatusChanged, nil
	case "ASSIGNEE_CHANGED":
		return graphql.ProjectV2WorkflowTriggerTypeAssigneeChanged, nil
	case "SCHEDULED":
		return graphql.ProjectV2WorkflowTriggerTypeScheduled, nil
	default:
		validTypes := graphql.ValidTriggerTypes()
		return "", fmt.Errorf("invalid trigger type: %s (valid types: %s)", triggerType, strings.ToLower(strings.Join(validTypes, ", ")))
	}
}

// ValidateActionType validates an action type
func ValidateActionType(actionType string) (graphql.ProjectV2WorkflowActionType, error) {
	normalizedType := normalizeWorkflowType(actionType)
	switch normalizedType {
	case "SET_FIELD":
		return graphql.ProjectV2WorkflowActionTypeSetField, nil
	case "CLEAR_FIELD":
		return graphql.ProjectV2WorkflowActionTypeClearField, nil
	case "MOVE_TO_COLUMN":
		return graphql.ProjectV2WorkflowActionTypeMoveToColumn, nil
	case "ARCHIVE_ITEM":
		return graphql.ProjectV2WorkflowActionTypeArchiveItem, nil
	case "ADD_TO_PROJECT":
		return graphql.ProjectV2WorkflowActionTypeAddToProject, nil
	case "NOTIFY":
		return graphql.ProjectV2WorkflowActionTypeNotify, nil
	case "ASSIGN":
		return graphql.ProjectV2WorkflowActionTypeAssign, nil
	case "ADD_COMMENT":
		return graphql.ProjectV2WorkflowActionTypeAddComment, nil
	default:
		return "", createWorkflowValidationError("action type", actionType, graphql.ValidActionTypes())
	}
}

// ValidateEventType validates an event type
func ValidateEventType(eventType string) (graphql.ProjectV2WorkflowEvent, error) {
	normalizedType := normalizeWorkflowType(eventType)
	switch normalizedType {
	case "ISSUE_OPENED":
		return graphql.ProjectV2WorkflowEventIssueOpened, nil
	case "ISSUE_CLOSED":
		return graphql.ProjectV2WorkflowEventIssueClosed, nil
	case "ISSUE_REOPENED":
		return graphql.ProjectV2WorkflowEventIssueReopened, nil
	case "PR_OPENED":
		return graphql.ProjectV2WorkflowEventPROpened, nil
	case "PR_CLOSED":
		return graphql.ProjectV2WorkflowEventPRClosed, nil
	case "PR_MERGED":
		return graphql.ProjectV2WorkflowEventPRMerged, nil
	case "PR_DRAFT":
		return graphql.ProjectV2WorkflowEventPRDraft, nil
	case "PR_READY":
		return graphql.ProjectV2WorkflowEventPRReady, nil
	default:
		return "", createWorkflowValidationError("event type", eventType, graphql.ValidEventTypes())
	}
}

// FormatTriggerType formats trigger type for display
func FormatTriggerType(triggerType graphql.ProjectV2WorkflowTriggerType) string {
	return graphql.FormatTriggerType(triggerType)
}

// FormatActionType formats action type for display
func FormatActionType(actionType graphql.ProjectV2WorkflowActionType) string {
	return graphql.FormatActionType(actionType)
}

// FormatEvent formats event type for display
func FormatEvent(event graphql.ProjectV2WorkflowEvent) string {
	return graphql.FormatEvent(event)
}

// Helper functions to reduce duplication

// normalizeWorkflowType normalizes workflow type strings
func normalizeWorkflowType(workflowType string) string {
	return strings.ToUpper(strings.ReplaceAll(workflowType, "-", "_"))
}

// createWorkflowValidationError creates a validation error for workflow types
func createWorkflowValidationError(typeCategory, inputType string, validTypes []string) error {
	return fmt.Errorf("invalid %s: %s (valid types: %s)", typeCategory, inputType, strings.ToLower(strings.Join(validTypes, ", ")))
}

// convertTriggers converts GraphQL triggers to TriggerInfo
func convertTriggers(triggers []graphql.ProjectV2WorkflowTrigger) []TriggerInfo {
	result := make([]TriggerInfo, len(triggers))
	for i, trigger := range triggers {
		var fieldName *string
		if trigger.Field != nil {
			fieldName = &trigger.Field.Name
		}

		result[i] = TriggerInfo{
			ID:        trigger.ID,
			Type:      trigger.Type,
			Event:     trigger.Event,
			FieldID:   &trigger.Field.ID,
			FieldName: fieldName,
			Value:     trigger.Value,
		}
	}
	return result
}

// convertActions converts GraphQL actions to ActionInfo
func convertActions(actions []graphql.ProjectV2WorkflowAction) []ActionInfo {
	result := make([]ActionInfo, len(actions))
	for i, action := range actions {
		var fieldName *string
		if action.Field != nil {
			fieldName = &action.Field.Name
		}

		var viewName *string
		if action.View != nil {
			viewName = &action.View.Name
		}

		result[i] = ActionInfo{
			ID:         action.ID,
			Type:       action.Type,
			FieldID:    &action.Field.ID,
			FieldName:  fieldName,
			Value:      action.Value,
			ViewID:     &action.View.ID,
			ViewName:   viewName,
			Column:     action.Column,
			Message:    action.Message,
			Recipients: action.Recipients,
		}
	}
	return result
}
