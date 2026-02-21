package tui

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
)

func TestGetItemIcon(t *testing.T) {
	tests := []struct {
		name     string
		item     items.ItemInterface
		expected string
	}{
		{
			name:     "task with todo status",
			item:     &items.Task{Status: items.Todo},
			expected: "☐ ",
		},
		{
			name:     "task with done status",
			item:     &items.Task{Status: items.Done},
			expected: "✓ ",
		},
		{
			name:     "task with in-progress status",
			item:     &items.Task{Status: items.InProgress},
			expected: "◐ ",
		},
		{
			name:     "task with cancelled status",
			item:     &items.Task{Status: items.Cancelled},
			expected: "x ",
		},
		{
			name:     "note",
			item:     &items.Note{},
			expected: "i ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetItemIcon(tc.item)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRenderDonePercentage(t *testing.T) {
	tests := []struct {
		name          string
		items         []items.ItemInterface
		expectZero    bool
		expectPercent string
	}{
		{
			name:          "empty list",
			items:         []items.ItemInterface{},
			expectZero:    true,
			expectPercent: "0%",
		},
		{
			name: "all done tasks",
			items: []items.ItemInterface{
				&items.Task{Status: items.Done},
				&items.Task{Status: items.Done},
			},
			expectZero:    false,
			expectPercent: "100%",
		},
		{
			name: "half done tasks",
			items: []items.ItemInterface{
				&items.Task{Status: items.Done},
				&items.Task{Status: items.Todo},
			},
			expectZero:    false,
			expectPercent: "50%",
		},
		{
			name: "mixed with notes",
			items: []items.ItemInterface{
				&items.Task{Status: items.Done},
				&items.Note{},
			},
			expectZero:    false,
			expectPercent: "100%",
		},
		{
			name: "only notes",
			items: []items.ItemInterface{
				&items.Note{},
				&items.Note{},
			},
			expectZero:    true,
			expectPercent: "0%",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			counts := make(map[items.Status]int)
			for _, item := range tc.items {
				if task, ok := item.(*items.Task); ok {
					counts[task.Status]++
				} else if _, ok := item.(*items.Note); ok {
					counts[items.NoteType]++
				}
			}
			result := renderDonePercentage(tc.items, counts)
			assert.Contains(t, result, tc.expectPercent)
		})
	}
}

func TestRenderContentCount(t *testing.T) {
	tests := []struct {
		name          string
		items         []items.ItemInterface
		expectedCount int
	}{
		{
			name:          "empty list",
			items:         []items.ItemInterface{},
			expectedCount: 0,
		},
		{
			name: "all with content",
			items: []items.ItemInterface{
				&items.Task{Item: items.Item{Body: "content"}},
				&items.Note{Item: items.Item{Body: "content"}},
			},
			expectedCount: 2,
		},
		{
			name: "none with content",
			items: []items.ItemInterface{
				&items.Task{Item: items.Item{Body: ""}},
				&items.Note{Item: items.Item{Body: ""}},
			},
			expectedCount: 0,
		},
		{
			name: "mixed content",
			items: []items.ItemInterface{
				&items.Task{Item: items.Item{Body: "has content"}},
				&items.Note{Item: items.Item{Body: ""}},
			},
			expectedCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := renderContentCount(tc.items)
			switch tc.expectedCount {
			case 0:
				assert.Contains(t, result, "0 items with content")
			case 1:
				assert.Contains(t, result, "1 items with content")
			default:
				assert.Contains(t, result, "2 items with content")
			}
		})
	}
}
