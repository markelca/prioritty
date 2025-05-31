package tasks

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

type Repository interface {
	FindAll() ([]Task, error)
	UpdateStatus(Task, Status) error
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
		if err := repo.initSQLiteDatabase(); err != nil {
			db.Close()
			return nil, err
		}
		if viper.GetBool("demo") {
			if err := repo.insertDefaultData(); err != nil {
				db.Close()
				return nil, err
			}
		}
	}

	return repo, nil
}

func (r SQLiteRepository) initSQLiteDatabase() error {
	// Create status table
	createStatusTable := `
	CREATE TABLE status (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);`

	if _, err := r.db.Exec(createStatusTable); err != nil {
		log.Printf("Error creating status table: %v", err)
		return err
	}

	// Create task table
	createTaskTable := `
	CREATE TABLE task (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		body TEXT,
		status_id INTEGER NOT NULL,
		FOREIGN KEY (status_id) REFERENCES status(id)
	);`

	if _, err := r.db.Exec(createTaskTable); err != nil {
		log.Printf("Error creating task table: %v", err)
		return err
	}

	// Insert initial status data
	insertStatuses := `
	INSERT INTO status (id, name) VALUES
		(0, 'Pending'),
		(1, 'In Progress'),
		(2, 'Completed'),
		(3, 'Cancelled');`

	if _, err := r.db.Exec(insertStatuses); err != nil {
		log.Printf("Error inserting status data: %v", err)
		return err
	}

	// log.Println("Database initialized successfully")
	return nil
}

func (r SQLiteRepository) insertDefaultData() error {
	// Insert initial task data
	insertTasks := `
	INSERT INTO task (title, body, status_id) VALUES 
		('Complete project documentation', 'Write comprehensive documentation covering all API endpoints, authentication methods, and usage examples. Include code samples in multiple languages and ensure all examples are tested and working. The documentation should be organized into clear sections: Getting Started, Authentication, Core API Reference, Advanced Features, and Troubleshooting. Each endpoint should include request/response examples, parameter descriptions, and common error codes. Add interactive examples where possible and ensure the documentation is accessible to both beginner and advanced developers.', 2),
		('Review code changes', NULL, 1),
		('Fix bug in authentication', NULL, 0),
		('Deploy to production', NULL, 2),
		('Write unit tests', NULL, 0),
		('Update dependencies', NULL, 3);`

	if _, err := r.db.Exec(insertTasks); err != nil {
		log.Printf("Error inserting task data: %v", err)
		return err
	}
	return nil
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
