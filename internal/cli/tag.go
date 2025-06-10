package cli

import (
	"log"
	"strconv"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagUnsetCmd)
}

var tagCmd = &cobra.Command{
	Use:     "tag {tag} {id...}",
	Aliases: []string{},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Sets the tag for one or more tasks",
	Long:    `Sets the tag for one or more tasks by providing the tag name and their IDs`,
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.InitialModel(false)
		tag := args[0]

		for _, arg := range args[1:] {
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

			err = m.Service.SetTag(item, tag)
			if err != nil {
				log.Printf("Error setting tag for task %d: %v\n", index, err)
				continue
			}
		}
	},
}

var tagUnsetCmd = &cobra.Command{
	Use:   "unset {id...}",
	Args:  cobra.MinimumNArgs(1),
	Short: "Unsets the tag for one or more tasks or notes",
	Long:  `Unsets the tag for one or more tasks or notes by providing their IDs`,
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

			err = m.Service.UnsetTag(item)
			if err != nil {
				log.Printf("Error unsetting tag for task %d: %v\n", index, err)
				continue
			}
		}
	},
}
