package crossword

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"unicode"
)

// Position represents a 2D coordinate in the crossword grid.
type Position struct {
	X, Y int
}

// CrosswordModel represents the state of a crossword puzzle game.
type CrosswordModel struct {
	// Game data
	date        string
	acrossClues []string
	downClues   []string
	answer      [15][15]byte
	grid        [15][15]byte

	// Grid metadata
	gridNums    [15][15]int
	clueIndices [15][15]Position

	// Current state
	clue           int
	cursor         Position
	isAcross       bool
	isAcrossSolved []bool
	isDownSolved   []bool
	movementAxis   *int
	orthoAxis      *int

	// Game state
	incorrect    [15][15]bool
	correctCount int
	filledCount  int
	autoCheck    bool
	message      string
}

// InitCrosswordModel creates and initializes a new crossword model.
// It loads puzzle data from file and sets up the initial game state.
func InitCrosswordModel() *CrosswordModel {
	date, err := GetLatestDate()
	if err != nil {
		fmt.Println("Failed to get latest date:", err)
	}

	m, err := LoadFromFile(date)
	if err != nil {
		fmt.Println("Failed to load crossword:", err)
	}

	// Set initial movement direction
	m.movementAxis = &m.cursor.X
	m.orthoAxis = &m.cursor.Y

	return &m
}

// Init implements the Bubble Tea interface for initialization.
func (m CrosswordModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress events and updates the model state accordingly.
func (m *CrosswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input
		switch msg.String() {
		// Game control keys
		case "ctrl+r":
			m.handleReset()

		case "ctrl+l":
			m.handleCheckLetter(m.cursor.X, m.cursor.Y)

		case "ctrl+w":
			m.handleCheckWord()

		case "ctrl+p":
			m.handleCheckPuzzle()

		case "ctrl+a":
			// Toggle auto-check mode
			m.autoCheck = !m.autoCheck
			if m.autoCheck {
				m.handleCheckPuzzle() // Check all entries when enabling
			} else {
				m.incorrect = [15][15]bool{} // Clear incorrect markers when disabling
			}

		// Navigation keys
		case "up":
			if m.movementAxis == &m.cursor.X {
				// If currently moving horizontally, switch to vertical movement
				m.switchAxis(&m.cursor.Y, &m.cursor.X, false)
			} else {
				// Move up in the grid
				m.handleMoveBackward()
			}

		case "down":
			if m.movementAxis == &m.cursor.X {
				// If currently moving horizontally, switch to vertical movement
				m.switchAxis(&m.cursor.Y, &m.cursor.X, false)
			} else {
				// Move down in the grid
				m.handleMoveForward()
			}

		case "left":
			if m.movementAxis == &m.cursor.Y {
				// If currently moving vertically, switch to horizontal movement
				m.switchAxis(&m.cursor.X, &m.cursor.Y, true)
			} else {
				// Move left in the grid
				m.handleMoveBackward()
			}

		case "right":
			if m.movementAxis == &m.cursor.Y {
				// If currently moving vertically, switch to horizontal movement
				m.switchAxis(&m.cursor.X, &m.cursor.Y, true)
			} else {
				// Move right in the grid
				m.handleMoveForward()
			}

		case " ":
			// Space advances to next cell if possible
			if *m.movementAxis < 14 && m.cellAt(1) != '.' {
				*m.movementAxis++
			}

		case "tab", "enter":
			m.handleNextWord()

		case "shift+tab", "shift+enter":
			m.handlePrevWord()

		case "backspace":
			m.handleDelete()

		default:
			m.handleInput(msg)
		}
	}

	// Update the current clue based on cursor position
	m.clue = m.clueAt(m.cursor.X, m.cursor.Y)

	return m, nil
}

// handleReset resets the puzzle to its initial state.
func (m *CrosswordModel) handleReset() {
	// Reset counters
	m.correctCount = 0
	m.filledCount = 0

	// Clear the grid
	for i := range m.grid {
		m.grid[i] = [15]byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	}

	// Restore black cells and update counters
	for i := range 225 {
		row, col := i/15, i%15

		// Preserve black cells from the answer grid
		if m.answer[row][col] == '.' {
			m.grid[row][col] = '.'
		}

		// Update counters for filled and correct cells
		if m.grid[row][col] == m.answer[row][col] {
			m.correctCount++
		}
		if m.grid[row][col] != ' ' {
			m.filledCount++
		}
	}

	// Clear incorrect markers
	m.incorrect = [15][15]bool{}

	// Save the reset state
	m.SaveToFile()
}

// handleCheckLetter checks if the letter at the specified position is correct.
func (m *CrosswordModel) handleCheckLetter(x, y int) {
	// Only mark as incorrect if the cell is filled and doesn't match the answer
	if m.grid[y][x] != m.answer[y][x] && m.grid[y][x] != ' ' {
		m.incorrect[y][x] = true
	}
}

