package obsidian

import (
	"testing"
	"time"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/markdown"
	"github.com/stretchr/testify/assert"
)

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
