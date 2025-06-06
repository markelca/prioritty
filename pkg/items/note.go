package items

var _ Renderable = (*Note)(nil)
var _ Base = (*Note)(nil)

type Note struct {
	Item
}

func (n Note) GetBody() string {
	return n.Item.Body
}

func (n Note) GetTitle() string {
	return n.Item.Title
}

func (n Note) Render(r Renderer) string {
	return r.Render(n)
}
