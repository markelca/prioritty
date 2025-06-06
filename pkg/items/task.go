package items

type Status int

const (
	Todo Status = iota
	InProgress
	Done
	Cancelled
	NoteType
)

var _ ItemInterface = (*Task)(nil)

type Task struct {
	Item
	Status Status
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
