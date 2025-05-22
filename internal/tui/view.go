package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

var successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#04B575"))

var greyStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#595a70"))

func (m model) View() string {
	s := ""

	for i, choice := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		if _, ok := m.done[i]; ok {
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				successStyle.Render("✔"),
				greyStyle.Render(choice))
		} else {
			s += fmt.Sprintf("%s %s %s %s\n",
				defaultStyle.Render(cursor),
				greyStyle.Render(fmt.Sprintf("%d.", i+1)),
				defaultStyle.Render("☐"),
				choice)
		}
	}

	s += greyStyle.Render("\nPress q to quit.\n")
	return s
}
