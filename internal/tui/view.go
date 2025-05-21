package tui

func (m model) View() string {
	return appStyle.Render(m.list.View())
}
