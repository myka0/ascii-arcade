package tetris

import (
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris/modes/single"
)

const (
	gameLevel     = 1
	maxLevel      = 15
	nextQueueSize = 5
)

// tickMsg is the internal message used to drive gravity.
type tickMsg struct {
	seq int
}

// TetrisModel represents the state of a Tetris game.
type TetrisModel struct {
	game *single.Game
	rng  *rand.Rand

	paused       bool
	gameOver     bool
	hasStarted   bool
	softDrop     bool
	tickSeq      int
	fallInterval time.Duration

	startTime      time.Time
	pausedAt       time.Time
	pausedDuration time.Duration
}

// InitTetrisModel creates a new Tetris game model with a freshly seeded random source.
func InitTetrisModel() *TetrisModel {
	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	g, _ := single.NewGame(&single.Input{
		Level:         gameLevel,
		MaxLevel:      maxLevel,
		IncreaseLevel: true,
		GhostEnabled:  true,
		Rand:          rng,
	})

	m := &TetrisModel{
		game:           g,
		rng:            rng,
		tickSeq:        1,
		startTime:      time.Now(),
		pausedAt:       time.Time{},
		pausedDuration: 0,
	}

	return m
}

// Init implements the Bubble Tea interface for initialization.
func (m *TetrisModel) Init() tea.Cmd {
	return nil
}

// Update handles ticks, keypress, and mouse events.
func (m *TetrisModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return m.handleTick(msg)
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)
	}

	return m, nil
}

// handleTick processes a gravity tick.
func (m *TetrisModel) handleTick(msg tickMsg) (tea.Model, tea.Cmd) {
	// Do not process stale ticks
	if msg.seq != m.tickSeq {
		return m, nil
	}

	// Do not process ticks when game shouldn't be running
	if !m.hasStarted || m.paused || m.gameOver {
		return m, nil
	}

	// Get the current piece so we can detect whether this tick caused the active piece to lock down
	curPiece := nextPieceValue(m.game)

	gameOver, _ := m.game.TickLower()
	if gameOver {
		m.gameOver = true
		return m, nil
	}

	// Stop soft drop if the piece just locked down
	lockedDown := nextPieceValue(m.game) != curPiece
	if lockedDown && m.softDrop {
		m.game.ToggleSoftDrop()
		m.softDrop = false
	}

	// Refresh the cached fall interval from the engine
	if !m.softDrop {
		m.fallInterval = m.game.GetDefaultFallInterval()
	} else {
		m.fallInterval = m.game.GetDefaultFallInterval() / 15
	}

	return m, m.scheduleTick()
}

// handleKeyPress handles keyboard input.
func (m *TetrisModel) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Always allow ctrl+r to reset
	if key == "ctrl+r" {
		return InitTetrisModel(), nil
	}

	// Game over keybinds
	if m.gameOver {
		if key == "enter" {
			return m, func() tea.Msg { return "home" }
		}
		return m, nil
	}

	// Pause toggle must work while paused and playing
	if key == "p" || key == "esc" {
		return m.togglePause()
	}

	// Ignore gameplay keys while paused
	if m.paused {
		return m, nil
	}

	switch key {
	case "left", "a":
		m.game.MoveLeft()
		return m, m.startIfNeeded()
	case "right", "d":
		m.game.MoveRight()
		return m, m.startIfNeeded()
	case "down", "s":
		return m.handleSoftDrop()
	case "up", "w", "space":
		return m.handleHardDrop()
	case "x":
		_ = m.game.Rotate(true)
		return m, m.startIfNeeded()
	case "z":
		_ = m.game.Rotate(false)
		return m, m.startIfNeeded()
	case "c":
		m.gameOver, _ = m.game.Hold()
		return m, m.startIfNeeded()
	}

	return m, nil
}

// handleMouseClick handles mouse interactions.
func (m *TetrisModel) handleMouseClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only respond to left clicks
	if msg.Mouse().Button != tea.MouseLeft {
		return m, nil
	}

	// Handle game over UI
	if m.gameOver {
		switch {
		case zone.Get("reset").InBounds(msg):
			return InitTetrisModel(), nil
		case zone.Get("exit").InBounds(msg):
			return m, func() tea.Msg { return "home" }
		}
	}

	return m, nil
}

// handleSoftDrop toggles soft drop and sets the new interval to use.
func (m *TetrisModel) handleSoftDrop() (tea.Model, tea.Cmd) {
	m.fallInterval = m.game.ToggleSoftDrop()
	m.softDrop = !m.softDrop

	// Start the game if this is the first move
	if !m.hasStarted {
		m.hasStarted = true
		m.startTime = time.Now()
	}

	// Trigger an immediate tick so the piece responds right away
	m.tickSeq++
	return m, m.scheduleTick()
}

// handleHardDrop performs a hard drop and checks for game over.
func (m *TetrisModel) handleHardDrop() (tea.Model, tea.Cmd) {
	if gameOver, _ := m.game.HardDrop(); gameOver {
		m.gameOver = true
		return m, nil
	}

	m.fallInterval = m.game.GetDefaultFallInterval()
	m.softDrop = false

	// Start the game if this is the first move
	if !m.hasStarted {
		m.hasStarted = true
		m.startTime = time.Now()
	}

	// Trigger an immediate tick so the piece responds right away
	m.tickSeq++
	return m, m.scheduleTick()
}

// scheduleTick returns a tea.Cmd that fires a tickMsg after the current fall interval has elapsed.
// The tick is tagged with the current tickSeq so any older ticks can be invalidated by bumping the counter.
func (m *TetrisModel) scheduleTick() tea.Cmd {
	seq := m.tickSeq
	interval := m.fallInterval
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return tickMsg{seq: seq}
	})
}

// nextPieceValue returns the value of the piece at the head of the engine's piece queue.
func nextPieceValue(g *single.Game) byte {
	bag := g.GetBagTetriminos()
	if len(bag) == 0 {
		return 0
	}
	return bag[0].Value
}

// togglePause flips the paused state.
func (m *TetrisModel) togglePause() (tea.Model, tea.Cmd) {
	m.paused = !m.paused

	// Bump tickSeq so any in-flight tick is dropped
	m.tickSeq++

	// Set pause time
	if m.paused {
		m.pausedAt = time.Now()
		return m, nil
	}

	// Add the pause span to the pause running total
	m.pausedDuration += time.Since(m.pausedAt)
	m.pausedAt = time.Time{}

	// Schedule a fresh tick on resume
	return m, m.scheduleTick()
}

// startIfNeeded schedules the first tick when the player makes their first move.
func (m *TetrisModel) startIfNeeded() tea.Cmd {
	if !m.hasStarted {
		m.hasStarted = true
		m.tickSeq++
		m.startTime = time.Now()

		return m.scheduleTick()
	}

	return nil
}
