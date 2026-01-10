package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// WorkflowOptions holds options for workflow commands
type WorkflowOptions struct {
	ProjectID string
	Format    string
}

// NewWorkflowCmd creates the workflow command group
func NewWorkflowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage project workflows and automation",
		Long: `Manage GitHub Project workflows and automation rules.

Workflows allow you to automate project management tasks based on triggers and conditions.

Subcommands:
  list    - List all workflows for a project
  create  - Create a new workflow
  update  - Update an existing workflow
  delete  - Delete a workflow
  status  - Show workflow status and statistics`,
	}

	cmd.AddCommand(
		NewWorkflowListCmd(),
		NewWorkflowCreateCmd(),
		NewWorkflowUpdateCmd(),
		NewWorkflowDeleteCmd(),
		NewWorkflowStatusCmd(),
	)

	return cmd
}

// NewWorkflowListCmd creates the workflow list command
func NewWorkflowListCmd() *cobra.Command {
	opts := &WorkflowOptions{}

	cmd := &cobra.Command{
		Use:   "list PROJECT_ID",
		Short: "List project workflows",
		Long: `List all automation workflows configured for a GitHub Project.

This shows workflow rules, triggers, and current status.

Examples:
  ghx project workflow list myorg/123
  ghx project workflow list user/456 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectID = args[0]
			return runWorkflowList(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

// NewWorkflowCreateCmd creates the workflow create command
func NewWorkflowCreateCmd() *cobra.Command {
	var (
		projectID string
		name      string
		trigger   string
		action    string
		condition string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new project workflow",
		Long: `Create a new automation workflow for a GitHub Project.

Workflows consist of triggers, conditions, and actions that automate project management tasks.

Available Triggers:
  issue.opened          - When an issue is opened
  issue.closed          - When an issue is closed  
  pull_request.opened   - When a PR is opened
  pull_request.merged   - When a PR is merged
  item.added           - When an item is added to project
  field.changed        - When a field value changes

Available Actions:
  add_to_project       - Add item to project
  set_field            - Set field value
  move_to_status       - Move to status column
  assign_user          - Assign to user
  add_label            - Add label to item

Examples:
  # Auto-add new issues to project
  ghx project workflow create --project-id myorg/123 --name "Auto-add issues" --trigger "issue.opened" --action "add_to_project"
  
  # Set priority on critical issues  
  ghx project workflow create --project-id myorg/123 --name "Critical priority" --trigger "issue.labeled" --condition "label=critical" --action "set_field:Priority=High"`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if projectID == "" || name == "" || trigger == "" || action == "" {
				return fmt.Errorf("--project-id, --name, --trigger, and --action are required")
			}

			return runWorkflowCreate(cmd.Context(), WorkflowCreateOptions{
				ProjectID: projectID,
				Name:      name,
				Trigger:   trigger,
				Action:    action,
				Condition: condition,
			})
		},
	}

	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Workflow name (required)")
	cmd.Flags().StringVar(&trigger, "trigger", "", "Workflow trigger event (required)")
	cmd.Flags().StringVar(&action, "action", "", "Workflow action to perform (required)")
	cmd.Flags().StringVar(&condition, "condition", "", "Optional condition for triggering")

	_ = cmd.MarkFlagRequired("project-id")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("trigger")
	_ = cmd.MarkFlagRequired("action")

	return cmd
}

// NewWorkflowUpdateCmd creates the workflow update command
func NewWorkflowUpdateCmd() *cobra.Command {
	var (
		workflowID string
		name       string
		enabled    bool
		disabled   bool
	)

	cmd := &cobra.Command{
		Use:   "update WORKFLOW_ID",
		Short: "Update an existing workflow",
		Long: `Update an existing project workflow.

You can modify workflow properties like name and enabled status.

Examples:
  ghx project workflow update workflow_123 --name "Updated name"
  ghx project workflow update workflow_123 --enabled
  ghx project workflow update workflow_123 --disabled`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID = args[0]
			return runWorkflowUpdate(cmd.Context(), WorkflowUpdateOptions{
				WorkflowID: workflowID,
				Name:       name,
				Enabled:    enabled && !disabled,
				Disabled:   disabled,
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Update workflow name")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "Enable the workflow")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "Disable the workflow")

	return cmd
}

// NewWorkflowDeleteCmd creates the workflow delete command
func NewWorkflowDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete WORKFLOW_ID",
		Short: "Delete a project workflow",
		Long: `Delete an automation workflow from a GitHub Project.

This permanently removes the workflow and stops all automation.

Examples:
  ghx project workflow delete workflow_123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			return runWorkflowDelete(cmd.Context(), workflowID)
		},
	}

	return cmd
}

