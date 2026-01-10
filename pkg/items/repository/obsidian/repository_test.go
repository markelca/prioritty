package obsidian

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObsidianRepository(t *testing.T) {
	repo := NewObsidianRepository("/path/to/vault")

	assert.Equal(t, "/path/to/vault", repo.VaultPath())
}

func TestObsidianRepository_VaultPath(t *testing.T) {
	repo := NewObsidianRepository("/test/vault")

	assert.Equal(t, "/test/vault", repo.VaultPath())
}

func TestObsidianRepository_Reset(t *testing.T) {
	vaultDir := setupTestVault(t)

	createTestMarkdownFile(t, vaultDir, "task1.md", "---\ntype: task\ntitle: Task 1\n---\n")
	createTestMarkdownFile(t, vaultDir, "note1.md", "---\ntype: note\ntitle: Note 1\n---\n")
	createTestMarkdownFile(t, vaultDir, "keep.txt", "should not be removed")
	require.NoError(t, os.WriteFile(filepath.Join(vaultDir, ".obsidian", "config.json"), []byte("{}"), 0644))

	repo := NewObsidianRepository(vaultDir)
	err := repo.Reset()

	require.NoError(t, err)

	files, err := os.ReadDir(vaultDir)
	require.NoError(t, err)

	// Markdown files removed; .obsidian and keep.txt remain.
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	assert.Contains(t, names, ".obsidian")
	assert.Contains(t, names, "keep.txt")
	assert.NotContains(t, names, "task1.md")
	assert.NotContains(t, names, "note1.md")
}

func TestObsidianRepository_Reset_EmptyVault(t *testing.T) {
	vaultDir := setupTestVault(t)

	repo := NewObsidianRepository(vaultDir)
	err := repo.Reset()

	require.NoError(t, err)
}
