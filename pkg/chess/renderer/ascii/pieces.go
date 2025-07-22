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
		style.Render(`       `),
		style.Render(`  ( )  `),
		style.Render(`  ) (  `),
		style.Render(` (___) `),
	}, "\n")
}

func (r AsciiRenderer) rook(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(`[ U U ]`),
		style.Render(` |   | `),
		style.Render(` |   | `),
		style.Render(`[_____]`),
	}, "\n")
}

func (r AsciiRenderer) knight(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(` /\v/\ `),
		style.Render(`/(o o)\`),
		style.Render(`  | |  `),
		style.Render(`  (∞)  `),
	}, "\n")
}

func (r AsciiRenderer) bishop(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(` /^\ `),
		style.Render(` ( ) `),
		style.Render(` ) ( `),
		style.Render(`(___)`),
	}, "\n")
}

func (r AsciiRenderer) queen(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(`  /o\  `),
		style.Render(` (   ) `),
		style.Render(`  ) (  `),
		style.Render(`(_____)`),
	}, "\n")
}

func (r AsciiRenderer) king(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(`  (+)  `),
		style.Render(` (   ) `),
		style.Render(`  ) (  `),
		style.Render(`(_____)`),
	}, "\n")
}
