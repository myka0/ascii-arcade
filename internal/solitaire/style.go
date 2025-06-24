package solitaire

import (
	"ascii-arcade/internal/colors"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	Red      = colors.Red
	White    = colors.Light2
	Empty    = colors.Medium1
	Selected = colors.Blue

	FGRed      = lipgloss.NewStyle().Foreground(Red)
	FGWhite    = lipgloss.NewStyle().Foreground(White)
	FGEmpty    = lipgloss.NewStyle().Foreground(Empty)
	FGSelected = lipgloss.NewStyle().Foreground(Selected)
)
