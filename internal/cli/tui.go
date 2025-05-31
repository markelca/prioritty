package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(tuiCmd)
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	Long:  `Launch the interactive Terminal User Interface for managing tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		model := tui.InitialModel(true)
		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
		)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(2)
		}
		isDemo := viper.GetBool("demo")
		fmt.Println(isDemo)
		if isDemo {
			model.DestroyDemo()
		}
	},
}
