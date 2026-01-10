package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	t.Run("Root command exists", func(t *testing.T) {
		cmd := NewRootCmd()
		assert.NotNil(t, cmd)
		assert.Equal(t, "ghp", cmd.Use)
	})

	t.Run("Root command has description", func(t *testing.T) {
		cmd := NewRootCmd()
		assert.NotEmpty(t, cmd.Short)
		assert.NotEmpty(t, cmd.Long)
	})

	t.Run("Version flag works", func(t *testing.T) {
		SetVersionInfo("1.0.0", "abc123", "2024-01-01")
		cmd := NewRootCmd()

		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"--version"})

		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "1.0.0")
		assert.Contains(t, output, "abc123")
		assert.Contains(t, output, "2024-01-01")
	})

	t.Run("Help flag works", func(t *testing.T) {
		cmd := NewRootCmd()

		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"--help"})

		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "command-line interface for GitHub features")
		assert.Contains(t, output, "ghp")
	})

	t.Run("Invalid command shows error", func(t *testing.T) {
		cmd := NewRootCmd()

		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"invalid-command"})

		err := cmd.Execute()
		// Should return error for invalid command
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown command")
	})
}

func TestExecute(t *testing.T) {
	t.Run("Execute function works", func(t *testing.T) {
		// Save original args
		oldArgs := rootCmd.Commands()
		defer func() {
			// Reset root command
			rootCmd = NewRootCmd()
			for _, cmd := range oldArgs {
				rootCmd.AddCommand(cmd)
			}
		}()

		// Test with help flag to avoid side effects
		rootCmd = NewRootCmd()
		rootCmd.SetArgs([]string{"--help"})

		err := Execute()
		assert.NoError(t, err)
	})
}
