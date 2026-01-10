package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

// TestDiscussionList tests the 'ghp discussion list' command
func TestDiscussionList(t *testing.T) {
	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	t.Run("list discussions", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg)
		result.AssertSuccess(t)
		// Output should contain table header or "No discussions found"
		if !strings.Contains(result.Stdout, "NUMBER") && !strings.Contains(result.Stdout, "No discussions found") {
			t.Errorf("Expected discussion list output, got: %s", result.Stdout)
		}
	})

	t.Run("list discussions with limit", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--limit", "5")
		result.AssertSuccess(t)
	})

	t.Run("list discussions with json format", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--format", "json")
		result.AssertSuccess(t)
		// JSON output should start with [ or contain "No discussions"
		if !strings.HasPrefix(strings.TrimSpace(result.Stdout), "[") && !strings.Contains(result.Stdout, "No discussions found") {
			t.Errorf("Expected JSON array output, got: %s", result.Stdout)
		}
	})

	t.Run("list discussions with state filter", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--state", "open")
		result.AssertSuccess(t)
	})

	t.Run("list discussions with invalid repo format", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", "invalid-repo-format")
		result.AssertFailure(t)
		result.AssertErrorContains(t, "invalid repository format")
	})
}

// TestDiscussionCategoryList tests the 'ghp discussion category list' command
func TestDiscussionCategoryList(t *testing.T) {
	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	t.Run("list categories", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "category", "list", repoArg)
		result.AssertSuccess(t)
		// Output should contain table header or "No discussion categories"
		if !strings.Contains(result.Stdout, "NAME") && !strings.Contains(result.Stdout, "No discussion categories") {
			t.Errorf("Expected category list output, got: %s", result.Stdout)
		}
	})

	t.Run("list categories with json format", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "category", "list", repoArg, "--format", "json")
		result.AssertSuccess(t)
	})
}

// TestDiscussionView tests the 'ghp discussion view' command
func TestDiscussionView(t *testing.T) {
	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	// First, get a valid discussion number from list
	listResult := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--format", "json")
	if listResult.Err != nil || strings.Contains(listResult.Stdout, "No discussions found") {
		t.Skip("No discussions available in test repository")
	}

	// Extract first discussion number from JSON
	numberRegex := regexp.MustCompile(`"number":\s*(\d+)`)
	matches := numberRegex.FindStringSubmatch(listResult.Stdout)
	if len(matches) < 2 {
		t.Skip("Could not find discussion number in list output")
	}
	discussionNumber := matches[1]

	t.Run("view discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "view", repoArg, discussionNumber)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "#"+discussionNumber)
	})

	t.Run("view discussion with json format", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "view", repoArg, discussionNumber, "--format", "json")
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "\"number\"")
	})

	t.Run("view non-existent discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "view", repoArg, "999999")
		result.AssertFailure(t)
	})

	t.Run("view discussion with invalid number", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "view", repoArg, "not-a-number")
		result.AssertFailure(t)
		result.AssertErrorContains(t, "invalid discussion number")
	})
}

