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

// renderedCells is a precomputed lookup table for rendered cell strings, indexed by the engine's cell byte value.
var renderedCells [128]string

func init() {
	renderedEmpty := emptyStyle.Render(cellEmpty)
	renderedGhost := ghostStyle.Render(cellGhost)

	// Default every cell to empty so unknown bytes render as blank space
	for i := range renderedCells {
		renderedCells[i] = renderedEmpty
	}

	// Overwrite entries for known piece bytes
	for b, style := range cellStyles {
		renderedCells[b] = style.Render(cellFilled)
	}

	// Ghost piece gets its own style
	renderedCells[ghost] = renderedGhost
}

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

	holdLabel = PanelLabel.Render("Hold")
	nextLabel = PanelLabel.Render("Next")
	infoLabel = PanelLabel.Render("Info")

	InfoLabel = lipgloss.NewStyle().Foreground(colors.Light2).Bold(true)
	InfoValue = lipgloss.NewStyle().Foreground(colors.Purple).Bold(true)

	HoldBox = SidePanel.Width(12).Height(5)
	NextBox = SidePanel.Width(12)
	InfoBox = SidePanel.Width(12)
)
