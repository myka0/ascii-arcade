package gogame

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// GoModel represents the state of the Go game.
type GoModel struct {
	board        Board
	boardHistory map[string]struct{}
	boardSize    int
	lastMove     *Position
	koPoint      *Position
	cursor       Position
	turn         int8
	score        Score
	message      string

	blackCaptures int
	whiteCaptures int

	gameOver    bool
	blackPassed bool
	whitePassed bool
	showLabels  bool

	markingDeadStones bool
	deadStones        map[Position]bool
}

// InitGoModel creates and initializes a new Go game model.
func InitGoModel() *GoModel {
	size := DefaultSize
	m := &GoModel{
		board:        NewBoard(size),
		boardSize:    size,
		lastMove:     nil,
		turn:         Black,
		showLabels:   true,
		cursor:       Position{X: size / 2, Y: size / 2},
		boardHistory: make(map[string]struct{}),
		deadStones:   make(map[Position]bool),
	}

	m.boardHistory[fmt.Sprint(m.board.Cells)] = struct{}{}

	return m
}

// Init implements the Bubble Tea interface for initialization.
func (m *GoModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress and mouse events.
func (m *GoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)
	}

	return m, nil
}

// handleKeyPress handles keyboard input.
func (m *GoModel) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// Always allow ctrl+r to reset
	if msg.String() == "ctrl+r" {
		return handleReset(m.boardSize), nil
	}

	// Allow enter to confirm score when marking dead stones
	if m.markingDeadStones && msg.String() == "enter" {
		return m.handleConfirmScore()
	}

	// Disable all other keypresses on game over screen
	if m.gameOver {
		return m, nil
	}

	switch msg.String() {
	case "up", "w":
		if m.cursor.Y > 0 {
			m.cursor.Y--
		}
	case "down", "s":
		if m.cursor.Y < m.boardSize-1 {
			m.cursor.Y++
		}
	case "left", "a":
		if m.cursor.X > 0 {
			m.cursor.X--
		}
	case "right", "d":
		if m.cursor.X < m.boardSize-1 {
			m.cursor.X++
		}
	case "space", "enter":
		return m.handlePlaceStone(m.cursor)
	case "p":
		return m.handlePass()
	case "r":
		return m.handleResign()
	case "ctrl+v":
		m.cycleBoardSize()
	case "l":
		m.showLabels = !m.showLabels
	}

	return m, nil
}

// handleMouseClick handles mouse interactions.
func (m *GoModel) handleMouseClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only respond to left clicks
	if msg.Mouse().Button != tea.MouseLeft {
		return m, nil
	}

	// Handle dead stone marking phase
	if m.gameOver && m.markingDeadStones {
		// Check if score button was clicked
		if zone.Get("score").InBounds(msg) {
			return m.handleConfirmScore()
		}

		// Check if a board intersection was clicked to toggle dead stone
		for y := 0; y < m.boardSize; y++ {
			for x := 0; x < m.boardSize; x++ {
				label := fmt.Sprintf("%d_%d", x, y)
				if !zone.Get(label).InBounds(msg) {
					continue
				}

				pos := Position{X: x, Y: y}
				return m.handleDeadStoneMarking(pos)
			}
		}

		return m, nil
	}

	// Handle game over UI after scoring is confirmed
	if m.gameOver {
		switch {
		case zone.Get("reset").InBounds(msg):
			return InitGoModel(), nil
		case zone.Get("exit").InBounds(msg):
			return m, func() tea.Msg { return "home" }
		default:
			return m, nil
		}
	}

	// Check pass button
	if zone.Get("pass").InBounds(msg) {
		return m.handlePass()
	}

	// Check resign button
	if zone.Get("resign").InBounds(msg) {
		return m.handleResign()
	}

	// Check if a board intersection was clicked
	for y := 0; y < m.boardSize; y++ {
		for x := 0; x < m.boardSize; x++ {
			label := fmt.Sprintf("%d_%d", x, y)
			if !zone.Get(label).InBounds(msg) {
				continue
			}

			m.cursor = Position{X: x, Y: y}
			return m.handlePlaceStone(m.cursor)
		}
	}

	return m, nil
}

// handlePlaceStone attempts to place a stone at the given position.
func (m *GoModel) handlePlaceStone(pos Position) (tea.Model, tea.Cmd) {
	// Cannot place on occupied position
	if m.board.Cells[pos.Y][pos.X] != Empty {
		m.message = "That intersection is already occupied."
		return m, nil
	}

	// Simulate the placement to check captures
	boardCopy := m.board.Clone()
	boardCopy.PlaceStone(pos, m.turn)

	// Check for suicide
	if boardCopy.GetGroup(pos).Liberties == 0 {
		m.message = fmt.Sprintf(
			"Suicide: stone at %s would have no liberties.",
			formatPosition(pos, m.boardSize),
		)
		return m, nil
	}

	// Check simple ko
	if m.koPoint != nil && m.koPoint.X == pos.X && m.koPoint.Y == pos.Y {
		m.message = "Ko: That move would recapture the ko immediately."
		return m, nil
	}

	// Check superko: no repetition of any previous board position
	if _, exists := m.boardHistory[fmt.Sprint(boardCopy.Cells)]; exists {
		m.message = "Ko: That move would repeat a previous board position (superko)."
		return m, nil
	}

	// Apply the move to the real board
	captured := m.board.PlaceStone(pos, m.turn)
	m.lastMove = &Position{X: pos.X, Y: pos.Y}
	m.boardHistory[fmt.Sprint(m.board.Cells)] = struct{}{}

	if len(captured) > 0 {
		if m.turn == Black {
			m.blackCaptures += len(captured)
		} else {
			m.whiteCaptures += len(captured)
		}

		m.message = fmt.Sprintf("Captured %d stone(s)!", len(captured))
	} else {
		m.message = ""
	}

	// Determine if a ko was created:
	// 1 stone captured AND the capturing stone has 1 liberty
	m.koPoint = nil
	if len(captured) == 1 {
		group := m.board.GetGroup(pos)
		if group.Liberties == 1 {
			m.koPoint = &Position{X: captured[0].X, Y: captured[0].Y}
		}
	}

	// Reset pass state and switch turn
	m.blackPassed = false
	m.whitePassed = false
	m.turn = m.turn * -1

	return m, nil
}

