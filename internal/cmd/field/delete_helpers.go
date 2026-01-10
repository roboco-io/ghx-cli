package field

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
)

// CommonDeleteOptions represents common options for delete operations
type CommonDeleteOptions struct {
	ID    string
	Force bool
}

// DeleteCommandConfig holds configuration for creating delete commands
type DeleteCommandConfig struct {
	ServiceAction func(context.Context, *api.Client, string) error
	Use           string
	Short         string
	Long          string
	ItemType      string
}

// createDeleteCmd creates a standardized delete command
func createDeleteCmd(config DeleteCommandConfig) *cobra.Command {
	opts := &CommonDeleteOptions{}

	cmd := &cobra.Command{
		Use:   config.Use,
		Short: config.Short,
		Long:  config.Long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]
			return executeDelete(cmd.Context(), opts, config.ItemType, config.ServiceAction)
		},
	}

	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	return cmd
}

// executeDelete handles the common delete workflow
func executeDelete(
	ctx context.Context,
	opts *CommonDeleteOptions,
	itemType string,
	serviceAction func(context.Context, *api.Client, string) error,
) error {
	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client
	client := api.NewClient(token)

	// Show confirmation unless --force is used
	if !opts.Force {
		fmt.Printf("⚠️  You are about to delete %s: %s\n", itemType, opts.ID)
		fmt.Printf("\nThis action cannot be undone. All %s data will be permanently lost.\n", itemType)
		fmt.Printf("Type 'DELETE' to confirm: ")

		var confirmation string
		_, scanErr := fmt.Scanln(&confirmation)
		if scanErr != nil {
			fmt.Println("❌ Failed to read confirmation.")
			return scanErr
		}

		if confirmation != "DELETE" {
			fmt.Println("❌ Deletion canceled.")
			return nil
		}
	}

	// Execute the service action
	err = serviceAction(ctx, client, opts.ID)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", itemType, err)
	}

	fmt.Printf("✅ %s %s deleted successfully.\n", itemType, opts.ID)
	return nil
}
