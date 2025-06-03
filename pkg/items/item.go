package items

type ItemType int

const (
	TaskItem ItemType = iota
	NoteItem
)

type Item interface {
	GetId() int
	GetTitle() string
	GetBody() string
	GetType() ItemType
}
