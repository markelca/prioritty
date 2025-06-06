package repository

import (
	"log"
	"time"

	"github.com/markelca/prioritty/pkg/items"
)

func (r *SQLiteRepository) GetTasks() ([]items.Task, error) {
	query := `
		SELECT t.id, t.title, t.body, t.status_id, t.created_at
		FROM task t
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		return []items.Task{}, err
	}
	defer rows.Close()

	var tasks []items.Task

	for rows.Next() {
		var task items.Task
		var body *string
		var statusId int
		var createdAtStr string

		err := rows.Scan(&task.Id, &task.Title, &body, &statusId, &createdAtStr)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at string: %v", err)
			continue
		}

		if body != nil {
			task.Body = *body
		}
		task.Status = items.Status(statusId)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *SQLiteRepository) UpdateTask(t items.Task) error {
	query := `
		UPDATE task
		SET title = ?, body = ?, status_id = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, t.Title, t.Body, t.Status, t.Id)
	return err
}

func (r *SQLiteRepository) UpdateTaskStatus(t items.Task, s items.Status) error {
	query := `
		UPDATE task
		SET status_id = ?
		WHERE id = ?
	`
	r.db.Exec(query, s, t.Id)
	return nil
}

func (r *SQLiteRepository) CreateTask(t items.Task) error {
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
