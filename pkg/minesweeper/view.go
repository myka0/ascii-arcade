package minesweeper

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// View renders the entire Minesweeper UI.
func (m *MinesweeperModel) View() tea.View {
	if !m.hasSelected {
		return tea.NewView(m.viewSelection())
	}
	return tea.NewView(m.viewGame())
}

// viewSelection renders the difficulty selection menu.
func (m *MinesweeperModel) viewSelection() string {
	entries := make([]string, len(Difficulties))
	for i, d := range Difficulties {
		stats := fmt.Sprintf("(%d×%d, %d mines)\n", d.Rows, d.Cols, d.Mines)
		if m.cursor.Y == i {
			entries[i] = SelectedListEntry.Render("> " + d.Name + "\n  " + stats)
		} else {
			entries[i] = ListEntry.Render(d.Name + "\n" + stats)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		LabelStyle.Render("Minesweeper"),
		lipgloss.JoinVertical(lipgloss.Left, entries...),
	)
}

// viewGame renders the main game board and surrounding UI elements.
func (m *MinesweeperModel) viewGame() string {
	rows := make([]string, m.difficulty.Rows+1)

	// First row is special - it has the top margin
	rows[0] = lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewTopMargin()),
		lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewGridRow(0)),
	)

	// Middle rows have connecting pieces between cells
	for y := 1; y < m.difficulty.Rows; y++ {
		rows[y] = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewMarginRow(y)),
			lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewGridRow(y)),
		)
	}

	// Last row is the bottom margin
	rows[m.difficulty.Rows] = m.viewBottomMargin(m.difficulty.Rows - 1)

	// Combine all elements vertically
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"\n",
		m.viewTitleBar(),
		lipgloss.JoinVertical(lipgloss.Center, rows...),
		MessageStyle.MarginTop(1).Render(m.message),
		viewButtons(m.gameOver),
	)
}

// viewTitleBar builds the title bar showing game difficulty, mines left, and game timer.
func (m *MinesweeperModel) viewTitleBar() string {
	// Compute elapsed time components
	var elapsed time.Duration
	switch {
	case m.gameOver:
		elapsed = m.endTime.Sub(m.startTime)
	case m.hasStarted:
		elapsed = time.Since(m.startTime)
	}
	hours := int(elapsed.Hours())
	mins := int(elapsed.Minutes()) % 60
	secs := int(elapsed.Seconds()) % 60

	// Format labels
	title := LabelStyle.Render(fmt.Sprintf("Minesweeper — %s", m.difficulty.Name))
	minesLeft := fmt.Sprintf("Mines Left: %2d", m.difficulty.Mines-m.flagsPlaced)
	timer := fmt.Sprintf("Time: %02d:%02d:%02d", hours, mins, secs)

	// Hardcoded indentation values for difficulties
	var indent int
	switch m.difficulty.Name {
	case "Beginner":
		indent = 13
	case "Intermediate":
		indent = 18
	case "Expert":
		indent = 94
	}

	// Beginner is narrow so stats wrap onto a second line beneath the title
	if m.difficulty.Cols == Difficulties[Beginner].Cols {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			lipgloss.JoinHorizontal(lipgloss.Top,
				LabelStyle.MarginBottom(0).Render(minesLeft),
				LabelStyle.MarginBottom(0).MarginLeft(indent).Render(timer),
			),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		LabelStyle.MarginLeft(indent).Render(fmt.Sprintf("Mines Left: %2d", m.difficulty.Mines-m.flagsPlaced)),
		LabelStyle.MarginLeft(2).Render(timer),
	)
}

// viewButtons renders the reset and exit buttons shown after game over.
func viewButtons(gameOver bool) string {
	if !gameOver {
		return "\n"
	}

	resetButton := zone.Mark("reset", ButtonStyle.Render("Reset"))
	exitButton := zone.Mark("exit", ButtonStyle.MarginLeft(5).Render("Exit"))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		resetButton,
		exitButton,
	)
}

// viewTopMargin renders the top margin of a row in the grid.
func (m *MinesweeperModel) viewTopMargin() string {
	top := make([]string, m.difficulty.Cols)
	for x, cell := range m.board[0] {
		isEven := x%2 == 0
		isRevealed := cell.revealed
		isCursor := x == m.cursor.X && 0 == m.cursor.Y
		isHit := x == m.mineHit.X && 0 == m.mineHit.Y

		// Apply appropriate styling based on cell state
		var rendered string
		switch {
		case isHit:
			rendered = HitLowerBar
		case isCursor:
			rendered = CursorLowerBar
		case isRevealed:
			rendered = RevealedLowerBar
		default:
			rendered = switchStyle(isEven, EmptyBGLowerBar)
		}

		top[x] = rendered
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, top[:]...)
}

