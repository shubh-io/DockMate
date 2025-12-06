package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// ============================================================================
// Keyboard shortcuts
// ============================================================================

// all the keybindings
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Start    key.Binding
	Stop     key.Binding
	Restart  key.Binding
	Logs     key.Binding
	Exec     key.Binding
	Remove   key.Binding
	Refresh  key.Binding
	PageUp   key.Binding
	NextPage key.Binding
	PrevPage key.Binding
	Terminal key.Binding
	PageDown key.Binding
	Quit     key.Binding
	Help     key.Binding
}

// global keymap
// supports vim keys (hjkl) and arrows
var keys = keyMap{
	Up:       key.NewBinding(key.WithKeys("up", "k")),
	Down:     key.NewBinding(key.WithKeys("down", "j")),
	Start:    key.NewBinding(key.WithKeys("s")),
	Stop:     key.NewBinding(key.WithKeys("x")),
	Logs:     key.NewBinding(key.WithKeys("l")),
	Exec:     key.NewBinding(key.WithKeys("e")),
	Restart:  key.NewBinding(key.WithKeys("r")),
	Remove:   key.NewBinding(key.WithKeys("d")),
	Refresh:  key.NewBinding(key.WithKeys("f5", "R")),
	PageUp:   key.NewBinding(key.WithKeys("pgup", "left")),
	NextPage: key.NewBinding(key.WithKeys("n", "pagedown")),
	PrevPage: key.NewBinding(key.WithKeys("p", "pageup")),
	Terminal: key.NewBinding(key.WithKeys("t")),
	PageDown: key.NewBinding(key.WithKeys("pgdown", "right")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c", "f10")),
	Help:     key.NewBinding(key.WithKeys("f1", "?")),
}
