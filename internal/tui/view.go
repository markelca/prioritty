package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/pkg/tasks"
)

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

var successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#04B575"))

var greyStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#595a70"))

var blueStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7aa0df"))

func (m Model) View() string {
	s := ""

	for i, task := range m.tasks {
		cursor := " "
		if m.withTui && m.cursor == i {
			cursor = ">"
		}

		switch task.Status {
		case tasks.Done:
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				successStyle.Render("✔"),
				greyStyle.Render(task.Title))

		case tasks.Todo:
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				defaultStyle.Render("☐"),
				task.Title)
		case tasks.InProgress:
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				blueStyle.Render("◐"),
				task.Title)
		}

	}

	var helpView string
	if m.withTui {
		helpView = m.help.View(m.Keys)

	}

	return s + "\n" + helpView
}
