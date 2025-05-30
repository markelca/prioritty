package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/pkg/tasks"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	task := &m.state.tasks[m.state.cursor]
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Help):
			Help.ShowAll = !Help.ShowAll

		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.state.cursor > 0 {
				m.state.cursor--
			}

		case key.Matches(msg, keys.Down):
			if m.state.cursor < len(m.state.tasks)-1 {
				m.state.cursor++
			}

		case key.Matches(msg, keys.InProgress):
			m.Service.UpdateStatus(task, tasks.InProgress)

		case key.Matches(msg, keys.ToDo):
			m.Service.UpdateStatus(task, tasks.Todo)

		case key.Matches(msg, keys.Done),
			key.Matches(msg, keys.Check):
			m.Service.UpdateStatus(task, tasks.Done)

		case key.Matches(msg, keys.Cancelled):
			m.Service.UpdateStatus(task, tasks.Cancelled)
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
