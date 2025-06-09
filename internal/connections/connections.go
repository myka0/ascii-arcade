package connections

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
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
	date              string
	wordGroups        [4]WordGroup
	board             [16]string
	selectedTiles     []string
	guessHistory      [][4]string
	mistakesRemaining int
	message           string
}

// InitConnectionsModel initializes a new connections model.
func InitConnectionsModel() *ConnectionsModel {
	date, err := GetLatestDate()
	if err != nil {
		fmt.Println("Failed to get latest date:", err)
	}

	m, err := LoadFromFile(date)
	if err != nil {
		fmt.Println("Failed to load crossword:", err)
	}

	// Initialize the board
	board := []string{}
	for _, wordGroup := range m.wordGroups {
		board = append(board, wordGroup.Members[:]...)
	}
	copy(m.board[:], board)

	m.shuffle()

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
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			m.shuffle()
		}

	// Handle mouse input
	case tea.MouseMsg:
		// Only respond to left clicks and if the game is still active
		if msg.Action != tea.MouseActionRelease ||
			msg.Button != tea.MouseButtonLeft ||
			m.mistakesRemaining == 0 {
			return m, nil
		}

		// Check if a word was clicked
		for _, word := range m.board {
			if zone.Get(word).InBounds(msg) {
				// If the word is already selected, deselect it
				if i := slices.Index(m.selectedTiles, word); i != -1 {
					m.selectedTiles = slices.Delete(m.selectedTiles, i, i+1)

					// Otherwise, add it to the selection
				} else if len(m.selectedTiles) < 4 {
					m.selectedTiles = append(m.selectedTiles, word)
				}

				return m, nil
			}
		}
	}

	return m, nil
}

// shuffle shuffles the board randomly.
func (m *ConnectionsModel) shuffle() {
	// Use the maphash package to generate a random seed
	generator := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	generator.Shuffle(len(m.board), func(i, j int) {
		m.board[i], m.board[j] = m.board[j], m.board[i]
	})
}
