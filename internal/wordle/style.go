package wordle

import (
  "crossword/internal/colors"
  "github.com/charmbracelet/lipgloss"
)

var (
  KeyCorrect  = style.Purple
  KeyPresent  = style.Pink
  KeyAbsent   = style.Medium
  Text        = style.Light2

  FGKeyCorrect  = lipgloss.NewStyle().Foreground(KeyCorrect)
  FGKeyPresent  = lipgloss.NewStyle().Foreground(KeyPresent)
  FGKeyAbsent   = lipgloss.NewStyle().Foreground(KeyAbsent)
  FGText        = lipgloss.NewStyle().Foreground(Text)

	Border = lipgloss.NewStyle().
    Padding(0, 1).
    Border(lipgloss.NormalBorder())
)
