package ui

import (
	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	stateMenu sessionState = iota
	stateInput
	stateLoading
	stateDashboard
	stateTree
	stateFileEdit
	stateSettings
	stateHelp
	stateHistory
	stateFavorites
	stateCompareInput
	stateCompareLoading
	stateCompareResult
	stateCloneInput
	stateCloning
)

// Message types for sub-models
type AnalyzeRepoMsg struct {
	repoName string
}

type BackToMenuMsg struct{}

type SwitchToInputMsg struct{}

type ErrorMsg error

type MainModel struct {
	state sessionState

	// Sub-models for different UI states
	menu           MenuModel
	input          InputModel
	loading        LoadingModel
	compareInput   CompareInputModel
	compareLoading CompareLoadingModel
	compareResult  CompareResultModel
	settings       SettingsModel
	help           HelpModel
	history        HistoryModel
	favorites      FavoritesModel
	cloneInput     CloneInputModel
	cloning        CloningModel

	// Shared models
	dashboard DashboardModel
	tree      TreeModel
	fileEdit  FileEditModel

	// Shared state
	windowWidth  int
	windowHeight int
	cache        *cache.Cache
	appConfig    *config.AppSettings

	// Additional state fields
	historyCursor   int
	animTick        int
	compareStep     int
	compareInput1   string
	compareInput2   string
	inTokenInput    bool
	tokenInput      string
	settingsOption  string
	favoritesCursor int
	analysisType    string
	helpContent     string
	err             error
}

// NewMainModel creates a new MainModel with initialized sub-models
func NewMainModel(cache *cache.Cache, config *config.AppSettings) MainModel {
	return MainModel{
		state:     stateMenu,
		menu:      NewMenuModel(),
		input:     NewInputModel(),
		cache:     cache,
		appConfig: config,
	}
}

// Init initializes the Bubble Tea program
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// Update sub-models with window size
		m.menu.width = msg.Width
		m.menu.height = msg.Height

	case BackToMenuMsg:
		m.state = stateMenu
		return m, nil

	case SwitchToInputMsg:
		m.state = stateInput
		return m, nil

	case AnalyzeRepoMsg:
		m.state = stateLoading
		// TODO: Start analysis command
		return m, nil

	case ErrorMsg:
		m.err = error(msg)
		return m, nil
	}

	// Delegate to current state's sub-model
	switch m.state {
	case stateMenu:
		menuModel, cmd := m.menu.Update(msg)
		m.menu = menuModel.(MenuModel)
		if m.menu.Done {
			switch m.menu.SelectedOption {
			case 0: // Analyze Repository
				m.state = stateInput
			case 1: // Favorites
				m.state = stateFavorites
			case 2: // Compare Repositories
				m.state = stateCompareInput
			case 3: // View History
				m.state = stateHistory
			case 4: // Clone Repository
				m.state = stateCloneInput
			case 5: // Settings
				m.state = stateSettings
			case 6: // Help
				m.state = stateHelp
			case 7: // Exit
				return m, tea.Quit
			}
			m.menu.Done = false
		}
		return m, cmd

	case stateInput:
		inputModel, cmd := m.input.Update(msg)
		m.input = inputModel
		return m, cmd

	default:
		return m, nil
	}
}

// View renders the current UI state
func (m MainModel) View() string {
	switch m.state {
	case stateMenu:
		return m.menu.View()
	case stateInput:
		return m.input.View()
	default:
		return "State not implemented"
	}
}

// Run starts the Bubble Tea program
func Run(cache *cache.Cache, config *config.AppSettings) error {
	model := NewMainModel(cache, config)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
