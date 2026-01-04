package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/elliotchance/orderedmap/v3"
	"github.com/markelca/prioritty/pkg/items"
)

type State struct {
	cursor int
	items  []items.ItemInterface
	item   ItemContent
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

// getDisplayOrderedItems returns items in the same order they appear in the view
// (grouped by tag, with null-tag items under "My Board").
func (s State) getDisplayOrderedItems() []items.ItemInterface {
	itemsByTag := orderedmap.NewOrderedMap[items.Tag, []items.ItemInterface]()

	for _, item := range s.items {
		tag := item.GetTag()
		if tag != nil {
			tagItems, _ := itemsByTag.Get(*tag)
			itemsByTag.Set(*tag, append(tagItems, item))
		} else {
			nullTag := items.Tag{}
			tagItems, _ := itemsByTag.Get(nullTag)
			itemsByTag.Set(nullTag, append(tagItems, item))
		}
	}

	var result []items.ItemInterface
	for _, itemList := range itemsByTag.AllFromFront() {
		result = append(result, itemList...)
	}
	return result
}

func (s State) GetCurrentItem() items.ItemInterface {
	displayItems := s.getDisplayOrderedItems()
	if s.cursor < 0 || s.cursor >= len(displayItems) {
		return nil
	}
	return displayItems[s.cursor]
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
