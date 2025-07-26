package checkers

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		`██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗███████╗██████╗ ███████╗
██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝██╔════╝██╔══██╗██╔════╝
██║     ███████║█████╗  ██║     █████╔╝ █████╗  ██████╔╝███████╗
██║     ██╔══██║██╔══╝  ██║     ██╔═██╗ ██╔══╝  ██╔══██╗╚════██║
╚██████╗██║  ██║███████╗╚██████╗██║  ██╗███████╗██║  ██║███████║
 ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝`,
	)

	Intro = `Capture all your opponent’s pieces or block them
from making a move in this classic strategy game.

• Pieces move diagonally forward one space at a time.
• Capture by jumping over an opponent’s piece.
• Multiple jumps are allowed in a single turn.
• Reach the far side of the board to crown a King,
  which can move both forward and backward.

The game ends when one player has no pieces left or
no legal moves remaining.`

	Rendering = `There are 3 different styles of rendering:
• Block - Pieces drawn using block characters
• ASCII - Pieces are created using ASCII art
• Nerdfont – Uses Nerdfont chess icons`

	HeaderStyle = lipgloss.NewStyle().
			Foreground(colors.Dark1).
			Background(colors.Purple).
			Margin(1, 0).
			Padding(0, 1).
			Bold(true)

	PieceBox = lipgloss.NewStyle().
			Margin(0, 4, 1, 4)

	PieceStyle = lipgloss.NewStyle().Foreground(colors.Purple)
)

// Help returns the Checkers help screen UI
func (m CheckersModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)

	// Combine all help menu sections vertically
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
		createRendererExamples(),
	)

	// Define keybindings specific to the game
	keybinds := []components.Keybind{
		{Key: "ctrl+r", Action: "reset"},
		{Key: "ctrl+v", Action: "renderer"},
		{Key: "click", Action: "select"},
	}

	return components.CreateHelpMenu(Header, menu, components.GameKeybinds(keybinds))
}

// createRendererExamples returns the renderer examples section for the Chess help screen.
func createRendererExamples() string {
	block := PieceBox.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			HeaderStyle.Render("Block"),
			blockKing(),
		),
	)

	ascii := PieceBox.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			HeaderStyle.Render("ASCII"),
			asciiKing(),
		),
	)

	nerdfont := PieceBox.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			HeaderStyle.Render("Nerdfont"),
			nerdfontKing(),
		),
	)

	pieces := lipgloss.JoinHorizontal(
		lipgloss.Top,
		block,
		ascii,
		nerdfont,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		Rendering,
		pieces,
		"The stlye can be changed by pressing ctrl+v.",
	)

	return components.Section("Rendering Styles", content)
}

// blockKing returns the block pawn example.
func blockKing() string {
	return strings.Join([]string{
		PieceStyle.Render(`  ▄▄▄  `),
		PieceStyle.Render(`▄▀█▀█▀▄`),
		PieceStyle.Render(`▀█▄▄▄█▀`),
		PieceStyle.Render(`  ▀▀▀  `),
	}, "\n")
}

// asciiKing returns the ASCII pawn example.
func asciiKing() string {
	return strings.Join([]string{
		PieceStyle.Render(` ,gPPRg, `),
		PieceStyle.Render(`dP' K 'Yb`),
		PieceStyle.Render(`Yb  K  dP`),
		PieceStyle.Render(` "8ggg8" `),
	}, "\n")
}

// nerdfontKing returns the Nerdfont pawn example.
func nerdfontKing() string {
	return PieceStyle.Margin(1).Render("󱟜 ")
}
