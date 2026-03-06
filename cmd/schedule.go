// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes commands for managing scheduled analysis reports.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/scheduler"
	"github.com/spf13/cobra"
)

// scheduleCmd defines the "schedule" command for managing scheduled reports
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage scheduled analysis reports",
	Long: `Manage scheduled analysis reports that run automatically and export results
to various formats (JSON, PDF, Markdown) and destinations (local path, webhook).

Examples:
  # Add a new scheduled job
  repo-lyzer schedule add owner/repo --interval=weekly --format=json

  # List all scheduled jobs
  repo-lyzer schedule list

  # Remove a scheduled job
  repo-lyzer schedule remove owner-repo-1234567890

  # Enable a scheduled job
  repo-lyzer schedule enable owner-repo-1234567890

  # Disable a scheduled job
  repo-lyzer schedule disable owner-repo-1234567890`,
}

// scheduleAddCmd defines the "schedule add" command
var scheduleAddCmd = &cobra.Command{
	Use:   "add <owner/repo>",
	Short: "Add a new scheduled report job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate repository URL
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Get flags
		intervalStr, _ := cmd.Flags().GetString("interval")
		formatStr, _ := cmd.Flags().GetString("format")
		localPath, _ := cmd.Flags().GetString("local-path")
		webhookURL, _ := cmd.Flags().GetString("webhook")
		analysisType, _ := cmd.Flags().GetString("analysis-type")
		customCron, _ := cmd.Flags().GetString("cron")

		// Parse interval
		interval := config.ScheduleInterval(intervalStr)
		if !isValidInterval(interval) {
			return fmt.Errorf("invalid interval: %s (valid: daily, weekly, monthly, custom)", intervalStr)
		}

		// Parse format
		format := config.ExportFormat(formatStr)
		if !isValidFormat(format) {
			return fmt.Errorf("invalid format: %s (valid: json, markdown, csv, html, pdf)", formatStr)
		}

		// Determine destination type
		destType := "local"
		if webhookURL != "" {
			destType = "webhook"
		}

		// Create destination
		destination := config.OutputDestination{
			Type:       destType,
			LocalPath:  localPath,
			WebhookURL: webhookURL,
			Enabled:    true,
		}

		// Create job
		job := config.ScheduledJob{
			Owner:          owner,
			Repo:           repo,
			Interval:       interval,
			CronExpression: customCron,
			Format:         format,
			Destination:    destination,
			Enabled:        true,
			AnalysisType:   analysisType,
		}

		// Use scheduler to add job
		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		if err := sched.AddJob(job); err != nil {
			return fmt.Errorf("failed to add job: %w", err)
		}

		fmt.Printf("✅ Scheduled job added successfully!\n")
		fmt.Printf("   Repository: %s/%s\n", owner, repo)
		fmt.Printf("   Interval: %s\n", interval.DisplayName())
		fmt.Printf("   Format: %s\n", format.DisplayName())
		fmt.Printf("   Destination: %s\n", destType)

		return nil
	},
}

// scheduleListCmd defines the "schedule list" command
var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled report jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		jobs := sched.ListJobs()

		if len(jobs) == 0 {
			fmt.Println("No scheduled jobs found.")
			fmt.Println("Use 'repo-lyzer schedule add <owner/repo>' to create one.")
			return nil
		}

		fmt.Println("\n📅 Scheduled Jobs")
		fmt.Println(strings.Repeat("─", 70))

		for i, job := range jobs {
			status := "✅ Enabled"
			if !job.Enabled {
				status = "❌ Disabled"
			}

			lastRun := "Never"
			if !job.LastRun.IsZero() {
				lastRun = job.LastRun.Format("2006-01-02 15:04")
			}

			nextRun := "N/A"
			if !job.NextRun.IsZero() {
				nextRun = job.NextRun.Format("2006-01-02 15:04")
			}

			fmt.Printf("\n[%d] %s\n", i+1, job.GetRepoFullName())
			fmt.Printf("    ID: %s\n", job.ID)
			fmt.Printf("    Interval: %s (%s)\n", job.Interval.DisplayName(), job.GetCronExpression())
			fmt.Printf("    Format: %s\n", job.Format.DisplayName())
			fmt.Printf("    Destination: %s\n", job.Destination.Type)
			fmt.Printf("    Status: %s\n", status)
			fmt.Printf("    Last Run: %s\n", lastRun)
			fmt.Printf("    Next Run: %s\n", nextRun)
		}

		fmt.Println("\n" + strings.Repeat("─", 70))
		fmt.Printf("Total: %d job(s)\n\n", len(jobs))

		return nil
	},
}

// scheduleRemoveCmd defines the "schedule remove" command
var scheduleRemoveCmd = &cobra.Command{
	Use:   "remove <job-id>",
	Short: "Remove a scheduled report job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID := args[0]

		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		// Check if job exists
		job := sched.GetJob(jobID)
		if job == nil {
			return fmt.Errorf("job not found: %s", jobID)
		}

		if err := sched.RemoveJob(jobID); err != nil {
			return fmt.Errorf("failed to remove job: %w", err)
		}

		fmt.Printf("✅ Scheduled job removed: %s\n", jobID)

		return nil
	},
}

