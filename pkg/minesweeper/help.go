package minesweeper

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"

	"charm.land/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		`███╗   ███╗██╗███╗   ██╗███████╗███████╗██╗    ██╗███████╗███████╗██████╗ ███████╗██████╗
████╗ ████║██║████╗  ██║██╔════╝██╔════╝██║    ██║██╔════╝██╔════╝██╔══██╗██╔════╝██╔══██╗
██╔████╔██║██║██╔██╗ ██║█████╗  ███████╗██║ █╗ ██║█████╗  █████╗  ██████╔╝█████╗  ██████╔╝
██║╚██╔╝██║██║██║╚██╗██║██╔══╝  ╚════██║██║███╗██║██╔══╝  ██╔══╝  ██╔═══╝ ██╔══╝  ██╔══██╗
██║ ╚═╝ ██║██║██║ ╚████║███████╗███████║╚███╔███╔╝███████╗███████╗██║     ███████╗██║  ██║
╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝     ╚══════╝╚═╝  ╚═╝`,
	)

	Intro = `Reveal all safe cells on the board without detonating a mine.

• Numbers indicate how many adjacent cells contain mines.
• Flag cells you suspect contain mines to keep track.
• Use chord to reveal all unflagged neighbors around a number
  when the correct number of flags are placed around it.
• The first cell you reveal is always safe.`
)

// Help returns the Minesweeper help screen UI.
func (m *MinesweeperModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)

	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
	)

	movementKeybinds := []components.Keybind{
		{Key: "↑ / w", Action: "up"},
		{Key: "↓ / s", Action: "down"},
		{Key: "← / a", Action: "left"},
		{Key: "→ / d", Action: "right"},
	}

	keyboardKeybinds := []components.Keybind{
		{Key: "space / enter", Action: "reveal cell"},
		{Key: "f", Action: "toggle flag"},
		{Key: "c", Action: "chord"},
		{Key: "1 / 2 / 3", Action: "difficulty"},
	}

	mouseKeybinds := []components.Keybind{
		{Key: "click", Action: "reveal cell"},
		{Key: "r-click", Action: "toggle flag"},
		{Key: "m-click", Action: "chord"},
	}

	keyboard := components.ViewWideKeybinds("Keyboard", keyboardKeybinds)
	mouse := components.ViewWideKeybinds("Mouse", mouseKeybinds)
	movement := components.ViewKeybinds("Movement", movementKeybinds)

	keybinds := lipgloss.JoinVertical(
		lipgloss.Center,
		components.JoinKeybinds(keyboard, mouse),
		components.JoinKeybinds(movement, components.GlobalKeybinds()),
	)

	return components.CreateHelpMenu(Header, menu, keybinds)
}