package main

import "charm.land/lipgloss/v2"

func main() {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(2).
		PaddingLeft(4).
		Width(22)

	lipgloss.Println(style.Render("Hello, kitty"))
}
