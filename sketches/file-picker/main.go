package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func main() {
	tm, _ := tea.NewProgram(newModel()).Run()
	mm := tm.(model)
	fmt.Println("\n  You selected: " + mm.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
}

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	help         help.Model
	quitting     bool
	err          error
}

type clearErrorMsg struct{}

type keyMap filepicker.KeyMap

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Back}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Back},
		{k.PageUp, k.PageDown, k.GoToTop, k.GoToLast},
		{k.Select, k.Open,
			key.NewBinding(
				key.WithKeys("ctrl+c", "esc"),
				key.WithHelp("ctrl+c/esc", "exit program"),
			)},
	}
}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func newModel() model {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	h := help.New()
	h.ShowAll = true

	return model{
		filepicker: fp,
		help:       h,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.filepicker.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.filepicker, cmd = m.filepicker.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	var s strings.Builder
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: ")
		s.WriteString(m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n")
	s.WriteString(m.help.View(keyMap(m.filepicker.KeyMap)))
	s.WriteString("\n\n")
	s.WriteString(m.filepicker.View())
	s.WriteString("\n")

	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}
