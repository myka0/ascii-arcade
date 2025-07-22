package wordle

import (
	"ascii-arcade/internal/colors"
	"github.com/charmbracelet/lipgloss/v2"
)

var (
	KeyCorrect = colors.Purple
	KeyPresent = colors.Pink
	KeyAbsent  = colors.Medium2
	Text       = colors.Light2

	FGKeyCorrect = lipgloss.NewStyle().Foreground(KeyCorrect)
	FGKeyPresent = lipgloss.NewStyle().Foreground(KeyPresent)
	FGKeyAbsent  = lipgloss.NewStyle().Foreground(KeyAbsent)
	FGText       = lipgloss.NewStyle().Foreground(Text)

	Border = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder())
)
