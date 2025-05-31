package tui

import (
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/pkg/tasks"
	"github.com/spf13/viper"
)

type TaskContentState struct {
	content  string
	ready    bool
	viewport viewport.Model
}

type State struct {
	cursor      int
	tasks       []tasks.Task
	taskContent TaskContentState
}

func (s State) GetCurrentTask() *tasks.Task {
	if s.cursor+1 > len(s.tasks) {
		return &tasks.Task{}
	}
	return &s.tasks[s.cursor]
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
	isDemo := viper.GetBool("demo")
	var dbFilePath string

	if isDemo {
		dbFilePath = path.Join(os.TempDir(), "prioritty_demo.db")
	} else {
		dbFilePath = viper.GetString(config.CONF_DATABASE_PATH)
	}

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

	taskContent := TaskContentState{}
	return Model{
		state:   State{tasks: tasks, taskContent: taskContent},
		params:  Params{withTui: withTui},
		Service: service,
	}
}

func (m Model) DestroyDemo() {
	err := m.Service.DestroyDemo()
	if err != nil {
		fmt.Println("Error - Failed destroy the demo data", err)
		os.Exit(5)
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
