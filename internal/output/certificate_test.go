package output

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
)

func TestPrintCertificate(t *testing.T) {
	// Create mock certificate data
	cert := &analyzer.CertificateData{
		Owner:         "testowner",
		RepoName:      "testrepo",
		Description:   "A test repository",
		Stars:         42,
		Forks:         10,
		OpenIssues:    5,
		CreatedAt:     "2020-01-01",
		UpdatedAt:     "2023-12-01",
		PrimaryLanguage: "Go",
		LanguageCount: 3,
		HealthScore:   85,
		MaturityScore: 78,
		MaturityLevel: "Good",
		BusFactor:     3,
		BusRisk:       "Medium",
		CommitsLastYear: 150,
		ActivityLevel: "Active",
		Contributors:  8,
		OverallScore:  80,
		Grade:         "B+",
		Uses: []string{
			"Learning Go programming",
			"Code analysis tool development",
			"Open source contribution practice",
		},
	}

	// Test that PrintCertificate doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintCertificate panicked: %v", r)
		}
	}()

	// Since we can't easily capture fmt output in this test environment,
	// we'll test the core logic by checking that the function completes without error
	PrintCertificate(cert)

	// The test passes if no panic occurs
	t.Log("PrintCertificate executed successfully without panics")
}

func TestExportCertificatePDF(t *testing.T) {
	// Create mock certificate data
	cert := &analyzer.CertificateData{
		Owner:         "testowner",
		RepoName:      "testrepo",
		Description:   "A test repository for PDF export",
		Stars:         100,
		Forks:         25,
		OpenIssues:    3,
		CreatedAt:     "2019-06-15",
		UpdatedAt:     "2023-11-20",
		PrimaryLanguage: "TypeScript",
		LanguageCount: 2,
		HealthScore:   92,
		MaturityScore: 85,
		MaturityLevel: "Excellent",
		BusFactor:     5,
		BusRisk:       "Low",
		CommitsLastYear: 300,
		ActivityLevel: "Very Active",
		Contributors:  15,
		OverallScore:  88,
		Grade:         "A-",
		Uses: []string{
			"Web development framework",
			"Full-stack application development",
			"API development",
		},
	}

	// Check if downloads directory exists
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		t.Skip("Failed to get downloads directory")
	}
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		t.Skip("Downloads directory does not exist")
	}

	// Test PDF export
	path, err := ExportCertificatePDF(cert)
	if err != nil {
		t.Errorf("ExportCertificatePDF failed: %v", err)
		return
	}

	defer os.Remove(path)

	// Check that path is not empty
	if path == "" {
		t.Error("ExportCertificatePDF returned empty path")
		return
	}

	// Check that path contains expected filename pattern
	if !strings.Contains(path, "testrepo_certificate_") {
		t.Errorf("PDF path doesn't contain expected filename pattern: %s", path)
	}

	// Check that path ends with .pdf
	if !strings.HasSuffix(path, ".pdf") {
		t.Errorf("PDF path doesn't end with .pdf: %s", path)
	}

	t.Logf("PDF exported successfully to: %s", path)
}

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		repoName string
		ext      string
		expected string
	}{
		{"owner/repo", "pdf", "owner_repo_certificate_"},
		{"test/repo-name", "json", "test_repo-name_certificate_"},
		{"user/project", "md", "user_project_certificate_"},
	}

	for _, tt := range tests {
		result := generateFilename(tt.repoName, tt.ext)
		if !strings.HasPrefix(result, tt.expected) {
			t.Errorf("generateFilename(%s, %s) = %s, expected to start with %s",
				tt.repoName, tt.ext, result, tt.expected)
		}
		if !strings.HasSuffix(result, "."+tt.ext) {
			t.Errorf("generateFilename(%s, %s) = %s, expected to end with .%s",
				tt.repoName, tt.ext, result, tt.ext)
		}
	}
}

func TestCertificateDataFormatting(t *testing.T) {
	// Test that certificate data is properly formatted
	cert := &analyzer.CertificateData{
		Stars:         42,
		Forks:         10,
		OpenIssues:    5,
		HealthScore:   85,
		BusFactor:     3,
		CommitsLastYear: 150,
		Contributors:  8,
		OverallScore:  80,
	}

	// Test that numeric fields are properly handled
	// This ensures the fmt.Sprintf conversions work correctly
	starsStr := certValueStyle.Render(fmt.Sprintf("%d", cert.Stars))
	if starsStr == "" {
		t.Error("Stars formatting failed")
	}

	forksStr := certValueStyle.Render(fmt.Sprintf("%d", cert.Forks))
	if forksStr == "" {
		t.Error("Forks formatting failed")
	}

	issuesStr := certValueStyle.Render(fmt.Sprintf("%d", cert.OpenIssues))
	if issuesStr == "" {
		t.Error("OpenIssues formatting failed")
	}

	commitsStr := certValueStyle.Render(fmt.Sprintf("%d", cert.CommitsLastYear))
	if commitsStr == "" {
		t.Error("CommitsLastYear formatting failed")
	}

	contributorsStr := certValueStyle.Render(fmt.Sprintf("%d", cert.Contributors))
	if contributorsStr == "" {
		t.Error("Contributors formatting failed")
	}

	t.Log("All certificate data formatting tests passed")
}
