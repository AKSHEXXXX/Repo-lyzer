package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Favorite struct {
	RepoName  string    `json:"repo_name"`
	UseCount  int       `json:"use_count"`
	LastUsed  time.Time `json:"last_used"`
	AddedAt   time.Time `json:"added_at"`
	Notes     string    `json:"notes"`
}

type Favorites struct {
	Items []Favorite `json:"items"`
}

func (f *Favorites) Add(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items[i].UseCount++
			f.Items[i].LastUsed = time.Now()
			return
		}
	}
	// Add new favorite
	f.Items = append(f.Items, Favorite{
		RepoName: repoName,
		UseCount: 1,
		LastUsed: time.Now(),
		AddedAt:  time.Now(),
	})
}

func (f *Favorites) Remove(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items = append(f.Items[:i], f.Items[i+1:]...)
			return
		}
	}
}

func (f *Favorites) UpdateUsage(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items[i].UseCount++
			f.Items[i].LastUsed = time.Now()
			return
		}
	}
}

func (f *Favorites) IsFavorite(repoName string) bool {
	for _, item := range f.Items {
		if item.RepoName == repoName {
			return true
		}
	}
	return false
}

func (f *Favorites) GetTopFavorites(n int) []Favorite {
	if len(f.Items) <= n {
		return f.Items
	}
	// Sort by UseCount descending
	sorted := make([]Favorite, len(f.Items))
	copy(sorted, f.Items)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].UseCount > sorted[i].UseCount {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted[:n]
}

func (f *Favorites) Clear() error {
	f.Items = []Favorite{}
	return f.Save()
}

func (f *Favorites) Save() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	favoritesPath := filepath.Join(configDir, "favorites.json")
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(favoritesPath, data, 0644)
}

func LoadFavorites() (*Favorites, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return &Favorites{Items: []Favorite{}}, nil
	}

	favoritesPath := filepath.Join(configDir, "favorites.json")
	data, err := os.ReadFile(favoritesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Favorites{Items: []Favorite{}}, nil
		}
		return nil, err
	}

	var favorites Favorites
	if err := json.Unmarshal(data, &favorites); err != nil {
		return &Favorites{Items: []Favorite{}}, nil
	}

	return &favorites, nil
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".repo-lyzer")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

type FavoritesModel struct {
	favorites *Favorites
	cursor    int
}

func NewFavoritesModel() FavoritesModel {
	return FavoritesModel{}
}

func (m FavoritesModel) Update(msg tea.Msg) (FavoritesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.favorites != nil && m.cursor < len(m.favorites.Items)-1 {
				m.cursor++
			}
		case "enter":
			// Analyze selected favorite
			if m.favorites != nil && len(m.favorites.Items) > 0 {
				repoName := m.favorites.Items[m.cursor].RepoName
				m.favorites.UpdateUsage(repoName)
				if err := m.favorites.Save(); err != nil {
					log.Printf("Failed to save favorites: %v", err)
					return m, func() tea.Msg { return StatusMsg{Message: "Failed to save favorites: " + err.Error(), IsError: true} }
				}
				return m, func() tea.Msg { return AnalyzeRepoMsg{repoName: repoName} }
			}
		case "d":
			// Remove from favorites
			if m.favorites != nil && len(m.favorites.Items) > 0 {
				m.favorites.Remove(m.favorites.Items[m.cursor].RepoName)
				if err := m.favorites.Save(); err != nil {
					log.Printf("Failed to save favorites: %v", err)
					return m, func() tea.Msg { return StatusMsg{Message: "Failed to save favorites: " + err.Error(), IsError: true} }
				}
				if m.cursor >= len(m.favorites.Items) && m.cursor > 0 {
					m.cursor--
				}
			}
		case "a":
			// Add new favorite (go to input)
			return m, func() tea.Msg { return SwitchToInputMsg{} }
		case "q", "esc":
			return m, func() tea.Msg { return BackToMenuMsg{} }
		}
	}
	return m, nil
}

func (m FavoritesModel) View(width, height int) string {
	header := TitleStyle.Render("⭐ Favorite Repositories")

	if m.favorites == nil || len(m.favorites.Items) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render("No favorites yet!\n\nAnalyze a repository and press 'b' to bookmark it."),
			SubtleStyle.Render("a: add new • q/ESC: back to menu"),
		)

		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	// Build favorites list
	var lines []string
	lines = append(lines, fmt.Sprintf("%-35s │ %-10s │ %s", "Repository", "Uses", "Last Used"))
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("─").Repeat(65))

	for i, fav := range m.favorites.Items {
		prefix := "  "
		if i == m.cursor {
			prefix = "▶ "
		}
		line := fmt.Sprintf("%s%-33s │ %-10d │ %s",
			prefix,
			fav.RepoName,
			fav.UseCount,
			fav.LastUsed.Format("2006-01-02"),
		)
		if i == m.cursor {
			lines = append(lines, SelectedStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	tableBox := BoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	footer := SubtleStyle.Render("↑↓: navigate • Enter: analyze • d: remove • a: add new • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
		footer,
	)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m *FavoritesModel) SetFavorites(favorites *Favorites) {
	m.favorites = favorites
}

func (m *FavoritesModel) GetCursor() int {
	return m.cursor
}

func (m *FavoritesModel) SetCursor(cursor int) {
	m.cursor = cursor
}
