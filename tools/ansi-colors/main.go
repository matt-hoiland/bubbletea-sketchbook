package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var colors = map[ansi.BasicColor]string{
	ansi.Black:         "black",
	ansi.Red:           "red",
	ansi.Green:         "green",
	ansi.Yellow:        "yellow",
	ansi.Blue:          "blue",
	ansi.Magenta:       "magenta",
	ansi.Cyan:          "cyan",
	ansi.White:         "white",
	ansi.BrightBlack:   "bright black",
	ansi.BrightRed:     "bright red",
	ansi.BrightGreen:   "bright green",
	ansi.BrightYellow:  "bright yellow",
	ansi.BrightBlue:    "bright blue",
	ansi.BrightMagenta: "bright magenta",
	ansi.BrightCyan:    "bright cyan",
	ansi.BrightWhite:   "bright white",
}

func main() {
	for i := range 16 {
		c := ansi.BasicColor(i)
		name := colors[c]
		style(c, name)
	}
}

func style(c ansi.BasicColor, name string) {
	fmt.Println(lipgloss.NewStyle().
		Background(c).
		Bold(true).
		Render(fmt.Sprintf("   %2d\t%s  ", c, name)))
}