// handleCheckWord checks all letters in the current word.
func (m *CrosswordModel) handleCheckWord() {
	// Find the start of the word
	offset := 0

	// Move backward until we hit the start of the word or the grid edge
	for pos := *m.movementAxis + offset; pos > 0 && m.cellAt(offset) != '.'; {
		offset--
		pos = *m.movementAxis + offset
	}

	// If we hit a black cell, move forward one to get to the start of the word
	if m.cellAt(offset) == '.' {
		offset++
	}

	// Check each letter in the word
	for *m.movementAxis+offset < 15 && m.cellAt(offset) != '.' {
		x := m.cursor.X
		y := m.cursor.Y

		// Adjust coordinates based on direction
		if m.isAcross {
			x = m.cursor.X + offset
		} else {
			y = m.cursor.Y + offset
		}

		m.handleCheckLetter(x, y)
		offset++
	}
}

// handleCheckPuzzle checks all letters in the entire puzzle.
func (m *CrosswordModel) handleCheckPuzzle() {
	for i := range 225 {
		x := i % 15
		y := i / 15
		m.handleCheckLetter(x, y)
	}
}

// incrementCursor moves the cursor forward along the current axis.
func (m *CrosswordModel) incrementCursor() {
	*m.movementAxis++

	// If we've reached the end of the grid, wrap around
	if *m.movementAxis > 14 {
		*m.movementAxis = 0

		// If we've reached the end of the puzzle, wrap around
		if *m.orthoAxis >= 14 {
			*m.orthoAxis = 0
		} else {
			*m.orthoAxis++
		}
	}
}

// decrementCursor moves the cursor backward along the current axis.
func (m *CrosswordModel) decrementCursor() {
	*m.movementAxis--

	// If we've reached the start of the grid, wrap around
	if *m.movementAxis < 0 {
		*m.movementAxis = 14

		// If we've reached the start of the puzzle, wrap around
		if *m.orthoAxis <= 0 {
			*m.orthoAxis = 14
		} else {
			*m.orthoAxis--
		}
	}
}

// handleMoveForward moves the cursor forward, skipping black cells.
func (m *CrosswordModel) handleMoveForward() {
	m.incrementCursor()

	// Skip black cells
	for m.cellAt(0) == '.' {
		m.incrementCursor()
	}
}

// handleMoveBackward moves the cursor backward, skipping black cells.
func (m *CrosswordModel) handleMoveBackward() {
	m.decrementCursor()

	// Skip black cells
	for m.cellAt(0) == '.' {
		m.decrementCursor()
	}
}

// handleNextWord moves to the next word in the current direction.
func (m *CrosswordModel) handleNextWord() {
	// Select the appropriate clue list based on direction
	clues := m.downClues
	if m.isAcross {
		clues = m.acrossClues
	}

	// Move to the next clue
	m.clue = (m.clue + 1) % len(clues)

	// Find the grid position for this clue
	for m.clueAt(m.cursor.X, m.cursor.Y) != m.clue {
		m.incrementCursor()
	}

	// If the puzzle is complete, don't try to find empty cells
	if m.filledCount == 225 {
		return
	}

	// Find the first empty cell in this word
	for m.grid[m.cursor.Y][m.cursor.X] != ' ' {
		*m.movementAxis++

		// If we hit the end of the word or a black cell, move to the next word
		if *m.movementAxis >= 15 || m.grid[m.cursor.Y][m.cursor.X] == '.' {
			*m.movementAxis--
			m.handleNextWord() // Recursively find the next word
			break
		}
	}
}

// handlePrevWord moves to the previous word in the current direction.
func (m *CrosswordModel) handlePrevWord() {
	// Select the appropriate clue list based on direction
	clues := m.downClues
	if m.isAcross {
		clues = m.acrossClues
	}

	// Move to the previous clue
	m.clue = (m.clue - 1 + len(clues)) % len(clues)

	// Find the grid position for this clue
	for m.clueAt(m.cursor.X, m.cursor.Y) != m.clue {
		m.decrementCursor()
	}

	// Move to the start of the word
	for *m.movementAxis > 0 && m.cellAt(-1) != '.' {
		*m.movementAxis--
	}

	// If the puzzle is complete, don't try to find empty cells
	if m.filledCount == 225 {
		return
	}

	// Find the first empty cell in this word
	for m.grid[m.cursor.Y][m.cursor.X] != ' ' {
		*m.movementAxis++

		// If we hit the end of the word or a black cell, move to the previous word
		if *m.movementAxis >= 15 || m.grid[m.cursor.Y][m.cursor.X] == '.' {
			*m.movementAxis--
			m.handlePrevWord() // Recursively find the previous word
			break
		}
	}
}

