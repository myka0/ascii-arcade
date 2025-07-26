package connections

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"slices"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// WordGroup represents a group of words that are connected to each other.
type WordGroup struct {
	Members    [4]string `json:"members"`
	Clue       string    `json:"clue"`
	Color      int       `json:"color"`
	IsRevealed bool      `json:"isRevealed"`
}

// ConnectionsModel represents the state of the connections game.
type ConnectionsModel struct {
	date               string
	wordGroups         [4]WordGroup
	board              [16]string
	selectedTiles      []string
	guessHistory       [][]string
	revealedWordGroups [][]string
	mistakesRemaining  int
	message            string
}

// InitConnectionsModel initializes a new connections model.
func InitConnectionsModel() *ConnectionsModel {
	date, err := GetLatestDate()
	if err != nil {
		fmt.Println("Failed to get latest date:", err)
	}

	m, err := LoadFromFile(date)
	if err != nil {
		fmt.Println("Failed to load connections:", err)
	}

	m.initBoard()

	return &m
}

// Init implements the Bubble Tea interface for initialization.
func (m ConnectionsModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress and mouse events to update the Connections game state.
func (m *ConnectionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle keyboard input
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			m.handleShuffle()
		case "ctrl+d":
			m.selectedTiles = []string{}
		case "enter":
			m.handleSubmit()
		case "ctrl+r":
			m.handleReset()
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
func (m *ConnectionsModel) handleMouseClick(msg tea.MouseMsg) {
	// Only respond to left clicks
	if msg.Mouse().Button != tea.MouseLeft {
		return
	}

	if m.mistakesRemaining == 0 {
		m.message = "Game over"
		return
	} else if m.isGameSolved() {
		m.message = "🎉 Congratulations! You win! 🎉"
		return
	} else {
		m.message = "Create four groups of four!"
	}

	// Handle button clicks
	switch {
	case zone.Get("Shuffle").InBounds(msg):
		m.handleShuffle()
		return

	case zone.Get("Deselect All").InBounds(msg):
		m.selectedTiles = nil
		return

	case zone.Get("Submit").InBounds(msg):
		m.handleSubmit()
		return
	}

	// Check if a word was clicked
	for _, word := range m.board {
		if zone.Get(word).InBounds(msg) {
			// If the word is already selected, deselect it
			if i := slices.Index(m.selectedTiles, word); i != -1 {
				m.selectedTiles = slices.Delete(m.selectedTiles, i, i+1)

			} else if len(m.selectedTiles) < 4 {
				// Otherwise, add it to the selection
				m.selectedTiles = append(m.selectedTiles, word)
			}

			return
		}
	}
}

// handleShuffle shuffles the board randomly.
func (m *ConnectionsModel) handleShuffle() {
	// Build a set of revealed wordGroups
	revealedSet := make(map[string]bool)
	for _, wordGroup := range m.revealedWordGroups {
		for _, word := range wordGroup {
			revealedSet[word] = true
		}
	}

	// Use the maphash package to generate a random seed
	generator := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	generator.Shuffle(len(m.board), func(i, j int) {
		// Don't shuffle revealed words
		if revealedSet[m.board[i]] || revealedSet[m.board[j]] {
			return
		}

		m.board[i], m.board[j] = m.board[j], m.board[i]
	})
}

// handleSubmit processes the current guess.
func (m *ConnectionsModel) handleSubmit() {
	// Check if the user selected the correct number of tiles
	if len(m.selectedTiles) != 4 {
		m.message = "Select four tiles"
		return
	}

	// Check for duplicate guesses
	for _, guess := range m.guessHistory {
		if stringSlicesEqual(guess, m.selectedTiles) {
			m.message = "Already guessed"
			return
		}
	}

	// Count number of selected words per word group
	var wordGroupCounts [4]int
	for _, guess := range m.selectedTiles {
		wordGroup := m.getWordGroup(guess)
		wordGroupCounts[wordGroup]++
	}

	// Find the word group with the most selected tiles
	largestCount := 0
	wordGroup := 0
	for i, count := range wordGroupCounts {
		if count > largestCount {
			largestCount = count
			wordGroup = i
		}
	}

	switch largestCount {
	case 4:
		// Correct guess
		m.wordGroups[wordGroup].IsRevealed = true
		m.revealedWordGroups = append(m.revealedWordGroups, m.wordGroups[wordGroup].Members[:])
		m.selectedTiles = []string{}
		m.initBoard()

		if m.isGameSolved() {
			m.message = "🎉 Congratulations! You win! 🎉"
		}
		return

	case 3:
		// One word away
		m.message = "One away..."
		m.mistakesRemaining--

	default:
		// Incorrect guess
		m.message = "Incorrect"
		m.mistakesRemaining--
	}

	// Record guess to history
	guess := make([]string, len(m.selectedTiles))
	copy(guess, m.selectedTiles)
	m.guessHistory = append(m.guessHistory, guess)

	// Check if the game is over
	if m.mistakesRemaining == 0 {
		m.message = "Game over"
		m.selectedTiles = []string{}
	}
}

// handleReset resets the game to the initial state.
func (m *ConnectionsModel) handleReset() {
	// Mark all word groups as unrevealed
	for i := range m.wordGroups {
		m.wordGroups[i].IsRevealed = false
	}

	// Reset the game state
	m.selectedTiles = []string{}
	m.guessHistory = [][]string{}
	m.revealedWordGroups = [][]string{}
	m.mistakesRemaining = 4
	m.message = "Create four groups of four!"

	m.initBoard()
	m.SaveToFile()
}

// initBoard initializes the board with the revealed words at the top and the unrevealed words below.
func (m *ConnectionsModel) initBoard() {
	m.message = "Create four groups of four!"
	board := []string{}

	// Initialize the board with revealed words at the top
	for _, wordGroup := range m.revealedWordGroups {
		board = append(board, wordGroup[:]...)
	}

	// Initialize the board with unrevealed words
	for _, wordGroup := range m.wordGroups {
		if !wordGroup.IsRevealed {
			board = append(board, wordGroup.Members[:]...)
		}
	}

	copy(m.board[:], board)
	m.handleShuffle()
}

// getWordGroup returns the index of the word group that contains the specified word.
func (m *ConnectionsModel) getWordGroup(word string) int {
	for i, wordGroup := range m.wordGroups {
		if slices.Contains(wordGroup.Members[:], word) {
			return i
		}
	}

	return 0
}

// isGameSolved returns true if all word groups are revealed.
func (m *ConnectionsModel) isGameSolved() bool {
	for _, wordGroup := range m.wordGroups {
		if !wordGroup.IsRevealed {
			return false
		}
	}

	return true
}

// stringSlicesEqual returns true if the two slices contain the same elements.
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
