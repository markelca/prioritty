package testutil

import (
	"errors"
	"testing"
	"time"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFakeRepository(t *testing.T) {
	repo := NewFakeRepository()

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.tasks)
	assert.NotNil(t, repo.notes)
	assert.NotNil(t, repo.tags)
	assert.NotNil(t, repo.Calls)
	assert.Empty(t, repo.Calls)
	assert.Equal(t, 0, repo.TaskCount())
	assert.Equal(t, 0, repo.NoteCount())
	assert.Equal(t, 0, repo.TagCount())
}

func TestFakeRepository_CallTracking(t *testing.T) {
	repo := NewFakeRepository()

	t.Run("HasCall returns false for untracked method", func(t *testing.T) {
		assert.False(t, repo.HasCall("GetTasks"))
	})

	t.Run("HasCall returns true after method is called", func(t *testing.T) {
		_, _ = repo.GetTasks()
		assert.True(t, repo.HasCall("GetTasks"))
	})

	t.Run("CallCount returns correct count", func(t *testing.T) {
		_, _ = repo.GetTasks()
		_, _ = repo.GetTasks()
		assert.Equal(t, 3, repo.CallCount("GetTasks"))
	})

	t.Run("ResetCalls clears tracking", func(t *testing.T) {
		repo.ResetCalls()
		assert.Empty(t, repo.Calls)
		assert.False(t, repo.HasCall("GetTasks"))
	})
}

