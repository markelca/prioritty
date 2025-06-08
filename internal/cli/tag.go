package cli

import (
	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tagCmd)
}

var tagCmd = &cobra.Command{
	Use:     "tag {id}",
	Aliases: []string{},
	Args:    cobra.ExactArgs(1),
	Short:   "Sets the tag for a task",
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.InitialModel(false)
		m.Service.AddTask(args[0])
	},
}
