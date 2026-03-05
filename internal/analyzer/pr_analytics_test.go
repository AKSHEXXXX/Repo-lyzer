package analyzer

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestAnalyzePullRequests_EmptyList(t *testing.T) {
	prs := []github.PullRequest{}
	reviews := make(map[int][]github.Review)

	analytics := AnalyzePullRequests(prs, reviews)

	if analytics.TotalPRs != 0 {
		t.Errorf("Expected TotalPRs to be 0, got %d", analytics.TotalPRs)
	}
	if analytics.MergedPRs != 0 {
		t.Errorf("Expected MergedPRs to be 0, got %d", analytics.MergedPRs)
	}
}

func TestAnalyzePullRequests_OnlyOpenPRs(t *testing.T) {
	now := time.Now()
	prs := []github.PullRequest{
		{
			Number:    1,
			State:     "open",
			CreatedAt: now.Add(-24 * time.Hour),
			User:      github.User{Login: "user1"},
			Additions: 50,
			Deletions: 20,
		},
		{
			Number:    2,
			State:     "open",
			CreatedAt: now.Add(-48 * time.Hour),
			User:      github.User{Login: "user2"},
			Additions: 200,
			Deletions: 100,
		},
	}
	reviews := make(map[int][]github.Review)

	analytics := AnalyzePullRequests(prs, reviews)

	if analytics.TotalPRs != 2 {
		t.Errorf("Expected TotalPRs to be 2, got %d", analytics.TotalPRs)
	}
	if analytics.OpenPRs != 2 {
		t.Errorf("Expected OpenPRs to be 2, got %d", analytics.OpenPRs)
	}
	if analytics.MergedPRs != 0 {
		t.Errorf("Expected MergedPRs to be 0, got %d", analytics.MergedPRs)
	}
	if analytics.AverageTimeToMerge != 0 {
		t.Errorf("Expected AverageTimeToMerge to be 0, got %v", analytics.AverageTimeToMerge)
	}
}

func TestAnalyzePullRequests_MergedPRs(t *testing.T) {
	now := time.Now()
	created := now.Add(-72 * time.Hour)
	merged := now.Add(-24 * time.Hour)

	prs := []github.PullRequest{
		{
			Number:    1,
			State:     "closed",
			CreatedAt: created,
			MergedAt:  &merged,
			User:      github.User{Login: "user1"},
			Additions: 50,
			Deletions: 20,
		},
	}
	reviews := make(map[int][]github.Review)

	analytics := AnalyzePullRequests(prs, reviews)

	if analytics.TotalPRs != 1 {
		t.Errorf("Expected TotalPRs to be 1, got %d", analytics.TotalPRs)
	}
	if analytics.MergedPRs != 1 {
		t.Errorf("Expected MergedPRs to be 1, got %d", analytics.MergedPRs)
	}
	if analytics.AverageTimeToMerge <= 0 {
		t.Errorf("Expected positive AverageTimeToMerge, got %v", analytics.AverageTimeToMerge)
	}

	expectedDuration := merged.Sub(created)
	if analytics.AverageTimeToMerge != expectedDuration {
		t.Errorf("Expected AverageTimeToMerge to be %v, got %v", expectedDuration, analytics.AverageTimeToMerge)
	}
}

func TestAnalyzePullRequests_AbandonedRatio(t *testing.T) {
	now := time.Now()
	closed := now.Add(-24 * time.Hour)
	merged := now.Add(-24 * time.Hour)

	prs := []github.PullRequest{
		{
			Number:    1,
			State:     "closed",
			CreatedAt: now.Add(-72 * time.Hour),
			MergedAt:  &merged,
			User:      github.User{Login: "user1"},
		},
		{
			Number:    2,
			State:     "closed",
			CreatedAt: now.Add(-72 * time.Hour),
			ClosedAt:  &closed,
			User:      github.User{Login: "user2"},
		},
		{
			Number:    3,
			State:     "closed",
			CreatedAt: now.Add(-72 * time.Hour),
			ClosedAt:  &closed,
			User:      github.User{Login: "user3"},
		},
	}
	reviews := make(map[int][]github.Review)

	analytics := AnalyzePullRequests(prs, reviews)

	if analytics.MergedPRs != 1 {
		t.Errorf("Expected MergedPRs to be 1, got %d", analytics.MergedPRs)
	}
	if analytics.ClosedPRs != 2 {
		t.Errorf("Expected ClosedPRs to be 2, got %d", analytics.ClosedPRs)
	}

	// 2 out of 3 closed PRs were abandoned (66.67%)
	expectedRatio := 66.67
	if analytics.AbandonedRatio < expectedRatio-0.1 || analytics.AbandonedRatio > expectedRatio+0.1 {
		t.Errorf("Expected AbandonedRatio around %.2f%%, got %.2f%%", expectedRatio, analytics.AbandonedRatio)
	}
}

