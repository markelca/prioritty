package obsidian

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObsidianRepository_GetTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Task 1
tag: work
---
`)

	tag, err := repo.GetTag("work")

	require.NoError(t, err)
	assert.Equal(t, "work", tag.Name)
}

func TestObsidianRepository_GetTag_NotFound(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	_, err := repo.GetTag("nonexistent")

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestObsidianRepository_GetTags(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Task 1
tag: work
---
`)
	createTestMarkdownFile(t, vaultDir, "note1.md", `---
type: note
title: Note 1
tag: work
---
`)
	createTestMarkdownFile(t, vaultDir, "task2.md", `---
type: task
title: Task 2
tag: urgent
---
`)

	tags, err := repo.GetTags()

	require.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "urgent", tags[0].Name)
	assert.Equal(t, "work", tags[1].Name)
}

func TestObsidianRepository_GetTags_Empty(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	tags, err := repo.GetTags()

	require.NoError(t, err)
	assert.Empty(t, tags)
}

func TestObsidianRepository_CreateTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	tag, err := repo.CreateTag("newtag")

	require.NoError(t, err)
	assert.Equal(t, "newtag", tag.Name)
}

func TestObsidianRepository_RemoveTag_Exists(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Task 1
tag: work
---
`)

	err := repo.RemoveTag("work")

	assert.NoError(t, err)
}

func TestObsidianRepository_RemoveTag_NotFound(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	err := repo.RemoveTag("nonexistent")

	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestObsidianRepository_GetItemsWithTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Work Task
tag: work
---
`)
	createTestMarkdownFile(t, vaultDir, "note1.md", `---
type: note
title: Work Note
tag: work
---
`)
	createTestMarkdownFile(t, vaultDir, "task2.md", `---
type: task
title: Other Task
tag: other
---
`)

	got, err := repo.GetItemsWithTag("work")

	require.NoError(t, err)
	require.Len(t, got, 2)

	var hasTask, hasNote bool
	for _, it := range got {
		switch v := it.(type) {
		case *items.Task:
			hasTask = v.Title == "Work Task"
		case *items.Note:
			hasNote = v.Title == "Work Note"
		}
	}
	assert.True(t, hasTask)
	assert.True(t, hasNote)
}

func TestObsidianRepository_GetItemsWithTag_Empty(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Task
tag: used
---
`)

	got, err := repo.GetItemsWithTag("unused")

	require.NoError(t, err)
	assert.Empty(t, got)
}
