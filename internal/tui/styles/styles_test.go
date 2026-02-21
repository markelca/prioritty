package styles

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderDeleteDialog(t *testing.T) {
	tests := []struct {
		name            string
		itemTitle       string
		expectedTitle   string
		shouldTruncate  bool
		truncatedLength int
	}{
		{
			name:           "short title - no truncation",
			itemTitle:      "Short task",
			expectedTitle:  "Short task",
			shouldTruncate: false,
		},
		{
			name:           "title exactly 30 chars - no truncation",
			itemTitle:      "This is exactly thirty chars!",
			expectedTitle:  "This is exactly thirty chars!",
			shouldTruncate: false,
		},
		{
			name:            "title 31 chars - truncation",
			itemTitle:       "This is exactly thirty-one char",
			expectedTitle:   "This is exactly thirty-o...",
			shouldTruncate:  true,
			truncatedLength: 30,
		},
		{
			name:            "very long title - truncation",
			itemTitle:       "This is a very long task title that needs to be truncated because it exceeds the limit",
			expectedTitle:   "This is a very long task...",
			shouldTruncate:  true,
			truncatedLength: 30,
		},
		{
			name:           "empty title",
			itemTitle:      "",
			expectedTitle:  "",
			shouldTruncate: false,
		},
		{
			name:           "single char title",
			itemTitle:      "X",
			expectedTitle:  "X",
			shouldTruncate: false,
		},
		{
			name:           "title with special characters",
			itemTitle:      "Task: fix bug #123 @work",
			expectedTitle:  "Task: fix bug #123 @work",
			shouldTruncate: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := RenderDeleteDialog(tc.itemTitle)

			require.NotEmpty(t, result, "Dialog should not be empty")
			assert.Contains(t, result, "Delete item?", "Should contain dialog title")
			assert.Contains(t, result, "y", "Should contain confirm key")
			assert.Contains(t, result, "n", "Should contain cancel key")

			if tc.shouldTruncate {
				truncatedTitle := tc.itemTitle[:27] + "..."
				assert.Contains(t, result, "\""+truncatedTitle+"\"", "Should contain truncated title in quotes")
			} else {
				assert.Contains(t, result, "\""+tc.expectedTitle+"\"", "Should contain full title in quotes")
			}
		})
	}
}

func TestRenderDeleteDialog_ContainsAllElements(t *testing.T) {
	result := RenderDeleteDialog("Test Item")

	assert.Contains(t, result, "Delete item?")
	assert.Contains(t, result, "Test Item")
	assert.Contains(t, result, "y")
	assert.Contains(t, result, "n")
	assert.Contains(t, result, "to confirm")
	assert.Contains(t, result, "to cancel")
}

func TestRenderDeleteDialog_TruncationExact(t *testing.T) {
	tests := []struct {
		name      string
		itemTitle string
		maxLen    int
	}{
		{
			name:      "31 characters should be truncated",
			itemTitle: "1234567890123456789012345678901",
			maxLen:    30,
		},
		{
			name:      "50 characters should be truncated",
			itemTitle: "12345678901234567890123456789012345678901234567890",
			maxLen:    30,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := RenderDeleteDialog(tc.itemTitle)

			assert.True(t, len(tc.itemTitle) > 30, "Input should be longer than 30 chars")

			expectedInResult := tc.itemTitle[:27] + "..."
			assert.Contains(t, result, expectedInResult)
		})
	}
}

func TestRenderDeleteDialog_NoTruncationWhenExactOrLess(t *testing.T) {
	tests := []struct {
		name      string
		itemTitle string
	}{
		{
			name:      "exactly 30 characters",
			itemTitle: "123456789012345678901234567890",
		},
		{
			name:      "29 characters",
			itemTitle: "12345678901234567890123456789",
		},
		{
			name:      "10 characters",
			itemTitle: "1234567890",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := RenderDeleteDialog(tc.itemTitle)

			assert.NotContains(t, result, "...", "Should not contain truncation ellipsis")
			assert.Contains(t, result, tc.itemTitle)
		})
	}
}

func TestStyleVariables_AreNotEmpty(t *testing.T) {
	assert.NotEmpty(t, NoteIcon)
	assert.NotEmpty(t, DoneIcon)
	assert.NotEmpty(t, TodoIcon)
	assert.NotEmpty(t, InProgressIcon)
	assert.NotEmpty(t, CancelledIcon)
	assert.NotEmpty(t, ContentIcon)
}

func TestStyleVariables_ContainExpectedIcons(t *testing.T) {
	assert.Contains(t, NoteIcon, "i")
	assert.Contains(t, DoneIcon, "✓")
	assert.Contains(t, TodoIcon, "☐")
	assert.Contains(t, InProgressIcon, "◐")
	assert.Contains(t, CancelledIcon, "x")
	assert.Contains(t, ContentIcon, "≣")
}

func TestStyleVariables_HavePadding(t *testing.T) {
	assert.True(t, strings.HasSuffix(NoteIcon, " "), "NoteIcon should have right padding")
	assert.True(t, strings.HasSuffix(DoneIcon, " "), "DoneIcon should have right padding")
	assert.True(t, strings.HasSuffix(TodoIcon, " "), "TodoIcon should have right padding")
	assert.True(t, strings.HasSuffix(InProgressIcon, " "), "InProgressIcon should have right padding")
	assert.True(t, strings.HasSuffix(CancelledIcon, " "), "CancelledIcon should have right padding")
	assert.True(t, strings.HasSuffix(ContentIcon, " "), "ContentIcon should have right padding")
}
