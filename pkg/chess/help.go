package chess

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		` ██████╗██╗  ██╗███████╗███████╗███████╗
██╔════╝██║  ██║██╔════╝██╔════╝██╔════╝
██║     ███████║█████╗  ███████╗███████╗
██║     ██╔══██║██╔══╝  ╚════██║╚════██║
╚██████╗██║  ██║███████╗███████║███████║
 ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝`,
	)

	Intro = `Checkmate your opponent’s king in a turn based
strategy game.

• Each piece moves in its own unique way.
• Select a piece to see its valid moves highlighted.
• Move your pieces to control the board and capture
  your opponent’s.`

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

// Help returns the Chess help screen UI
func (m ChessModel) Help() string {
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
			blockPawn(),
		),
	)

	ascii := PieceBox.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			HeaderStyle.Render("ASCII"),
			asciiPawn(),
		),
	)

	nerdfont := PieceBox.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			HeaderStyle.Render("Nerdfont"),
			nerdfontPawn(),
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

// blockPawn returns the block pawn example.
func blockPawn() string {
	return strings.Join([]string{
		PieceStyle.Render(`  ▄  `),
		PieceStyle.Render(`█████`),
		PieceStyle.Render(`▄███▄`),
	}, "\n")
}

// asciiPawn returns the ASCII pawn example.
func asciiPawn() string {
	return strings.Join([]string{
		PieceStyle.Render(` ( ) `),
		PieceStyle.Render(` ) ( `),
		PieceStyle.Render(`(___)`),
	}, "\n")
}

// nerdfontPawn returns the Nerdfont pawn example.
func nerdfontPawn() string {
	return PieceStyle.Margin(1).Render(" ")
}
