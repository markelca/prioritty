package tui

import (
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/renderer"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/markelca/prioritty/pkg/items/service"
	"github.com/spf13/viper"
)

type TaskContentState struct {
	content  string
	ready    bool
	viewport viewport.Model
}

type State struct {
	cursor      int
	items       []items.ItemInterface
	taskContent TaskContentState
}

func (s State) GetCurrentItem() items.ItemInterface {
	if s.cursor+1 > len(s.items) {
		return nil
	}
	return s.items[s.cursor]
}

type Params struct {
	withTui bool
}

type Model struct {
	params   Params
	state    State
	Service  service.Service
	renderer renderer.CliRendererer
}

func (m Model) GetItemAt(index int) items.ItemInterface {
	if index <= 0 || len(m.state.items)-1 < index {
		return nil
	} else {
		return m.state.items[index]
	}
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

	repo, err := repository.NewSQLiteRepository(dbFilePath)
	if err != nil {
		fmt.Println("Error - Failed to create repository:", err)
		os.Exit(3)
	}

	service := service.NewService(repo)

	itemList, err := service.GetAll()
	if err != nil {
		fmt.Println("Error - Failed to get the tasks:", err)
		os.Exit(4)
	}

	taskContent := TaskContentState{}
	return Model{
		state:    State{taskContent: taskContent, items: itemList},
		params:   Params{withTui: withTui},
		Service:  service,
		renderer: renderer.CliRendererer{},
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
