package cli

import (
	"fmt"
	"strconv"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove {id}",
	Aliases: []string{"rm", "delete"},
	Args:    cobra.ExactArgs(1),
	Short:   "Removes a task by ID",
	Long:    `Removes a task from the list by providing its ID`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: Invalid task ID '%s'. Please provide a valid number.\n", args[0])
			return
		}

		m := tui.InitialModel(false)
		err = m.Service.RemoveTask(id)
		if err != nil {
			fmt.Printf("Error removing task: %v\n", err)
			return
		}
	},
}

