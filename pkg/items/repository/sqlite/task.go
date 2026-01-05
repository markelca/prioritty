package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/markelca/prioritty/pkg/items"
)

func (r *SQLiteRepository) GetTasks() ([]items.Task, error) {
	query := `
		SELECT t.id, t.title, t.body, t.status_id, t.created_at, tag.id, tag.name
		FROM task t
			LEFT JOIN tag on t.tag_id = tag.id
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
		var taskId int
		var statusId int
		var tagId sql.NullInt64
		var tagName sql.NullString
		var createdAtStr string

		err := rows.Scan(&taskId, &task.Title, &body, &statusId, &createdAtStr, &tagId, &tagName)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		task.Id = strconv.Itoa(taskId)

		task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at string: %v", err)
			continue
		}

		if body != nil {
			task.Body = *body
		}

		if tagId.Valid {
			tag := items.Tag{
				Id:   strconv.FormatInt(tagId.Int64, 10),
				Name: tagName.String,
			}
			task.Tag = &tag
		} else {
			task.Tag = nil
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

func (r *SQLiteRepository) CreateTask(t *items.Task) error {
	query := `
		INSERT INTO task (title, body, status_id)
		VALUES (?, ?, ?)
	`
	result, err := r.db.Exec(query, t.Title, t.Body, t.Status)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.Id = fmt.Sprintf("%d", id)
	return nil
}

func (r *SQLiteRepository) RemoveTask(id string) error {
	query := `
		DELETE FROM task
		WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SQLiteRepository) SetTaskTag(t items.Task, tag items.Tag) error {
	query := `
		UPDATE task
		SET tag_id = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, tag.Id, t.Id)
	if err != nil {
		log.Printf("Error setting tag to task: %v", err)
		return err
	}
	return nil
}

func (r *SQLiteRepository) UnsetTaskTag(t items.Task) error {
	query := `
		UPDATE task
		SET tag_id = NULL
		WHERE id = ?
	`
	_, err := r.db.Exec(query, t.Id)
	if err != nil {
		log.Printf("Error unsetting tag from task: %v", err)
		return err
	}
	return nil
}
