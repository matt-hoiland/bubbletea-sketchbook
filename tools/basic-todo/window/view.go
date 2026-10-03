package window

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) View() tea.View {
	var (
		helpView  = m.help.View(m.keys)
		titleView = lipgloss.NewStyle().
				Width(m.width).
				Background(lipgloss.Cyan).
				Foreground(lipgloss.Black).
				Bold(true).
				Render(m.title)
		windowView = lipgloss.NewStyle().
				Width(m.width).
				Height(m.height - lipgloss.Height(titleView) - lipgloss.Height(helpView)).
				Render("Hello!")
	)

	v := tea.NewView(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleView,
			windowView,
			helpView,
		),
	)

	v.AltScreen = true

	return v
}
