package connections

import (
	"fmt"
	"hash/maphash"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
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
	selectedTiles     [4]string
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

// Update handles keypress events and updates the model state accordingly.
func (m *ConnectionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			m.shuffle()
		}
	}

	return m, nil
}

// isSelected checks if the specified cell is selected.
func (m ConnectionsModel) isSelected(word string) bool {
	for _, selectedWord := range m.selectedTiles {
		if selectedWord == word {
			return true
		}
	}

	return false
}

// shuffle shuffles the board randomly.
func (m *ConnectionsModel) shuffle() {
	// Use the maphash package to generate a random seed
	generator := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
	generator.Shuffle(len(m.board), func(i, j int) {
		m.board[i], m.board[j] = m.board[j], m.board[i]
	})
}
