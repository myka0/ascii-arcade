package minesweeper

import (
	"fmt"
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

const (
	Beginner = iota
	Intermediate
	Expert
)

// Difficulty defines the grid size and mine count for a game mode.
type Difficulty struct {
	Name       string
	Rows, Cols int
	Mines      int
}

// Difficulties contains the available game modes.
var Difficulties = []Difficulty{
	{"Beginner", 9, 9, 10},
	{"Intermediate", 16, 16, 40},
	{"Expert", 16, 30, 99},
}

// timerTickMsg is the internal message used to refresh the game timer display.
type timerTickMsg struct {
	seq int
}

// Cell represents a single tile on the board.
type Cell struct {
	mine     bool
	revealed bool
	flagged  bool
	adjacent uint8
}

// Position represents a coordinate in the 2D grid.
type Position struct {
	X, Y int
}

// MinesweeperModel represents the state of the Minesweeper game.
type MinesweeperModel struct {
	difficulty   Difficulty
	board        [][]Cell
	cursor       Position
	mineHit      Position
	flagsPlaced  int
	revealedSafe int
	message      string
	gameOver     bool
	hasStarted   bool
	hasSelected  bool
	startTime    time.Time
	endTime      time.Time
	timerSeq     int
}

// InitMinesweeperModel creates and initializes a new minesweeper model.
func InitMinesweeperModel() *MinesweeperModel {
	return &MinesweeperModel{
		message: "Select a difficulty.",
	}
}

// createGame initializes a new game with the given difficulty settings.
func createGame(difficulty Difficulty) *MinesweeperModel {
	m := &MinesweeperModel{
		difficulty:   difficulty,
		board:        newBoard(difficulty),
		cursor:       Position{X: difficulty.Cols / 2, Y: difficulty.Rows / 2},
		mineHit:      Position{X: -1, Y: -1},
		flagsPlaced:  0,
		revealedSafe: 0,
		message:      "Select any cell to begin.",
		gameOver:     false,
		hasStarted:   false,
		hasSelected:  true,
	}

	return m
}

// newBoard initializes a new game board based on the given difficulty.
func newBoard(difficulty Difficulty) [][]Cell {
	board := make([][]Cell, difficulty.Rows)
	for y := range board {
		board[y] = make([]Cell, difficulty.Cols)
		for x := range difficulty.Cols {
			board[y][x] = Cell{mine: false, revealed: false, flagged: false, adjacent: 0}
		}
	}

	return board
}

// Init implements the Bubble Tea interface for initialization.
func (m *MinesweeperModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress and mouse events.
func (m *MinesweeperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timerTickMsg:
		return m.handleTimerTick(msg)
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)
	}

	return m, nil
}

// handleTimerTick processes a timer tick, triggering a re-render and scheduling the next tick.
func (m *MinesweeperModel) handleTimerTick(msg timerTickMsg) (tea.Model, tea.Cmd) {
	if msg.seq != m.timerSeq || m.gameOver || !m.hasStarted {
		return m, nil
	}
	return m, m.scheduleTimerTick()
}

// scheduleTimerTick returns a tea.Cmd that fires a timerTickMsg after one second.
func (m *MinesweeperModel) scheduleTimerTick() tea.Cmd {
	seq := m.timerSeq
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return timerTickMsg{seq: seq}
	})
}

// handleKeyPress handles keyboard input.
func (m *MinesweeperModel) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Always allow resetting and changing game
	switch key {
	case "ctrl+r":
		return createGame(m.difficulty), nil
	case "1":
		return createGame(Difficulties[Beginner]), nil
	case "2":
		return createGame(Difficulties[Intermediate]), nil
	case "3":
		return createGame(Difficulties[Expert]), nil
	}

	// Disable all other keypresses on game over screen
	if m.gameOver {
		return m, nil
	}

	// Difficulty selection menu
	if !m.hasSelected {
		switch key {
		case "up", "w":
			m.cursor.Y = (m.cursor.Y - 1 + len(Difficulties)) % len(Difficulties)
		case "down", "s":
			m.cursor.Y = (m.cursor.Y + 1) % len(Difficulties)
		case "enter":
			return createGame(Difficulties[m.cursor.Y]), nil
		}
		return m, nil
	}

	switch key {
	case "up", "w":
		if m.cursor.Y > 0 {
			m.cursor.Y--
		}
	case "down", "s":
		if m.cursor.Y < m.difficulty.Rows-1 {
			m.cursor.Y++
		}
	case "left", "a":
		if m.cursor.X > 0 {
			m.cursor.X--
		}
	case "right", "d":
		if m.cursor.X < m.difficulty.Cols-1 {
			m.cursor.X++
		}
	case "space", "enter":
		return m.handleRevealCell(m.cursor)
	case "f":
		return m.handleToggleFlag(m.cursor)
	case "c":
		return m.handleRevealChord(m.cursor)
	}

	return m, nil
}

