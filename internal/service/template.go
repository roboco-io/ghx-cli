package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/roboco-io/ghx-cli/internal/api"
)

// TemplateService handles template-related operations
type TemplateService struct {
	client *api.Client
}

// NewTemplateService creates a new template service
func NewTemplateService(client *api.Client) *TemplateService {
	return &TemplateService{
		client: client,
	}
}

// TemplateInfo represents template information for display
type TemplateInfo struct {
	ID          string
	Name        string
	Description string
	Category    string
	Tags        []string
	Fields      []TemplateField
	Views       []TemplateView
	Workflows   []TemplateWorkflow
	CreatedAt   string
	UpdatedAt   string
}

// TemplateField represents a field in a template
type TemplateField struct {
	Name     string
	Type     string
	Required bool
	Options  []string
}

// TemplateView represents a view in a template
type TemplateView struct {
	Name   string
	Type   string
	Fields []string
	Filter string
}

// TemplateWorkflow represents a workflow in a template
type TemplateWorkflow struct {
	Name      string
	Trigger   string
	Action    string
	Condition string
}

// CreateTemplateInput represents input for creating a template
type CreateTemplateInput struct {
	Name        string
	Description string
	ProjectID   string
	Category    string
	Tags        []string
}

// UpdateTemplateInput represents input for updating a template
type UpdateTemplateInput struct {
	TemplateID  string
	Name        string
	Description string
	Category    string
	Tags        []string
}

// ApplyTemplateInput represents input for applying a template
type ApplyTemplateInput struct {
	TemplateID  string
	ProjectName string
	Owner       string
	Org         bool
	Customize   bool
}

// ApplyTemplateResult represents the result of applying a template
type ApplyTemplateResult struct {
	ID            string
	Name          string
	Owner         string
	URL           string
	FieldCount    int
	ViewCount     int
	WorkflowCount int
}

// ExportTemplateInput represents input for exporting a template
type ExportTemplateInput struct {
	TemplateID string
	Output     string
	Format     string
}

// ImportTemplateInput represents input for importing a template
type ImportTemplateInput struct {
	File   string
	Name   string
	Update bool
}

