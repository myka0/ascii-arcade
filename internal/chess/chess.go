package chess

import (
	"fmt"
	"slices"

	t "ascii-arcade/internal/chess/types"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

const (
	Empty = iota
	Move
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

const (
	White = -1
	Black = 1
)

// TODO Win condition and stalemate
// TODO Pawn romotion
type ChessModel struct {
	renderer   int
	board      [][]t.Piece
	selected   t.Position
	validMoves []t.Position
	turn       int8

	whiteCapturedPieces []t.Piece
	blackCapturedPieces []t.Piece

	whiteCastleKingside  bool
	whiteCastleQueenside bool
	blackCastleKingside  bool
	blackCastleQueenside bool

	enPassantTarget *t.Position

	isWhiteKingInCheck bool
	isBlackKingInCheck bool
}

// InitChessModel creates and initializes a new chess model.
func InitChessModel() *ChessModel {
	m := ChessModel{
		renderer: Block,
		board:    InitChessBoard(),
		selected: pos(-1, -1),
		turn:     White,

		whiteCastleKingside:  true,
		whiteCastleQueenside: true,
		blackCastleKingside:  true,
		blackCastleQueenside: true,

		isWhiteKingInCheck: false,
		isBlackKingInCheck: false,
	}

	return &m
}

// InitChessBoard initializes and returns the standard starting position of a chessboard.
func InitChessBoard() [][]t.Piece {
	pieces := make([][]t.Piece, 8)

	// Black pieces
	pieces[0] = []t.Piece{
		{Color: Black, Value: Rook},
		{Color: Black, Value: Knight},
		{Color: Black, Value: Bishop},
		{Color: Black, Value: Queen},
		{Color: Black, Value: King},
		{Color: Black, Value: Bishop},
		{Color: Black, Value: Knight},
		{Color: Black, Value: Rook},
	}

	// Black pawns
	for range 8 {
		pieces[1] = append(pieces[1], t.Piece{Color: Black, Value: Pawn})
	}

	// Empty squares
	for y := 2; y < 6; y++ {
		for range 8 {
			pieces[y] = append(pieces[y], newEmptyPiece())
		}
	}

	// White pawns
	for range 8 {
		pieces[6] = append(pieces[6], t.Piece{Color: White, Value: Pawn})
	}

	// White pieces
	pieces[7] = []t.Piece{
		{Color: White, Value: Rook},
		{Color: White, Value: Knight},
		{Color: White, Value: Bishop},
		{Color: White, Value: Queen},
		{Color: White, Value: King},
		{Color: White, Value: Bishop},
		{Color: White, Value: Knight},
		{Color: White, Value: Rook},
	}

	return pieces
}

// Init implements the Bubble Tea interface for initialization.
func (m ChessModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress and mouse events to update the Chess game state.
func (m *ChessModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle keyboard input
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			m.renderer = (m.renderer + 1) % 3
			return m, nil
		}

	// Handle mouse input
	case tea.MouseMsg:
		switch msg := msg.(type) {
		case tea.MouseClickMsg:
			m.handleMouseClick(msg)
		}
	}

	return m, nil
}

// handleMouseClick handles mouse interactions.
func (m *ChessModel) handleMouseClick(msg tea.MouseMsg) {
	// Only respond to left clicks
	if msg.Mouse().Button != tea.MouseLeft {
		return
	}

	// Check if a piece was clicked
	for y := range m.board {
		for x := range m.board[y] {
			// Each square is labeled by its index
			label := fmt.Sprint(y*len(m.board) + x)
			if !zone.Get(label).InBounds(msg) {
				continue
			}

			piece := m.board[y][x]
			clicked := pos(x, y)

			// If a piece is already selected
			if m.selected.X != -1 && m.selected.Y != -1 {
				// Move if clicked on a valid move position
				if slices.Contains(m.validMoves, clicked) {
					m.handleMovePiece(m.selected, clicked)
					return
				}

				// Deselect if clicked the same square again
				if m.selected.X == x && m.selected.Y == y {
					m.selected = pos(-1, -1)
					m.validMoves = nil
					return
				}

				// Select a different piece
				if piece.Color == m.turn {
					m.selected = clicked
					m.validMoves = nil
					m.generateValidMoves(clicked)
				}

				return
			}

			// If no piece is selected piece, select the clicked piece
			if piece.Value != Empty && piece.Color == m.turn {
				m.selected = clicked
				m.validMoves = nil
				m.generateValidMoves(clicked)
			}

			return
		}
	}
}

// handleMovePiece executes a move from one position to another.
func (m *ChessModel) handleMovePiece(from, to t.Position) {
	piece := m.board[from.Y][from.X]

	// Handle En Passant capture
	if m.enPassantTarget != nil && to == *m.enPassantTarget {
		if piece.Value == Pawn {
			// Remove the captured pawn
			m.addTakenPiece(m.board[from.Y][to.X])
			m.board[from.Y][to.X] = newEmptyPiece()
		}
	}

	// Clear previous en passant target
	m.enPassantTarget = nil

	// Set new en passant target if pawn moves two spaces forward
	if piece.Value == Pawn && abs(to.Y-from.Y) == 2 {
		passY := (from.Y + to.Y) / 2
		m.enPassantTarget = &t.Position{X: from.X, Y: passY}
	}

	// Disable castling rights for king if it moves
	if piece.Value == King {
		if piece.Color == White {
			m.whiteCastleKingside = false
			m.whiteCastleQueenside = false
		} else {
			m.blackCastleKingside = false
			m.blackCastleQueenside = false
		}
	}

	// Disable castling rights for rooks if they move
	if piece.Value == Rook {
		switch {
		case piece.Color == White && from.Y == 7 && from.X == 0:
			m.whiteCastleQueenside = false
		case piece.Color == White && from.Y == 7 && from.X == 7:
			m.whiteCastleKingside = false
		case piece.Color == Black && from.Y == 0 && from.X == 0:
			m.blackCastleQueenside = false
		case piece.Color == Black && from.Y == 0 && from.X == 7:
			m.blackCastleKingside = false
		}
	}

	// Handle castling movement
	if piece.Value == King {
		switch {
		// White kingside castle
		case piece.Color == White && from.Y == 7 && from.X == 4 && to.X == 6:
			m.board[7][5] = m.board[7][7]
			m.board[7][7] = newEmptyPiece()

		// White queenside castle
		case piece.Color == White && from.Y == 7 && from.X == 4 && to.X == 2:
			m.board[7][3] = m.board[7][0]
			m.board[7][0] = newEmptyPiece()

		// Black kingside castle
		case piece.Color == Black && from.Y == 0 && from.X == 4 && to.X == 6:
			m.board[0][5] = m.board[0][7]
			m.board[0][7] = newEmptyPiece()

		// Black queenside castle
		case piece.Color == Black && from.Y == 0 && from.X == 4 && to.X == 2:
			m.board[0][3] = m.board[0][0]
			m.board[0][0] = newEmptyPiece()
		}
	}

	m.addTakenPiece(m.board[to.Y][to.X])

	// Move the piece
	m.board[to.Y][to.X] = piece
	m.board[from.Y][from.X] = newEmptyPiece()

	// Clear values
	m.selected = pos(-1, -1)
	m.validMoves = nil
	m.turn = m.turn * -1

	// Update king check status
	m.isWhiteKingInCheck = m.isKingInCheck(White)
	m.isBlackKingInCheck = m.isKingInCheck(Black)
}

// addTakenPiece adds a piece to the appropriate captured piece list.
func (m *ChessModel) addTakenPiece(piece t.Piece) {
	if piece.Value != Empty {
		if piece.Color == White {
			m.blackCapturedPieces = append(m.blackCapturedPieces, piece)
		} else {
			m.whiteCapturedPieces = append(m.whiteCapturedPieces, piece)
		}
	}
}

// newEmptyPiece returns a new empty piece.
func newEmptyPiece() t.Piece {
	return t.Piece{Color: Empty, Value: Empty}
}

// pos returns a new position.
func pos(x, y int) t.Position {
	return t.Position{X: x, Y: y}
}
