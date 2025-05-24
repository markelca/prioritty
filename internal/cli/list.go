package cli

import (
	"fmt"
	"log"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/markelca/prioritty/pkg/tasks"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Shows all the tasks",
	Long:    `[Long description]`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := tasks.NewSQLiteRepository("data/test.db")
		if err != nil {
			log.Fatal("Failed to create repository:", err)
		}

		tasks := repo.FindAll()

		m := tui.InitialModel(false)
		m.Tasks = tasks
		fmt.Print(m.View())
	},
}
