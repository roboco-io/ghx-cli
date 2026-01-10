package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/auth"
)

// StatusOptions holds options for the status command
type StatusOptions struct {
	Format string
}

// NewStatusCmd creates the status command
func NewStatusCmd() *cobra.Command {
	opts := &StatusOptions{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long: `Display detailed information about the current authentication status.

This command checks:
• GitHub CLI installation
• Token availability (from gh CLI or environment)
• Token validity with GitHub API
• Required scopes for GitHub Projects

Examples:
  ghx auth status                 # Show status in table format
  ghx auth status --format json  # Show status as JSON`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runStatus(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "table", "Output format: table, json")

	return cmd
}

func runStatus(opts *StatusOptions) error {
	authManager := auth.NewAuthManager()
	status := authManager.GetAuthenticationStatus()

	switch opts.Format {
	case "json":
		return outputStatusJSON(status)
	case "table":
		return outputStatusTable(status)
	default:
		return fmt.Errorf("unknown format: %s", opts.Format)
	}
}

func outputStatusTable(status auth.Status) error {
	fmt.Printf("GitHub CLI Authentication Status\n")
	fmt.Printf("================================\n\n")

	// Overall status
	if status.IsReady() {
		fmt.Printf("✅ Status: Ready\n")
	} else {
		fmt.Printf("❌ Status: Not Ready\n")
	}

	// Details
	fmt.Printf("\nDetails:\n")
	fmt.Printf("--------\n")

	if status.GHCLIInstalled {
		fmt.Printf("✅ GitHub CLI: Installed\n")
	} else {
		fmt.Printf("❌ GitHub CLI: Not installed\n")
	}

	if status.HasEnvToken {
		fmt.Printf("✅ Environment Token: Available\n")
	} else {
		fmt.Printf("ℹ️  Environment Token: Not set\n")
	}

	if status.TokenAvailable {
		fmt.Printf("✅ Token: Available\n")
	} else {
		fmt.Printf("❌ Token: Not available\n")
	}

	if status.TokenValid {
		fmt.Printf("✅ Token Validity: Valid\n")
	} else if status.TokenAvailable {
		fmt.Printf("❌ Token Validity: Invalid\n")
	} else {
		fmt.Printf("➖ Token Validity: N/A\n")
	}

	if status.HasRequiredScopes {
		fmt.Printf("✅ Required Scopes: Available\n")
	} else if status.TokenValid {
		fmt.Printf("❌ Required Scopes: Missing\n")
	} else {
		fmt.Printf("➖ Required Scopes: N/A\n")
	}

	// Scopes information
	if len(status.Scopes) > 0 {
		fmt.Printf("\nAvailable Scopes: %v\n", status.Scopes)
	}
	if len(status.RequiredScopes) > 0 {
		fmt.Printf("Required Scopes: %v\n", status.RequiredScopes)
	}

	// Error information
	if status.Error != "" {
		fmt.Printf("\nError: %s\n", status.Error)
	}

	// Recommendation
	fmt.Printf("\nRecommendation:\n")
	fmt.Printf("---------------\n")
	fmt.Printf("%s\n", status.GetRecommendation())

	return nil
}

func outputStatusJSON(status auth.Status) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(status)
}
