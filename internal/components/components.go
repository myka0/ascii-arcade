package components

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type Keybind struct {
	Key    string
	Action string
}

// CreateHelpMenu renders the full help menu UI.
func CreateHelpMenu(header, menu, keybinds string) string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		menu,
		keybinds,
		Continue(),
	)
}

// Section formats a titled section with a styled header and body content.
func Section(header, content string) string {
	header = Header.Render(header)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		Text.Render(content),
	)
}

// ViewKeybinds displays a list of keybindings under a title, formatted into two columns.
func ViewKeybinds(title string, keybinds []Keybind) string {
	var keybindsStr []string

	// Format each keybind
	for _, keybind := range keybinds {
		entry := KeyStyle.Render(keybind.Key) + KeyActionStyle.Render(keybind.Action)
		keybindsStr = append(keybindsStr, entry)
	}

	// Split into two balanced columns
	middle := int(math.Ceil(float64(len(keybindsStr)) / 2))
	left := strings.Join(keybindsStr[:middle], "\n")
	right := strings.Join(keybindsStr[middle:], "\n")

	// Render left and right columns side by side
	menu := lipgloss.JoinHorizontal(
		lipgloss.Top,
		KeyBindBox.Align(lipgloss.Left).Render(left),
		KeyBindBox.Align(lipgloss.Right).Render(right),
	)

	// Combine header and keybinds
	return lipgloss.JoinVertical(
		lipgloss.Center,
		Header.Render(title),
		KeyBindMenu.Render(menu),
	)
}

// JoinKeybinds joins two keybind menus side by side.
func JoinKeybinds(left, right string) string {
	// Estimate visual height of each block
	leftHeight := len(strings.Split(left, "\n")) - 4
	rightHeight := len(strings.Split(right, "\n")) - 4
	maxKeybindsHeight := max(leftHeight, rightHeight)

	// Create vertical divider with height matching the tallest column
	divider := Divider.Render(strings.Repeat("│\n", maxKeybindsHeight-1) + "│")

	// Horizontally join the keybinds and divider
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		divider,
		right,
	)
}

// GameKeybinds renders both game specific and global keybindings side by side.
func GameKeybinds(keybinds []Keybind) string {
	gameKeybindsView := ViewKeybinds("Game Keybinds", keybinds)
	globalKeybindsView := GlobalKeybinds()

	return JoinKeybinds(gameKeybindsView, globalKeybindsView)
}

// GlobalKeybinds returns a rendered section with global application shortcuts.
func GlobalKeybinds() string {
	keybinds := []Keybind{
		{Key: "ctrl+h", Action: "home"},
		{Key: "ctrl+c", Action: "quit"},
		{Key: "?", Action: "help"},
	}

	return ViewKeybinds("Global Keybinds", keybinds)
}

// Continue renders a styled button and hint to advance the screen.
func Continue() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		zone.Mark("continue", ButtonStyle.Render("Continue")),
		LightText.Render("Press enter to continue..."),
	)
}
