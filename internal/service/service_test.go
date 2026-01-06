package service

import (
	"errors"
	"testing"
	"time"

	"github.com/markelca/prioritty/internal/editor"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	mockRepo := repository.NewMockRepository()

	svc := NewService(mockRepo)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.TaskService)
	assert.NotNil(t, svc.NoteService)
}

func TestService_GetAll(t *testing.T) {
	t.Run("returns tasks and notes combined", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()

		// Add test data
		task := items.Task{
			Item:   items.Item{Id: "task-1", Title: "Test Task", CreatedAt: time.Now()},
			Status: items.Todo,
		}
		note := items.Note{
			Item: items.Item{Id: "note-1", Title: "Test Note", CreatedAt: time.Now()},
		}
		mockRepo.AddTask(task)
		mockRepo.AddNote(note)

		svc := NewService(mockRepo)

		allItems, err := svc.GetAll()

		require.NoError(t, err)
		assert.Len(t, allItems, 2)
		assert.True(t, mockRepo.HasCall("GetNotes"))
		assert.True(t, mockRepo.HasCall("GetTasks"))
	})

	t.Run("returns empty slice when no items", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		allItems, err := svc.GetAll()

		require.NoError(t, err)
		assert.Empty(t, allItems)
	})

	t.Run("returns items sorted correctly", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()

		tag := items.Tag{Id: "tag-1", Name: "work"}
		earlier := time.Now().Add(-1 * time.Hour)
		later := time.Now()

		// Item with tag should come before item without tag
		taskWithTag := items.Task{
			Item:   items.Item{Id: "task-1", Title: "Tagged Task", CreatedAt: earlier, Tag: &tag},
			Status: items.Todo,
		}
		taskWithoutTag := items.Task{
			Item:   items.Item{Id: "task-2", Title: "Untagged Task", CreatedAt: later},
			Status: items.Todo,
		}
		mockRepo.AddTask(taskWithTag)
		mockRepo.AddTask(taskWithoutTag)

		svc := NewService(mockRepo)

		allItems, err := svc.GetAll()

		require.NoError(t, err)
		assert.Len(t, allItems, 2)
		// Tagged item should be first
		assert.Equal(t, "task-1", allItems[0].GetId())
	})

	t.Run("propagates notes error", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		mockRepo.GetNotesError = errors.New("notes error")
		svc := NewService(mockRepo)

		_, err := svc.GetAll()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "notes error")
	})

	t.Run("propagates tasks error", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		mockRepo.GetTasksError = errors.New("tasks error")
		svc := NewService(mockRepo)

		_, err := svc.GetAll()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "tasks error")
	})
}

func TestService_RemoveItem(t *testing.T) {
	t.Run("removes task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{Item: items.Item{Id: "task-1", Title: "Test"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.RemoveItem(&task)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("RemoveTask"))
		assert.Equal(t, 0, mockRepo.TaskCount())
	})

	t.Run("removes note", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		note := items.Note{Item: items.Item{Id: "note-1", Title: "Test"}}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		err := svc.RemoveItem(&note)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("RemoveNote"))
		assert.Equal(t, 0, mockRepo.NoteCount())
	})

	t.Run("returns error for unknown type", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		err := svc.RemoveItem(&unknownItem{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Cannot remove item")
	})
}

func TestService_SetTag(t *testing.T) {
	t.Run("creates new tag and sets on task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{Item: items.Item{Id: "task-1", Title: "Test"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.SetTag(&task, "work")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("GetTag"))
		assert.True(t, mockRepo.HasCall("CreateTag"))
		assert.True(t, mockRepo.HasCall("SetTaskTag"))
	})

	t.Run("uses existing tag", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		existingTag := items.Tag{Id: "tag-1", Name: "work"}
		mockRepo.AddTag(existingTag)
		task := items.Task{Item: items.Item{Id: "task-1", Title: "Test"}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.SetTag(&task, "work")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("GetTag"))
		assert.False(t, mockRepo.HasCall("CreateTag")) // Should not create
		assert.True(t, mockRepo.HasCall("SetTaskTag"))
	})

	t.Run("sets tag on note", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		note := items.Note{Item: items.Item{Id: "note-1", Title: "Test"}}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		err := svc.SetTag(&note, "docs")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("SetNoteTag"))
	})

	t.Run("returns error for unknown type", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		err := svc.SetTag(&unknownItem{}, "tag")

		require.Error(t, err)
	})
}

