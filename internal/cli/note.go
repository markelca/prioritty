package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/tui"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(noteCmd)
}

var noteCmd = &cobra.Command{
	Use:     "note [title]",
	Aliases: []string{},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Adds a new note",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Create with editor
			m := tui.CreateModel(items.ItemTypeNote)
			tea.NewProgram(m).Run()
		} else {
			// Create with title only
			m := tui.InitialModel(false)
			m.Service.AddNote(args[0])
		}
	},
}
