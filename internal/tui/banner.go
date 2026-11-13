package tui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// The wordmark is a 4-row block-letter "aup". While a run is in progress a
// colour sweep travels across it, one column per spinner frame.
var wordmark = []string{
	"  ▄▄   █ █  ▄▄▄ ",
	" █▄▄█  █ █  █ █ ",
	" █  █  █ █  █▄█ ",
	" ▀▀▀▀  ▀▀▀  ▀   ",
}

const (
	bannerCols = 17
	bannerRows = 4
)

var bannerPalette = []color.Color{
	lipgloss.Color("#6E5A9A"),
	lipgloss.Color("#B392D7"),
	lipgloss.Color("#E8D6FF"),
	lipgloss.Color("#5FD7A7"),
	lipgloss.Color("#B392D7"),
	lipgloss.Color("#6E5A9A"),
}

func banner(frame int, live bool) string {
	colours := lipgloss.Blend1D(bannerCols, bannerPalette...)
	var lines []string
	for y, row := range wordmark {
		var b strings.Builder
		for x, r := range []rune(row) {
			if r == ' ' {
				b.WriteRune(' ')
				continue
			}
			i := x
			if live {
				i = (x + frame) % bannerCols
			}
			b.WriteString(lipgloss.NewStyle().Foreground(colours[i]).Render(string(r)))
		}
		if y == 0 {
			b.WriteString("  " + faint.Render("agent harness updater"))
		}
		lines = append(lines, b.String())
	}
	return strings.Join(lines, "\n")
}
