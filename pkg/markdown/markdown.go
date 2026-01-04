package markdown

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/markelca/prioritty/pkg/items"
	"gopkg.in/yaml.v3"
)

const Delimiter = "---"

// unquotedString is a string type that marshals to YAML without quotes.
type unquotedString string

func (s unquotedString) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: string(s),
	}, nil
}

// TaskFrontmatter represents the YAML frontmatter for tasks.
type TaskFrontmatter struct {
	Title  unquotedString `yaml:"title"`
	Status unquotedString `yaml:"status"`
	Tag    unquotedString `yaml:"tag"`
}

// NoteFrontmatter represents the YAML frontmatter for notes.
type NoteFrontmatter struct {
	Title unquotedString `yaml:"title"`
	Tag   unquotedString `yaml:"tag"`
}

// ItemInput contains the data to serialize an item to markdown.
type ItemInput struct {
	ItemType items.ItemType
	Title    string
	Body     string
	Status   string
	Tag      string
}

// Parse extracts frontmatter and body from markdown content.
// The frontmatter struct must be passed as a pointer.
func Parse[T any](content string, fm *T) (body string, err error) {
	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, Delimiter) {
		return content, fmt.Errorf("no frontmatter found")
	}

	// Find the closing delimiter
	rest := content[len(Delimiter):]
	endIdx := strings.Index(rest, "\n"+Delimiter)
	if endIdx == -1 {
		return content, fmt.Errorf("unclosed frontmatter")
	}

	// Extract frontmatter YAML (skip the leading newline if present)
	fmContent := rest[:endIdx]
	fmContent = strings.TrimPrefix(fmContent, "\n")

	// Parse YAML
	if err := yaml.Unmarshal([]byte(fmContent), fm); err != nil {
		return content, fmt.Errorf("invalid frontmatter YAML: %w", err)
	}

	// Extract body (after closing delimiter and newline)
	bodyStart := len(Delimiter) + endIdx + len("\n"+Delimiter)
	body = content[bodyStart:]
	body = strings.TrimPrefix(body, "\n")

	return body, nil
}

// SerializeFrontmatter creates markdown content from a frontmatter struct and body.
func SerializeFrontmatter[T any](fm T, body string) ([]byte, error) {
	var buf bytes.Buffer

	// Write opening delimiter
	buf.WriteString(Delimiter)
	buf.WriteString("\n")

	// Marshal frontmatter to YAML
	yamlBytes, err := yaml.Marshal(&fm)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
	}
	buf.Write(yamlBytes)

	// Write closing delimiter
	buf.WriteString(Delimiter)
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

// Serialize creates markdown content with frontmatter from an ItemInput.
func Serialize(input ItemInput) (string, error) {
	var content []byte
	var err error

	if input.ItemType == items.ItemTypeTask {
		fm := TaskFrontmatter{
			Title:  unquotedString(input.Title),
			Status: unquotedString(input.Status),
			Tag:    unquotedString(input.Tag),
		}
		content, err = SerializeFrontmatter(fm, input.Body)
	} else {
		fm := NoteFrontmatter{
			Title: unquotedString(input.Title),
			Tag:   unquotedString(input.Tag),
		}
		content, err = SerializeFrontmatter(fm, input.Body)
	}

	if err != nil {
		return "", err
	}

	return string(content), nil
}
