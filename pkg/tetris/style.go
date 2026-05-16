package tetris

import (
	"ascii-arcade/internal/colors"

	"charm.land/lipgloss/v2"
)

const (
	cellFilled = "██"
	cellEmpty  = "  "
	cellGhost  = "░░"
)

const (
	pieceI byte = 'I'
	pieceO byte = 'O'
	pieceT byte = 'T'
	pieceS byte = 'S'
	pieceZ byte = 'Z'
	pieceJ byte = 'J'
	pieceL byte = 'L'
	ghost  byte = 'G'
	empty  byte = '0'
)

var cellStyles = map[byte]lipgloss.Style{
	pieceI: lipgloss.NewStyle().Foreground(colors.Cyan),
	pieceO: lipgloss.NewStyle().Foreground(colors.Yellow),
	pieceT: lipgloss.NewStyle().Foreground(colors.Purple),
	pieceS: lipgloss.NewStyle().Foreground(colors.Green),
	pieceZ: lipgloss.NewStyle().Foreground(colors.Red),
	pieceJ: lipgloss.NewStyle().Foreground(colors.Blue),
	pieceL: lipgloss.NewStyle().Foreground(colors.Orange),
}

var (
	ghostStyle = lipgloss.NewStyle().Foreground(colors.Medium1)
	emptyStyle = lipgloss.NewStyle().Foreground(colors.Dark1)
)

var (
	Playfield = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colors.Medium1).
			Padding(0, 1)

	SidePanel = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colors.Medium1).
			Padding(0, 1).
			MarginRight(1).
			MarginLeft(1)

	PanelLabel = lipgloss.NewStyle().
			Foreground(colors.Light2).
			Bold(true).
			Underline(true)

	InfoLabel = lipgloss.NewStyle().Foreground(colors.Light2).Bold(true)
	InfoValue = lipgloss.NewStyle().Foreground(colors.Purple).Bold(true)

	HoldBox = SidePanel.Width(12).Height(5)
	NextBox = SidePanel.Width(12)
	InfoBox = SidePanel.Width(12)
)
