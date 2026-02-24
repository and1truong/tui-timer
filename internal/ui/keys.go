package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Toggle    key.Binding
	Reset     key.Binding
	Skip      key.Binding
	Quit      key.Binding
	Config    key.Binding
	TimeUp    key.Binding
	TimeDown  key.Binding
	TimeRight key.Binding
	TimeLeft  key.Binding
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
		TimeUp: key.NewBinding(
			key.WithKeys("shift+up"),
			key.WithHelp("shift+↑", "+1 min"),
		),
		TimeDown: key.NewBinding(
			key.WithKeys("shift+down"),
			key.WithHelp("shift+↓", "-1 min"),
		),
		TimeRight: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+→", "+10 min"),
		),
		TimeLeft: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+←", "-10 min"),
		),
	}
}
