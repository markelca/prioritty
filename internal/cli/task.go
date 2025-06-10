package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(taskCmd)
}

var taskCmd = &cobra.Command{
	Use:     "task [title]",
	Aliases: []string{},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Adds a new task",
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.InitialModel(false)
		if len(args) == 0 {
			// Create with editor
			m.SetCreateMode("task")
			tea.NewProgram(m).Run()
		} else {
			// Create with title only
			m.Service.AddTask(args[0])
		}
	},
}
