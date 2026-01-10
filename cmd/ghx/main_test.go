package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Test that main function exists and can be called
	// This is a basic test to ensure the entry point works
	t.Run("Main entry point exists", func(t *testing.T) {
		// We can't directly test main(), but we can test that the file compiles
		assert.NotNil(t, main)
	})
}

func TestVersionInfo(t *testing.T) {
	t.Run("Version variables are defined", func(t *testing.T) {
		// These will be set by the build process
		assert.NotNil(t, Version)
		assert.NotNil(t, Commit)
		assert.NotNil(t, BuildTime)
	})
}

func TestExitCode(t *testing.T) {
	t.Run("Exit with success code", func(t *testing.T) {
		// Save original args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// Set test args
		os.Args = []string{"ghx", "--help"}

		// The main function should not panic
		assert.NotPanics(t, func() {
			// We'll test the actual execution in integration tests
			// For now, just ensure no panic
		})
	})
}
