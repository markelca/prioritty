package tui

import (
	"log"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/config"
	obsidianMigrations "github.com/markelca/prioritty/internal/migrations/obsidian"
	sqliteMigrations "github.com/markelca/prioritty/internal/migrations/sqlite"
	"github.com/markelca/prioritty/internal/render"
	"github.com/markelca/prioritty/internal/service"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/spf13/viper"
)

var Help = help.New()

type Params struct {
	withTui    bool
	CreateMode string // "task" or "note" for creation mode
	EditMode   bool   // true when editing an existing item
}

type Model struct {
	params   Params
	state    State
	Service  service.Service
	renderer render.CLI
}

func InitialModel(withTui bool) Model {
	isDemo := viper.GetBool("demo")
	repoType := viper.GetString(config.CONF_REPOSITORY_TYPE)

	var repo repository.Repository
	var err error

	switch repoType {
	case "obsidian":
		var vaultPath string
		if isDemo {
			vaultPath = path.Join(os.TempDir(), "prioritty_demo_vault")
		} else {
			vaultPath = viper.GetString(config.CONF_DATABASE_PATH)
		}
		repo, err = obsidianMigrations.NewObsidianRepository(vaultPath)
	default: // "sqlite" or any other value defaults to sqlite
		var dbFilePath string
		if isDemo {
			dbFilePath = path.Join(os.TempDir(), "prioritty_demo.db")
		} else {
			dbFilePath = viper.GetString(config.CONF_DATABASE_PATH)
		}
		repo, err = sqliteMigrations.NewSQLiteRepository(dbFilePath)
	}

	if err != nil {
		log.Println("Error - Failed to create repository:", err)
		os.Exit(3)
	}

	service := service.NewService(repo)

	itemList, err := service.GetAll()
	if err != nil {
		log.Println("Error - Failed to get the tasks:", err)
		os.Exit(4)
	}

	taskContent := ItemContent{}
	return Model{
		state:    State{item: taskContent, items: itemList},
		params:   Params{withTui: withTui},
		Service:  service,
		renderer: render.CLI{},
	}
}

func (m Model) Init() tea.Cmd {
	// If in creation mode, immediately open the editor
	if m.params.CreateMode != "" {
		cmd, err := m.Service.CreateWithEditor(m.params.CreateMode)
		if err != nil {
			log.Println("Error opening editor:", err)
			return tea.Quit
		}
		return cmd
	}
	// If in edit mode, immediately open the editor with current item
	if m.params.EditMode {
		item := m.state.GetCurrentItem()
		cmd, err := m.Service.EditWithEditor(item)
		if err != nil {
			log.Println("Error opening editor:", err)
			return tea.Quit
		}
		return cmd
	}
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) GetItemAt(index int) items.ItemInterface {
	if index < 0 || len(m.state.items)-1 < index {
		return nil
	} else {
		return m.state.items[index]
	}
}

func (m Model) DestroyDemo() {
	err := m.Service.DestroyDemo()
	if err != nil {
		log.Println("Error - Failed destroy the demo data", err)
		os.Exit(5)
	}
}

func (m *Model) SetCreateMode(mode string) {
	m.params.CreateMode = mode
}

func EditModel(item items.ItemInterface) Model {
	// Create a minimal model for editing
	m := InitialModel(false)
	m.state.items = []items.ItemInterface{item}
	m.state.cursor = 0
	m.params.EditMode = true
	return m
}
