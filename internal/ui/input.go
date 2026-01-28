package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputModel struct {
	input string
	err   error
}

func NewInputModel() InputModel {
	return InputModel{}
}

func (m InputModel) Init() tea.Cmd {
	return nil
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cleanInput := sanitizeRepoInput(m.input)
			if cleanInput != "" {
				m.input = cleanInput
				m.err = nil
				return m, func() tea.Msg { return AnalyzeRepoMsg{repoName: cleanInput} }
			} else {
				m.err = fmt.Errorf("please enter a valid repository (owner/repo or GitHub URL)")
			}
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyRunes:
			m.input += string(msg.Runes)
		case tea.KeyEsc:
			return m, func() tea.Msg { return BackToMenuMsg{} }
		case tea.KeyCtrlU:
			m.input = ""
		case tea.KeyCtrlW:
			m.input = strings.TrimRight(m.input, " ")
			if idx := strings.LastIndex(m.input, " "); idx >= 0 {
				m.input = m.input[:idx+1]
			} else {
				m.input = ""
			}
		}
	}
	return m, nil
}

func (m InputModel) View() string {
	inputContent :=
		TitleStyle.Render("📥 ENTER REPOSITORY") + "\n\n" +
			InputStyle.Render("> "+m.input) + "\n\n" +
			SubtleStyle.Render("Format: owner/repo or GitHub URL  •  Press Enter to analyze")

	if m.err != nil {
		inputContent += "\n\n" + ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	box := BoxStyle.Render(inputContent)

	return box
}

func (m *InputModel) SetInput(input string) {
	m.input = input
}

func (m *InputModel) GetInput() string {
	return m.input
}

func (m *InputModel) ClearError() {
	m.err = nil
}

func (m *InputModel) SetError(err error) {
	m.err = err
}

func sanitizeRepoInput(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	// Handle GitHub URLs
	if strings.HasPrefix(input, "https://github.com/") {
		parts := strings.Split(input, "/")
		if len(parts) >= 5 {
			return parts[3] + "/" + parts[4]
		}
	}
	// Handle owner/repo format
	if strings.Contains(input, "/") && len(strings.Split(input, "/")) == 2 {
		return input
	}
	return ""
}
