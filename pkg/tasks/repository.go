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
	RemoveTask(int) error
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
		SELECT t.id, t.title, t.body, t.status_id 
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
		var body *string
		var statusId int

		err := rows.Scan(&task.Id, &task.Title, &body, &statusId)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}

		if body != nil {
			task.Body = *body
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

func (r *SQLiteRepository) RemoveTask(id int) error {
	query := `
		DELETE FROM task
		WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
}
