package gogame

import (
	"fmt"

	"ascii-arcade/internal/colors"
	"ascii-arcade/pkg/overlay"

	"image/color"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
	"github.com/muesli/reflow/ansi"
)

// View renders the entire Go board UI.
func (m *GoModel) View() tea.View {
	if m.gameOver && !m.markingDeadStones {
		return tea.NewView(m.viewGameOver())
	}
	return tea.NewView(m.viewBoard())
}

// viewBoard assembles the full board UI.
func (m *GoModel) viewBoard() string {
	grid := m.viewGrid()
	borderStyle := Border
	if m.showLabels {
		borderStyle = BorderLabels
	}

	board := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Right, m.viewStatus(), grid),
	)

	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewTitle(),
		board,
		m.viewButtons(),
	)

	message := m.viewMessage()
	return lipgloss.JoinVertical(lipgloss.Center, ui, message)
}

// viewTitle builds the title bar showing whose turn it is and the board size.
func (m *GoModel) viewTitle() string {
	var title string
	if m.markingDeadStones {
		title = TitleStyle.Render("Go — Mark Dead Stones")
	} else {
		title = TitleStyle.Render(fmt.Sprintf("Go — %s's Turn", colorName(m.turn)))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		viewBoardSizeIndicator(m.boardSize, m.markingDeadStones),
	)
}