// ListTemplates gets all available templates
func (s *TemplateService) ListTemplates(ctx context.Context) ([]TemplateInfo, error) {
	// Simplified implementation with mock data
	templates := []TemplateInfo{
		{
			ID:          "template_dev_001",
			Name:        "Software Development Sprint",
			Description: "Template for agile software development sprints with user stories, bugs, and tasks",
			Category:    "development",
			Tags:        []string{"agile", "sprint", "development"},
			Fields: []TemplateField{
				{Name: "Priority", Type: "single_select", Required: true, Options: []string{"Critical", "High", "Medium", "Low"}},
				{Name: "Story Points", Type: "number", Required: false},
				{Name: "Sprint", Type: "iteration", Required: true},
				{Name: "Epic", Type: "single_select", Required: false},
			},
			Views: []TemplateView{
				{Name: "Sprint Board", Type: "board", Fields: []string{"Status", "Priority", "Assignee"}, Filter: "sprint:current"},
				{Name: "Backlog", Type: "table", Fields: []string{"Title", "Priority", "Story Points", "Epic"}, Filter: "status:backlog"},
			},
			Workflows: []TemplateWorkflow{
				{Name: "Auto-assign critical bugs", Trigger: "issue.labeled", Action: "set_field:Priority=Critical", Condition: "label=critical-bug"},
				{Name: "Move completed items", Trigger: "issue.closed", Action: "move_to_column:Done", Condition: ""},
			},
			CreatedAt: "2024-01-15T10:00:00Z",
			UpdatedAt: "2024-01-15T10:00:00Z",
		},
		{
			ID:          "template_mkt_001",
			Name:        "Marketing Campaign Tracker",
			Description: "Template for tracking marketing campaigns with deliverables, timelines, and metrics",
			Category:    "marketing",
			Tags:        []string{"marketing", "campaign", "deliverables"},
			Fields: []TemplateField{
				{Name: "Campaign", Type: "single_select", Required: true, Options: []string{"Q1 Launch", "Brand Awareness", "Product Demo"}},
				{Name: "Channel", Type: "single_select", Required: true, Options: []string{"Social Media", "Email", "Content", "Paid Ads"}},
				{Name: "Launch Date", Type: "date", Required: true},
				{Name: "Budget", Type: "number", Required: false},
			},
			Views: []TemplateView{
				{Name: "Campaign Timeline", Type: "roadmap", Fields: []string{"Title", "Campaign", "Launch Date", "Status"}, Filter: ""},
				{Name: "Deliverables", Type: "table", Fields: []string{"Title", "Channel", "Status", "Assignee"}, Filter: "type:deliverable"},
			},
			Workflows: []TemplateWorkflow{
				{Name: "Notify on launch", Trigger: "field.changed", Action: "notify", Condition: "field:Launch Date"},
				{Name: "Archive completed campaigns", Trigger: "field.changed", Action: "archive_item", Condition: "status:completed"},
			},
			CreatedAt: "2024-01-16T14:30:00Z",
			UpdatedAt: "2024-01-16T14:30:00Z",
		},
		{
			ID:          "template_bug_001",
			Name:        "Bug Tracking System",
			Description: "Comprehensive bug tracking template with severity levels, reproduction steps, and resolution workflow",
			Category:    "development",
			Tags:        []string{"bug", "tracking", "qa", "testing"},
			Fields: []TemplateField{
				{Name: "Severity", Type: "single_select", Required: true, Options: []string{"Critical", "High", "Medium", "Low"}},
				{Name: "Component", Type: "single_select", Required: true, Options: []string{"Frontend", "Backend", "Database", "API"}},
				{Name: "Browser", Type: "single_select", Required: false, Options: []string{"Chrome", "Firefox", "Safari", "Edge"}},
				{Name: "Resolution", Type: "single_select", Required: false, Options: []string{"Fixed", "Won't Fix", "Duplicate", "Not a Bug"}},
			},
			Views: []TemplateView{
				{Name: "Triage Board", Type: "board", Fields: []string{"Status", "Severity", "Component"}, Filter: "status:new,investigating"},
				{Name: "Critical Bugs", Type: "table", Fields: []string{"Title", "Severity", "Component", "Assignee"}, Filter: "severity:critical"},
			},
			Workflows: []TemplateWorkflow{
				{Name: "Auto-assign critical bugs", Trigger: "issue.opened", Action: "set_field:Priority=Urgent", Condition: "label=critical"},
				{Name: "Close resolved bugs", Trigger: "field.changed", Action: "move_to_column:Closed", Condition: "field:Resolution"},
			},
			CreatedAt: "2024-01-17T09:15:00Z",
			UpdatedAt: "2024-01-17T09:15:00Z",
		},
		{
			ID:          "template_design_001",
			Name:        "Design System Project",
			Description: "Template for managing design systems with components, documentation, and review workflows",
			Category:    "design",
			Tags:        []string{"design", "components", "ui", "system"},
			Fields: []TemplateField{
				{Name: "Component Type", Type: "single_select", Required: true, Options: []string{"Layout", "Form", "Navigation", "Feedback", "Display"}},
				{Name: "Design Status", Type: "single_select", Required: true, Options: []string{"Draft", "Review", "Approved", "Implemented"}},
				{Name: "Platform", Type: "single_select", Required: false, Options: []string{"Web", "Mobile", "Both"}},
				{Name: "Complexity", Type: "single_select", Required: false, Options: []string{"Simple", "Medium", "Complex"}},
			},
			Views: []TemplateView{
				{Name: "Component Roadmap", Type: "roadmap", Fields: []string{"Title", "Component Type", "Design Status"}, Filter: ""},
				{Name: "Review Queue", Type: "table", Fields: []string{"Title", "Component Type", "Complexity", "Assignee"}, Filter: "status:review"},
			},
			Workflows: []TemplateWorkflow{
				{Name: "Notify on design completion", Trigger: "field.changed", Action: "notify", Condition: "field:Design Status=Approved"},
				{Name: "Move to implementation", Trigger: "field.changed", Action: "move_to_column:Implementation", Condition: "design_status:approved"},
			},
			CreatedAt: "2024-01-18T11:20:00Z",
			UpdatedAt: "2024-01-18T11:20:00Z",
		},
	}

	return templates, nil
}

