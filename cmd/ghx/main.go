package main

import (
	"fmt"
	"os"

	"github.com/roboco-io/ghx-cli/cmd"
)

// Version information set by build flags
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

func main() {
	// Set version info for the root command
	cmd.SetVersionInfo(Version, Commit, BuildTime)

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
