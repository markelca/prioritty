package cli

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit <index>",
	Short: "Edit a task or note by index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: invalid index '%s'. Please provide a valid number.\n", args[0])
			return
		}

		// First get the item to edit
		m := tui.InitialModel(false)
		item := m.GetItemAt(index - 1) // Convert to 0-based index
		if item == nil {
			fmt.Printf("Error: no item found at index %d\n", index)
			return
		}

		// Create edit model and run
		editModel := tui.EditModel(item)
		tea.NewProgram(editModel).Run()
	},
}