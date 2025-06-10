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
	Use:     "remove {id...}",
	Aliases: []string{"rm", "delete"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Removes one or more tasks by ID",
	Long:    `Removes one or more tasks from the list by providing their IDs`,
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.InitialModel(false)

		for _, arg := range args {
			index, err := strconv.Atoi(arg)
			if err != nil {
				log.Printf("Error: Invalid task ID '%s'. Please provide a valid number.\n", arg)
				continue
			}

			item := m.GetItemAt(index - 1)

			if item == nil {
				log.Printf("Task at index %d does not exist", index)
				continue
			}

			err = m.Service.RemoveItem(item)
			if err != nil {
				log.Printf("Error removing task: %v\n", err)
				continue
			}
		}
	},
}
