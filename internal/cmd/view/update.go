package view

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// UpdateOptions holds options for the update command
type UpdateOptions struct {
	ViewID string
	Name   string
	Filter string
	Format string
}

// NewUpdateCmd creates the update command
func NewUpdateCmd() *cobra.Command {
	opts := &UpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <view-id>",
		Short: "Update a project view",
		Long: `Update properties of an existing project view.

You can update the view name and filter. At least one property must be specified.
The view layout cannot be changed after creation.

Examples:
  ghx view update view-id --name "Updated Dashboard"
  ghx view update view-id --filter "status:todo"
  ghx view update view-id --name "Sprint Board" --filter "milestone:sprint-1"
  ghx view update view-id --name "Bug Tracking" --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ViewID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runUpdate(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "New name for the view")
	cmd.Flags().StringVar(&opts.Filter, "filter", "", "Filter expression for the view")

	return cmd
}

func runUpdate(ctx context.Context, opts *UpdateOptions) error {
	// Validate at least one field is provided
	if opts.Name == "" && opts.Filter == "" {
		return fmt.Errorf("at least one of --name or --filter must be provided")
	}

	// Validate name if provided
	if opts.Name != "" {
		if err := service.ValidateViewName(opts.Name); err != nil {
			return err
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
	viewService := service.NewViewService(client)

	// Prepare input
	input := service.UpdateViewInput{
		ViewID: opts.ViewID,
	}

	if opts.Name != "" {
		input.Name = &opts.Name
	}
	if opts.Filter != "" {
		input.Filter = &opts.Filter
	}

	// Update view
	view, err := viewService.UpdateView(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update view: %w", err)
	}

	// Output updated view
	return outputUpdatedView(view, opts.Format)
}

func outputUpdatedView(view *graphql.ProjectV2View, format string) error {
	switch format {
	case "json":
		return outputUpdatedViewJSON(view)
	case "table":
		return outputUpdatedViewTable(view)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputUpdatedViewTable(view *graphql.ProjectV2View) error {
	fmt.Printf("✅ View '%s' updated successfully\n\n", view.Name)

	fmt.Printf("View Details:\n")
	fmt.Printf("  ID: %s\n", view.ID)
	fmt.Printf("  Name: %s\n", view.Name)
	fmt.Printf("  Layout: %s\n", service.FormatViewLayout(view.Layout))
	fmt.Printf("  Number: %d\n", view.Number)

	if view.Filter != nil && *view.Filter != "" {
		fmt.Printf("  Filter: %s\n", *view.Filter)
	}

	if len(view.GroupBy) > 0 {
		fmt.Printf("  Group By:\n")
		for _, gb := range view.GroupBy {
			fmt.Printf("    - %s (%s)\n", gb.Field.Name, service.FormatSortDirection(gb.Direction))
		}
	}

	if len(view.SortBy) > 0 {
		fmt.Printf("  Sort By:\n")
		for _, sb := range view.SortBy {
			fmt.Printf("    - %s (%s)\n", sb.Field.Name, service.FormatSortDirection(sb.Direction))
		}
	}

	return nil
}

func outputUpdatedViewJSON(view *graphql.ProjectV2View) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"id\": \"%s\",\n", view.ID)
	fmt.Printf("  \"name\": \"%s\",\n", view.Name)
	fmt.Printf("  \"layout\": \"%s\",\n", view.Layout)
	fmt.Printf("  \"number\": %d", view.Number)

	if view.Filter != nil {
		fmt.Printf(",\n  \"filter\": \"%s\"", *view.Filter)
	}

	if len(view.GroupBy) > 0 {
		fmt.Printf(",\n  \"groupBy\": [\n")
		for i, gb := range view.GroupBy {
			fmt.Printf("    {\n")
			fmt.Printf("      \"fieldId\": \"%s\",\n", gb.Field.ID)
			fmt.Printf("      \"fieldName\": \"%s\",\n", gb.Field.Name)
			fmt.Printf("      \"direction\": \"%s\"\n", gb.Direction)
			fmt.Printf("    }")
			if i < len(view.GroupBy)-1 {
				fmt.Printf(",")
			}
			fmt.Printf("\n")
		}
		fmt.Printf("  ]")
	}

	if len(view.SortBy) > 0 {
		fmt.Printf(",\n  \"sortBy\": [\n")
		for i, sb := range view.SortBy {
			fmt.Printf("    {\n")
			fmt.Printf("      \"fieldId\": \"%s\",\n", sb.Field.ID)
			fmt.Printf("      \"fieldName\": \"%s\",\n", sb.Field.Name)
			fmt.Printf("      \"direction\": \"%s\"\n", sb.Direction)
			fmt.Printf("    }")
			if i < len(view.SortBy)-1 {
				fmt.Printf(",")
			}
			fmt.Printf("\n")
		}
		fmt.Printf("  ]")
	}

	fmt.Printf("\n}\n")

	return nil
}
