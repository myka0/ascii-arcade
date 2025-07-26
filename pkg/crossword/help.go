package crossword

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		` ██████╗██████╗  ██████╗ ███████╗███████╗██╗    ██╗ ██████╗ ██████╗ ██████╗
██╔════╝██╔══██╗██╔═══██╗██╔════╝██╔════╝██║    ██║██╔═══██╗██╔══██╗██╔══██╗
██║     ██████╔╝██║   ██║███████╗███████╗██║ █╗ ██║██║   ██║██████╔╝██║  ██║
██║     ██╔══██╗██║   ██║╚════██║╚════██║██║███╗██║██║   ██║██╔══██╗██║  ██║
╚██████╗██║  ██║╚██████╔╝███████║███████║╚███╔███╔╝╚██████╔╝██║  ██║██████╔╝
 ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝ ╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═╝╚═════╝ `,
	)

	Intro = `Solve the daily Crossword by filling in all the blank
squares with words that match the clues.

• Arrow into a square to select it.
• Type to fill in your guess, and use tab or enter to
  move between clues.
• Words must fit both across and down clues.

• The game fetches the latest Crossword puzzle from NYT.
• Your progress is saved automatically.`
)

// Help returns the Crossword help screen UI
func (m CrosswordModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)

	// Combine all help menu sections vertically
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
	)

	// Define keybindings specific to the game
	gameKeybinds := []components.Keybind{
		{Key: "<char>", Action: "input"},
		{Key: "bksp", Action: "erase"},
		{Key: "ctrl+r", Action: "reset"},
	}

	// Define movement keybindings
	movementKeybinds := []components.Keybind{
		{Key: "↑", Action: "up"},
		{Key: "↓", Action: "down"},
		{Key: "←", Action: "left"},
		{Key: "→", Action: "right"},
		{Key: "tab", Action: "next"},
		{Key: "enter", Action: "next"},
		{Key: "s+tab", Action: "prev"},
		{Key: "s+enter", Action: "prev"},
	}

	// Define check keybindings
	checkKeybinds := []components.Keybind{
		{Key: "ctrl+l", Action: "letter"},
		{Key: "ctrl+w", Action: "word"},
		{Key: "ctrl+p", Action: "puzzle"},
		{Key: "ctrl+a", Action: "auto"},
	}

	movement := components.ViewKeybinds("Movement", movementKeybinds)
	check := components.ViewKeybinds("Check", checkKeybinds)

	keybinds := lipgloss.JoinVertical(
		lipgloss.Left,
		components.JoinKeybinds(movement, check),
		components.GameKeybinds(gameKeybinds),
	)

	return components.CreateHelpMenu(Header, menu, keybinds)
}
