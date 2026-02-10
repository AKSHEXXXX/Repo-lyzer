package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
)

var (
	prTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00E5FF"))
	prLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7AE7C7"))
	prValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB000"))
	prBarStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
)

// PrintPRAnalytics displays PR analytics in a formatted terminal output
func PrintPRAnalytics(analytics *analyzer.PRAnalytics) {
	fmt.Println(SectionStyle.Render("\n📊 Pull Request Analytics"))

	// Summary section
	fmt.Println(prTitleStyle.Render("\n═══ Summary ═══"))
	fmt.Printf("%s: %s\n", prLabelStyle.Render("Total PRs"), prValueStyle.Render(fmt.Sprintf("%d", analytics.TotalPRs)))
	fmt.Printf("%s: %s\n", prLabelStyle.Render("Merged"), prValueStyle.Render(fmt.Sprintf("%d", analytics.MergedPRs)))
	fmt.Printf("%s: %s\n", prLabelStyle.Render("Closed (not merged)"), prValueStyle.Render(fmt.Sprintf("%d", analytics.ClosedPRs)))
	fmt.Printf("%s: %s\n", prLabelStyle.Render("Open"), prValueStyle.Render(fmt.Sprintf("%d", analytics.OpenPRs)))

	// Timing metrics
	fmt.Println(prTitleStyle.Render("\n═══ Merge Time ═══"))
	if analytics.AverageTimeToMerge > 0 {
		fmt.Printf("%s: %s\n", prLabelStyle.Render("Average Time to Merge"), prValueStyle.Render(formatDuration(analytics.AverageTimeToMerge)))
		fmt.Printf("%s: %s\n", prLabelStyle.Render("Median Time to Merge"), prValueStyle.Render(formatDuration(analytics.MedianTimeToMerge)))
	} else {
		fmt.Printf("%s\n", prLabelStyle.Render("No merged PRs to calculate merge time"))
	}

	// Review participation
	fmt.Println(prTitleStyle.Render("\n═══ Review Participation ═══"))
	fmt.Printf("%s: %s (PRs with 2+ reviewers)\n",
		prLabelStyle.Render("Review Participation"),
		prValueStyle.Render(fmt.Sprintf("%.1f%%", analytics.ReviewParticipation)))

	// Visual bar for review participation
	barLen := int(analytics.ReviewParticipation / 5)
	if barLen > 20 {
		barLen = 20
	}
	bar := prBarStyle.Render(strings.Repeat("█", barLen))
	fmt.Printf("%s\n", bar)

	// PR size distribution
	fmt.Println(prTitleStyle.Render("\n═══ PR Size Distribution ═══"))
	printSizeDistribution(analytics.PRSizeDistribution)

	// Abandoned ratio
	fmt.Println(prTitleStyle.Render("\n═══ Quality Metrics ═══"))
	fmt.Printf("%s: %s\n",
		prLabelStyle.Render("Abandoned PR Ratio"),
		prValueStyle.Render(fmt.Sprintf("%.1f%%", analytics.AbandonedRatio)))

	// First-time contributor metrics
	fmt.Println(prTitleStyle.Render("\n═══ First-Time Contributor Friendliness ═══"))
	if analytics.FirstTimeContributorMetrics.TotalFirstTimePRs > 0 {
		fmt.Printf("%s: %s\n",
			prLabelStyle.Render("First-Time PRs"),
			prValueStyle.Render(fmt.Sprintf("%d", analytics.FirstTimeContributorMetrics.TotalFirstTimePRs)))
		fmt.Printf("%s: %s\n",
			prLabelStyle.Render("Acceptance Rate"),
			prValueStyle.Render(fmt.Sprintf("%.1f%%", analytics.FirstTimeContributorMetrics.AcceptanceRate)))

		if analytics.FirstTimeContributorMetrics.AvgTimeToFirstReview > 0 {
			fmt.Printf("%s: %s\n",
				prLabelStyle.Render("Avg Time to First Review"),
				prValueStyle.Render(formatDuration(analytics.FirstTimeContributorMetrics.AvgTimeToFirstReview)))
		}
	} else {
		fmt.Printf("%s\n", prLabelStyle.Render("No first-time contributors found"))
	}

	fmt.Println()
}

// printSizeDistribution displays PR size distribution with bars
func printSizeDistribution(dist map[string]int) {
	sizes := []string{"small", "medium", "large", "xlarge"}
	labels := map[string]string{
		"small":  "Small (<100 changes)  ",
		"medium": "Medium (100-500)     ",
		"large":  "Large (500-1000)     ",
		"xlarge": "X-Large (>1000)      ",
	}

	total := 0
	for _, count := range dist {
		total += count
	}

	if total == 0 {
		fmt.Println(prLabelStyle.Render("No PR size data available"))
		return
	}

	maxCount := 0
	for _, count := range dist {
		if count > maxCount {
			maxCount = count
		}
	}

	for _, size := range sizes {
		count := dist[size]
		percentage := float64(count) / float64(total) * 100

		barLen := 0
		if maxCount > 0 {
			barLen = int(float64(count) / float64(maxCount) * 20)
		}

		bar := prBarStyle.Render(strings.Repeat("█", barLen))
		fmt.Printf("%s | %s %s (%.0f%%)\n",
			prLabelStyle.Render(labels[size]),
			bar,
			prValueStyle.Render(fmt.Sprintf("%d", count)),
			percentage)
	}
}

// formatDuration formats a duration into a human-readable string
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "N/A"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}

	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%.0fs", d.Seconds())
}

// FormatPRAnalyticsJSON formats PR analytics as JSON
func FormatPRAnalyticsJSON(analytics *analyzer.PRAnalytics) (string, error) {
	jsonData, err := json.MarshalIndent(analytics, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal PR analytics to JSON: %w", err)
	}
	return string(jsonData), nil
}
