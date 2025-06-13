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
	Mistake    = colors.Purple

	Color1 = colors.Purple
	Color2 = colors.Pink
	Color3 = colors.Blue
	Color4 = colors.Yellow

	BGSelected   = lipgloss.NewStyle().Background(Selected)
	BGBackground = lipgloss.NewStyle().Background(Background)
	FGLightText  = lipgloss.NewStyle().Foreground(LightText)

	Border = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder())

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

	MistakeCell = lipgloss.NewStyle().
			Foreground(Mistake).
			Bold(true)

	RevealedLine = lipgloss.NewStyle().
			Foreground(Background).
			Bold(true)
)
