package field

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// NewDeleteOptionCmd creates the delete-option command
func NewDeleteOptionCmd() *cobra.Command {
	config := DeleteCommandConfig{
		Use:   "delete-option <option-id>",
		Short: "Delete single select field option",
		Long: `Delete an option from a single select field.

⚠️  WARNING: This action is irreversible. Items that currently have this
option selected will lose their field value. Use with caution.

By default, this command will prompt for confirmation. Use --force to skip
the confirmation prompt.

Examples:
  ghx field delete-option option-id
  ghx field delete-option option-id --force`,
		ItemType: "field option",
		ServiceAction: func(ctx context.Context, client *api.Client, optionID string) error {
			fieldService := service.NewFieldService(client)
			input := service.DeleteFieldOptionInput{
				OptionID: optionID,
			}
			return fieldService.DeleteFieldOption(ctx, input)
		},
	}

	return createDeleteCmd(config)
}
