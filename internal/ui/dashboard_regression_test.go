package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestDashboardQualityView_KeepsTopTabsVisibleWhenContentOverflows(t *testing.T) {
	model := NewDashboardModel()
	model.currentView = viewQualityDashboard
	model.width = 80
	model.height = 12

	repo := &github.Repo{
		FullName:      "owner/repo",
		Description:   "Regression test repository",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/owner/repo",
		CreatedAt:     time.Now(),
		PushedAt:      time.Now(),
	}

	model.SetData(AnalysisResult{
		Repo: repo,
		QualityDashboard: &analyzer.QualityDashboard{
			OverallScore: 69,
			RiskLevel:    "Medium",
			QualityGrade: "C",
			KeyMetrics: analyzer.DashboardMetrics{
				HealthScore:      90,
				SecurityScore:    100,
				MaturityLevel:    "Prototype",
				BusFactor:        2,
				ActivityLevel:    "Low",
				ContributorCount: 17,
			},
			ProblemHotspots: []analyzer.ProblemHotspot{
				{Area: "Bus Factor", Severity: "High", Description: "Very low contributor diversity"},
				{Area: "Security", Severity: "Medium", Description: "Security checks need hardening"},
				{Area: "Activity", Severity: "Medium", Description: "Low activity in the last 90 days"},
				{Area: "Testing", Severity: "Low", Description: "Missing coverage on critical paths"},
				{Area: "Docs", Severity: "Low", Description: "Insufficient onboarding docs"},
			},
			Recommendations: []string{
				"👥 Encourage more contributors to reduce bus factor risk",
				"📚 Improve documentation to enable easier onboarding",
				"🧪 Add regression coverage for UI layout behavior",
				"🔒 Improve security checks in CI",
				"📈 Increase regular maintenance cadence",
			},
		},
	})

	view := model.View()

	if !strings.Contains(view, "Overview") || !strings.Contains(view, "Quality") {
		t.Fatalf("top tabs not visible in quality view output:\n%s", view)
	}
}
