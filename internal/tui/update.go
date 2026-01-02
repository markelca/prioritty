package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/editor"
	"github.com/markelca/prioritty/pkg/items"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	item := m.state.GetCurrentItem()

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Help):
			Help.ShowAll = !Help.ShowAll

		case key.Matches(msg, keys.Quit):
			if m.state.item.ready {
				m.state.item.ready = false
			} else {
				return m, tea.Quit
			}

		case key.Matches(msg, keys.MenuQuit):
			m.state.item.ready = false

		case key.Matches(msg, keys.HardQuit):
			return m, tea.Quit

		case key.Matches(msg, keys.Up),
			key.Matches(msg, keys.Down):
			m.move(msg)

		case key.Matches(msg, keys.InProgress),
			key.Matches(msg, keys.ToDo),
			key.Matches(msg, keys.Done),
			key.Matches(msg, keys.Cancelled):
			m.updateStatus(msg, item)

		case key.Matches(msg, keys.Show):
			m.state.item.show(item)
		case key.Matches(msg, keys.Edit):
			msg, err := m.Service.EditWithEditor(item)
			if err != nil {
				log.Println(err)
			}
			return m, msg
		}
	case tea.WindowSizeMsg:
		m.state.item.init(ItemContentDimensions{
			width:        msg.Width,
			height:       msg.Height,
			headerHeight: lipgloss.Height(m.headerView()),
			footerHeight: lipgloss.Height(m.footerView()),
		})

	case editor.TaskEditorFinishedMsg:
		// Check if the editor operation was cancelled (no content)
		if msg.Err != nil {
			if m.params.CreateMode != "" {
				// Creation mode - just quit without creating
				return m, tea.Quit
			} else if m.params.EditMode {
				// Standalone edit mode - quit without updating
				return m, tea.Quit
			} else {
				// Interactive TUI mode - return to list without updating
				return m, tea.ClearScreen
			}
		}

		if m.params.CreateMode != "" {
			// Creation mode
			var err error
			if m.params.CreateMode == "task" {
				err = m.Service.CreateTaskFromEditorMsg(msg)
			} else if m.params.CreateMode == "note" {
				err = m.Service.CreateNoteFromEditorMsg(msg)
			}
			if err != nil {
				log.Println("Error creating item:", err)
			}
			return m, tea.Quit
		} else {
			// Edit mode
			t := m.state.GetCurrentItem()
			err := m.Service.UpdateItemFromEditorMsg(t, msg)
			if err != nil {
				log.Println("Error - ", err)
			}
			if m.params.EditMode {
				// Standalone edit mode - quit after editing
				return m, tea.Quit
			} else {
				// Interactive TUI mode - return to list
				return m, tea.ClearScreen
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	// return m, nil
	m.state.item.viewport, cmd = m.state.item.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) move(msg tea.KeyMsg) {
	style := lipgloss.NewStyle().Width(m.state.item.viewport.Width)
	switch {
	case key.Matches(msg, keys.Up):
		if m.state.cursor == 0 {
			m.state.cursor = len(m.state.items) - 1
		} else {
			m.state.cursor--
		}
	case key.Matches(msg, keys.Down):
		if m.state.cursor == len(m.state.items)-1 {
			m.state.cursor = 0
		} else {
			m.state.cursor++
		}
	}
	item := m.state.GetCurrentItem()
	content := style.Render(item.GetBody())
	m.state.item.viewport.SetContent(content)
}

func (m *Model) updateStatus(msg tea.KeyMsg, item items.ItemInterface) error {
	var (
		task *items.Task
		ok   bool
		s    items.Status
	)

	if task, ok = item.(*items.Task); !ok {
		return nil
	}

	switch {
	case key.Matches(msg, keys.InProgress):
		s = items.InProgress
	case key.Matches(msg, keys.ToDo):
		s = items.Todo
	case key.Matches(msg, keys.Cancelled):
		s = items.Cancelled
	case key.Matches(msg, keys.Done):
		s = items.Done
	default:
		return fmt.Errorf("Error - the message cannot be mapped to a task status (%v)", task)
	}

	m.Service.UpdateStatus(task, s)
	return nil
}
