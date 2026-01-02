package items

type Renderer interface {
	Render(Renderable) string
}

type Renderable interface {
	Render(Renderer) string
}
