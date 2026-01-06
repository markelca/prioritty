package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/tui/styles"
	"github.com/markelca/prioritty/pkg/items"
)

var taskIcons = map[items.Status]string{
	items.Done:       styles.DoneIcon,
	items.InProgress: styles.InProgressIcon,
	items.Cancelled:  styles.CancelledIcon,
	items.Todo:       styles.TodoIcon,
}

func (m Model) View() string {
	view := ""
	counts := make(map[items.Status]int)

	if len(m.state.items) == 0 {
		view += "No items found!"
		if m.params.IsTUI {
			view += styles.Default.
				MarginTop(1).
				SetString(Help.View(keys)).
				Render()
		}
		return view
	}

	// Track current tag to print headers when it changes
	var currentTag *string
	tagStats := make(map[string]struct{ completed, total int })

	// First pass: calculate stats per tag
	for _, item := range m.state.items {
		tagKey := ""
		if tag := item.GetTag(); tag != nil {
			tagKey = tag.Name
		}
		if task, ok := item.(*items.Task); ok {
			stats := tagStats[tagKey]
			stats.total++
			if task.Status == items.Done || task.Status == items.Cancelled {
				stats.completed++
			}
			tagStats[tagKey] = stats
		}
	}

	// Second pass: render items
	for index, item := range m.state.items {
		tagKey := ""
		if tag := item.GetTag(); tag != nil {
			tagKey = tag.Name
		}

		// Print tag header when tag changes
		if currentTag == nil || *currentTag != tagKey {
			currentTag = &tagKey
			var tagName string
			if tagKey == "" {
				tagName = "My Board"
			} else {
				tagName = "@" + tagKey
			}

			view += "\n  " + styles.Default.Underline(true).Render(tagName)

			if stats := tagStats[tagKey]; stats.total > 0 {
				view += styles.Secondary.Render(fmt.Sprintf(" [%d/%d]", stats.completed, stats.total))
			}

			view += "\n"
		}

		// Count for summary
		view += "  "
		switch v := item.(type) {
		case *items.Note:
			counts[items.NoteType] += 1
		case *items.Task:
			counts[v.Status] += 1
		}

		cursor := " "
		if m.params.IsTUI && m.state.cursor == index {
			cursor = ">"
		}

		var padding string
		if len(m.state.items) >= 10 {
			padding = "2"
		} else {
			padding = "1"
		}
		view += cursor
		view += styles.Secondary.
			SetString(fmt.Sprintf(" %"+padding+"d. ", index+1)).
			Render()
		view += item.Render(m.renderer)
	}

	view += renderDonePercentage(m.state.items, counts)
	view += renderSummary(counts)

	if m.params.IsTUI {
		view += styles.Default.
			MarginTop(1).
			SetString(Help.View(keys)).
			Render()
	}

	if m.state.contentView.ready {
		view = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.state.contentView.viewport.View(), m.footerView())
	}

	// Show delete confirmation dialog
	if m.state.Mode == ModeDeleteConfirm && m.state.pendingDelete != nil {
		view += "\n" + styles.RenderDeleteDialog(m.state.pendingDelete.GetTitle())
	}

	return view
}

func renderDonePercentage(taskList []items.ItemInterface, counts map[items.Status]int) string {
	var taskCount int
	for _, t := range taskList {
		if _, ok := t.(*items.Task); ok {
			taskCount += 1
		}
	}
	var donePercentage float64
	if taskCount > 0 {
		donePercentage = float64(counts[items.Done]+counts[items.Cancelled]) / float64(taskCount) * 100
	} else {
		donePercentage = 0
	}

	var doneStyle lipgloss.Style
	if donePercentage > 0 {
		doneStyle = styles.Done
	} else {
		doneStyle = styles.Secondary
	}

	return fmt.Sprintf("\n  %s %s",
		doneStyle.Render(
			fmt.Sprintf("%.f%%", donePercentage),
		),
		styles.Secondary.Render("of all tasks completed."),
	)
}

func renderSummary(counts map[items.Status]int) string {
	return fmt.Sprintf("\n  %s %s %s %s %s %s %s %s %s %s\n",
		styles.Done.Render(fmt.Sprintf("%d", counts[items.Done])),
		styles.Secondary.Render("done ·"),
		styles.InProgress.Render(fmt.Sprintf("%d", counts[items.InProgress])),
		styles.Secondary.Render("in-progress ·"),
		styles.Default.Render(fmt.Sprintf("%d", counts[items.Todo])),
		styles.Secondary.Render("pending ·"),
		styles.Cancelled.Render(fmt.Sprintf("%d", counts[items.Cancelled])),
		styles.Secondary.Render("cancelled ·"),
		styles.InProgress.Render(fmt.Sprintf("%d", counts[items.NoteType])),
		styles.Secondary.Render("notes"),
	)
}

func GetItemIcon(item items.ItemInterface) string {
	var icon string
	if t, ok := item.(*items.Task); ok {
		icon = taskIcons[t.Status]
	} else if _, ok := item.(*items.Note); ok {
		icon = styles.NoteIcon
	}
	return icon
}

func (m Model) headerView() string {
	item := m.state.GetCurrentItem()
	if item == nil {
		return ""
	}
	icon := GetItemIcon(item)
	title := styles.TitleStyle.Render(icon + item.GetTitle())
	line := strings.Repeat("─", max(0, m.state.contentView.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.state.contentView.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.state.contentView.viewport.Width-lipgloss.Width(info)))
	footer := lipgloss.JoinHorizontal(lipgloss.Center, line, info, "\n")
	return footer
}