// handleMouseClick handles mouse interactions.
func (m *MinesweeperModel) handleMouseClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only respond to left, right, or middle clicks
	if msg.Mouse().Button != tea.MouseLeft && msg.Mouse().Button != tea.MouseRight && msg.Mouse().Button != tea.MouseMiddle {
		return m, nil
	}

	// Handle game over buttons
	if m.gameOver {
		switch {
		case zone.Get("reset").InBounds(msg):
			return createGame(m.difficulty), nil
		case zone.Get("exit").InBounds(msg):
			return m, func() tea.Msg { return "home" }
		default:
			return m, nil
		}
	}

	// Check if a board intersection was clicked
	for y := 0; y < m.difficulty.Rows; y++ {
		for x := 0; x < m.difficulty.Cols; x++ {
			label := fmt.Sprintf("%d_%d", x, y)
			if !zone.Get(label).InBounds(msg) {
				continue
			}

			m.cursor = Position{X: x, Y: y}

			if msg.Mouse().Button == tea.MouseRight {
				return m.handleToggleFlag(m.cursor)
			}
			if msg.Mouse().Button == tea.MouseMiddle {
				return m.handleRevealChord(m.cursor)
			}
			return m.handleRevealCell(m.cursor)
		}
	}

	return m, nil
}

// handleRevealCell reveals the cell at p.
func (m *MinesweeperModel) handleRevealCell(p Position) (tea.Model, tea.Cmd) {
	// First reveal seeds the mine layout and starts the clock
	startingTimer := false
	if !m.hasStarted {
		m.placeMines(p)
		m.computeAdjacency()
		m.hasStarted = true
		m.startTime = time.Now()
		m.message = ""
		startingTimer = true
	}

	cell := &m.board[p.Y][p.X]
	if cell.revealed || cell.flagged {
		return m, nil
	}

	// Hitting a mine reveals it and ends the game
	if cell.mine {
		cell.revealed = true
		m.mineHit = p
		m.gameOver = true
		m.endTime = time.Now()
		m.timerSeq++
		m.message = "Game Over! You hit a mine."
		m.forEachCell(func(x, y int, c *Cell) {
			if c.mine && !c.flagged {
				c.revealed = true
			}
		})
		return m, nil
	}

	m.floodReveal(p)

	// Check win
	safeCells := m.difficulty.Rows*m.difficulty.Cols - m.difficulty.Mines
	if m.revealedSafe < safeCells {
		if startingTimer {
			return m, m.scheduleTimerTick()
		}
		return m, nil
	}

	m.gameOver = true
	m.timerSeq++
	m.endTime = time.Now()
	m.message = "You win!"
	m.flagsPlaced = m.difficulty.Mines

	// Flag all remaining unflagged mines
	m.forEachCell(func(x, y int, c *Cell) {
		if c.mine && !c.flagged {
			c.flagged = true
		}
	})

	return m, nil
}

// handleToggleFlag toggles a flag on the cell at p.
func (m *MinesweeperModel) handleToggleFlag(p Position) (tea.Model, tea.Cmd) {
	if !m.hasStarted {
		return m, nil
	}

	cell := &m.board[p.Y][p.X]
	if cell.revealed {
		return m, nil
	}

	// Removing a flag is always allowed
	if cell.flagged {
		cell.flagged = false
		m.flagsPlaced--
		m.message = ""
		return m, nil
	}

	// Placing a flag is gated by the flag budget
	if m.flagsPlaced >= m.difficulty.Mines {
		m.message = "All flags placed. Unflag a cell to place a different flag."
		return m, nil
	}

	cell.flagged = true
	m.flagsPlaced++
	m.message = ""
	return m, nil
}

