package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var wordmark = []string{
	"▄▄▄▄▄▄▄ ▄     ▄ ▄▄▄▄▄▄▄",
	"▓     ▓ ▓     ▓ ▓     ▓",
	"▒▀▀▀▀▀▒ ▒     ▒ ▒ ▀▀▀▀▀",
	"░     ░ ░▄▄▄▄▄░ ░      ",
}

const (
	bannerCols = 23
	bannerRows = 4
)

var bannerStyles = func() []lipgloss.Style {
	styles := make([]lipgloss.Style, bannerCols)
	for i, c := range lipgloss.Blend1D(bannerCols,
		lipgloss.Color("#6E5A9A"),
		accent,
		lipgloss.Color("#E8D6FF"),
		done,
		accent,
		lipgloss.Color("#6E5A9A"),
	) {
		styles[i] = lipgloss.NewStyle().Foreground(c)
	}
	return styles
}()

func banner(frame int) string {
	var lines []string
	for y, row := range wordmark {
		var b strings.Builder
		for x, r := range []rune(row) {
			if r == ' ' {
				b.WriteRune(' ')
				continue
			}
			b.WriteString(bannerStyles[(x+frame)%bannerCols].Render(string(r)))
		}
		if y == 0 {
			b.WriteString("  " + faint.Render("agent harness updater"))
		}
		lines = append(lines, b.String())
	}
	return strings.Join(lines, "\n")
}
