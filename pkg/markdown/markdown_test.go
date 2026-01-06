package markdown

import (
	"strings"
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_ValidFrontmatter(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		expectedFM   Frontmatter
		expectedBody string
	}{
		{
			name: "task with all fields",
			content: `---
title: Test Task
type: task
status: done
tag: work
created_at: 2024-01-15
---
This is the body.
`,
			expectedFM: Frontmatter{
				Title:     "Test Task",
				Type:      "task",
				Status:    "done",
				Tag:       "work",
				CreatedAt: "2024-01-15",
			},
			expectedBody: "This is the body.\n",
		},
		{
			name: "note without status",
			content: `---
title: Test Note
type: note
tag: docs
---
Note body content.
`,
			expectedFM: Frontmatter{
				Title: "Test Note",
				Type:  "note",
				Tag:   "docs",
			},
			expectedBody: "Note body content.\n",
		},
		{
			name: "minimal frontmatter",
			content: `---
title: Minimal
---
Body.
`,
			expectedFM: Frontmatter{
				Title: "Minimal",
			},
			expectedBody: "Body.\n",
		},
		{
			name: "empty body",
			content: `---
title: No Body
type: task
status: todo
---
`,
			expectedFM: Frontmatter{
				Title:  "No Body",
				Type:   "task",
				Status: "todo",
			},
			expectedBody: "",
		},
		{
			name: "multiline body",
			content: `---
title: Multiline
---
Line 1
Line 2
Line 3
`,
			expectedFM: Frontmatter{
				Title: "Multiline",
			},
			expectedBody: "Line 1\nLine 2\nLine 3\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var fm Frontmatter
			body, err := Parse(tc.content, &fm)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedFM, fm)
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}

func TestParse_InvalidFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectedErr string
	}{
		{
			name:        "no frontmatter delimiter",
			content:     "Just plain content",
			expectedErr: "no frontmatter found",
		},
		{
			name:        "unclosed frontmatter",
			content:     "---\ntitle: Test\n",
			expectedErr: "unclosed frontmatter",
		},
		{
			name:        "invalid YAML",
			content:     "---\ntitle: [invalid\n---\n",
			expectedErr: "invalid frontmatter YAML",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var fm Frontmatter
			_, err := Parse(tc.content, &fm)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestSerialize_Task(t *testing.T) {
	input := ItemInput{
		ItemType: items.ItemTypeTask,
		Title:    "Test Task",
		Body:     "Task body content",
		Status:   "in-progress",
		Tag:      "work",
	}

	result, err := Serialize(input)

	require.NoError(t, err)
	assert.Contains(t, result, "---")
	assert.Contains(t, result, "title: Test Task")
	assert.Contains(t, result, "type: task")
	assert.Contains(t, result, "status: in-progress")
	assert.Contains(t, result, "tag: work")
	assert.Contains(t, result, "Task body content")
}

func TestSerialize_Note(t *testing.T) {
	input := ItemInput{
		ItemType: items.ItemTypeNote,
		Title:    "Test Note",
		Body:     "Note body content",
		Tag:      "docs",
	}

	result, err := Serialize(input)

	require.NoError(t, err)
	assert.Contains(t, result, "title: Test Note")
	assert.Contains(t, result, "type: note")
	assert.Contains(t, result, "tag: docs")
	assert.NotContains(t, result, "status:") // Notes don't have status
	assert.Contains(t, result, "Note body content")
}

func TestSerialize_EmptyBody(t *testing.T) {
	input := ItemInput{
		ItemType: items.ItemTypeTask,
		Title:    "No Body Task",
		Body:     "",
		Status:   "todo",
	}

	result, err := Serialize(input)

	require.NoError(t, err)
	assert.Contains(t, result, "title: No Body Task")
	// Should end with closing delimiter
	assert.True(t, strings.HasSuffix(strings.TrimSpace(result), "---"))
}

func TestSerializeForEditor_Task(t *testing.T) {
	input := ItemInput{
		ItemType: items.ItemTypeTask,
		Title:    "Editor Task",
		Body:     "Body content",
		Status:   "done",
		Tag:      "coding",
	}

	result, err := SerializeForEditor(input)

	require.NoError(t, err)
	assert.Contains(t, result, "title: Editor Task")
	assert.Contains(t, result, "type: task")
	assert.Contains(t, result, "status: done")
	assert.Contains(t, result, "tag: coding")
}

func TestSerializeForEditor_TaskDefaultsToTodo(t *testing.T) {
	input := ItemInput{
		ItemType: items.ItemTypeTask,
		Title:    "New Task",
		Status:   "", // Empty status
	}

	result, err := SerializeForEditor(input)

	require.NoError(t, err)
	assert.Contains(t, result, "status: todo") // Should default to todo
}

func TestSerializeForEditor_Note(t *testing.T) {
	input := ItemInput{
		ItemType: items.ItemTypeNote,
		Title:    "Editor Note",
		Body:     "Note content",
		Tag:      "docs",
	}

	result, err := SerializeForEditor(input)

	require.NoError(t, err)
	assert.Contains(t, result, "title: Editor Note")
	assert.Contains(t, result, "type: note")
	assert.Contains(t, result, "tag: docs")
	assert.NotContains(t, result, "status:") // Notes don't have status field
}

func TestRoundTrip_Task(t *testing.T) {
	original := ItemInput{
		ItemType: items.ItemTypeTask,
		Title:    "Round Trip Task",
		Body:     "Body content here",
		Status:   "in-progress",
		Tag:      "testing",
	}

	// Serialize
	serialized, err := Serialize(original)
	require.NoError(t, err)

	// Parse back
	var fm Frontmatter
	body, err := Parse(serialized, &fm)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, original.Title, fm.Title)
	assert.Equal(t, string(original.ItemType), fm.Type)
	assert.Equal(t, original.Status, fm.Status)
	assert.Equal(t, original.Tag, fm.Tag)
	assert.Equal(t, original.Body+"\n", body) // Body gets newline appended
}

func TestRoundTrip_Note(t *testing.T) {
	original := ItemInput{
		ItemType: items.ItemTypeNote,
		Title:    "Round Trip Note",
		Body:     "Note body here",
		Tag:      "documentation",
	}

	// Serialize
	serialized, err := Serialize(original)
	require.NoError(t, err)

	// Parse back
	var fm Frontmatter
	body, err := Parse(serialized, &fm)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, original.Title, fm.Title)
	assert.Equal(t, string(original.ItemType), fm.Type)
	assert.Equal(t, original.Tag, fm.Tag)
	assert.Empty(t, fm.Status) // Notes don't have status
	assert.Equal(t, original.Body+"\n", body)
}

func TestFrontmatter_Serialize(t *testing.T) {
	fm := Frontmatter{
		Title:  "Direct Serialize",
		Type:   "task",
		Status: "done",
	}

	result, err := fm.Serialize("Body text")

	require.NoError(t, err)
	assert.Contains(t, string(result), "title: Direct Serialize")
	assert.Contains(t, string(result), "Body text")
}

func TestSerialize_SpecialCharactersInTitle(t *testing.T) {
	input := ItemInput{
		ItemType: items.ItemTypeTask,
		Title:    "Task with: colon and 'quotes'",
		Status:   "todo",
	}

	result, err := Serialize(input)
	require.NoError(t, err)

	// Parse it back to verify it's valid
	var fm Frontmatter
	_, err = Parse(result, &fm)
	require.NoError(t, err)

	// The title should be parseable (YAML may quote it)
	assert.Equal(t, "Task with: colon and 'quotes'", fm.Title)
}

func TestDelimiterConstant(t *testing.T) {
	assert.Equal(t, "---", Delimiter)
}
