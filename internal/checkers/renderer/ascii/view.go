package ascii

import (
	t "ascii-arcade/internal/checkers/types"
	"fmt"
	"image/color"

	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// TODO Create ascii background for each piece
type AsciiRenderer struct {
	t.RenderContext
}

// View renders the full ascii style checkers board.
func (r AsciiRenderer) View() string {
	boardHeight := len(r.Board)
	pieces := make([][]string, boardHeight)

	// Render each piece on the board with the appropriate style
	for y := range r.Board {
		for x := range r.Board[y] {
			piece := r.Board[y][x]
			pieces[y] = append(pieces[y], r.viewPiece(piece.Value, pieceStyle(piece.Color)))
		}
	}

	// Highlight the currently selected piece
	sel := r.Selected
	if sel.X != -1 && sel.Y != -1 {
		pieces[sel.Y][sel.X] = r.viewPiece(r.Board[sel.Y][sel.X].Value, SelectedPiece)
	}

	// Highlight valid move destinations
	for _, move := range r.ValidMoves {
		dest := r.Board[move.Y][move.X]

		if dest.Value != Empty {
			// Show takeable piece with a different color
			pieces[move.Y][move.X] = r.viewPiece(dest.Value, TakePiece)
		} else {
			// Show possible move marker
			pieces[move.Y][move.X] = r.viewPiece(Move, SelectedPiece)
		}
	}

	// Highlight valid captures
	for _, move := range r.CaptureMoves {
		pieces[move.To.Y][move.To.X] = r.viewPiece(Move, SelectedPiece)
		from := r.Board[move.From.Y][move.From.X]
		pieces[move.From.Y][move.From.X] = r.viewPiece(from.Value, SelectedPiece)
	}

	// Assemble the full view
	return r.viewBoard(pieces)
}

// viewBoard renders a grid of strings into a styled checkersboard.
func (r AsciiRenderer) viewBoard(board [][]string) string {
	rows := make([]string, len(board))

	for y := range board {
		var row []string
		for x, content := range board[y] {
			// Determine if cell is even for checkerboard pattern
			isEven := (x+y)%2 == 0
			style := styleCell(isEven, EvenCell, OddCell)

			// Use a unique zone label for mouse interactivity
			label := fmt.Sprint(y*len(board) + x)
			cell := zone.Mark(label, style.Render(content))

			row = append(row, cell)
		}

		rows[y] = lipgloss.JoinHorizontal(lipgloss.Top, row...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// viewPiece renders a single given checkers piece.
func (r AsciiRenderer) viewPiece(piece int8, style lipgloss.Style) string {
	switch piece {
	case Move:
		return r.move(style)
	case Pawn:
		return r.pawn(style)
	case King:
		return r.king(style)
	default:
		return ""
	}
}

// ViewStyledPiece renders a single fully styled checkers piece.
func (r AsciiRenderer) ViewStyledPiece(piece int8, color int8, background color.Color) string {
	return EmptyCell.Background(background).Render(r.viewPiece(piece, pieceStyle(color)))
}

// pieceStyle returns the appropriate style for the given piece color.
func pieceStyle(color int8) lipgloss.Style {
	if color == White {
		return WhitePiece
	}
	return BlackPiece
}

// styleCell selects between two styles based on a boolean condition.
// Used to implement checkerboard patterns in the grid.
func styleCell(parity bool, firstStyle, secondStyle lipgloss.Style) lipgloss.Style {
	if parity {
		return firstStyle
	}
	return secondStyle
}
