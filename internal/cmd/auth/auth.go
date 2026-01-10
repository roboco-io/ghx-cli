package auth

import (
	"github.com/spf13/cobra"
)

// NewAuthCmd creates the auth command group
func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Manage authentication",
		Long: `Manage GitHub authentication for ghp-cli.

This command group provides authentication management capabilities including:

• Check authentication status and token validity
• View available and required scopes
• Get recommendations for authentication setup

Authentication is handled through GitHub CLI integration with fallback to
environment variables (GITHUB_TOKEN or GH_TOKEN).

For initial setup, authenticate with GitHub CLI:
  gh auth login

For more information about GitHub CLI authentication:
https://docs.github.com/en/github-cli/github-cli/about-github-cli`,
		Example: `  ghx auth status                     # Check authentication status
  ghx auth status --format json       # Show status as JSON`,
	}

	// Add subcommands
	cmd.AddCommand(NewStatusCmd())

	return cmd
}
