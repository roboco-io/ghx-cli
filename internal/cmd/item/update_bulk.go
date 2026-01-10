package item

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// NewUpdateBulkCmd creates the update-bulk command
func NewUpdateBulkCmd() *cobra.Command {
	var (
		filter    string
		items     string
		fieldName string
		value     string
	)

	cmd := &cobra.Command{
		Use:   "update-bulk PROJECT_ID",
		Short: "Update multiple project items in bulk",
		Long: `Update field values for multiple project items in bulk.

This command allows you to update the same field for multiple items at once using:
• Filter by label or other criteria
• Item number range

Examples:
  # Update all items with specific label
  ghx item update-bulk myorg/123 --filter "label:epic" --field "Status" --value "Todo"
  
  # Update items by number range
  ghx item update-bulk myorg/123 --items 34-46 --field "Status" --value "In Progress"
  
  # Update all items matching a filter
  ghx item update-bulk myorg/123 --filter "assignee:@me" --field "Priority" --value "High"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateBulk(cmd.Context(), args[0], filter, items, fieldName, value)
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "Filter items to update (e.g., 'label:epic')")
	cmd.Flags().StringVar(&items, "items", "", "Item number range (e.g., 34-46)")
	cmd.Flags().StringVar(&fieldName, "field", "", "Field name to update")
	cmd.Flags().StringVar(&value, "value", "", "Value to set for the field")

	_ = cmd.MarkFlagRequired("field")
	_ = cmd.MarkFlagRequired("value")

	return cmd
}

func runUpdateBulk(ctx context.Context, projectRef, filter, items, fieldName, value string) error {
	// Validate required flags
	if fieldName == "" || value == "" {
		return fmt.Errorf("--field and --value are required")
	}

	if filter == "" && items == "" {
		return fmt.Errorf("either --filter or --items must be specified")
	}

	// Parse project reference (owner/number)
	parts := strings.Split(projectRef, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid project format: %s (expected owner/number)", projectRef)
	}

	owner := parts[0]
	projectNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid project number: %s", parts[1])
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	itemService := service.NewItemService(client)
	projectService := service.NewProjectService(client)

	// Get project to get project ID
	project, err := projectService.GetProjectWithOwnerDetection(ctx, owner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	var itemsToUpdate []string

	// Handle filter
	if filter != "" {
		filtered, filterErr := itemService.GetItemsByFilter(ctx, project.ID, filter)
		if filterErr != nil {
			return fmt.Errorf("failed to get items by filter: %w", filterErr)
		}
		itemsToUpdate = append(itemsToUpdate, filtered...)
	}

	// Handle item range
	if items != "" {
		itemRange, rangeErr := service.ParseNumberRange(items)
		if rangeErr != nil {
			return fmt.Errorf("invalid item range: %w", rangeErr)
		}
		itemsToUpdate = append(itemsToUpdate, itemRange...)
	}

	// Remove duplicates
	itemsToUpdate = service.RemoveDuplicates(itemsToUpdate)

	fmt.Printf("Updating %d items in project %s...\n", len(itemsToUpdate), projectRef)
	fmt.Printf("Setting field '%s' to '%s'\n\n", fieldName, value)

	// Update items using service
	input := service.BulkUpdateInput{
		ProjectID: project.ID,
		ItemIDs:   itemsToUpdate,
		FieldName: fieldName,
		Value:     value,
	}

	result, err := itemService.BulkUpdateItems(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update items: %w", err)
	}

	fmt.Printf("\n✅ Successfully updated %d items", result.Updated)
	if result.Failed > 0 {
		fmt.Printf(" (%d failed)", result.Failed)
		for _, errMsg := range result.Errors {
			fmt.Printf("\n  Error: %s", errMsg)
		}
	}
	fmt.Printf("\n")

	return nil
}
