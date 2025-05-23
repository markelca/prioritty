package tui

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/pkg/tasks"
)

type Model struct {
	withTui    bool
	Keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	tasks      []tasks.Task     // items on the to-do list
	cursor     int              // which to-do list item our cursor is pointing at
	done       map[int]struct{} // which to-do items are checked
}

func InitialModel(withTui bool) Model {
	return Model{
		withTui:    withTui,
		Keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		// Our to-do list is a grocery list
		tasks: []tasks.Task{
			{Title: "Build TUI", Status: tasks.Done},
			{Title: "Add commands", Status: tasks.InProgress},
			{Title: "Organize packages", Status: tasks.Todo},
		},
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
