package crossword

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

// View renders the complete crossword puzzle UI.
func (m CrosswordModel) View() string {
	rows := make([]string, m.height+1)

	// First row is special - it has the top margin
	rows[0] = lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Bottom, m.viewMargin(0, "▄▄▄▄▄")),
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
	rows[m.height] = m.viewMargin(m.height-1, "▀▀▀▀▀")

	// Combine all elements vertically
	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, rows[:]...),
		m.viewCluesBox(),
		m.message,
	)
}

// viewCluesBox renders the box containing across and down clues.
// It displays the current clue in the middle with surrounding clues above and below.
func (m CrosswordModel) viewCluesBox() string {
	var clues [13]string

	// Get the across clue index for the current cursor position
	clueStartIdx := m.clueIndices[m.cursor.Y][m.cursor.X].X
	viewClues(m.acrossClues, clues[:], m.isAcrossSolved, clueStartIdx, true)

	// Add a vertical separator between across and down clues
	border := FGBorder.Render("│")
	for i := range 13 {
		clues[i] += " " + border + " "
	}

	// Get the down clue index for the current cursor position
	clueStartIdx = m.clueIndices[m.cursor.Y][m.cursor.X].Y
	if clueStartIdx <= -1 {
		clueStartIdx = 0
	}

	viewClues(m.downClues, clues[:], m.isDownSolved, clueStartIdx, false)

	// Create the header for the clues box
	acrossStyle := AcrossClue.Align(lipgloss.Center)
	downStyle := DownClue.Align(lipgloss.Center)
	top := acrossStyle.Render("Across") + "   " +
		downStyle.Render("Down") + "\n"

	return BorderStyle.Render(top + lipgloss.JoinVertical(lipgloss.Center, clues[:]...))
}

// viewClues renders a set of clues into the provided slice.
// It places the current clue in the middle and surrounding clues above and below.
func viewClues(clues, viewClues []string, isSolved []bool, startIdx int, isAcross bool) {
	// Select the appropriate style based on clue direction
	currentClueStyle := AcrossClue
	if !isAcross {
		currentClueStyle = DownClue
	}

	// Wrap the current clue
	lines := splitClue(clues[startIdx])
	n := len(lines)

	// Where to place the first line of the current clue so it’s visually centered:
	topOffset, bottomOffset := n/2, (n-1)/2
	startRow := 6 - topOffset

	// Render current clue lines
	for i, line := range lines {
		viewClues[startRow+i] += currentClueStyle.Render(line)
	}

	// Add clues that come before the current clue
	topSlots := 6 - topOffset
	cluesBefore := viewSurroundingClues(clues, isSolved, startIdx, -1, topSlots)
	for i, clue := range cluesBefore {
		viewClues[startRow-1-i] += clue
	}

	// Add clues that come after the current clue
	bottomSlots := 6 - bottomOffset
	cluesAfter := viewSurroundingClues(clues, isSolved, startIdx, +1, bottomSlots)
	for i, clue := range cluesAfter {
		viewClues[startRow+n+i] += clue
	}
}

// viewSurroundingClues returns up to maxLines of rendered clue lines in the given direction
// relative to startIdx. For direction -1 (above) the lines are ordered nearest-first topward.
// For +1 (below) it’s nearest-first downward.
func viewSurroundingClues(clues []string, isSolved []bool, startIdx, direction, maxLines int) []string {
	viewClues := make([]string, 0, maxLines)

	for step := 1; len(viewClues) < maxLines; step++ {
		// Calculate the index with wrapping
		wrappedIdx := (startIdx + direction*step + len(clues)) % len(clues)

		// Choose style based on whether the clue is solved
		style := NormalClue
		if isSolved[wrappedIdx] {
			style = SolvedClue
		}

		lines := splitClue(clues[wrappedIdx])

		if direction == -1 {
			// For clues above, reverse and prepend
			for i := len(lines) - 1; i >= 0 && len(viewClues) < maxLines; i-- {
				viewClues = append(viewClues, style.Render(lines[i]))
			}
		} else {
			// For clues below, append in order
			for i := 0; i < len(lines) && len(viewClues) < maxLines; i++ {
				viewClues = append(viewClues, style.Render(lines[i]))
			}
		}
	}

	return viewClues
}

// splitClue wraps a clue into lines of at most ClueWidth characters.
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

		// If the remaining clue is less than the maximum width, add it entirely and break
		if len(remaining) < ClueWidth {
			lines = append(lines, remaining)
			break
		}

		// Otherwise, find the last space in the remaining clue and split there
		splitIdx := strings.LastIndex(remaining[:ClueWidth], " ")
		lines = append(lines, remaining[:splitIdx])
		remaining = remaining[splitIdx+1:]
	}

	return lines
}

// viewMargin renders the top or bottom margin of a row in the grid.
func (m CrosswordModel) viewMargin(y int, cell string) string {
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
			top[x] = FGCursor.Render(cell)
		case isIncorrect:
			top[x] = FGIncorrect.Render(cell)
		case isAcross:
			top[x] = FGAcross.Render(cell)
		case isDown:
			top[x] = FGDown.Render(cell)
		case isEmpty:
			top[x] = "     " // Empty space for black cells
		default:
			top[x] = styleCell(isEven, FGEven, FGOdd).Render(cell)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, top[:]...)
}

