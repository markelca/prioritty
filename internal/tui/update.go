package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/pkg/tasks"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	task := &m.state.tasks[m.state.cursor]
	// task := &m.state.tasks[m.state.cursor]
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Help):
			Help.ShowAll = !Help.ShowAll

		case key.Matches(msg, keys.Quit):
			if m.state.taskContent.ready {
				m.state.taskContent.ready = false
			} else {
				return m, tea.Quit
			}

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
		case key.Matches(msg, keys.Show):
			width := m.state.taskContent.viewport.Width
			content := lipgloss.NewStyle().Width(width).Render(task.Body)

			m.state.taskContent.viewport.SetContent(content)

			m.state.taskContent.ready = true

		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.state.taskContent.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.state.taskContent.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.state.taskContent.viewport.YPosition = headerHeight
			// m.state.taskContent.viewport.SetContent(m.state.taskContent.content)
		} else {
			m.state.taskContent.viewport.Width = msg.Width
			m.state.taskContent.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	// return m, nil
	m.state.taskContent.viewport, cmd = m.state.taskContent.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
