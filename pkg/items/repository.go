package items

import (
	"database/sql"
)

type Repository interface {
	GetAll() ([]Item, error)
	GetTasks()
	GetNotes()
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
