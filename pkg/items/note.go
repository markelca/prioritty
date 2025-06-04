package items

type Note struct {
	Item
}

func (n Note) Render(r Renderer) string {
	return r.RenderNote(n)
}
