package components

import (
	"ascii-arcade/internal/colors"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Text      = lipgloss.NewStyle().Foreground(colors.Light2)
	LightText = lipgloss.NewStyle().Foreground(colors.Medium2)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(colors.Dark1).
			Background(colors.Purple).
			Margin(2, 0, 1, 0).
			Padding(0, 1).
			Bold(true)

	Header = lipgloss.NewStyle().
		Foreground(colors.Dark1).
		Background(colors.Purple).
		Margin(2, 0, 1, 0).
		Padding(0, 1).
		Bold(true)

	// Keybinding styles
	KeyBindMenuWidth = 15

	KeyStyle = lipgloss.NewStyle().
			MarginRight(1).
			Foreground(colors.Medium1)

	KeyActionStyle = lipgloss.NewStyle().
			Foreground(colors.Medium2)

	KeyBindMenu = lipgloss.NewStyle().
			Align(lipgloss.Center)

	KeyBindBox = lipgloss.NewStyle().
			Width(KeyBindMenuWidth)

	Divider = lipgloss.NewStyle().
		Foreground(colors.Medium2).
		Margin(4, 2, 0, 2)
)
