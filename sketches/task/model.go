package main

import (
	"io"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	catppuccingo "github.com/catppuccin/go"
	"github.com/davecgh/go-spew/spew"
)

var Mocha = catppuccingo.Mocha

type bounds struct {
	min int
	max int
}

type dimensions struct {
	height bounds
	width  bounds
}

type Model struct {
	dump io.Writer

	window tea.WindowSizeMsg

	dimensions
}

type cfg struct {
	dimensions
}

type Option func(*cfg)

func WithHeight(height int) Option {
	return func(c *cfg) {
		c.dimensions.height.max = height
	}
}

func WithWidth(width int) Option {
	return func(c *cfg) {
		c.dimensions.width.max = width
	}
}

func New(opts ...Option) Model {
	var c cfg
	for _, opt := range opts {
		opt(&c)
	}

	return Model{
		dimensions: c.dimensions,
	}
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
		m.window = msg
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	content := lipgloss.NewStyle().
		Width(m.dimensions.width.max).
		Height(m.height.max).
		Render(`[ ] Complete this task!
 | - A second line of text
 | - A third line of text`)
	return m.boxView(content)
}

func (m Model) boxView(content string) tea.View {
	v := tea.NewView(
		m.workspaceStyle(
			m.boxStyle(content),
		),
	)
	v.AltScreen = true
	return v
}

func (m Model) boxStyle(content string) string {
	var (
		borderColor = Mocha.Surface1()
	)

	return lipgloss.NewStyle().
		Border(lipgloss.BlockBorder()).
		BorderForeground(borderColor).
		Margin(0, 1).
		MarginBackground(borderColor).
		Render(content)

}

func (m Model) workspaceStyle(content string) string {
	return lipgloss.NewStyle().
		Background(catppuccingo.Mocha.Crust()).
		Align(lipgloss.Center, lipgloss.Center).
		Height(m.window.Height).
		Width(m.window.Width).
		Render(content)
}
