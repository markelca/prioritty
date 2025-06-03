package renderer

import "github.com/markelca/prioritty/pkg/items"

type Renderer interface {
	RenderItem(item items.Item) string
}
