package items

import "strings"

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in-progress"
	Done       Status = "done"
	Cancelled  Status = "cancelled"
	NoteType   Status = "note" // Used for counting notes in UI
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

// ParseStatus converts a string to Status, handling alternate spellings.
func ParseStatus(s string) Status {
	switch strings.ToLower(s) {
	case "todo":
		return Todo
	case "in-progress", "inprogress":
		return InProgress
	case "done":
		return Done
	case "cancelled", "canceled":
		return Cancelled
	default:
		return Todo
	}
}
