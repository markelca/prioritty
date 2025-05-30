package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/pkg/tasks"
	"github.com/spf13/viper"
)

type State struct {
	cursor int
	tasks  []tasks.Task
}

type Params struct {
	withTui bool
}

type Model struct {
	params  Params
	state   State
	Service tasks.Service
}

var Help = help.New()

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
		state:   State{tasks: tasks},
		params:  Params{withTui: withTui},
		Service: service,
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
