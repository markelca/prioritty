package sqlite

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
)

// Schema SQL for creating the test database
const testSchemaSQL = `
CREATE TABLE status (
   id INTEGER PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE tag (
   id INTEGER PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE task (
   id INTEGER PRIMARY KEY,
   title TEXT NOT NULL,
   body TEXT,
   status_id INTEGER NOT NULL,
   tag_id INTEGER,
   created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (status_id) REFERENCES status(id)
   FOREIGN KEY (tag_id) REFERENCES tag(id)
);

CREATE TABLE note (
   id INTEGER PRIMARY KEY,
   title TEXT NOT NULL,
   body TEXT,
   tag_id INTEGER,
   created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (tag_id) REFERENCES tag(id)
);

INSERT INTO status (id, name) VALUES
   (0, 'Pending'),
   (1, 'In Progress'),
   (2, 'Completed'),
   (3, 'Cancelled');
`

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *SQLiteRepository {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(testSchemaSQL)
	require.NoError(t, err)

	repo := NewSQLiteRepository(db, ":memory:")

	t.Cleanup(func() {
		db.Close()
	})

	return repo
}

// --- Task Tests ---

func TestSQLiteRepository_CreateTask(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{
		Item:   items.Item{Title: "Test Task", Body: "Task body"},
		Status: items.Todo,
	}

	err := repo.CreateTask(task)

	require.NoError(t, err)
	assert.NotEmpty(t, task.Id, "ID should be assigned after creation")
}

func TestSQLiteRepository_GetTasks(t *testing.T) {
	repo := setupTestDB(t)

	// Create some tasks
	task1 := &items.Task{Item: items.Item{Title: "Task 1"}, Status: items.Todo}
	task2 := &items.Task{Item: items.Item{Title: "Task 2"}, Status: items.Done}
	require.NoError(t, repo.CreateTask(task1))
	require.NoError(t, repo.CreateTask(task2))

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Len(t, tasks, 2)
}

