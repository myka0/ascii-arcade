package crossword

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// View renders the complete crossword puzzle UI.
func (m *CrosswordModel) View() tea.View {
	rows := make([]string, m.height+1)

	// First row is special - it has the top margin
	rows[0] = lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewTopMargin()),
		lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewGridRow(0)),
	)

	// Middle rows have connecting pieces between cells
	for y := 1; y < len(m.grid); y++ {
		rows[y] = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewTopRow(y)),
			lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewGridRow(y)),
		)
	}

	// Last row is the bottom margin
	rows[m.height] = m.viewBottomMargin(m.height - 1)

	// Render the crossword grid
	gridView := lipgloss.JoinVertical(lipgloss.Center, rows...)

	// Choose the clues layout based on puzzle kind
	var puzzleView string
	if m.kind == KindDaily {
		puzzleView = lipgloss.JoinVertical(
			lipgloss.Center,
			gridView,
			"",
			m.viewCluesBox(),
		)
	} else {
		puzzleView = lipgloss.JoinHorizontal(
			lipgloss.Center,
			gridView,
			m.viewMiniCluesBox(),
		)
	}

	// Combine all elements vertically
	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Center,
			puzzleView,
			MessageStyle.Render(m.message),
		),
	)
}

// viewCluesBox renders the box containing across and down clues.
// It displays the current clue in the middle with surrounding clues above and below.
func (m *CrosswordModel) viewCluesBox() string {
	// Get the across clue index for the current cursor position
	clueStartIdx := m.clueIndices[m.cursor.Y][m.cursor.X].X
	acrossLines := append(
		[]string{AcrossClue.Align(lipgloss.Center).Render("Across")},
		viewClues(m.acrossClues, m.isAcrossSolved, clueStartIdx, AcrossClue)...,
	)
	acrossClues := lipgloss.JoinVertical(lipgloss.Left, acrossLines...)

	// Get the down clue index for the current cursor position
	clueStartIdx = max(m.clueIndices[m.cursor.Y][m.cursor.X].Y, 0)

	downLines := append(
		[]string{DownClue.Align(lipgloss.Center).Render("Down")},
		viewClues(m.downClues, m.isDownSolved, clueStartIdx, DownClue)...,
	)
	downClues := lipgloss.JoinVertical(lipgloss.Left, downLines...)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().MarginRight(3).Render(acrossClues),
		downClues,
	)
}

// viewClues renders a set of clues into the provided slice.
// It places the current clue in the middle and surrounding clues above and below.
func viewClues(clues []string, isSolved []bool, startIdx int, activeStyle lipgloss.Style) []string {
	rows := make([]string, CluesVisibleRows)

	// Wrap the current clue
	lines := splitClue(clues[startIdx])
	n := len(lines)

	// Where to place the first line of the current clue so it’s visually centered:
	topOffset, bottomOffset := n/2, (n-1)/2
	startRow := CluesCenterRow - topOffset

	// Render current clue lines
	for i, line := range lines {
		rows[startRow+i] = activeStyle.Render(line)
	}

	// Add clues that come before the current clue
	topSlots := CluesCenterRow - topOffset
	cluesBefore := viewSurroundingClues(clues, isSolved, startIdx, -1, topSlots)
	for i, clue := range cluesBefore {
		rows[startRow-1-i] = clue
	}

	// Add clues that come after the current clue
	bottomSlots := CluesCenterRow - bottomOffset
	cluesAfter := viewSurroundingClues(clues, isSolved, startIdx, +1, bottomSlots)
	for i, clue := range cluesAfter {
		rows[startRow+n+i] = clue
	}

	return rows
}

