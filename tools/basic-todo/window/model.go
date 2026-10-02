package window

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type keyMap struct {
	ExitProgram key.Binding
}

// FullHelp implements [help.KeyMap].
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.ExitProgram, k.ExitProgram}}

}

// ShortHelp implements [help.KeyMap].
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ExitProgram, k.ExitProgram}
}

var keys = keyMap{
	ExitProgram: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q, esc, ctrl+c :", "exit"),
	),
}

type Model struct {
	keys keyMap
	help help.Model
}

func New() *Model {
	m := &Model{
		keys: keys,
		help: help.New(),
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.ExitProgram):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	return tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			"Hello!",
			m.help.View(m.keys),
		),
	)
}
