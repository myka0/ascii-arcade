package checkers

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"
	"ascii-arcade/pkg/checkers/renderer/ascii"
	"ascii-arcade/pkg/checkers/renderer/block"
	"ascii-arcade/pkg/checkers/renderer/nerdfont"
	t "ascii-arcade/pkg/checkers/types"
	"image/color"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
func (m CheckersModel) View() tea.View {
	renderer := m.getPieceRenderer(m.renderer)

	if m.gameOver {
		return tea.NewView(m.viewGameOver(renderer))
	}

	return tea.NewView(renderer.View())
}

// viewGameOver renders the end of game UI.
func (m CheckersModel) viewGameOver(renderer PieceRenderer) string {
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
	return components.GameOver(color, renderer.View(), winner)
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
