package tetris

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"

	"charm.land/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		`████████╗███████╗████████╗██████╗ ██╗███████╗
╚══██╔══╝██╔════╝╚══██╔══╝██╔══██╗██║██╔════╝
   ██║   █████╗     ██║   ██████╔╝██║███████╗
   ██║   ██╔══╝     ██║   ██╔══██╗██║╚════██║
   ██║   ███████╗   ██║   ██║  ██║██║███████║
   ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝`,
	)

	Intro = `Score as many points as possible by clearing lines of blocks
without letting the stack reach the top of the playing field.

• Clear 1, 2, 3, or 4 lines at a time for increasing points.
• A four-line clear (a "Tetris") and back-to-back clears award
  bonus points.
• The level increases as you clear lines, and pieces fall faster.
• A faded "ghost" piece shows where the active tetrimino will
  land if dropped from its current column.
• You can hold one piece in reserve and swap it in later, but
  only once per piece.`
)

// Help returns the Tetris help screen UI.
func (m *TetrisModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)

	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
	)

	// Define movement keybindings
	movementKeybinds := []components.Keybind{
		{Key: "← / a", Action: "move left"},
		{Key: "→ / d", Action: "move right"},
		{Key: "↓ / s", Action: "soft drop"},
		{Key: "↑ / w / space", Action: "hard drop"},
	}

	gameKeybinds := []components.Keybind{
		{Key: "z", Action: "rotate ccw"},
		{Key: "x", Action: "rotate cw"},
		{Key: "c", Action: "hold"},
		{Key: "p / esc", Action: "pause"},
		{Key: "ctrl+r", Action: "reset"},
	}

	keybinds := lipgloss.JoinVertical(
		lipgloss.Center,
		components.ViewWideKeybinds("Movement", movementKeybinds),
		components.GameKeybinds(gameKeybinds),
	)

	return components.CreateHelpMenu(Header, menu, keybinds)
}
