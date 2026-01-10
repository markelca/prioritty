package obsidian

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObsidianRepository_CreateTask(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "New Task", Body: "Task content"},
		Status: items.Todo,
	}

	err := repo.CreateTask(task)

	require.NoError(t, err)
	assert.NotEmpty(t, task.Id)
	assert.Contains(t, task.Id, ".md")
}

func TestObsidianRepository_GetTasks(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Task 1
status: todo
---
Task 1 body`)

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task 1", tasks[0].Title)
}

func TestObsidianRepository_GetTasks_Empty(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestObsidianRepository_GetTasks_IgnoresNotes(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "note1.md", `---
type: note
title: Note 1
---
`)

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestObsidianRepository_UpdateTask(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Original", Body: "body"},
		Status: items.Todo,
	}
	require.NoError(t, repo.CreateTask(task))

	task.Title = "Updated"
	task.Body = "New body"
	task.Status = items.Done
	err := repo.UpdateTask(*task)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Equal(t, "Updated", tasks[0].Title)
	assert.Equal(t, "New body\n", tasks[0].Body)
	assert.Equal(t, items.Done, tasks[0].Status)
}

func TestObsidianRepository_RemoveTask(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task to remove"},
		Status: items.Todo,
	}
	require.NoError(t, repo.CreateTask(task))

	err := repo.RemoveTask(task.Id)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Empty(t, tasks)
}

func TestObsidianRepository_UpdateTaskStatus(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task"},
		Status: items.Todo,
	}
	require.NoError(t, repo.CreateTask(task))

	err := repo.UpdateTaskStatus(*task, items.InProgress)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Equal(t, items.InProgress, tasks[0].Status)
}

func TestObsidianRepository_SetTaskTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task"},
		Status: items.Todo,
	}
	require.NoError(t, repo.CreateTask(task))

	err := repo.SetTaskTag(*task, items.Tag{Id: "work", Name: "work"})

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	require.NotNil(t, tasks[0].Tag)
	assert.Equal(t, "work", tasks[0].Tag.Name)
}

func TestObsidianRepository_UnsetTaskTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task"},
		Status: items.Todo,
	}
	require.NoError(t, repo.CreateTask(task))
	require.NoError(t, repo.SetTaskTag(*task, items.Tag{Id: "work", Name: "work"}))

	err := repo.UnsetTaskTag(*task)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Nil(t, tasks[0].Tag)
}

func TestObsidianRepository_TaskCRUDCycle(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "CRUD Task", Body: "Initial body"},
		Status: items.Todo,
	}

	err := repo.CreateTask(task)
	require.NoError(t, err)

	tasks, err := repo.GetTasks()
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	assert.Equal(t, "CRUD Task", tasks[0].Title)

	tasks[0].Title = "Updated Task"
	tasks[0].Body = "Updated body"
	err = repo.UpdateTask(tasks[0])
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, "Updated Task", tasks[0].Title)

	err = repo.UpdateTaskStatus(tasks[0], items.Done)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, items.Done, tasks[0].Status)

	err = repo.RemoveTask(tasks[0].Id)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Empty(t, tasks)
}
