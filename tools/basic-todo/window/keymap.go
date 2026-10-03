package window

import "charm.land/bubbles/v2/key"

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
