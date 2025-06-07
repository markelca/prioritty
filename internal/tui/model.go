package tui

import (
	"log"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/renderer"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/markelca/prioritty/pkg/items/service"
	"github.com/spf13/viper"
)

type ItemContent struct {
	content  string
	ready    bool
	viewport viewport.Model
}

type ItemContentDimensions struct {
	width        int
	height       int
	headerHeight int
	footerHeight int
}

func (itemContent *ItemContent) init(dimensions ItemContentDimensions) {
	var (
		width        = dimensions.width
		height       = dimensions.height
		headerHeight = dimensions.headerHeight
		footerHeight = dimensions.footerHeight
	)
	verticalMarginHeight := headerHeight + footerHeight
	if !itemContent.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		itemContent.viewport = viewport.New(width, height-verticalMarginHeight)
		itemContent.viewport.YPosition = headerHeight
	} else {
		itemContent.viewport.Width = width
		itemContent.viewport.Height = height - verticalMarginHeight
	}

}

func (content *ItemContent) show(item items.ItemInterface) {
	style := lipgloss.NewStyle().Width(content.viewport.Width)
	if content.ready {
		content.ready = false
	} else {
		body := item.GetBody()
		contentStr := style.Render(body)
		content.viewport.SetContent(contentStr)
		content.ready = true
	}

}

type State struct {
	cursor int
	items  []items.ItemInterface
	item   ItemContent
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
		renderer: renderer.CliRendererer{},
	}
}

func (m Model) DestroyDemo() {
	err := m.Service.DestroyDemo()
	if err != nil {
		log.Println("Error - Failed destroy the demo data", err)
		os.Exit(5)
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
