package sqlite

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestSQLiteRepository_SetNoteTag(t *testing.T) {
	repo := setupTestDB(t)

	note := &items.Note{Item: items.Item{Title: "Note"}}
	require.NoError(t, repo.CreateNote(note))

	tag, err := repo.CreateTag("docs")
	require.NoError(t, err)

	err = repo.SetNoteTag(*note, *tag)

	require.NoError(t, err)

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

	notes, _ := repo.GetNotes()
	assert.Nil(t, notes[0].Tag)
}

func TestSQLiteRepository_NoteCRUDCycle(t *testing.T) {
	repo := setupTestDB(t)

	note := &items.Note{
		Item: items.Item{Title: "CRUD Note", Body: "Initial body"},
	}

	err := repo.CreateNote(note)
	require.NoError(t, err)
	noteId := note.Id

	notes, err := repo.GetNotes()
	require.NoError(t, err)
	require.Len(t, notes, 1)
	assert.Equal(t, "CRUD Note", notes[0].Title)

	notes[0].Title = "Updated CRUD Note"
	notes[0].Body = "Updated body"
	err = repo.UpdateNote(notes[0])
	require.NoError(t, err)

	notes, _ = repo.GetNotes()
	assert.Equal(t, "Updated CRUD Note", notes[0].Title)

	err = repo.RemoveNote(noteId)
	require.NoError(t, err)

	notes, _ = repo.GetNotes()
	assert.Empty(t, notes)
}
