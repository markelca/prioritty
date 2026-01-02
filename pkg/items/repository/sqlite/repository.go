package sqlite

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
)

type SQLiteRepository struct {
	db       *sql.DB
	filepath string
}

func NewSQLiteRepository(db *sql.DB, filepath string) *SQLiteRepository {
	return &SQLiteRepository{db: db, filepath: filepath}
}

func (r *SQLiteRepository) Reset() error {
	return os.Remove(r.filepath)
}

func (r *SQLiteRepository) GetTag(name string) (*items.Tag, error) {
	var tag items.Tag
	var id int
	query := `
		SELECT id, name
		FROM tag
		WHERE name = ?
	`
	row := r.db.QueryRow(query, name)

	err := row.Scan(&id, &tag.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		log.Printf("Error scanning tag: %v", err)
		return nil, err
	}
	tag.Id = strconv.Itoa(id)

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
		var id int
		err := rows.Scan(&id, &tag.Name)
		if err != nil {
			log.Printf("Error scanning tag: %v", err)
			return nil, err
		}
		tag.Id = strconv.Itoa(id)
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
		log.Printf("Error inserting tag: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last inserted tag id: %v", err)
		return nil, err
	}

	tag := items.Tag{
		Id:   strconv.FormatInt(id, 10),
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
		return repository.ErrNotFound
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
		var taskId int
		var tagId int
		var tagName string

		err := rows.Scan(&taskId, &task.Title, &body, &statusId, &createdAtStr, &tagId, &tagName)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			return nil, err
		}
		task.Id = strconv.Itoa(taskId)

		task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at string: %v", err)
			return nil, err
		}

		if body != nil {
			task.Body = *body
		}

		tag := items.Tag{
			Id:   strconv.Itoa(tagId),
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
		var noteId int
		var tagId int
		var tagName string

		err := rows.Scan(&noteId, &note.Title, &body, &createdAtStr, &tagId, &tagName)
		if err != nil {
			log.Printf("Error scanning note: %v", err)
			return nil, err
		}
		note.Id = strconv.Itoa(noteId)

		note.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at string: %v", err)
			return nil, err
		}

		if body != nil {
			note.Body = *body
		}

		tag := items.Tag{
			Id:   strconv.Itoa(tagId),
			Name: tagName,
		}
		note.Tag = &tag

		allItems = append(allItems, &note)
	}

	return allItems, nil
}
