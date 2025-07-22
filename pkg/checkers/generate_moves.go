package checkers

import (
	t "ascii-arcade/pkg/checkers/types"
)

// generateValidMoves determines all legal moves for the selected piece
func (m *CheckersModel) generateValidMoves(selected t.Position) {
	// Determine valid directions based on piece type
	var directions []t.Position
	switch m.board[selected.Y][selected.X].Value {
	case Pawn:
		dy := int(m.board[selected.Y][selected.X].Color)
		directions = []t.Position{
			{X: -1, Y: dy}, {X: 1, Y: dy},
		}
	case King:
		directions = []t.Position{
			{X: -1, Y: -1}, {X: 1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: 1},
		}
	}

	// Add all valid moves to the list
	m.validMoves = m.addValidMoves(selected, directions)
}

// addValidMoves returns all single step non-capturing moves for a given piece
func (m *CheckersModel) addValidMoves(from t.Position, directions []t.Position) []t.Position {
	validMoves := make([]t.Position, 0)

	for _, d := range directions {
		x, y := from.X+d.X, from.Y+d.Y

		if inBounds(x, y) && m.board[y][x].Value == Empty {
			validMoves = append(validMoves, pos(x, y))
		}
	}

	return validMoves
}

// generateAllCaptureMoves finds and stores all possible capture moves for a given color.
func (m *CheckersModel) generateAllCaptureMoves(color int8) {
	captures := make([]t.CaptureMove, 0)

	for y := range 8 {
		for x := range 8 {
			p := m.board[y][x]

			// Skip empty or opponent pieces
			if p.Value == Empty || p.Color != color {
				continue
			}

			from := pos(x, y)

			// Determine valid directions based on piece type
			var directions []t.Position
			switch p.Value {
			case Pawn:
				dy := int(m.board[from.Y][from.X].Color)
				directions = []t.Position{
					{X: -1, Y: dy}, {X: 1, Y: dy},
				}
			case King:
				directions = []t.Position{
					{X: -1, Y: -1}, {X: 1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: 1},
				}
			}

			// Append any valid capture moves from this piece
			captures = append(captures, m.addValidCaptures(from, directions)...)
		}
	}

	m.captureMoves = captures
}

// addValidCaptures returns capture moves from a position in specified directions.
func (m *CheckersModel) addValidCaptures(from t.Position, directions []t.Position) []t.CaptureMove {
	moves := make([]t.CaptureMove, 0)

	for _, d := range directions {
		mx, my := from.X+d.X, from.Y+d.Y
		jx, jy := from.X+2*d.X, from.Y+2*d.Y

		if !inBounds(mx, my) || !inBounds(jx, jy) {
			continue
		}

		mid := m.board[my][mx]
		dest := m.board[jy][jx]

		if mid.Value != Empty && mid.Color != m.board[from.Y][from.X].Color && dest.Value == Empty {
			moves = append(moves, t.CaptureMove{
				From: from,
				To:   pos(jx, jy),
			})
		}
	}

	return moves
}

// inBounds returns true if the given position is within the bounds of the board.
func inBounds(x, y int) bool {
	return x >= 0 && x < 8 && y >= 0 && y < 8
}
