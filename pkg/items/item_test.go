package items

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseItemType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ItemType
	}{
		{"task lowercase", "task", ItemTypeTask},
		{"task uppercase", "TASK", ItemTypeTask},
		{"task mixed case", "Task", ItemTypeTask},
		{"note lowercase", "note", ItemTypeNote},
		{"note uppercase", "NOTE", ItemTypeNote},
		{"note mixed case", "Note", ItemTypeNote},
		{"invalid returns empty", "invalid", ""},
		{"empty returns empty", "", ""},
		{"random string returns empty", "something", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseItemType(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestItemTypeConstants(t *testing.T) {
	assert.Equal(t, ItemType("task"), ItemTypeTask)
	assert.Equal(t, ItemType("note"), ItemTypeNote)
}

func TestItem_GetId(t *testing.T) {
	item := Item{Id: "test-id-123"}
	assert.Equal(t, "test-id-123", item.GetId())
}

func TestItem_GetTitle(t *testing.T) {
	item := Item{Title: "Test Title"}
	assert.Equal(t, "Test Title", item.GetTitle())
}

func TestItem_GetBody(t *testing.T) {
	item := Item{Body: "Test body content"}
	assert.Equal(t, "Test body content", item.GetBody())
}

func TestItem_GetCreatedAt(t *testing.T) {
	now := time.Now()
	item := Item{CreatedAt: now}
	assert.Equal(t, now, item.GetCreatedAt())
}

func TestItem_GetTag(t *testing.T) {
	t.Run("with tag", func(t *testing.T) {
		tag := &Tag{Id: "tag-1", Name: "work"}
		item := Item{Tag: tag}
		assert.Equal(t, tag, item.GetTag())
	})

	t.Run("without tag", func(t *testing.T) {
		item := Item{Tag: nil}
		assert.Nil(t, item.GetTag())
	})
}

func TestItem_After(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)
	later := now.Add(1 * time.Hour)

	tag := &Tag{Id: "tag-1", Name: "work"}

	tests := []struct {
		name     string
		itemI    Item
		itemU    ItemInterface
		expected bool
	}{
		{
			name:     "item with tag comes before item without tag",
			itemI:    Item{Tag: tag, CreatedAt: now},
			itemU:    &Task{Item: Item{Tag: nil, CreatedAt: now}},
			expected: false, // i should come before u, so After returns false
		},
		{
			name:     "item without tag comes after item with tag",
			itemI:    Item{Tag: nil, CreatedAt: now},
			itemU:    &Task{Item: Item{Tag: tag, CreatedAt: now}},
			expected: true, // i should come after u
		},
		{
			name:     "both have tags - earlier created comes after",
			itemI:    Item{Tag: tag, CreatedAt: earlier},
			itemU:    &Task{Item: Item{Tag: tag, CreatedAt: later}},
			expected: true, // earlier is before later, so i comes after u
		},
		{
			name:     "both have tags - later created comes before",
			itemI:    Item{Tag: tag, CreatedAt: later},
			itemU:    &Task{Item: Item{Tag: tag, CreatedAt: earlier}},
			expected: false, // later is after earlier, so i comes before u
		},
		{
			name:     "neither has tag - earlier created comes after",
			itemI:    Item{Tag: nil, CreatedAt: earlier},
			itemU:    &Task{Item: Item{Tag: nil, CreatedAt: later}},
			expected: true,
		},
		{
			name:     "neither has tag - later created comes before",
			itemI:    Item{Tag: nil, CreatedAt: later},
			itemU:    &Task{Item: Item{Tag: nil, CreatedAt: earlier}},
			expected: false,
		},
		{
			name:     "same time both with tags",
			itemI:    Item{Tag: tag, CreatedAt: now},
			itemU:    &Task{Item: Item{Tag: tag, CreatedAt: now}},
			expected: false, // same time, Before returns false
		},
		{
			name:     "same time both without tags",
			itemI:    Item{Tag: nil, CreatedAt: now},
			itemU:    &Task{Item: Item{Tag: nil, CreatedAt: now}},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.itemI.After(tc.itemU)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestItem_After_WithNotes(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)
	tag := &Tag{Id: "tag-1", Name: "work"}

	t.Run("task vs note - tag comparison", func(t *testing.T) {
		taskWithTag := Item{Tag: tag, CreatedAt: now}
		noteWithoutTag := &Note{Item: Item{Tag: nil, CreatedAt: now}}

		assert.False(t, taskWithTag.After(noteWithoutTag), "Item with tag should come before item without tag")
	})

	t.Run("task vs note - time comparison", func(t *testing.T) {
		taskEarlier := Item{Tag: nil, CreatedAt: earlier}
		noteLater := &Note{Item: Item{Tag: nil, CreatedAt: now}}

		assert.True(t, taskEarlier.After(noteLater), "Earlier item should come after later item (most recent first)")
	})
}
