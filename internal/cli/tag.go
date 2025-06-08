package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tagCmd)
}

var tagCmd = &cobra.Command{
	Use:     "tag {id} {tag}",
	Aliases: []string{},
	Args:    cobra.ExactArgs(2),
	Short:   "Sets the tag for a task",
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.InitialModel(false)
		i, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error converting string index to int")
		}

		item := m.GetItemAt(i - 1)
		if item == nil {
			fmt.Printf("Error - No item found with index %d", i)
		}
		err = m.Service.SetTag(item, args[1])
		if err != nil {
			log.Printf("Error setting the task: %v", err)
		}
	},
}
