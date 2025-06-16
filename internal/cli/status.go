package cli

import (
	"fmt"
	"strconv"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(todoCmd)
	rootCmd.AddCommand(cancelCmd)
	rootCmd.AddCommand(startCmd)
}

func updateTaskStatus(args []string, status items.Status) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide at least one task ID")
	}

	m := tui.InitialModel(false)

	for _, arg := range args {
		i, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("Invalid task ID: %s\n", arg)
			continue
		}

		// Find item by ID
		item := m.GetItemAt(i - 1)
		switch v := item.(type) {
		case *items.Task:
			if item == nil {
				fmt.Printf("Task with ID %d not found\n", i)
				continue
			}

			err = m.Service.UpdateStatus(v, status)
			if err != nil {
				fmt.Printf("Failed to update task %d: %v\n", i, err)
				continue
			}
		default:
			fmt.Printf("Failed to update status, item %d has to be a task", i)
			continue
		}

	}

	return nil
}

var doneCmd = &cobra.Command{
	Use:   "done [task_ids...]",
	Short: "Mark tasks as done",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateTaskStatus(args, items.Done)
	},
}

var todoCmd = &cobra.Command{
	Use:   "todo [task_ids...]",
	Short: "Mark tasks as todo",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateTaskStatus(args, items.Todo)
	},
}

var cancelCmd = &cobra.Command{
	Use:   "cancel [task_ids...]",
	Short: "Mark tasks as cancelled",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateTaskStatus(args, items.Cancelled)
	},
}

var startCmd = &cobra.Command{
	Use:     "start [task_ids...]",
	Aliases: []string{"progress", "pg"},
	Short:   "Mark tasks as in progress",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateTaskStatus(args, items.InProgress)
	},
}
