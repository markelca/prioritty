package obsidian

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObsidianRepository_CreateNote(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{
		Item: items.Item{Title: "New Note", Body: "Note content"},
	}

	err := repo.CreateNote(note)

	require.NoError(t, err)
	assert.NotEmpty(t, note.Id)
	assert.Contains(t, note.Id, ".md")
}

func TestObsidianRepository_GetNotes(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "note1.md", `---
type: note
title: Note 1
---
Note 1 body`)

	notes, err := repo.GetNotes()

	require.NoError(t, err)
	assert.Len(t, notes, 1)
	assert.Equal(t, "Note 1", notes[0].Title)
}

func TestObsidianRepository_GetNotes_Empty(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	notes, err := repo.GetNotes()

	require.NoError(t, err)
	assert.Empty(t, notes)
}

func TestObsidianRepository_GetNotes_IgnoresTasks(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Task 1
---
`)

	notes, err := repo.GetNotes()

	require.NoError(t, err)
	assert.Empty(t, notes)
}

func TestObsidianRepository_UpdateNote(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{
		Item: items.Item{Title: "Original", Body: "body"},
	}
	require.NoError(t, repo.CreateNote(note))

	note.Title = "Updated"
	note.Body = "New body"
	err := repo.UpdateNote(*note)

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	assert.Equal(t, "Updated", notes[0].Title)
	assert.Equal(t, "New body\n", notes[0].Body)
}

func TestObsidianRepository_RemoveNote(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{Item: items.Item{Title: "Note to remove"}}
	require.NoError(t, repo.CreateNote(note))

	err := repo.RemoveNote(note.Id)

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	assert.Empty(t, notes)
}

func TestObsidianRepository_SetNoteTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{Item: items.Item{Title: "Note"}}
	require.NoError(t, repo.CreateNote(note))

	err := repo.SetNoteTag(*note, items.Tag{Id: "docs", Name: "docs"})

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	require.NotNil(t, notes[0].Tag)
	assert.Equal(t, "docs", notes[0].Tag.Name)
}

func TestObsidianRepository_UnsetNoteTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{Item: items.Item{Title: "Note"}}
	require.NoError(t, repo.CreateNote(note))
	require.NoError(t, repo.SetNoteTag(*note, items.Tag{Id: "docs", Name: "docs"}))

	err := repo.UnsetNoteTag(*note)

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	assert.Nil(t, notes[0].Tag)
}

func TestObsidianRepository_NoteCRUDCycle(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{Item: items.Item{Title: "CRUD Note", Body: "Initial body"}}

	err := repo.CreateNote(note)
	require.NoError(t, err)

	notes, err := repo.GetNotes()
	require.NoError(t, err)
	require.Len(t, notes, 1)
	assert.Equal(t, "CRUD Note", notes[0].Title)

	notes[0].Title = "Updated Note"
	notes[0].Body = "Updated body"
	err = repo.UpdateNote(notes[0])
	require.NoError(t, err)

	notes, _ = repo.GetNotes()
	assert.Equal(t, "Updated Note", notes[0].Title)

	err = repo.RemoveNote(notes[0].Id)
	require.NoError(t, err)

	notes, _ = repo.GetNotes()
	assert.Empty(t, notes)
}