// viewBottomMargin renders the bottom margin of a row in the grid.
func (m *MinesweeperModel) viewBottomMargin(y int) string {
	bottom := make([]string, m.difficulty.Cols)
	for x, cell := range m.board[y] {
		isEven := (x+y)%2 == 0
		isRevealed := cell.revealed
		isCursor := x == m.cursor.X && y == m.cursor.Y
		isHit := x == m.mineHit.X && y == m.mineHit.Y

		// Apply appropriate styling based on cell state
		var rendered string
		switch {
		case isHit:
			rendered = HitUpperBar
		case isCursor:
			rendered = CursorUpperBar
		case isRevealed:
			rendered = RevealedUpperBar
		default:
			rendered = switchStyle(isEven, EmptyBGUpperBar)
		}

		bottom[x] = rendered
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, bottom[:]...)
}

// viewGridRow renders a row of cells in the crossword grid.
func (m *MinesweeperModel) viewGridRow(y int) string {
	cells := make([]string, m.difficulty.Cols)
	for x, cell := range m.board[y] {
		isEven := (x+y)%2 == 0
		isRevealed := cell.revealed
		isFlagged := cell.flagged
		isMine := cell.mine
		isCursor := x == m.cursor.X && y == m.cursor.Y
		isHit := x == m.mineHit.X && y == m.mineHit.Y

		// Apply appropriate styling based on cell state
		var rendered string
		switch {
		case isHit:
			rendered = HitCell
		case isCursor && isFlagged:
			rendered = CursorFlaggedBar
		case isCursor && isRevealed:
			rendered = CursorAdjacentStyles[cell.adjacent]
		case isCursor:
			rendered = CursorFullBar
		case isRevealed && isMine:
			rendered = MineCell
		case isRevealed:
			rendered = ColoredAdjacentStyles[cell.adjacent]
		case isFlagged && !isMine && m.gameOver:
			rendered = switchStyle(isEven, WrongFlaggedCell)
		case isFlagged:
			rendered = switchStyle(isEven, FlaggedCell)
		default:
			rendered = switchStyle(isEven, NormalBar)
		}

		cells[x] = zone.Mark(fmt.Sprintf("%d_%d", x, y), rendered)
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, cells[:]...)
}

// viewMarginRow renders the connecting row between two grid rows.
func (m *MinesweeperModel) viewMarginRow(y int) string {
	margin := make([]string, m.difficulty.Cols)
	for x, cell := range m.board[y] {
		isEven := (x+y)%2 == 0

		isRevealed := cell.revealed
		isRevealedAbove := m.board[y-1][x].revealed

		isCursor := x == m.cursor.X && y == m.cursor.Y
		isCursorAbove := x == m.cursor.X && y-1 == m.cursor.Y

		isHit := x == m.mineHit.X && y == m.mineHit.Y
		isHitAbove := x == m.mineHit.X && y-1 == m.mineHit.Y

		// Apply appropriate styling based on cell state
		var rendered string
		switch {
		case isHit && isRevealedAbove:
			rendered = HitRevealedMarginLowerBar
		case isHit:
			rendered = switchStyle(!isEven, HitMarginLowerBar)
		case isHitAbove && isRevealed:
			rendered = HitRevealedMarginUpperBar
		case isHitAbove:
			rendered = switchStyle(isEven, HitMarginUpperBar)

		case isCursor && isRevealedAbove:
			rendered = CursorRevealedMarginLowerBar
		case isCursor:
			rendered = switchStyle(!isEven, CursorMarginLowerBar)
		case isCursorAbove && isRevealed:
			rendered = CursorRevealedMarginUpperBar
		case isCursorAbove:
			rendered = switchStyle(isEven, CursorMarginUpperBar)

		case isRevealed && isRevealedAbove:
			rendered = RevealedFullBar
		case isRevealed:
			rendered = switchStyle(!isEven, RevealedMarginLowerBar)
		case isRevealedAbove:
			rendered = switchStyle(isEven, RevealedMarginUpperBar)

		default:
			rendered = switchStyle(isEven, MarginLowerBar)
		}

		margin[x] = rendered
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, margin[:]...)
}

// switchStyle selects between two styles based on a boolean condition.
// Used to implement checkerboard patterns in the grid.
func switchStyle(parity bool, styles [2]string) string {
	if parity {
		return styles[0]
	}
	return styles[1]
}