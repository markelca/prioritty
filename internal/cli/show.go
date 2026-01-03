package cli

import (
	"fmt"
	"strconv"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/markelca/prioritty/internal/tui/styles"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/markdown"
	"github.com/spf13/cobra"
)

var rawOutput bool

func init() {
	showCmd.Flags().BoolVar(&rawOutput, "raw", false, "Show item with frontmatter (markdown format)")
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

		if rawOutput {
			var input markdown.ItemInput
			input.Title = item.GetTitle()
			input.Body = item.GetBody()
			if tag := item.GetTag(); tag != nil {
				input.Tag = tag.Name
			}

			if task, ok := item.(*items.Task); ok {
				input.ItemType = items.ItemTypeTask
				input.Status = string(task.Status)
			} else {
				input.ItemType = items.ItemTypeNote
			}

			content, err := markdown.Serialize(input)
			if err != nil {
				fmt.Printf("Error: Could not serialize item: %v\n", err)
				return
			}

			fmt.Print(content)
			return
		}

		// Default output: icon + title + tag + body
		icon := tui.GetItemIcon(item)
		title := icon + item.GetTitle()
		if tag := item.GetTag(); tag != nil {
			title += " " + styles.Secondary.Render("@"+tag.Name)
		}
		fmt.Println(title)
		if item.GetBody() != "" {
			fmt.Printf("\n" + item.GetBody())
		}
	},
}

