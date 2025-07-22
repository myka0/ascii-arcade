package checkers

import (
	"fmt"
	"slices"

	t "ascii-arcade/pkg/checkers/types"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

const (
	Empty = iota
	Move
	Pawn
	King
)

const (
	White = -1
	Black = 1
)

type CheckersModel struct {
	renderer     int
	board        [][]t.Piece
	selected     t.Position
	validMoves   []t.Position
	captureMoves []t.CaptureMove
	turn         int8

	whitePiecesLeft int
	blackPiecesLeft int

	whiteWins bool
	blackWins bool
	gameOver  bool
}

// InitCheckersModel creates and initializes a new checkers model.
func InitCheckersModel() *CheckersModel {
	m := CheckersModel{
		renderer: Ascii,
		board:    InitCheckersBoard(),
		selected: pos(-1, -1),
		turn:     White,

		whitePiecesLeft: 12,
		blackPiecesLeft: 12,

		whiteWins: false,
		blackWins: false,
		gameOver:  false,
	}

	return &m
}

// InitCheckersBoard initializes and returns the standard starting position of a checkersboard.
func InitCheckersBoard() [][]t.Piece {
	board := make([][]t.Piece, 8)

	// Set up black pieces
	for y := range 3 {
		for x := range 8 {
			// Place pawns on odd squares
			if (x+y)%2 == 1 {
				board[y] = append(board[y], t.Piece{Value: Pawn, Color: Black})
			} else {
				board[y] = append(board[y], newEmptyPiece())
			}
		}
	}

	// Set up empty rows
	for range 8 {
		board[3] = append(board[3], newEmptyPiece())
		board[4] = append(board[4], newEmptyPiece())
	}

	// Set up white pieces
	for y := 5; y < 8; y++ {
		for x := range 8 {
			// Place pawns on odd squares
			if (x+y)%2 == 1 {
				board[y] = append(board[y], t.Piece{Value: Pawn, Color: White})
			} else {
				board[y] = append(board[y], newEmptyPiece())
			}
		}
	}

	return board
}

// Init implements the Bubble Tea interface for initialization.
func (m CheckersModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress and mouse events to update the Checkers game state.
func (m *CheckersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle keyboard input
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			return InitCheckersModel(), nil
		case "ctrl+v":
			m.renderer = (m.renderer + 1) % 3
			return m, nil
		}

	// Handle mouse input
	case tea.MouseMsg:
		switch msg := msg.(type) {
		case tea.MouseClickMsg:
			return m.handleMouseClick(msg)
		}
	}

	return m, nil
}

// handleMouseClick handles mouse interactions.
func (m *CheckersModel) handleMouseClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only respond to left clicks
	if msg.Mouse().Button != tea.MouseLeft {
		return m, nil
	}

	// Handle game over UI
	if m.gameOver {
		switch {
		case zone.Get("reset").InBounds(msg):
			return InitCheckersModel(), nil
		case zone.Get("exit").InBounds(msg):
			return m, func() tea.Msg { return "home" }
		default:
			return m, nil
		}
	}

	// Handle capture moves
	if len(m.captureMoves) > 0 {
		for _, move := range m.captureMoves {
			label := fmt.Sprint(move.To.Y*len(m.board) + move.To.X)
			if zone.Get(label).InBounds(msg) {
				m.handleMovePiece(move.From, move.To)
				return m, nil
			}
		}

		return m, nil
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
					return m, nil
				}

				// Deselect if clicked the same square again
				if m.selected.X == x && m.selected.Y == y {
					m.selected = pos(-1, -1)
					m.validMoves = nil
					return m, nil
				}

				// Select a different piece
				if piece.Color == m.turn {
					m.selected = clicked
					m.validMoves = nil
					m.generateValidMoves(clicked)
				}

				return m, nil
			}

			// If no piece is selected piece, select the clicked piece
			if piece.Value != Empty && piece.Color == m.turn {
				m.selected = clicked
				m.validMoves = nil
				m.generateValidMoves(clicked)
			}

			return m, nil
		}
	}

	return m, nil
}

func (m *CheckersModel) handleMovePiece(from, to t.Position) {
	piece := m.board[from.Y][from.X]

	// If a pawn moves to the end of the board, it becomes king
	if piece.Value == Pawn && ((m.turn == White && to.Y == 0) || (m.turn == Black && to.Y == 7)) {
		piece.Value = King
	}

	// Move the piece
	m.board[to.Y][to.X] = piece
	m.board[from.Y][from.X] = newEmptyPiece()

	// Clear values
	m.selected = pos(-1, -1)
	m.validMoves = nil

	if len(m.captureMoves) > 0 {
		// Calculate the midpoint
		midX := (from.X + to.X) / 2
		midY := (from.Y + to.Y) / 2

		// Remove captured piece
		captured := m.board[midY][midX]
		m.board[midY][midX] = newEmptyPiece()

		// Update piece counts
		switch captured.Color {
		case White:
			m.whitePiecesLeft--
		case Black:
			m.blackPiecesLeft--
		}

		// Check if the game has ended
		m.whiteWins = m.blackPiecesLeft == 0
		m.blackWins = m.whitePiecesLeft == 0
		m.gameOver = m.whiteWins || m.blackWins

		// Check if more captures are possible
		m.generateAllCaptureMoves(piece.Color)
		if len(m.captureMoves) > 0 {
			return
		}
	}

	// Swap turns and generate possible captures
	m.turn = m.turn * -1
	m.generateAllCaptureMoves(m.turn)
}

// newEmptyPiece returns a new empty piece.
func newEmptyPiece() t.Piece {
	return t.Piece{Color: Empty, Value: Empty}
}

// pos returns a new position.
func pos(x, y int) t.Position {
	return t.Position{X: x, Y: y}
}
