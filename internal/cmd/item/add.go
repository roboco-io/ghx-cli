package item

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

const (
	maxDisplayTitleLength  = 80
	addTitleTruncateLength = 77
)

// AddOptions holds options for the add command
type AddOptions struct {
	ProjectRef string
	ItemRef    string
	Title      string
	Body       string
	Format     string
	Draft      bool
}

// NewAddCmd creates the add command
func NewAddCmd() *cobra.Command {
	opts := &AddOptions{}

	cmd := &cobra.Command{
		Use:   "add <project> <item>",
		Short: "Add an item to a project",
		Long: `Add an existing issue, pull request, or create a draft issue in a project.

Item references can be in the following formats:
• owner/repo#123 (issue or PR reference)
• https://github.com/owner/repo/issues/123 (GitHub issue URL)
• https://github.com/owner/repo/pull/456 (GitHub PR URL)

Project references should be in owner/number format (e.g., octocat/1).

Examples:
  ghx item add octocat/1 octocat/Hello-World#123     # Add issue to project
  ghx item add myorg/2 myorg/repo#456 --format json  # Add PR with JSON output
  ghx item add octocat/1 --draft --title "New task"  # Create draft issue`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			if len(args) > 1 {
				opts.ItemRef = args[1]
			}
			return runAdd(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Draft, "draft", false, "Create a draft issue instead of adding existing item")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Title for draft issue (required when --draft is used)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Body for draft issue")
	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

func validateAddOptions(opts *AddOptions) error {
	if opts.Draft {
		if opts.Title == "" {
			return fmt.Errorf("title is required when creating draft issue (use --title)")
		}
		if opts.ItemRef != "" {
			return fmt.Errorf("cannot specify both --draft and item reference")
		}
	} else if opts.ItemRef == "" {
		return fmt.Errorf("item reference is required (or use --draft to create draft issue)")
	}
	return nil
}

func setupAddServices(_ context.Context) (*api.Client, *service.ItemService, *service.ProjectService, error) {
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("authentication failed: %w", err)
	}

	client := api.NewClient(token)
	itemService := service.NewItemService(client)
	projectService := service.NewProjectService(client)

	return client, itemService, projectService, nil
}

func addDraftIssue(ctx context.Context, itemService *service.ItemService, projectID, title string, body *string, format string) error {
	item, err := itemService.CreateDraftIssue(ctx, projectID, title, body)
	if err != nil {
		return fmt.Errorf("failed to create draft issue: %w", err)
	}

	fmt.Printf("✅ Draft issue created and added to project!\n\n")
	return outputAddedItem(item, format, "DraftIssue", title)
}

func addExistingItem(ctx context.Context, itemService *service.ItemService, projectID, itemRef, format string) error {
	itemOwner, itemRepo, itemNumber, err := service.ParseItemReference(itemRef)
	if err != nil {
		return fmt.Errorf("invalid item reference: %w", err)
	}

	var contentID, itemType, itemTitle string

	// Try to get as issue first, then as PR
	issue, err := itemService.GetIssue(ctx, itemOwner, itemRepo, itemNumber)
	if err == nil {
		contentID = issue.ID
		itemType = "Issue"
		itemTitle = issue.Title
	} else {
		pr, prErr := itemService.GetPullRequest(ctx, itemOwner, itemRepo, itemNumber)
		if prErr != nil {
			return fmt.Errorf("failed to find issue or pull request: %w", prErr)
		}
		contentID = pr.ID
		itemType = "PullRequest"
		itemTitle = pr.Title
	}

	item, err := itemService.AddItemToProject(ctx, projectID, contentID)
	if err != nil {
		return fmt.Errorf("failed to add item to project: %w", err)
	}

	fmt.Printf("✅ %s added to project!\n\n", itemType)
	return outputAddedItem(item, format, itemType, itemTitle)
}

func runAdd(ctx context.Context, opts *AddOptions) error {
	if err := validateAddOptions(opts); err != nil {
		return err
	}

	// Parse project reference
	projectOwner, projectNumber, err := service.ParseProjectReference(opts.ProjectRef)
	if err != nil {
		return fmt.Errorf("invalid project reference: %w", err)
	}

	// Setup services
	_, itemService, projectService, err := setupAddServices(ctx)
	if err != nil {
		return err
	}

	// Get project details to obtain project ID
	project, err := projectService.GetProjectWithOwnerDetection(ctx, projectOwner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	if opts.Draft {
		var body *string
		if opts.Body != "" {
			body = &opts.Body
		}
		return addDraftIssue(ctx, itemService, project.ID, opts.Title, body, opts.Format)
	}

	return addExistingItem(ctx, itemService, project.ID, opts.ItemRef, opts.Format)
}

func outputAddedItem(item interface{}, format, itemType, title string) error {
	switch format {
	case formatJSON:
		return outputAddedItemJSON(item)
	case formatTable:
		return outputAddedItemTable(itemType, title)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputAddedItemTable(itemType, title string) error {
	fmt.Printf("Type: %s\n", itemType)

	if title != "" {
		displayTitle := title
		if len(displayTitle) > maxDisplayTitleLength {
			displayTitle = displayTitle[:addTitleTruncateLength] + "..."
		}
		fmt.Printf("Title: %s\n", displayTitle)
	}

	return nil
}

func outputAddedItemJSON(_ interface{}) error {
	// In a real implementation, we'd properly serialize the item
	fmt.Printf("{\n")
	fmt.Printf("  \"status\": \"added\",\n")
	fmt.Printf("  \"item\": {\n")
	fmt.Printf("    \"id\": \"<item-id>\",\n")
	fmt.Printf("    \"type\": \"<item-type>\"\n")
	fmt.Printf("  }\n")
	fmt.Printf("}\n")
	return nil
}
