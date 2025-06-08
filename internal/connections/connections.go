package connections

import (
	"hash/maphash"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WordGroup represents a group of words that are connected to each other.
type WordGroup struct {
	members    [4]string
	clue       string
	color      int
	isRevealed bool
}

// ConnectionsModel represents the state of the connections game.
type ConnectionsModel struct {
	wordGroups        [4]WordGroup
	board             [16]string
	selectedTiles     [4]string
	guessHistory      [][4]string
	mistakesRemaining int
	message           string
}

// InitConnectionsModel initializes a new connections model.
func InitConnectionsModel() *ConnectionsModel {
	// Initialize the word groups
	// TODO: Load from file
	wordGroups := [4]WordGroup{
		{
			members:    [4]string{"BEAR", "BULL", "DOVE", "HAWK"},
			clue:       "Animal metaphors in economics",
			color:      1,
			isRevealed: false,
		},
		{
			members:    [4]string{"HOLD", "LAST", "STAND", "STAY"},
			clue:       "Persist",
			color:      2,
			isRevealed: false,
		},
		{
			members:    [4]string{"BORN", "EDUCATION", "OCCUPATION", "SPOUSE"},
			clue:       "Sidebar info on a person’s Wikipedia page",
			color:      3,
			isRevealed: false,
		},
		{
			members:    [4]string{"BRED", "CACHE", "DOE", "LOOT"},
			clue:       "Homophones of slang for money",
			color:      4,
			isRevealed: false,
		},
	}

	// Initialize the model
	m := ConnectionsModel{
		wordGroups:        wordGroups,
		board:             [16]string{},
		selectedTiles:     [4]string{"BEAR", "BULL", "DOVE", "HAWK"},
		guessHistory:      [][4]string{},
		mistakesRemaining: 4,
		message:           "",
	}

	// Initialize the board
	board := []string{}
	for _, wordGroup := range wordGroups {
		board = append(board, wordGroup.members[:]...)
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

// View renders the connections game board.
func (m ConnectionsModel) View() string {
	// Render the board rows
	var rows [4]string
	for i := range 4 {
		rows[i] = m.viewBoardRow(i)
	}

	// Combine all elements vertically
	board := lipgloss.JoinVertical(lipgloss.Center, rows[:]...)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		board,
		m.message+"\n",
		"Create four groups of four!",
	)
}

// viewBoardRow renders a row of cells in the connections grid.
func (m ConnectionsModel) viewBoardRow(row int) string {
	const cellWidth = 16
	var cells [4]string

	start := row * 4

	// Render each cell in the row
	for col := range 4 {
		index := start + col
		word := m.board[index]

		// Determine cell style based on selection state
		style := NormalCell
		if m.isSelected(word) {
			style = SelectedCell
		}

		// Style content
		cells[col] = styleCell(word, cellWidth, style)
	}

	// Join all 4 cells horizontally
	return lipgloss.JoinHorizontal(lipgloss.Bottom, cells[:]...) + "\n"
}

// styleCell centers the text in the specified cell and styles it.
func styleCell(text string, width int, style lipgloss.Style) string {
	// Calculate padding
	padding := width - len(text)
	left := padding / 2
	right := padding - left

	// Create the margin style content
	margin := style.Render(strings.Repeat(" ", width))
	content := style.Render(strings.Repeat(" ", left) + text + strings.Repeat(" ", right))

	return lipgloss.JoinVertical(lipgloss.Left, margin, content, margin)
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
