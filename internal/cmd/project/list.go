package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// ListOptions holds options for the list command
type ListOptions struct {
	Owner  string
	State  string
	Format string
	Limit  int
	Org    bool
	User   bool
}

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list [owner]",
		Short: "List projects",
		Long: `List projects for a user or organization.

Examples:
  ghx project list              # List projects for authenticated user
  ghx project list octocat      # List projects for user octocat
  ghx project list --org myorg  # List projects for organization myorg`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd.Context(), opts, args)
		},
	}

	cmd.Flags().BoolVar(&opts.Org, "org", false, "List organization projects")
	cmd.Flags().BoolVar(&opts.User, "user", false, "List user projects (default)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", defaultListLimit, "Maximum number of projects to list")
	cmd.Flags().StringVar(&opts.State, "state", "all", "Filter by state: open, closed, all")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

func runList(ctx context.Context, opts *ListOptions, args []string) error {
	// Determine owner
	if len(args) > 0 {
		opts.Owner = args[0]
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

	// Get current user if no owner specified
	if opts.Owner == "" {
		// For now, we'll require an owner to be specified
		// In the future, we can implement a viewer query to get the current user
		return fmt.Errorf("owner must be specified")
	}

	var projects []service.ProjectInfo

	// List projects based on type
	if opts.Org {
		listOpts := service.ListOrgProjectsOptions{
			Login: opts.Owner,
			First: opts.Limit,
		}
		projects, err = projectService.ListOrgProjects(ctx, listOpts)
	} else {
		listOpts := service.ListUserProjectsOptions{
			Login: opts.Owner,
			First: opts.Limit,
		}
		projects, err = projectService.ListUserProjects(ctx, listOpts)
	}

	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	// Filter by state if specified
	if opts.State != "all" {
		projects = filterProjectsByState(projects, opts.State)
	}

	// Output results
	return outputProjects(projects, opts.Format)
}

func filterProjectsByState(projects []service.ProjectInfo, state string) []service.ProjectInfo {
	var filtered []service.ProjectInfo

	for _, project := range projects {
		switch state {
		case statusOpen:
			if !project.Closed {
				filtered = append(filtered, project)
			}
		case statusClosed:
			if project.Closed {
				filtered = append(filtered, project)
			}
		default:
			filtered = append(filtered, project)
		}
	}

	return filtered
}

func outputProjects(projects []service.ProjectInfo, format string) error {
	if len(projects) == 0 {
		fmt.Println("No projects found")
		return nil
	}

	switch format {
	case "json":
		return outputProjectsJSON(projects)
	case "table":
		return outputProjectsTable(projects)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputProjectsTable(projects []service.ProjectInfo) error {
	// Print header
	fmt.Printf("%-8s %-30s %-10s %-15s %-10s %-10s\n",
		"NUMBER", "TITLE", "STATE", "OWNER", "ITEMS", "FIELDS")
	fmt.Println(strings.Repeat("-", tableSeparatorWidth))

	// Print projects
	for _, project := range projects {
		state := "OPEN"
		if project.Closed {
			state = "CLOSED"
		}

		title := project.Title
		if len(title) > titleMaxLength {
			title = title[:25] + "..."
		}

		owner := project.Owner
		if len(owner) > ownerMaxLength {
			owner = owner[:10] + "..."
		}

		fmt.Printf("%-8d %-30s %-10s %-15s %-10d %-10d\n",
			project.Number, title, state, owner,
			project.ItemCount, project.FieldCount)
	}

	return nil
}

func outputProjectsJSON(projects []service.ProjectInfo) error {
	// For now, simple JSON-like output
	// In production, we'd use proper JSON marshaling
	fmt.Println("[")
	for i, project := range projects {
		state := "open"
		if project.Closed {
			state = "closed"
		}

		description := "null"
		if project.Description != nil {
			description = fmt.Sprintf("%q", *project.Description)
		}

		fmt.Printf("  {\n")
		fmt.Printf("    \"id\": \"%s\",\n", project.ID)
		fmt.Printf("    \"number\": %d,\n", project.Number)
		fmt.Printf("    \"title\": \"%s\",\n", project.Title)
		fmt.Printf("    \"description\": %s,\n", description)
		fmt.Printf("    \"url\": \"%s\",\n", project.URL)
		fmt.Printf("    \"state\": \"%s\",\n", state)
		fmt.Printf("    \"owner\": \"%s\",\n", project.Owner)
		fmt.Printf("    \"itemCount\": %d,\n", project.ItemCount)
		fmt.Printf("    \"fieldCount\": %d\n", project.FieldCount)

		if i < len(projects)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")

	return nil
}
