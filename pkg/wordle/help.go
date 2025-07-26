package wordle

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		`██╗    ██╗ ██████╗ ██████╗ ██████╗ ██╗     ███████╗
██║    ██║██╔═══██╗██╔══██╗██╔══██╗██║     ██╔════╝
██║ █╗ ██║██║   ██║██████╔╝██║  ██║██║     █████╗
██║███╗██║██║   ██║██╔══██╗██║  ██║██║     ██╔══╝
╚███╔███╔╝╚██████╔╝██║  ██║██████╔╝███████╗███████╗
 ╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═╝╚═════╝ ╚══════╝╚══════╝`,
	)

	Intro = `Guess the Wordle in 6 tries.

• Each guess must be a valid 5-letter word.
• The color of the tiles will change to show how
  close your guess was to the word.

• The game fetches the latest Wordle puzzle from NYT.
• Your progress is saved automatically.`
)

// Help returns the Wordle help screen UI
func (m WordleModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)

	// Combine all help menu sections vertically
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
		createExamples(),
	)

	// Define keybindings specific to the game
	keybinds := []components.Keybind{
		{Key: "<char>", Action: "input"},
		{Key: "bksp", Action: "erase"},
		{Key: "enter", Action: "submit"},
		{Key: "ctrl+r", Action: "reset"},
	}

	return components.CreateHelpMenu(Header, menu, components.GameKeybinds(keybinds))
}

// createExamples returns the visual examples section for the Wordle help screen.
func createExamples() string {
	wordy := lipgloss.JoinHorizontal(
		lipgloss.Top,
		FGKeyCorrect.Render(Border.Render("W")),
		FGText.Render(Border.Render("O")),
		FGText.Render(Border.Render("R")),
		FGText.Render(Border.Render("D")),
		FGText.Render(Border.Render("Y")),
	)
	ex1 := lipgloss.JoinVertical(
		lipgloss.Left,
		wordy,
		" W is in the word and in the correct spot.\n",
	)

	light := lipgloss.JoinHorizontal(
		lipgloss.Top,
		FGText.Render(Border.Render("L")),
		FGKeyPresent.Render(Border.Render("I")),
		FGText.Render(Border.Render("G")),
		FGText.Render(Border.Render("H")),
		FGText.Render(Border.Render("T")),
	)
	ex2 := lipgloss.JoinVertical(
		lipgloss.Left,
		light,
		" I is in the word but in the wrong spot.\n",
	)

	rogue := lipgloss.JoinHorizontal(
		lipgloss.Top,
		FGText.Render(Border.Render("R")),
		FGText.Render(Border.Render("O")),
		FGText.Render(Border.Render("U")),
		FGKeyAbsent.Render(Border.Render("G")),
		FGText.Render(Border.Render("E")),
	)
	ex3 := lipgloss.JoinVertical(
		lipgloss.Left,
		rogue,
		" U is not in the word in any spot.",
	)

	return components.Section(
		"Examples",
		lipgloss.JoinVertical(lipgloss.Left, ex1, ex2, ex3),
	)
}
