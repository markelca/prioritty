package repository

import (
	"database/sql"
	"log"
	"os"

	"github.com/markelca/prioritty/migrations"
	"github.com/markelca/prioritty/pkg/items"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
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
	dbExists := false
	if _, err := os.Stat(dbPath); err == nil {
		dbExists = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteRepository{db: db, filepath: dbPath}

	if !dbExists {
		if _, err := db.Exec(migrations.SchemaSQL); err != nil {
			log.Printf("Error executing schema: %v", err)
			db.Close()
			return nil, err
		}
		if viper.GetBool("demo") {
			if _, err := db.Exec(migrations.SeedSQL); err != nil {
				db.Close()
				log.Printf("Error executing seed data: %v", err)
				return nil, err
			}
		}
	}
	return repo, nil
}

func (r *SQLiteRepository) DropSchema() error {
	return os.Remove(r.filepath)
}
