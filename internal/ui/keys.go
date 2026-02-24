package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Toggle key.Binding
	Reset  key.Binding
	Skip   key.Binding
	Quit   key.Binding
	Config key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "start/pause"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		Skip: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "skip"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Config: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "config"),
		),
	}
}
