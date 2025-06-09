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
	SetTaskTag(items.Task, items.Tag) error
	UnsetTaskTag(items.Task) error
}

type NoteRepository interface {
	GetNotes() ([]items.Note, error)
	UpdateNote(items.Note) error
	CreateNote(items.Note) error
	RemoveNote(int) error
	SetNoteTag(items.Note, items.Tag) error
	UnsetNoteTag(items.Note) error
}

type Repository interface {
	TaskRepository
	NoteRepository
	GetTag(string) (*items.Tag, error)
	CreateTag(string) (*items.Tag, error)
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

func (r *SQLiteRepository) GetTag(name string) (*items.Tag, error) {
	var tag items.Tag
	query := `
		SELECT id, name
		FROM tag
		WHERE name = ?
	`
	row := r.db.QueryRow(query, name)

	err := row.Scan(&tag.Id, &tag.Name)
	if err != nil {
		log.Printf("Error scanning task: %v", err)
		return nil, err
	}

	return &tag, nil
}

func (r *SQLiteRepository) CreateTag(name string) (*items.Tag, error) {
	query := `
		INSERT INTO tag (name)
		VALUES (?)
	`
	result, err := r.db.Exec(query, name)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last inserted tag id: %v", err)
		return nil, err
	}

	tag := items.Tag{
		Id:   int(id),
		Name: name,
	}

	return &tag, nil
}
