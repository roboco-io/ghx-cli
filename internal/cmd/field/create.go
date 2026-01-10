package field

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// CreateOptions holds options for the create command
type CreateOptions struct {
	ProjectRef string
	ProjectID  string
	Owner      string
	Name       string
	FieldType  string
	Format     string
	Options    []string
	Duration   string
	Number     int
	Org        bool
}

// NewCreateCmd creates the create command
func NewCreateCmd() *cobra.Command {
	opts := &CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create [owner/number] [name] [type]",
		Short: "Create a new project field",
		Long: `Create a new custom field in a GitHub Project.

Fields allow you to track additional metadata for your project items.
Different field types support different kinds of data:

Field Types:
  text         - Text field for arbitrary text input
  number       - Numeric field for numbers and calculations  
  date         - Date field for deadlines and milestones
  single_select - Single select field with predefined options
  iteration    - Iteration field for sprint/cycle planning

For single select fields, you can provide initial options using --options.
For iteration fields, you can specify duration using --duration.

Examples:
  # Traditional syntax
  ghx field create octocat/123 "Priority" text
  
  # New syntax with flags (Issue #18)
  ghx field create --project-id PROJECT_ID --name "Priority" --type single-select --options "Critical,High,Medium,Low"
  ghx field create --project-id PROJECT_ID --name "Sprint" --type iteration --duration 2w
  ghx field create octocat/123 "Story Points" number
  ghx field create octocat/123 "Due Date" date
  ghx field create octocat/123 "Status" single_select --options "Todo,In Progress,Done"
  ghx field create --org myorg/456 "Sprint" iteration`,

		Args: cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Support both traditional args and new flag-based syntax
			if len(args) == 3 {
				opts.ProjectRef = args[0]
				opts.Name = args[1]
				opts.FieldType = args[2]
			}
			opts.Format = cmd.Flag("format").Value.String()
			return runCreate(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Org, "org", false, "Project belongs to an organization")
	cmd.Flags().StringSliceVar(&opts.Options, "options", []string{}, "Options for single select field (comma-separated)")

	// New flags for Issue #18 syntax
	cmd.Flags().StringVar(&opts.ProjectID, "project-id", "", "Project ID (alternative to owner/number)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Field name")
	cmd.Flags().StringVar(&opts.FieldType, "type", "", "Field type (text, number, date, single_select, iteration)")
	cmd.Flags().StringVar(&opts.Duration, "duration", "", "Duration for iteration field (e.g., 2w, 1m)")

	return cmd
}

func runCreate(ctx context.Context, opts *CreateOptions) error {
	// Support both traditional args and new flag-based syntax
	var err error
	var projectID string

	if opts.ProjectID != "" {
		// New syntax: --project-id flag
		projectID = opts.ProjectID
		if opts.Name == "" {
			return fmt.Errorf("--name is required when using --project-id")
		}
		if opts.FieldType == "" {
			return fmt.Errorf("--type is required when using --project-id")
		}
	} else if opts.ProjectRef != "" {
		// Traditional syntax: positional args
		if strings.Contains(opts.ProjectRef, "/") {
			opts.Owner, opts.Number, err = service.ParseProjectReference(opts.ProjectRef)
			if err != nil {
				return fmt.Errorf("invalid project reference: %w", err)
			}
		} else {
			return fmt.Errorf("project reference must be in format owner/number")
		}
	} else {
		return fmt.Errorf("either project reference (owner/number) or --project-id must be provided")
	}

	// Validate field name
	if validateErr := service.ValidateFieldName(opts.Name); validateErr != nil {
		return validateErr
	}

	// Validate field type
	dataType, err := service.ValidateFieldType(opts.FieldType)
	if err != nil {
		return err
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	fieldService := service.NewFieldService(client)
	projectService := service.NewProjectService(client)

	var project *graphql.ProjectV2

	if opts.ProjectID != "" {
		// New syntax: use project ID directly
		// For new syntax, we'll create a mock project for output
		project = &graphql.ProjectV2{
			ID:    projectID,
			Title: fmt.Sprintf("Project %s", projectID),
		}
	} else {
		// Traditional syntax: get project by owner/number
		project, err = projectService.GetProject(ctx, opts.Owner, opts.Number, opts.Org)
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}
		projectID = project.ID
	}

	// Create field
	input := service.CreateFieldInput{
		ProjectID:           projectID,
		Name:                opts.Name,
		DataType:            dataType,
		SingleSelectOptions: opts.Options,
		Duration:            opts.Duration,
	}

	field, err := fieldService.CreateField(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create field: %w", err)
	}

	// Output created field
	return outputCreatedField(field, project.Title, opts.Format)
}

func outputCreatedField(field *graphql.ProjectV2Field, projectName, format string) error {
	switch format {
	case formatJSON:
		return outputCreatedFieldJSON(field)
	case formatTable:
		return outputCreatedFieldTable(field, projectName)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputCreatedFieldTable(field *graphql.ProjectV2Field, projectName string) error {
	fmt.Printf("✅ Field '%s' created successfully in project '%s'\n\n", field.Name, projectName)

	fmt.Printf("Field Details:\n")
	fmt.Printf("  ID: %s\n", field.ID)
	fmt.Printf("  Name: %s\n", field.Name)
	fmt.Printf("  Type: %s\n", service.FormatFieldDataType(field.DataType))

	// Show options if single select field
	if len(field.Options.Nodes) > 0 {
		fmt.Printf("  Options:\n")
		for _, option := range field.Options.Nodes {
			fmt.Printf("    • %s (%s)", option.Name, service.FormatColor(option.Color))
			if option.Description != nil && *option.Description != "" {
				fmt.Printf(" - %s", *option.Description)
			}
			fmt.Printf("\n")
		}
	}

	return nil
}

func outputCreatedFieldJSON(field *graphql.ProjectV2Field) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", field.ID)
	fmt.Printf("  \"name\": \"%s\",\n", field.Name)
	fmt.Printf("  \"dataType\": \"%s\"", field.DataType)

	if len(field.Options.Nodes) > 0 {
		fmt.Printf(",\n  \"options\": [\n")
		for i, option := range field.Options.Nodes {
			fmt.Printf("    {\n")
			fmt.Printf("      \"id\": \"%s\",\n", option.ID)
			fmt.Printf("      \"name\": \"%s\",\n", option.Name)
			fmt.Printf("      \"color\": \"%s\"", option.Color)
			if option.Description != nil {
				fmt.Printf(",\n      \"description\": \"%s\"", *option.Description)
			}
			fmt.Printf("\n    }")
			if i < len(field.Options.Nodes)-1 {
				fmt.Printf(",")
			}
			fmt.Printf("\n")
		}
		fmt.Printf("  ]\n")
	} else {
		fmt.Printf("\n")
	}
	fmt.Printf("}\n")

	return nil
}
