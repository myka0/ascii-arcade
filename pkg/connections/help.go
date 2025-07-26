package connections

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Header = lipgloss.NewStyle().Foreground(colors.Purple).Render(
		` ██████╗ ██████╗ ███╗   ██╗███╗   ██╗███████╗ ██████╗████████╗██╗ ██████╗ ███╗   ██╗███████╗
██╔════╝██╔═══██╗████╗  ██║████╗  ██║██╔════╝██╔════╝╚══██╔══╝██║██╔═══██╗████╗  ██║██╔════╝
██║     ██║   ██║██╔██╗ ██║██╔██╗ ██║█████╗  ██║        ██║   ██║██║   ██║██╔██╗ ██║███████╗
██║     ██║   ██║██║╚██╗██║██║╚██╗██║██╔══╝  ██║        ██║   ██║██║   ██║██║╚██╗██║╚════██║
╚██████╗╚██████╔╝██║ ╚████║██║ ╚████║███████╗╚██████╗   ██║   ██║╚██████╔╝██║ ╚████║███████║
 ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝ ╚═════╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝`,
	)

	Intro = `Find groups of four items that share something in common.

• Select four items and tap 'Submit' to check if your
  guess is correct.
• Find the groups without making 4 mistakes!

• The game fetches the latest Connections puzzle from NYT.
• Your progress is saved automatically.`

	Examples = `• FISH: Bass, Flounder, Salmon, Trout
• FIRE ___: Ant, Drill, Island, Opal

Categories will always be more specific than
"5-LETTER-WORDS," "NAMES" or "VERBS."

Each puzzle has exactly one solution.
Watch out for words that seem to belong to multiple categories!

Each group is assigned a color, which will be revealed as you solve:
`
)

// Help returns the Connections help screen UI
func (m ConnectionsModel) Help() string {
	howToPlay := components.Section("How To Play", Intro)

	// Combine all help menu sections vertically
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		howToPlay,
		createExamples(),
	)

	// Define keybindings specific to the game
	keybinds := []components.Keybind{
		{Key: "ctrl+s", Action: "shuffle"},
		{Key: "crtl+d", Action: "deselect"},
		{Key: "ctrl+r", Action: "reset"},
		{Key: "enter", Action: "submit"},
		{Key: "click", Action: "select"},
	}

	return components.CreateHelpMenu(Header, menu, components.GameKeybinds(keybinds))
}

// createExamples returns the examples section for the Connections help screen.
func createExamples() string {
	exampleStyle := lipgloss.NewStyle().MarginLeft(2)

	examples := lipgloss.JoinVertical(
		lipgloss.Left,
		Examples,
		exampleStyle.Background(Color1).Render("  ")+FGLightText.Render("  Straightforward"),
		exampleStyle.Background(Color2).Render("  ")+FGLightText.Render("    │"),
		exampleStyle.Background(Color3).Render("  ")+FGLightText.Render("    ▼"),
		exampleStyle.Background(Color4).Render("  ")+FGLightText.Render("  Tricky"),
	)

	return components.Section("Category Examples", examples)
}