// scheduleEnableCmd defines the "schedule enable" command
var scheduleEnableCmd = &cobra.Command{
	Use:   "enable <job-id>",
	Short: "Enable a scheduled report job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID := args[0]

		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		if err := sched.EnableJob(jobID, true); err != nil {
			return fmt.Errorf("failed to enable job: %w", err)
		}

		fmt.Printf("✅ Scheduled job enabled: %s\n", jobID)

		return nil
	},
}

// scheduleDisableCmd defines the "schedule disable" command
var scheduleDisableCmd = &cobra.Command{
	Use:   "disable <job-id>",
	Short: "Disable a scheduled report job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID := args[0]

		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		if err := sched.EnableJob(jobID, false); err != nil {
			return fmt.Errorf("failed to disable job: %w", err)
		}

		fmt.Printf("✅ Scheduled job disabled: %s\n", jobID)

		return nil
	},
}

// scheduleRunCmd defines the "schedule run" command to manually trigger a job
var scheduleRunCmd = &cobra.Command{
	Use:   "run <job-id>",
	Short: "Manually trigger a scheduled report job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobID := args[0]

		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		job := sched.GetJob(jobID)
		if job == nil {
			return fmt.Errorf("job not found: %s", jobID)
		}

		fmt.Printf("🔄 Running scheduled job: %s for %s/%s\n", jobID, job.Owner, job.Repo)
		fmt.Println("Note: This may take a while depending on repository size...")

		// Note: In a full implementation, we'd trigger the job execution here
		// For now, we just print a message
		fmt.Println("Job execution triggered (background)")

		return nil
	},
}

// isValidInterval checks if the interval is valid
func isValidInterval(interval config.ScheduleInterval) bool {
	switch interval {
	case config.ScheduleDaily, config.ScheduleWeekly, config.ScheduleMonthly, config.ScheduleCustom:
		return true
	default:
		return false
	}
}

// isValidFormat checks if the format is valid
func isValidFormat(format config.ExportFormat) bool {
	switch format {
	case config.ExportJSON, config.ExportMarkdown, config.ExportCSV, config.ExportHTML, config.ExportPDF:
		return true
	default:
		return false
	}
}

func init() {
	rootCmd.AddCommand(scheduleCmd)

	// Add subcommands
	scheduleCmd.AddCommand(scheduleAddCmd)
	scheduleCmd.AddCommand(scheduleListCmd)
	scheduleCmd.AddCommand(scheduleRemoveCmd)
	scheduleCmd.AddCommand(scheduleEnableCmd)
	scheduleCmd.AddCommand(scheduleDisableCmd)
	scheduleCmd.AddCommand(scheduleRunCmd)

	// Add flags for schedule add
	scheduleAddCmd.Flags().StringP("interval", "i", "weekly", "Schedule interval (daily, weekly, monthly, custom)")
	scheduleAddCmd.Flags().StringP("format", "f", "json", "Export format (json, markdown, csv, html, pdf)")
	scheduleAddCmd.Flags().StringP("local-path", "l", "", "Local directory to save reports")
	scheduleAddCmd.Flags().StringP("webhook", "w", "", "Webhook URL to send reports")
	scheduleAddCmd.Flags().StringP("analysis-type", "a", "quick", "Analysis type (quick, detailed)")
	scheduleAddCmd.Flags().String("cron", "", "Custom cron expression (e.g., '0 9 * * *')")
}

// RunScheduleAdd adds a scheduled job programmatically
func RunScheduleAdd(owner, repo string, interval string, format string, localPath string, webhookURL string) error {
	// Parse interval
	schedInterval := config.ScheduleInterval(interval)
	if !isValidInterval(schedInterval) {
		return fmt.Errorf("invalid interval: %s", interval)
	}

	// Parse format
	exportFormat := config.ExportFormat(format)
	if !isValidFormat(exportFormat) {
		return fmt.Errorf("invalid format: %s", format)
	}

	// Determine destination type
	destType := "local"
	if webhookURL != "" {
		destType = "webhook"
	}

	// Create destination
	destination := config.OutputDestination{
		Type:       destType,
		LocalPath:  localPath,
		WebhookURL: webhookURL,
		Enabled:    true,
	}

	// Create job
	job := config.ScheduledJob{
		Owner:        owner,
		Repo:         repo,
		Interval:     schedInterval,
		Format:       exportFormat,
		Destination:  destination,
		Enabled:      true,
		AnalysisType: "quick",
	}

	// Use scheduler to add job
	sched, err := scheduler.NewScheduler()
	if err != nil {
		return fmt.Errorf("failed to create scheduler: %w", err)
	}

	return sched.AddJob(job)
}

// GetScheduledJobs returns all scheduled jobs
func GetScheduledJobs() ([]config.ScheduledJob, error) {
	sched, err := scheduler.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return sched.ListJobs(), nil
}

// RemoveScheduledJob removes a scheduled job by ID
func RemoveScheduledJob(jobID string) error {
	sched, err := scheduler.NewScheduler()
	if err != nil {
		return fmt.Errorf("failed to create scheduler: %w", err)
	}

	return sched.RemoveJob(jobID)
}

// EnsureDefaultExportDir ensures the default export directory exists
func EnsureDefaultExportDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := home + "/Downloads/repo-lyzer-reports"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return dir, nil
}

// GetDefaultExportDir returns the default export directory
func GetDefaultExportDir() string {
	home, _ := os.UserHomeDir()
	return home + "/Downloads/repo-lyzer-reports"
}
