package solitaire

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		`███████╗ ██████╗ ██╗     ██╗████████╗ █████╗ ██╗██████╗ ███████╗
██╔════╝██╔═══██╗██║     ██║╚══██╔══╝██╔══██╗██║██╔══██╗██╔════╝
███████╗██║   ██║██║     ██║   ██║   ███████║██║██████╔╝█████╗
╚════██║██║   ██║██║     ██║   ██║   ██╔══██║██║██╔══██╗██╔══╝
███████║╚██████╔╝███████╗██║   ██║   ██║  ██║██║██║  ██║███████╗
╚══════╝ ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚══════╝`,
	)

	Intro = `Sort all cards into four foundation piles by suit,
from Ace to King.

• Cards can be moved between columns in descending order,
  alternating colors.
• Empty columns can be filled with Kings or valid sequences.
• Draw from the stock when no more moves are available.

You can play using either the mouse or keyboard shortcuts.
Both methods will perform the first valid move available
for the selected deck.`

	Shortcuts = `• Press space to draw from the stock.
• Press w to play the top card from the waste pile.
• Press 1–7 to play from tableau columns.
• Press shift + 1-4 to play from foundations (♠, ♣, ♥, ♦).
• Press u to undo your last move.

Clear all the cards to win!`
)

// Help returns the Solitaire help screen UI
func (m SolitaireModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)
	shortcuts := components.Section("Shortcuts", Shortcuts)

	// Combine all help menu sections vertically
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
		shortcuts,
	)

	// Define keybindings specific to the game
	keybinds := []components.Keybind{
		{Key: "click", Action: "select"},
		{Key: "r-click", Action: "undo"},
		{Key: "ctrl+r", Action: "reset"},
		{Key: "space", Action: "draw"},
		{Key: "w", Action: "waste"},
		{Key: "u", Action: "undo"},
		{Key: "1-7", Action: "tableau"},
		{Key: "!", Action: "♠ foundation"},
		{Key: "@", Action: "♣ foundation"},
		{Key: "#", Action: "♥ foundation"},
		{Key: "$", Action: "♦ foundation"},
	}

	return components.CreateHelpMenu(Header, menu, components.GameKeybinds(keybinds))
}
