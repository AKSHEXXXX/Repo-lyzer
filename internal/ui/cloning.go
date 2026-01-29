package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CloningModel struct {
	spinner  spinner.Model
	repoName string
}

func NewCloningModel() CloningModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return CloningModel{
		spinner: s,
	}
}

func (m CloningModel) Update(msg tea.Msg) (CloningModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case cloneResult:
		return m, nil // Handled by parent
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m CloningModel) View(width, height int) string {
	header := TitleStyle.Render("📥 CLONING REPOSITORY")

	content := fmt.Sprintf(
		"%s Cloning %s to Desktop...\n\n"+
			"Please wait while the repository is being cloned.",
		m.spinner.View(),
		m.repoName,
	)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render(content),
		),
	)
}

func (m *CloningModel) SetRepoName(repoName string) {
	m.repoName = repoName
}
