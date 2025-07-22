package checkers

import (
	"ascii-arcade/pkg/checkers/renderer/ascii"
	"ascii-arcade/pkg/checkers/renderer/block"
	"ascii-arcade/pkg/checkers/renderer/nerdfont"
	t "ascii-arcade/pkg/checkers/types"
	"ascii-arcade/internal/colors"
	"ascii-arcade/pkg/overlay"

	"image/color"

	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type PieceRenderer interface {
	View() string
	ViewStyledPiece(piece int8, color int8, background color.Color) string
}

const (
	Block = iota
	Ascii
	Nerdfont
)

// View renders the entire Chess board.
func (m CheckersModel) View() string {
	renderer := m.getPieceRenderer(m.renderer)

	if m.gameOver {
		return m.viewGameOver(renderer)
	}

	return renderer.View()
}

// viewGameOver renders the end of game UI.
func (m CheckersModel) viewGameOver(renderer PieceRenderer) string {
	mainView := renderer.View()
	background := colors.Dark2

	// Determine game outcome and assign appropriate styling
	var winner string
	var color color.Color
	switch {
	case m.whiteWins:
		winner = "White wins!"
		color = colors.Orange
	case m.blackWins:
		winner = "Black wins!"
		color = colors.Purple
	default:
		winner = "Stalemate!"
		color = colors.Blue
	}

	winner = lipgloss.NewStyle().Foreground(color).Render(winner)

	buttonStyle := lipgloss.NewStyle().
		Foreground(background).
		Background(color).
		Padding(0, 1)

	// Box used to align buttons
	buttonBox := lipgloss.NewStyle().
		Background(background).
		Width(12)

	// Create interactive buttons and join them side by side
	resetButton := zone.Mark("reset", buttonStyle.Align(lipgloss.Left).Render("Reset"))
	exitButton := zone.Mark("exit", buttonStyle.Align(lipgloss.Right).Render("Exit"))
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		buttonBox.Align(lipgloss.Left).Render(resetButton),
		buttonBox.Align(lipgloss.Right).Render(exitButton),
	)

	return overlay.PlaceNotification(
		mainView,
		"Game over.",
		winner,
		buttons,
	)
}

// getPieceRenderer returns the appropriate piece renderer.
func (m CheckersModel) getPieceRenderer(renderer int) PieceRenderer {
	context := t.RenderContext{
		Board:        m.board,
		Selected:     m.selected,
		ValidMoves:   m.validMoves,
		CaptureMoves: m.captureMoves,
	}

	switch renderer {
	case Block:
		return block.BlockRenderer{RenderContext: context}
	case Ascii:
		return ascii.AsciiRenderer{RenderContext: context}
	case Nerdfont:
		return nerdfont.NerdfontRenderer{RenderContext: context}
	}

	return nil
}