func TestService_UnsetTag(t *testing.T) {
	t.Run("unsets tag from task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		tag := items.Tag{Id: "tag-1", Name: "work"}
		task := items.Task{Item: items.Item{Id: "task-1", Title: "Test", Tag: &tag}, Status: items.Todo}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.UnsetTag(&task)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("UnsetTaskTag"))
	})

	t.Run("unsets tag from note", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		tag := items.Tag{Id: "tag-1", Name: "docs"}
		note := items.Note{Item: items.Item{Id: "note-1", Title: "Test", Tag: &tag}}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		err := svc.UnsetTag(&note)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("UnsetNoteTag"))
	})
}

func TestService_GetTags(t *testing.T) {
	mockRepo := repository.NewMockRepository()
	tag1 := items.Tag{Id: "1", Name: "work"}
	tag2 := items.Tag{Id: "2", Name: "docs"}
	mockRepo.AddTag(tag1)
	mockRepo.AddTag(tag2)
	svc := NewService(mockRepo)

	tags, err := svc.GetTags()

	require.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.True(t, mockRepo.HasCall("GetTags"))
}

func TestService_RemoveTag(t *testing.T) {
	t.Run("removes tag with no items", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		tag := items.Tag{Id: "1", Name: "unused"}
		mockRepo.AddTag(tag)
		svc := NewService(mockRepo)

		err := svc.RemoveTag("unused")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("RemoveTag"))
	})

	t.Run("prevents removal when tag is in use", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		tag := items.Tag{Id: "1", Name: "work"}
		mockRepo.AddTag(tag)
		task := items.Task{
			Item:   items.Item{Id: "task-1", Title: "Test", Tag: &tag},
			Status: items.Todo,
		}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		err := svc.RemoveTag("work")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot remove tag")
		assert.Contains(t, err.Error(), "1 item(s)")
	})
}

func TestService_GetItemsWithTag(t *testing.T) {
	mockRepo := repository.NewMockRepository()
	tag := items.Tag{Id: "1", Name: "work"}
	mockRepo.AddTag(tag)
	task := items.Task{
		Item:   items.Item{Id: "task-1", Title: "Work Task", Tag: &tag},
		Status: items.Todo,
	}
	mockRepo.AddTask(task)
	svc := NewService(mockRepo)

	items, err := svc.GetItemsWithTag("work")

	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.True(t, mockRepo.HasCall("GetItemsWithTag"))
}

func TestService_CreateTaskFromEditorMsg(t *testing.T) {
	t.Run("creates task without tag", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			Title:  "New Task",
			Body:   "Task body",
			Status: "todo",
		}

		err := svc.CreateTaskFromEditorMsg(msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("CreateTask"))
		assert.Equal(t, 1, mockRepo.TaskCount())
	})

	t.Run("creates task with tag", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			Title:  "Tagged Task",
			Body:   "Body",
			Status: "in-progress",
			Tag:    "work",
		}

		err := svc.CreateTaskFromEditorMsg(msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("CreateTask"))
		assert.True(t, mockRepo.HasCall("SetTaskTag"))
	})

	t.Run("propagates create error", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		mockRepo.CreateTaskError = errors.New("create failed")
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{Title: "Task"}

		err := svc.CreateTaskFromEditorMsg(msg)

		require.Error(t, err)
	})
}