// handleRevealChord performs the chord action at p.
// https://minesweeper.fandom.com/wiki/Chording
func (m *MinesweeperModel) handleRevealChord(p Position) (tea.Model, tea.Cmd) {
	if !m.hasStarted {
		return m, nil
	}

	cell := m.board[p.Y][p.X]
	if !cell.revealed || cell.adjacent == 0 {
		return m, nil
	}

	// Chord only fires when flagged neighbors exactly match the cell's number
	var flagged uint8
	m.forEachNeighbor(p, func(n Position) {
		if m.board[n.Y][n.X].flagged {
			flagged++
		}
	})
	if flagged != cell.adjacent {
		return m, nil
	}

	// Collect neighbors first
	var targets []Position
	m.forEachNeighbor(p, func(n Position) {
		nc := m.board[n.Y][n.X]
		if !nc.flagged && !nc.revealed {
			targets = append(targets, n)
		}
	})

	// Reveal each neighbor
	for _, t := range targets {
		m.handleRevealCell(t)
		if m.gameOver {
			return m, nil
		}
	}

	return m, nil
}

// floodReveal performs an iterative BFS reveal starting at p.
func (m *MinesweeperModel) floodReveal(start Position) {
	queue := []Position{start}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		cell := &m.board[p.Y][p.X]
		if cell.revealed {
			continue
		}

		cell.revealed = true
		m.revealedSafe++

		// Numbered cells terminate the flood
		if cell.adjacent != 0 {
			continue
		}

		m.forEachNeighbor(p, func(n Position) {
			nc := m.board[n.Y][n.X]
			if !nc.revealed && !nc.flagged && !nc.mine {
				queue = append(queue, n)
			}
		})
	}
}

// forEachNeighbor invokes fn for each in-bounds cell adjacent to p (8-way).
func (m *MinesweeperModel) forEachNeighbor(p Position, fn func(Position)) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			n := Position{X: p.X + dx, Y: p.Y + dy}
			if n.Y >= 0 && n.Y < m.difficulty.Rows && n.X >= 0 && n.X < m.difficulty.Cols {
				fn(n)
			}
		}
	}
}

// forEachCell invokes fn for every cell on the board.
func (m *MinesweeperModel) forEachCell(fn func(x, y int, c *Cell)) {
	for y := range m.board {
		for x := range m.board[y] {
			fn(x, y, &m.board[y][x])
		}
	}
}

// placeMines randomly distributes mines across the board, guaranteeing that
// the cell at avoid and its 8 neighbors are mine-free.
func (m *MinesweeperModel) placeMines(avoid Position) {
	rows := m.difficulty.Rows
	cols := m.difficulty.Cols
	totalCells := rows * cols

	// Build the safe set
	safe := make(map[int]struct{}, 9)
	safe[avoid.Y*cols+avoid.X] = struct{}{}
	m.forEachNeighbor(avoid, func(n Position) {
		safe[n.Y*cols+n.X] = struct{}{}
	})

	// Build the pool of candidate cell indices
	candidates := make([]int, 0, totalCells-len(safe))
	for i := range totalCells {
		if _, isSafe := safe[i]; isSafe {
			continue
		}
		candidates = append(candidates, i)
	}

	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	// Partial Fisher-Yates: pick mines distinct indices from candidates
	for i := range m.difficulty.Mines {
		j := i + rng.IntN(len(candidates)-i)
		candidates[i], candidates[j] = candidates[j], candidates[i]
		idx := candidates[i]
		m.board[idx/cols][idx%cols].mine = true
	}
}

// computeAdjacency fills in the adjacent mine count for every cell.
func (m *MinesweeperModel) computeAdjacency() {
	m.forEachCell(func(x, y int, c *Cell) {
		var count uint8
		m.forEachNeighbor(Position{X: x, Y: y}, func(n Position) {
			if m.board[n.Y][n.X].mine {
				count++
			}
		})
		c.adjacent = count
	})
}