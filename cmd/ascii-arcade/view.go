package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

var implementedGames = map[string]bool{
	"Tetris":       true,
	"Snake":        false,
	"Solitaire":    true,
	"Minesweeper":  false,
	"Crossword":    true,
	"Mini":         true,
	"Wordle":       true,
	"Connections":  true,
	"Sudoku":       false,
	"Go":           true,
	"Chess":        true,
	"Checkers":     true,
	"Connect Four": false,
}

var keyBindingsView string

func init() {
	keyBindings := []string{
		KeyStyle.Render("↑ ") + KeyActionStyle.Render("up"),
		KeyStyle.Render("↓ ") + KeyActionStyle.Render("down"),
		KeyStyle.Render("<char> ") + KeyActionStyle.Render("filter"),
		KeyStyle.Render("bksp ") + KeyActionStyle.Render("erase"),
		KeyStyle.Render("esc ") + KeyActionStyle.Render("clear search"),
		KeyStyle.Render("enter ") + KeyActionStyle.Render("select"),
		KeyStyle.Render("ctrl+h ") + KeyActionStyle.Render("main menu"),
		KeyStyle.Render("ctrl+c ") + KeyActionStyle.Render("quit"),
	}

	left := strings.Join(keyBindings[:4], "\n")
	right := strings.Join(keyBindings[4:], "\n")

	keyBindingsView = lipgloss.JoinHorizontal(
		lipgloss.Top,
		KeyBindMenu.Render(left),
		KeyBindMenu.Render(right),
	)
}

// View renders the full UI centered in the terminal.
func (m model) View() tea.View {
	var view string
	switch {
	case m.isHelpSelected:
		// If a help page is selected, render the help page
		view = m.activeModel.(ViewModel).Help()
	case m.isGameSelected:
		// If a game is selected, render the game UI
		view = m.activeModel.(ViewModel).View().Content
	default:
		// Render the home menu
		keyBindings := lipgloss.NewStyle().MarginLeft(8).Render(keyBindingsView)
		message := CenteredText.Render(m.message) + "\n"

		view = MenuStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				Header,
				m.viewGameList(),
				message,
				keyBindings,
			),
		)
	}

	v := tea.NewView(
		zone.Scan(
			lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center, view),
		),
	)

	v.AltScreen = true
	if !m.noMouse {
		v.MouseMode = tea.MouseModeCellMotion
	}

	return v
}

// viewGameList returns the rendered list of game titles.
func (m model) viewGameList() string {
	var list string

	// Choose the appropriate view based on whether the user has entered a search query.
	if m.searchQuery == "" {
		list = m.viewDefaultList()
	} else {
		list = m.viewFilteredList()
	}

	return ListStyle.Render(list)
}

// viewDefaultList constructs the full, grouped list of all games.
func (m model) viewDefaultList() string {
	var b strings.Builder
	b.Grow(2048)

	for _, gamesList := range Games {
		// Add the header for the current list.
		b.WriteString(ListHeader.Render(gamesList.Header))
		b.WriteByte('\n')

		// Add games from the current list and highlight the selected game.
		for _, game := range gamesList.Games {
			switch {
			case m.games[m.selectedGameIdx] == game:
				b.WriteString(SelectedListEntry.Render("> " + game))
			case !implementedGames[game]:
				b.WriteString(UnimplementedListEntry.Render(game))
			default:
				b.WriteString(ListEntry.Render(game))
			}
			b.WriteByte('\n')
		}
	}

	return b.String()
}

// viewFilteredList constructs the flat list of games that match the current search query.
func (m model) viewFilteredList() string {
	var b strings.Builder
	b.Grow(512)
	b.WriteString(ListHeader.Render("Query: " + m.searchQuery))
	b.WriteByte('\n')

	if len(m.games) == 0 {
		b.WriteString(Text.Render("No results found"))
		b.WriteByte('\n')
		return b.String()
	}

	// Add games from the current list and highlight the selected game.
	for i, game := range m.games {
		switch {
		case m.selectedGameIdx == i:
			b.WriteString(SelectedListEntry.Render("> " + game))
		case !implementedGames[game]:
			b.WriteString(UnimplementedListEntry.Render(game))
		default:
			b.WriteString(ListEntry.Render(game))
		}
		b.WriteByte('\n')
	}

	return b.String()
}