func TestFakeRepository_TaskOperations(t *testing.T) {
	repo := NewFakeRepository()
	now := time.Now()

	task := items.Task{
		Item: items.Item{
			Id:        "task-1",
			Title:     "Test Task",
			Body:      "Task body",
			CreatedAt: now,
		},
		Status: items.Todo,
	}

	t.Run("CreateTask", func(t *testing.T) {
		err := repo.CreateTask(&task)
		require.NoError(t, err)
		assert.True(t, repo.HasCall("CreateTask"))
		assert.Equal(t, 1, repo.TaskCount())
	})

	t.Run("GetTasks", func(t *testing.T) {
		tasks, err := repo.GetTasks()
		require.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, task.Id, tasks[0].Id)
		assert.Equal(t, task.Title, tasks[0].Title)
	})

	t.Run("GetTaskById", func(t *testing.T) {
		found, exists := repo.GetTaskById("task-1")
		assert.True(t, exists)
		assert.Equal(t, "Test Task", found.Title)
	})

	t.Run("GetTaskById not found", func(t *testing.T) {
		_, exists := repo.GetTaskById("nonexistent")
		assert.False(t, exists)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		updatedTask := task
		updatedTask.Title = "Updated Task"
		err := repo.UpdateTask(updatedTask)
		require.NoError(t, err)

		found, _ := repo.GetTaskById("task-1")
		assert.Equal(t, "Updated Task", found.Title)
	})

	t.Run("UpdateTask not found", func(t *testing.T) {
		nonExistentTask := items.Task{Item: items.Item{Id: "nonexistent"}}
		err := repo.UpdateTask(nonExistentTask)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("UpdateTaskStatus", func(t *testing.T) {
		err := repo.UpdateTaskStatus(task, items.Done)
		require.NoError(t, err)

		found, _ := repo.GetTaskById("task-1")
		assert.Equal(t, items.Done, found.Status)
	})

	t.Run("UpdateTaskStatus not found", func(t *testing.T) {
		nonExistentTask := items.Task{Item: items.Item{Id: "nonexistent"}}
		err := repo.UpdateTaskStatus(nonExistentTask, items.Done)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("RemoveTask", func(t *testing.T) {
		err := repo.RemoveTask("task-1")
		require.NoError(t, err)
		assert.Equal(t, 0, repo.TaskCount())
	})

	t.Run("RemoveTask not found", func(t *testing.T) {
		err := repo.RemoveTask("nonexistent")
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestFakeRepository_TaskTagOperations(t *testing.T) {
	repo := NewFakeRepository()
	tag := items.Tag{Id: "tag-1", Name: "work"}
	task := items.Task{
		Item:   items.Item{Id: "task-1", Title: "Task"},
		Status: items.Todo,
	}
	repo.AddTask(task)

	t.Run("SetTaskTag", func(t *testing.T) {
		err := repo.SetTaskTag(task, tag)
		require.NoError(t, err)

		found, _ := repo.GetTaskById("task-1")
		require.NotNil(t, found.Tag)
		assert.Equal(t, "work", found.Tag.Name)
	})

	t.Run("UnsetTaskTag", func(t *testing.T) {
		err := repo.UnsetTaskTag(task)
		require.NoError(t, err)

		found, _ := repo.GetTaskById("task-1")
		assert.Nil(t, found.Tag)
	})

	t.Run("SetTaskTag not found", func(t *testing.T) {
		nonExistentTask := items.Task{Item: items.Item{Id: "nonexistent"}}
		err := repo.SetTaskTag(nonExistentTask, tag)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("UnsetTaskTag not found", func(t *testing.T) {
		nonExistentTask := items.Task{Item: items.Item{Id: "nonexistent"}}
		err := repo.UnsetTaskTag(nonExistentTask)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestFakeRepository_NoteOperations(t *testing.T) {
	repo := NewFakeRepository()
	now := time.Now()

	note := items.Note{
		Item: items.Item{
			Id:        "note-1",
			Title:     "Test Note",
			Body:      "Note body",
			CreatedAt: now,
		},
	}

	t.Run("CreateNote", func(t *testing.T) {
		err := repo.CreateNote(&note)
		require.NoError(t, err)
		assert.True(t, repo.HasCall("CreateNote"))
		assert.Equal(t, 1, repo.NoteCount())
	})

	t.Run("GetNotes", func(t *testing.T) {
		notes, err := repo.GetNotes()
		require.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, note.Id, notes[0].Id)
	})

	t.Run("GetNoteById", func(t *testing.T) {
		found, exists := repo.GetNoteById("note-1")
		assert.True(t, exists)
		assert.Equal(t, "Test Note", found.Title)
	})

	t.Run("GetNoteById not found", func(t *testing.T) {
		_, exists := repo.GetNoteById("nonexistent")
		assert.False(t, exists)
	})

	t.Run("UpdateNote", func(t *testing.T) {
		updatedNote := note
		updatedNote.Title = "Updated Note"
		err := repo.UpdateNote(updatedNote)
		require.NoError(t, err)

		found, _ := repo.GetNoteById("note-1")
		assert.Equal(t, "Updated Note", found.Title)
	})

	t.Run("UpdateNote not found", func(t *testing.T) {
		nonExistentNote := items.Note{Item: items.Item{Id: "nonexistent"}}
		err := repo.UpdateNote(nonExistentNote)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("RemoveNote", func(t *testing.T) {
		err := repo.RemoveNote("note-1")
		require.NoError(t, err)
		assert.Equal(t, 0, repo.NoteCount())
	})

	t.Run("RemoveNote not found", func(t *testing.T) {
		err := repo.RemoveNote("nonexistent")
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestFakeRepository_NoteTagOperations(t *testing.T) {
	repo := NewFakeRepository()
	tag := items.Tag{Id: "tag-1", Name: "personal"}
	note := items.Note{
		Item: items.Item{Id: "note-1", Title: "Note"},
	}
	repo.AddNote(note)

	t.Run("SetNoteTag", func(t *testing.T) {
		err := repo.SetNoteTag(note, tag)
		require.NoError(t, err)

		found, _ := repo.GetNoteById("note-1")
		require.NotNil(t, found.Tag)
		assert.Equal(t, "personal", found.Tag.Name)
	})

	t.Run("UnsetNoteTag", func(t *testing.T) {
		err := repo.UnsetNoteTag(note)
		require.NoError(t, err)

		found, _ := repo.GetNoteById("note-1")
		assert.Nil(t, found.Tag)
	})

	t.Run("SetNoteTag not found", func(t *testing.T) {
		nonExistentNote := items.Note{Item: items.Item{Id: "nonexistent"}}
		err := repo.SetNoteTag(nonExistentNote, tag)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("UnsetNoteTag not found", func(t *testing.T) {
		nonExistentNote := items.Note{Item: items.Item{Id: "nonexistent"}}
		err := repo.UnsetNoteTag(nonExistentNote)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestFakeRepository_TagOperations(t *testing.T) {
	repo := NewFakeRepository()

	t.Run("CreateTag", func(t *testing.T) {
		tag, err := repo.CreateTag("work")
		require.NoError(t, err)
		assert.Equal(t, "work", tag.Name)
		assert.Contains(t, tag.Id, "mock-tag-")
		assert.Equal(t, 1, repo.TagCount())
	})

	t.Run("GetTag", func(t *testing.T) {
		tag, err := repo.GetTag("work")
		require.NoError(t, err)
		assert.Equal(t, "work", tag.Name)
	})

	t.Run("GetTag not found", func(t *testing.T) {
		_, err := repo.GetTag("nonexistent")
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("GetTags", func(t *testing.T) {
		_, _ = repo.CreateTag("personal")
		tags, err := repo.GetTags()
		require.NoError(t, err)
		assert.Len(t, tags, 2)
	})

	t.Run("RemoveTag", func(t *testing.T) {
		err := repo.RemoveTag("work")
		require.NoError(t, err)
		assert.Equal(t, 1, repo.TagCount())
	})

	t.Run("RemoveTag not found", func(t *testing.T) {
		err := repo.RemoveTag("nonexistent")
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestFakeRepository_GetItemsWithTag(t *testing.T) {
	repo := NewFakeRepository()

	tag := items.Tag{Id: "tag-1", Name: "work"}
	repo.AddTag(tag)

	task1 := items.Task{
		Item:   items.Item{Id: "task-1", Title: "Task 1", Tag: &tag},
		Status: items.Todo,
	}
	task2 := items.Task{
		Item:   items.Item{Id: "task-2", Title: "Task 2"},
		Status: items.Todo,
	}
	note1 := items.Note{
		Item: items.Item{Id: "note-1", Title: "Note 1", Tag: &tag},
	}

	repo.AddTask(task1)
	repo.AddTask(task2)
	repo.AddNote(note1)

	items, err := repo.GetItemsWithTag("work")
	require.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestFakeRepository_Reset(t *testing.T) {
	repo := NewFakeRepository()

	repo.AddTask(items.Task{Item: items.Item{Id: "task-1"}})
	repo.AddNote(items.Note{Item: items.Item{Id: "note-1"}})
	repo.AddTag(items.Tag{Id: "tag-1", Name: "test"})

	err := repo.Reset()
	require.NoError(t, err)

	assert.Equal(t, 0, repo.TaskCount())
	assert.Equal(t, 0, repo.NoteCount())
	assert.Equal(t, 0, repo.TagCount())
}

func TestFakeRepository_ErrorInjection(t *testing.T) {
	repo := NewFakeRepository()
	expectedErr := errors.New("injected error")

	t.Run("GetTasksError", func(t *testing.T) {
		repo.GetTasksError = expectedErr
		_, err := repo.GetTasks()
		assert.Equal(t, expectedErr, err)
		repo.GetTasksError = nil
	})

	t.Run("GetNotesError", func(t *testing.T) {
		repo.GetNotesError = expectedErr
		_, err := repo.GetNotes()
		assert.Equal(t, expectedErr, err)
		repo.GetNotesError = nil
	})

	t.Run("CreateTaskError", func(t *testing.T) {
		repo.CreateTaskError = expectedErr
		err := repo.CreateTask(&items.Task{Item: items.Item{Id: "t1"}})
		assert.Equal(t, expectedErr, err)
		repo.CreateTaskError = nil
	})

	t.Run("CreateNoteError", func(t *testing.T) {
		repo.CreateNoteError = expectedErr
		err := repo.CreateNote(&items.Note{Item: items.Item{Id: "n1"}})
		assert.Equal(t, expectedErr, err)
		repo.CreateNoteError = nil
	})

	t.Run("UpdateTaskError", func(t *testing.T) {
		repo.AddTask(items.Task{Item: items.Item{Id: "t1"}})
		repo.UpdateTaskError = expectedErr
		err := repo.UpdateTask(items.Task{Item: items.Item{Id: "t1"}})
		assert.Equal(t, expectedErr, err)
		repo.UpdateTaskError = nil
	})

	t.Run("UpdateNoteError", func(t *testing.T) {
		repo.AddNote(items.Note{Item: items.Item{Id: "n1"}})
		repo.UpdateNoteError = expectedErr
		err := repo.UpdateNote(items.Note{Item: items.Item{Id: "n1"}})
		assert.Equal(t, expectedErr, err)
		repo.UpdateNoteError = nil
	})

	t.Run("RemoveTaskError", func(t *testing.T) {
		repo.RemoveTaskError = expectedErr
		err := repo.RemoveTask("t1")
		assert.Equal(t, expectedErr, err)
		repo.RemoveTaskError = nil
	})

	t.Run("RemoveNoteError", func(t *testing.T) {
		repo.RemoveNoteError = expectedErr
		err := repo.RemoveNote("n1")
		assert.Equal(t, expectedErr, err)
		repo.RemoveNoteError = nil
	})

	t.Run("UpdateTaskStatusError", func(t *testing.T) {
		repo.UpdateTaskStatusError = expectedErr
		err := repo.UpdateTaskStatus(items.Task{Item: items.Item{Id: "t1"}}, items.Done)
		assert.Equal(t, expectedErr, err)
		repo.UpdateTaskStatusError = nil
	})

	t.Run("SetTaskTagError", func(t *testing.T) {
		repo.SetTaskTagError = expectedErr
		err := repo.SetTaskTag(items.Task{Item: items.Item{Id: "t1"}}, items.Tag{})
		assert.Equal(t, expectedErr, err)
		repo.SetTaskTagError = nil
	})

	t.Run("UnsetTaskTagError", func(t *testing.T) {
		repo.UnsetTaskTagError = expectedErr
		err := repo.UnsetTaskTag(items.Task{Item: items.Item{Id: "t1"}})
		assert.Equal(t, expectedErr, err)
		repo.UnsetTaskTagError = nil
	})

	t.Run("SetNoteTagError", func(t *testing.T) {
		repo.SetNoteTagError = expectedErr
		err := repo.SetNoteTag(items.Note{Item: items.Item{Id: "n1"}}, items.Tag{})
		assert.Equal(t, expectedErr, err)
		repo.SetNoteTagError = nil
	})

	t.Run("UnsetNoteTagError", func(t *testing.T) {
		repo.UnsetNoteTagError = expectedErr
		err := repo.UnsetNoteTag(items.Note{Item: items.Item{Id: "n1"}})
		assert.Equal(t, expectedErr, err)
		repo.UnsetNoteTagError = nil
	})

	t.Run("GetTagError", func(t *testing.T) {
		repo.GetTagError = expectedErr
		_, err := repo.GetTag("test")
		assert.Equal(t, expectedErr, err)
		repo.GetTagError = nil
	})

	t.Run("GetTagsError", func(t *testing.T) {
		repo.GetTagsError = expectedErr
		_, err := repo.GetTags()
		assert.Equal(t, expectedErr, err)
		repo.GetTagsError = nil
	})

	t.Run("CreateTagError", func(t *testing.T) {
		repo.CreateTagError = expectedErr
		_, err := repo.CreateTag("test")
		assert.Equal(t, expectedErr, err)
		repo.CreateTagError = nil
	})

	t.Run("RemoveTagError", func(t *testing.T) {
		repo.RemoveTagError = expectedErr
		err := repo.RemoveTag("test")
		assert.Equal(t, expectedErr, err)
		repo.RemoveTagError = nil
	})

	t.Run("GetItemsWithTagError", func(t *testing.T) {
		repo.GetItemsWithTagError = expectedErr
		_, err := repo.GetItemsWithTag("test")
		assert.Equal(t, expectedErr, err)
		repo.GetItemsWithTagError = nil
	})

	t.Run("ResetError", func(t *testing.T) {
		repo.ResetError = expectedErr
		err := repo.Reset()
		assert.Equal(t, expectedErr, err)
		repo.ResetError = nil
	})
}

func TestFakeRepository_TestHelperMethods(t *testing.T) {
	repo := NewFakeRepository()

	t.Run("AddTask bypasses CreateTask", func(t *testing.T) {
		task := items.Task{Item: items.Item{Id: "task-1", Title: "Direct"}}
		repo.AddTask(task)
		assert.Equal(t, 1, repo.TaskCount())
		assert.False(t, repo.HasCall("CreateTask"))
	})

	t.Run("AddNote bypasses CreateNote", func(t *testing.T) {
		note := items.Note{Item: items.Item{Id: "note-1", Title: "Direct"}}
		repo.AddNote(note)
		assert.Equal(t, 1, repo.NoteCount())
		assert.False(t, repo.HasCall("CreateNote"))
	})

	t.Run("AddTag bypasses CreateTag", func(t *testing.T) {
		tag := items.Tag{Id: "tag-1", Name: "direct"}
		repo.AddTag(tag)
		assert.Equal(t, 1, repo.TagCount())
		assert.False(t, repo.HasCall("CreateTag"))
	})
}

func TestFakeRepository_ImplementsRepository(t *testing.T) {
	var _ repository.Repository = NewFakeRepository()
}

func TestFakeRepository_MultipleTasksAndNotes(t *testing.T) {
	repo := NewFakeRepository()

	for i := 0; i < 5; i++ {
		repo.AddTask(items.Task{Item: items.Item{Id: string(rune('a' + i))}})
		repo.AddNote(items.Note{Item: items.Item{Id: string(rune('A' + i))}})
	}

	assert.Equal(t, 5, repo.TaskCount())
	assert.Equal(t, 5, repo.NoteCount())

	tasks, _ := repo.GetTasks()
	assert.Len(t, tasks, 5)

	notes, _ := repo.GetNotes()
	assert.Len(t, notes, 5)
}
