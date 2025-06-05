package items

var _ Renderable = (*Note)(nil)

type Note struct {
	Item
}

func (n Note) Render(r Renderer) string {
	return r.Render(n)
}
