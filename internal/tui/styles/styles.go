package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/pkg/tasks"
)

var (
	special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	Default = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

	Done = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1"))

	Secondary = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#595a70"))

	InProgress = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa0df"))

	Cancelled = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f28aa8"))

	DoneIcon = Done.SetString("✓").
			PaddingRight(1).
			String()

	TodoIcon = Default.SetString("☐").
			PaddingRight(1).
			String()

	InProgressIcon = InProgress.SetString("◐").
			PaddingRight(1).
			String()

	CancelledIcon = Cancelled.SetString("x").
			PaddingRight(1).
			String()

	DoneTitle = Secondary.
			Strikethrough(true)

	ContentIcon = Default.SetString("📄").
			PaddingRight(1).
			String()

	RenderTask = func(t tasks.Task) string {
		var title string
		var icon string
		var cIcon string
		var style lipgloss.Style

		switch t.Status {
		case tasks.Done:
			icon = DoneIcon
			style = DoneTitle
		case tasks.Cancelled:
			icon = CancelledIcon
			style = DoneTitle
		case tasks.InProgress:
			icon = InProgressIcon
			style = Default
		case tasks.Todo:
			icon = TodoIcon
			style = Default
		}

		if len(t.Body) > 1 {
			cIcon = ContentIcon
		}

		title += style.
			// PaddingBottom(1).
			Render(t.Title)

		return icon + cIcon + title + "\n"
	}

	TitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	InfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return TitleStyle.BorderStyle(b)
	}()
)
