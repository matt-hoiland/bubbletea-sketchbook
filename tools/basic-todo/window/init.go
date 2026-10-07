package window

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.helpInit(),
	)
}

func (m *Model) helpInit() tea.Cmd {
	s := m.help.Styles

	var (
		keyStyle = s.ShortKey.
				Foreground(lipgloss.BrightYellow).
				Bold(true)
		descStyle = s.ShortKey.
				Foreground(lipgloss.Yellow)
		sepStyle = s.ShortSeparator.
				Foreground(lipgloss.Black)
	)

	m.help.Styles = help.Styles{
		ShortKey:       keyStyle,
		ShortDesc:      descStyle,
		ShortSeparator: sepStyle,
		Ellipsis:       sepStyle,
		FullKey:        keyStyle,
		FullDesc:       descStyle,
		FullSeparator:  sepStyle,
	}

	return nil
}
