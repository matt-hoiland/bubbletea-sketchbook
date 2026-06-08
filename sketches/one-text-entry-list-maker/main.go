package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		slog.Error("bubbletea runtime error", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

type model struct {
	textInput textinput.Model
	items     []string
	err       error
	quitting  bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return model{
		textInput: ti,
	}
}

// Init implements [tea.Model].
func (model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements [tea.Model].
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.items = append(m.items, m.textInput.Value())
			m.textInput.SetValue("")
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View implements [tea.Model].
func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.itemsView()))
	}

	str := lipgloss.JoinVertical(
		lipgloss.Top,
		m.headerView(),
		m.itemsView(),
		m.textInput.View(),
		m.footerView(),
	)

	if m.quitting {
		str += "\n"
	}

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m model) headerView() string { return "What's  your favorite Pokémon?\n" }
func (m model) footerView() string { return "\n(esc to quit)" }

func (m model) itemsView() string {
	var bob strings.Builder

	fmt.Fprintf(&bob, "-----------------\n")

	for _, item := range m.items {
		fmt.Fprintf(&bob, "* %s\n", item)
	}

	fmt.Fprintf(&bob, "-----------------\n")

	return bob.String()
}