func TestService_CreateNoteFromEditorMsg(t *testing.T) {
	t.Run("creates note without tag", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			Title: "New Note",
			Body:  "Note body",
		}

		err := svc.CreateNoteFromEditorMsg(msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("CreateNote"))
		assert.Equal(t, 1, mockRepo.NoteCount())
	})

	t.Run("creates note with tag", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			Title: "Tagged Note",
			Body:  "Body",
			Tag:   "docs",
		}

		err := svc.CreateNoteFromEditorMsg(msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("CreateNote"))
		assert.True(t, mockRepo.HasCall("SetNoteTag"))
	})
}

func TestService_UpdateItemFromEditorMsg(t *testing.T) {
	t.Run("updates task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{
			Item:   items.Item{Id: "task-1", Title: "Old Title"},
			Status: items.Todo,
		}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			ItemType: items.ItemTypeTask,
			Title:    "New Title",
			Body:     "New Body",
			Status:   "done",
		}

		err := svc.UpdateItemFromEditorMsg(&task, msg)

		require.NoError(t, err)
		assert.Equal(t, "New Title", task.Title)
		assert.Equal(t, "New Body", task.Body)
		assert.Equal(t, items.Done, task.Status)
	})

	t.Run("updates note", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		note := items.Note{
			Item: items.Item{Id: "note-1", Title: "Old Title"},
		}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			ItemType: items.ItemTypeNote,
			Title:    "New Title",
			Body:     "New Body",
		}

		err := svc.UpdateItemFromEditorMsg(&note, msg)

		require.NoError(t, err)
		assert.Equal(t, "New Title", note.Title)
		assert.Equal(t, "New Body", note.Body)
	})

	t.Run("converts task to note", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{
			Item:   items.Item{Id: "task-1", Title: "Task"},
			Status: items.Todo,
		}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			Id:       "task-1", // Original ID
			ItemType: items.ItemTypeNote,
			Title:    "Now a Note",
			Body:     "Converted",
		}

		err := svc.UpdateItemFromEditorMsg(&task, msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("RemoveTask"))
		assert.True(t, mockRepo.HasCall("CreateNote"))
	})

	t.Run("converts note to task", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		note := items.Note{
			Item: items.Item{Id: "note-1", Title: "Note"},
		}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			Id:       "note-1", // Original ID
			ItemType: items.ItemTypeTask,
			Title:    "Now a Task",
			Status:   "todo",
		}

		err := svc.UpdateItemFromEditorMsg(&note, msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("RemoveNote"))
		assert.True(t, mockRepo.HasCall("CreateTask"))
	})

	t.Run("updates tag when changed", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		task := items.Task{
			Item:   items.Item{Id: "task-1", Title: "Task"},
			Status: items.Todo,
		}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			ItemType: items.ItemTypeTask,
			Title:    "Task",
			Tag:      "newTag",
		}

		err := svc.UpdateItemFromEditorMsg(&task, msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("SetTaskTag"))
	})

	t.Run("removes tag when empty", func(t *testing.T) {
		mockRepo := repository.NewMockRepository()
		tag := items.Tag{Id: "1", Name: "work"}
		task := items.Task{
			Item:   items.Item{Id: "task-1", Title: "Task", Tag: &tag},
			Status: items.Todo,
		}
		mockRepo.AddTask(task)
		svc := NewService(mockRepo)

		msg := editor.EditorFinishedMsg{
			ItemType: items.ItemTypeTask,
			Title:    "Task",
			Tag:      "", // Empty to remove
		}

		err := svc.UpdateItemFromEditorMsg(&task, msg)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("UnsetTaskTag"))
	})
}

// unknownItem is a test type that implements ItemInterface but isn't Task or Note
type unknownItem struct{}

func (u *unknownItem) GetId() string               { return "unknown" }
func (u *unknownItem) GetTitle() string            { return "Unknown" }
func (u *unknownItem) GetBody() string             { return "" }
func (u *unknownItem) GetTag() *items.Tag          { return nil }
func (u *unknownItem) GetCreatedAt() time.Time     { return time.Time{} }
func (u *unknownItem) After(items.ItemInterface) bool { return false }
func (u *unknownItem) Render(r items.Renderer) string { return r.Render(u) }
