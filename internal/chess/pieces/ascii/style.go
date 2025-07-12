package ascii

import (
	"ascii-arcade/internal/colors"

	"github.com/charmbracelet/lipgloss/v2"
)

const (
	Empty = iota
	Move
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

const (
	White = -1
	Black = 1
)

var (
	Width  = 9
	Height = 4

	CWhite    = colors.Orange
	CBlack    = colors.Purple
	CSelected = colors.Blue
	CTake     = colors.Pink
	CCheck    = colors.Red

	WhitePiece    = lipgloss.NewStyle().Foreground(CWhite).Bold(true)
	BlackPiece    = lipgloss.NewStyle().Foreground(CBlack).Bold(true)
	SelectedPiece = lipgloss.NewStyle().Foreground(CSelected).Bold(true)
	TakePiece     = lipgloss.NewStyle().Foreground(CTake).Bold(true)
	CheckPiece    = lipgloss.NewStyle().Foreground(CCheck).Bold(true)

	EmptyCell = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(Width).
			Height(Height)

	OddCell = EmptyCell.
		Background(colors.Light2)

	EvenCell = EmptyCell.
			Background(colors.Dark)
)
