package chess

import t "ascii-arcade/pkg/chess/types"

// generateValidMoves determines all legal moves for the selected piece
// and stores them in m.validMoves, filtering out any that leave the king in check.
func (m *ChessModel) generateValidMoves(selected t.Position) {
	var moves []t.Position

	switch m.board[selected.Y][selected.X].Value {
	case Pawn:
		moves = m.generatePawnMoves(selected)

	case Rook:
		moves = m.generateLinearMoves(selected, []t.Position{
			{X: 0, Y: -1}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 1, Y: 0}, // Vertical & horizontal
		})

	case Knight:
		moves = m.generateStepMoves(selected, []t.Position{
			{X: 2, Y: 1}, {X: 1, Y: 2}, {X: -1, Y: 2}, {X: -2, Y: 1},
			{X: -2, Y: -1}, {X: -1, Y: -2}, {X: 1, Y: -2}, {X: 2, Y: -1}, // L-shaped jumps
		})

	case Bishop:
		moves = m.generateLinearMoves(selected, []t.Position{
			{X: -1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: -1}, {X: 1, Y: 1}, // Diagonals
		})

	case Queen:
		moves = m.generateLinearMoves(selected, []t.Position{
			{X: 0, Y: -1}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 1, Y: 0}, // Straight
			{X: -1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: -1}, {X: 1, Y: 1}, // Diagonal
		})

	case King:
		moves = m.generateStepMoves(selected, []t.Position{
			{X: 0, Y: -1}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 1, Y: 0},
			{X: -1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: -1}, {X: 1, Y: 1}, // Adjacent squares
		})

		// Add castling options if available
		moves = append(moves, m.generateCastlingMoves(selected)...)
	}

	// Filter out any moves that would place own king in check
	m.addMovesIfNotCheck(selected, moves)
}

// generatePawnMoves returns all valid pawn moves from the given position.
func (m *ChessModel) generatePawnMoves(pawn t.Position) []t.Position {
	validMoves := make([]t.Position, 0)

	x, y := pawn.X, pawn.Y
	dir := int(m.turn)

	startRow := 6
	if m.turn == Black {
		startRow = 1
	}

	// Single forward move
	if inBounds(x, y+dir) && m.board[y+dir][x].Value == Empty {
		validMoves = append(validMoves, pos(x, y+dir))

		// Double forward move from starting row
		if y == startRow && m.board[y+2*dir][x].Value == Empty {
			validMoves = append(validMoves, pos(x, y+2*dir))
		}
	}

	// Captures
	for _, dx := range []int{-1, 1} {
		nx, ny := x+dx, y+dir
		if inBounds(nx, ny) {
			target := m.board[ny][nx]
			if target.Value != Empty && target.Color != m.board[pawn.Y][pawn.X].Color {
				validMoves = append(validMoves, pos(nx, ny))
			}
		}
	}

	// En passant
	if m.enPassantTarget != nil {
		target := *m.enPassantTarget

		// Ensure the target is diagonally forward by 1 square
		if target.Y == pawn.Y+dir && abs(target.X-pawn.X) == 1 {
			m.generateEnPassant(pawn, target)
		}
	}

	return validMoves
}

// generateEnPassant creates an en passant capture and adds the move if it doesn't leave the king in check.
func (m *ChessModel) generateEnPassant(pawn t.Position, target t.Position) {
	piece := m.board[pawn.Y][pawn.X]

	// Remove captured pawn
	capturedPawn := m.board[pawn.Y][target.X]
	m.board[pawn.Y][target.X] = newEmptyPiece()

	// Move pawn to en passant square
	m.board[target.Y][target.X] = piece
	m.board[pawn.Y][pawn.X] = newEmptyPiece()

	// Only allow move if it doesn't leave king in check
	if !m.isKingInCheck(piece.Color) {
		m.validMoves = append(m.validMoves, target)
	}

	// Restore original board state
	m.board[pawn.Y][pawn.X] = piece
	m.board[pawn.Y][target.X] = capturedPawn
	m.board[target.Y][target.X] = newEmptyPiece()
}

// generateStepMoves generates valid single step moves in specified directions.
func (m *ChessModel) generateStepMoves(start t.Position, directions []t.Position) []t.Position {
	var validMoves []t.Position

	for _, d := range directions {
		// Calculate the target square
		x, y := start.X+d.X, start.Y+d.Y
		if !inBounds(x, y) {
			continue
		}

		// Add move if the square is empty or can be captured
		target := m.board[y][x]
		if target.Value == Empty || target.Color != m.board[start.Y][start.X].Color {
			validMoves = append(validMoves, pos(x, y))
		}
	}

	return validMoves
}

// generateLinearMoves generates valid moves along one or more directions.
func (m *ChessModel) generateLinearMoves(start t.Position, directions []t.Position) []t.Position {
	validMoves := make([]t.Position, 0)

	for _, d := range directions {
		x, y := start.X, start.Y

		for {
			x += d.X
			y += d.Y

			if !inBounds(x, y) {
				break
			}

			target := m.board[y][x]

			// If square is empty, add move and continue searching in that direction
			if target.Value == Empty {
				validMoves = append(validMoves, pos(x, y))
				continue
			}

			// If it's an opponent's piece, capture is possible
			if target.Color != m.board[start.Y][start.X].Color {
				validMoves = append(validMoves, pos(x, y)) // capture
			}

			break // blocked
		}
	}

	return validMoves
}

// generateCastlingMoves returns castling positions if the king is allowed to castle.
func (m *ChessModel) generateCastlingMoves(king t.Position) []t.Position {
	var validMoves []t.Position
	color := m.board[king.Y][king.X].Color
	y := king.Y

	// Kingside castling
	if ((color == White && m.whiteCastleKingside) || (color == Black && m.blackCastleKingside)) &&
		m.isPathClearTo(king, pos(7, y), 1, 0) &&
		!m.isKingInCheck(color) &&
		!m.wouldBeInCheck(king, pos(5, y)) &&
		!m.wouldBeInCheck(king, pos(6, y)) {

		validMoves = append(validMoves, pos(6, y))
	}

	// Queenside castling
	if ((color == White && m.whiteCastleQueenside) || (color == Black && m.blackCastleQueenside)) &&
		m.isPathClearTo(king, pos(0, y), -1, 0) &&
		!m.isKingInCheck(color) &&
		!m.wouldBeInCheck(king, pos(3, y)) &&
		!m.wouldBeInCheck(king, pos(2, y)) {

		validMoves = append(validMoves, pos(2, y))
	}

	return validMoves
}

// hasValidMoves checks if the given color has at least one legal move.
func (m *ChessModel) hasValidMoves(color int8) bool {
	for y := range 8 {
		for x := range 8 {
			p := m.board[y][x]

			// Skip empty squares or pieces of the opposite color
			if p.Color != color || p.Value == Empty {
				continue
			}

			// Generate and check if there are valid moves
			m.generateValidMoves(pos(x, y))
			if len(m.validMoves) > 0 {
				m.validMoves = nil
				return true
			}
		}
	}

	return false
}

// addMovesIfNotCheck filters the given moves to only include those that don't
// leave the player's king in check.
func (m *ChessModel) addMovesIfNotCheck(selected t.Position, moves []t.Position) {
	for _, move := range moves {
		if !m.wouldBeInCheck(selected, move) {
			m.validMoves = append(m.validMoves, move)
		}
	}
}

// inBounds returns true if the given position is within the bounds of the board.
func inBounds(x, y int) bool {
	return x >= 0 && x < 8 && y >= 0 && y < 8
}