// TestDiscussionCRUDWorkflow tests the full create-read-update-delete workflow
// This test requires write access to the repository
func TestDiscussionCRUDWorkflow(t *testing.T) {
	// Skip unless explicitly enabled with write access
	if os.Getenv("GHP_E2E_WRITE_TESTS") != "1" {
		t.Skip("Write tests disabled (set GHP_E2E_WRITE_TESTS=1 to enable)")
	}

	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	// First, find an available category
	catResult := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "category", "list", repoArg, "--format", "json")
	catResult.AssertSuccess(t)

	// Extract first category slug
	slugRegex := regexp.MustCompile(`"slug":\s*"([^"]+)"`)
	matches := slugRegex.FindStringSubmatch(catResult.Stdout)
	if len(matches) < 2 {
		t.Skip("No categories available in test repository")
	}
	categorySlug := matches[1]

	// Generate unique test ID
	testID := GenerateTestID()
	testTitle := fmt.Sprintf("E2E Test Discussion %s", testID)
	testBody := fmt.Sprintf("This is an automated E2E test discussion created at %s.\nTest ID: %s", time.Now().Format(time.RFC3339), testID)

	var discussionNumber string

	// Step 1: Create discussion
	t.Run("create discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout,
			"discussion", "create", repoArg,
			"--category", categorySlug,
			"--title", testTitle,
			"--body", testBody,
		)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "Created discussion")

		// Extract discussion number from output
		numberRegex := regexp.MustCompile(`#(\d+)`)
		matches := numberRegex.FindStringSubmatch(result.Stdout)
		if len(matches) < 2 {
			t.Fatalf("Could not extract discussion number from create output: %s", result.Stdout)
		}
		discussionNumber = matches[1]
		t.Logf("Created discussion #%s", discussionNumber)
	})

	if discussionNumber == "" {
		t.Fatal("Discussion was not created, cannot continue workflow tests")
	}

	// Cleanup function to delete the discussion
	defer func() {
		t.Run("cleanup: delete discussion", func(t *testing.T) {
			result := cfg.RunGHPWithTimeout(commandTimeout,
				"discussion", "delete", repoArg, discussionNumber, "--force",
			)
			if result.Err != nil {
				t.Logf("Warning: Failed to cleanup discussion #%s: %v", discussionNumber, result.Err)
			}
		})
	}()

	// Step 2: View the created discussion
	t.Run("view created discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "view", repoArg, discussionNumber)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, testTitle)
	})

	// Step 3: Edit discussion
	t.Run("edit discussion title", func(t *testing.T) {
		newTitle := testTitle + " (Edited)"
		result := cfg.RunGHPWithTimeout(commandTimeout,
			"discussion", "edit", repoArg, discussionNumber,
			"--title", newTitle,
		)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "Updated discussion")
	})

	// Step 4: Add a comment
	t.Run("add comment to discussion", func(t *testing.T) {
		commentBody := fmt.Sprintf("E2E test comment at %s", time.Now().Format(time.RFC3339))
		result := cfg.RunGHPWithTimeout(commandTimeout,
			"discussion", "comment", repoArg, discussionNumber,
			"--body", commentBody,
		)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "Added comment")
	})

	// Step 5: Close discussion
	t.Run("close discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout,
			"discussion", "close", repoArg, discussionNumber,
			"--reason", "resolved",
		)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "Closed discussion")
	})

	// Step 6: Reopen discussion
	t.Run("reopen discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout,
			"discussion", "reopen", repoArg, discussionNumber,
		)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "Reopened discussion")
	})

	// Step 7: Lock discussion
	t.Run("lock discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout,
			"discussion", "lock", repoArg, discussionNumber,
		)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "Locked discussion")
	})

	// Step 8: Unlock discussion
	t.Run("unlock discussion", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout,
			"discussion", "unlock", repoArg, discussionNumber,
		)
		result.AssertSuccess(t)
		result.AssertOutputContains(t, "Unlocked discussion")
	})

	// Step 9: Delete is handled by defer cleanup
}

// TestDiscussionCommandHelp tests that all discussion commands have proper help
func TestDiscussionCommandHelp(t *testing.T) {
	cfg := GetTestConfig(t)

	commands := []string{
		"discussion",
		"discussion list",
		"discussion view",
		"discussion create",
		"discussion edit",
		"discussion delete",
		"discussion close",
		"discussion reopen",
		"discussion lock",
		"discussion unlock",
		"discussion comment",
		"discussion answer",
		"discussion category",
		"discussion category list",
	}

	for _, cmd := range commands {
		t.Run(fmt.Sprintf("help for %q", cmd), func(t *testing.T) {
			args := append(strings.Split(cmd, " "), "--help")
			result := cfg.RunGHPWithTimeout(commandTimeout, args...)
			result.AssertSuccess(t)
			// Help output should contain "Usage:" or "Available Commands:"
			if !strings.Contains(result.Stdout, "Usage:") && !strings.Contains(result.Stdout, "Available Commands:") {
				t.Errorf("Expected help output for %q, got: %s", cmd, result.Stdout)
			}
		})
	}
}

