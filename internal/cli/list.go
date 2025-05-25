package cli

import (
	"fmt"

	"github.com/markelca/prioritty/internal/tui"
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
		m := tui.InitialModel(false)
		fmt.Print(m.View())
	},
}
