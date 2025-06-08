package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/markelca/prioritty/pkg/items"
)

func (r *SQLiteRepository) GetNotes() ([]items.Note, error) {
	query := `
		SELECT n.id, n.title, n.body, n.created_at, tag.id, tag.name
		FROM note n
			LEFT JOIN tag on n.tag_id = tag.id
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
		var tagId sql.NullInt64
		var tagName sql.NullString
		var createdAtStr string

		err := rows.Scan(&note.Id, &note.Title, &body, &createdAtStr, &tagId, &tagName)
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

		if tagId.Valid {
			tag := items.Tag{
				Id:   int(tagId.Int64),
				Name: tagName.String,
			}
			note.Tag = &tag // assuming Task.Tag is *Tag
		} else {
			note.Tag = nil
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

func (r *SQLiteRepository) SetNoteTag(n items.Note, tag items.Tag) error {
	query := `
		UPDATE note
		SET tag_id = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, tag.Id, n.Id)
	if err != nil {
		log.Printf("Error setting tag to note: %v", err)
		return err
	}
	return nil
}
