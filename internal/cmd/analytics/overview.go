package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/roboco-io/ghx-cli/internal/api"
	"github.com/roboco-io/ghx-cli/internal/auth"
	"github.com/roboco-io/ghx-cli/internal/service"
)

const (
	overviewPercentageMultiplier = 100.0
)

// OverviewOptions holds options for the overview command
type OverviewOptions struct {
	ProjectRef string
	Format     string
}

// NewOverviewCmd creates the overview command
func NewOverviewCmd() *cobra.Command {
	opts := &OverviewOptions{}

	cmd := &cobra.Command{
		Use:   "overview <owner/project-number>",
		Short: "Generate project overview analytics",
		Long: `Generate comprehensive overview analytics for a GitHub Project.

The overview report includes:
• Project basic information (title, item count, field count, view count)
• Item distribution by status with counts and percentages
• Item distribution by assignee with workload analysis
• Item distribution by labels and milestones
• Basic velocity metrics and timeline information

This report provides a high-level view of project health and progress,
making it easy to understand the current state and identify potential
bottlenecks or areas that need attention.

Examples:
  ghx analytics overview octocat/123
  ghx analytics overview octocat/123 --format json
  ghx analytics overview --org myorg/456 --format table`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectRef = args[0]
			opts.Format = cmd.Flag("format").Value.String()
			return runOverview(cmd.Context(), opts)
		},
	}

	cmd.Flags().Bool("org", false, "Target organization project")

	return cmd
}

func runOverview(ctx context.Context, opts *OverviewOptions) error {
	// Parse project reference
	parts := strings.Split(opts.ProjectRef, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid project reference format. Use: owner/project-number")
	}

	owner := parts[0]
	projectNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid project number: %s", parts[1])
	}

	// Initialize authentication
	authManager := auth.NewAuthManager()
	token, err := authManager.GetValidatedToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Create client and services
	client := api.NewClient(token)
	projectService := service.NewProjectService(client)
	analyticsService := service.NewAnalyticsService(client)

	// Get project to validate access and get project ID
	project, err := projectService.GetProjectWithOwnerDetection(ctx, owner, projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Get analytics data
	analytics, err := analyticsService.GetProjectAnalytics(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("failed to get project analytics: %w", err)
	}

	// Format and output analytics
	analyticsInfo := service.FormatAnalytics(analytics)
	return outputOverview(analyticsInfo, opts.Format)
}

