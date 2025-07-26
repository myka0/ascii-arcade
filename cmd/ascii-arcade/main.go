package main

import (
	"ascii-arcade/pkg/checkers"
	"ascii-arcade/pkg/chess"
	"ascii-arcade/pkg/connections"
	"ascii-arcade/pkg/crossword"
	"ascii-arcade/pkg/solitaire"
	"ascii-arcade/pkg/wordle"

	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// GamesList holds the list of games.
type GamesList struct {
	Header string
	Games  []string
}

var Games = []GamesList{
	{
		Header: "Classics",
		Games: []string{
			"Tetris",
			"Snake",
			"Solitaire",
			"Minesweeper",
		},
	},
	{
		Header: "New York Times",
		Games: []string{
			"Crossword",
			"Wordle",
			"Connections",
			"Sudoku",
		},
	},
	{
		Header: "Strategy Games",
		Games: []string{
			"Go",
			"Chess",
			"Checkers",
			"Connect Four",
		},
	},
}

// Saver defines a game model that can persist state.
type Saver interface {
	SaveToFile() error
}

// ViewModel defines a view model that can render itself.
type ViewModel interface {
	View() string
	Help() string
}

// model holds global app state.
type model struct {
	windowHeight    int
	windowWidth     int
	isGameSelected  bool
	isHelpSelected  bool
	activeModel     tea.Model
	selectedGame    string
	selectedGameIdx int
	games           []string
	searchQuery     string
	message         string
}

// Creates the initial model with connections as default.
func initialModel(startGame string) *model {
	m := model{}
	m.games = handleSearch("")

	// If a start game is specified, initialize it
	if startGame != "" {
		m.selectedGame = strings.ToUpper(startGame[0:1]) + strings.ToLower(startGame[1:])
		m.handleSwitchModel()
	}

	return &m
}

// Init implements the Bubble Tea interface for initialization.
func (m model) Init() tea.Cmd {
	if m.isGameSelected && m.activeModel != nil {
		return m.activeModel.Init()
	}
	return nil
}

// Update handles keypress events and updates the model state accordingly.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global key bindings
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.handleSaveGame()
			return m, tea.Quit

		case "ctrl+h":
			m.isGameSelected = false
			m.isHelpSelected = false
			return m, nil

		case "?":
			if m.isGameSelected {
				m.isHelpSelected = !m.isHelpSelected
			}
			return m, nil
		}

	// If the window is resized, store its new dimensions
	case tea.WindowSizeMsg:
		return m, m.handleResize(msg)
	}

	// If the help is selected use the help key bindings
	if m.isHelpSelected {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.isHelpSelected = false
				return m, nil
			}

		case tea.MouseMsg:
			switch msg := msg.(type) {
			case tea.MouseClickMsg:
				if msg.Mouse().Button != tea.MouseLeft {
					return m, nil
				}

				if zone.Get("continue").InBounds(msg) {
					m.isHelpSelected = false
					return m, nil
				}
			}
		}

		return m, nil
	}

	// If the game is selected, pass the keypress to the active model
	if m.isGameSelected {
		var cmd tea.Cmd
		m.activeModel, cmd = m.activeModel.Update(msg)

		// If the active game wants to go home, exit the game
		if cmd != nil && cmd() == "home" {
			m.isGameSelected = false
			return m, nil
		}

		return m, cmd
	}

	return m.handleHomeMenuInput(msg)
}

// handleHomeMenuInput handles keypress events for the home menu.
func (m *model) handleHomeMenuInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.message = ""
			m.selectedGameIdx = (m.selectedGameIdx + 1) % len(m.games)
			m.selectedGame = m.games[m.selectedGameIdx]

		case "up":
			m.message = ""
			m.selectedGameIdx = (m.selectedGameIdx - 1 + len(m.games)) % len(m.games)
			m.selectedGame = m.games[m.selectedGameIdx]

		case "enter":
			m.selectedGame = m.games[m.selectedGameIdx]
			m.handleSwitchModel()
			if m.activeModel == nil {
				m.message = "Selected game not implemented yet."
				return m, nil
			}
			return m.handleInitGame()

		case "backspace":
			m.message = ""
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.games = handleSearch(m.searchQuery)
			}

		case "esc":
			m.message = ""
			m.searchQuery = ""
			m.games = handleSearch(m.searchQuery)

		// Handle search input
		default:
			m.message = ""
			if len(msg.String()) == 1 && len(m.games) != 0 {
				m.searchQuery += msg.String()
				m.games = handleSearch(m.searchQuery)

				if len(m.games) == 0 {
					m.selectedGameIdx = 0
				} else if m.selectedGameIdx >= len(m.games) {
					m.selectedGameIdx = len(m.games) - 1
				}
			}
		}
	}

	return m, nil
}

// handleSwitchModel swaps in a new game model based on selected tab.
func (m *model) handleSwitchModel() tea.Model {
	switch m.selectedGame {
	case "Solitaire":
		m.activeModel = solitaire.InitSolitaireModel()
	case "Crossword":
		m.activeModel = crossword.InitCrosswordModel()
	case "Wordle":
		m.activeModel = wordle.InitWordleModel()
	case "Connections":
		m.activeModel = connections.InitConnectionsModel()
	case "Checkers":
		m.activeModel = checkers.InitCheckersModel()
	case "Chess":
		m.activeModel = chess.InitChessModel()
	default:
		return m
	}

	m.isGameSelected = true
	m.isHelpSelected = true
	return m
}

// handleInitGame initializes the game model.
func (m *model) handleInitGame() (tea.Model, tea.Cmd) {
	return m, m.activeModel.Init()
}

// handleSaveGame updates window size on resize.
func (m *model) handleResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.windowHeight = msg.Height
	m.windowWidth = msg.Width
	return nil
}

// handleSaveGame saves the current model if it implements Saver.
func (m *model) handleSaveGame() {
	if saver, ok := m.activeModel.(Saver); ok {
		if err := saver.SaveToFile(); err != nil {
			fmt.Println("Auto-save failed:", err)
		}
	}
}

// handleSearch filters for games that contain the search query.
func handleSearch(query string) []string {
	var matches []string

	// Return all games if query is empty
	if query == "" {
		for _, gamesList := range Games {
			matches = append(matches, gamesList.Games...)
		}
		return matches
	}

	// Match games whose names start with the query
	for _, list := range Games {
		for _, game := range list.Games {
			if strings.HasPrefix(strings.ToLower(game), strings.ToLower(query)) {
				matches = append(matches, game)
			}
		}
	}

	return matches
}

// Entry point of the application.
func main() {
	noMouse := flag.Bool("no-mouse", false, "Run without mouse support")
	startGame := flag.String("game", "", "Start with a specific game")
	flag.Parse()

	zone.NewGlobal()

	var opts []tea.ProgramOption
	if !*noMouse {
		opts = append(opts, tea.WithAltScreen())
		opts = append(opts, tea.WithMouseCellMotion())
	}

	p := tea.NewProgram(initialModel(*startGame), opts...)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
