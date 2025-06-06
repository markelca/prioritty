package items

type Status int

const (
	Todo Status = iota
	InProgress
	Done
	Cancelled
	NoteType
)

var _ Renderable = (*Task)(nil)
var _ Base = (*Task)(nil)

type Task struct {
	Item
	Status Status
}

func (t Task) GetTitle() string {
	return t.Item.Title
}

func (t Task) GetBody() string {
	return t.Item.Body
}

func (t *Task) SetStatus(s Status) {
	if t.Status == s {
		t.Status = Todo
	} else {
		t.Status = s
	}
}

func (t Task) Render(r Renderer) string {
	return r.Render(t)
}
