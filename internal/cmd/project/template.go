package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// TemplateOptions holds options for template commands
type TemplateOptions struct {
	Format string
}

// NewTemplateCmd creates the template command group
func NewTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage project templates",
		Long: `Manage GitHub Project templates for quick project setup.

Templates allow you to create reusable project configurations including:
• Custom fields and their configurations
• Default views and layouts  
• Workflow automation rules
• Standard project settings

Subcommands:
  list    - List available project templates
  create  - Create a new project template
  apply   - Apply a template to create a new project
  update  - Update an existing template
  delete  - Delete a template
  export  - Export template configuration
  import  - Import template from file`,
	}

	cmd.AddCommand(
		NewTemplateListCmd(),
		NewTemplateCreateCmd(),
		NewTemplateApplyCmd(),
		NewTemplateUpdateCmd(),
		NewTemplateDeleteCmd(),
		NewTemplateExportCmd(),
		NewTemplateImportCmd(),
	)

	return cmd
}

// NewTemplateListCmd creates the template list command
func NewTemplateListCmd() *cobra.Command {
	opts := &TemplateOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List project templates",
		Long: `List available project templates.

Templates are reusable configurations that can be applied to new projects.

Examples:
  ghx project template list
  ghx project template list --format json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTemplateList(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

// NewTemplateCreateCmd creates the template create command
func NewTemplateCreateCmd() *cobra.Command {
	var (
		name        string
		description string
		projectID   string
		category    string
		tags        []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new project template",
		Long: `Create a new project template from an existing project.

This captures the project's configuration including fields, views, and workflows
as a reusable template.

Available Categories:
  development  - Software development projects
  marketing    - Marketing campaign projects  
  design       - Design and creative projects
  research     - Research and analysis projects
  general      - General purpose projects

Examples:
  # Create template from existing project
  ghx project template create --name "Sprint Planning" --project-id myorg/123 --category development
  
  # Create template with description and tags
  ghx project template create --name "Bug Tracking" --description "Template for bug tracking" --project-id myorg/456 --category development --tags bug,tracking,support`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" || projectID == "" {
				return fmt.Errorf("--name and --project-id are required")
			}

			return runTemplateCreate(cmd.Context(), TemplateCreateOptions{
				Name:        name,
				Description: description,
				ProjectID:   projectID,
				Category:    category,
				Tags:        tags,
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Template name (required)")
	cmd.Flags().StringVar(&description, "description", "", "Template description")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Source project ID (required)")
	cmd.Flags().StringVar(&category, "category", "general", "Template category")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Template tags (comma-separated)")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("project-id")

	return cmd
}

// NewTemplateApplyCmd creates the template apply command
func NewTemplateApplyCmd() *cobra.Command {
	var (
		templateID  string
		projectName string
		owner       string
		org         bool
		customize   bool
	)

	cmd := &cobra.Command{
		Use:   "apply TEMPLATE_ID",
		Short: "Apply a template to create a new project",
		Long: `Apply a project template to create a new project.

This creates a new project with all the configuration from the template
including fields, views, and workflow automation.

Examples:
  # Apply template to create new project
  ghx project template apply template_123 --name "Q1 Sprint Planning" --owner myorg
  
  # Apply template with customization prompt
  ghx project template apply template_456 --name "Bug Tracking" --owner myuser --customize`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateID = args[0]

			if projectName == "" || owner == "" {
				return fmt.Errorf("--name and --owner are required")
			}

			return runTemplateApply(cmd.Context(), TemplateApplyOptions{
				TemplateID:  templateID,
				ProjectName: projectName,
				Owner:       owner,
				Org:         org,
				Customize:   customize,
			})
		},
	}

	cmd.Flags().StringVar(&projectName, "name", "", "New project name (required)")
	cmd.Flags().StringVar(&owner, "owner", "", "Project owner (user or organization) (required)")
	cmd.Flags().BoolVar(&org, "org", false, "Owner is an organization")
	cmd.Flags().BoolVar(&customize, "customize", false, "Prompt for template customization")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("owner")

	return cmd
}

// NewTemplateUpdateCmd creates the template update command
func NewTemplateUpdateCmd() *cobra.Command {
	var (
		name        string
		description string
		category    string
		tags        []string
	)

	cmd := &cobra.Command{
		Use:   "update TEMPLATE_ID",
		Short: "Update an existing template",
		Long: `Update an existing project template.

You can modify template metadata like name, description, category, and tags.

Examples:
  ghx project template update template_123 --name "Updated Sprint Planning"
  ghx project template update template_456 --description "Enhanced bug tracking template" --tags bug,enhanced`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateID := args[0]

			return runTemplateUpdate(cmd.Context(), TemplateUpdateOptions{
				TemplateID:  templateID,
				Name:        name,
				Description: description,
				Category:    category,
				Tags:        tags,
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Update template name")
	cmd.Flags().StringVar(&description, "description", "", "Update template description")
	cmd.Flags().StringVar(&category, "category", "", "Update template category")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Update template tags (comma-separated)")

	return cmd
}

// NewTemplateDeleteCmd creates the template delete command
func NewTemplateDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete TEMPLATE_ID",
		Short: "Delete a project template",
		Long: `Delete a project template.

This permanently removes the template and its configuration.

Examples:
  ghx project template delete template_123
  ghx project template delete template_456 --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateID := args[0]
			return runTemplateDelete(cmd.Context(), templateID, force)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force deletion without confirmation")

	return cmd
}

