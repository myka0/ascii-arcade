package chess

import (
	"ascii-arcade/internal/chess/renderer/ascii"
	"ascii-arcade/internal/chess/renderer/block"
	"ascii-arcade/internal/chess/renderer/nerdfont"
	t "ascii-arcade/internal/chess/types"
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/overlay"

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

	if m.pawnPromotionTarget != nil {
		return m.viewPawnPromotion(renderer)
	}

	return renderer.View()
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

	// Build the promotion UI
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"Select a piece below.\n",
		lipgloss.JoinHorizontal(lipgloss.Top, knight, bishop, rook, queen),
	)

	return overlay.PlaceNotification(content, mainView)
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