// viewGrid builds the complete grid of intersections, stones, and labels.
func (m *GoModel) viewGrid() string {
	var rows []string

	// Top edge row
	rows = append(rows, m.viewGridRow(gridTopLeft, gridTopIntersect, gridTopRight, 0))

	// Middle rows with connectors between them
	for y := 1; y < m.boardSize-1; y++ {
		rows = append(rows, m.viewGridConnectors())
		rows = append(rows, m.viewGridRow(gridMidLeft, gridMidIntersect, gridMidRight, y))
	}

	// Bottom edge row
	rows = append(rows, m.viewGridConnectors())
	rows = append(rows, m.viewGridRow(gridBotLeft, gridBotIntersect, gridBotRight, m.boardSize-1))

	// Column labels below the board
	if m.showLabels {
		rows = append(rows, m.viewColumnLabels())
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// viewColumnLabels renders the letter coordinate labels below the board.
func (m *GoModel) viewColumnLabels() string {
	var row []string
	row = append(row, LabelStyle.Render(" "))
	for _, l := range columnLabels(m.boardSize) {
		row = append(row, LabelStyle.Render(fmt.Sprintf("  %s ", l)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, row...)
}

// viewMessage returns the rendered status message or an empty styled string.
func (m *GoModel) viewMessage() string {
	if m.message != "" {
		if ansi.PrintableRuneWidth(m.message)%2 == 0 {
			m.message += " "
		}
		return MessageStyle.Render("> " + m.message)
	}
	return MessageStyle.Render("")
}

// viewGridRow renders a single row of intersection cells in the grid.
func (m *GoModel) viewGridRow(leftCorner, intersection, rightCorner string, y int) string {
	connector := BoardLine.Render(gridHorizBar)
	var row []string

	// Row number label on the left side
	if m.showLabels {
		label := fmt.Sprintf("%2d", m.boardSize-y)
		row = append(row, LabelStyle.Render(label))
	}

	row = append(row, m.buildInteractiveCell(leftCorner, 0, y))

	for i := 1; i < m.boardSize-1; i++ {
		row = append(row, connector)
		row = append(row, m.buildInteractiveCell(intersection, i, y))
	}

	row = append(row, connector)
	row = append(row, m.buildInteractiveCell(rightCorner, m.boardSize-1, y))

	return lipgloss.JoinHorizontal(lipgloss.Top, row...)
}

// buildInteractiveCell creates a single intersection and wraps it in a clickable zone label.
func (m *GoModel) buildInteractiveCell(intersection string, x, y int) string {
	cell := m.viewIntersection(Position{X: x, Y: y}, intersection)
	return zone.Mark(fmt.Sprintf("%d_%d", x, y), cell)
}

// viewGridConnectors renders a row of vertical bar connectors between grid rows.
func (m *GoModel) viewGridConnectors() string {
	var row []string
	if m.showLabels {
		row = append(row, BoardLine.Render("  "+gridVertBar))
	} else {
		row = append(row, BoardLine.Render(gridVertBar))
	}
	for i := 1; i < m.boardSize; i++ {
		row = append(row, BoardLine.Render(" "+gridVertBar))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, row...)
}

// viewIntersection renders a single intersection or stone on the board.
func (m *GoModel) viewIntersection(pos Position, defaultIntersection string) string {
	cell := m.board.Cells[pos.Y][pos.X]

	// Render marked dead stones
	if m.markingDeadStones && m.deadStones[pos] {
		return DeadStoneStyle.Render("▐█▌")
	}

	// Highlight the cursor position
	if m.cursor.X == pos.X && m.cursor.Y == pos.Y && !m.markingDeadStones {
		return CursorStyle.Render("▐█▌")
	}

	// Render the base cell character
	switch cell {
	case Black:
		return BlackStoneStyle.Render("▐█▌")
	case White:
		return WhiteStoneStyle.Render("▐█▌")
	default:
		return BoardLine.Render(defaultIntersection)
	}
}

// viewStatus renders the capture counts line.
func (m *GoModel) viewStatus() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		LabelStyle.Render("Captures  "),
		BlackStoneStyle.Render("●"),
		LabelStyle.Render(fmt.Sprintf(" %d  ", m.blackCaptures)),
		WhiteStoneStyle.Render("●"),
		LabelStyle.Render(fmt.Sprintf(" %d ", m.whiteCaptures)),
	)
}

// viewButtons renders the bottom row of buttons.
func (m *GoModel) viewButtons() string {
	// Show Score button during stone marking phase
	if m.markingDeadStones {
		scoreBtn := zone.Mark("score", ButtonStyle.Render("Score"))
		return lipgloss.JoinHorizontal(lipgloss.Top, scoreBtn)
	}

	passBtn := zone.Mark("pass", ButtonStyle.Render("Pass"))
	resignBtn := zone.Mark("resign", ButtonStyle.Render("Resign"))
	return lipgloss.JoinHorizontal(lipgloss.Top, passBtn, resignBtn)
}

// viewBoardSizeIndicator returns a styled board size label
// positioned to align with the right edge of the board grid.
func viewBoardSizeIndicator(size int, markingDeadStones bool) string {
	// Hardcoded indentation values for board sizes
	indent := 0
	switch size {
	case 9:
		indent = 19
	case 13:
		indent = 33
	case 19:
		indent = 57
	}

	// Decrease the indentation to fit marking stones label
	if markingDeadStones {
		indent -= 4
	}

	return TitleStyle.MarginLeft(indent).Render(fmt.Sprintf("%d×%d", size, size))
}

// viewGameOver renders the end of game UI.
func (m *GoModel) viewGameOver() string {
	mainView := m.viewBoard()
	bgColor := colors.Dark2

	// Determine the game outcome
	var (
		winnerText string
		winColor   color.Color
	)
	switch m.score.Winner {
	case Black:
		winnerText = fmt.Sprintf("Black wins!  %.1f vs %.1f", m.score.BlackScore, m.score.WhiteScore)
		winColor = colors.Orange
	case White:
		winnerText = fmt.Sprintf("White wins!  %.1f vs %.1f", m.score.WhiteScore, m.score.BlackScore)
		winColor = colors.Purple
	default:
		winnerText = fmt.Sprintf("Draw!  %.1f each", m.score.BlackScore)
		winColor = colors.Blue
	}
	winnerStyled := lipgloss.NewStyle().Foreground(winColor).Render(winnerText)

	// Style for interactive buttons
	buttonStyle := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(winColor).
		Padding(0, 1)

	// Button alignment box
	buttonBox := lipgloss.NewStyle().
		Background(bgColor).
		Width(12)

	resetButton := zone.Mark("reset", buttonStyle.Render("Reset"))
	exitButton := zone.Mark("exit", buttonStyle.Render("Exit"))
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		buttonBox.Align(lipgloss.Left).Render(resetButton),
		buttonBox.Align(lipgloss.Right).Render(exitButton),
	)

	// Last move notation
	details := ""
	if m.lastMove != nil {
		details = lipgloss.NewStyle().
			Foreground(colors.Light2).
			Render("Last move: " + formatPosition(*m.lastMove, m.boardSize))
	}

	return overlay.PlaceNotification(
		mainView,
		"Game over.",
		winnerStyled,
		details,
		buttons,
	)
}

// formatPosition converts a board Position into standard Go notation
func formatPosition(pos Position, boardSize int) string {
	col := 'A' + rune(pos.X)
	// Skip I to avoid confusion with 1
	if col >= 'I' {
		col++
	}
	row := boardSize - pos.Y
	return fmt.Sprintf("%c%d", col, row)
}

// columnLabels returns column labels A-J (skipping I) for the given board size.
func columnLabels(size int) []string {
	labels := make([]string, size)
	letter := 'A'
	for i := range size {
		// Skip I to avoid confusion with 1
		if letter == 'I' {
			letter++
		}
		labels[i] = string(letter)
		letter++
	}
	return labels
}