func TestSQLiteRepository_GetTasks_Empty(t *testing.T) {
	repo := setupTestDB(t)

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestSQLiteRepository_UpdateTask(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{Item: items.Item{Title: "Original"}, Status: items.Todo}
	require.NoError(t, repo.CreateTask(task))

	task.Title = "Updated"
	task.Body = "New body"
	err := repo.UpdateTask(*task)

	require.NoError(t, err)

	// Verify update
	tasks, _ := repo.GetTasks()
	assert.Equal(t, "Updated", tasks[0].Title)
	assert.Equal(t, "New body", tasks[0].Body)
}

func TestSQLiteRepository_UpdateTaskStatus(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{Item: items.Item{Title: "Task"}, Status: items.Todo}
	require.NoError(t, repo.CreateTask(task))

	err := repo.UpdateTaskStatus(*task, items.Done)

	require.NoError(t, err)

	// Verify status change
	tasks, _ := repo.GetTasks()
	assert.Equal(t, items.Done, tasks[0].Status)
}

func TestSQLiteRepository_RemoveTask(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{Item: items.Item{Title: "Task"}, Status: items.Todo}
	require.NoError(t, repo.CreateTask(task))

	err := repo.RemoveTask(task.Id)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Empty(t, tasks)
}

// --- Note Tests ---

func TestSQLiteRepository_CreateNote(t *testing.T) {
	repo := setupTestDB(t)

	note := &items.Note{
		Item: items.Item{Title: "Test Note", Body: "Note body"},
	}

	err := repo.CreateNote(note)

	require.NoError(t, err)
	assert.NotEmpty(t, note.Id, "ID should be assigned after creation")
}

func TestSQLiteRepository_GetNotes(t *testing.T) {
	repo := setupTestDB(t)

	note1 := &items.Note{Item: items.Item{Title: "Note 1"}}
	note2 := &items.Note{Item: items.Item{Title: "Note 2"}}
	require.NoError(t, repo.CreateNote(note1))
	require.NoError(t, repo.CreateNote(note2))

	notes, err := repo.GetNotes()

	require.NoError(t, err)
	assert.Len(t, notes, 2)
}

func TestSQLiteRepository_GetNotes_Empty(t *testing.T) {
	repo := setupTestDB(t)

	notes, err := repo.GetNotes()

	require.NoError(t, err)
	assert.Empty(t, notes)
}

func TestSQLiteRepository_UpdateNote(t *testing.T) {
	repo := setupTestDB(t)

	note := &items.Note{Item: items.Item{Title: "Original"}}
	require.NoError(t, repo.CreateNote(note))

	note.Title = "Updated"
	note.Body = "New body"
	err := repo.UpdateNote(*note)

	require.NoError(t, err)

	// Verify update
	notes, _ := repo.GetNotes()
	assert.Equal(t, "Updated", notes[0].Title)
	assert.Equal(t, "New body", notes[0].Body)
}

func TestSQLiteRepository_RemoveNote(t *testing.T) {
	repo := setupTestDB(t)

	note := &items.Note{Item: items.Item{Title: "Note"}}
	require.NoError(t, repo.CreateNote(note))

	err := repo.RemoveNote(note.Id)

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	assert.Empty(t, notes)
}

// --- Tag Tests ---

func TestSQLiteRepository_CreateTag(t *testing.T) {
	repo := setupTestDB(t)

	tag, err := repo.CreateTag("work")

	require.NoError(t, err)
	assert.NotEmpty(t, tag.Id)
	assert.Equal(t, "work", tag.Name)
}

func TestSQLiteRepository_GetTag(t *testing.T) {
	repo := setupTestDB(t)

	created, err := repo.CreateTag("work")
	require.NoError(t, err)

	found, err := repo.GetTag("work")

	require.NoError(t, err)
	assert.Equal(t, created.Id, found.Id)
	assert.Equal(t, "work", found.Name)
}

func TestSQLiteRepository_GetTag_NotFound(t *testing.T) {
	repo := setupTestDB(t)

	_, err := repo.GetTag("nonexistent")

	require.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestSQLiteRepository_GetTags(t *testing.T) {
	repo := setupTestDB(t)

	repo.CreateTag("work")
	repo.CreateTag("docs")
	repo.CreateTag("urgent")

	tags, err := repo.GetTags()

	require.NoError(t, err)
	assert.Len(t, tags, 3)
}

func TestSQLiteRepository_RemoveTag(t *testing.T) {
	repo := setupTestDB(t)

	repo.CreateTag("work")

	err := repo.RemoveTag("work")

	require.NoError(t, err)

	_, err = repo.GetTag("work")
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestSQLiteRepository_RemoveTag_NotFound(t *testing.T) {
	repo := setupTestDB(t)

	err := repo.RemoveTag("nonexistent")

	require.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

// --- Task Tag Tests ---

func TestSQLiteRepository_SetTaskTag(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{Item: items.Item{Title: "Task"}, Status: items.Todo}
	require.NoError(t, repo.CreateTask(task))

	tag, err := repo.CreateTag("work")
	require.NoError(t, err)

	err = repo.SetTaskTag(*task, *tag)

	require.NoError(t, err)

	// Verify tag is set
	tasks, _ := repo.GetTasks()
	require.NotNil(t, tasks[0].Tag)
	assert.Equal(t, "work", tasks[0].Tag.Name)
}

func TestSQLiteRepository_UnsetTaskTag(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{Item: items.Item{Title: "Task"}, Status: items.Todo}
	require.NoError(t, repo.CreateTask(task))

	tag, _ := repo.CreateTag("work")
	repo.SetTaskTag(*task, *tag)

	err := repo.UnsetTaskTag(*task)

	require.NoError(t, err)

	// Verify tag is removed
	tasks, _ := repo.GetTasks()
	assert.Nil(t, tasks[0].Tag)
}

// --- Note Tag Tests ---

func TestSQLiteRepository_SetNoteTag(t *testing.T) {
	repo := setupTestDB(t)

	note := &items.Note{Item: items.Item{Title: "Note"}}
	require.NoError(t, repo.CreateNote(note))

	tag, err := repo.CreateTag("docs")
	require.NoError(t, err)

	err = repo.SetNoteTag(*note, *tag)

	require.NoError(t, err)

	// Verify tag is set
	notes, _ := repo.GetNotes()
	require.NotNil(t, notes[0].Tag)
	assert.Equal(t, "docs", notes[0].Tag.Name)
}

func TestSQLiteRepository_UnsetNoteTag(t *testing.T) {
	repo := setupTestDB(t)

	note := &items.Note{Item: items.Item{Title: "Note"}}
	require.NoError(t, repo.CreateNote(note))

	tag, _ := repo.CreateTag("docs")
	repo.SetNoteTag(*note, *tag)

	err := repo.UnsetNoteTag(*note)

	require.NoError(t, err)

	// Verify tag is removed
	notes, _ := repo.GetNotes()
	assert.Nil(t, notes[0].Tag)
}

// --- GetItemsWithTag Tests ---

func TestSQLiteRepository_GetItemsWithTag(t *testing.T) {
	repo := setupTestDB(t)

	// Create tag
	tag, _ := repo.CreateTag("work")

	// Create task with tag
	task := &items.Task{Item: items.Item{Title: "Work Task"}, Status: items.Todo}
	repo.CreateTask(task)
	repo.SetTaskTag(*task, *tag)

	// Create note with tag
	note := &items.Note{Item: items.Item{Title: "Work Note"}}
	repo.CreateNote(note)
	repo.SetNoteTag(*note, *tag)

	// Create items without tag
	otherTask := &items.Task{Item: items.Item{Title: "Other Task"}, Status: items.Todo}
	repo.CreateTask(otherTask)

	items, err := repo.GetItemsWithTag("work")

	require.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestSQLiteRepository_GetItemsWithTag_Empty(t *testing.T) {
	repo := setupTestDB(t)

	repo.CreateTag("unused")

	items, err := repo.GetItemsWithTag("unused")

	require.NoError(t, err)
	assert.Empty(t, items)
}

// --- Status Mapping Tests ---

func TestStatusIdToStatus(t *testing.T) {
	tests := []struct {
		id       int
		expected items.Status
	}{
		{0, items.Todo},
		{1, items.InProgress},
		{2, items.Done},
		{3, items.Cancelled},
		{99, items.Todo}, // Unknown defaults to Todo
		{-1, items.Todo}, // Negative defaults to Todo
	}

	for _, tc := range tests {
		result := statusIdToStatus(tc.id)
		assert.Equal(t, tc.expected, result, "statusIdToStatus(%d) should return %s", tc.id, tc.expected)
	}
}

// --- Full CRUD Cycle Test ---

func TestSQLiteRepository_TaskCRUDCycle(t *testing.T) {
	repo := setupTestDB(t)

	// Create
	task := &items.Task{
		Item:   items.Item{Title: "CRUD Task", Body: "Initial body"},
		Status: items.Todo,
	}
	err := repo.CreateTask(task)
	require.NoError(t, err)
	taskId := task.Id

	// Read
	tasks, err := repo.GetTasks()
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	assert.Equal(t, "CRUD Task", tasks[0].Title)

	// Update
	task.Title = "Updated CRUD Task"
	task.Body = "Updated body"
	err = repo.UpdateTask(*task)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, "Updated CRUD Task", tasks[0].Title)

	// Update Status
	err = repo.UpdateTaskStatus(*task, items.Done)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, items.Done, tasks[0].Status)

	// Delete
	err = repo.RemoveTask(taskId)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Empty(t, tasks)
}

func TestSQLiteRepository_NoteCRUDCycle(t *testing.T) {
	repo := setupTestDB(t)

	// Create
	note := &items.Note{
		Item: items.Item{Title: "CRUD Note", Body: "Initial body"},
	}
	err := repo.CreateNote(note)
	require.NoError(t, err)
	noteId := note.Id

	// Read
	notes, err := repo.GetNotes()
	require.NoError(t, err)
	require.Len(t, notes, 1)
	assert.Equal(t, "CRUD Note", notes[0].Title)

	// Update
	note.Title = "Updated CRUD Note"
	note.Body = "Updated body"
	err = repo.UpdateNote(*note)
	require.NoError(t, err)

	notes, _ = repo.GetNotes()
	assert.Equal(t, "Updated CRUD Note", notes[0].Title)

	// Delete
	err = repo.RemoveNote(noteId)
	require.NoError(t, err)

	notes, _ = repo.GetNotes()
	assert.Empty(t, notes)
}