// NewTemplateExportCmd creates the template export command
func NewTemplateExportCmd() *cobra.Command {
	var (
		output string
		format string
	)

	cmd := &cobra.Command{
		Use:   "export TEMPLATE_ID",
		Short: "Export template configuration",
		Long: `Export a project template configuration to a file.

This creates a backup file that can be shared or imported elsewhere.

Examples:
  ghx project template export template_123 --output sprint-template.json
  ghx project template export template_456 --output bug-template.yaml --format yaml`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateID := args[0]

			if output == "" {
				return fmt.Errorf("--output is required")
			}

			return runTemplateExport(cmd.Context(), TemplateExportOptions{
				TemplateID: templateID,
				Output:     output,
				Format:     format,
			})
		},
	}

	cmd.Flags().StringVar(&output, "output", "", "Output file path (required)")
	cmd.Flags().StringVar(&format, "format", "json", "Export format: json, yaml")

	_ = cmd.MarkFlagRequired("output")

	return cmd
}

// NewTemplateImportCmd creates the template import command
func NewTemplateImportCmd() *cobra.Command {
	var (
		file   string
		name   string
		update bool
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import template from file",
		Long: `Import a project template from an exported configuration file.

This creates a new template from a previously exported template file.

Examples:
  ghx project template import --file sprint-template.json --name "Imported Sprint Template"
  ghx project template import --file bug-template.yaml --name "Bug Tracking" --update`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if file == "" || name == "" {
				return fmt.Errorf("--file and --name are required")
			}

			return runTemplateImport(cmd.Context(), TemplateImportOptions{
				File:   file,
				Name:   name,
				Update: update,
			})
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Template file to import (required)")
	cmd.Flags().StringVar(&name, "name", "", "Template name (required)")
	cmd.Flags().BoolVar(&update, "update", false, "Update existing template if it exists")

	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

// Command implementations

type TemplateCreateOptions struct {
	Name        string
	Description string
	ProjectID   string
	Category    string
	Tags        []string
}

type TemplateApplyOptions struct {
	TemplateID  string
	ProjectName string
	Owner       string
	Org         bool
	Customize   bool
}

type TemplateUpdateOptions struct {
	TemplateID  string
	Name        string
	Description string
	Category    string
	Tags        []string
}

type TemplateExportOptions struct {
	TemplateID string
	Output     string
	Format     string
}

type TemplateImportOptions struct {
	File   string
	Name   string
	Update bool
}

func runTemplateList(ctx context.Context, opts *TemplateOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	templateService := service.NewTemplateService(client)

	templates, err := templateService.ListTemplates(ctx)
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No project templates found")
		return nil
	}

	switch opts.Format {
	case "json":
		return outputTemplatesJSON(templates)
	case "table":
		return outputTemplatesTable(templates)
	default:
		return fmt.Errorf("unknown format: %s", opts.Format)
	}
}

func runTemplateCreate(ctx context.Context, opts TemplateCreateOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	templateService := service.NewTemplateService(client)

	template, err := templateService.CreateTemplate(ctx, service.CreateTemplateInput{
		Name:        opts.Name,
		Description: opts.Description,
		ProjectID:   opts.ProjectID,
		Category:    opts.Category,
		Tags:        opts.Tags,
	})
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	fmt.Printf("✅ Template '%s' created successfully\n\n", template.Name)
	fmt.Printf("Template Details:\n")
	fmt.Printf("  ID: %s\n", template.ID)
	fmt.Printf("  Name: %s\n", template.Name)
	fmt.Printf("  Description: %s\n", template.Description)
	fmt.Printf("  Category: %s\n", template.Category)
	if len(template.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", joinStrings(template.Tags, ", "))
	}
	fmt.Printf("  Fields: %d\n", len(template.Fields))
	fmt.Printf("  Views: %d\n", len(template.Views))
	fmt.Printf("  Workflows: %d\n", len(template.Workflows))

	return nil
}

func runTemplateApply(ctx context.Context, opts TemplateApplyOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	templateService := service.NewTemplateService(client)

	project, err := templateService.ApplyTemplate(ctx, service.ApplyTemplateInput{
		TemplateID:  opts.TemplateID,
		ProjectName: opts.ProjectName,
		Owner:       opts.Owner,
		Org:         opts.Org,
		Customize:   opts.Customize,
	})
	if err != nil {
		return fmt.Errorf("failed to apply template: %w", err)
	}

	fmt.Printf("✅ Template applied successfully\n\n")
	fmt.Printf("New Project Details:\n")
	fmt.Printf("  ID: %s\n", project.ID)
	fmt.Printf("  Name: %s\n", project.Name)
	fmt.Printf("  Owner: %s\n", project.Owner)
	fmt.Printf("  URL: %s\n", project.URL)
	fmt.Printf("  Fields Created: %d\n", project.FieldCount)
	fmt.Printf("  Views Created: %d\n", project.ViewCount)
	fmt.Printf("  Workflows Created: %d\n", project.WorkflowCount)

	return nil
}

func runTemplateUpdate(ctx context.Context, opts TemplateUpdateOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	templateService := service.NewTemplateService(client)

	template, err := templateService.UpdateTemplate(ctx, service.UpdateTemplateInput{
		TemplateID:  opts.TemplateID,
		Name:        opts.Name,
		Description: opts.Description,
		Category:    opts.Category,
		Tags:        opts.Tags,
	})
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	fmt.Printf("✅ Template updated successfully\n\n")
	fmt.Printf("  ID: %s\n", template.ID)
	fmt.Printf("  Name: %s\n", template.Name)
	fmt.Printf("  Description: %s\n", template.Description)
	fmt.Printf("  Category: %s\n", template.Category)

	return nil
}

func runTemplateDelete(ctx context.Context, templateID string, force bool) error {
	if !force {
		fmt.Printf("Are you sure you want to delete template %s? This action cannot be undone. (y/N): ", templateID)
		var response string
		_, _ = fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Template deletion canceled")
			return nil
		}
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	templateService := service.NewTemplateService(client)

	err = templateService.DeleteTemplate(ctx, templateID)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	fmt.Printf("✅ Template %s deleted successfully\n", templateID)
	return nil
}

func runTemplateExport(ctx context.Context, opts TemplateExportOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	templateService := service.NewTemplateService(client)

	err = templateService.ExportTemplate(ctx, service.ExportTemplateInput{
		TemplateID: opts.TemplateID,
		Output:     opts.Output,
		Format:     opts.Format,
	})
	if err != nil {
		return fmt.Errorf("failed to export template: %w", err)
	}

	fmt.Printf("✅ Template exported successfully to %s\n", opts.Output)
	return nil
}

func runTemplateImport(ctx context.Context, opts TemplateImportOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	templateService := service.NewTemplateService(client)

	template, err := templateService.ImportTemplate(ctx, service.ImportTemplateInput{
		File:   opts.File,
		Name:   opts.Name,
		Update: opts.Update,
	})
	if err != nil {
		return fmt.Errorf("failed to import template: %w", err)
	}

	fmt.Printf("✅ Template imported successfully\n\n")
	fmt.Printf("  ID: %s\n", template.ID)
	fmt.Printf("  Name: %s\n", template.Name)
	fmt.Printf("  Category: %s\n", template.Category)

	return nil
}

// Output functions
func outputTemplatesTable(templates []service.TemplateInfo) error {
	fmt.Printf("Project Templates:\n\n")
	fmt.Printf("%-15s %-30s %-15s %-10s %-20s\n", "ID", "NAME", "CATEGORY", "FIELDS", "CREATED")
	fmt.Printf("%-15s %-30s %-15s %-10s %-20s\n", "──", "────", "────────", "──────", "───────")

	for _, template := range templates {
		fmt.Printf("%-15s %-30s %-15s %-10d %-20s\n",
			truncateString(template.ID, 15),
			truncateString(template.Name, 30),
			truncateString(template.Category, 15),
			len(template.Fields),
			template.CreatedAt,
		)
	}

	fmt.Printf("\n%d templates total\n", len(templates))
	return nil
}

func outputTemplatesJSON(templates []service.TemplateInfo) error {
	fmt.Printf("[\n")
	for i, template := range templates {
		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": \"%s\",\n", template.ID)
		fmt.Printf("    \"name\": \"%s\",\n", template.Name)
		fmt.Printf("    \"description\": \"%s\",\n", template.Description)
		fmt.Printf("    \"category\": \"%s\",\n", template.Category)
		fmt.Printf("    \"tags\": [%s],\n", formatTagsJSON(template.Tags))
		fmt.Printf("    \"field_count\": %d,\n", len(template.Fields))
		fmt.Printf("    \"view_count\": %d,\n", len(template.Views))
		fmt.Printf("    \"workflow_count\": %d,\n", len(template.Workflows))
		fmt.Printf("    \"created_at\": \"%s\",\n", template.CreatedAt)
		fmt.Printf("    \"updated_at\": \"%s\"\n", template.UpdatedAt)
		fmt.Printf("  }")
		if i < len(templates)-1 {
			fmt.Printf(",")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("]\n")
	return nil
}

// Helper functions
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}

func formatTagsJSON(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	result := "\"" + tags[0] + "\""
	for _, tag := range tags[1:] {
		result += ", \"" + tag + "\""
	}
	return result
}
