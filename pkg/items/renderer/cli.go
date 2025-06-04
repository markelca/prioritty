package renderer

import (
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

func (r CliRendererer) RenderTask(t items.Task) string {
	var contentIcon string

	icon := taskIcons[t.Status]
	style := taskTitleStyle[t.Status]

	if len(t.Body) > 1 {
		contentIcon = styles.ContentIcon
	}

	title := style.Render(t.Title)

	return icon + contentIcon + title + "\n"
}

func (r CliRendererer) RenderNote(t items.Note) string {
	return ""
}
