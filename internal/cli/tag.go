package cli

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/markelca/prioritty/internal/tui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tagCmd)
	rootCmd.AddCommand(tagsCmd)
	tagCmd.AddCommand(tagUnsetCmd)
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagRmCmd)
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

func listTags(cmd *cobra.Command, args []string) {
	m := tui.InitialModel(false)

	tags, err := m.Service.GetTags()
	if err != nil {
		log.Printf("Error getting tags: %v\n", err)
		return
	}

	if len(tags) == 0 {
		fmt.Println("No tags found")
		return
	}

	for _, tag := range tags {
		fmt.Println(tag.Name)
	}
}

var tagListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Short:   "Lists all available tags",
	Long:    `Lists all available tags in the system`,
	Run:     listTags,
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Args:  cobra.NoArgs,
	Short: "Lists all available tags",
	Long:  `Lists all available tags in the system`,
	Run:   listTags,
}

var tagRmCmd = &cobra.Command{
	Use:   "rm {tag}",
	Args:  cobra.ExactArgs(1),
	Short: "Removes a tag from the database",
	Long:  `Removes a tag from the database if it's not assigned to any tasks or notes`,
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.InitialModel(false)
		tagName := args[0]

		err := m.Service.RemoveTag(tagName)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("Tag '%s' not found\n", tagName)
				return
			}
			fmt.Printf("Error removing tag: %v\n", err)
			return
		}

		fmt.Printf("Tag '%s' removed successfully\n", tagName)
	},
}
