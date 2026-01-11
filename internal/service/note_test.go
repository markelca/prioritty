package service

import (
	"errors"
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/markelca/prioritty/pkg/items/repository/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteService_GetNotes(t *testing.T) {
	t.Run("returns all notes", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		note1 := items.Note{Item: items.Item{Id: "1", Title: "Note 1"}}
		note2 := items.Note{Item: items.Item{Id: "2", Title: "Note 2"}}
		mockRepo.AddNote(note1)
		mockRepo.AddNote(note2)
		svc := NewService(mockRepo)

		notes, err := svc.GetNotes()

		require.NoError(t, err)
		assert.Len(t, notes, 2)
		assert.True(t, mockRepo.HasCall("GetNotes"))
	})

	t.Run("returns empty slice when no notes", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		svc := NewService(mockRepo)

		notes, err := svc.GetNotes()

		require.NoError(t, err)
		assert.Empty(t, notes)
	})

	t.Run("propagates error", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		mockRepo.GetNotesError = errors.New("db error")
		svc := NewService(mockRepo)

		_, err := svc.GetNotes()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestNoteService_UpdateNote(t *testing.T) {
	t.Run("updates note", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		note := items.Note{Item: items.Item{Id: "1", Title: "Original"}}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		note.Title = "Updated"
		err := svc.UpdateNote(note)

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("UpdateNote"))
	})

	t.Run("propagates error", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		mockRepo.UpdateNoteError = errors.New("update error")
		note := items.Note{Item: items.Item{Id: "1"}}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		err := svc.UpdateNote(note)

		require.Error(t, err)
	})
}

func TestNoteService_AddNote(t *testing.T) {
	t.Run("creates note with title", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		svc := NewService(mockRepo)

		err := svc.AddNote("New Note")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("CreateNote"))
		assert.Equal(t, 1, mockRepo.NoteCount())
	})

	t.Run("propagates error", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		mockRepo.CreateNoteError = errors.New("create error")
		svc := NewService(mockRepo)

		err := svc.AddNote("Note")

		require.Error(t, err)
	})
}

func TestNoteService_removeNote(t *testing.T) {
	t.Run("removes existing note", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		note := items.Note{Item: items.Item{Id: "1", Title: "Note"}}
		mockRepo.AddNote(note)
		svc := NewService(mockRepo)

		err := svc.removeNote("1")

		require.NoError(t, err)
		assert.True(t, mockRepo.HasCall("RemoveNote"))
		assert.Equal(t, 0, mockRepo.NoteCount())
	})

	t.Run("returns error for non-existent note", func(t *testing.T) {
		mockRepo := testutils.NewMockRepository()
		svc := NewService(mockRepo)

		err := svc.removeNote("non-existent")

		require.Error(t, err)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}
