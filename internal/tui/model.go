package tui

import (
	"log"
	"os"

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

// sortItemsByTag groups items by tag in the order tags first appear.
// Items without tags come first (under "My Board").
func sortItemsByTag(itemList []items.ItemInterface) []items.ItemInterface {
	var result []items.ItemInterface
	tagOrder := []string{}
	itemsByTag := make(map[string][]items.ItemInterface)

	for _, item := range itemList {
		var tagKey string
		if tag := item.GetTag(); tag != nil {
			tagKey = tag.Name
		}

		if _, exists := itemsByTag[tagKey]; !exists {
			tagOrder = append(tagOrder, tagKey)
		}
		itemsByTag[tagKey] = append(itemsByTag[tagKey], item)
	}

	for _, tagKey := range tagOrder {
		result = append(result, itemsByTag[tagKey]...)
	}
	return result
}

// Mode represents the current operation mode of the TUI
type Mode string

const (
	ModeList   Mode = "list"   // default: browsing items
	ModeCreate Mode = "create" // creating a new item
	ModeEdit   Mode = "edit"   // editing an existing item
)

// Params controls the behavior of the TUI model
type Params struct {
	IsTUI bool // true = return to list after action, false = quit after action
	Mode  Mode // current operation mode
}

type Model struct {
	params   Params
	state    State
	Service  service.Service
	renderer render.CLI
	initCmd  tea.Cmd // command to execute on Init(), used for CLI create/edit
}

func InitialModel(isTUI bool) Model {
	isDemo := viper.GetBool("demo")
	repoType := viper.GetString(config.CONF_REPOSITORY_TYPE)

	var repo repository.Repository
	dbPath, err := repository.GetDatabasePath(repoType, isDemo)
	if err != nil {
		log.Printf("Error - %s:", err)
		os.Exit(ExitCodeRepositoryNotSupported)
	}

	switch repoType {
	case repository.RepoTypeObsidian:
		repo, err = obsidianMigrations.NewObsidianRepository(dbPath)
	case repository.RepoTypeSQLite:
		repo, err = sqliteMigrations.NewSQLiteRepository(dbPath)
	default:
		log.Println("Error - Repository type not supported: ", repoType)
		os.Exit(ExitCodeRepositoryNotSupported)
	}

	if err != nil {
		log.Println("Error - Failed to create repository:", err)
		os.Exit(ExitCodeRepositoryCreate)
	}

	service := service.NewService(repo)

	itemList, err := service.GetAll()
	if err != nil {
		log.Println("Error - Failed to get the tasks:", err)
		os.Exit(ExitCodeGetItems)
	}
	itemList = sortItemsByTag(itemList)

	taskContent := ItemContent{}
	return Model{
		state:    State{item: taskContent, items: itemList},
		params:   Params{IsTUI: isTUI},
		Service:  service,
		renderer: render.CLI{},
	}
}

func (m Model) Init() tea.Cmd {
	// Return any command set during model creation (used for CLI create/edit)
	return m.initCmd
}

func (m Model) GetItemAt(index int) items.ItemInterface {
	if index < 0 || index >= len(m.state.items) {
		return nil
	}
	return m.state.items[index]
}

func (m Model) DestroyDemo() {
	err := m.Service.DestroyDemo()
	if err != nil {
		log.Println("Error - Failed destroy the demo data", err)
		os.Exit(ExitCodeDestroyDemo)
	}
}

// refreshItems reloads the item list from the service.
func (m *Model) refreshItems() {
	itemList, err := m.Service.GetAll()
	if err != nil {
		log.Println("Error refreshing items:", err)
		return
	}
	m.state.items = sortItemsByTag(itemList)
}

// CreateModel returns a model configured for CLI item creation
func CreateModel(itemType items.ItemType) Model {
	m := InitialModel(false)
	m.params.Mode = ModeCreate
	cmd, err := m.Service.CreateWithEditor(itemType)
	if err != nil {
		log.Println("Error opening editor:", err)
		m.initCmd = tea.Quit
	} else {
		m.initCmd = cmd
	}
	return m
}

// EditModel returns a model configured for CLI item editing
func EditModel(item items.ItemInterface) Model {
	m := InitialModel(false)
	m.state.items = []items.ItemInterface{item}
	m.state.cursor = 0
	m.params.Mode = ModeEdit
	cmd, err := m.Service.EditWithEditor(item)
	if err != nil {
		log.Println("Error opening editor:", err)
		m.initCmd = tea.Quit
	} else {
		m.initCmd = cmd
	}
	return m
}
