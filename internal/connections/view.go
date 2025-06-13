package connections

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

const (
	gap       = 2
	cellWidth = 16
	clueWidth = cellWidth*4 + gap*3
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
		FGLightText.Render("Create four groups of four!\n"),
		board,
		m.viewMistakesRemaining()+"\n",
		FGLightText.Render(m.message),
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

	// Pad and join words
	var words [4]string
	for i, word := range group.Members {
		padding := cellWidth - len(word)
		left := padding / 2
		right := padding - left
		words[i] = strings.Repeat(" ", left) + word + strings.Repeat(" ", right)
	}
	wordRow := RevealedLine.Background(color).Render(strings.Join(words[:], "  "))

	// Pad and render clue
	padding := clueWidth - len(group.Clue)
	left := padding / 2
	right := padding - left
	clueRow := RevealedLine.Background(color).Render(
		strings.Repeat(" ", left) + group.Clue + strings.Repeat(" ", right),
	)

	// Top and bottom margin
	margin := RevealedLine.Background(color).Render(strings.Repeat(" ", 70))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		margin,
		wordRow,
		clueRow,
		margin) + "\n"
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
		cells[col] = styleCell(word, cellWidth, style)
	}

	// Join all 4 cells horizontally
	return lipgloss.JoinHorizontal(lipgloss.Bottom, cells[:]...) + "\n"
}

func (m ConnectionsModel) viewMistakesRemaining() string {
	mistakes := MistakeCell.Render(strings.Repeat("● ", m.mistakesRemaining))
	return FGLightText.Render("Mistakes Remaining: ") + mistakes
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

	return zone.Mark(text, lipgloss.JoinVertical(lipgloss.Left, margin, content, margin))
}