// NewWorkflowStatusCmd creates the workflow status command
func NewWorkflowStatusCmd() *cobra.Command {
	opts := &WorkflowOptions{}

	cmd := &cobra.Command{
		Use:   "status PROJECT_ID",
		Short: "Show workflow status and statistics",
		Long: `Show status and execution statistics for project workflows.

This displays workflow performance metrics, success rates, and recent executions.

Examples:
  ghx project workflow status myorg/123
  ghx project workflow status user/456 --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectID = args[0]
			return runWorkflowStatus(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

// Command implementations

type WorkflowCreateOptions struct {
	ProjectID string
	Name      string
	Trigger   string
	Action    string
	Condition string
}

type WorkflowUpdateOptions struct {
	WorkflowID string
	Name       string
	Enabled    bool
	Disabled   bool
}

func runWorkflowList(ctx context.Context, opts *WorkflowOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	workflows, err := workflowService.ListWorkflows(ctx, opts.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	if len(workflows) == 0 {
		fmt.Printf("No workflows found for project %s\n", opts.ProjectID)
		return nil
	}

	switch opts.Format {
	case "json":
		return outputWorkflowsJSON(workflows)
	case "table":
		return outputWorkflowsTable(workflows)
	default:
		return fmt.Errorf("unknown format: %s", opts.Format)
	}
}

func runWorkflowCreate(ctx context.Context, opts WorkflowCreateOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	workflow, err := workflowService.CreateWorkflow(ctx, service.CreateWorkflowInput{
		ProjectID: opts.ProjectID,
		Name:      opts.Name,
		Trigger:   opts.Trigger,
		Action:    opts.Action,
		Condition: opts.Condition,
	})
	if err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}

	fmt.Printf("✅ Workflow '%s' created successfully\n\n", workflow.Name)
	fmt.Printf("Workflow Details:\n")
	fmt.Printf("  ID: %s\n", workflow.ID)
	fmt.Printf("  Name: %s\n", workflow.Name)
	fmt.Printf("  Trigger: %s\n", workflow.Trigger)
	fmt.Printf("  Action: %s\n", workflow.Action)
	if workflow.Condition != "" {
		fmt.Printf("  Condition: %s\n", workflow.Condition)
	}
	fmt.Printf("  Status: %s\n", workflow.Status)

	return nil
}

func runWorkflowUpdate(ctx context.Context, opts WorkflowUpdateOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	workflow, err := workflowService.UpdateWorkflow(ctx, service.UpdateWorkflowInput{
		WorkflowID: opts.WorkflowID,
		Name:       opts.Name,
		Enabled:    opts.Enabled,
		Disabled:   opts.Disabled,
	})
	if err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
	}

	fmt.Printf("✅ Workflow updated successfully\n\n")
	fmt.Printf("  ID: %s\n", workflow.ID)
	fmt.Printf("  Name: %s\n", workflow.Name)
	fmt.Printf("  Status: %s\n", workflow.Status)

	return nil
}

func runWorkflowDelete(ctx context.Context, workflowID string) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	err = workflowService.DeleteWorkflow(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	fmt.Printf("✅ Workflow %s deleted successfully\n", workflowID)
	return nil
}

func runWorkflowStatus(ctx context.Context, opts *WorkflowOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	workflowService := service.NewWorkflowService(client)

	status, err := workflowService.GetWorkflowStatus(ctx, opts.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to get workflow status: %w", err)
	}

	switch opts.Format {
	case "json":
		return outputWorkflowStatusJSON(status)
	case "table":
		return outputWorkflowStatusTable(status)
	default:
		return fmt.Errorf("unknown format: %s", opts.Format)
	}
}

// Output functions
func outputWorkflowsTable(workflows []service.WorkflowInfo) error {
	fmt.Printf("Project Workflows:\n\n")
	fmt.Printf("%-20s %-30s %-20s %-15s %-10s\n", "ID", "NAME", "TRIGGER", "ACTION", "STATUS")
	fmt.Printf("%-20s %-30s %-20s %-15s %-10s\n", "──", "────", "───────", "──────", "──────")

	for _, workflow := range workflows {
		fmt.Printf("%-20s %-30s %-20s %-15s %-10s\n",
			truncateString(workflow.ID, 20),
			truncateString(workflow.Name, 30),
			truncateString(workflow.Trigger, 20),
			truncateString(workflow.Action, 15),
			workflow.Status,
		)
	}

	fmt.Printf("\n%d workflows total\n", len(workflows))
	return nil
}

func outputWorkflowsJSON(workflows []service.WorkflowInfo) error {
	fmt.Printf("[\n")
	for i, workflow := range workflows {
		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": \"%s\",\n", workflow.ID)
		fmt.Printf("    \"name\": \"%s\",\n", workflow.Name)
		fmt.Printf("    \"trigger\": \"%s\",\n", workflow.Trigger)
		fmt.Printf("    \"action\": \"%s\",\n", workflow.Action)
		fmt.Printf("    \"condition\": \"%s\",\n", workflow.Condition)
		fmt.Printf("    \"status\": \"%s\",\n", workflow.Status)
		fmt.Printf("    \"created_at\": \"%s\",\n", workflow.CreatedAt)
		fmt.Printf("    \"updated_at\": \"%s\"\n", workflow.UpdatedAt)
		fmt.Printf("  }")
		if i < len(workflows)-1 {
			fmt.Printf(",")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("]\n")
	return nil
}

func outputWorkflowStatusTable(status *service.WorkflowStatus) error {
	fmt.Printf("Workflow Status for Project %s:\n\n", status.ProjectID)
	fmt.Printf("Total Workflows: %d\n", status.TotalWorkflows)
	fmt.Printf("Active Workflows: %d\n", status.ActiveWorkflows)
	fmt.Printf("Total Executions: %d\n", status.TotalExecutions)
	fmt.Printf("Success Rate: %.1f%%\n\n", status.SuccessRate)

	if len(status.RecentExecutions) > 0 {
		fmt.Printf("Recent Executions:\n")
		fmt.Printf("%-20s %-30s %-15s %-10s %-20s\n", "WORKFLOW", "TRIGGER", "STATUS", "DURATION", "EXECUTED AT")
		fmt.Printf("%-20s %-30s %-15s %-10s %-20s\n", "────────", "───────", "──────", "────────", "───────────")

		for _, execution := range status.RecentExecutions {
			fmt.Printf("%-20s %-30s %-15s %-10s %-20s\n",
				truncateString(execution.WorkflowName, 20),
				truncateString(execution.Trigger, 30),
				execution.Status,
				execution.Duration,
				execution.ExecutedAt,
			)
		}
	}

	return nil
}

func outputWorkflowStatusJSON(status *service.WorkflowStatus) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"project_id\": \"%s\",\n", status.ProjectID)
	fmt.Printf("  \"total_workflows\": %d,\n", status.TotalWorkflows)
	fmt.Printf("  \"active_workflows\": %d,\n", status.ActiveWorkflows)
	fmt.Printf("  \"total_executions\": %d,\n", status.TotalExecutions)
	fmt.Printf("  \"success_rate\": %.1f,\n", status.SuccessRate)
	fmt.Printf("  \"recent_executions\": [\n")

	for i, execution := range status.RecentExecutions {
		fmt.Printf("    {\n")
		fmt.Printf("      \"workflow_name\": \"%s\",\n", execution.WorkflowName)
		fmt.Printf("      \"trigger\": \"%s\",\n", execution.Trigger)
		fmt.Printf("      \"status\": \"%s\",\n", execution.Status)
		fmt.Printf("      \"duration\": \"%s\",\n", execution.Duration)
		fmt.Printf("      \"executed_at\": \"%s\"\n", execution.ExecutedAt)
		fmt.Printf("    }")
		if i < len(status.RecentExecutions)-1 {
			fmt.Printf(",")
		}
		fmt.Printf("\n")
	}

	fmt.Printf("  ]\n")
	fmt.Printf("}\n")
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
