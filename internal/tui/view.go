package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/elliotchance/orderedmap/v3"
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
		if m.params.withTui {
			view += styles.Default.
				MarginTop(1).
				SetString(Help.View(keys)).
				Render()
		}
		return view
	}

	itemsByTag := orderedmap.NewOrderedMap[items.Tag, []items.ItemInterface]()
	nullTag := items.Tag{}

	for _, item := range m.state.items {
		tag := item.GetTag()
		if tag != nil {
			items, _ := itemsByTag.Get(*tag)
			itemsByTag.Set(*item.GetTag(), append(items, item))
		} else {
			items, _ := itemsByTag.Get(nullTag)
			itemsByTag.Set(nullTag, append(items, item))
		}
	}
	var index int
	myBoardItems, ok := itemsByTag.Get(nullTag)
	if ok {
		view += "\n  " + styles.Default.Underline(true).Render("My Board") + "\n"
		view, index = m.renderItemList(myBoardItems, counts, view, index)
	}

	for tag, itemList := range itemsByTag.AllFromFront() {
		if tag == nullTag {
			continue
		}
		tagName := "@" + tag.Name
		view += "\n  " + styles.Default.Underline(true).Render(tagName) + "\n"

		view, index = m.renderItemList(itemList, counts, view, index)
	}

	view += renderDonePercentage(m.state.items, counts)
	view += renderSummary(counts)

	if m.params.withTui {
		view += styles.Default.
			MarginTop(1).
			SetString(Help.View(keys)).
			Render()
	}

	if m.state.item.ready {
		view = fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.state.item.viewport.View(), m.footerView())
	}

	return view
}

func (m Model) renderItemList(list []items.ItemInterface, counts map[items.Status]int, view string, index int) (string, int) {
	for _, item := range list {
		view += "  "
		switch v := item.(type) {
		case *items.Note:
			counts[items.NoteType] += 1
		case *items.Task:
			counts[v.Status] += 1
		}
		cursor := " "

		if m.params.withTui && m.state.cursor == index {
			cursor = ">"
		}

		if m.params.withTui && m.state.cursor == index {
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

		index += 1
	}
	return view, index
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

func (m Model) headerView() string {
	item := m.state.GetCurrentItem()
	if item == nil {
		return ""
	}
	var icon string
	if t, ok := item.(*items.Task); ok {
		icon = taskIcons[t.Status]
	} else if _, ok := item.(*items.Note); ok {
		icon = styles.NoteIcon
	}
	title := styles.TitleStyle.Render(icon + item.GetTitle())
	line := strings.Repeat("─", max(0, m.state.item.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", m.state.item.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.state.item.viewport.Width-lipgloss.Width(info)))
	footer := lipgloss.JoinHorizontal(lipgloss.Center, line, info, "\n")
	return footer
}
