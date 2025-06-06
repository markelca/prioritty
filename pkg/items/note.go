package items

var _ ItemInterface = (*Task)(nil)

type Note struct {
	Item
}

func (n Note) Render(r Renderer) string {
	return r.Render(n)
}
