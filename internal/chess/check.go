package chess

import t "ascii-arcade/internal/chess/types"

// isKingInCheck returns true if the king of the given color is under attack.
func (m *ChessModel) isKingInCheck(color int8) bool {
	var king t.Position

	// Find the king
	for y := range 8 {
		for x := range 8 {
			p := m.board[y][x]
			if p.Value == King && p.Color == color {
				king = pos(x, y)
				break
			}
		}
	}

	// Check all opponent pieces for attacks on the king
	for y := range 8 {
		for x := range 8 {
			p := m.board[y][x]
			if p.Color != color && p.Value != Empty {
				if m.canPieceReach(pos(x, y), king) {
					return true
				}
			}
		}
	}

	return false
}

// wouldBeInCheck simulates moving the king to a square and checks if it would be in check.
func (m *ChessModel) wouldBeInCheck(from, to t.Position) bool {
	piece := m.board[from.Y][from.X]
	captured := m.board[to.Y][to.X]

	// Simulate the move
	m.board[from.Y][from.X] = newEmptyPiece()
	m.board[to.Y][to.X] = piece

	// Check if the king is in check after the move
	isCheck := m.isKingInCheck(piece.Color)

	// Revert the simulated move
	m.board[to.Y][to.X] = captured
	m.board[from.Y][from.X] = piece

	return isCheck
}

// canPieceReach returns true if the given piece can legally reach the given position.
func (m *ChessModel) canPieceReach(from, to t.Position) bool {
	dx := sign(to.X - from.X)
	dy := sign(to.Y - from.Y)
	absX := abs(to.X - from.X)
	absY := abs(to.Y - from.Y)

	switch m.board[from.Y][from.X].Value {
	case Pawn:
		return absX == 1 && to.Y-from.Y == int(m.board[from.Y][from.X].Color)
	case Knight:
		return (absX == 1 && absY == 2) || (absX == 2 && absY == 1)
	case Bishop:
		return absX == absY && m.isPathClearTo(from, to, dx, dy)
	case Rook:
		return (dx == 0 || dy == 0) && m.isPathClearTo(from, to, dx, dy)
	case Queen:
		return (dx == 0 || dy == 0 || absX == absY) && m.isPathClearTo(from, to, dx, dy)
	case King:
		return absX <= 1 && absY <= 1
	}

	return false
}

// isPathClearTo returns true if the path from the given piece to the given position is clear.
func (m *ChessModel) isPathClearTo(from, to t.Position, dx, dy int) bool {
	x, y := from.X+dx, from.Y+dy
	for x != to.X || y != to.Y {
		if !inBounds(x, y) || m.board[y][x].Value != Empty {
			return false
		}
		x += dx
		y += dy
	}
	return true
}

// abs returns the absolute value of the given integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// sign reduces any integer to -1, 0, or 1 based on its sign.
func sign(n int) int {
	if n < 0 {
		return -1
	}
	if n > 0 {
		return 1
	}
	return 0
}
