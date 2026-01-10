package discussion

import (
	"github.com/spf13/cobra"
)

// NewDiscussionCmd creates the discussion command group
func NewDiscussionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discussion <command>",
		Short: "Manage GitHub Discussions",
		Long: `Work with GitHub Discussions in repositories.

GitHub Discussions provide a collaborative communication forum within your repository.
This command group provides comprehensive discussion management capabilities including:

- List, view, create, edit, and delete discussions
- Close, reopen, lock, and unlock discussions
- Add comments and mark answers
- Manage discussion categories

For more information about GitHub Discussions, visit:
https://docs.github.com/en/discussions`,
		Example: `  ghp discussion list owner/repo              # List discussions
  ghp discussion view owner/repo 123          # View discussion #123
  ghp discussion create owner/repo            # Create a discussion
  ghp discussion close owner/repo 123         # Close discussion #123
  ghp discussion comment owner/repo 123       # Add a comment
  ghp discussion category list owner/repo     # List categories`,
		Aliases: []string{"disc", "discussions"},
	}

	// Add subcommands
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewViewCmd())
	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewEditCmd())
	cmd.AddCommand(NewDeleteCmd())
	cmd.AddCommand(NewCloseCmd())
	cmd.AddCommand(NewReopenCmd())
	cmd.AddCommand(NewLockCmd())
	cmd.AddCommand(NewUnlockCmd())
	cmd.AddCommand(NewCommentCmd())
	cmd.AddCommand(NewAnswerCmd())
	cmd.AddCommand(NewCategoryCmd())

	return cmd
}
