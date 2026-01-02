package items

import "strings"

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

// StatusToString converts a Status enum to its string value.
func StatusToString(s Status) string {
	switch s {
	case Todo:
		return "todo"
	case InProgress:
		return "in-progress"
	case Done:
		return "done"
	case Cancelled:
		return "cancelled"
	default:
		return "todo"
	}
}

// StringToStatus converts a string to Status enum.
func StringToStatus(s string) Status {
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
