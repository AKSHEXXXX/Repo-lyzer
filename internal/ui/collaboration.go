// Package ui provides user interface components for Repo-lyzer.
// This file implements collaboration insights UI components.
package ui

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
)

// CollaborationDashboard represents the collaboration insights dashboard
type CollaborationDashboard struct {
	Metrics *analyzer.CollaborationMetrics
}

// NewCollaborationDashboard creates a new collaboration dashboard
func NewCollaborationDashboard(metrics *analyzer.CollaborationMetrics) *CollaborationDashboard {
	return &CollaborationDashboard{
		Metrics: metrics,
	}
}

// Render renders the collaboration dashboard
func (d *CollaborationDashboard) Render() string {
	output := fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║           TEAM COLLABORATION INSIGHTS                      ║
╠══════════════════════════════════════════════════════════════╣
║  Collaboration Score: %.1f/100                             ║
║  Contributors: %d | Active Reviewers: %d                    ║
╠══════════════════════════════════════════════════════════════╣
`, d.Metrics.CollaborationScore, d.Metrics.TotalContributors, d.Metrics.ReviewParticipation.UniqueReviewers)

	// Review Participation
	output += fmt.Sprintf(`║  👀 REVIEW PARTICIPATION                                   ║
║     Total Reviews: %d                                      ║
║     Unique Reviewers: %d                                   ║
`, d.Metrics.ReviewParticipation.TotalReviews, d.Metrics.ReviewParticipation.UniqueReviewers)

	// Network Metrics
	output += fmt.Sprintf(`║  🔗 NETWORK METRICS                                        ║
║     Network Density: %.2f                                  ║
║     Clustering Coefficient: %.2f                             ║
║     Connected Components: %d                                 ║
`, d.Metrics.NetworkMetrics.NetworkDensity, d.Metrics.NetworkMetrics.ClusteringCoefficient, d.Metrics.NetworkMetrics.ConnectedComponents)

	// Review Timing
	output += fmt.Sprintf(`║  ⏱️  REVIEW TIMING                                          ║
║     Avg Time to First Review: %s                            ║
║     Avg Time to Merge: %s                                   ║
║     Reviews per PR: %.1f                                    ║
`, formatDurationUI(d.Metrics.ReviewTiming.AverageTimeToFirstReview),
		formatDurationUI(d.Metrics.ReviewTiming.AverageTimeToMerge),
		d.Metrics.ReviewTiming.ReviewsPerPR)

	// Knowledge Silos
	if len(d.Metrics.KnowledgeSilos) > 0 {
		output += "║  ⚠️  KNOWLEDGE SILOS                                        ║\n"
		for i, silo := range d.Metrics.KnowledgeSilos {
			if i >= 3 {
				break
			}
			output += fmt.Sprintf("║     %s - Owner: %s                                     ║\n", silo.Filename, silo.Owner)
		}
	}

	// Recommendations
	output += "║  💡 RECOMMENDATIONS                                         ║\n"
	for _, rec := range d.Metrics.Recommendations {
		if len(rec) > 50 {
			rec = rec[:50] + "..."
		}
		output += fmt.Sprintf("║     %s                                    ║\n", rec)
	}

	output += "╚══════════════════════════════════════════════════════════════╝\n"

	return output
}

// NetworkVisualization provides ASCII network visualization
func (d *CollaborationDashboard) NetworkVisualization() string {
	if len(d.Metrics.ReviewParticipation.ReviewNetwork) == 0 {
		return "No network data available.\n"
	}

	output := "\n🔗 Collaboration Network:\n\n"

	// Simple node-edge representation
	edges := d.Metrics.ReviewParticipation.ReviewNetwork
	if len(edges) > 15 {
		edges = edges[:15]
	}

	for _, edge := range edges {
		output += fmt.Sprintf("   %s → %s (%d reviews)\n",
			edge.Reviewer, edge.Author, edge.Count)
	}

	return output
}

// InteractiveNetworkExplorer provides interactive exploration of the network
type InteractiveNetworkExplorer struct {
	Metrics *analyzer.CollaborationMetrics
	Current int
}

// NewInteractiveNetworkExplorer creates a new interactive network explorer
func NewInteractiveNetworkExplorer(metrics *analyzer.CollaborationMetrics) *InteractiveNetworkExplorer {
	return &InteractiveNetworkExplorer{
		Metrics: metrics,
		Current: 0,
	}
}

// Next shows the next page of the network
func (e *InteractiveNetworkExplorer) Next() string {
	edges := e.Metrics.ReviewParticipation.ReviewNetwork
	if len(edges) == 0 {
		return "No network data available.\n"
	}

	pageSize := 10
	start := e.Current * pageSize
	end := start + pageSize
	if start >= len(edges) {
		e.Current = 0
		start = 0
		end = pageSize
	}
	if end > len(edges) {
		end = len(edges)
	}

	output := fmt.Sprintf("\n📄 Page %d of %d:\n\n", e.Current+1, (len(edges)+pageSize-1)/pageSize)

	for i := start; i < end; i++ {
		edge := edges[i]
		output += fmt.Sprintf("   %d. %s → %s (%d reviews)\n", i+1, edge.Reviewer, edge.Author, edge.Count)
	}

	e.Current++
	return output
}

// GetContributorDetails gets details about a specific contributor
func (e *InteractiveNetworkExplorer) GetContributorDetails(login string) string {
	output := fmt.Sprintf("\n👤 Contributor: %s\n\n", login)

	// Find reviews given by this contributor
	reviewsGiven := 0
	reviewsReceived := 0

	for _, edge := range e.Metrics.ReviewParticipation.ReviewNetwork {
		if edge.Reviewer == login {
			reviewsGiven += edge.Count
		}
		if edge.Author == login {
			reviewsReceived += edge.Count
		}
	}

	output += fmt.Sprintf("   Reviews Given: %d\n", reviewsGiven)
	output += fmt.Sprintf("   Reviews Received: %d\n", reviewsReceived)

	// Centrality
	centrality, ok := e.Metrics.NetworkMetrics.DegreeCentrality[login]
	if ok {
		output += fmt.Sprintf("   Degree Centrality: %.2f\n", centrality)
	}

	betweenness, ok := e.Metrics.NetworkMetrics.BetweennessCentrality[login]
	if ok {
		output += fmt.Sprintf("   Betweenness Centrality: %.2f\n", betweenness)
	}

	return output
}

// formatDurationUI formats duration for UI display
func formatDurationUI(d interface{}) string {
	// This is a placeholder - in production would use time.Duration formatting
	switch v := d.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", d)
	}
}
