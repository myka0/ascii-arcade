package chess

import (
	"ascii-arcade/internal/chess/pieces/ascii"
	"ascii-arcade/internal/chess/pieces/block"
	"ascii-arcade/internal/chess/pieces/nerdfont"
)

type PieceRenderer interface {
	View() string
}

const (
	Block = iota
	Ascii
	Nerdfont
)

// View renders the entire Chess board.
func (m ChessModel) View() string {
	return m.getPieceRenderer(m.renderer).View()
}

// getPieceRenderer returns the appropriate piece renderer.
func (m ChessModel) getPieceRenderer(renderer int) PieceRenderer {
	switch renderer {
	case Block:
		return block.BlockRenderer{
			Board:      m.board,
			Selected:   m.selected,
			ValidMoves: m.validMoves,

			WhiteCapturedPieces: m.whiteCapturedPieces,
			BlackCapturedPieces: m.blackCapturedPieces,

			IsWhiteKingInCheck: m.isWhiteKingInCheck,
			IsBlackKingInCheck: m.isBlackKingInCheck,
		}

	case Ascii:
		return ascii.AsciiRenderer{
			Board:      m.board,
			Selected:   m.selected,
			ValidMoves: m.validMoves,

			WhiteCapturedPieces: m.whiteCapturedPieces,
			BlackCapturedPieces: m.blackCapturedPieces,

			IsWhiteKingInCheck: m.isWhiteKingInCheck,
			IsBlackKingInCheck: m.isBlackKingInCheck,
		}

	case Nerdfont:
		return nerdfont.NerdfontRenderer{
			Board:      m.board,
			Selected:   m.selected,
			ValidMoves: m.validMoves,

			WhiteCapturedPieces: m.whiteCapturedPieces,
			BlackCapturedPieces: m.blackCapturedPieces,

			IsWhiteKingInCheck: m.isWhiteKingInCheck,
			IsBlackKingInCheck: m.isBlackKingInCheck,
		}
	}

	return nil
}
