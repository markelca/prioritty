package obsidian

import (
	"time"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/markdown"
)

const timeFormat = time.RFC3339

// Frontmatter represents the YAML frontmatter in a markdown file.
type Frontmatter struct {
	Title     string `yaml:"title"`
	Type      string `yaml:"type"`
	Status    string `yaml:"status,omitempty"`
	Tag       string `yaml:"tag,omitempty"`
	CreatedAt string `yaml:"created_at,omitempty"`
}

// parseFrontmatter extracts frontmatter and body from markdown content.
func parseFrontmatter(content []byte) (Frontmatter, string, error) {
	var fm Frontmatter
	body, err := markdown.Parse(string(content), &fm)
	return fm, body, err
}

// serializeFrontmatter creates markdown content from frontmatter and body.
func serializeFrontmatter(fm Frontmatter, body string) ([]byte, error) {
	return markdown.SerializeFrontmatter(fm, body)
}

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
func taskFromFrontmatter(fm Frontmatter, body, id string) items.Task {
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
func noteFromFrontmatter(fm Frontmatter, body, id string) items.Note {
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

// frontmatterFromTask creates frontmatter from a Task.
func frontmatterFromTask(t items.Task) Frontmatter {
	fm := Frontmatter{
		Title:     t.Title,
		Type:      string(items.ItemTypeTask),
		Status:    string(t.Status),
		CreatedAt: formatCreatedAt(t.CreatedAt),
	}
	if t.Tag != nil {
		fm.Tag = t.Tag.Name
	}
	return fm
}

// frontmatterFromNote creates frontmatter from a Note.
func frontmatterFromNote(n items.Note) Frontmatter {
	fm := Frontmatter{
		Title:     n.Title,
		Type:      string(items.ItemTypeNote),
		CreatedAt: formatCreatedAt(n.CreatedAt),
	}
	if n.Tag != nil {
		fm.Tag = n.Tag.Name
	}
	return fm
}
