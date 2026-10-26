package main

import (
	"io"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/davecgh/go-spew/spew"
)

type Model struct {
	dump io.Writer

	height, width int
}

func New() Model {
	return Model{}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	quitCopy := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Height(m.height).
		Width(m.width).
		Render("ctrl+c to exit")
	v := tea.NewView(quitCopy)
	v.AltScreen = true
	return v
}
