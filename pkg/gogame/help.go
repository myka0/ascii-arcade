package gogame

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"

	"charm.land/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		` ██████╗  ██████╗
██╔════╝ ██╔═══██╗
██║  ███╗██║   ██║
██║   ██║██║   ██║
╚██████╔╝╚██████╔╝
 ╚═════╝  ╚═════╝`,
	)

	Intro = `Go is an ancient abstract strategy board game where two
players take turns placing stones on a grid.

• Place stones to surround territory and capture opponent's stones.
• A group of stones is captured when it has no liberties
  (adjacent empty intersections).
• The player with the most area (territory surrounded by their
  stones + their remaining stones on the board) wins.
• A stone cannot be placed where it would be captured
  immediately (suicide), unless it captures opponent stones.`

	ScoringInfo = `Scoring:
• Area scoring is used: Territory + Stones on board.
• White receives 7.5 komi (compensation for going second).
• The game ends when both players pass consecutively.
• After the game ends, click stones to mark them as dead,
  then press Score to remove them and calculate the final score.
• Dead stones no longer count toward either score.`
)

// Help returns the Go help screen UI.
func (m *GoModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)

	// Combine all help menu sections vertically
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
		components.Section("Scoring", ScoringInfo),
	)

	shortcutKeybinds := []components.Keybind{
		{Key: "space/enter", Action: "place stone"},
		{Key: "click", Action: "place stone"},
		{Key: "ctrl+v", Action: "board size"},
		{Key: "ctrl+r", Action: "reset"},
	}

	gameKeybinds := []components.Keybind{
		{Key: "↑/w", Action: "up"},
		{Key: "↓/s", Action: "down"},
		{Key: "←/a", Action: "left"},
		{Key: "→/d", Action: "right"},
		{Key: "r", Action: "resign"},
		{Key: "l", Action: "labels"},
		{Key: "p", Action: "pass"},
	}

	game := components.ViewKeybinds("Game Keybinds", gameKeybinds)
	shortcuts := components.ViewWideKeybinds("Shortcuts", shortcutKeybinds)

	keybinds := lipgloss.JoinVertical(
		lipgloss.Center,
		components.JoinKeybinds(game, shortcuts),
		components.GlobalKeybinds(),
	)

	return components.CreateHelpMenu(Header, menu, keybinds)
}
