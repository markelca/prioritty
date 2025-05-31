package tasks

import (
	"database/sql"
	"log"
	"os"

	"github.com/markelca/prioritty/migrations"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

type Repository interface {
	FindAll() ([]Task, error)
	UpdateStatus(Task, Status) error
	UpdateTask(Task) error
	CreateTask(Task) error
	RemoveTask(int) error
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

func (r *SQLiteRepository) UpdateTask(t Task) error {
	query := `
		UPDATE task
		SET title = ?, body = ?, status_id = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, t.Title, t.Body, t.Status, t.Id)
	return err
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
