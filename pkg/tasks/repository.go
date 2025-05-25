package tasks

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	FindAll() ([]Task, error)
	UpdateStatus(Task, Status) error
	CreateTask(Task) error
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) FindAll() ([]Task, error) {
	query := `
		SELECT t.id, t.title, t.status_id 
		FROM task t
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		return []Task{}, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var statusId int

		err := rows.Scan(&task.Id, &task.Title, &statusId)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}

		task.Status = Status(statusId)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *SQLiteRepository) UpdateStatus(t Task, s Status) error {
	query := `
		UPDATE task
		SET status_id = ?
		WHERE id = ?
	`
	r.db.Exec(query, s, t.Id)
	return nil
}

func (r *SQLiteRepository) CreateTask(t Task) error {
	query := `
		INSERT INTO task (title, status_id)
		VALUES (?, ?)
	`
	_, err := r.db.Exec(query, t.Title, t.Status)
	return err
}
