package window

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) View() tea.View {
	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.titleView(),
			m.windowView(),
			m.helpView(),
		),
	)

	v.AltScreen = true

	return v
}

func (m *Model) helpView() string {
	return lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.NormalBorder(), true, false, false).
		BorderForeground(lipgloss.Color("#e48e1e")).
		Padding(0, 1).
		Render(m.help.View(m.keys))
}

func (m *Model) titleView() string {
	return lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1).
		Background(lipgloss.Color("#8c00ff")).
		Bold(true).
		Render(m.title)
}

func (m *Model) windowView() string {
	return lipgloss.NewStyle().
		Height(m.height-lipgloss.Height(m.titleView())-lipgloss.Height(m.helpView())).
		Width(m.width).
		Padding(0, 1).
		Render("Hello!")
}
