package overlay

import (
	"ascii-arcade/internal/colors"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	charmansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/reflow/ansi"
)

// Position represents a position along a horizontal or vertical axis.
type Position float64

func (p Position) value() float64 {
	return math.Min(1, math.Max(0, float64(p)))
}

// Position aliases.
const (
	Top    Position = 0.0
	Bottom Position = 1.0
	Center Position = 0.5
	Left   Position = 0.0
	Right  Position = 1.0
)

// Place overlays the foreground (fg) string on top of the background (bg)
// string, aligned according to hPos and vPos. It preserves ANSI styling.
func Place(hPos, vPos Position, fg, bg string) string {
	bgLines, bgWidth := getLines(bg)
	fgLines, fgWidth := getLines(fg)

	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	vGap := bgHeight - fgHeight
	hGap := bgWidth - fgWidth

	// If overlay is larger, just return it
	if vGap <= 0 || hGap <= 0 {
		return fg
	}

	var b strings.Builder
	b.WriteRune('\n')

	// Compute vertical split
	top := int(math.Round(float64(vGap) * vPos.value()))
	bottom := vGap - top

	// Top background lines
	for i := range top {
		b.WriteString(bgLines[i])
		b.WriteRune('\n')
	}

	// Overlay region
	for i := range fgHeight {
		fgLine := fgLines[i]
		bgLine := bgLines[top+i]

		hSplit := int(math.Round(float64(bgWidth-fgWidth) * hPos.value()))

		// Left background portion
		left := charmansi.Truncate(bgLine, hSplit, "")
		leftLength := ansi.PrintableRuneWidth(left)

		// Right portion after the overlay
		right := charmansi.TruncateLeft(bgLine, leftLength+fgWidth, "")

		b.WriteString(left)
		b.WriteString(fgLine)
		b.WriteString(right)
		b.WriteRune('\n')
	}

	// Bottom background lines
	for i := range bottom {
		b.WriteString(bgLines[bgHeight-bottom+i])
		b.WriteRune('\n')
	}

	return b.String()
}

// NewNotification creates a styled new notification.
func NewNotification(content string) string {
	return lipgloss.NewStyle().
		Padding(1, 4, 0, 4).
		Foreground(colors.Light2).
		Background(colors.Dark2).
		Render(content)
}

// PlaceNotification creates a notification and places it centered on the main view.
func PlaceNotification(mainView string, notifContent ...string) string {
	maxWidth := 0
	for _, line := range notifContent {
		_, w := getLines(line)
		if maxWidth < w {
			maxWidth = w
		}
	}

	block := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Background(colors.Dark2).
		PaddingBottom(1).
		Width(maxWidth)

	var content []string
	for _, line := range notifContent {
		content = append(content, block.Foreground(colors.Light2).Render(line))
	}

	notifBody := lipgloss.JoinVertical(
		lipgloss.Center,
		content...,
	)

	notification := NewNotification(notifBody)
	return Place(Center, Center, notification, mainView)
}

// Obtained from https://github.com/charmbracelet/lipgloss/blob/master/get.go
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := ansi.PrintableRuneWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
