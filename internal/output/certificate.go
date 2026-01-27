package output

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
	"github.com/jung-kurt/gofpdf"
)

var (
	certTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			Align(lipgloss.Center).
			Margin(1, 0)

	certSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FF87")).
				MarginTop(1)

	certKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00E5FF")).
			Bold(true)

	certValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	scoreStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB000")).
			Bold(true)

	gradeStyle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center)
)

// PrintCertificate displays a formatted repository certificate
func PrintCertificate(cert *analyzer.CertificateData) {
	// Title
	fmt.Println(certTitleStyle.Render("🏆 REPOSITORY CERTIFICATE 🏆"))
	fmt.Printf("Repository: %s/%s\n", cert.Owner, cert.RepoName)
	fmt.Println(strings.Repeat("=", 60))

	// Repository Information
	fmt.Println(certSectionStyle.Render("📋 Repository Information"))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Description"), certValueStyle.Render(cert.Description))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Stars"), certValueStyle.Render(fmt.Sprintf("%d", cert.Stars)))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Forks"), certValueStyle.Render(fmt.Sprintf("%d", cert.Forks)))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Open Issues"), certValueStyle.Render(fmt.Sprintf("%d", cert.OpenIssues)))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Created"), certValueStyle.Render(cert.CreatedAt))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Last Updated"), certValueStyle.Render(cert.UpdatedAt))
	fmt.Printf("%s: %s (%d languages)\n", certKeyStyle.Render("Primary Language"), certValueStyle.Render(cert.PrimaryLanguage), cert.LanguageCount)

	// Scores
	fmt.Println(certSectionStyle.Render("📊 Scores & Metrics"))
	fmt.Printf("%s: %s/100\n", certKeyStyle.Render("Health Score"), scoreStyle.Render(fmt.Sprintf("%d", cert.HealthScore)))
	fmt.Printf("%s: %s/100 (%s)\n", certKeyStyle.Render("Maturity Score"), scoreStyle.Render(fmt.Sprintf("%d", cert.MaturityScore)), certValueStyle.Render(cert.MaturityLevel))
	fmt.Printf("%s: %s (%s)\n", certKeyStyle.Render("Bus Factor"), scoreStyle.Render(fmt.Sprintf("%d", cert.BusFactor)), certValueStyle.Render(cert.BusRisk))
	fmt.Printf("%s: %s (%s)\n", certKeyStyle.Render("Commits (Last Year)"), certValueStyle.Render(fmt.Sprintf("%d", cert.CommitsLastYear)), certValueStyle.Render(cert.ActivityLevel))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Contributors"), certValueStyle.Render(fmt.Sprintf("%d", cert.Contributors)))

	// Overall Assessment
	fmt.Println(certSectionStyle.Render("🎯 Overall Assessment"))
	fmt.Printf("%s: %s\n", certKeyStyle.Render("Overall Score"), scoreStyle.Render(fmt.Sprintf("%d/100", cert.OverallScore)))
	fmt.Printf("%s: ", certKeyStyle.Render("Grade"))
	switch cert.Grade {
	case "A+", "A":
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#00FF00")).Render(cert.Grade))
	case "B+", "B":
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#FFFF00")).Render(cert.Grade))
	case "C+", "C":
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#FFA500")).Render(cert.Grade))
	default:
		fmt.Println(gradeStyle.Foreground(lipgloss.Color("#FF0000")).Render(cert.Grade))
	}

	// Potential Uses
	fmt.Println(certSectionStyle.Render("💡 Potential Uses"))
	for i, use := range cert.Uses {
		fmt.Printf("%d. %s\n", i+1, certValueStyle.Render(use))
	}

	fmt.Println(strings.Repeat("=", 60))
}

// getDownloadsDir returns the user's Downloads folder
func getDownloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Downloads"), nil
}

// openFileManager opens the file manager to show the exported file
func openFileManager(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows: use explorer with /select to highlight the file
		cmd = exec.Command("explorer", "/select,"+filePath)
	case "darwin":
		// macOS: use open with -R to reveal in Finder
		cmd = exec.Command("open", "-R", filePath)
	case "linux":
		// Linux: use xdg-open on the directory
		dir := filepath.Dir(filePath)
		cmd = exec.Command("xdg-open", dir)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// generateFilename creates a filename with repo name and timestamp
func generateFilename(repoName, ext string) string {
	safeName := strings.ReplaceAll(repoName, "/", "_")
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s_certificate_%s.%s", safeName, timestamp, ext)
}

// ExportCertificatePDF generates a PDF certificate and saves it to Downloads
func ExportCertificatePDF(cert *analyzer.CertificateData) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	// Ensure the downloads directory exists
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		return "", err
	}

	filename := filepath.Join(downloadsDir, generateFilename(cert.Owner+"/"+cert.RepoName, "pdf"))

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 15, "Repository Certificate")
	pdf.Ln(20)

	// Repository info
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Repository Information")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(50, 8, "Repository:")
	pdf.Cell(0, 8, cert.Owner+"/"+cert.RepoName)
	pdf.Ln(6)
	pdf.Cell(50, 8, "Description:")
	pdf.MultiCell(0, 8, cert.Description, "", "", false)
	pdf.Ln(6)
	pdf.Cell(50, 8, "Stars:")
	pdf.Cell(0, 8, fmt.Sprintf("%d", cert.Stars))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Forks:")
	pdf.Cell(0, 8, fmt.Sprintf("%d", cert.Forks))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Open Issues:")
	pdf.Cell(0, 8, fmt.Sprintf("%d", cert.OpenIssues))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Created:")
	pdf.Cell(0, 8, cert.CreatedAt)
	pdf.Ln(6)
	pdf.Cell(50, 8, "Last Updated:")
	pdf.Cell(0, 8, cert.UpdatedAt)
	pdf.Ln(6)
	pdf.Cell(50, 8, "Primary Language:")
	pdf.Cell(0, 8, fmt.Sprintf("%s (%d languages)", cert.PrimaryLanguage, cert.LanguageCount))
	pdf.Ln(15)

	// Scores & Metrics
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Scores & Metrics")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(50, 8, "Health Score:")
	pdf.Cell(0, 8, fmt.Sprintf("%d/100", cert.HealthScore))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Maturity Score:")
	pdf.Cell(0, 8, fmt.Sprintf("%d/100 (%s)", cert.MaturityScore, cert.MaturityLevel))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Bus Factor:")
	pdf.Cell(0, 8, fmt.Sprintf("%d (%s)", cert.BusFactor, cert.BusRisk))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Commits (Last Year):")
	pdf.Cell(0, 8, fmt.Sprintf("%d (%s)", cert.CommitsLastYear, cert.ActivityLevel))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Contributors:")
	pdf.Cell(0, 8, fmt.Sprintf("%d", cert.Contributors))
	pdf.Ln(15)

	// Overall Assessment
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Overall Assessment")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(50, 8, "Overall Score:")
	pdf.Cell(0, 8, fmt.Sprintf("%d/100", cert.OverallScore))
	pdf.Ln(6)
	pdf.Cell(50, 8, "Grade:")
	pdf.Cell(0, 8, cert.Grade)
	pdf.Ln(15)

	// Potential Uses
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Potential Uses")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	for i, use := range cert.Uses {
		pdf.Cell(10, 8, fmt.Sprintf("%d.", i+1))
		pdf.MultiCell(0, 8, use, "", "", false)
		pdf.Ln(2)
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(0, 8, fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02 15:04:05")))

	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}

	_ = openFileManager(filename)

	return filename, nil
}
