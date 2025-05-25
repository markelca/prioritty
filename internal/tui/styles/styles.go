package styles

import "github.com/charmbracelet/lipgloss"

var (
	Default = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

	Done = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1"))

	Secondary = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#595a70"))

	InProgress = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa0df"))

	Cancelled = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f28aa8"))

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}
)
