package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/cli"
	tui "github.com/markelca/prioritty/internal/tui"
)

func main() {
	if len(os.Args) < 2 {
		p := tea.NewProgram(tui.InitialModel(true))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	} else {
		cli.Execute()
	}
}
