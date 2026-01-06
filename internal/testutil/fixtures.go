package testutil

import (
	"time"

	"github.com/markelca/prioritty/pkg/items"
)

// TaskOption is a functional option for configuring test tasks
type TaskOption func(*items.Task)

// WithStatus sets the status of a test task
func WithStatus(s items.Status) TaskOption {
	return func(t *items.Task) {
		t.Status = s
	}
}

// WithBody sets the body of a test task
func WithBody(body string) TaskOption {
	return func(t *items.Task) {
		t.Body = body
	}
}

// WithTag sets the tag of a test task
func WithTag(tag *items.Tag) TaskOption {
	return func(t *items.Task) {
		t.Tag = tag
	}
}

// WithCreatedAt sets the creation time of a test task
func WithCreatedAt(t time.Time) TaskOption {
	return func(task *items.Task) {
		task.CreatedAt = t
	}
}

// WithId sets the ID of a test task
func WithId(id string) TaskOption {
	return func(t *items.Task) {
		t.Id = id
	}
}

// NewTestTask creates a Task with sensible defaults for testing
func NewTestTask(title string, opts ...TaskOption) *items.Task {
	task := &items.Task{
		Item: items.Item{
			Id:        "test-task-id",
			Title:     title,
			Body:      "",
			CreatedAt: time.Now(),
			Tag:       nil,
		},
		Status: items.Todo,
	}

	for _, opt := range opts {
		opt(task)
	}

	return task
}

// NoteOption is a functional option for configuring test notes
type NoteOption func(*items.Note)

// WithNoteBody sets the body of a test note
func WithNoteBody(body string) NoteOption {
	return func(n *items.Note) {
		n.Body = body
	}
}

// WithNoteTag sets the tag of a test note
func WithNoteTag(tag *items.Tag) NoteOption {
	return func(n *items.Note) {
		n.Tag = tag
	}
}

// WithNoteCreatedAt sets the creation time of a test note
func WithNoteCreatedAt(t time.Time) NoteOption {
	return func(n *items.Note) {
		n.CreatedAt = t
	}
}

// WithNoteId sets the ID of a test note
func WithNoteId(id string) NoteOption {
	return func(n *items.Note) {
		n.Id = id
	}
}

// NewTestNote creates a Note with sensible defaults for testing
func NewTestNote(title string, opts ...NoteOption) *items.Note {
	note := &items.Note{
		Item: items.Item{
			Id:        "test-note-id",
			Title:     title,
			Body:      "",
			CreatedAt: time.Now(),
			Tag:       nil,
		},
	}

	for _, opt := range opts {
		opt(note)
	}

	return note
}

// NewTestTag creates a Tag with sensible defaults for testing
func NewTestTag(name string) *items.Tag {
	return &items.Tag{
		Id:   "test-tag-id",
		Name: name,
	}
}

// NewTestTagWithId creates a Tag with a specific ID
func NewTestTagWithId(id, name string) *items.Tag {
	return &items.Tag{
		Id:   id,
		Name: name,
	}
}
