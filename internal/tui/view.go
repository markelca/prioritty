package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/pkg/tasks"
)

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

var successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#a6e3a1"))

var greyStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#595a70"))

var blueStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7aa0df"))

func (m Model) View() string {
	s := ""

	countToDo := 0
	countInProgress := 0
	countDone := 0

	for i, task := range m.tasks {
		cursor := " "
		if m.withTui && m.cursor == i {
			cursor = ">"
		}

		switch task.Status {
		case tasks.Done:
			countDone += 1
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				successStyle.Render("✔"),
				greyStyle.Render(task.Title))

		case tasks.Todo:
			countToDo += 1
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				defaultStyle.Render("☐"),
				task.Title)
		case tasks.InProgress:
			countInProgress += 1
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				blueStyle.Render("◐"),
				task.Title)
		}

	}

	if len(m.tasks) > 0 {

	}

	s += fmt.Sprintf("\n  %s %s",
		successStyle.Render(
			fmt.Sprintf("%.f%%", float64(countDone)/float64(len(m.tasks))*100),
		),
		greyStyle.Render("of all tasks completed."),
	)
	s += fmt.Sprintf("\n  %s %s %s %s %s %s\n",
		successStyle.Render(fmt.Sprintf("%d", countDone)),
		greyStyle.Render("done ·"),
		blueStyle.Render(fmt.Sprintf("%d", countInProgress)),
		greyStyle.Render("in-progress ·"),
		defaultStyle.Render(fmt.Sprintf("%d", countToDo)),
		greyStyle.Render("pending"),
	)

	var helpView string
	if m.withTui {
		helpView = m.help.View(m.Keys)

	}

	return s + "\n" + helpView
}