// CreateTemplate creates a new template from an existing project
func (s *TemplateService) CreateTemplate(ctx context.Context, input CreateTemplateInput) (*TemplateInfo, error) {
	// Validate input
	if err := validateTemplateName(input.Name); err != nil {
		return nil, err
	}

	if err := validateTemplateCategory(input.Category); err != nil {
		return nil, err
	}

	// Simplified implementation
	template := &TemplateInfo{
		ID:          fmt.Sprintf("template_%s_%03d", strings.ToLower(input.Category)[:3], len(input.Name)),
		Name:        input.Name,
		Description: input.Description,
		Category:    input.Category,
		Tags:        input.Tags,
		Fields:      []TemplateField{},    // Would be extracted from source project
		Views:       []TemplateView{},     // Would be extracted from source project
		Workflows:   []TemplateWorkflow{}, // Would be extracted from source project
		CreatedAt:   "2024-01-20T10:00:00Z",
		UpdatedAt:   "2024-01-20T10:00:00Z",
	}

	return template, nil
}

// UpdateTemplate updates an existing template
func (s *TemplateService) UpdateTemplate(ctx context.Context, input UpdateTemplateInput) (*TemplateInfo, error) {
	// Validate input
	if input.Name != "" {
		if err := validateTemplateName(input.Name); err != nil {
			return nil, err
		}
	}

	if input.Category != "" {
		if err := validateTemplateCategory(input.Category); err != nil {
			return nil, err
		}
	}

	// Simplified implementation
	template := &TemplateInfo{
		ID:          input.TemplateID,
		Name:        input.Name,
		Description: input.Description,
		Category:    input.Category,
		Tags:        input.Tags,
		UpdatedAt:   "2024-01-20T10:00:00Z",
	}

	return template, nil
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(ctx context.Context, templateID string) error {
	// Simplified implementation
	return nil
}

// ApplyTemplate applies a template to create a new project
func (s *TemplateService) ApplyTemplate(ctx context.Context, input ApplyTemplateInput) (*ApplyTemplateResult, error) {
	// Simplified implementation
	result := &ApplyTemplateResult{
		ID:            fmt.Sprintf("project_%d", len(input.ProjectName)),
		Name:          input.ProjectName,
		Owner:         input.Owner,
		URL:           fmt.Sprintf("https://github.com/%s/projects/%s", input.Owner, input.ProjectName),
		FieldCount:    4, // Mock data
		ViewCount:     2, // Mock data
		WorkflowCount: 2, // Mock data
	}

	return result, nil
}

// ExportTemplate exports a template to a file
func (s *TemplateService) ExportTemplate(ctx context.Context, input ExportTemplateInput) error {
	// Simplified implementation
	// In real implementation, this would:
	// 1. Get template by ID
	// 2. Serialize to JSON/YAML
	// 3. Write to file
	return nil
}

// ImportTemplate imports a template from a file
func (s *TemplateService) ImportTemplate(ctx context.Context, input ImportTemplateInput) (*TemplateInfo, error) {
	// Simplified implementation
	template := &TemplateInfo{
		ID:        fmt.Sprintf("template_imported_%d", len(input.Name)),
		Name:      input.Name,
		Category:  "imported",
		CreatedAt: "2024-01-20T10:00:00Z",
		UpdatedAt: "2024-01-20T10:00:00Z",
	}

	return template, nil
}

// Validation functions
func validateTemplateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("template name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("template name cannot exceed 100 characters")
	}
	return nil
}

func validateTemplateCategory(category string) error {
	validCategories := []string{"development", "marketing", "design", "research", "general"}

	category = strings.ToLower(strings.TrimSpace(category))
	for _, valid := range validCategories {
		if category == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid category: %s (valid categories: %s)", category, strings.Join(validCategories, ", "))
}

// FormatTemplateCategory formats template category for display
func FormatTemplateCategory(category string) string {
	switch strings.ToLower(category) {
	case "development":
		return "Development"
	case "marketing":
		return "Marketing"
	case "design":
		return "Design"
	case "research":
		return "Research"
	case "general":
		return "General"
	default:
		return strings.Title(strings.ToLower(category))
	}
}
