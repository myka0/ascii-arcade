package chess

import (
	"ascii-arcade/pkg/chess/renderer/ascii"
	"ascii-arcade/pkg/chess/renderer/block"
	"ascii-arcade/pkg/chess/renderer/nerdfont"
	t "ascii-arcade/pkg/chess/types"
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
func (m ChessModel) View() string {
	renderer := m.getPieceRenderer(m.renderer)

	if m.gameOver {
		return m.viewGameOver(renderer)
	}

	if m.pawnPromotionTarget != nil {
		return m.viewPawnPromotion(renderer)
	}

	return renderer.View()
}

// viewGameOver renders the end of game UI.
func (m ChessModel) viewGameOver(renderer PieceRenderer) string {
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

// viewPawnPromotion renders the pawn promotion UI.
func (m ChessModel) viewPawnPromotion(renderer PieceRenderer) string {
	mainView := renderer.View()
	color := m.turn * -1
	background := colors.Dark2

	// Render each promotion piece option
	knight := zone.Mark("knight", renderer.ViewStyledPiece(Knight, color, background))
	bishop := zone.Mark("bishop", renderer.ViewStyledPiece(Bishop, color, background))
	rook := zone.Mark("rook", renderer.ViewStyledPiece(Rook, color, background))
	queen := zone.Mark("queen", renderer.ViewStyledPiece(Queen, color, background))

	pieces := lipgloss.JoinHorizontal(lipgloss.Top, knight, bishop, rook, queen)

	return overlay.PlaceNotification(
		mainView,
		"Select a piece below.",
		pieces,
	)
}

// getPieceRenderer returns the appropriate piece renderer.
func (m ChessModel) getPieceRenderer(renderer int) PieceRenderer {
	context := t.RenderContext{
		Board:      m.board,
		Selected:   m.selected,
		ValidMoves: m.validMoves,

		WhiteCapturedPieces: m.whiteCapturedPieces,
		BlackCapturedPieces: m.blackCapturedPieces,

		IsWhiteKingInCheck: m.isWhiteKingInCheck,
		IsBlackKingInCheck: m.isBlackKingInCheck,
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
