package sqlite

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestSQLiteRepository_SetTaskTag(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{Item: items.Item{Title: "Task"}, Status: items.Todo}
	require.NoError(t, repo.CreateTask(task))

	tag, err := repo.CreateTag("work")
	require.NoError(t, err)

	err = repo.SetTaskTag(*task, *tag)

	require.NoError(t, err)

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

	tasks, _ := repo.GetTasks()
	assert.Nil(t, tasks[0].Tag)
}

func TestSQLiteRepository_TaskCRUDCycle(t *testing.T) {
	repo := setupTestDB(t)

	task := &items.Task{
		Item:   items.Item{Title: "CRUD Task", Body: "Initial body"},
		Status: items.Todo,
	}

	err := repo.CreateTask(task)
	require.NoError(t, err)
	taskId := task.Id

	tasks, err := repo.GetTasks()
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	assert.Equal(t, "CRUD Task", tasks[0].Title)

	tasks[0].Title = "Updated CRUD Task"
	tasks[0].Body = "Updated body"
	err = repo.UpdateTask(tasks[0])
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, "Updated CRUD Task", tasks[0].Title)

	err = repo.UpdateTaskStatus(tasks[0], items.Done)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, items.Done, tasks[0].Status)

	err = repo.RemoveTask(taskId)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Empty(t, tasks)
}
