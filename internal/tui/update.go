package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/pkg/tasks"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ", "d":
			task := &m.tasks[m.cursor]
			if task.Status == tasks.Done {
				task.Status = tasks.Todo
			} else {
				task.Status = tasks.Done
			}
		case "p":
			task := &m.tasks[m.cursor]
			if task.Status == tasks.InProgress {
				task.Status = tasks.Todo
			} else {
				task.Status = tasks.InProgress
			}

		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
