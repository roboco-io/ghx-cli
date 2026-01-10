package analytics

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

const (
	bulkPercentageMultiplier = 100.0
)

// BulkUpdateOptions holds options for the bulk-update command
type BulkUpdateOptions struct {
	Updates    map[string]interface{}
	ProjectRef string
	Format     string
	ItemIDs    []string
}

// NewBulkUpdateCmd creates the bulk-update command
func NewBulkUpdateCmd() *cobra.Command {
	opts := &BulkUpdateOptions{
		Updates: make(map[string]interface{}),
	}

	cmd := &cobra.Command{
		Use:   "bulk-update <owner/project-number>",
		Short: "Bulk update project items",
		Long: `Perform bulk updates on multiple project items simultaneously.

This command allows you to update multiple project items at once, which is
much more efficient than updating items individually. You can update any
field values including status, assignees, labels, milestones, and custom fields.

The bulk update operation runs asynchronously, so you can check the status
using the 'operation-status' command with the operation ID returned.

Update Options:
  --items              Comma-separated list of item IDs to update
  --field-<name>       Set field value (e.g., --field-status Done, --field-priority High)
  --status             Set status field value
  --assignee           Set assignee field value
  --labels             Set labels field value (comma-separated)
  --milestone          Set milestone field value

Examples:
  ghx analytics bulk-update octocat/123 --items item1,item2,item3 --status Done
  ghx analytics bulk-update octocat/123 --items item1,item2 --assignee octocat --priority High
  ghx analytics bulk-update octocat/123 --items item1,item2,item3 --labels bug,urgent --format json
  ghx analytics bulk-update --org myorg/456 --items item1,item2 --field-custom-field "Custom Value"`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Format = cmd.Flag("format").Value.String()

			// Get item IDs
			if itemsStr, _ := cmd.Flags().GetString("items"); itemsStr != "" {
				opts.ItemIDs = strings.Split(itemsStr, ",")
				for i, id := range opts.ItemIDs {
					opts.ItemIDs[i] = strings.TrimSpace(id)
				}
			}

			// Get field updates
			if status, _ := cmd.Flags().GetString("status"); status != "" {
				opts.Updates["status"] = status
			}
			if assignee, _ := cmd.Flags().GetString("assignee"); assignee != "" {
				opts.Updates["assignee"] = assignee
			}
			if labels, _ := cmd.Flags().GetString("labels"); labels != "" {
				labelList := strings.Split(labels, ",")
				for i, label := range labelList {
					labelList[i] = strings.TrimSpace(label)
				}
				opts.Updates["labels"] = labelList
			}
			if milestone, _ := cmd.Flags().GetString("milestone"); milestone != "" {
				opts.Updates["milestone"] = milestone
			}
			if priority, _ := cmd.Flags().GetString("priority"); priority != "" {
				opts.Updates["priority"] = priority
			}

			return runBulkUpdate(cmd.Context(), opts)
		},
	}

	cmd.Flags().String("items", "", "Comma-separated list of item IDs")
	cmd.Flags().String("status", "", "Status field value")
	cmd.Flags().String("assignee", "", "Assignee field value")
	cmd.Flags().String("labels", "", "Labels field value (comma-separated)")
	cmd.Flags().String("milestone", "", "Milestone field value")
	cmd.Flags().String("priority", "", "Priority field value")
	cmd.Flags().Bool("org", false, "Target organization project")

	// Make items flag required
	_ = cmd.MarkFlagRequired("items")

	return cmd
}

func runBulkUpdate(ctx context.Context, opts *BulkUpdateOptions) error {
	// Validate input
	if len(opts.ItemIDs) == 0 {
		return fmt.Errorf("no items specified for update")
	}

	if len(opts.Updates) == 0 {
		return fmt.Errorf("no field updates specified")
	}

	// Parse project reference
	parts := strings.Split(opts.ProjectRef, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid project reference format. Use: owner/project-number")
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
	projectService := service.NewProjectService(client)
	analyticsService := service.NewAnalyticsService(client)

	// Get project to validate access and get project ID (with automatic owner detection)
	project, err := projectService.GetProjectWithOwnerDetection(ctx, owner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Prepare bulk update input
	input := service.BulkUpdateItemsInput{
		ProjectID: project.ID,
		ItemIDs:   opts.ItemIDs,
		Updates:   opts.Updates,
	}

	// Start bulk update operation
	operation, err := analyticsService.BulkUpdateItems(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to start bulk update: %w", err)
	}

	// Output operation result
	return outputBulkOperation(operation, "update", opts.Format)
}

func outputBulkOperation(operation *service.BulkOperation, operationType, format string) error {
	switch format {
	case FormatJSON:
		return outputBulkOperationJSON(operation)
	case FormatTable:
		return outputBulkOperationTable(operation, operationType)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputBulkOperationTable(operation *service.BulkOperation, operationType string) error {
	fmt.Printf("✅ Bulk %s operation started successfully\n\n", operationType)

	fmt.Printf("Operation Details:\n")
	fmt.Printf("  Operation ID: %s\n", operation.ID)
	fmt.Printf("  Type: %s\n", service.FormatBulkOperationType(operation.Type))
	fmt.Printf("  Status: %s\n", service.FormatBulkOperationStatus(operation.Status))
	fmt.Printf("  Progress: %.1f%%\n", operation.Progress*bulkPercentageMultiplier)
	fmt.Printf("  Total Items: %d\n", operation.TotalItems)
	fmt.Printf("  Processed Items: %d\n", operation.ProcessedItems)
	if operation.FailedItems > 0 {
		fmt.Printf("  Failed Items: %d\n", operation.FailedItems)
	}
	fmt.Printf("  Created At: %s\n", operation.CreatedAt.Format("2006-01-02 15:04:05"))

	if operation.CompletedAt != nil {
		fmt.Printf("  Completed At: %s\n", operation.CompletedAt.Format("2006-01-02 15:04:05"))
	}

	if operation.ErrorMessage != nil {
		fmt.Printf("  Error: %s\n", *operation.ErrorMessage)
	}

	fmt.Printf("\n💡 Use 'ghp analytics operation-status %s' to check operation progress\n", operation.ID)

	return nil
}

func outputBulkOperationJSON(operation *service.BulkOperation) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"success\": true,\n")
	fmt.Printf("  \"operationId\": \"%s\",\n", operation.ID)
	fmt.Printf("  \"type\": \"%s\",\n", operation.Type)
	fmt.Printf("  \"status\": \"%s\",\n", operation.Status)
	fmt.Printf("  \"progress\": %.3f,\n", operation.Progress)
	fmt.Printf("  \"totalItems\": %d,\n", operation.TotalItems)
	fmt.Printf("  \"processedItems\": %d,\n", operation.ProcessedItems)
	fmt.Printf("  \"failedItems\": %d,\n", operation.FailedItems)
	fmt.Printf("  \"createdAt\": \"%s\"", operation.CreatedAt.Format("2006-01-02T15:04:05Z"))

	if operation.CompletedAt != nil {
		fmt.Printf(",\n  \"completedAt\": \"%s\"", operation.CompletedAt.Format("2006-01-02T15:04:05Z"))
	}

	if operation.ErrorMessage != nil {
		fmt.Printf(",\n  \"errorMessage\": \"%s\"", *operation.ErrorMessage)
	}

	fmt.Printf("\n}\n")
	return nil
}
