package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/tui/styles"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/tasks"
)

func (m Model) View() string {
	view := ""
	// counts := make(map[tasks.Status]int)

	if len(m.state.tasks) == 0 {
		view += "No tasks found!"
		if m.params.withTui {
			view += styles.Default.
				MarginTop(1).
				SetString(Help.View(keys)).
				Render()
		}
		return view
	}

	for _, item := range m.state.items {
		view += item.Render(m.renderer)
	}

	// for i, task := range m.state.tasks {
	// 	counts[task.Status] += 1
	// 	cursor := " "
	// 	if m.params.withTui && m.state.cursor == i {
	// 		cursor = ">"
	// 	}
	//
	// 	var padding string
	// 	if len(m.state.tasks) >= 10 {
	// 		padding = "2"
	// 	} else {
	// 		padding = "1"
	// 	}
	// 	view += cursor
	// 	view += styles.Secondary.
	// 		SetString(fmt.Sprintf(" %"+padding+"d. ", i+1)).
	// 		Render()
	// 	view += renderTask(task)
	// }

	// view += renderDonePercentage(m.state.tasks, counts)
	// view += renderSummary(counts)

	if m.params.withTui {
		view += styles.Default.
			MarginTop(1).
			SetString(Help.View(keys)).
			Render()
	}

	if m.state.taskContent.ready {
		view = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.state.taskContent.viewport.View(), m.footerView())
	}

	return view
}

func renderDonePercentage(taskList []tasks.Task, counts map[tasks.Status]int) string {
	var donePercentage float64
	if len(taskList) > 0 {
		donePercentage = float64(counts[tasks.Done]+counts[tasks.Cancelled]) / float64(len(taskList)) * 100
	} else {
		donePercentage = 0
	}

	return fmt.Sprintf("\n  %s %s",
		styles.Done.Render(
			fmt.Sprintf("%.f%%", donePercentage),
		),
		styles.Secondary.Render("of all tasks completed."),
	)
}

func renderSummary(counts map[tasks.Status]int) string {
	return fmt.Sprintf("\n  %s %s %s %s %s %s %s %s\n",
		styles.Done.Render(fmt.Sprintf("%d", counts[tasks.Done])),
		styles.Secondary.Render("done ·"),
		styles.InProgress.Render(fmt.Sprintf("%d", counts[tasks.InProgress])),
		styles.Secondary.Render("in-progress ·"),
		styles.Default.Render(fmt.Sprintf("%d", counts[tasks.Done])),
		styles.Secondary.Render("pending ·"),
		styles.Cancelled.Render(fmt.Sprintf("%d", counts[tasks.Cancelled])),
		styles.Secondary.Render("cancelled"),
	)
}

var taskIcons = map[items.Status]string{
	items.Done:       styles.DoneIcon,
	items.InProgress: styles.InProgressIcon,
	items.Cancelled:  styles.CancelledIcon,
	items.Todo:       styles.TodoIcon,
}

func (m Model) headerView() string {
	task := m.state.GetCurrentTask()
	icon := taskIcons[task.Status]
	title := styles.TitleStyle.Render(icon + task.Title)
	line := strings.Repeat("─", max(0, m.state.taskContent.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.state.taskContent.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.state.taskContent.viewport.Width-lipgloss.Width(info)))
	footer := lipgloss.JoinHorizontal(lipgloss.Center, line, info, "\n")
	return footer
}
