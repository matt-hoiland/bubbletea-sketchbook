package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/matt-hoiland/bubbletea-sketchbook/tools/vanity/model"
)

func main() {
	if _, err := tea.NewProgram(model.NewWindow()).Run(); err != nil {
		header := lipgloss.NewStyle().
			Foreground(lipgloss.Red).
			Bold(true).
			Render("[ ERROR ]")
		lipgloss.Printf("%s: %v", header, err)
	}
}
