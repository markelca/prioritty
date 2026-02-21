package tui

import (
	"testing"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/stretchr/testify/assert"
)

func TestSortItemsByTag(t *testing.T) {
	tagA := &items.Tag{Id: "1", Name: "alpha"}
	tagB := &items.Tag{Id: "2", Name: "beta"}

	tests := []struct {
		name     string
		input    []items.ItemInterface
		expected []string
	}{
		{
			name:     "empty list",
			input:    []items.ItemInterface{},
			expected: []string{},
		},
		{
			name: "single item without tag",
			input: []items.ItemInterface{
				&items.Task{Item: items.Item{Title: "task1"}},
			},
			expected: []string{"task1"},
		},
		{
			name: "items grouped by tag order",
			input: []items.ItemInterface{
				&items.Task{Item: items.Item{Title: "task1", Tag: tagA}, Status: items.Todo},
				&items.Task{Item: items.Item{Title: "task2"}, Status: items.Todo},
				&items.Task{Item: items.Item{Title: "task3", Tag: tagB}, Status: items.Todo},
				&items.Task{Item: items.Item{Title: "task4", Tag: tagA}, Status: items.Todo},
			},
			expected: []string{"task1", "task4", "task2", "task3"},
		},
		{
			name: "items without tags appear first in tag order",
			input: []items.ItemInterface{
				&items.Task{Item: items.Item{Title: "no-tag"}, Status: items.Todo},
				&items.Task{Item: items.Item{Title: "tagged", Tag: tagA}, Status: items.Todo},
			},
			expected: []string{"no-tag", "tagged"},
		},
		{
			name: "mixed tasks and notes",
			input: []items.ItemInterface{
				&items.Task{Item: items.Item{Title: "task", Tag: tagA}, Status: items.Todo},
				&items.Note{Item: items.Item{Title: "note", Tag: tagA}},
			},
			expected: []string{"task", "note"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := sortItemsByTag(tc.input)
			if len(tc.expected) == 0 {
				assert.Empty(t, result)
				return
			}
			var titles []string
			for _, item := range result {
				titles = append(titles, item.GetTitle())
			}
			assert.Equal(t, tc.expected, titles)
		})
	}
}

func TestGetItemAt(t *testing.T) {
	items := []items.ItemInterface{
		&items.Task{Item: items.Item{Id: "1", Title: "task1"}, Status: items.Todo},
		&items.Task{Item: items.Item{Id: "2", Title: "task2"}, Status: items.Todo},
	}

	m := Model{
		state: State{items: items},
	}

	tests := []struct {
		name     string
		index    int
		expected string
		isNil    bool
	}{
		{"first item", 0, "task1", false},
		{"second item", 1, "task2", false},
		{"negative index", -1, "", true},
		{"out of bounds", 2, "", true},
		{"large out of bounds", 100, "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := m.GetItemAt(tc.index)
			if tc.isNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.expected, result.GetTitle())
			}
		})
	}
}

func TestState_GetCurrentItem(t *testing.T) {
	itemList := []items.ItemInterface{
		&items.Task{Item: items.Item{Id: "1", Title: "task1"}, Status: items.Todo},
		&items.Task{Item: items.Item{Id: "2", Title: "task2"}, Status: items.Todo},
	}

	tests := []struct {
		name     string
		cursor   int
		items    []items.ItemInterface
		expected string
		isNil    bool
	}{
		{"first item", 0, itemList, "task1", false},
		{"second item", 1, itemList, "task2", false},
		{"negative cursor", -1, itemList, "", true},
		{"cursor out of bounds", 2, itemList, "", true},
		{"empty list", 0, []items.ItemInterface{}, "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := State{cursor: tc.cursor, items: tc.items}
			result := s.GetCurrentItem()
			if tc.isNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.expected, result.GetTitle())
			}
		})
	}
}

func TestMode_Constants(t *testing.T) {
	assert.Equal(t, Mode("list"), ModeList)
	assert.Equal(t, Mode("create"), ModeCreate)
	assert.Equal(t, Mode("edit"), ModeEdit)
	assert.Equal(t, Mode("delete_confirm"), ModeDeleteConfirm)
}
