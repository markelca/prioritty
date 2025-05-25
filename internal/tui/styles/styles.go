package styles

import "github.com/charmbracelet/lipgloss"

var Default = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

var Done = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#a6e3a1"))

var Secondary = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#595a70"))

var InProgress = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7aa0df"))

var Cancelled = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#f28aa8"))
