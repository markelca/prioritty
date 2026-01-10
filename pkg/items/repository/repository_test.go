package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandTilde(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde path expands to home",
			input:    "~/documents/test.txt",
			expected: home + "/documents/test.txt",
		},
		{
			name:     "tilde at end",
			input:    "~/",
			expected: home,
		},
		{
			name:     "absolute path unchanged",
			input:    "/absolute/path/file.txt",
			expected: "/absolute/path/file.txt",
		},
		{
			name:     "relative path unchanged",
			input:    "relative/path/file.txt",
			expected: "relative/path/file.txt",
		},
		{
			name:     "empty string unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "tilde not at start unchanged",
			input:    "path/to/~/file.txt",
			expected: "path/to/~/file.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := expandTilde(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExpandTilde_ActualHomeDir(t *testing.T) {
	result := expandTilde("~/test")
	home, _ := os.UserHomeDir()
	assert.Equal(t, home+"/test", result)
}

func TestGetDatabasePath_SQLite(t *testing.T) {
	t.Run("demo returns temp path", func(t *testing.T) {
		dbPath, err := GetDatabasePath(RepoTypeSQLite, true)

		assert.NoError(t, err)
		assert.Contains(t, dbPath, "prioritty_demo.db")
		assert.Contains(t, dbPath, "/tmp")
	})
}

func TestGetDatabasePath_Obsidian(t *testing.T) {
	t.Run("demo returns temp path", func(t *testing.T) {
		dbPath, err := GetDatabasePath(RepoTypeObsidian, true)

		assert.NoError(t, err)
		assert.Contains(t, dbPath, "prioritty_demo_vault")
		assert.Contains(t, dbPath, "/tmp")
	})
}

func TestGetDatabasePath_Unsupported(t *testing.T) {
	t.Run("unsupported repo type returns error", func(t *testing.T) {
		_, err := GetDatabasePath("unsupported", false)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "repository type not supported")
	})
}

func TestRepoTypeConstants(t *testing.T) {
	assert.Equal(t, "obsidian", RepoTypeObsidian)
	assert.Equal(t, "sqlite", RepoTypeSQLite)
}

func TestErrNotFound(t *testing.T) {
	assert.Error(t, ErrNotFound)
	assert.Equal(t, "not found", ErrNotFound.Error())
}
