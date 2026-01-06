package render

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
)

func TestCLI_ImplementsRenderer(t *testing.T) {
	var _ items.Renderer = (*CLI)(nil)
}

func TestCLI_RenderTask(t *testing.T) {
	renderer := CLI{}

	tests := []struct {
		name           string
		task           items.Task
		expectContains []string
	}{
		{
			name: "todo task",
			task: items.Task{
				Item:   items.Item{Title: "Todo Task"},
				Status: items.Todo,
			},
			expectContains: []string{"Todo Task"},
		},
		{
			name: "done task",
			task: items.Task{
				Item:   items.Item{Title: "Done Task"},
				Status: items.Done,
			},
			expectContains: []string{"Done Task"},
		},
		{
			name: "in-progress task",
			task: items.Task{
				Item:   items.Item{Title: "In Progress Task"},
				Status: items.InProgress,
			},
			expectContains: []string{"In Progress Task"},
		},
		{
			name: "cancelled task",
			task: items.Task{
				Item:   items.Item{Title: "Cancelled Task"},
				Status: items.Cancelled,
			},
			expectContains: []string{"Cancelled Task"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := renderer.Render(tc.task)

			for _, expected := range tc.expectContains {
				assert.Contains(t, result, expected)
			}
			// Should end with newline
			assert.True(t, len(result) > 0 && result[len(result)-1] == '\n')
		})
	}
}

func TestCLI_RenderNote(t *testing.T) {
	renderer := CLI{}

	note := items.Note{
		Item: items.Item{
			Title: "Test Note",
		},
	}

	result := renderer.Render(note)

	assert.Contains(t, result, "Test Note")
	assert.True(t, len(result) > 0 && result[len(result)-1] == '\n')
}

func TestCLI_RenderTaskWithBody(t *testing.T) {
	renderer := CLI{}

	task := items.Task{
		Item: items.Item{
			Title: "Task with content",
			Body:  "This is a longer body content",
		},
		Status: items.Todo,
	}

	result := renderer.Render(task)

	// Should contain title
	assert.Contains(t, result, "Task with content")
	// Should have content icon (body > 1 char)
	// The exact icon depends on styles, but result should be longer
	assert.Greater(t, len(result), len("Task with content\n"))
}

func TestCLI_RenderNoteWithBody(t *testing.T) {
	renderer := CLI{}

	note := items.Note{
		Item: items.Item{
			Title: "Note with content",
			Body:  "This is note body",
		},
	}

	result := renderer.Render(note)

	assert.Contains(t, result, "Note with content")
}

func TestCLI_RenderUnknownType(t *testing.T) {
	renderer := CLI{}

	// Create a custom type that implements Renderable but isn't Task or Note
	unknown := &unknownRenderable{}

	result := renderer.Render(unknown)

	assert.Contains(t, result, "Unknown renderable item type")
}

// unknownRenderable is a test type for testing unknown type handling
type unknownRenderable struct{}

func (u *unknownRenderable) Render(r items.Renderer) string {
	return r.Render(u)
}

func TestTaskIcons_AllStatusesCovered(t *testing.T) {
	statuses := []items.Status{
		items.Todo,
		items.InProgress,
		items.Done,
		items.Cancelled,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			icon, exists := taskIcons[status]
			assert.True(t, exists, "Icon should exist for status %s", status)
			assert.NotEmpty(t, icon, "Icon should not be empty for status %s", status)
		})
	}
}

func TestTaskTitleStyle_AllStatusesCovered(t *testing.T) {
	statuses := []items.Status{
		items.Todo,
		items.InProgress,
		items.Done,
		items.Cancelled,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			_, exists := taskTitleStyle[status]
			assert.True(t, exists, "Title style should exist for status %s", status)
		})
	}
}
