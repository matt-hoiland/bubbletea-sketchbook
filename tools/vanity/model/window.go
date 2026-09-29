package model

import (
	"sort"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var borders = map[string]lipgloss.Border{
	"ASCII Border":            lipgloss.ASCIIBorder(),
	"Block Border":            lipgloss.BlockBorder(),
	"Double Border":           lipgloss.DoubleBorder(),
	"Hidden Border":           lipgloss.HiddenBorder(),
	"Inner Half Block Border": lipgloss.InnerHalfBlockBorder(),
	"Markdown Border":         lipgloss.MarkdownBorder(),
	"Normal Border":           lipgloss.NormalBorder(),
	"Outer Half Block Border": lipgloss.OuterHalfBlockBorder(),
	"Rounded Border":          lipgloss.RoundedBorder(),
	"Thick Border":            lipgloss.ThickBorder(),
}

type keyMap struct {
	NextBorderStyle key.Binding
	Quit            key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextBorderStyle, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextBorderStyle, k.Quit},
	}
}

var keys = keyMap{
	NextBorderStyle: key.NewBinding(
		key.WithKeys("tab", "enter"),
		key.WithHelp("tab", "next border style"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
}

type Window struct {
	keys          keyMap
	help          help.Model
	height, width int

	borderSelected int
	borderOptions  []string

	spinner spinner.Model
}

func NewWindow() Window {
	w := Window{
		keys: keys,
		help: help.New(),
		spinner: spinner.New(spinner.WithSpinner(spinner.Spinner{
			Frames: []string{
				"[.    ]",
				"[.o   ]",
				"[.oO  ]",
				"[ oOo ]",
				"[  Oo.]",
				"[   o.]",
				"[    .]",
			},
			FPS: time.Second / 8,
		})),
	}

	for style := range borders {
		w.borderOptions = append(w.borderOptions, style)
	}
	sort.Strings(w.borderOptions)

	return w
}

// Init implements [tea.Model].
func (w Window) Init() tea.Cmd {
	return w.spinner.Tick
}

// Update implements [tea.Model].
func (w Window) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return w.handleWindowSizeMsg(msg)
	case tea.KeyPressMsg:
		return w.handleKeyPressMsg(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		w.spinner, cmd = w.spinner.Update(msg)
		return w, cmd
	}

	return w, nil
}

func (w Window) handleKeyPressMsg(msg tea.KeyPressMsg) (Window, tea.Cmd) {
	switch {
	case key.Matches(msg, w.keys.Quit):
		return w, tea.Quit
	case key.Matches(msg, w.keys.NextBorderStyle):
		return w.nextBorderStyle()
	}
	return w, nil
}

func (w Window) handleWindowSizeMsg(msg tea.WindowSizeMsg) (Window, tea.Cmd) {
	w.height, w.width = msg.Height, msg.Width
	return w, nil
}

func (w Window) nextBorderStyle() (Window, tea.Cmd) {
	w.borderSelected++
	if w.borderSelected >= len(w.borderOptions) {
		w.borderSelected = 0
	}
	return w, nil
}

// View implements [tea.Model].
func (w Window) View() tea.View {
	helpView := w.help.View(w.keys)
	lipgloss.Height(helpView)

	content := lipgloss.NewStyle().
		Height(w.height - lipgloss.Height(helpView)).
		Width(w.width).
		Border(borders[w.borderOptions[w.borderSelected]]).
		BorderForeground(lipgloss.Magenta).
		Render(" " + w.borderOptions[w.borderSelected] + " " + w.spinner.View())

	v := tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		helpView,
	))
	v.AltScreen = true

	return v
}