// handlePass records a pass for the current player.
func (m *GoModel) handlePass() (tea.Model, tea.Cmd) {
	if m.turn == Black {
		m.blackPassed = true
	} else {
		m.whitePassed = true
	}

	player := "Black"
	if m.turn == White {
		player = "White"
	}
	m.message = fmt.Sprintf("%s passes.", player)

	// Clear ko point on pass
	m.koPoint = nil

	// End game if both players have passed consecutively
	if m.blackPassed && m.whitePassed {
		m.gameOver = true
		if m.board.HasStones() {
			// Enter dead stone marking phase
			m.markingDeadStones = true
			m.deadStones = make(map[Position]bool)
			m.message = "Mark dead stones by clicking them, then press Enter or [Score] to finish."
		} else {
			// Empty board, calculate score directly
			m.score = CalculateScore(m.board)
		}
		return m, nil
	}

	m.turn = m.turn * -1
	return m, nil
}

// handleResign handles a player resigning.
func (m *GoModel) handleResign() (tea.Model, tea.Cmd) {
	winner := m.turn * -1
	m.gameOver = true

	// The non-resigning player gets the maximum score
	maxScore := float64(m.boardSize * m.boardSize)
	if winner == Black {
		m.score.BlackScore = maxScore
		m.score.WhiteScore = 0
	} else {
		m.score.BlackScore = 0
		m.score.WhiteScore = maxScore + Komi
	}
	m.score.Winner = winner

	return m, nil
}

// handleConfirmScore removes dead stones and calculates the final score.
func (m *GoModel) handleConfirmScore() (tea.Model, tea.Cmd) {
	// Remove dead stones from the board and count them
	removedBlack := 0
	removedWhite := 0
	for pos := range m.deadStones {
		switch m.board.Cells[pos.Y][pos.X] {
		case Black:
			removedBlack++
		case White:
			removedWhite++
		}
		m.board.Cells[pos.Y][pos.X] = Empty
	}

	// Dead stones are effectively captured
	m.blackCaptures += removedWhite
	m.whiteCaptures += removedBlack

	// Calculate final score with dead stones removed
	m.score = CalculateScore(m.board)
	m.markingDeadStones = false

	// Show what was removed
	if removedBlack > 0 || removedWhite > 0 {
		m.message = fmt.Sprintf(
			"Removed %d dead stone(s) (● %d, ○ %d).",
			removedBlack+removedWhite,
			removedBlack,
			removedWhite,
		)
	} else {
		m.message = "No dead stones marked. Final score calculated."
	}

	return m, nil
}

// handleDeadStoneMarking toggles the alive/dead status of a stone.
func (m *GoModel) handleDeadStoneMarking(pos Position) (tea.Model, tea.Cmd) {
	if m.board.Cells[pos.Y][pos.X] != Empty {
		if m.deadStones[pos] {
			// Mark the stone as alive
			delete(m.deadStones, pos)
			m.message = fmt.Sprintf(
				"Unmarked %s as alive. %d dead stone(s) marked.",
				formatPosition(pos, m.boardSize),
				len(m.deadStones),
			)
		} else {
			// Mark the stone as dead
			m.deadStones[pos] = true
			m.message = fmt.Sprintf(
				"Marked %s as dead. %d dead stone(s) marked.",
				formatPosition(pos, m.boardSize),
				len(m.deadStones),
			)
		}
	}

	return m, nil
}

// handleReset resets the game to its initial state.
func handleReset(boardSize int) *GoModel {
	newModel := InitGoModel()
	newModel.boardSize = boardSize
	newModel.board = NewBoard(boardSize)
	newModel.boardHistory = make(map[string]struct{})
	newModel.boardHistory[fmt.Sprint(newModel.board.Cells)] = struct{}{}
	newModel.cursor = Position{X: boardSize / 2, Y: boardSize / 2}
	return newModel
}

// cycleBoardSize cycles through 9x9, 13x13, and 19x19 board sizes.
func (m *GoModel) cycleBoardSize() {
	switch m.boardSize {
	case BoardSize9:
		*m = *handleReset(BoardSize13)
	case BoardSize13:
		*m = *handleReset(BoardSize19)
	case BoardSize19:
		*m = *handleReset(BoardSize9)
	}
}

// colorName returns the name of a color constant.
func colorName(color int8) string {
	switch color {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "Unknown"
	}
}
