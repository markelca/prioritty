package tui

import "github.com/charmbracelet/bubbles/key"

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Help       key.Binding
	Quit       key.Binding
	HardQuit   key.Binding
	MenuQuit   key.Binding
	InProgress key.Binding
	Done       key.Binding
	ToDo       key.Binding
	Cancelled  key.Binding
	Show       key.Binding
	Edit       key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.InProgress, k.ToDo, k.Done, k.Cancelled},
		{k.Show, k.Edit},
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	HardQuit: key.NewBinding(
		key.WithKeys("ctrl+c"),
	),
	MenuQuit: key.NewBinding(
		key.WithKeys("esc"),
	),
	InProgress: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "In progress"),
	),
	Done: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Done"),
	),
	ToDo: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "To do"),
	),
	Cancelled: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Cancelled"),
	),
	Show: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Show"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit"),
	),
}
