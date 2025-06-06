package items

type Base interface {
	GetId() int
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

func (i Item) GetId() int {
	return i.Id
}

func (i Item) GetBody() string {
	return i.Body
}

func (i Item) GetTitle() string {
	return i.Title
}
