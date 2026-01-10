package field

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/service"
)

// NewDeleteCmd creates the delete command
func NewDeleteCmd() *cobra.Command {
	config := DeleteCommandConfig{
		Use:   "delete <field-id>",
		Short: "Delete a project field",
		Long: `Delete a project field and all its data.

⚠️  WARNING: This action is irreversible. All field data for project items
will be permanently lost. Use with caution.

By default, this command will prompt for confirmation. Use --force to skip
the confirmation prompt.

Examples:
  ghx field delete field-id
  ghx field delete field-id --force`,
		ItemType: "field",
		ServiceAction: func(ctx context.Context, client *api.Client, fieldID string) error {
			fieldService := service.NewFieldService(client)
			input := service.DeleteFieldInput{
				FieldID: fieldID,
			}
			return fieldService.DeleteField(ctx, input)
		},
	}

	return createDeleteCmd(config)
}
