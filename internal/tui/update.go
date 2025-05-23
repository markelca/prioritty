package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/pkg/tasks"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Help       key.Binding
	Quit       key.Binding
	Check      key.Binding
	InProgress key.Binding
	Done       key.Binding
	ToDo       key.Binding
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
		{k.Check, k.InProgress, k.ToDo, k.Done},
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
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Check: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Check/uncheck task"),
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
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.Keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.Keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.Keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, m.Keys.Down):
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case key.Matches(msg, m.Keys.Check):
			task := &m.tasks[m.cursor]
			if task.Status == tasks.Done {
				task.Status = tasks.Todo
			} else {
				task.Status = tasks.Done
			}

		case key.Matches(msg, m.Keys.InProgress):
			task := &m.tasks[m.cursor]
			if task.Status == tasks.InProgress {
				task.Status = tasks.Todo
			} else {
				task.Status = tasks.InProgress
			}

		case key.Matches(msg, m.Keys.ToDo):
			task := &m.tasks[m.cursor]
			task.Status = tasks.Todo

		case key.Matches(msg, m.Keys.Done):
			task := &m.tasks[m.cursor]
			task.Status = tasks.Done
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
