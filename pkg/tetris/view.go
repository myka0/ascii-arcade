package tetris

import (
	"ascii-arcade/internal/colors"
	"ascii-arcade/internal/components"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
)

// View renders the entire Tetris game view.
func (m *TetrisModel) View() tea.View {
	if m.gameOver {
		return tea.NewView(m.viewGameOver())
	}
	return tea.NewView(m.viewGame())
}

// viewGame renders the active gameplay layout.
func (m *TetrisModel) viewGame() string {
	// While paused, hide the side panels and show only the board
	if m.paused {
		return m.viewBoard()
	}

	// Stack the Hold and Info panels vertically on the left
	left := lipgloss.JoinVertical(lipgloss.Left,
		m.viewHold(),
		m.viewInfo(),
	)

	board := m.viewBoard()
	right := m.viewNext()

	// Join the three columns
	return lipgloss.JoinHorizontal(lipgloss.Top, left, board, right)
}

// viewBoard renders the 20x10 visible matrix.
func (m *TetrisModel) viewBoard() string {
	matrix, _ := m.game.GetVisibleMatrix()

	// Render each row of cells with newlines in between
	var board strings.Builder
	for row := range matrix {
		for col := range matrix[row] {
			// Each engine cell becomes a two-character "pixel" via renderCell
			board.WriteString(renderCell(matrix[row][col]))
		}
		if row < len(matrix)-1 {
			board.WriteByte('\n')
		}
	}

	return Playfield.Render(board.String())
}

// viewHold renders the held tetrimino inside frame.
func (m *TetrisModel) viewHold() string {
	label := PanelLabel.Render("Hold")
	piece := renderSingleTetrimino(m.game.GetHoldTetrimino())
	return HoldBox.Render(lipgloss.JoinVertical(lipgloss.Left, label, piece))
}

// viewNext renders the next few tetriminoes from the engine's bag.
func (m *TetrisModel) viewNext() string {
	label := PanelLabel.Render("Next")

	bag := m.game.GetBagTetriminos()
	rows := []string{label}
	for i := 0; i < len(bag) && i < nextQueueSize; i++ {
		t := bag[i]
		rows = append(rows, renderSingleTetrimino(&t)+"\n")
	}
	return NextBox.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

// viewInfo renders Score, Level, Lines, and Time stacked vertically.
func (m *TetrisModel) viewInfo() string {
	label := PanelLabel.Render("Info")

	// Get elapsed time, accounting for pauses
	var elapsed time.Duration
	if m.hasStarted {
		elapsed = time.Since(m.startTime) - m.pausedDuration
		if m.paused && !m.pausedAt.IsZero() {
			elapsed -= time.Since(m.pausedAt)
		}
	}
	mins := int(elapsed.Minutes())
	secs := int(elapsed.Seconds()) % 60

	rows := []string{
		label,
		infoRow("Score", fmt.Sprintf("%d", m.game.GetTotalScore())),
		infoRow("Lines", fmt.Sprintf("%d", m.game.GetLinesCleared())),
		infoRow("Level", fmt.Sprintf("%d", m.game.GetLevel())),
		"",
		infoRow("Time", fmt.Sprintf("%02d:%02d", mins, secs)),
	}
	return InfoBox.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

// infoRow formats a "label: value" line for the info panel.
func infoRow(label, value string) string {
	return InfoLabel.Render(label+":") + " " + InfoValue.Render(value)
}

// viewGameOver renders the end of game UI.
func (m *TetrisModel) viewGameOver() string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		infoRow("Score", fmt.Sprintf("%d", m.game.GetTotalScore())),
		infoRow("Lines", fmt.Sprintf("%d", m.game.GetLinesCleared())),
		infoRow("Level", fmt.Sprintf("%d", m.game.GetLevel())),
	)

	return components.GameOver(colors.Purple, m.viewGame(), content)
}

// renderSingleTetrimino draws a single tetrimino.
func renderSingleTetrimino(t *tetris.Tetrimino) string {
	if t == nil || t.Value == 0 || len(t.Cells) == 0 {
		// Render a blank 2x4 preview so the panel keeps its shape
		blank := strings.Repeat(cellEmpty, 4)
		return emptyStyle.Render(blank + "\n" + blank)
	}

	// Build the tetrimino string row by row from its cell grid
	var b strings.Builder
	for row := range t.Cells {
		for col := range t.Cells[row] {
			if t.Cells[row][col] {
				// Active: render with the piece's color
				b.WriteString(renderCell(t.Value))
			} else {
				// Inactive: render as empty space
				b.WriteString(emptyStyle.Render(cellEmpty))
			}
		}
		// Separate rows with newlines, but omit trailing newline
		if row < len(t.Cells)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// renderCell converts a byte from the engine's visible matrix into a styled two-character string.
func renderCell(cell byte) string {
	if cell == ghost {
		return ghostStyle.Render(cellGhost)
	} else if style, ok := cellStyles[cell]; ok {
		return style.Render(cellFilled)
	}

	return emptyStyle.Render(cellEmpty)
}
