package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/tui/styles"
	"github.com/markelca/prioritty/pkg/tasks"
)

func (m Model) View() string {
	view := ""
	counts := make(map[tasks.Status]int)

	for i, task := range m.state.tasks {
		counts[task.Status] += 1
		cursor := " "
		if m.params.withTui && m.state.cursor == i {
			cursor = ">"
		}

		view += cursor
		view += styles.Secondary.
			SetString(fmt.Sprintf(" %d. ", i+1)).
			Render()
		view += renderTask(task)
	}

	view += renderDonePercentage(m.state.tasks, counts)
	view += renderSummary(counts)

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

func renderTask(t tasks.Task) string {
	var title string
	var icon string
	var contentIcon string
	var style lipgloss.Style

	switch t.Status {
	case tasks.Done:
		icon = styles.DoneIcon
		style = styles.DoneTitle
	case tasks.Cancelled:
		icon = styles.CancelledIcon
		style = styles.DoneTitle
	case tasks.InProgress:
		icon = styles.InProgressIcon
		style = styles.Default
	case tasks.Todo:
		icon = styles.TodoIcon
		style = styles.Default
	}

	if len(t.Body) > 1 {
		contentIcon = styles.ContentIcon
	}

	title += style.
		// PaddingBottom(1).
		Render(t.Title)

	return icon + contentIcon + title + "\n"
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

func (m Model) headerView() string {
	task := &m.state.tasks[m.state.cursor]
	title := styles.TitleStyle.Render(task.Title)
	line := strings.Repeat("─", max(0, m.state.taskContent.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.state.taskContent.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.state.taskContent.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
