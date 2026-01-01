package repository

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/markelca/prioritty/internal/migrations"
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
	GetTags() ([]items.Tag, error)
	CreateTag(string) (*items.Tag, error)
	RemoveTag(string) error
	GetItemsWithTag(string) ([]items.ItemInterface, error)
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

func (r *SQLiteRepository) GetTags() ([]items.Tag, error) {
	var tags []items.Tag
	query := `
		SELECT id, name
		FROM tag
		ORDER BY name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying tags: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag items.Tag
		err := rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			log.Printf("Error scanning tag: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over tags: %v", err)
		return nil, err
	}

	return tags, nil
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

func (r *SQLiteRepository) RemoveTag(name string) error {
	query := `DELETE FROM tag WHERE name = ?`
	result, err := r.db.Exec(query, name)
	if err != nil {
		log.Printf("Error removing tag: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLiteRepository) GetItemsWithTag(tagName string) ([]items.ItemInterface, error) {
	var allItems []items.ItemInterface

	tasksQuery := `
		SELECT t.id, t.title, t.body, t.status_id, t.created_at, tg.id, tg.name
		FROM task t
		JOIN tag tg ON t.tag_id = tg.id
		WHERE tg.name = ?
	`
	rows, err := r.db.Query(tasksQuery, tagName)
	if err != nil {
		log.Printf("Error querying tasks with tag: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task items.Task
		var body *string
		var statusId int
		var createdAtStr string
		var tagId int
		var tagName string

		err := rows.Scan(&task.Id, &task.Title, &body, &statusId, &createdAtStr, &tagId, &tagName)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			return nil, err
		}

		task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at string: %v", err)
			return nil, err
		}

		if body != nil {
			task.Body = *body
		}

		tag := items.Tag{
			Id:   tagId,
			Name: tagName,
		}
		task.Tag = &tag
		task.Status = items.Status(statusId)

		allItems = append(allItems, &task)
	}

	notesQuery := `
		SELECT n.id, n.title, n.body, n.created_at, tg.id, tg.name
		FROM note n
		JOIN tag tg ON n.tag_id = tg.id
		WHERE tg.name = ?
	`
	rows, err = r.db.Query(notesQuery, tagName)
	if err != nil {
		log.Printf("Error querying notes with tag: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var note items.Note
		var body *string
		var createdAtStr string
		var tagId int
		var tagName string

		err := rows.Scan(&note.Id, &note.Title, &body, &createdAtStr, &tagId, &tagName)
		if err != nil {
			log.Printf("Error scanning note: %v", err)
			return nil, err
		}

		note.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at string: %v", err)
			return nil, err
		}

		if body != nil {
			note.Body = *body
		}

		tag := items.Tag{
			Id:   tagId,
			Name: tagName,
		}
		note.Tag = &tag

		allItems = append(allItems, &note)
	}

	return allItems, nil
}
