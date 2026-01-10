package obsidian

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/markdown"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestVault(t *testing.T) string {
	vaultDir, err := os.MkdirTemp("", "prioritty_obsidian_test_*")
	require.NoError(t, err)

	err = os.Mkdir(filepath.Join(vaultDir, ".obsidian"), 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(vaultDir)
	})

	return vaultDir
}

func createTestMarkdownFile(t *testing.T, vaultDir, filename, content string) {
	filePath := filepath.Join(vaultDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)
}

// --- Helper Tests ---

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
	os.Mkdir(filepath.Join(vaultDir, "subdir"), 0755)
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

	os.WriteFile(path1, []byte("content"), 0644)

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

// --- Repository Setup Tests ---

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
	os.WriteFile(filepath.Join(vaultDir, ".obsidian", "config.json"), []byte("{}"), 0644)

	repo := NewObsidianRepository(vaultDir)
	err := repo.Reset()

	require.NoError(t, err)

	files, _ := os.ReadDir(vaultDir)
	assert.Len(t, files, 2)
	assert.Contains(t, files[0].Name(), ".obsidian")
	assert.Contains(t, files[1].Name(), "keep.txt")
}

func TestObsidianRepository_Reset_EmptyVault(t *testing.T) {
	vaultDir := setupTestVault(t)

	repo := NewObsidianRepository(vaultDir)
	err := repo.Reset()

	require.NoError(t, err)
}

// --- Markdown Parsing Tests ---

func TestParseCreatedAt(t *testing.T) {
	t.Run("valid RFC3339", func(t *testing.T) {
		result := parseCreatedAt("2024-01-15T10:30:00Z")
		assert.Equal(t, time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), result)
	})

	t.Run("empty string returns current time", func(t *testing.T) {
		before := time.Now()
		result := parseCreatedAt("")
		after := time.Now()
		assert.True(t, result.After(before.Add(-time.Second)))
		assert.True(t, result.Before(after.Add(time.Second)))
	})

	t.Run("invalid format returns current time", func(t *testing.T) {
		before := time.Now()
		result := parseCreatedAt("invalid")
		after := time.Now()
		assert.True(t, result.After(before.Add(-time.Second)))
		assert.True(t, result.Before(after.Add(time.Second)))
	})
}

func TestFormatCreatedAt(t *testing.T) {
	t.Run("with time", func(t *testing.T) {
		result := formatCreatedAt(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))
		assert.Equal(t, "2024-01-15T10:30:00Z", result)
	})

	t.Run("zero time returns current time", func(t *testing.T) {
		before := time.Now()
		result := formatCreatedAt(time.Time{})
		after := time.Now()
		parsed, _ := time.Parse(time.RFC3339, result)
		assert.True(t, parsed.After(before.Add(-time.Second)))
		assert.True(t, parsed.Before(after.Add(time.Second)))
	})
}

func TestTaskFromFrontmatter(t *testing.T) {
	fm := markdown.Frontmatter{
		Type:      "task",
		Title:     "Test Task",
		Status:    "todo",
		Tag:       "work",
		CreatedAt: "2024-01-15T10:30:00Z",
	}

	task := taskFromFrontmatter(fm, "Task body", "test.md")

	assert.Equal(t, "test.md", task.Id)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "Task body", task.Body)
	assert.Equal(t, items.Todo, task.Status)
	assert.NotNil(t, task.Tag)
	assert.Equal(t, "work", task.Tag.Name)
}

func TestTaskFromFrontmatter_NoTag(t *testing.T) {
	fm := markdown.Frontmatter{
		Type:      "task",
		Title:     "Test Task",
		Status:    "done",
		CreatedAt: "",
	}

	task := taskFromFrontmatter(fm, "", "test.md")

	assert.Nil(t, task.Tag)
}

func TestNoteFromFrontmatter(t *testing.T) {
	fm := markdown.Frontmatter{
		Type:      "note",
		Title:     "Test Note",
		Tag:       "ideas",
		CreatedAt: "2024-01-15T10:30:00Z",
	}

	note := noteFromFrontmatter(fm, "Note body", "test.md")

	assert.Equal(t, "test.md", note.Id)
	assert.Equal(t, "Test Note", note.Title)
	assert.Equal(t, "Note body", note.Body)
	assert.NotNil(t, note.Tag)
	assert.Equal(t, "ideas", note.Tag.Name)
}

func TestNoteFromFrontmatter_NoTag(t *testing.T) {
	fm := markdown.Frontmatter{
		Type:      "note",
		Title:     "Test Note",
		CreatedAt: "",
	}

	note := noteFromFrontmatter(fm, "", "test.md")

	assert.Nil(t, note.Tag)
}

func TestItemInputFromTask(t *testing.T) {
	task := items.Task{
		Item:   items.Item{Title: "Task", Body: "body", CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		Status: items.InProgress,
	}
	task.Tag = &items.Tag{Id: "work", Name: "work"}

	input := itemInputFromTask(task)

	assert.Equal(t, items.ItemTypeTask, input.ItemType)
	assert.Equal(t, "Task", input.Title)
	assert.Equal(t, "body", input.Body)
	assert.Equal(t, "in-progress", input.Status)
	assert.Equal(t, "work", input.Tag)
}

func TestItemInputFromTask_NoTag(t *testing.T) {
	task := items.Task{
		Item:   items.Item{Title: "Task"},
		Status: items.Todo,
	}

	input := itemInputFromTask(task)

	assert.Equal(t, "", input.Tag)
}

func TestItemInputFromNote(t *testing.T) {
	note := items.Note{
		Item: items.Item{Title: "Note", Body: "body", CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
	}
	note.Tag = &items.Tag{Id: "docs", Name: "docs"}

	input := itemInputFromNote(note)

	assert.Equal(t, items.ItemTypeNote, input.ItemType)
	assert.Equal(t, "Note", input.Title)
	assert.Equal(t, "body", input.Body)
	assert.Equal(t, "docs", input.Tag)
}

func TestItemInputFromNote_NoTag(t *testing.T) {
	note := items.Note{
		Item: items.Item{Title: "Note"},
	}

	input := itemInputFromNote(note)

	assert.Equal(t, "", input.Tag)
}

// --- Task CRUD Tests ---

func TestObsidianRepository_CreateTask(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "New Task", Body: "Task content"},
		Status: items.Todo,
	}

	err := repo.CreateTask(task)

	require.NoError(t, err)
	assert.NotEmpty(t, task.Id)
	assert.Contains(t, task.Id, ".md")
}