func TestAnalyzePullRequests_ReviewParticipation(t *testing.T) {
	now := time.Now()

	prs := []github.PullRequest{
		{Number: 1, State: "open", CreatedAt: now, User: github.User{Login: "author1"}},
		{Number: 2, State: "open", CreatedAt: now, User: github.User{Login: "author2"}},
		{Number: 3, State: "open", CreatedAt: now, User: github.User{Login: "author3"}},
		{Number: 4, State: "open", CreatedAt: now, User: github.User{Login: "author4"}},
	}

	reviews := map[int][]github.Review{
		1: {
			{User: github.User{Login: "reviewer1"}, SubmittedAt: now},
			{User: github.User{Login: "reviewer2"}, SubmittedAt: now},
		},
		2: {
			{User: github.User{Login: "reviewer1"}, SubmittedAt: now},
		},
		3: {
			{User: github.User{Login: "reviewer1"}, SubmittedAt: now},
			{User: github.User{Login: "reviewer2"}, SubmittedAt: now},
			{User: github.User{Login: "reviewer3"}, SubmittedAt: now},
		},
		// PR 4 has no reviews
	}

	analytics := AnalyzePullRequests(prs, reviews)

	// 2 out of 4 PRs have 2+ reviewers (50%)
	expectedParticipation := 50.0
	if analytics.ReviewParticipation != expectedParticipation {
		t.Errorf("Expected ReviewParticipation to be %.1f%%, got %.1f%%", expectedParticipation, analytics.ReviewParticipation)
	}
}

func TestAnalyzePullRequests_PRSizeDistribution(t *testing.T) {
	now := time.Now()

	prs := []github.PullRequest{
		{Number: 1, State: "open", CreatedAt: now, User: github.User{Login: "user1"}, Additions: 30, Deletions: 20},    // small (50)
		{Number: 2, State: "open", CreatedAt: now, User: github.User{Login: "user2"}, Additions: 200, Deletions: 150},  // medium (350)
		{Number: 3, State: "open", CreatedAt: now, User: github.User{Login: "user3"}, Additions: 600, Deletions: 200},  // large (800)
		{Number: 4, State: "open", CreatedAt: now, User: github.User{Login: "user4"}, Additions: 1500, Deletions: 500}, // xlarge (2000)
		{Number: 5, State: "open", CreatedAt: now, User: github.User{Login: "user5"}, Additions: 50, Deletions: 30},    // small (80)
	}
	reviews := make(map[int][]github.Review)

	analytics := AnalyzePullRequests(prs, reviews)

	if analytics.PRSizeDistribution["small"] != 2 {
		t.Errorf("Expected 2 small PRs, got %d", analytics.PRSizeDistribution["small"])
	}
	if analytics.PRSizeDistribution["medium"] != 1 {
		t.Errorf("Expected 1 medium PR, got %d", analytics.PRSizeDistribution["medium"])
	}
	if analytics.PRSizeDistribution["large"] != 1 {
		t.Errorf("Expected 1 large PR, got %d", analytics.PRSizeDistribution["large"])
	}
	if analytics.PRSizeDistribution["xlarge"] != 1 {
		t.Errorf("Expected 1 xlarge PR, got %d", analytics.PRSizeDistribution["xlarge"])
	}
}

func TestAnalyzePullRequests_FirstTimeContributors(t *testing.T) {
	now := time.Now()
	merged := now.Add(-24 * time.Hour)
	closed := now.Add(-68 * time.Hour) // Add closed timestamp

	prs := []github.PullRequest{
		// first-time, accepted
		{Number: 1, State: "closed", CreatedAt: now.Add(-72 * time.Hour), MergedAt: &merged, User: github.User{Login: "user1"}, AuthorAssociation: "FIRST_TIME_CONTRIBUTOR"},
		// not first-time (member)
		{Number: 2, State: "closed", CreatedAt: now.Add(-90 * time.Hour), MergedAt: &merged, User: github.User{Login: "user2"}, AuthorAssociation: "MEMBER"},
		{Number: 3, State: "closed", CreatedAt: now.Add(-85 * time.Hour), MergedAt: &merged, User: github.User{Login: "user2"}, AuthorAssociation: "MEMBER"},
		// first-time, not accepted
		{Number: 4, State: "closed", CreatedAt: now.Add(-70 * time.Hour), ClosedAt: &closed, User: github.User{Login: "user3"}, AuthorAssociation: "FIRST_TIME_CONTRIBUTOR"},
	}

	firstReview := now.Add(-70 * time.Hour)
	reviews := map[int][]github.Review{
		1: {{User: github.User{Login: "reviewer1"}, SubmittedAt: firstReview}},
	}

	analytics := AnalyzePullRequests(prs, reviews)

	// 2 first-time contributors (user1 and user3)
	if analytics.FirstTimeContributorMetrics.TotalFirstTimePRs != 2 {
		t.Errorf("Expected 2 first-time PRs, got %d", analytics.FirstTimeContributorMetrics.TotalFirstTimePRs)
	}

	// 1 accepted (user1)
	if analytics.FirstTimeContributorMetrics.AcceptedFirstTimePRs != 1 {
		t.Errorf("Expected 1 accepted first-time PR, got %d", analytics.FirstTimeContributorMetrics.AcceptedFirstTimePRs)
	}

	// Acceptance rate should be 50%
	expectedRate := 50.0
	if analytics.FirstTimeContributorMetrics.AcceptanceRate != expectedRate {
		t.Errorf("Expected acceptance rate %.1f%%, got %.1f%%", expectedRate, analytics.FirstTimeContributorMetrics.AcceptanceRate)
	}
}

func TestClassifyPRSize(t *testing.T) {
	tests := []struct {
		name      string
		additions int
		deletions int
		expected  string
	}{
		{"small PR", 30, 20, "small"},
		{"medium PR", 200, 150, "medium"},
		{"large PR", 600, 200, "large"},
		{"xlarge PR", 1500, 500, "xlarge"},
		{"boundary small", 50, 49, "small"},
		{"boundary medium", 250, 249, "medium"},
		{"boundary large", 500, 499, "large"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyPRSize(tt.additions, tt.deletions)
			if result != tt.expected {
				t.Errorf("classifyPRSize(%d, %d) = %s, want %s", tt.additions, tt.deletions, result, tt.expected)
			}
		})
	}
}