// viewGridRow renders a row of cells in the crossword grid.
func (m CrosswordModel) viewGridRow(y int) string {
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
		cells[x] = fmt.Sprintf("%s%c  ", gridNum, char)

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
			cells[x] = AcrossCell.Render(cells[x])
		case isDown:
			cells[x] = DownCell.Render(cells[x])
		case isEmpty:
			cells[x] = "     "
		default:
			cells[x] = styleCell(isEven, EvenCell, OddCell).Render(cells[x])
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, cells[:]...)
}

// viewTopRow renders the connecting row between two grid rows.
func (m CrosswordModel) viewTopRow(y int) string {
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
			top[x] = FGCursor.Render("▄▄▄▄▄")
		case isCursor && isIncorrectAbove:
			top[x] = IncorrectTopCursor.Render("▀▀▀▀▀")
		case isCursor && !m.isAcross:
			top[x] = CursorDownTop.Render("▄▄▄▄▄")
		case isCursor:
			top[x] = styleCell(!isEven, CursorTopEven, CursorTopOdd).Render("▄▄▄▄▄")
		case isEmpty && isCursorAbove:
			top[x] = FGCursor.Render("▀▀▀▀▀")
		case isCursorAbove && isIncorrect:
			top[x] = IncorrectTopCursor.Render("▄▄▄▄▄")
		case isCursorAbove && !m.isAcross:
			top[x] = CursorDownTop.Render("▀▀▀▀▀")
		case isCursorAbove:
			top[x] = styleCell(isEven, CursorTopEven, CursorTopOdd).Render("▀▀▀▀▀")

		// Incorrect cell cases
		case isIncorrect && isEmptyAbove:
			top[x] = FGCursor.Render("▄▄▄▄▄")
		case isIncorrect && isIncorrectAbove:
			top[x] = FGIncorrect.Render("█████")
		case isIncorrect && isAcrossAbove:
			top[x] = IncorrectTopAcross.Render("▄▄▄▄▄")
		case isIncorrect && isDownAbove:
			top[x] = IncorrectTopDown.Render("▄▄▄▄▄")
		case isIncorrect:
			top[x] = styleCell(!isEven, IncorrectTopEven, IncorrectTopOdd).Render("▄▄▄▄▄")
		case isEmpty && isIncorrectAbove:
			top[x] = FGIncorrect.Render("▀▀▀▀▀")
		case isCursor && isIncorrectAbove:
			top[x] = IncorrectTopCursor.Render("▀▀▀▀▀")
		case isAcross && isIncorrectAbove:
			top[x] = IncorrectTopAcross.Render("▀▀▀▀▀")
		case isDown && isIncorrectAbove:
			top[x] = IncorrectTopDown.Render("▀▀▀▀▀")
		case isIncorrectAbove:
			top[x] = styleCell(isEven, IncorrectTopEven, IncorrectTopOdd).Render("▀▀▀▀▀")

		// Across clue cases
		case isAcross && isEmptyAbove:
			top[x] = FGAcross.Render("▄▄▄▄▄")
		case isAcross:
			top[x] = styleCell(!isEven, AcrossTopEven, AcrossTopOdd).Render("▄▄▄▄▄")
		case isEmpty && isAcrossAbove:
			top[x] = FGAcross.Render("▀▀▀▀▀")
		case isAcrossAbove:
			top[x] = styleCell(isEven, AcrossTopEven, AcrossTopOdd).Render("▀▀▀▀▀")

		// Down clue cases
		case isDown && isEmptyAbove:
			top[x] = FGDown.Render("▄▄▄▄▄")
		case isDown && !m.isAcross:
			top[x] = FGDown.Render("█████")
		case isDown:
			top[x] = styleCell(!isEven, DownTopEven, DownTopOdd).Render("▄▄▄▄▄")
		case isEmpty && isDownAbove:
			top[x] = FGDown.Render("▀▀▀▀▀")
		case isDownAbove:
			top[x] = styleCell(isEven, DownTopEven, DownTopOdd).Render("▀▀▀▀▀")

		// Empty cell cases
		case isEmpty && isEmptyAbove:
			top[x] = "     " // Space between black cells
		case isEmpty:
			top[x] = styleCell(isEven, FGOdd, FGEven).Render("▀▀▀▀▀")
		case isEmptyAbove:
			top[x] = styleCell(isEven, FGEven, FGOdd).Render("▄▄▄▄▄")

		// Normal cell connection
		default:
			top[x] = styleCell(isEven, TopEven, TopOdd).Render("▄▄▄▄▄")
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, top[:]...)
}

// viewGridNum renders the grid number for a cell as superscript characters.
func (m CrosswordModel) viewGridNum(x, y int) string {
	gridNum := m.gridNums[y][x]

	// If no grid number, return spaces
	if gridNum == 0 {
		return "  "
	}

	// Map of regular digits to superscript characters
	superscriptMap := map[rune]rune{
		'0': '⁰', '1': '¹', '2': '²', '3': '³', '4': '⁴',
		'5': '⁵', '6': '⁶', '7': '⁷', '8': '⁸', '9': '⁹',
	}

	// Convert number to string
	numStr := strconv.Itoa(gridNum)

	// Convert each digit to its superscript equivalent
	var result string
	for _, char := range numStr {
		if supChar, ok := superscriptMap[char]; ok {
			result += string(supChar)
		}
	}

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

// styleCell selects between two styles based on a boolean condition.
// Used to implement checkerboard patterns in the grid.
func styleCell(parity bool, firstStyle, secondStyle lipgloss.Style) lipgloss.Style {
	if parity {
		return firstStyle
	}
	return secondStyle
}