// handleDelete removes the letter at the current cursor position.
func (m *CrosswordModel) handleDelete() {
	// If current cell is empty and we're not at the start of a word, move back
	if m.grid[m.cursor.Y][m.cursor.X] == ' ' && *m.movementAxis != 0 {
		*m.movementAxis--
	}

	// If we hit a black cell, move forward
	if m.grid[m.cursor.Y][m.cursor.X] == '.' {
		*m.movementAxis++
	}

	// Update correctCount if we're deleting a correct letter
	if m.grid[m.cursor.Y][m.cursor.X] == m.answer[m.cursor.Y][m.cursor.X] {
		m.correctCount--
	}

	// Update filled count if we're deleting a letter
	if m.grid[m.cursor.Y][m.cursor.X] != ' ' {
		m.filledCount--
	}

	// Clear the cell and any incorrect marking
	m.grid[m.cursor.Y][m.cursor.X] = ' '
	m.incorrect[m.cursor.Y][m.cursor.X] = false

	// Get the clue indices for this cell
	acrossClue := m.clueIndices[m.cursor.Y][m.cursor.X].X
	downClue := m.clueIndices[m.cursor.Y][m.cursor.X].Y

	// Update the solved status for affected clues
	m.isAcrossSolved[acrossClue], m.isDownSolved[downClue] = m.isClueSolved()
}

// handleInput processes letter input from the keyboard.
func (m *CrosswordModel) handleInput(msg tea.KeyMsg) {
	// Only process single character inputs
	if len(msg.String()) != 1 {
		return
	}

	// Only accept letters
	input := rune(msg.String()[0])
	if !unicode.IsLetter(input) {
		return
	}

	// If the cell already has a letter, delete it first
	if m.grid[m.cursor.Y][m.cursor.X] != ' ' {
		m.handleDelete()
	}

	// Convert to uppercase
	input = unicode.ToUpper(input)

	// Add the letter to the grid
	m.grid[m.cursor.Y][m.cursor.X] = byte(input)
	m.incorrect[m.cursor.Y][m.cursor.X] = false

	// Get the clue indices for this cell
	acrossClue := m.clueIndices[m.cursor.Y][m.cursor.X].X
	downClue := m.clueIndices[m.cursor.Y][m.cursor.X].Y

	// Update the solved status for affected clues
	m.isAcrossSolved[acrossClue], m.isDownSolved[downClue] = m.isClueSolved()

	// Update counters
	m.filledCount++
	if m.grid[m.cursor.Y][m.cursor.X] == m.answer[m.cursor.Y][m.cursor.X] {
		m.correctCount++
	} else if m.autoCheck {
		// Mark as incorrect if auto-check is enabled
		m.incorrect[m.cursor.Y][m.cursor.X] = true
	}

	// Check for win condition
	if m.correctCount == 225 {
		m.message = "🎉 Congratulations! You solved the crossword! 🎉"
	}

	// Advance cursor if possible
	if *m.movementAxis < 14 && m.cellAt(1) != '.' {
		*m.movementAxis++
	}
}

// isClueSolved checks if the across and down clues at the current position are solved.
func (m *CrosswordModel) isClueSolved() (bool, bool) {
	// Start at the current position
	x := m.cursor.X
	y := m.cursor.Y
	isAcrossSolved := true
	isDownSolved := true

	// Find the start of the across word
	for x > 0 && m.grid[m.cursor.Y][x-1] != '.' {
		x--
	}

	// Find the start of the down word
	for y > 0 && m.grid[y-1][m.cursor.X] != '.' {
		y--
	}

	// Check if the across word is completely filled
	for x < 15 && m.grid[m.cursor.Y][x] != '.' {
		if m.grid[m.cursor.Y][x] == ' ' {
			isAcrossSolved = false
			break
		}
		x++
	}

	// Check if the down word is completely filled
	for y < 15 && m.grid[y][m.cursor.X] != '.' {
		if m.grid[y][m.cursor.X] == ' ' {
			isDownSolved = false
			break
		}
		y++
	}

	return isAcrossSolved, isDownSolved
}

// switchAxis changes the current movement direction.
func (m *CrosswordModel) switchAxis(newAxis, newOrtho *int, isAcross bool) {
	m.movementAxis = newAxis
	m.orthoAxis = newOrtho
	m.isAcross = isAcross
}

// cellAt returns the cell value at an offset from the current cursor position.
func (m *CrosswordModel) cellAt(offset int) byte {
	if m.movementAxis == &m.cursor.X {
		// Moving horizontally
		return m.grid[m.cursor.Y][m.cursor.X+offset]
	}
	// Moving vertically
	return m.grid[m.cursor.Y+offset][m.cursor.X]
}

// clueAt returns the clue index for the cell at the specified position.
func (m *CrosswordModel) clueAt(x, y int) int {
	if m.isAcross {
		return m.clueIndices[y][x].X
	}
	return m.clueIndices[y][x].Y
}
