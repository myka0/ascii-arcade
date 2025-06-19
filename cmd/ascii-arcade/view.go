package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// View renders the full UI centered in the terminal.
func (m model) View() string {
	// If a game is selected, render the game UI
	if m.isGameSelected {
		game := m.activeModel.(ViewModel).View()
		return zone.Scan(
			lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center, game),
		)
	}

	keyBindings := lipgloss.NewStyle().MarginLeft(8).Render(m.viewKeyBindings())
	message := CenteredText.Render(m.message) + "\n"

	// Render the home menu
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		Header,
		m.viewGameList(),
		message,
		keyBindings,
	)

	return lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center,
		MenuStyle.Render(menu),
	)
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
	var list string

	for _, gamesList := range Games {
		// Add the header for the current list.
		list += ListHeader.Render(gamesList.Header) + "\n"

		// Add games from the current list and highlight the selected game.
		for _, game := range gamesList.Games {
			if m.games[m.selectedGameIdx] == game {
				list += SelectedListEntry.Render("> "+game) + "\n"
			} else {
				list += ListEntry.Render(game) + "\n"
			}
		}
	}

	return list
}

// viewFilteredList constructs the flat list of games that match the current search query.
func (m model) viewFilteredList() string {
	list := ListHeader.Render("Query: "+m.searchQuery) + "\n"

	if len(m.games) == 0 {
		list += Text.Render("No results found") + "\n"
		return list
	}

	// Add games from the current list and highlight the selected game.
	for i, game := range m.games {
		if m.selectedGameIdx == i {
			list += SelectedListEntry.Render("> "+game) + "\n"
		} else {
			list += ListEntry.Render(game) + "\n"
		}
	}

	return list
}

// viewKeyBindings renders the key bindings for the main menu.
func (m model) viewKeyBindings() string {
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

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		KeyBindMenu.Render(left),
		KeyBindMenu.Render(right),
	)
}
