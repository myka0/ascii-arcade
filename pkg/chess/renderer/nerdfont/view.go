package nerdfont

import (
	t "ascii-arcade/pkg/chess/types"
	"fmt"
	"image/color"

	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type NerdfontRenderer struct {
	t.RenderContext
}

// View renders the full block style chess board.
func (r NerdfontRenderer) View() string {
	boardHeight := len(r.Board)
	pieces := make([][]string, boardHeight)

	// Render each piece on the board with the appropriate style
	for y := range r.Board {
		for x := range r.Board[y] {
			piece := r.Board[y][x]

			// Highlight kings in check
			if piece.Value == King {
				if (piece.Color == White && r.IsWhiteKingInCheck) ||
					(piece.Color == Black && r.IsBlackKingInCheck) {
					pieces[y] = append(pieces[y], r.viewPiece(piece.Value, CheckPiece))
					continue
				}
			}

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

	// Assemble the full view
	return lipgloss.JoinVertical(
		lipgloss.Left,
		r.viewTakenPieces(r.BlackCapturedPieces),
		r.viewBoard(pieces),
		r.viewTakenPieces(r.WhiteCapturedPieces),
	)
}

// viewBoard renders a grid of strings into a styled chessboard.
func (r NerdfontRenderer) viewBoard(board [][]string) string {
	var rows []string

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

		rows = append(rows, r.viewMarginRow(y))
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	rows = append(rows, r.viewMarginRow(len(board)))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (r NerdfontRenderer) viewMarginRow(y int) string {
	top := make([]string, len(r.Board))
	for x := range len(r.Board) {
		isEven := (x+y)%2 == 0

		if y == 0 {
			top[x] = styleCell(isEven, EndMarginEven, EndMarginOdd).Render("▄▄▄▄▄")
		} else if y == len(r.Board) {
			top[x] = styleCell(!isEven, EndMarginEven, EndMarginOdd).Render("▀▀▀▀▀")
		} else {
			top[x] = styleCell(isEven, MarginEven, MarginOdd).Render("▄▄▄▄▄")
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, top[:]...)
}

// viewTakenPieces renders the captured pieces for the given player.
func (r NerdfontRenderer) viewTakenPieces(takenPieces []t.Piece) string {
	cells := make([]string, 16)

	// Render each piece using appropriate style
	for i, piece := range takenPieces {
		cells[i] = r.viewPiece(piece.Value, pieceStyle(piece.Color))
	}

	// Apply cell styling and separate into two rows based on index parity
	var top, bot []string
	for i, content := range cells {
		cell := EmptyCell.Render(content)
		if i%2 == 0 {
			top = append(top, cell)
		} else {
			bot = append(bot, cell)
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, top...)+"\n",
		lipgloss.JoinHorizontal(lipgloss.Top, bot...),
	)
}

// viewPiece renders a single given chess piece.
func (r NerdfontRenderer) viewPiece(piece int8, style lipgloss.Style) string {
	switch piece {
	case Move:
		return r.move(style)
	case Pawn:
		return r.pawn(style)
	case Rook:
		return r.rook(style)
	case Knight:
		return r.knight(style)
	case Bishop:
		return r.bishop(style)
	case Queen:
		return r.queen(style)
	case King:
		return r.king(style)
	default:
		return ""
	}
}

// ViewStyledPiece renders a single fully styled chess piece.
func (r NerdfontRenderer) ViewStyledPiece(piece int8, color int8, background color.Color) string {
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
