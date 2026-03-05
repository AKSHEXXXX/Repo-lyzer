package analyzer

import (
	"sort"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// PRAnalytics contains comprehensive pull request analytics
type PRAnalytics struct {
	// Overall statistics
	TotalPRs  int
	MergedPRs int
	ClosedPRs int
	OpenPRs   int

	// Timing metrics
	AverageTimeToMerge time.Duration
	MedianTimeToMerge  time.Duration

	// Review metrics
	ReviewParticipation float64 // % of PRs with 2+ reviewers

	// Size distribution
	PRSizeDistribution map[string]int

	// Quality metrics
	AbandonedRatio float64 // % of PRs closed without merge

	// Contributor metrics
	FirstTimeContributorMetrics FirstTimeMetrics
}

// FirstTimeMetrics tracks metrics specific to first-time contributors
type FirstTimeMetrics struct {
	TotalFirstTimePRs    int
	AcceptedFirstTimePRs int
	AcceptanceRate       float64
	AvgTimeToFirstReview time.Duration
}

// AnalyzePullRequests analyzes pull requests and returns comprehensive metrics
func AnalyzePullRequests(prs []github.PullRequest, reviews map[int][]github.Review) *PRAnalytics {
	analytics := &PRAnalytics{
		TotalPRs:           len(prs),
		PRSizeDistribution: make(map[string]int),
	}

	if len(prs) == 0 {
		return analytics
	}

	var mergeTimes []time.Duration
	var reviewCounts []int
	var reviewTimesForFirstTime []time.Duration

	// Analyze each PR
	for _, pr := range prs {
		// Count states
		if pr.State == "open" {
			analytics.OpenPRs++
		} else if pr.MergedAt != nil {
			analytics.MergedPRs++

			// Calculate time to merge
			mergeTime := pr.MergedAt.Sub(pr.CreatedAt)
			mergeTimes = append(mergeTimes, mergeTime)
		} else if pr.ClosedAt != nil {
			analytics.ClosedPRs++
		}

		// Classify PR size
		size := classifyPRSize(pr.Additions, pr.Deletions)
		analytics.PRSizeDistribution[size]++

		// Check if first-time contributor using GitHub's author association
		if pr.AuthorAssociation == "FIRST_TIME_CONTRIBUTOR" {
			analytics.FirstTimeContributorMetrics.TotalFirstTimePRs++

			if pr.MergedAt != nil {
				analytics.FirstTimeContributorMetrics.AcceptedFirstTimePRs++
			}

			// Calculate time to first review for first-time contributors
			if prReviews, exists := reviews[pr.Number]; exists && len(prReviews) > 0 {
				firstReviewTime := prReviews[0].SubmittedAt.Sub(pr.CreatedAt)
				reviewTimesForFirstTime = append(reviewTimesForFirstTime, firstReviewTime)
			}
		}

		// Count reviews for participation metric
		if prReviews, exists := reviews[pr.Number]; exists {
			// Count unique reviewers
			uniqueReviewers := make(map[string]bool)
			for _, review := range prReviews {
				// Don't count the PR author as a reviewer
				if review.User.Login != pr.User.Login {
					uniqueReviewers[review.User.Login] = true
				}
			}
			reviewCounts = append(reviewCounts, len(uniqueReviewers))
		} else {
			reviewCounts = append(reviewCounts, 0)
		}
	}

	// Calculate average time to merge
	if len(mergeTimes) > 0 {
		var totalMergeTime time.Duration
		for _, t := range mergeTimes {
			totalMergeTime += t
		}
		analytics.AverageTimeToMerge = totalMergeTime / time.Duration(len(mergeTimes))

		// Calculate median time to merge
		sort.Slice(mergeTimes, func(i, j int) bool {
			return mergeTimes[i] < mergeTimes[j]
		})
		mid := len(mergeTimes) / 2
		if len(mergeTimes)%2 == 0 {
			analytics.MedianTimeToMerge = (mergeTimes[mid-1] + mergeTimes[mid]) / 2
		} else {
			analytics.MedianTimeToMerge = mergeTimes[mid]
		}
	}

	// Calculate review participation (% with 2+ reviewers)
	prsWithTwoOrMoreReviewers := 0
	for _, count := range reviewCounts {
		if count >= 2 {
			prsWithTwoOrMoreReviewers++
		}
	}
	if len(reviewCounts) > 0 {
		analytics.ReviewParticipation = float64(prsWithTwoOrMoreReviewers) / float64(len(reviewCounts)) * 100
	}

	// Calculate abandoned ratio (closed without merge)
	closedWithoutMerge := analytics.ClosedPRs
	totalClosed := analytics.MergedPRs + analytics.ClosedPRs
	if totalClosed > 0 {
		analytics.AbandonedRatio = float64(closedWithoutMerge) / float64(totalClosed) * 100
	}

	// Calculate first-time contributor metrics
	if analytics.FirstTimeContributorMetrics.TotalFirstTimePRs > 0 {
		analytics.FirstTimeContributorMetrics.AcceptanceRate =
			float64(analytics.FirstTimeContributorMetrics.AcceptedFirstTimePRs) /
				float64(analytics.FirstTimeContributorMetrics.TotalFirstTimePRs) * 100
	}

	if len(reviewTimesForFirstTime) > 0 {
		var totalReviewTime time.Duration
		for _, t := range reviewTimesForFirstTime {
			totalReviewTime += t
		}
		analytics.FirstTimeContributorMetrics.AvgTimeToFirstReview =
			totalReviewTime / time.Duration(len(reviewTimesForFirstTime))
	}

	return analytics
}

// classifyPRSize categorizes a PR by its size
func classifyPRSize(additions, deletions int) string {
	totalChanges := additions + deletions

	switch {
	case totalChanges < 100:
		return "small"
	case totalChanges < 500:
		return "medium"
	case totalChanges < 1000:
		return "large"
	default:
		return "xlarge"
	}
}
