package notes

import (
	"database/sql"
	"log"
)

type Repository interface {
	FindAll() ([]Note, error)
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

func (r *SQLiteRepository) FindAll() ([]Note, error) {
	query := `
		SELECT n.id, n.title, n.body
		FROM note n
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying notes: %v", err)
		return []Note{}, err
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var note Note
		var body *string

		err := rows.Scan(&note.Id, &note.Title, &body)
		if err != nil {
			log.Printf("Error scanning note: %v", err)
			continue
		}

		if body != nil {
			note.Body = *body
		}
		notes = append(notes, note)
	}

	return notes, nil
}
