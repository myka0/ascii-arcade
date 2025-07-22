package nerdfont

import (
	"github.com/charmbracelet/lipgloss/v2"
)

func (r NerdfontRenderer) move(style lipgloss.Style) string {
	return style.Render("")
}

func (r NerdfontRenderer) pawn(style lipgloss.Style) string {
	return style.Render("")
}

func (r NerdfontRenderer) rook(style lipgloss.Style) string {
	return style.Render("  ")
}

func (r NerdfontRenderer) knight(style lipgloss.Style) string {
	return style.Render("  ")
}

func (r NerdfontRenderer) bishop(style lipgloss.Style) string {
	return style.Render("")
}

func (r NerdfontRenderer) queen(style lipgloss.Style) string {
	return style.Render("  ")
}

func (r NerdfontRenderer) king(style lipgloss.Style) string {
	return style.Render("  ")
}
