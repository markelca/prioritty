package repository

import (
	"database/sql"
	"os"

	"github.com/markelca/prioritty/pkg/items"
)

type TaskRepository interface {
	GetTasks() ([]items.Task, error)
	UpdateTask(items.Task) error
	CreateTask(items.Task) error
	RemoveTask(int) error
	UpdateTaskStatus(items.Task, items.Status) error
}

type NoteRepository interface {
	GetNotes() ([]items.Note, error)
	UpdateNote(items.Note) error
	CreateNote(items.Note) error
	RemoveNote(int) error
}

type Repository interface {
	TaskRepository
	NoteRepository
	DropSchema() error
}

type SQLiteRepository struct {
	db       *sql.DB
	filepath string
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteRepository{db: db, filepath: dbPath}

	return repo, nil
}

func (r *SQLiteRepository) DropSchema() error {
	return os.Remove(r.filepath)
}
