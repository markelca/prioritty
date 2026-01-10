package sqlite

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusIdToStatus(t *testing.T) {
	tests := []struct {
		id       int
		expected items.Status
	}{
		{0, items.Todo},
		{1, items.InProgress},
		{2, items.Done},
		{3, items.Cancelled},
		{99, items.Todo},
		{-1, items.Todo},
	}

	for _, tc := range tests {
		result := statusIdToStatus(tc.id)
		assert.Equal(t, tc.expected, result, "statusIdToStatus(%d) should return %s", tc.id, tc.expected)
	}
}

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

func TestSQLiteRepository_GetItemsWithTag(t *testing.T) {
	repo := setupTestDB(t)

	tag, _ := repo.CreateTag("work")

	task := &items.Task{Item: items.Item{Title: "Work Task"}, Status: items.Todo}
	repo.CreateTask(task)
	repo.SetTaskTag(*task, *tag)

	note := &items.Note{Item: items.Item{Title: "Work Note"}}
	repo.CreateNote(note)
	repo.SetNoteTag(*note, *tag)

	otherTask := &items.Task{Item: items.Item{Title: "Other Task"}, Status: items.Todo}
	repo.CreateTask(otherTask)

	got, err := repo.GetItemsWithTag("work")

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestSQLiteRepository_GetItemsWithTag_Empty(t *testing.T) {
	repo := setupTestDB(t)

	repo.CreateTag("unused")

	got, err := repo.GetItemsWithTag("unused")

	require.NoError(t, err)
	assert.Empty(t, got)
}
