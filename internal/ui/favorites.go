package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type FavoriteItem struct {
	RepoName string    `json:"repo_name"`
	UseCount int       `json:"use_count"`
	LastUsed time.Time `json:"last_used"`
	AddedAt  time.Time `json:"added_at"`
	Notes    string    `json:"notes"`
}

type FavoritesModel struct {
	Items  []FavoriteItem `json:"items"`
	width  int
	height int
}

func NewFavoritesModel() *FavoritesModel {
	return &FavoritesModel{
		Items: []FavoriteItem{},
	}
}

func (f *FavoritesModel) Add(repoName string) {
	// Check if already exists
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items[i].UseCount++
			f.Items[i].LastUsed = time.Now()
			return
		}
	}

	// Add new item
	f.Items = append(f.Items, FavoriteItem{
		RepoName: repoName,
		UseCount: 1,
		LastUsed: time.Now(),
		AddedAt:  time.Now(),
	})
}

func (f *FavoritesModel) Remove(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items = append(f.Items[:i], f.Items[i+1:]...)
			return
		}
	}
}

func (f *FavoritesModel) UpdateUsage(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items[i].UseCount++
			f.Items[i].LastUsed = time.Now()
			return
		}
	}
}

func (f *FavoritesModel) IsFavorite(repoName string) bool {
	for _, item := range f.Items {
		if item.RepoName == repoName {
			return true
		}
	}
	return false
}

func (f *FavoritesModel) GetTopFavorites(limit int) []FavoriteItem {
	if limit <= 0 {
		return []FavoriteItem{}
	}
	if limit >= len(f.Items) {
		return f.Items
	}
	return f.Items[:limit]
}

func (f *FavoritesModel) Clear() {
	f.Items = []FavoriteItem{}
}

func (f *FavoritesModel) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".repo-lyzer")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(configDir, "favorites.json")
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func LoadFavorites() (*FavoritesModel, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return NewFavoritesModel(), err
	}

	filePath := filepath.Join(home, ".repo-lyzer", "favorites.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFavoritesModel(), nil
		}
		return NewFavoritesModel(), err
	}

	var favorites FavoritesModel
	if err := json.Unmarshal(data, &favorites); err != nil {
		return NewFavoritesModel(), err
	}

	// Sort by last used (most recent first)
	sort.Slice(favorites.Items, func(i, j int) bool {
		return favorites.Items[i].LastUsed.After(favorites.Items[j].LastUsed)
	})

	return &favorites, nil
}

func (f *FavoritesModel) View() string {
	header := TitleStyle.Render("⭐ Favorite Repositories")

	if len(f.Items) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render("No favorites yet!\n\nAnalyze a repository and press 'b' to bookmark it."),
			SubtleStyle.Render("a: add new • q/ESC: back to menu"),
		)

		if f.width == 0 {
			return content
		}

		return lipgloss.Place(
			f.width,
			f.height,
			lipgloss.Center,
			lipgloss.Center,
			content,
		)
	}

	// Build favorites list
	var lines []string
	lines = append(lines, fmt.Sprintf("%-35s │ %-10s │ %s", "Repository", "Uses", "Last Used"))
	lines = append(lines, strings.Repeat("─", 65))

	for i, fav := range f.Items {
		prefix := "  "
		if i == 0 { // Assuming cursor is 0 for now, since we don't have cursor in this model
			prefix = "▶ "
		}
		line := fmt.Sprintf("%s%-33s │ %-10d │ %s",
			prefix,
			fav.RepoName,
			fav.UseCount,
			fav.LastUsed.Format("2006-01-02"),
		)
		if i == 0 {
			lines = append(lines, SelectedStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	tableBox := BoxStyle.Render(strings.Join(lines, "\n"))
	footer := SubtleStyle.Render("↑↓: navigate • Enter: analyze • d: remove • a: add new • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
		footer,
	)

	if f.width == 0 {
		return content
	}

	return lipgloss.Place(
		f.width,
		f.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
