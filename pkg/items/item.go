package items

type Base interface {
	GetTitle() string
	GetBody() string
}

type ItemInterface interface {
	Base
	Renderable
}

type Item struct {
	Id    int
	Title string
	Body  string
}
