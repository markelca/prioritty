package items

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/tui/styles"
)

type CliRendererer struct {
}

var taskIcons = map[Status]string{
	Done:       styles.DoneIcon,
	InProgress: styles.InProgressIcon,
	Cancelled:  styles.CancelledIcon,
	Todo:       styles.TodoIcon,
}

var taskTitleStyle = map[Status]lipgloss.Style{
	Done:       styles.DoneTitle,
	InProgress: styles.Default,
	Cancelled:  styles.DoneTitle,
	Todo:       styles.Default,
}

func (r CliRendererer) RenderTask(t Task) string {
	var contentIcon string

	icon := taskIcons[t.Status]
	style := taskTitleStyle[t.Status]

	if len(t.Body) > 1 {
		contentIcon = styles.ContentIcon
	}

	title := style.Render(t.Title)

	return icon + contentIcon + title + "\n"
}

func (r CliRendererer) RenderNote(t Note) string {
	return ""
}
