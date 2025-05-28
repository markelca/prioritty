package tasks

type Status int

const (
	Todo Status = iota
	InProgress
	Done
	Cancelled
)

type Task struct {
	Id     int
	Title  string
	Body   string
	Status Status
}

func (t *Task) SetStatus(s Status) {
	if t.Status == s {
		t.Status = Todo
	} else {
		t.Status = s
	}
}
