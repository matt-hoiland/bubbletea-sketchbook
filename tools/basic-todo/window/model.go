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

	width, height int

	title string
}

type Option func(*conf)

type conf struct {
	title string
}

func New(opts ...Option) *Model {
	c := conf{
		title: "program",
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&c)
		}
	}

	m := &Model{
		keys: keys,
		help: help.New(),

		title: c.title,
	}

	return m
}

func WithTitle(title string) Option {
	return func(c *conf) {
		c.title = title
	}
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
	case tea.WindowSizeMsg:
		m.height, m.width = msg.Height, msg.Width
	}
	return m, nil
}

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
