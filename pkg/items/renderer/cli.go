package renderer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/tui/styles"
	"github.com/markelca/prioritty/pkg/items"
)

type CliRendererer struct {
}

var taskIcons = map[items.Status]string{
	items.Done:       styles.DoneIcon,
	items.InProgress: styles.InProgressIcon,
	items.Cancelled:  styles.CancelledIcon,
	items.Todo:       styles.TodoIcon,
}

var taskTitleStyle = map[items.Status]lipgloss.Style{
	items.Done:       styles.DoneTitle,
	items.InProgress: styles.Default,
	items.Cancelled:  styles.DoneTitle,
	items.Todo:       styles.Default,
}

func (r CliRendererer) Render(item items.Renderable) string {
	switch v := item.(type) {
	case items.Task:
		return r.renderTask(v)
	case items.Note:
		return r.renderNote(v)
	default:
		return fmt.Sprintf("Unknown renderable item type: %T\n", v)
	}
}

func (r CliRendererer) renderTask(t items.Task) string {
	var contentIcon string

	icon := taskIcons[t.Status]
	style := taskTitleStyle[t.Status]

	if len(t.Body) > 1 {
		contentIcon = styles.ContentIcon
	}

	title := style.Render(t.Title)

	return icon + contentIcon + title + "\n"
}

func (r CliRendererer) renderNote(t items.Note) string {
	var contentIcon string

	if len(t.Body) > 1 {
		contentIcon = styles.ContentIcon
	}

	return styles.NoteIcon + contentIcon + t.Title + "\n"
}

func RenderTask(t items.Task) string {
	var title string
	var icon string
	var cIcon string
	var style lipgloss.Style

	switch t.Status {
	case items.Done:
		icon = styles.DoneIcon
		style = styles.DoneTitle
	case items.Cancelled:
		icon = styles.CancelledIcon
		style = styles.DoneTitle
	case items.InProgress:
		icon = styles.InProgressIcon
		style = styles.Default
	case items.Todo:
		icon = styles.TodoIcon
		style = styles.Default
	}

	if len(t.Body) > 1 {
		cIcon = styles.ContentIcon
	}

	title += style.
		// PaddingBottom(1).
		Render(t.Title)

	return icon + cIcon + title + "\n"
}
