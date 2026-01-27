package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const historyFile = "exports/history.json"
const maxHistoryItems = 50

// HistoryEntry represents a saved analysis
type HistoryEntry struct {
	RepoName      string    `json:"repo_name"`
	AnalyzedAt    time.Time `json:"analyzed_at"`
	HealthScore   int       `json:"health_score"`
	Stars         int       `json:"stars"`
	Forks         int       `json:"forks"`
	MaturityLevel string    `json:"maturity_level"`
}

// History holds all history entries
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// HistoryModel wraps History and implements tea.Model
type HistoryModel struct {
	*History
	cursor int
}

// NewHistoryModel creates a new HistoryModel
func NewHistoryModel() HistoryModel {
	return HistoryModel{
		History: &History{Entries: []HistoryEntry{}},
		cursor:  0,
	}
}

// LoadHistory loads history from file and returns HistoryModel
func LoadHistory() (HistoryModel, error) {
	history := &History{Entries: []HistoryEntry{}}

	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return NewHistoryModel(), nil
		}
		return NewHistoryModel(), err
	}

	if err := json.Unmarshal(data, history); err != nil {
		return NewHistoryModel(), nil // Return empty on parse error
	}

	return HistoryModel{
		History: history,
		cursor:  0,
	}, nil
}

// Save saves history to file
func (h HistoryModel) Save() error {
	if err := os.MkdirAll(filepath.Dir(historyFile), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(h.History, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

// AddEntry adds a new entry to history
func (h *HistoryModel) AddEntry(data AnalysisResult) {
	entry := HistoryEntry{
		RepoName:      data.Repo.FullName,
		AnalyzedAt:    time.Now(),
		HealthScore:   data.HealthScore,
		Stars:         data.Repo.Stars,
		Forks:         data.Repo.Forks,
		MaturityLevel: data.MaturityLevel,
	}

	// Remove duplicate if exists
	for i, e := range h.Entries {
		if e.RepoName == entry.RepoName {
			h.Entries = append(h.Entries[:i], h.Entries[i+1:]...)
			break
		}
	}

	// Add to front
	h.Entries = append([]HistoryEntry{entry}, h.Entries...)

	// Trim to max size
	if len(h.Entries) > maxHistoryItems {
		h.Entries = h.Entries[:maxHistoryItems]
	}
}

// Delete removes a specific entry
func (h *HistoryModel) Delete(index int) {
	if index >= 0 && index < len(h.Entries) {
		h.Entries = append(h.Entries[:index], h.Entries[index+1:]...)
	}
}

// Clear removes all history
func (h *HistoryModel) Clear() {
	h.Entries = []HistoryEntry{}
}

// Init implements tea.Model
func (h HistoryModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (h HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if h.cursor > 0 {
				h.cursor--
			}
		case "down", "j":
			if h.cursor < len(h.Entries)-1 {
				h.cursor++
			}
		}
	}
	return h, nil
}

// View implements tea.Model
func (h HistoryModel) View() string {
	return h.ViewWithSize(0, 0)
}

// ViewWithSize returns the view with specified dimensions
func (h HistoryModel) ViewWithSize(width, height int) string {
	header := TitleStyle.Render("📜 Analysis History")

	if len(h.Entries) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render("No history yet. Analyze a repository to get started!"),
			SubtleStyle.Render("q/ESC: back to menu"),
		)

		if width == 0 {
			return content
		}

		return lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			content,
		)
	}

	// Build history list
	var lines []string
	lines = append(lines, fmt.Sprintf("%-30s │ %-8s │ %-5s │ %-12s │ %s", "Repository", "Stars", "Health", "Maturity", "Analyzed"))
	lines = append(lines, strings.Repeat("─", 85))

	for i, entry := range h.Entries {
		prefix := "  "
		if i == h.cursor {
			prefix = "▶ "
		}
		line := fmt.Sprintf("%s%-28s │ ⭐%-6d │ 💚%-3d │ %-12s │ %s",
			prefix,
			entry.RepoName,
			entry.Stars,
			entry.HealthScore,
			entry.MaturityLevel,
			entry.AnalyzedAt.Format("2006-01-02 15:04"),
		)
		if i == h.cursor {
			lines = append(lines, SelectedStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	tableBox := BoxStyle.Render(strings.Join(lines, "\n"))

	footer := SubtleStyle.Render("↑↓: navigate • Enter: re-analyze • d: delete • c: clear all • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
		footer,
	)

	if width == 0 {
		return content
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// GetRecent returns the most recent entries
func (h *History) GetRecent(count int) []HistoryEntry {
	if count > len(h.Entries) {
		count = len(h.Entries)
	}
	return h.Entries[:count]
}

// SortByDate sorts entries by date (newest first)
func (h *History) SortByDate() {
	sort.Slice(h.Entries, func(i, j int) bool {
		return h.Entries[i].AnalyzedAt.After(h.Entries[j].AnalyzedAt)
	})
}

// FormatEntry formats a history entry for display
func (e HistoryEntry) Format() string {
	return fmt.Sprintf("%-30s │ ⭐%-6d │ 💚%-3d │ %s │ %s",
		e.RepoName,
		e.Stars,
		e.HealthScore,
		e.MaturityLevel,
		e.AnalyzedAt.Format("2006-01-02 15:04"),
	)
}
