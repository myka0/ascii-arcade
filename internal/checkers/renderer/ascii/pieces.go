package ascii

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

func (r AsciiRenderer) move(style lipgloss.Style) string {
	style = style.MarginTop(1)
	return style.Render("o")
}

func (r AsciiRenderer) pawn(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(` ,gPPRg, `),
		style.Render(`dP'   'Yb`),
		style.Render(`Yb     dP`),
		style.Render(` "8ggg8" `),
	}, "\n")
}

func (r AsciiRenderer) king(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(` ,gPPRg, `),
		style.Render(`dP' K 'Yb`),
		style.Render(`Yb  K  dP`),
		style.Render(` "8ggg8" `),
	}, "\n")
}
