package repository

import (
	"errors"

	"github.com/markelca/prioritty/pkg/items"
)

var ErrNotFound = errors.New("not found")

type TaskRepository interface {
	GetTasks() ([]items.Task, error)
	UpdateTask(items.Task) error
	CreateTask(items.Task) error
	RemoveTask(string) error
	UpdateTaskStatus(items.Task, items.Status) error
	SetTaskTag(items.Task, items.Tag) error
	UnsetTaskTag(items.Task) error
}

type NoteRepository interface {
	GetNotes() ([]items.Note, error)
	UpdateNote(items.Note) error
	CreateNote(items.Note) error
	RemoveNote(string) error
	SetNoteTag(items.Note, items.Tag) error
	UnsetNoteTag(items.Note) error
}

type Repository interface {
	TaskRepository
	NoteRepository
	GetTag(string) (*items.Tag, error)
	GetTags() ([]items.Tag, error)
	CreateTag(string) (*items.Tag, error)
	RemoveTag(string) error
	GetItemsWithTag(string) ([]items.ItemInterface, error)
	Reset() error
}
