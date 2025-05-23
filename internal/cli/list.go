package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows all the tasks",
	Long:  `[Long description]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("onueth")

		p := tea.NewProgram(tui.InitialModel(), tea.WithoutRenderer())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}
