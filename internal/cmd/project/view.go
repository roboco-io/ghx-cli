package project

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// ViewOptions holds options for the view command
type ViewOptions struct {
	Owner  string
	Format string
	Number int
	Org    bool
	Fields bool
	Items  bool
	Web    bool
}

// NewViewCmd creates the view command
func NewViewCmd() *cobra.Command {
	opts := &ViewOptions{}

	cmd := &cobra.Command{
		Use:   "view {<owner>/<number> | <number>}",
		Short: "View a project",
		Long: `View details of a specific project.

Examples:
  ghx project view 123               # View project 123 in current repository context
  ghx project view octocat/123       # View project 123 owned by octocat
  ghx project view --org myorg/456   # View project 456 owned by organization myorg`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runView(cmd.Context(), opts, args)
		},
	}

	cmd.Flags().BoolVar(&opts.Org, "org", false, "Project belongs to an organization")
	cmd.Flags().StringVar(&opts.Format, "format", "details", "Output format: details, json")
	cmd.Flags().BoolVar(&opts.Fields, "fields", false, "Show project fields")
	cmd.Flags().BoolVar(&opts.Items, "items", false, "Show project items")
	cmd.Flags().BoolVar(&opts.Web, "web", false, "Open project in web browser")

	return cmd
}

func runView(ctx context.Context, opts *ViewOptions, args []string) error {
	// Parse project reference
	projectRef := args[0]
	var err error

	if strings.Contains(projectRef, "/") {
		opts.Owner, opts.Number, err = service.ParseProjectReference(projectRef)
		if err != nil {
			return fmt.Errorf("invalid project reference: %w", err)
		}
	} else {
		// Just a number, need to determine owner from context
		opts.Number, err = strconv.Atoi(projectRef)
		if err != nil {
			return fmt.Errorf("invalid project number: %s", projectRef)
		}

		// For now, require owner to be specified
		// In the future, we can infer from git context
		return fmt.Errorf("owner must be specified in format owner/number")
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)

	// Get project details
	project, err := projectService.GetProject(ctx, opts.Owner, opts.Number, opts.Org)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Open in web browser if requested
	if opts.Web {
		fmt.Printf("Opening project in browser: %s\n", project.URL)
		// In a real implementation, we'd use a library to open the browser
		return nil
	}

	// Output project details
	return outputProjectDetails(project, opts)
}

func outputProjectDetails(project *graphql.ProjectV2, opts *ViewOptions) error {
	switch opts.Format {
	case "json":
		return outputProjectDetailsJSON(project)
	case "details":
		return outputProjectDetailsTable(project, opts)
	default:
		return fmt.Errorf("unknown format: %s", opts.Format)
	}
}

func outputProjectDetailsTable(project *graphql.ProjectV2, opts *ViewOptions) error {
	// Basic project information
	fmt.Printf("Project #%d\n", project.Number)
	fmt.Printf("Title: %s\n", project.Title)

	if project.Description != nil {
		fmt.Printf("Description: %s\n", *project.Description)
	}

	fmt.Printf("URL: %s\n", project.URL)
	fmt.Printf("Owner: %s (%s)\n", project.Owner.Login, project.Owner.Type)

	state := "Open"
	if project.Closed {
		state = "Closed"
	}
	fmt.Printf("State: %s\n", state)

	fmt.Printf("Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", project.UpdatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("Items: %d\n", len(project.Items.Nodes))
	fmt.Printf("Fields: %d\n", len(project.Fields.Nodes))

	// Show fields if requested
	if opts.Fields && len(project.Fields.Nodes) > 0 {
		fmt.Printf("\nFields:\n")
		fmt.Printf("%-20s %-15s %-10s\n", "NAME", "TYPE", "OPTIONS")
		fmt.Println(strings.Repeat("-", fieldsTableWidth))

		for _, field := range project.Fields.Nodes {
			optionCount := len(field.Options.Nodes)
			optionsStr := ""
			if optionCount > 0 {
				optionsStr = fmt.Sprintf("%d options", optionCount)
			}

			fmt.Printf("%-20s %-15s %-10s\n", field.Name, field.DataType, optionsStr)
		}
	}

	// Show items if requested
	if opts.Items && len(project.Items.Nodes) > 0 {
		fmt.Printf("\nItems:\n")
		fmt.Printf("%-10s %-30s %-10s %-15s\n", "TYPE", "TITLE", "STATE", "URL")
		fmt.Println(strings.Repeat("-", itemsTableWidth))

		for i := range project.Items.Nodes {
			item := &project.Items.Nodes[i]
			var title, state, url string

			switch item.Content.TypeName {
			case "Issue":
				title = item.Content.IssueTitle
				state = item.Content.IssueState
				url = item.Content.IssueURL
			case "PullRequest":
				title = item.Content.PRTitle
				state = item.Content.PRState
				url = item.Content.PRURL
			case "DraftIssue":
				title = item.Content.DraftTitle
				state = "Draft"
				url = "-"
			default:
				title = "Unknown"
				state = "-"
				url = "-"
			}

			if len(title) > titleMaxLength {
				title = title[:25] + "..."
			}

			fmt.Printf("%-10s %-30s %-10s %-15s\n",
				item.Content.TypeName, title, state, url)
		}
	}

	return nil
}

func outputProjectDetailsJSON(project *graphql.ProjectV2) error {
	// Simplified JSON output
	state := "open"
	if project.Closed {
		state = "closed"
	}

	description := "null"
	if project.Description != nil {
		description = fmt.Sprintf("%q", *project.Description)
	}

	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", project.ID)
	fmt.Printf("  \"number\": %d,\n", project.Number)
	fmt.Printf("  \"title\": \"%s\",\n", project.Title)
	fmt.Printf("  \"description\": %s,\n", description)
	fmt.Printf("  \"url\": \"%s\",\n", project.URL)
	fmt.Printf("  \"state\": \"%s\",\n", state)
	fmt.Printf("  \"owner\": {\n")
	fmt.Printf("    \"id\": \"%s\",\n", project.Owner.ID)
	fmt.Printf("    \"login\": \"%s\",\n", project.Owner.Login)
	fmt.Printf("    \"type\": \"%s\"\n", project.Owner.Type)
	fmt.Printf("  },\n")
	fmt.Printf("  \"createdAt\": \"%s\",\n", project.CreatedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("  \"updatedAt\": \"%s\",\n", project.UpdatedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Printf("  \"itemCount\": %d,\n", len(project.Items.Nodes))
	fmt.Printf("  \"fieldCount\": %d\n", len(project.Fields.Nodes))
	fmt.Printf("}\n")

	return nil
}
