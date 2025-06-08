package connections

import (
	"crossword/internal/colors"
	"github.com/charmbracelet/lipgloss"
)

var (
	Selected   = colors.Light2
	Background = colors.Dark
	LightText  = colors.Light2
	DarkText   = colors.Dark
	GreyText   = colors.Medium

	BGSelected   = lipgloss.NewStyle().Background(Selected)
	BGBackground = lipgloss.NewStyle().Background(Background)
	FGLightText  = lipgloss.NewStyle().Foreground(LightText)

	NormalCell = lipgloss.NewStyle().
			Foreground(LightText).
			Background(Background).
			MarginLeft(1).
			MarginRight(1).
			Bold(true)

	SelectedCell = lipgloss.NewStyle().
			Foreground(DarkText).
			Background(Selected).
			MarginLeft(1).
			MarginRight(1).
			Bold(true)
)
