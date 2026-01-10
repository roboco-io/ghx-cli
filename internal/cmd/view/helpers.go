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

// outputViewDetailsTable outputs common view details in table format
func outputViewDetailsTable(view *graphql.ProjectV2View) {
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
}

// outputViewDetailsJSON outputs common view details in JSON format
func outputViewDetailsJSON(view *graphql.ProjectV2View) {
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
			fmt.Printf("      \"field\": \"%s\",\n", gb.Field.Name)
			fmt.Printf("      \"direction\": \"%s\"\n", gb.Direction)
			if i < len(view.GroupBy)-1 {
				fmt.Printf("    },\n")
			} else {
				fmt.Printf("    }\n")
			}
		}
		fmt.Printf("  ]")
	}

	if len(view.SortBy) > 0 {
		fmt.Printf(",\n  \"sortBy\": [\n")
		for i, sb := range view.SortBy {
			fmt.Printf("    {\n")
			fmt.Printf("      \"field\": \"%s\",\n", sb.Field.Name)
			fmt.Printf("      \"direction\": \"%s\"\n", sb.Direction)
			if i < len(view.SortBy)-1 {
				fmt.Printf("    },\n")
			} else {
				fmt.Printf("    }\n")
			}
		}
		fmt.Printf("  ]")
	}

	fmt.Printf("\n")
}

// ConfigurationOptions represents common options for view configuration commands
type ConfigurationOptions struct {
	ViewID    string
	FieldID   string
	Direction string
	Format    string
	Clear     bool
}

// ConfigurationConfig holds configuration for view configuration commands
type ConfigurationConfig struct {
	UpdateFunction func(ctx context.Context, service *service.ViewService, input interface{}) error
	CreateInput    func(viewID, fieldID string, direction graphql.ProjectV2ViewSortDirection) interface{}
	Use            string
	Short          string
	Long           string
	OperationType  string
}

// createViewConfigurationCmd creates a view configuration command with shared logic
func createViewConfigurationCmd(config *ConfigurationConfig) *cobra.Command {
	opts := &ConfigurationOptions{}

	cmd := &cobra.Command{
		Use:   config.Use,
		Short: config.Short,
		Long:  config.Long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ViewID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runViewConfiguration(cmd.Context(), opts, config)
		},
	}

	cmd.Flags().StringVar(&opts.FieldID, "field", "", "Field ID to "+config.OperationType+" by")
	cmd.Flags().StringVar(&opts.Direction, "direction", "asc", config.OperationType+" direction (asc, desc)")
	cmd.Flags().BoolVar(&opts.Clear, "clear", false, "Clear "+config.OperationType+" from the view")

	return cmd
}

