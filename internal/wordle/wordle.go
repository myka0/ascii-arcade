package wordle

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"unicode"
)

// Possible states for each key on the keyboard.
const (
	keyUntried = 0
	keyAbsent  = 1
	keyPresent = 2
	keyCorrect = 3
)

// WordleModel represents the state of a Wordle game.
type WordleModel struct {
	date     string
	answer   [5]byte
	guesses  [6][5]byte
	cursorX  int
	cursorY  int
	keyboard map[byte]int
	message  string
}

// InitWordleModel creates and initializes a new wordle model.
// It loads puzzle data from file and sets up the initial game state.
func InitWordleModel() *WordleModel {
	m, err := LoadFromFile()
	if err != nil {
		m.message = err.Error()
	}

	return &m
}

// Init implements the Bubble Tea interface for initialization.
func (m WordleModel) Init() tea.Cmd {
	return nil
}

// Update handles keypress events and updates the model state accordingly.
func (m *WordleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			m.handleReset()

		case "backspace":
			m.handleDelete()

		case "enter":
			m.handleSubmit()

		default:
			m.handleInput(msg)
		}
	}

	return m, nil
}

// handleReset resets the game board and keyboard to initial state.
func (m *WordleModel) handleReset() {
	m.cursorX = 0
	m.cursorY = 0

	// Reset all guesses to empty space
	for i := range m.guesses {
		m.guesses[i] = [5]byte{' ', ' ', ' ', ' ', ' '}
	}

	// Reset keyboard states to keyUntried for all letters
	for c := 'A'; c <= 'Z'; c++ {
		m.keyboard[byte(c)] = keyUntried
	}

	m.SaveToFile()
}

// handleDelete removes the last entered letter in the current guess.
func (m *WordleModel) handleDelete() {
	// Only delete if the cursor is not at the beginning of the row
	if m.cursorX > 0 {
		m.cursorX--
		m.guesses[m.cursorY][m.cursorX] = ' '
		m.message = ""
	}
}

// handleSubmit validates and processes the current guess.
func (m *WordleModel) handleSubmit() {
	// Ensure the guess is 5 letters
	if m.cursorX < 5 {
		m.message = "❌ Not enough letters."
		return
	}

	// Check if the current guess is a valid word from the word list
	isValid, err := isValid(m.guesses[m.cursorY])
	if err != nil {
		m.message = err.Error()
	}

	// Display a message if the word is not in the word list
	if !isValid {
		m.message = "❌ Not in word list."
		return
	}

	// Update keyboard state and move to the next row
	m.updateKeyStates()
	m.cursorY++
	m.cursorX = 0

	// Check if the guess is correct
	if m.guesses[m.cursorY-1] == m.answer {
		m.message = "🎉 Congratulations! You guessed the word! 🎉"
		m.cursorY = 6
		return
	}

	// If all guesses have been used, end the game
	if m.cursorY == 6 {
		m.message = fmt.Sprintf("❌ Game Over! The word was \"%s\" ❌", string(m.answer[:]))
	}
}

// handleInput processes a letter key input.
func (m *WordleModel) handleInput(msg tea.KeyMsg) {
	// Ensure the input is a single character and that we are within bounds
	if len(msg.String()) != 1 || m.cursorX >= 5 || m.cursorY >= 6 {
		return
	}

	// Extract the character from the input message
	input := rune(msg.String()[0])

	// Process the input if it is a letter
	if unicode.IsLetter(input) {
		// Convert to uppercase if it's a lowercase letter
		input = unicode.ToUpper(input)

		// Store the input and move the cursor
		m.guesses[m.cursorY][m.cursorX] = byte(input)
		m.cursorX++
	}
}

// updateKeyStates updates the keyboard based on the most recent guess.
func (m WordleModel) updateKeyStates() {
	currentGuess := m.guesses[m.cursorY]
	for i, char := range currentGuess {
		switch {
		// Character is correct
		case char == m.answer[i]:
			m.keyboard[char] = keyCorrect

		// Character is present in word but at wrong position
		case findIndex(m.answer, char) != -1:
			if m.keyboard[char] < keyPresent {
				m.keyboard[char] = keyPresent
			}

		// Character is not in answer
		default:
			if m.keyboard[char] < keyAbsent {
				m.keyboard[char] = keyAbsent
			}
		}
	}
}