func TestObsidianRepository_GetTasks(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "task1.md", `---
type: task
title: Task 1
status: todo
---
Task 1 body`)

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "Task 1", tasks[0].Title)
}

func TestObsidianRepository_GetTasks_Empty(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestObsidianRepository_GetTasks_IgnoresNotes(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	createTestMarkdownFile(t, vaultDir, "note1.md", `---
type: note
title: Note 1
---
`)

	tasks, err := repo.GetTasks()

	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestObsidianRepository_UpdateTask(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Original", Body: "body"},
		Status: items.Todo,
	}
	repo.CreateTask(task)

	task.Title = "Updated"
	task.Body = "New body"
	task.Status = items.Done
	err := repo.UpdateTask(*task)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Equal(t, "Updated", tasks[0].Title)
	assert.Equal(t, "New body\n", tasks[0].Body)
	assert.Equal(t, items.Done, tasks[0].Status)
}

func TestObsidianRepository_RemoveTask(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task to remove"},
		Status: items.Todo,
	}
	repo.CreateTask(task)

	err := repo.RemoveTask(task.Id)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Empty(t, tasks)
}

func TestObsidianRepository_UpdateTaskStatus(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task"},
		Status: items.Todo,
	}
	repo.CreateTask(task)

	err := repo.UpdateTaskStatus(*task, items.InProgress)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Equal(t, items.InProgress, tasks[0].Status)
}

func TestObsidianRepository_SetTaskTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task"},
		Status: items.Todo,
	}
	repo.CreateTask(task)

	err := repo.SetTaskTag(*task, items.Tag{Id: "work", Name: "work"})

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.NotNil(t, tasks[0].Tag)
	assert.Equal(t, "work", tasks[0].Tag.Name)
}

func TestObsidianRepository_UnsetTaskTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "Task"},
		Status: items.Todo,
	}
	repo.CreateTask(task)
	repo.SetTaskTag(*task, items.Tag{Id: "work", Name: "work"})

	err := repo.UnsetTaskTag(*task)

	require.NoError(t, err)

	tasks, _ := repo.GetTasks()
	assert.Nil(t, tasks[0].Tag)
}

// --- Note CRUD Tests ---

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
	repo.CreateNote(note)

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

	note := &items.Note{
		Item: items.Item{Title: "Note to remove"},
	}
	repo.CreateNote(note)

	err := repo.RemoveNote(note.Id)

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	assert.Empty(t, notes)
}

func TestObsidianRepository_SetNoteTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{
		Item: items.Item{Title: "Note"},
	}
	repo.CreateNote(note)

	err := repo.SetNoteTag(*note, items.Tag{Id: "docs", Name: "docs"})

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	assert.NotNil(t, notes[0].Tag)
	assert.Equal(t, "docs", notes[0].Tag.Name)
}

func TestObsidianRepository_UnsetNoteTag(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{
		Item: items.Item{Title: "Note"},
	}
	repo.CreateNote(note)
	repo.SetNoteTag(*note, items.Tag{Id: "docs", Name: "docs"})

	err := repo.UnsetNoteTag(*note)

	require.NoError(t, err)

	notes, _ := repo.GetNotes()
	assert.Nil(t, notes[0].Tag)
}

// --- Tag Tests ---

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

	assert.Error(t, err)
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

	assert.Error(t, err)
}

// --- GetItemsWithTag Tests ---

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

	items, err := repo.GetItemsWithTag("work")

	require.NoError(t, err)
	assert.Len(t, items, 2)
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

	items, err := repo.GetItemsWithTag("unused")

	require.NoError(t, err)
	assert.Empty(t, items)
}

// --- Full CRUD Cycle Tests ---

func TestObsidianRepository_TaskCRUDCycle(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	task := &items.Task{
		Item:   items.Item{Title: "CRUD Task", Body: "Initial body"},
		Status: items.Todo,
	}

	err := repo.CreateTask(task)
	require.NoError(t, err)

	tasks, err := repo.GetTasks()
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	assert.Equal(t, "CRUD Task", tasks[0].Title)

	tasks[0].Title = "Updated Task"
	tasks[0].Body = "Updated body"
	err = repo.UpdateTask(tasks[0])
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, "Updated Task", tasks[0].Title)

	err = repo.UpdateTaskStatus(tasks[0], items.Done)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Equal(t, items.Done, tasks[0].Status)

	err = repo.RemoveTask(tasks[0].Id)
	require.NoError(t, err)

	tasks, _ = repo.GetTasks()
	assert.Empty(t, tasks)
}

func TestObsidianRepository_NoteCRUDCycle(t *testing.T) {
	vaultDir := setupTestVault(t)
	repo := NewObsidianRepository(vaultDir)

	note := &items.Note{
		Item: items.Item{Title: "CRUD Note", Body: "Initial body"},
	}

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