func runViewConfiguration(ctx context.Context, opts *ConfigurationOptions, config *ConfigurationConfig) error {
	// Validate input
	if opts.Clear && opts.FieldID != "" {
		return fmt.Errorf("cannot use --clear with --field")
	}

	if !opts.Clear && opts.FieldID == "" {
		return fmt.Errorf("must specify --field or --clear")
	}

	// Validate direction if not clearing
	var direction graphql.ProjectV2ViewSortDirection
	if !opts.Clear {
		var err error
		direction, err = service.ValidateSortDirection(opts.Direction)
		if err != nil {
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

	// Create input and execute update
	input := config.CreateInput(opts.ViewID, opts.FieldID, direction)
	err = config.UpdateFunction(ctx, viewService, input)
	if err != nil {
		return fmt.Errorf("failed to update view %s: %w", config.OperationType, err)
	}

	// Output result using shared helper
	return outputViewConfigurationResult(ctx, opts.ViewID, config.OperationType, opts.Clear, opts.Format)
}

// ViewConfigurationResult handles common view configuration result output
func outputViewConfigurationResult(ctx context.Context, viewID, operationType string, cleared bool, format string) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	viewService := service.NewViewService(client)

	// Get updated view for output
	viewInfo, err := viewService.GetView(ctx, viewID)
	if err != nil {
		return fmt.Errorf("failed to get updated view: %w", err)
	}

	// Output result based on format
	switch format {
	case formatJSON:
		return outputConfigurationResultJSON(viewInfo, operationType, cleared)
	case formatTable:
		return outputConfigurationResultTable(viewInfo, operationType, cleared)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputConfigurationResultTable(viewInfo *service.ViewInfo, operationType string, cleared bool) error {
	if cleared {
		fmt.Printf("✅ %s cleared from view '%s'\n", operationType, viewInfo.Name)
	} else {
		fmt.Printf("✅ View '%s' %s configuration updated\n", viewInfo.Name, operationType)
	}

	fmt.Printf("\nView Details:\n")
	fmt.Printf("  ID: %s\n", viewInfo.ID)
	fmt.Printf("  Name: %s\n", viewInfo.Name)
	fmt.Printf("  Layout: %s\n", service.FormatViewLayout(viewInfo.Layout))
	return nil
}

func outputConfigurationResultJSON(viewInfo *service.ViewInfo, operationType string, cleared bool) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"success\": true,\n")
	fmt.Printf("  \"operation\": \"%s\",\n", operationType)
	fmt.Printf("  \"cleared\": %t,\n", cleared)
	fmt.Printf("  \"viewId\": \"%s\",\n", viewInfo.ID)
	fmt.Printf("  \"viewName\": \"%s\",\n", viewInfo.Name)
	fmt.Printf("  \"layout\": \"%s\"\n", service.FormatViewLayout(viewInfo.Layout))
	fmt.Printf("}\n")
	return nil
}

// createGroupCmd creates the group command using shared configuration
func createGroupCmd() *cobra.Command {
	return createViewOperationCmd("group", "Configure view grouping", groupLongDescription(),
		func(ctx context.Context, viewService *service.ViewService, input interface{}) error {
			return viewService.UpdateViewGroup(ctx, input.(service.UpdateViewGroupInput))
		},
		func(viewID, fieldID string, direction graphql.ProjectV2ViewSortDirection) interface{} {
			input := service.UpdateViewGroupInput{
				ViewID:    viewID,
				Direction: direction,
			}
			if fieldID != "" {
				input.GroupByID = &fieldID
			}
			return input
		})
}

// createSortCmd creates the sort command using shared configuration
func createSortCmd() *cobra.Command {
	return createViewOperationCmd("sort", "Configure view sorting", sortLongDescription(),
		func(ctx context.Context, viewService *service.ViewService, input interface{}) error {
			return viewService.UpdateViewSort(ctx, input.(service.UpdateViewSortInput))
		},
		func(viewID, fieldID string, direction graphql.ProjectV2ViewSortDirection) interface{} {
			input := service.UpdateViewSortInput{
				ViewID:    viewID,
				Direction: direction,
			}
			if fieldID != "" {
				input.SortByID = &fieldID
			}
			return input
		})
}

// createViewOperationCmd creates a view operation command with the given parameters
func createViewOperationCmd(operation, short, long string,
	updateFunc func(context.Context, *service.ViewService, interface{}) error,
	createInputFunc func(string, string, graphql.ProjectV2ViewSortDirection) interface{}) *cobra.Command {
	config := &ConfigurationConfig{
		UpdateFunction: updateFunc,
		CreateInput:    createInputFunc,
		Use:            operation + " <view-id>",
		Short:          short,
		OperationType:  operation,
		Long:           long,
	}
	return createViewConfigurationCmd(config)
}

// groupLongDescription returns the long description for the group command
func groupLongDescription() string {
	return `Configure grouping for a project view.

You can set the field to group by and the group direction. Use --clear to
remove grouping from the view. Grouping is particularly useful for board
and roadmap views.

Group Directions:
  asc, ascending    - Group in ascending order (A-Z, 1-9, oldest first)
  desc, descending  - Group in descending order (Z-A, 9-1, newest first)

Examples:
  ghx view group view-id --field status-field-id --direction asc
  ghx view group view-id --field assignee-field-id --direction desc
  ghx view group view-id --clear
  ghx view group view-id --field priority-field-id --direction desc --format json`
}

// sortLongDescription returns the long description for the sort command
func sortLongDescription() string {
	return `Configure sorting for a project view.

You can set the field to sort by and the sort direction. Use --clear to
remove sorting from the view.

Sort Directions:
  asc, ascending    - Sort in ascending order (A-Z, 1-9, oldest first)
  desc, descending  - Sort in descending order (Z-A, 9-1, newest first)

Examples:
  ghx view sort view-id --field priority-field-id --direction desc
  ghx view sort view-id --field status-field-id --direction asc
  ghx view sort view-id --clear
  ghx view sort view-id --field due-date-field-id --direction asc --format json`
}
