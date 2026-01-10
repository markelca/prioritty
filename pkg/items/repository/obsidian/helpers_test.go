package obsidian

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestVault(t *testing.T) string {
	vaultDir, err := os.MkdirTemp("", "prioritty_obsidian_test_*")
	require.NoError(t, err)

	err = os.Mkdir(filepath.Join(vaultDir, ".obsidian"), 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(vaultDir)
	})

	return vaultDir
}

func createTestMarkdownFile(t *testing.T, vaultDir, filename, content string) {
	filePath := filepath.Join(vaultDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple text", "Hello World", "hello-world"},
		{"camelCase", "helloWorld", "helloworld"},
		{"snake_case", "hello_world", "hello-world"},
		{"spaces", "multiple   spaces", "multiple-spaces"},
		{"special chars", "Hello/World*Test", "helloworldtest"},
		{"numbers", "Task 123", "task-123"},
		{"empty", "", "untitled"},
		{"only special", "!!!", "untitled"},
		{"unicode", "Tarea Española", "tarea-espanola"},
		{"leading trailing spaces", "  Title  ", "title"},
		{"multiple hyphens", "test---title", "test-title"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := toKebabCase(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilenameFromTitle(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{"simple title", "My Task", "my-task.md"},
		{"with spaces", "Task One Two", "task-one-two.md"},
		{"empty", "", "untitled.md"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := filenameFromTitle(tc.title)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilenameFromTitle_Exported(t *testing.T) {
	result := FilenameFromTitle("Test Title")
	assert.Equal(t, "test-title.md", result)
}

func TestRelativeID(t *testing.T) {
	vault := "/vault"
	fullPath := "/vault/subdir/file.md"

	result := relativeID(vault, fullPath)

	assert.Equal(t, filepath.Join("subdir", "file.md"), result)
}

func TestRelativeID_Error(t *testing.T) {
	vault := "/vault"
	fullPath := "/other/file.md"

	result := relativeID(vault, fullPath)

	assert.Contains(t, result, "file.md")
}

func TestFullPathFromID(t *testing.T) {
	vault := "/vault"
	id := "subdir/file.md"

	result := fullPathFromID(vault, id)

	assert.Equal(t, filepath.Join("/vault", "subdir", "file.md"), result)
}

func TestScanMarkdownFiles(t *testing.T) {
	vaultDir := setupTestVault(t)

	createTestMarkdownFile(t, vaultDir, "task1.md", "---\ntype: task\ntitle: Task 1\n---\n")
	createTestMarkdownFile(t, vaultDir, "note1.md", "---\ntype: note\ntitle: Note 1\n---\n")
	createTestMarkdownFile(t, vaultDir, "readme.txt", "not markdown")
	require.NoError(t, os.Mkdir(filepath.Join(vaultDir, "subdir"), 0755))
	createTestMarkdownFile(t, vaultDir, "subdir/task2.md", "---\ntype: task\ntitle: Task 2\n---\n")

	files, err := scanMarkdownFiles(vaultDir)

	require.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Contains(t, files, filepath.Join(vaultDir, "task1.md"))
	assert.Contains(t, files, filepath.Join(vaultDir, "note1.md"))
}

func TestScanMarkdownFiles_EmptyVault(t *testing.T) {
	vaultDir := setupTestVault(t)

	files, err := scanMarkdownFiles(vaultDir)

	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestScanMarkdownFiles_NonExistentVault(t *testing.T) {
	_, err := scanMarkdownFiles("/nonexistent/path")
	assert.Error(t, err)
}

func TestUniqueFilename(t *testing.T) {
	vaultDir := setupTestVault(t)

	path1 := uniqueFilename(vaultDir, "My Task")
	assert.Equal(t, "my-task.md", filepath.Base(path1))

	require.NoError(t, os.WriteFile(path1, []byte("content"), 0644))

	path2 := uniqueFilename(vaultDir, "My Task")
	assert.Equal(t, "my-task-2.md", filepath.Base(path2))
	assert.NotEqual(t, path1, path2)
}

func TestUniqueFilename_ConflictWithCounter(t *testing.T) {
	vaultDir := setupTestVault(t)

	createTestMarkdownFile(t, vaultDir, "my-task.md", "")
	createTestMarkdownFile(t, vaultDir, "my-task-2.md", "")

	path := uniqueFilename(vaultDir, "My Task")

	assert.Contains(t, path, "my-task-3.md")
}
