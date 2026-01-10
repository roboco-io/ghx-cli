package view

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// DeleteOptions holds options for the delete command
type DeleteOptions struct {
	ViewID string
	Format string
	Force  bool
}

// NewDeleteCmd creates the delete command
func NewDeleteCmd() *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <view-id>",
		Short: "Delete a project view",
		Long: `Delete an existing project view.

This operation cannot be undone. By default, you will be prompted for
confirmation unless you use the --force flag.

WARNING: Deleting a view will remove all its configuration including
filters, sorting, and grouping settings.

Examples:
  ghx view delete view-id
  ghx view delete view-id --force
  ghx view delete view-id --format json`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ViewID = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runDelete(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(ctx context.Context, opts *DeleteOptions) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and service
	client := api.NewClient(token)
	viewService := service.NewViewService(client)

	// Get view details for confirmation
	viewInfo, err := viewService.GetView(ctx, opts.ViewID)
	if err != nil {
		return fmt.Errorf("failed to get view details: %w", err)
	}

	// Confirm deletion unless force flag is used
	if !opts.Force {
		fmt.Printf("Are you sure you want to delete view '%s' (%s)?\n", viewInfo.Name, service.FormatViewLayout(viewInfo.Layout))
		fmt.Printf("This action cannot be undone. [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, readErr := reader.ReadString('\n')
		if readErr != nil {
			return fmt.Errorf("failed to read confirmation: %w", readErr)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Deletion canceled.")
			return nil
		}
	}

	// Delete view
	input := service.DeleteViewInput{
		ViewID: opts.ViewID,
	}

	err = viewService.DeleteView(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete view: %w", err)
	}

	// Output confirmation
	return outputDeleteConfirmation(viewInfo, opts.Format)
}

func outputDeleteConfirmation(viewInfo *service.ViewInfo, format string) error {
	switch format {
	case formatJSON:
		return outputDeleteConfirmationJSON(viewInfo)
	case formatTable:
		return outputDeleteConfirmationTable(viewInfo)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputDeleteConfirmationTable(viewInfo *service.ViewInfo) error {
	fmt.Printf("✅ View '%s' deleted successfully\n", viewInfo.Name)
	return nil
}

func outputDeleteConfirmationJSON(viewInfo *service.ViewInfo) error {
	fmt.Printf("{\n")
	fmt.Printf("  \"deleted\": true,\n")
	fmt.Printf("  \"viewId\": \"%s\",\n", viewInfo.ID)
	fmt.Printf("  \"viewName\": \"%s\"\n", viewInfo.Name)
	fmt.Printf("}\n")
	return nil
}
