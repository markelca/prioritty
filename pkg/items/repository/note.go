package repository

import (
	"log"
	"time"

	"github.com/markelca/prioritty/pkg/items"
)

func (r *SQLiteRepository) GetNotes() ([]items.Note, error) {
	query := `
		SELECT n.id, n.title, n.body, n.created_at
		FROM note n
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		return []items.Note{}, err
	}
	defer rows.Close()

	var notes []items.Note

	for rows.Next() {
		var note items.Note
		var body *string
		var createdAtStr string

		err := rows.Scan(&note.Id, &note.Title, &body, &createdAtStr)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		note.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at string: %v", err)
			continue
		}
		if body != nil {
			note.Body = *body
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *SQLiteRepository) UpdateNote(n items.Note) error {
	query := `
		UPDATE note
		SET title = ?, body = ? 
		WHERE id = ?
	`
	_, err := r.db.Exec(query, n.Title, n.Body, n.Id)
	return err
}

func (r *SQLiteRepository) CreateNote(n items.Note) error {
	query := `
		INSERT INTO note (title)
		VALUES (?)
	`
	_, err := r.db.Exec(query, n.Title)
	return err
}

func (r *SQLiteRepository) RemoveNote(id int) error {
	query := `
		DELETE FROM note
		WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
}
