package block

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func (r BlockRenderer) move(style lipgloss.Style) string {
	style = style.MarginTop(1)
	return style.Render("●")
}

const text = `
▀	▁	▂	▃	▄	▅	▆	▇	█	▉	▊	▋	▌	▍	▎	▏

▐	░	▒	▓	▔	▕	▖	▗	▘	▙	▚	▛	▜	▝	▞	▟
`

func (r BlockRenderer) pawn(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(`  ▄▄▄  `),
		style.Render(`▄█████▄`),
		style.Render(`▀█████▀`),
		style.Render(`  ▀▀▀  `),
	}, "\n")
}

func (r BlockRenderer) king(style lipgloss.Style) string {
	return strings.Join([]string{
		style.Render(`  ▄▄▄  `),
		style.Render(`▄▀█▀█▀▄`),
		style.Render(`▀█▄▄▄█▀`),
		style.Render(`  ▀▀▀  `),
	}, "\n")
}
