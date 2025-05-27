package tui

import (
	"fmt"

	"github.com/markelca/prioritty/internal/tui/styles"
	"github.com/markelca/prioritty/pkg/tasks"
)

func (m Model) View() string {
	s := ""

	countToDo := 0
	countInProgress := 0
	countDone := 0
	countCancelled := 0

	for i, task := range m.tasks {
		cursor := " "
		if m.withTui && m.cursor == i {
			cursor = ">"
		}

		switch task.Status {
		case tasks.Done:
			countDone += 1
			s += fmt.Sprintf("%s %s %s %s\n",
				styles.Default.Render(cursor),
				styles.Secondary.Render(fmt.Sprintf("%d.", i+1)),
				styles.Done.Render("✔"),
				styles.Secondary.Render(task.Title))

		case tasks.Todo:
			countToDo += 1
			s += fmt.Sprintf("%s %s %s %s\n",
				styles.Default.Render(cursor),
				styles.Secondary.Render(fmt.Sprintf("%d.", i+1)),
				styles.Default.Render("☐"),
				task.Title)
		case tasks.InProgress:
			countInProgress += 1
			s += fmt.Sprintf("%s %s %s %s\n",
				styles.Default.Render(cursor),
				styles.Secondary.Render(fmt.Sprintf("%d.", i+1)),
				styles.InProgress.Render("◐"),
				task.Title)
		case tasks.Cancelled:
			countCancelled += 1
			s += fmt.Sprintf("%s %s %s %s\n",
				styles.Default.Render(cursor),
				styles.Secondary.Render(fmt.Sprintf("%d.", i+1)),
				styles.Cancelled.Render("x"),
				styles.Secondary.Render(task.Title))
		}

	}

	var donePercentage float64
	if len(m.tasks) > 0 {
		donePercentage = float64(countDone+countCancelled) / float64(len(m.tasks)) * 100
	} else {
		donePercentage = 0
	}

	s += fmt.Sprintf("\n  %s %s",
		styles.Done.Render(
			fmt.Sprintf("%.f%%", donePercentage),
		),
		styles.Secondary.Render("of all tasks completed."),
	)
	s += fmt.Sprintf("\n  %s %s %s %s %s %s %s %s\n",
		styles.Done.Render(fmt.Sprintf("%d", countDone)),
		styles.Secondary.Render("done ·"),
		styles.InProgress.Render(fmt.Sprintf("%d", countInProgress)),
		styles.Secondary.Render("in-progress ·"),
		styles.Default.Render(fmt.Sprintf("%d", countToDo)),
		styles.Secondary.Render("pending ·"),
		styles.Cancelled.Render(fmt.Sprintf("%d", countCancelled)),
		styles.Secondary.Render("cancelled"),
	)

	var helpView string
	if m.withTui {
		helpView = m.help.View(m.keys)

	}

	return s + "\n" + helpView
}
