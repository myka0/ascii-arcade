package connections

import (
	"crossword/internal/colors"
	"github.com/charmbracelet/lipgloss"
)

const (
	Gap       = 2
	CellWidth = 16
	ClueWidth = CellWidth*4 + Gap*3
)

var (
	Selected   = colors.Light2
	Background = colors.Dark
	LightText  = colors.Light2
	DarkText   = colors.Dark
	GreyText   = colors.Medium
	Special    = colors.Purple

	Color1 = colors.Orange
	Color2 = colors.Pink
	Color3 = colors.Blue
	Color4 = colors.Yellow

	BGSelected   = lipgloss.NewStyle().Background(Selected)
	BGBackground = lipgloss.NewStyle().Background(Background)
	FGSpecial    = lipgloss.NewStyle().Foreground(Special)
	FGLightText  = lipgloss.NewStyle().Foreground(LightText)

	Border = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder())

	NormalCell = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(LightText).
			Background(Background).
			Padding(1, 0).
			Margin(0, 1).
			Width(CellWidth).
			Bold(true)

	SelectedCell = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(DarkText).
			Background(Selected).
			Padding(1, 0).
			Margin(0, 1).
			Width(CellWidth).
			Bold(true)

	RevealedLine = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(Background).
			Padding(1, 0).
			Width(ClueWidth).
			Bold(true)

	MistakeCell = lipgloss.NewStyle().
			Foreground(Special).
			Bold(true)

	Button = lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(DarkText).
		Background(Special).
		Padding(0, 2).
		Margin(0, 4).
		Width(CellWidth).
		Bold(true)
)
