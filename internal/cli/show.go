package cli

import (
	"fmt"
	"strconv"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show [index]",
	Short: "Show task or note details by index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: Invalid index '%s'. Please provide a valid number.\n", args[0])
			return
		}

		m := tui.InitialModel(false)
		item := m.GetItemAt(index - 1) // Convert to 0-based index

		if item == nil {
			allItems, err := m.Service.GetAll()
			if err != nil {
				fmt.Printf("Error: Could not retrieve items: %v\n", err)
				return
			}
			fmt.Printf("Error: Index %d is out of range. Available items: 1-%d\n", index, len(allItems))
			return
		}

		fmt.Printf("Title: %s\n", item.GetTitle())
		if item.GetBody() != "" {
			fmt.Printf("Content:\n%s\n", item.GetBody())
		} else {
			fmt.Printf("Content: (empty)\n")
		}
	},
}