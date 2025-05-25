package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/pkg/tasks"
	"github.com/spf13/viper"
)

type Model struct {
	withTui    bool
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	Service    tasks.Service
	tasks      []tasks.Task // items on the to-do list
	cursor     int          // which to-do list item our cursor is pointing at
}

func InitialModel(withTui bool) Model {
	dbFilePath := viper.GetString(config.KEY_DATABASE_PATH)

	repo, err := tasks.NewSQLiteRepository(dbFilePath)
	if err != nil {
		fmt.Println("Error - Failed to create repository:", err)
		os.Exit(3)
	}

	service := tasks.NewService(repo)

	tasks, err := service.FindAll()
	if err != nil {
		fmt.Println("Error - Failed to get the tasks:", err)
		os.Exit(4)
	}

	return Model{
		withTui:    withTui,
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		Service:    service,
		tasks:      tasks,
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
