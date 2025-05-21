package tui

import tea "github.com/charmbracelet/bubbletea"

type model struct {
	tasks  []string         // items on the to-do list
	cursor int              // which to-do list item our cursor is pointing at
	done   map[int]struct{} // which to-do items are checked
}

func InitialModel() model {
	return model{
		// Our to-do list is a grocery list
		tasks: []string{"Build TUI", "Add commands", "Organize packages"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		done: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
