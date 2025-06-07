package items

import "time"

type Base interface {
	GetId() int
	GetTitle() string
	GetBody() string
	GetCreatedAt() time.Time
}

type ItemInterface interface {
	Base
	Renderable
}

type Item struct {
	Id        int
	Title     string
	Body      string
	CreatedAt time.Time
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

func (i Item) GetCreatedAt() time.Time {
	return i.CreatedAt
}