func outputOverview(analytics *service.AnalyticsInfo, format string) error {
	switch format {
	case FormatJSON:
		return outputOverviewJSON(analytics)
	case FormatTable:
		return outputOverviewTable(analytics)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}

func outputOverviewTable(analytics *service.AnalyticsInfo) error {
	fmt.Printf("📊 Project Overview: %s\n\n", analytics.Title)

	// Basic Statistics
	fmt.Printf("Basic Statistics:\n")
	fmt.Printf("  Project ID: %s\n", analytics.ProjectID)
	fmt.Printf("  Total Items: %d\n", analytics.ItemCount)
	fmt.Printf("  Total Fields: %d\n", analytics.FieldCount)
	fmt.Printf("  Total Views: %d\n", analytics.ViewCount)

	// Status Distribution
	if len(analytics.StatusStats) > 0 {
		fmt.Printf("\n📈 Item Distribution by Status:\n")
		for _, stat := range analytics.StatusStats {
			percentage := float64(stat.Count) / float64(analytics.ItemCount) * overviewPercentageMultiplier
			fmt.Printf("  %-20s %3d items (%.1f%%)\n", stat.Status, stat.Count, percentage)
		}
	}

	// Assignee Distribution
	if len(analytics.AssigneeStats) > 0 {
		fmt.Printf("\n👥 Item Distribution by Assignee:\n")
		for _, stat := range analytics.AssigneeStats {
			assignee := stat.Assignee
			if assignee == "" {
				assignee = "Unassigned"
			}
			percentage := float64(stat.Count) / float64(analytics.ItemCount) * overviewPercentageMultiplier
			fmt.Printf("  %-20s %3d items (%.1f%%)\n", assignee, stat.Count, percentage)
		}
	}

	// Velocity Information
	if analytics.VelocityData != nil {
		fmt.Printf("\n⚡ Velocity Metrics (%s):\n", analytics.VelocityData.Period)
		fmt.Printf("  Completed Items: %d\n", analytics.VelocityData.CompletedItems)
		fmt.Printf("  Added Items: %d\n", analytics.VelocityData.AddedItems)
		fmt.Printf("  Closure Rate: %.1f%%\n", analytics.VelocityData.ClosureRate*overviewPercentageMultiplier)
		fmt.Printf("  Average Lead Time: %s\n", analytics.VelocityData.LeadTime)
		fmt.Printf("  Average Cycle Time: %s\n", analytics.VelocityData.CycleTime)
	}

	// Timeline Information
	if analytics.TimelineData != nil {
		fmt.Printf("\n📅 Timeline Information:\n")
		if analytics.TimelineData.StartDate != nil {
			fmt.Printf("  Start Date: %s\n", *analytics.TimelineData.StartDate)
		}
		if analytics.TimelineData.EndDate != nil {
			fmt.Printf("  End Date: %s\n", *analytics.TimelineData.EndDate)
		}
		if analytics.TimelineData.Duration > 0 {
			fmt.Printf("  Duration: %d days\n", analytics.TimelineData.Duration)
		}
		fmt.Printf("  Milestones: %d\n", analytics.TimelineData.MilestoneCount)
		fmt.Printf("  Activities: %d\n", analytics.TimelineData.ActivityCount)
	}

	return nil
}

func outputOverviewJSON(analytics *service.AnalyticsInfo) error {
	return outputJSONObject(map[string]interface{}{
		"projectId":            analytics.ProjectID,
		"title":                analytics.Title,
		"itemCount":            analytics.ItemCount,
		"fieldCount":           analytics.FieldCount,
		"viewCount":            analytics.ViewCount,
		"statusDistribution":   formatStatusStats(analytics.StatusStats, analytics.ItemCount),
		"assigneeDistribution": formatAssigneeStats(analytics.AssigneeStats, analytics.ItemCount),
		"velocity":             formatVelocityData(analytics.VelocityData),
		"timeline":             formatTimelineData(analytics.TimelineData),
	})
}

func formatStatusStats(stats []service.StatusStat, total int) []map[string]interface{} {
	result := make([]map[string]interface{}, len(stats))
	for i, stat := range stats {
		percentage := float64(stat.Count) / float64(total) * overviewPercentageMultiplier
		result[i] = map[string]interface{}{
			"status":     stat.Status,
			"count":      stat.Count,
			"percentage": percentage,
		}
	}
	return result
}

func formatAssigneeStats(stats []service.AssigneeStat, total int) []map[string]interface{} {
	result := make([]map[string]interface{}, len(stats))
	for i, stat := range stats {
		assignee := stat.Assignee
		if assignee == "" {
			assignee = "Unassigned"
		}
		percentage := float64(stat.Count) / float64(total) * overviewPercentageMultiplier
		result[i] = map[string]interface{}{
			"assignee":   assignee,
			"count":      stat.Count,
			"percentage": percentage,
		}
	}
	return result
}

func formatVelocityData(velocity *service.VelocityInfo) interface{} {
	if velocity == nil {
		return nil
	}
	return map[string]interface{}{
		"period":         velocity.Period,
		"completedItems": velocity.CompletedItems,
		"addedItems":     velocity.AddedItems,
		"closureRate":    velocity.ClosureRate,
		"leadTime":       velocity.LeadTime,
		"cycleTime":      velocity.CycleTime,
	}
}

func formatTimelineData(timeline *service.TimelineInfo) interface{} {
	if timeline == nil {
		return nil
	}
	result := map[string]interface{}{
		"milestoneCount": timeline.MilestoneCount,
		"activityCount":  timeline.ActivityCount,
	}
	if timeline.StartDate != nil {
		result["startDate"] = *timeline.StartDate
	}
	if timeline.EndDate != nil {
		result["endDate"] = *timeline.EndDate
	}
	if timeline.Duration > 0 {
		result["duration"] = timeline.Duration
	}
	return result
}

func outputJSONObject(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}
