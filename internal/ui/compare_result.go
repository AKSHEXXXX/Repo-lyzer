package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CompareResultModel struct {
	result *CompareResult
	err    error
}

func NewCompareResultModel() CompareResultModel {
	return CompareResultModel{}
}

func (m CompareResultModel) Update(msg tea.Msg) (CompareResultModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, func() tea.Msg { return BackToMenuMsg{} }
		case "j":
			// Export comparison to JSON
			if m.result != nil && m.result.Repo1.Repo != nil && m.result.Repo2.Repo != nil {
				return m, func() tea.Msg {
					_, err := ExportCompareJSON(*m.result)
					if err != nil {
						return ErrorMsg(fmt.Errorf("failed to export JSON: %w", err))
					}
					return StatusMsg{Message: "✓ Exported comparison to JSON successfully", IsError: false}
				}
			} else {
				return m, func() tea.Msg {
					return ErrorMsg(fmt.Errorf("no comparison data available for export"))
				}
			}
		case "m":
			// Export comparison to Markdown
			if m.result != nil && m.result.Repo1.Repo != nil && m.result.Repo2.Repo != nil {
				return m, func() tea.Msg {
					_, err := ExportCompareMarkdown(*m.result)
					if err != nil {
						return ErrorMsg(fmt.Errorf("failed to export Markdown: %w", err))
					}
					return StatusMsg{Message: "✓ Exported comparison to Markdown successfully", IsError: false}
				}
			} else {
				return m, func() tea.Msg {
					return ErrorMsg(fmt.Errorf("no comparison data available for export"))
				}
			}
		}
	}
	return m, nil
}

func (m CompareResultModel) View(width, height int) string {
	if m.result == nil || m.result.Repo1.Repo == nil || m.result.Repo2.Repo == nil {
		return "No comparison data"
	}

	r1 := m.result.Repo1
	r2 := m.result.Repo2

	header := TitleStyle.Render(fmt.Sprintf("📊 Comparison: %s vs %s", r1.Repo.FullName, r2.Repo.FullName))

	// Check if repositories are identical
	if r1.Repo.Stars == r2.Repo.Stars &&
		r1.Repo.Forks == r2.Repo.Forks &&
		len(r1.Commits) == len(r2.Commits) &&
		len(r1.Contributors) == len(r2.Contributors) &&
		r1.BusFactor == r2.BusFactor &&
		r1.MaturityScore == r2.MaturityScore {

		noDiffBox := BoxStyle.Render("✅ No differences found between the two repositories.\nBoth repositories have identical metrics.")
		footer := SubtleStyle.Render("j: export JSON • m: export Markdown • q/ESC: back to menu")

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			noDiffBox,
			footer,
		)

		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	// Build comparison table
	rows := []string{
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "Metric", r1.Repo.FullName, r2.Repo.FullName),
		strings.Repeat("─", 75),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "⭐ Stars", r1.Repo.Stars, r2.Repo.Forks),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "🍴 Forks", r1.Repo.Forks, r2.Repo.Forks),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "📦 Commits (1y)", len(r1.Commits), len(r2.Commits)),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "👥 Contributors", len(r1.Contributors), len(r2.Contributors)),
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "💚 Health Score", fmt.Sprintf("%d", r1.HealthScore), fmt.Sprintf("%d", r2.HealthScore)),
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "⚠️ Bus Factor", fmt.Sprintf("%d (%s)", r1.BusFactor, r1.BusRisk), fmt.Sprintf("%d (%s)", r2.BusFactor, r2.BusRisk)),
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "🏗️ Maturity", fmt.Sprintf("%s (%d)", r1.MaturityLevel, r1.MaturityScore), fmt.Sprintf("%s (%d)", r2.MaturityLevel, r2.MaturityScore)),
	}

	tableContent := strings.Join(rows, "\n")
	tableBox := BoxStyle.Render(tableContent)

	// Verdict
	var verdict string
	if r1.MaturityScore > r2.MaturityScore {
		verdict = fmt.Sprintf("➡️ %s appears more mature and stable.", r1.Repo.FullName)
	} else if r2.MaturityScore > r1.MaturityScore {
		verdict = fmt.Sprintf("➡️ %s appears more mature and stable.", r2.Repo.FullName)
	} else {
		verdict = "➡️ Both repositories are similarly mature."
	}
	verdictBox := BoxStyle.Render("📌 Verdict\n" + verdict)

	footer := SubtleStyle.Render("j: export JSON • m: export Markdown • q/ESC: back to menu")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
		verdictBox,
		footer,
	)

	// Add error/status message if present
	if m.err != nil {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			"\n"+ErrorStyle.Render(fmt.Sprintf("Status: %v", m.err)),
		)
	}

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m *CompareResultModel) SetResult(result *CompareResult) {
	m.result = result
}

func (m *CompareResultModel) SetError(err error) {
	m.err = err
}
