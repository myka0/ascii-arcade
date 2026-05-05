package gogame

import (
	"ascii-arcade/internal/colors"

	"charm.land/lipgloss/v2"
)

const (
	gridTopLeft      = " ┏━"
	gridTopIntersect = "━┳━"
	gridTopRight     = "━┓ "
	gridMidLeft      = " ┣━"
	gridMidIntersect = "━╋━"
	gridMidRight     = "━┫ "
	gridBotLeft      = " ┗━"
	gridBotIntersect = "━┻━"
	gridBotRight     = "━┛ "
	gridHorizBar     = "━"
	gridVertBar      = " ┃ "
)

var (
	BoardLine  = lipgloss.NewStyle().Foreground(colors.Medium2).Background(colors.Tan)
	LabelStyle = BoardLine.Bold(true)

	BlackStoneStyle = lipgloss.NewStyle().Foreground(colors.Dark1).Background(colors.Tan)
	WhiteStoneStyle = lipgloss.NewStyle().Foreground(colors.Light1).Background(colors.Tan)
	DeadStoneStyle  = lipgloss.NewStyle().Foreground(colors.Red).Background(colors.Tan)

	CursorStyle = lipgloss.NewStyle().Foreground(colors.Blue).Background(colors.Tan)

	MessageStyle = lipgloss.NewStyle().Foreground(colors.Pink).MarginTop(1)

	Border       = lipgloss.NewStyle().Padding(1, 4, 2, 4).Background(colors.Tan)
	BorderLabels = lipgloss.NewStyle().Padding(1, 4, 1, 2).Background(colors.Tan)

	TitleStyle = lipgloss.NewStyle().
			Foreground(colors.Dark1).
			Background(colors.Purple).
			Padding(0, 1).
			Margin(1, 0).
			Bold(true)

	ButtonStyle = lipgloss.NewStyle().
			Foreground(colors.Dark1).
			Background(colors.Purple).
			Padding(0, 1).
			Margin(1, 2, 0, 0).
			Bold(true)
)
