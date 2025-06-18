package wordle

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/v2"
)

// View renders the entire game UI.
func (m WordleModel) View() string {
	// Generate each row of the Wordle grid
	var rows [6]string
	for y := range m.guesses {
		rows[y] = m.viewGridRow(y)
	}

	// Vertically join all rows with center alignment and compose the full view
	joindedRows := lipgloss.JoinVertical(lipgloss.Center, rows[:]...)
	output := lipgloss.JoinVertical(
		lipgloss.Center,
		joindedRows,
		m.viewKeyboard(),
		"\n"+m.message+"\n",
	)

	return FGText.Render(output)
}

// viewGridRow renders a single row of the Wordle grid based on its position.
func (m WordleModel) viewGridRow(y int) string {
	// Initialize keyStates to keyAbsent
	keyStates := [5]int{1, 1, 1, 1, 1}
	guess := m.guesses[y]
	answer := m.answer

	var cells [5]string
	for i, letter := range guess {
		// Mark letters after cursor as keyUntried
		if y >= m.cursorY {
			keyStates[i] = keyUntried
			answer[i] = 0

			// If the letter matches the answer at this position mark as keyCorrect
		} else if letter == answer[i] {
			keyStates[i] = keyCorrect
			answer[i] = 0

			// If the letter is found at a different position mark as keyPresent
		} else if foundIdx := findIndex(answer, letter); foundIdx != -1 {
			keyStates[i] = keyPresent
			answer[foundIdx] = 0
		}

		// Style the cell
		cellContent := fmt.Sprintf("%c", letter)
		cellStyle := m.styleCell(keyStates[i])
		cells[i] = cellStyle.Render(Border.Render(cellContent))
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, cells[:]...)
}

// viewKeyboard renders the on screen keyboard with styling.
func (m WordleModel) viewKeyboard() string {
	return Border.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.viewKeyboardRow([]byte{'Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'}),
			m.viewKeyboardRow([]byte{'A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'}),
			m.viewKeyboardRow([]byte{'Z', 'X', 'C', 'V', 'B', 'N', 'M'}),
		),
	)
}

// viewKeyboardRow renders a row of keys with their appropriate styles.
func (m WordleModel) viewKeyboardRow(letters []byte) string {
	keys := make([]string, len(letters))

	// Style each key in the keyboard row
	for i, key := range letters {
		cell := Border.Render(fmt.Sprintf("%c", key))
		keys[i] = m.styleCell(m.keyboard[key]).Render(cell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, keys[:]...)
}

// styleCell returns a style object based on the key state.
func (m WordleModel) styleCell(keyStyle int) lipgloss.Style {
	switch keyStyle {
	case keyAbsent:
		return FGKeyAbsent
	case keyPresent:
		return FGKeyPresent
	case keyCorrect:
		return FGKeyCorrect
	default:
		return FGText
	}
}

// findIndex searches for a character in a 5 letter word slice and returns its index or -1 if not found
func findIndex(word [5]byte, char byte) int {
	for i, c := range word {
		if c == char {
			return i
		}
	}

	return -1
}
