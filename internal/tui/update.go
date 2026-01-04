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
			if m.state.contentView.ready {
				m.state.contentView.ready = false
			} else {
				return m, tea.Quit
			}

		case key.Matches(msg, keys.MenuQuit):
			m.state.contentView.ready = false

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
			m.state.contentView.show(item)
		case key.Matches(msg, keys.Edit):
			m.state.Mode = ModeEdit
			cmd, err := m.Service.EditWithEditor(item)
			if err != nil {
				log.Println(err)
			}
			return m, cmd
		case key.Matches(msg, keys.Add):
			m.state.Mode = ModeCreate
			cmd, err := m.Service.AddWithEditor(items.ItemTypeTask)
			if err != nil {
				log.Println(err)
			}
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.state.contentView.init(ItemContentDimensions{
			width:        msg.Width,
			height:       msg.Height,
			headerHeight: lipgloss.Height(m.headerView()),
			footerHeight: lipgloss.Height(m.footerView()),
		})

	case editor.EditorFinishedMsg:
		// Check if the editor operation was cancelled (no content)
		if msg.Err != nil {
			if !m.params.IsTUI {
				// CLI mode - just quit
				return m, tea.Quit
			}
			// TUI mode - return to list
			m.state.Mode = ModeList
			return m, tea.ClearScreen
		}

		switch m.state.Mode {
		case ModeCreate:
			// Create the item using the type from editor (allows user to change type field)
			var err error
			if msg.ItemType == items.ItemTypeTask {
				err = m.Service.CreateTaskFromEditorMsg(msg)
			} else if msg.ItemType == items.ItemTypeNote {
				err = m.Service.CreateNoteFromEditorMsg(msg)
			}
			if err != nil {
				log.Println("Error creating item:", err)
			}
		case ModeEdit:
			// Update the existing item
			item := m.state.GetCurrentItem()
			err := m.Service.UpdateItemFromEditorMsg(item, msg)
			if err != nil {
				log.Println("Error updating item:", err)
			}
		}

		// After action: quit (CLI) or return to list (TUI)
		if !m.params.IsTUI {
			return m, tea.Quit
		}
		m.state.Mode = ModeList
		m.refreshItems()
		return m, tea.ClearScreen
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	// return m, nil
	m.state.contentView.viewport, cmd = m.state.contentView.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) move(msg tea.KeyMsg) {
	style := lipgloss.NewStyle().Width(m.state.contentView.viewport.Width)
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
	m.state.contentView.viewport.SetContent(content)
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
