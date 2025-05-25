package tui

import (
	"fmt"
	"image/color"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/markelca/prioritty/internal/tui/styles"
	"github.com/markelca/prioritty/pkg/tasks"
	"github.com/muesli/gamut"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	columnWidth = 30
)

var (
	subtle         = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	blends = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 50)
)

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}

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
				styles.Done.Render("✓"),
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
		donePercentage = float64(countDone) / float64(len(m.tasks)) * 100
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

	{
		okButton := activeButtonStyle.Render("Yes")
		cancelButton := buttonStyle.Render("Maybe")

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(rainbow(lipgloss.NewStyle(), "Are you sure you want to eat marmalade?", blends))
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		dialog := lipgloss.Place(width, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("猫"),
			lipgloss.WithWhitespaceForeground(subtle),
		)

		s += dialog + "\n\n"

		// doc.WriteString(dialog + "\n\n")
	}

	var helpView string
	if m.withTui {
		helpView = m.help.View(m.keys)

	}

	return s + "\n" + helpView
}
