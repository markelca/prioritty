package cli

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/tui"
	"github.com/markelca/prioritty/pkg/tasks"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tuiCmd)
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	Long:  `Launch the interactive Terminal User Interface for managing tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := tasks.NewSQLiteRepository("data/test.db")
		if err != nil {
			log.Fatal("Failed to create repository:", err)
		}

		tasks := repo.FindAll()

		m := tui.InitialModel(true)
		m.Tasks = tasks
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}
