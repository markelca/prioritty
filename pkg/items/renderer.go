package items

type Renderer interface {
	RenderTask(Task) string
	RenderNote(Note) string
}

type Renderable interface {
	Render(Renderer) string
}
