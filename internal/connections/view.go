package connections

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone/v2"
)

// View renders the connections game board.
func (m ConnectionsModel) View() string {
	// Render the board rows
	var rows [4]string
	for i := range 4 {
		wordGroupIndex := m.getWordGroup(m.board[i*4])
		if m.wordGroups[wordGroupIndex].IsRevealed {
			rows[i] = m.viewRevealedRow(wordGroupIndex)
		} else {
			rows[i] = m.viewBoardRow(i)
		}
	}

	// Combine all elements vertically
	board := lipgloss.JoinVertical(lipgloss.Center, rows[:]...)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		FGLightText.Render(m.message)+"\n",
		board,
		m.viewMistakesRemaining(),
		viewButtonRow(),
	)
}

// viewRevealedRow renders a fully revealed group row with color and clue styling.
func (m ConnectionsModel) viewRevealedRow(row int) string {
	group := m.wordGroups[row]

	// Determine color based on group color number
	var color lipgloss.Color
	switch group.Color {
	case 1:
		color = Color1
	case 2:
		color = Color2
	case 3:
		color = Color3
	case 4:
		color = Color4
	}

	wordRow := strings.Join(group.Members[:], "  ") + "\n"
	return RevealedLine.Background(color).Render(wordRow+group.Clue) + "\n"
}

// viewBoardRow renders a row of cells in the connections grid.
func (m ConnectionsModel) viewBoardRow(row int) string {
	var cells [4]string

	start := row * 4

	// Render each cell in the row
	for col := range 4 {
		index := start + col
		word := m.board[index]

		// Determine cell style based on selection state
		style := NormalCell
		if slices.Contains(m.selectedTiles, word) {
			style = SelectedCell
		}

		// Style content
		cells[col] = zone.Mark(word, style.Render(word))
	}

	// Join all 4 cells horizontally
	return lipgloss.JoinHorizontal(lipgloss.Bottom, cells[:]...) + "\n"
}

// viewMistakesRemaining renders the number of mistakes remaining.
func (m ConnectionsModel) viewMistakesRemaining() string {
	mistakes := MistakeCell.Render(strings.Repeat("● ", m.mistakesRemaining))
	return FGLightText.Render("Mistakes Remaining: ") + mistakes
}

// viewButtonRow renders the buttons at the bottom of the board.
func viewButtonRow() string {
	return lipgloss.JoinHorizontal(lipgloss.Center,
		viewButton("Shuffle"),
		viewButton("Deselect All"),
		viewButton("Submit"),
	)
}

// viewButton creates and styles a button with the specified name.
func viewButton(name string) string {
	top := FGSpecial.Render(strings.Repeat("▄", CellWidth))
	label := Button.Render(name)
	bottom := FGSpecial.Render(strings.Repeat("▀", CellWidth))

	button := lipgloss.JoinVertical(lipgloss.Center, top, label, bottom)
	return zone.Mark(name, button)
}
