package connections

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

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
		if slices.Contains(m.selectedTiles, word) {
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

	return zone.Mark(text, lipgloss.JoinVertical(lipgloss.Left, margin, content, margin))
}
