package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/pkg/editor"
	"github.com/markelca/prioritty/pkg/items"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	item := m.state.GetCurrentItem()
	if item == nil {
		log.Fatal("nthn")
	}

	contentStyle := lipgloss.NewStyle().Width(m.state.taskContent.viewport.Width)
	// fmt.Println(m.state.cursor)

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

		case key.Matches(msg, keys.MenuQuit):
			if m.state.taskContent.ready {
				m.state.taskContent.ready = false
			}

		case key.Matches(msg, keys.HardQuit):
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.state.cursor == 0 {
				m.state.cursor = len(m.state.items) - 1
			} else {
				m.state.cursor--
			}
			item = m.state.GetCurrentItem()
			content := contentStyle.Render(item.GetBody())
			m.state.taskContent.viewport.SetContent(content)

		case key.Matches(msg, keys.Down):
			if m.state.cursor == len(m.state.items)-1 {
				m.state.cursor = 0
			} else {
				m.state.cursor++
			}
			item = m.state.GetCurrentItem()
			content := contentStyle.Render(item.GetBody())
			m.state.taskContent.viewport.SetContent(content)

		case key.Matches(msg, keys.InProgress):
			if t, ok := item.(*items.Task); ok {
				if !m.state.taskContent.ready {
					m.Service.UpdateStatus(t, items.InProgress)
				}
			}

		case key.Matches(msg, keys.ToDo):
			if t, ok := item.(*items.Task); ok {
				m.Service.UpdateStatus(t, items.Todo)
			}

		case key.Matches(msg, keys.Done):
			var t *items.Task
			var ok bool
			if t, ok = item.(*items.Task); ok {
				if !m.state.taskContent.ready {
					m.Service.UpdateStatus(t, items.Done)
				}
			}

		case key.Matches(msg, keys.Cancelled):
			if t, ok := item.(*items.Task); ok {
				if !m.state.taskContent.ready {
					m.Service.UpdateStatus(t, items.Cancelled)
				}
			}

		case key.Matches(msg, keys.Show):
			if m.state.taskContent.ready {
				m.state.taskContent.ready = false
			} else {
				body := (m.state.GetCurrentItem()).GetBody()
				content := contentStyle.Render(body)
				m.state.taskContent.viewport.SetContent(content)
				m.state.taskContent.ready = true
			}
		case key.Matches(msg, keys.Edit):
			// msg, err := m.Service.EditWithEditor(item)
			// if err != nil {
			// fmt.Println(err)
			// }
			// return m, msg
		}
	case tea.WindowSizeMsg:
		// headerHeight := lipgloss.Height(m.headerView())
		// footerHeight := lipgloss.Height(m.footerView())
		// verticalMarginHeight := headerHeight + footerHeight

		if !m.state.taskContent.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			// m.state.taskContent.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			// m.state.taskContent.viewport.YPosition = headerHeight
		} else {
			// m.state.taskContent.viewport.Width = msg.Width
			// m.state.taskContent.viewport.Height = msg.Height - verticalMarginHeight
		}
	case editor.TaskEditorFinishedMsg:
		// t := m.state.GetCurrentItem()
		// m.Service.UpdateTaskFromEditorMsg(t, msg)
		return m, tea.ClearScreen
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	// return m, nil
	m.state.taskContent.viewport, cmd = m.state.taskContent.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
