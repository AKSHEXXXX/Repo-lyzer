// Package output provides formatting and output functions for Repo-lyzer.
// This file implements collaboration insights visualization.
package output

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
)

// PrintCollaborationInsights prints comprehensive collaboration analysis
func PrintCollaborationInsights(metrics *analyzer.CollaborationMetrics) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 TEAM COLLABORATION INSIGHTS")
	fmt.Println(strings.Repeat("=", 60))

	// Overall Score
	fmt.Printf("\n🎯 Collaboration Score: %.1f/100\n", metrics.CollaborationScore)
	fmt.Printf("   Contributors: %d | Active Reviewers: %d\n\n",
		metrics.TotalContributors, metrics.ReviewParticipation.UniqueReviewers)

	// Review Participation
	fmt.Println("👀 REVIEW PARTICIPATION")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("   Total Reviews: %d\n", metrics.ReviewParticipation.TotalReviews)
	fmt.Printf("   Unique Reviewers: %d\n", metrics.ReviewParticipation.UniqueReviewers)

	if metrics.ReviewParticipation.MostActiveReviewer != "" {
		fmt.Printf("   Most Active Reviewer: %s\n", metrics.ReviewParticipation.MostActiveReviewer)
	}
	if metrics.ReviewParticipation.LeastActiveReviewer != "" {
		fmt.Printf("   Least Active Reviewer: %s\n", metrics.ReviewParticipation.LeastActiveReviewer)
	}

	// Top Reviewers
	if len(metrics.ReviewParticipation.ReviewsPerContributor) > 0 {
		fmt.Println("\n   Top Reviewers:")
		sorted := sortMapByFloatValue(metrics.ReviewParticipation.ReviewsPerContributor)
		count := 0
		for _, reviewer := range sorted {
			if count >= 5 {
				break
			}
			fmt.Printf("      • %s: %d reviews\n", reviewer.Key, reviewer.Value)
			count++
		}
	}

	// Cross-Contributor Patterns
	fmt.Println("\n📁 CROSS-CONTRIBUTOR PATTERNS")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("   Files Analyzed: %d\n", metrics.CrossContributorPatterns.TotalFiles)
	fmt.Printf("   Files by Multiple Contributors: %d (%.1f%%)\n",
		metrics.CrossContributorPatterns.FilesByMultiple,
		metrics.CrossContributorPatterns.CollaborationRatio)
	fmt.Printf("   Files by Single Contributor: %d\n",
		metrics.CrossContributorPatterns.FilesBySingle)

	// Network Metrics
	fmt.Println("\n🔗 NETWORK METRICS")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("   Network Density: %.2f\n", metrics.NetworkMetrics.NetworkDensity)
	fmt.Printf("   Clustering Coefficient: %.2f\n", metrics.NetworkMetrics.ClusteringCoefficient)
	fmt.Printf("   Connected Components: %d\n", metrics.NetworkMetrics.ConnectedComponents)

	// Centrality (Top 5)
	if len(metrics.NetworkMetrics.DegreeCentrality) > 0 {
		fmt.Println("\n   Top Contributors by Centrality:")
		sortedCentrality := sortMapByFloat64Value(metrics.NetworkMetrics.DegreeCentrality)
		count := 0
		for _, contributor := range sortedCentrality {
			if count >= 5 {
				break
			}
			fmt.Printf("      • %s: %.2f\n", contributor.Key, contributor.Value)
			count++
		}
	}

	// Review Timing
	fmt.Println("\n⏱️  REVIEW TIMING")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("   Avg Time to First Review: %s\n",
		formatDuration(metrics.ReviewTiming.AverageTimeToFirstReview))
	fmt.Printf("   Median Time to First Review: %s\n",
		formatDuration(metrics.ReviewTiming.MedianTimeToFirstReview))
	fmt.Printf("   Avg Time to Merge: %s\n",
		formatDuration(metrics.ReviewTiming.AverageTimeToMerge))
	fmt.Printf("   Reviews per PR: %.1f\n", metrics.ReviewTiming.ReviewsPerPR)
	fmt.Printf("   PRs with Reviews: %.1f%%\n", metrics.ReviewTiming.PRsWithReviews)

	// Knowledge Silos
	if len(metrics.KnowledgeSilos) > 0 {
		fmt.Println("\n⚠️  KNOWLEDGE SILOS (High Risk)")
		fmt.Println(strings.Repeat("-", 40))
		for i, silo := range metrics.KnowledgeSilos {
			if i >= 5 {
				break
			}
			icon := "🟡"
			if silo.RiskLevel == "Critical" {
				icon = "🔴"
			} else if silo.RiskLevel == "High" {
				icon = "🟠"
			}
			fmt.Printf("   %s %s (Owner: %s)\n", icon, silo.Filename, silo.Owner)
		}
	}

	// Recommendations
	fmt.Println("\n💡 RECOMMENDATIONS")
	fmt.Println(strings.Repeat("-", 40))
	for _, rec := range metrics.Recommendations {
		fmt.Printf("   %s\n", rec)
	}

	fmt.Println()
}

// PrintCollaborationNetwork prints an ASCII network graph
func PrintCollaborationNetwork(metrics *analyzer.CollaborationMetrics) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🔗 COLLABORATION NETWORK (ASCII Visualization)")
	fmt.Println(strings.Repeat("=", 60))

	if len(metrics.ReviewParticipation.ReviewNetwork) == 0 {
		fmt.Println("   No review network data available.")
		fmt.Println()
		return
	}

	// Build edge list
	edges := metrics.ReviewParticipation.ReviewNetwork
	if len(edges) > 20 {
		edges = edges[:20]
	}

	fmt.Println("\n   Review Relationships (who reviewed whose PRs):")
	fmt.Println()

	// Print simple ASCII art
	for _, edge := range edges {
		stars := ""
		if edge.Count > 3 {
			stars = strings.Repeat("★", int(math.Min(float64(edge.Count), 5)))
		}
		fmt.Printf("   %s --[%d reviews]--> %s %s\n",
			padRight(edge.Reviewer, 15),
			edge.Count,
			padRight(edge.Author, 15),
			stars)
	}

	fmt.Println()

	// Legend
	fmt.Println("   Legend: ★ = frequent reviewer (4+ reviews)")
	fmt.Println()
}

// PrintCollaborationSummary prints a brief summary
func PrintCollaborationSummary(metrics *analyzer.CollaborationMetrics) {
	fmt.Println("\n📊 Collaboration Summary")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("   Score: %.1f/100\n", metrics.CollaborationScore)
	fmt.Printf("   Contributors: %d\n", metrics.TotalContributors)
	fmt.Printf("   Reviewers: %d\n", metrics.ReviewParticipation.UniqueReviewers)
	fmt.Printf("   Network Density: %.2f\n", metrics.NetworkMetrics.NetworkDensity)
	fmt.Printf("   Avg Review Time: %s\n", formatDuration(metrics.ReviewTiming.AverageTimeToFirstReview))
	fmt.Println()
}

// Helper functions

func sortMapByFloatValue(m map[string]int) []struct {
	Key   string
	Value int
} {
	var entries []struct {
		Key   string
		Value int
	}
	for k, v := range m {
		entries = append(entries, struct {
			Key   string
			Value int
		}{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
	return entries
}

func sortMapByFloat64Value(m map[string]float64) []struct {
	Key   string
	Value float64
} {
	var entries []struct {
		Key   string
		Value float64
	}
	for k, v := range m {
		entries = append(entries, struct {
			Key   string
			Value float64
		}{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
	return entries
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}
