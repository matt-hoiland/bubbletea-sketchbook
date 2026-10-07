package window

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		cmd := m.handleKeyPressMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.ToggleFullHelp):
		m.help.ShowAll = !m.help.ShowAll
	case key.Matches(msg, m.keys.ExitProgram):
		return tea.Quit
	}
	return nil
}
