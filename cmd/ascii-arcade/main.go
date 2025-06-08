package main

import (
	"crossword/internal/colors"
	"crossword/internal/connections"
	"crossword/internal/crossword"
	"crossword/internal/wordle"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	Text = lipgloss.NewStyle().Foreground(colors.Light1)

	ActiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	TabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	Tab = lipgloss.NewStyle().
		Border(TabBorder, true).
		BorderForeground(colors.Blue).
		Padding(0, 1)

	ActiveTab = Tab.Border(ActiveTabBorder, true)

	TabGap = Tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

// model holds global app state.
type model struct {
	game         int
	windowHeight int
	windowWidth  int
	activeModel  tea.Model
}

// Saver defines a game model that can persist state.
type Saver interface {
	SaveToFile() error
}

// Creates the initial model with connections as default.
func initialModel() model {
	m := model{
		game:        2,
		activeModel: connections.InitConnectionsModel(),
	}

	return m
}

// Init implements the Bubble Tea interface for initialization.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles keypress events and updates the model state accordingly.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Quit
		case "ctrl+c":
			m.handleSaveGame()
			return m, tea.Quit

		// Move to prev game
		case "esc":
			if m.game > 0 {
				m.handleSaveGame()
				m.game--
				return m.handleSwitchModel()
			}

		// Move to next game
		case "ctrl+]":
			if m.game < 2 {
				m.handleSaveGame()
				m.game++
				return m.handleSwitchModel()
			}
		}

	// If the window is resized, store its new dimensions.
	case tea.WindowSizeMsg:
		return m, m.handleResize(msg)
	}

	var cmd tea.Cmd
	m.activeModel, cmd = m.activeModel.Update(msg)
	return m, cmd
}

// handleSwitchModel swaps in a new game model based on selected tab.
func (m model) handleSwitchModel() (tea.Model, tea.Cmd) {
	switch m.game {
	case 0:
		m.activeModel = crossword.InitCrosswordModel()
	case 1:
		m.activeModel = wordle.InitWordleModel()
	case 2:
		m.activeModel = connections.InitConnectionsModel()
	}
	return m, m.activeModel.Init()
}

// handleSaveGame updates window size on resize.
func (m *model) handleResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.windowHeight = msg.Height
	m.windowWidth = msg.Width
	return nil
}

// handleSaveGame saves the current model if it implements Saver.
func (m model) handleSaveGame() {
	if saver, ok := m.activeModel.(Saver); ok {
		if err := saver.SaveToFile(); err != nil {
			fmt.Println("Auto-save failed:", err)
		}
	}
}

// View renders the full UI centered in the terminal.
func (m model) View() string {
	game := m.viewTabBar() + "\n" + m.activeModel.View()
	return lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center, game)
}

// viewTabBar renders the navigation tab bar with styling.
func (m model) viewTabBar() string {
	tabsAsString := []string{"Crossword", "Wordle", "Connections"}
	tabs := make([]string, len(tabsAsString))

	// Highlight the active tab
	for i, tabName := range tabsAsString {
		if m.game == i {
			tabs[i] = ActiveTab.Render(Text.Render(tabName))
		} else {
			tabs[i] = Tab.Render(Text.Render(tabName))
		}
	}

	renderedTabs := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := TabGap.Render(strings.Repeat(" ", max(0, 96-lipgloss.Width(renderedTabs)-2)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs, gap) + "\n"
}

// Entry point of the application.
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
