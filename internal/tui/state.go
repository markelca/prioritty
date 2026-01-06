package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/markelca/prioritty/pkg/items"
)

type State struct {
	cursor        int
	items         []items.ItemInterface
	contentView   ItemContent            // viewport for displaying item details
	Mode          Mode                   // current operation mode
	pendingDelete items.ItemInterface    // item awaiting deletion confirmation
}

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

func (s State) GetCurrentItem() items.ItemInterface {
	if s.cursor < 0 || s.cursor >= len(s.items) {
		return nil
	}
	return s.items[s.cursor]
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