// TestDiscussionAliases tests that command aliases work
func TestDiscussionAliases(t *testing.T) {
	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	t.Run("disc alias works", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "disc", "list", repoArg)
		result.AssertSuccess(t)
	})

	t.Run("discussions alias works", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussions", "list", repoArg)
		result.AssertSuccess(t)
	})
}

// TestDiscussionFilterCombinations tests various filter combinations
func TestDiscussionFilterCombinations(t *testing.T) {
	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	testCases := []struct {
		name string
		args []string
	}{
		{"state open", []string{"--state", "open"}},
		{"state closed", []string{"--state", "closed"}},
		{"state all", []string{"--state", "all"}},
		{"limit 1", []string{"--limit", "1"}},
		{"limit 100", []string{"--limit", "100"}},
		{"answered true", []string{"--answered"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"discussion", "list", repoArg}, tc.args...)
			result := cfg.RunGHPWithTimeout(commandTimeout, args...)
			result.AssertSuccess(t)
		})
	}
}

// TestDiscussionErrorHandling tests error handling for invalid inputs
func TestDiscussionErrorHandling(t *testing.T) {
	cfg := GetTestConfig(t)

	t.Run("missing repo argument", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list")
		result.AssertFailure(t)
	})

	t.Run("missing number argument for view", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "view", "owner/repo")
		result.AssertFailure(t)
	})

	t.Run("invalid state value", func(t *testing.T) {
		// This might succeed with filtering nothing, or fail - either is acceptable
		// The command should at least not panic
		_ = cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", "owner/repo", "--state", "invalid")
	})

	t.Run("invalid format value", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", "owner/repo", "--format", "xml")
		result.AssertFailure(t)
	})
}

// TestDiscussionOutputFormats tests different output formats
func TestDiscussionOutputFormats(t *testing.T) {
	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	// First check if there are discussions
	listResult := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--limit", "1")
	if strings.Contains(listResult.Stdout, "No discussions found") {
		t.Skip("No discussions available to test output formats")
	}

	t.Run("table format", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--format", "table")
		result.AssertSuccess(t)
		// Table format should have header row with column names
		result.AssertOutputContains(t, "NUMBER")
	})

	t.Run("json format", func(t *testing.T) {
		result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--format", "json")
		result.AssertSuccess(t)
		// JSON format should be valid JSON array
		if !strings.HasPrefix(strings.TrimSpace(result.Stdout), "[") {
			t.Errorf("Expected JSON array, got: %s", result.Stdout)
		}
	})
}

// TestDiscussionPagination tests pagination functionality
func TestDiscussionPagination(t *testing.T) {
	cfg := GetTestConfig(t)
	repoArg := fmt.Sprintf("%s/%s", cfg.Owner, cfg.Repo)

	// Get total count with high limit
	result := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--limit", "100", "--format", "json")
	result.AssertSuccess(t)

	if strings.Contains(result.Stdout, "No discussions found") {
		t.Skip("No discussions available to test pagination")
	}

	// Count discussions in full list
	fullCount := strings.Count(result.Stdout, "\"number\":")
	if fullCount < 2 {
		t.Skip("Not enough discussions to test pagination")
	}

	// Get with limit 1
	limitedResult := cfg.RunGHPWithTimeout(commandTimeout, "discussion", "list", repoArg, "--limit", "1", "--format", "json")
	limitedResult.AssertSuccess(t)

	limitedCount := strings.Count(limitedResult.Stdout, "\"number\":")
	if limitedCount != 1 {
		t.Errorf("Expected 1 discussion with --limit 1, got %d", limitedCount)
	}
}

// BenchmarkDiscussionList benchmarks the list command performance
func BenchmarkDiscussionList(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	token := os.Getenv("GHP_TEST_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		b.Skip("No GitHub token available")
	}

	repo := os.Getenv(envTestRepo)
	if repo == "" {
		repo = defaultTestRepo
	}

	binaryPath := "./bin/ghp"
	if _, err := os.Stat(binaryPath); err != nil {
		binaryPath = "../../bin/ghp"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command(binaryPath, "discussion", "list", repo, "--limit", "10")
		cmd.Env = append(os.Environ(), "GITHUB_TOKEN="+token)
		_ = cmd.Run()
	}
}