// viewSurroundingClues returns up to maxLines of rendered clue lines in the given direction
// relative to startIdx. For direction -1 (above) the lines are ordered nearest-first topward.
// For +1 (below) it’s nearest-first downward.
func viewSurroundingClues(clues []string, isSolved []bool, startIdx, direction, maxLines int) []string {
	lines := make([]string, 0, maxLines)

	for step := 1; len(lines) < maxLines; step++ {
		// Calculate the index with wrapping
		wrappedIdx := (startIdx + direction*step + len(clues)) % len(clues)

		// Choose style based on whether the clue is solved
		style := NormalClue
		if isSolved[wrappedIdx] {
			style = SolvedClue
		}

		clueLines := splitClue(clues[wrappedIdx])

		if direction == -1 {
			// For clues above, reverse and prepend
			for i := len(clueLines) - 1; i >= 0 && len(lines) < maxLines; i-- {
				lines = append(lines, style.Render(clueLines[i]))
			}
		} else {
			// For clues below, append in order
			for i := 0; i < len(clueLines) && len(lines) < maxLines; i++ {
				lines = append(lines, style.Render(clueLines[i]))
			}
		}
	}

	return lines
}

// viewMiniCluesBox renders Mini puzzle clues in a single vertical column.
func (m *CrosswordModel) viewMiniCluesBox() string {
	clueIndex := m.clueIndices[m.cursor.Y][m.cursor.X]

	acrossLines := viewMiniColumn(m.acrossClues, m.isAcrossSolved, clueIndex.X, AcrossClue)
	downLines := viewMiniColumn(m.downClues, m.isDownSolved, clueIndex.Y, DownClue)

	// Preallocate enough capacity
	lines := make([]string, 0, len(acrossLines)+len(downLines)+3)

	headerStyle := NormalClue.Align(lipgloss.Center)

	lines = append(lines, headerStyle.Render("Across"))
	lines = append(lines, acrossLines...)
	lines = append(lines, "")
	lines = append(lines, headerStyle.Render("Down"))
	lines = append(lines, downLines...)

	return lipgloss.NewStyle().
		MarginLeft(4).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

// viewMiniColumn returns a rendered lines for one side of the Mini clues.
func viewMiniColumn(clues []string, solved []bool, activeIdx int, activeStyle lipgloss.Style) []string {
	lines := make([]string, 0, len(clues))

	for i, clue := range clues {
		// Style the clue
		style := NormalClue
		if i == activeIdx {
			style = activeStyle
		} else if i < len(solved) && solved[i] {
			style = SolvedClue
		}

		// Wrap the clue and render it
		for _, line := range splitClue(clue) {
			lines = append(lines, style.Render(line))
		}
	}

	return lines
}

// splitClue wraps a clue into lines of at most ClueWidth display cells.
func splitClue(clue string) []string {
	var lines []string
	remaining := clue
	indent := strings.Repeat(" ", strings.Index(clue, " ")+1)

	// Split the clue into lines
	for idx := 0; len(remaining) > 0; idx++ {
		// Add indentation if this is not the first line
		if idx > 0 {
			remaining = indent + remaining
		}

		// If what's left fits in the box, emit it and stop
		if ansi.StringWidth(remaining) <= ClueWidth {
			lines = append(lines, remaining)
			break
		}

		// Otherwise, find the last space in the remaining clue and split there
		splitIdx := strings.LastIndex(remaining[:ClueWidth], " ")
		head, tail := remaining[:splitIdx], remaining[splitIdx+1:]
		lines = append(lines, head)
		remaining = tail
	}

	return lines
}

// viewTopMargin renders the top margin of a row in the grid.
func (m *CrosswordModel) viewTopMargin() string {
	top := make([]string, m.width)
	for x, char := range m.grid[0] {
		isEven := x%2 == 0
		isIncorrect := m.incorrect[0][x]
		isCursor := x == m.cursor.X && 0 == m.cursor.Y
		isAcross := m.clue == m.clueAt(x, 0) && m.isAcross
		isDown := m.clue == m.clueAt(x, 0) && !m.isAcross && m.clueIndices[0][x].Y != -1
		isEmpty := char == '.'

		// Apply appropriate styling based on cell state
		switch {
		case isCursor:
			top[x] = CursorLowerBar
		case isIncorrect:
			top[x] = IncorrectLowerBar
		case isAcross:
			top[x] = AcrossLowerBar
		case isDown:
			top[x] = DownLowerBar
		case isEmpty:
			top[x] = Blank
		default:
			top[x] = switchStyledBar(isEven, EmptyBGLowerBar)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, top[:]...)
}

// viewBottomMargin renders the bottom margin of a row in the grid.
func (m *CrosswordModel) viewBottomMargin(y int) string {
	top := make([]string, m.width)
	for x, char := range m.grid[y] {
		isEven := (x+y)%2 == 0
		isIncorrect := m.incorrect[y][x]
		isCursor := x == m.cursor.X && y == m.cursor.Y
		isAcross := m.clue == m.clueAt(x, y) && m.isAcross
		isDown := m.clue == m.clueAt(x, y) && !m.isAcross && m.clueIndices[y][x].Y != -1
		isEmpty := char == '.'

		// Apply appropriate styling based on cell state
		switch {
		case isCursor:
			top[x] = CursorUpperBar
		case isIncorrect:
			top[x] = IncorrectUpperBar
		case isAcross:
			top[x] = AcrossUpperBar
		case isDown:
			top[x] = DownUpperBar
		case isEmpty:
			top[x] = Blank
		default:
			top[x] = switchStyledBar(isEven, EmptyBGUpperBar)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, top[:]...)
}

// viewGridRow renders a row of cells in the crossword grid.
func (m *CrosswordModel) viewGridRow(y int) string {
	cells := make([]string, m.width)
	for x, char := range m.grid[y] {
		isEven := (x+y)%2 == 0
		isIncorrect := m.incorrect[y][x]
		isCursor := x == m.cursor.X && y == m.cursor.Y
		isAcross := m.clue == m.clueAt(x, y) && m.isAcross
		isDown := m.clue == m.clueAt(x, y) && !m.isAcross && m.clueIndices[y][x].Y != -1
		isEmpty := char == '.'

		// Get the grid number for this cell
		gridNum := m.viewGridNum(x, y)
		cellContent := fmt.Sprintf("%s%c  ", gridNum, char)
		cells[x] = cellContent

		// Apply appropriate styling based on cell state
		switch {
		case isIncorrect && isCursor:
			cell := CursorCell.Underline(true).Italic(true).Render(string(char))
			cells[x] = CursorCell.Render(gridNum) + cell + CursorCell.Render("  ")
		case isIncorrect:
			cell := IncorrectCell.Italic(true).Render(string(char))
			cells[x] = IncorrectCell.Render(gridNum) + cell + IncorrectCell.Render("  ")
		case isCursor:
			cell := CursorCell.Underline(true).Render(string(char))
			cells[x] = CursorCell.Render(gridNum) + cell + CursorCell.Render("  ")
		case isAcross:
			cells[x] = AcrossCell.Render(cellContent)
		case isDown:
			cells[x] = DownCell.Render(cellContent)
		case isEmpty:
			cells[x] = Blank
		case isEven:
			cells[x] = EvenCell.Render(cellContent)
		default:
			cells[x] = OddCell.Render(cellContent)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, cells[:]...)
}

// viewTopRow renders the connecting row between two grid rows.
func (m *CrosswordModel) viewTopRow(y int) string {
	top := make([]string, m.width)
	for x, char := range m.grid[y] {
		isEven := (x+y)%2 == 0

		isIncorrect := m.incorrect[y][x]
		isIncorrectAbove := m.incorrect[y-1][x]

		isCursor := x == m.cursor.X && y == m.cursor.Y
		isCursorAbove := x == m.cursor.X && y-1 == m.cursor.Y

		isAcross := m.clue == m.clueAt(x, y) && m.isAcross
		isAcrossAbove := m.clue == m.clueAt(x, y-1) && m.isAcross

		isDown := m.clue == m.clueAt(x, y) && !m.isAcross && m.clueIndices[y][x].Y != -1
		isDownAbove := m.clue == m.clueAt(x, y-1) && !m.isAcross && m.clueIndices[y-1][x].Y != -1

		isEmpty := char == '.'
		isEmptyAbove := m.grid[y-1][x] == '.'

		// Apply appropriate styling based on the complex state combinations
		// This large switch statement handles all possible combinations of cell states
		switch {
		// Cursor related cases
		case isCursor && isEmptyAbove:
			top[x] = CursorLowerBar
		case isCursor && isIncorrectAbove:
			top[x] = IncorrectTopCursorUpperBar
		case isCursor && !m.isAcross:
			top[x] = CursorDownTopLowerBar
		case isCursor:
			top[x] = switchStyledBar(!isEven, CursorTopLowerBar)
		case isEmpty && isCursorAbove:
			top[x] = CursorUpperBar
		case isCursorAbove && isIncorrect:
			top[x] = IncorrectTopCursorLowerBar
		case isCursorAbove && !m.isAcross:
			top[x] = CursorDownTopUpperBar
		case isCursorAbove:
			top[x] = switchStyledBar(isEven, CursorTopUpperBar)

		// Incorrect cell cases
		case isIncorrect && isEmptyAbove:
			top[x] = CursorLowerBar
		case isIncorrect && isIncorrectAbove:
			top[x] = IncorrectFullBar
		case isIncorrect && isAcrossAbove:
			top[x] = IncorrectTopAcrossLowerBar
		case isIncorrect && isDownAbove:
			top[x] = IncorrectTopDownLowerBar
		case isIncorrect:
			top[x] = switchStyledBar(!isEven, IncorrectTopLowerBar)
		case isEmpty && isIncorrectAbove:
			top[x] = IncorrectUpperBar
		case isAcross && isIncorrectAbove:
			top[x] = IncorrectTopAcrossUpperBar
		case isDown && isIncorrectAbove:
			top[x] = IncorrectTopDownUpperBar
		case isIncorrectAbove:
			top[x] = switchStyledBar(isEven, IncorrectTopUpperBar)

		// Across clue cases
		case isAcross && isEmptyAbove:
			top[x] = AcrossLowerBar
		case isAcross:
			top[x] = switchStyledBar(!isEven, AcrossTopLowerBar)
		case isEmpty && isAcrossAbove:
			top[x] = AcrossUpperBar
		case isAcrossAbove:
			top[x] = switchStyledBar(isEven, AcrossTopUpperBar)

		// Down clue cases
		case isDown && isEmptyAbove:
			top[x] = DownLowerBar
		case isDown && !m.isAcross:
			top[x] = DownFullBar
		case isDown:
			top[x] = switchStyledBar(!isEven, DownTopLowerBar)
		case isEmpty && isDownAbove:
			top[x] = DownUpperBar
		case isDownAbove:
			top[x] = switchStyledBar(isEven, DownTopUpperBar)

		// Empty cell cases
		case isEmpty && isEmptyAbove:
			top[x] = Blank // Space between black cells
		case isEmpty:
			top[x] = switchStyledBar(!isEven, EmptyBGUpperBar)
		case isEmptyAbove:
			top[x] = switchStyledBar(isEven, EmptyBGLowerBar)

		// Normal cell connection
		default:
			top[x] = switchStyledBar(isEven, TopLowerBar)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, top[:]...)
}

// viewGridNum renders the grid number for a cell as superscript characters.
func (m *CrosswordModel) viewGridNum(x, y int) string {
	gridNum := m.gridNums[y][x]

	// If no grid number, return spaces
	if gridNum == 0 {
		return "  "
	}

	// Map of regular digits to superscript characters
	superscriptDigits := [10]rune{'⁰', '¹', '²', '³', '⁴', '⁵', '⁶', '⁷', '⁸', '⁹'}

	// Convert number to string
	numStr := strconv.Itoa(gridNum)

	// Convert each digit to its superscript equivalent
	var b strings.Builder
	b.Grow(len(numStr) * 3)
	for _, char := range numStr {
		b.WriteRune(superscriptDigits[char-'0'])
	}
	result := b.String()

	// Add padding for single digit numbers
	if len(numStr) == 1 {
		return result + " "
	}

	// TODO Handle 3 digit numbers properly
	if len(numStr) == 3 {
		runes := []rune(result)
		return string(runes[len(runes)-2:])
	}

	return result
}

// switchStyledBar selects between two styles based on a boolean condition.
// Used to implement checkerboard patterns in the grid.
func switchStyledBar(parity bool, styles [2]string) string {
	if parity {
		return styles[0]
	}
	return styles[1]
}