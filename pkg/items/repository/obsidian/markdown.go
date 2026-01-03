package obsidian

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/markelca/prioritty/pkg/items"
	"gopkg.in/yaml.v3"
)

const (
	frontmatterDelimiter = "---"
	typeTask             = "task"
	typeNote             = "note"
	timeFormat           = time.RFC3339
)

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

	str := string(content)

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(str, frontmatterDelimiter) {
		return fm, str, fmt.Errorf("no frontmatter found")
	}

	// Find the closing delimiter
	rest := str[len(frontmatterDelimiter):]
	endIdx := strings.Index(rest, "\n"+frontmatterDelimiter)
	if endIdx == -1 {
		return fm, str, fmt.Errorf("unclosed frontmatter")
	}

	// Extract frontmatter YAML (skip the leading newline if present)
	fmContent := rest[:endIdx]
	fmContent = strings.TrimPrefix(fmContent, "\n")

	// Parse YAML
	if err := yaml.Unmarshal([]byte(fmContent), &fm); err != nil {
		return fm, str, fmt.Errorf("invalid frontmatter YAML: %w", err)
	}

	// Extract body (after closing delimiter and newline)
	bodyStart := len(frontmatterDelimiter) + endIdx + len("\n"+frontmatterDelimiter)
	body := str[bodyStart:]
	body = strings.TrimPrefix(body, "\n")

	return fm, body, nil
}

// serializeFrontmatter creates markdown content from frontmatter and body.
func serializeFrontmatter(fm Frontmatter, body string) ([]byte, error) {
	var buf bytes.Buffer

	// Write opening delimiter
	buf.WriteString(frontmatterDelimiter)
	buf.WriteString("\n")

	// Marshal frontmatter to YAML
	yamlBytes, err := yaml.Marshal(&fm)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
	}
	buf.Write(yamlBytes)

	// Write closing delimiter
	buf.WriteString(frontmatterDelimiter)
	buf.WriteString("\n")

	// Write body
	if body != "" {
		buf.WriteString(body)
		// Ensure file ends with newline
		if !strings.HasSuffix(body, "\n") {
			buf.WriteString("\n")
		}
	}

	return buf.Bytes(), nil
}

// statusToString converts a Status enum to its frontmatter string value.
func statusToString(s items.Status) string {
	switch s {
	case items.Todo:
		return "todo"
	case items.InProgress:
		return "in-progress"
	case items.Done:
		return "done"
	case items.Cancelled:
		return "cancelled"
	default:
		return "todo"
	}
}

// stringToStatus converts a frontmatter string to Status enum.
func stringToStatus(s string) items.Status {
	switch strings.ToLower(s) {
	case "todo":
		return items.Todo
	case "in-progress", "inprogress":
		return items.InProgress
	case "done":
		return items.Done
	case "cancelled", "canceled":
		return items.Cancelled
	default:
		return items.Todo
	}
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
		Status: stringToStatus(fm.Status),
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
		Type:      typeTask,
		Status:    statusToString(t.Status),
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
		Type:      typeNote,
		CreatedAt: formatCreatedAt(n.CreatedAt),
	}
	if n.Tag != nil {
		fm.Tag = n.Tag.Name
	}
	return fm
}
