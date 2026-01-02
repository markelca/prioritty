package items

import "time"

type Base interface {
	GetId() string
	GetTitle() string
	GetBody() string
	GetTag() *Tag
	GetCreatedAt() time.Time
	After(ItemInterface) bool
}

type ItemInterface interface {
	Base
	Renderable
}

type Tag struct {
	Id   string
	Name string
}

type Item struct {
	Id        string
	Title     string
	Body      string
	CreatedAt time.Time
	Tag       *Tag
}

func (i Item) GetId() string {
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
func (i Item) GetTag() *Tag {
	return i.Tag
}

func (i Item) After(u ItemInterface) bool {
	// 1. Items with a tag come before items without a tag
	iHasTag := i.GetTag() != nil
	uHasTag := u.GetTag() != nil

	if iHasTag && !uHasTag {
		return false // i should come before u
	}
	if !iHasTag && uHasTag {
		return true // i should come after u
	}

	// 2. If both have (or don't have) tags, sort by CreatedAt (most recent first)
	return i.GetCreatedAt().Before(u.GetCreatedAt())
}
