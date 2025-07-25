package connections

import (
	"ascii-arcade/internal/colors"
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	Gap       = 2
	CellWidth = 16
	ClueWidth = CellWidth*4 + Gap*3
)

var (
	Selected   = colors.Light2
	Background = colors.Dark1
	LightText  = colors.Light2
	DarkText   = colors.Dark1
	Special    = colors.Purple

	Color1 = colors.Yellow
	Color2 = colors.Orange
	Color3 = colors.Blue
	Color4 = colors.Pink

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
