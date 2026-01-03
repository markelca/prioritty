package obsidian

import (
	"time"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/markdown"
)

const timeFormat = time.RFC3339

// parseCreatedAt parses the created_at string to time.Time.
func parseCreatedAt(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		return time.Now()
	}
	return t
}

// formatCreatedAt formats time.Time to the frontmatter string format.
func formatCreatedAt(t time.Time) string {
	if t.IsZero() {
		return time.Now().Format(timeFormat)
	}
	return t.Format(timeFormat)
}

// taskFromFrontmatter creates a Task from frontmatter data.
func taskFromFrontmatter(fm markdown.Frontmatter, body, id string) items.Task {
	var tag *items.Tag
	if fm.Tag != "" {
		tag = &items.Tag{
			Id:   fm.Tag,
			Name: fm.Tag,
		}
	}

	return items.Task{
		Item: items.Item{
			Id:        id,
			Title:     fm.Title,
			Body:      body,
			CreatedAt: parseCreatedAt(fm.CreatedAt),
			Tag:       tag,
		},
		Status: items.ParseStatus(fm.Status),
	}
}

// noteFromFrontmatter creates a Note from frontmatter data.
func noteFromFrontmatter(fm markdown.Frontmatter, body, id string) items.Note {
	var tag *items.Tag
	if fm.Tag != "" {
		tag = &items.Tag{
			Id:   fm.Tag,
			Name: fm.Tag,
		}
	}

	return items.Note{
		Item: items.Item{
			Id:        id,
			Title:     fm.Title,
			Body:      body,
			CreatedAt: parseCreatedAt(fm.CreatedAt),
			Tag:       tag,
		},
	}
}

// itemInputFromTask creates an ItemInput from a Task.
func itemInputFromTask(t items.Task) markdown.ItemInput {
	input := markdown.ItemInput{
		ItemType:  items.ItemTypeTask,
		Title:     t.Title,
		Body:      t.Body,
		Status:    string(t.Status),
		CreatedAt: formatCreatedAt(t.CreatedAt),
	}
	if t.Tag != nil {
		input.Tag = t.Tag.Name
	}
	return input
}

// itemInputFromNote creates an ItemInput from a Note.
func itemInputFromNote(n items.Note) markdown.ItemInput {
	input := markdown.ItemInput{
		ItemType:  items.ItemTypeNote,
		Title:     n.Title,
		Body:      n.Body,
		CreatedAt: formatCreatedAt(n.CreatedAt),
	}
	if n.Tag != nil {
		input.Tag = n.Tag.Name
	}
	return input
}
