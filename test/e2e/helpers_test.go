package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	// Environment variable names
	envTestRepo  = "GHP_TEST_REPO"  // e.g., "owner/repo"
	envTestToken = "GHP_TEST_TOKEN" // GitHub token for testing
	envSkipE2E   = "GHP_SKIP_E2E"   // Set to "1" to skip E2E tests

	// Default test repository (can be overridden via env)
	defaultTestRepo = "roboco-io/gh-project-cli"

	// Timeouts
	commandTimeout = 30 * time.Second
)

// TestConfig holds E2E test configuration
type TestConfig struct {
	Token      string
	Owner      string
	Repo       string
	BinaryPath string
}

// GetTestConfig returns E2E test configuration from environment
func GetTestConfig(t *testing.T) *TestConfig {
	t.Helper()

	// Skip if explicitly disabled
	if os.Getenv(envSkipE2E) == "1" {
		t.Skip("E2E tests disabled via GHP_SKIP_E2E=1")
	}

	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Get token (prefer test token, fallback to standard GitHub token)
	token := os.Getenv(envTestToken)
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token = os.Getenv("GHP_TOKEN")
	}
	if token == "" {
		// Try to get from gh CLI
		out, err := exec.Command("gh", "auth", "token").Output()
		if err == nil {
			token = strings.TrimSpace(string(out))
		}
	}
	if token == "" {
		t.Skip("No GitHub token available for E2E tests (set GHP_TEST_TOKEN, GITHUB_TOKEN, or login via gh auth)")
	}

	// Get test repository
	repo := os.Getenv(envTestRepo)
	if repo == "" {
		repo = defaultTestRepo
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		t.Fatalf("Invalid test repository format: %s (expected owner/repo)", repo)
	}

	// Find binary path
	binaryPath := findBinary(t)

	return &TestConfig{
		Token:      token,
		Owner:      parts[0],
		Repo:       parts[1],
		BinaryPath: binaryPath,
	}
}

// findBinary locates the ghp binary for testing
func findBinary(t *testing.T) string {
	t.Helper()

	// Try common locations
	paths := []string{
		"./bin/ghp",
		"../../bin/ghp",
		os.Getenv("GOPATH") + "/bin/ghp",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Try to build it
	t.Log("Building ghp binary for E2E tests...")
	cmd := exec.Command("go", "build", "-o", "../../bin/ghp", "../../cmd/ghp")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build ghp binary: %v", err)
	}

	return "../../bin/ghp"
}

// CommandResult holds the result of a command execution
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

// RunGHP executes the ghp command with given arguments
func (c *TestConfig) RunGHP(args ...string) *CommandResult {
	cmd := exec.Command(c.BinaryPath, args...)
	cmd.Env = append(os.Environ(), "GITHUB_TOKEN="+c.Token)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
		Err:      err,
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	} else if err != nil {
		result.ExitCode = -1
	}

	return result
}

// RunGHPWithTimeout executes the ghp command with timeout
func (c *TestConfig) RunGHPWithTimeout(timeout time.Duration, args ...string) *CommandResult {
	cmd := exec.Command(c.BinaryPath, args...)
	cmd.Env = append(os.Environ(), "GITHUB_TOKEN="+c.Token)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start command
	if err := cmd.Start(); err != nil {
		return &CommandResult{
			Err:      err,
			ExitCode: -1,
		}
	}

	// Wait with timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		result := &CommandResult{
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			ExitCode: 0,
			Err:      err,
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else if err != nil {
			result.ExitCode = -1
		}
		return result

	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return &CommandResult{
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			ExitCode: -1,
			Err:      fmt.Errorf("command timed out after %v", timeout),
		}
	}
}

// AssertSuccess asserts the command succeeded
func (r *CommandResult) AssertSuccess(t *testing.T) {
	t.Helper()
	if r.Err != nil {
		t.Errorf("Command failed: %v\nStdout: %s\nStderr: %s", r.Err, r.Stdout, r.Stderr)
	}
	if r.ExitCode != 0 {
		t.Errorf("Command exited with code %d\nStdout: %s\nStderr: %s", r.ExitCode, r.Stdout, r.Stderr)
	}
}

// AssertFailure asserts the command failed
func (r *CommandResult) AssertFailure(t *testing.T) {
	t.Helper()
	if r.Err == nil && r.ExitCode == 0 {
		t.Errorf("Expected command to fail but it succeeded\nStdout: %s", r.Stdout)
	}
}

// AssertOutputContains asserts stdout contains the given string
func (r *CommandResult) AssertOutputContains(t *testing.T, substring string) {
	t.Helper()
	if !strings.Contains(r.Stdout, substring) {
		t.Errorf("Expected output to contain %q, got:\n%s", substring, r.Stdout)
	}
}

// AssertOutputNotContains asserts stdout does not contain the given string
func (r *CommandResult) AssertOutputNotContains(t *testing.T, substring string) {
	t.Helper()
	if strings.Contains(r.Stdout, substring) {
		t.Errorf("Expected output to NOT contain %q, got:\n%s", substring, r.Stdout)
	}
}

// AssertErrorContains asserts stderr contains the given string
func (r *CommandResult) AssertErrorContains(t *testing.T, substring string) {
	t.Helper()
	combined := r.Stderr + r.Stdout // Some errors go to stdout
	if !strings.Contains(combined, substring) {
		t.Errorf("Expected error to contain %q, got:\nStderr: %s\nStdout: %s", substring, r.Stderr, r.Stdout)
	}
}

// GenerateTestID generates a unique test ID for resource naming
func GenerateTestID() string {
	return fmt.Sprintf("e2e-test-%d", time.Now().UnixNano())
}
