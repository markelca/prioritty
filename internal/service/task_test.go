package service

import (
	"errors"
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskService_GetTasks(t *testing.T) {
	t.Run("returns all tasks", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task1 := items.Task{Item: items.Item{Id: "1", Title: "Task 1"}, Status: items.Todo}
		task2 := items.Task{Item: items.Item{Id: "2", Title: "Task 2"}, Status: items.Done}
		mockRepo.AddTask(task1)
		mockRepo.AddTask(task2)
		svc := NewService(mockRepo)

		tasks, err := svc.GetTasks()

		require.NoError(t, err)
		assert.Len(t, tasks, 2)
		assert.True(t, mockRepo.HasCall("GetTasks"))
	})

	t.Run("returns empty slice when no tasks", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		tasks, err := svc.GetTasks()

		require.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("propagates error", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		mockRepo.GetTasksError = errors.New("db error")
		svc := NewService(mockRepo)

		_, err := svc.GetTasks()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	t.Run("updates task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{Item: items.Item{Id: "1", Title: "Original"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		task.Title = "Updated"
		err := svc.UpdateTask(task)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("UpdateTask"))
	})

	t.Run("propagates error", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		mockRepo.UpdateTaskError = errors.New("update error")
		task := items.Task{Item: items.Item{Id: "1"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.UpdateTask(task)

		require.Error(t, err)
	})
}

func TestTaskService_UpdateStatus(t *testing.T) {
	t.Run("updates status from todo to done", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{Item: items.Item{Id: "1", Title: "Task"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.UpdateStatus(&task, items.Done)

		require.NoError(t, err)
		assert.Equal(t, items.Done, task.Status)
		assert.True(t, mockRepo.HasCall("UpdateTaskStatus"))
	})

	t.Run("toggles back to todo when same status", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{Item: items.Item{Id: "1", Title: "Task"}, Status: items.Done}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.UpdateStatus(&task, items.Done) // Same status

		require.NoError(t, err)
		assert.Equal(t, items.Todo, task.Status) // Toggled to todo
	})

	t.Run("propagates repository error", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		mockRepo.UpdateTaskStatusError = errors.New("status error")
		task := items.Task{Item: items.Item{Id: "1"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.UpdateStatus(&task, items.Done)

		require.Error(t, err)
	})
}

func TestTaskService_AddTask(t *testing.T) {
	t.Run("creates task with title", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		err := svc.AddTask("New Task")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("CreateTask"))
		assert.Equal(t, 1, mockRepo.TaskCount())
	})

	t.Run("propagates error", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		mockRepo.CreateTaskError = errors.New("create error")
		svc := NewService(mockRepo)

		err := svc.AddTask("Task")

		require.Error(t, err)
	})
}

func TestTaskService_DestroyDemo(t *testing.T) {
	mockRepo := repository.NewMockRepository()
	svc := NewService(mockRepo)

	err := svc.DestroyDemo()

	require.NoError(t, err)
	assert.True(t, mockRepo.HasCall("Reset"))
}

func TestTaskService_removeTask(t *testing.T) {
	t.Run("removes existing task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{Item: items.Item{Id: "1", Title: "Task"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.removeTask("1")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("RemoveTask"))
		assert.Equal(t, 0, mockRepo.TaskCount())
	})

	t.Run("returns error for non-existent task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		err := svc.removeTask("non-existent")

		require.Error(t, err)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}
