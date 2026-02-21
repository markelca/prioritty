package editor

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEditorContent_ValidContent(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		itemType       items.ItemType
		expectedResult EditorFinishedMsg
	}{
		{
			name: "task with all fields",
			content: `---
title: Test Task
type: task
status: done
tag: work
---
This is the body content.
`,
			itemType: items.ItemTypeTask,
			expectedResult: EditorFinishedMsg{
				ItemType: items.ItemTypeTask,
				Title:    "Test Task",
				Body:     "This is the body content.",
				Status:   "done",
				Tag:      "work",
			},
		},
		{
			name: "note without status",
			content: `---
title: Test Note
type: note
tag: docs
---
Note body here.
`,
			itemType: items.ItemTypeNote,
			expectedResult: EditorFinishedMsg{
				ItemType: items.ItemTypeNote,
				Title:    "Test Note",
				Body:     "Note body here.",
				Tag:      "docs",
			},
		},
		{
			name: "task with minimal fields",
			content: `---
title: Minimal Task
---
`,
			itemType: items.ItemTypeTask,
			expectedResult: EditorFinishedMsg{
				ItemType: items.ItemTypeTask,
				Title:    "Minimal Task",
				Body:     "",
			},
		},
		{
			name: "type override - note to task",
			content: `---
title: Changed Type
type: task
status: todo
---
`,
			itemType: items.ItemTypeNote,
			expectedResult: EditorFinishedMsg{
				ItemType: items.ItemTypeTask,
				Title:    "Changed Type",
				Body:     "",
				Status:   "todo",
			},
		},
		{
			name: "type fallback with invalid type",
			content: `---
title: Invalid Type
type: invalid
---
`,
			itemType: items.ItemTypeTask,
			expectedResult: EditorFinishedMsg{
				ItemType: items.ItemTypeTask,
				Title:    "Invalid Type",
				Body:     "",
			},
		},
		{
			name: "body with multiple lines",
			content: `---
title: Multi-line Body
---
Line 1
Line 2
Line 3
`,
			itemType: items.ItemTypeTask,
			expectedResult: EditorFinishedMsg{
				ItemType: items.ItemTypeTask,
				Title:    "Multi-line Body",
				Body:     "Line 1\nLine 2\nLine 3",
			},
		},
		{
			name: "title with leading/trailing whitespace",
			content: `---
title:   Trimmed Title   
---
`,
			itemType: items.ItemTypeTask,
			expectedResult: EditorFinishedMsg{
				ItemType: items.ItemTypeTask,
				Title:    "Trimmed Title",
				Body:     "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseEditorContent(tc.content, tc.itemType)
			require.NoError(t, result.Err)
			assert.Equal(t, tc.expectedResult.ItemType, result.ItemType)
			assert.Equal(t, tc.expectedResult.Title, result.Title)
			assert.Equal(t, tc.expectedResult.Body, result.Body)
			assert.Equal(t, tc.expectedResult.Status, result.Status)
			assert.Equal(t, tc.expectedResult.Tag, result.Tag)
		})
	}
}

func TestParseEditorContent_Errors(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		itemType    items.ItemType
		expectedErr string
	}{
		{
			name:        "empty content",
			content:     "",
			itemType:    items.ItemTypeTask,
			expectedErr: "no content provided",
		},
		{
			name:        "whitespace only content",
			content:     "   \n\t  \n  ",
			itemType:    items.ItemTypeTask,
			expectedErr: "no content provided",
		},
		{
			name: "missing title",
			content: `---
type: task
status: todo
---
Body content.
`,
			itemType:    items.ItemTypeTask,
			expectedErr: "no title provided",
		},
		{
			name: "empty title",
			content: `---
title: ""
---
Body content.
`,
			itemType:    items.ItemTypeTask,
			expectedErr: "no title provided",
		},
		{
			name: "title with only whitespace",
			content: `---
title: "   "
---
Body content.
`,
			itemType:    items.ItemTypeTask,
			expectedErr: "no title provided",
		},
		{
			name: "invalid YAML frontmatter",
			content: `---
title: [invalid
---
Body.
`,
			itemType:    items.ItemTypeTask,
			expectedErr: "invalid frontmatter",
		},
		{
			name:        "no frontmatter delimiter",
			content:     "Just plain content without frontmatter",
			itemType:    items.ItemTypeTask,
			expectedErr: "invalid frontmatter",
		},
		{
			name: "unclosed frontmatter",
			content: `---
title: Test
`,
			itemType:    items.ItemTypeTask,
			expectedErr: "invalid frontmatter",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := parseEditorContent(tc.content, tc.itemType)
			require.Error(t, result.Err)
			assert.Contains(t, result.Err.Error(), tc.expectedErr)
		})
	}
}

func TestParseEditorContent_EdgeCases(t *testing.T) {
	t.Run("status field ignored for notes", func(t *testing.T) {
		content := `---
title: Note with Status
type: note
status: done
tag: personal
---
Body.
`
		result := parseEditorContent(content, items.ItemTypeNote)
		require.NoError(t, result.Err)
		assert.Equal(t, items.ItemTypeNote, result.ItemType)
		assert.Equal(t, "done", result.Status)
	})

	t.Run("body is trimmed", func(t *testing.T) {
		content := `---
title: Test
---
  
  Body with surrounding whitespace.  
  
`
		result := parseEditorContent(content, items.ItemTypeTask)
		require.NoError(t, result.Err)
		assert.Equal(t, "Body with surrounding whitespace.", result.Body)
	})
}

func TestEditorInput_Struct(t *testing.T) {
	input := EditorInput{
		Id:       "test-id",
		ItemType: items.ItemTypeTask,
		Title:    "Test Title",
		Body:     "Test Body",
		Status:   "in-progress",
		Tag:      "work",
	}

	assert.Equal(t, "test-id", input.Id)
	assert.Equal(t, items.ItemTypeTask, input.ItemType)
	assert.Equal(t, "Test Title", input.Title)
	assert.Equal(t, "Test Body", input.Body)
	assert.Equal(t, "in-progress", input.Status)
	assert.Equal(t, "work", input.Tag)
}

func TestEditorFinishedMsg_Struct(t *testing.T) {
	msg := EditorFinishedMsg{
		Id:       "msg-id",
		ItemType: items.ItemTypeNote,
		Title:    "Msg Title",
		Body:     "Msg Body",
		Status:   "todo",
		Tag:      "docs",
		Err:      nil,
	}

	assert.Equal(t, "msg-id", msg.Id)
	assert.Equal(t, items.ItemTypeNote, msg.ItemType)
	assert.Equal(t, "Msg Title", msg.Title)
	assert.Equal(t, "Msg Body", msg.Body)
	assert.Equal(t, "todo", msg.Status)
	assert.Equal(t, "docs", msg.Tag)
	assert.Nil(t, msg.Err)
}
