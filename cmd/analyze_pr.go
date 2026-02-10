package cmd

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/spf13/cobra"
)

var analyzePRCmd = &cobra.Command{
	Use:   "analyze-pr owner/repo",
	Short: "Analyze Pull Request metrics for a GitHub repository",
	Long: `Analyze pull request metrics including:
  • Average time to merge
  • Review participation (% of PRs with 2+ reviewers)
  • PR size distribution
  • Abandoned PR ratio
  • First-time contributor friendliness

Examples:
  repo-lyzer analyze-pr golang/go
  repo-lyzer analyze-pr microsoft/vscode --state closed
  repo-lyzer analyze-pr octocat/Hello-World --limit 50 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, _ := cmd.Flags().GetString("state")
		limit, _ := cmd.Flags().GetInt("limit")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Validate the repository URL format
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Record start time for analysis timing
		startTime := time.Now()

		// Initialize GitHub client
		client := github.NewClient()

		// Inform user about fetching
		if !jsonOutput {
			fmt.Printf("🔍 Fetching pull requests for %s/%s (state: %s)...\n", owner, repo, state)
		}

		// Fetch pull requests
		var prs []github.PullRequest
		if limit > 0 {
			prs, err = client.GetPullRequestsWithLimit(owner, repo, state, limit)
		} else {
			prs, err = client.GetPullRequests(owner, repo, state)
		}
		if err != nil {
			return fmt.Errorf("failed to fetch pull requests: %w", err)
		}

		if len(prs) == 0 {
			if !jsonOutput {
				fmt.Printf("No pull requests found for %s/%s with state '%s'\n", owner, repo, state)
			}
			return nil
		}

		if !jsonOutput {
			fmt.Printf("✓ Found %d pull requests\n", len(prs))
			fmt.Printf("🔍 Fetching reviews for pull requests...\n")
		}

		// Fetch reviews for each PR
		reviews := make(map[int][]github.Review)
		for i, pr := range prs {
			if !jsonOutput && i%10 == 0 {
				fmt.Printf("Progress: %d/%d PRs analyzed\r", i, len(prs))
			}

			prReviews, err := client.GetPullRequestReviews(owner, repo, pr.Number)
			if err != nil {
				// Log error but continue with other PRs
				if !jsonOutput {
					fmt.Printf("⚠️  Warning: failed to fetch reviews for PR #%d: %v\n", pr.Number, err)
				}
				continue
			}
			reviews[pr.Number] = prReviews
		}

		if !jsonOutput {
			fmt.Printf("✓ Completed fetching reviews                    \n\n")
		}

		// Analyze pull requests
		analytics := analyzer.AnalyzePullRequests(prs, reviews)

		// Output results
		if jsonOutput {
			jsonStr, err := output.FormatPRAnalyticsJSON(analytics)
			if err != nil {
				return fmt.Errorf("failed to format JSON: %w", err)
			}
			fmt.Println(jsonStr)
		} else {
			output.PrintPRAnalytics(analytics)

			// Display analysis time
			duration := time.Since(startTime)
			fmt.Printf("⏱️  Analysis completed in %.2f seconds\n", duration.Seconds())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzePRCmd)
	analyzePRCmd.Flags().String("state", "all", "Filter PRs by state: open, closed, or all")
	analyzePRCmd.Flags().Int("limit", 0, "Limit number of PRs to analyze (0 = no limit)")
	analyzePRCmd.Flags().Bool("json", false, "Output results as JSON")
}
