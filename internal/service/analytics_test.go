package service

import (
	"context"
	"testing"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/api/graphql"
)

var testCtx = context.Background()

func TestNewAnalyticsService(t *testing.T) {
	client := api.NewClient("test-token")
	service := NewAnalyticsService(client)

	if service == nil {
		t.Fatal("Expected analytics service to be created")
	}

	if service.client != client {
		t.Error("Expected client to be set correctly")
	}
}

func TestAnalyticsService_GetProjectAnalytics(t *testing.T) {
	client := api.NewClient("invalid-token")
	service := NewAnalyticsService(client)

	// Test with invalid token (should fail)
	_, err := service.GetProjectAnalytics(testCtx, "test-project-id")
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestAnalyticsService_ExportProject(t *testing.T) {
	client := api.NewClient("invalid-token")
	service := NewAnalyticsService(client)

	input := ExportProjectInput{
		ProjectID:        "test-project-id",
		Format:           graphql.ProjectV2ExportFormatJSON,
		IncludeItems:     true,
		IncludeFields:    true,
		IncludeViews:     true,
		IncludeWorkflows: true,
	}

	// Test with invalid token (should fail)
	_, err := service.ExportProject(testCtx, input)
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestAnalyticsService_ImportProject(t *testing.T) {
	client := api.NewClient("invalid-token")
	service := NewAnalyticsService(client)

	input := ImportProjectInput{
		ProjectID:     "test-project-id",
		Format:        graphql.ProjectV2ExportFormatJSON,
		Data:          `{"items": []}`,
		MergeStrategy: "merge",
	}

	// Test with invalid token (should fail)
	_, err := service.ImportProject(testCtx, input)
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestAnalyticsService_BulkUpdateItems(t *testing.T) {
	client := api.NewClient("invalid-token")
	service := NewAnalyticsService(client)

	input := BulkUpdateItemsInput{
		ProjectID: "test-project-id",
		ItemIDs:   []string{"item1", "item2"},
		Updates:   map[string]interface{}{"status": "Done"},
	}

	// Test with invalid token (should fail)
	_, err := service.BulkUpdateItems(testCtx, input)
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestAnalyticsService_BulkDeleteItems(t *testing.T) {
	client := api.NewClient("invalid-token")
	service := NewAnalyticsService(client)

	input := BulkDeleteItemsInput{
		ProjectID: "test-project-id",
		ItemIDs:   []string{"item1", "item2"},
	}

	// Test with invalid token (should fail)
	_, err := service.BulkDeleteItems(testCtx, input)
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestAnalyticsService_BulkArchiveItems(t *testing.T) {
	client := api.NewClient("invalid-token")
	service := NewAnalyticsService(client)

	input := BulkArchiveItemsInput{
		ProjectID: "test-project-id",
		ItemIDs:   []string{"item1", "item2"},
	}

	// Test with invalid token (should fail)
	_, err := service.BulkArchiveItems(testCtx, input)
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestAnalyticsService_GetBulkOperation(t *testing.T) {
	client := api.NewClient("invalid-token")
	service := NewAnalyticsService(client)

	// Test with invalid token (should fail)
	_, err := service.GetBulkOperation(testCtx, "test-operation-id")
	if err == nil {
		t.Error("Expected error with invalid token")
	}
}

func TestValidateExportFormat(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		expected    graphql.ProjectV2ExportFormat
		expectError bool
	}{
		{
			name:     "Valid JSON format",
			format:   "JSON",
			expected: graphql.ProjectV2ExportFormatJSON,
		},
		{
			name:     "Valid CSV format",
			format:   "CSV",
			expected: graphql.ProjectV2ExportFormatCSV,
		},
		{
			name:     "Valid XML format",
			format:   "XML",
			expected: graphql.ProjectV2ExportFormatXML,
		},
		{
			name:     "Valid JSON lowercase",
			format:   "json",
			expected: graphql.ProjectV2ExportFormatJSON,
		},
		{
			name:     "Valid CSV with spaces",
			format:   " csv ",
			expected: graphql.ProjectV2ExportFormatCSV,
		},
		{
			name:        "Invalid format",
			format:      "YAML",
			expectError: true,
		},
		{
			name:        "Empty format",
			format:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateExportFormat(tt.format)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestValidateBulkOperationType(t *testing.T) {
	tests := []struct {
		name        string
		opType      string
		expected    graphql.BulkOperationType
		expectError bool
	}{
		{
			name:     "Valid UPDATE type",
			opType:   "UPDATE",
			expected: graphql.BulkOperationTypeUpdate,
		},
		{
			name:     "Valid DELETE type",
			opType:   "DELETE",
			expected: graphql.BulkOperationTypeDelete,
		},
		{
			name:     "Valid IMPORT type",
			opType:   "IMPORT",
			expected: graphql.BulkOperationTypeImport,
		},
		{
			name:     "Valid EXPORT type",
			opType:   "EXPORT",
			expected: graphql.BulkOperationTypeExport,
		},
		{
			name:     "Valid ARCHIVE type",
			opType:   "ARCHIVE",
			expected: graphql.BulkOperationTypeArchive,
		},
		{
			name:     "Valid MOVE type",
			opType:   "MOVE",
			expected: graphql.BulkOperationTypeMove,
		},
		{
			name:     "Valid update lowercase",
			opType:   "update",
			expected: graphql.BulkOperationTypeUpdate,
		},
		{
			name:     "Valid delete with spaces",
			opType:   " delete ",
			expected: graphql.BulkOperationTypeDelete,
		},
		{
			name:        "Invalid operation type",
			opType:      "INVALID",
			expectError: true,
		},
		{
			name:        "Empty operation type",
			opType:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateBulkOperationType(tt.opType)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestValidateMergeStrategy(t *testing.T) {
	tests := []struct {
		name        string
		strategy    string
		expectError bool
	}{
		{
			name:     "Valid merge strategy",
			strategy: "merge",
		},
		{
			name:     "Valid replace strategy",
			strategy: "replace",
		},
		{
			name:     "Valid append strategy",
			strategy: "append",
		},
		{
			name:     "Valid skip_conflicts strategy",
			strategy: "skip_conflicts",
		},
		{
			name:     "Valid MERGE uppercase",
			strategy: "MERGE",
		},
		{
			name:     "Valid replace with spaces",
			strategy: " replace ",
		},
		{
			name:        "Invalid strategy",
			strategy:    "invalid",
			expectError: true,
		},
		{
			name:        "Empty strategy",
			strategy:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMergeStrategy(tt.strategy)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestFormatExportFormat(t *testing.T) {
	tests := []struct {
		format   graphql.ProjectV2ExportFormat
		expected string
	}{
		{graphql.ProjectV2ExportFormatJSON, "JSON"},
		{graphql.ProjectV2ExportFormatCSV, "CSV"},
		{graphql.ProjectV2ExportFormatXML, "XML"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			result := FormatExportFormat(tt.format)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatBulkOperationType(t *testing.T) {
	tests := []struct {
		opType   graphql.BulkOperationType
		expected string
	}{
		{graphql.BulkOperationTypeUpdate, "Update"},
		{graphql.BulkOperationTypeDelete, "Delete"},
		{graphql.BulkOperationTypeImport, "Import"},
		{graphql.BulkOperationTypeExport, "Export"},
		{graphql.BulkOperationTypeArchive, "Archive"},
		{graphql.BulkOperationTypeMove, "Move"},
	}

	for _, tt := range tests {
		t.Run(string(tt.opType), func(t *testing.T) {
			result := FormatBulkOperationType(tt.opType)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatBulkOperationStatus(t *testing.T) {
	tests := []struct {
		status   graphql.BulkOperationStatus
		expected string
	}{
		{graphql.BulkOperationStatusPending, "Pending"},
		{graphql.BulkOperationStatusRunning, "Running"},
		{graphql.BulkOperationStatusCompleted, "Completed"},
		{graphql.BulkOperationStatusFailed, "Failed"},
		{graphql.BulkOperationStatusCancelled, "Canceled"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := FormatBulkOperationStatus(tt.status)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
