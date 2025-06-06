package cli

import (
	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(taskCmd)
}

var taskCmd = &cobra.Command{
	Use:     "task {title}",
	Aliases: []string{},
	Args:    cobra.ExactArgs(1),
	Short:   "Adds a new task",
	Long:    `[Long description]`,
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.InitialModel(false)
		m.Service.AddTask(args[0])
	},
}
