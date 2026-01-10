package discussion

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/gh-project-cli/internal/api"
	"github.com/roboco-io/gh-project-cli/internal/auth"
	"github.com/roboco-io/gh-project-cli/internal/service"
)

// CategoryListOptions holds options for the category list command
type CategoryListOptions struct {
	Repo   string
	Format string
}

// NewCategoryCmd creates the category command group
func NewCategoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category <command>",
		Short: "Manage discussion categories",
		Long: `Manage discussion categories for a repository.

Note: GitHub only allows creating and modifying categories through the web UI.
This command can list existing categories.`,
		Example: `  ghp discussion category list owner/repo`,
	}

	cmd.AddCommand(NewCategoryListCmd())

	return cmd
}

// NewCategoryListCmd creates the category list command
func NewCategoryListCmd() *cobra.Command {
	opts := &CategoryListOptions{}

	cmd := &cobra.Command{
		Use:   "list <owner/repo>",
		Short: "List discussion categories",
		Long: `List all discussion categories for a repository.

Categories define the type of discussion and determine whether answers can be marked.`,
		Example: `  ghp discussion category list owner/repo
  ghp discussion category list owner/repo --format json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repo = args[0]
			return runCategoryList(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", formatTable, "Output format: table, json")

	return cmd
}

func runCategoryList(ctx context.Context, opts *CategoryListOptions) error {
	// Parse repository reference
	owner, repo, err := service.ParseRepositoryReference(opts.Repo)
	if err != nil {
		return err
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	discussionService := service.NewDiscussionService(client)

	// List categories
	categories, err := discussionService.ListCategories(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to list categories: %w", err)
	}

	return outputCategories(categories, opts.Format)
}

func outputCategories(categories []service.CategoryInfo, format string) error {
	if len(categories) == 0 {
		fmt.Println("No discussion categories found")
		return nil
	}

	switch format {
	case formatJSON:
		return outputCategoriesJSON(categories)
	case formatTable:
		return outputCategoriesTable(categories)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputCategoriesTable(categories []service.CategoryInfo) error {
	fmt.Printf("%-5s %-20s %-20s %-12s %s\n",
		"EMOJI", "NAME", "SLUG", "ANSWERABLE", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", categorySeparatorWidth))

	for _, c := range categories {
		answerable := "No"
		if c.IsAnswerable {
			answerable = "Yes"
		}

		description := c.Description
		if len(description) > descriptionTruncateLength {
			description = description[:descriptionTruncateLength-3] + "..."
		}

		fmt.Printf("%-5s %-20s %-20s %-12s %s\n",
			c.Emoji, c.Name, c.Slug, answerable, description)
	}

	return nil
}

func outputCategoriesJSON(categories []service.CategoryInfo) error {
	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
