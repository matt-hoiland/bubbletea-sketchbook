package main

import (
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct{}

func newModel() *model {
	return &model{}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl-c", "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() tea.View {
	return tea.NewView("Hello! Press q, esc, or ctrl-c to exit.")
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		slog.Error("program failed to start", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
