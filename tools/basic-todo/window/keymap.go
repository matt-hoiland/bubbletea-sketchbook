package window

import "charm.land/bubbles/v2/key"

type keyMap struct {
	ToggleFullHelp key.Binding
	ExitProgram    key.Binding
}

var keys = keyMap{
	ToggleFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	ExitProgram: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "exit"),
	),
}

// FullHelp implements [help.KeyMap].
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ExitProgram, k.ToggleFullHelp},
		{k.ExitProgram, k.ToggleFullHelp},
	}

}

// ShortHelp implements [help.KeyMap].
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ExitProgram, k.ToggleFullHelp}
}
