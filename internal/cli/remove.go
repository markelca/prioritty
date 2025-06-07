package cli

import (
	"log"
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
		index, err := strconv.Atoi(args[0])
		if err != nil {
			log.Printf("Error: Invalid task ID '%s'. Please provide a valid number.\n", args[0])
			return
		}

		m := tui.InitialModel(false)

		item := m.GetItemAt(index - 1)

		if item == nil {
			log.Printf("Task at index %d does not exist", index)
			return
		}

		err = m.Service.RemoveItem(item)
		if err != nil {
			log.Printf("Error removing task: %v\n", err)
			return
		}
	},
}
