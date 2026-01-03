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

// Frontmatter represents the YAML frontmatter for items.
type Frontmatter struct {
	Title     string `yaml:"title"`
	Type      string `yaml:"type,omitempty"`
	Status    string `yaml:"status,omitempty"`
	Tag       string `yaml:"tag,omitempty"`
	CreatedAt string `yaml:"created_at,omitempty"`
}

// unquotedFrontmatter is used internally for serialization to produce clean YAML without quotes.
type unquotedFrontmatter struct {
	Title     unquotedString `yaml:"title"`
	Type      unquotedString `yaml:"type,omitempty"`
	Status    unquotedString `yaml:"status,omitempty"`
	Tag       unquotedString `yaml:"tag,omitempty"`
	CreatedAt unquotedString `yaml:"created_at,omitempty"`
}

// toUnquoted converts a Frontmatter to unquotedFrontmatter for serialization.
func (fm Frontmatter) toUnquoted() unquotedFrontmatter {
	return unquotedFrontmatter{
		Title:     unquotedString(fm.Title),
		Type:      unquotedString(fm.Type),
		Status:    unquotedString(fm.Status),
		Tag:       unquotedString(fm.Tag),
		CreatedAt: unquotedString(fm.CreatedAt),
	}
}

// Serialize converts a Frontmatter to markdown content with body.
func (fm Frontmatter) Serialize(body string) ([]byte, error) {
	return SerializeFrontmatter(fm.toUnquoted(), body)
}

// ItemInput contains the data to serialize an item to markdown.
type ItemInput struct {
	ItemType  items.ItemType
	Title     string
	Body      string
	Status    string
	Tag       string
	CreatedAt string // Only populated when serializing for storage/display, not for editor
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
	fm := Frontmatter{
		Title:     input.Title,
		Type:      string(input.ItemType),
		Tag:       input.Tag,
		CreatedAt: input.CreatedAt,
	}

	// Only include status for tasks
	if input.ItemType == items.ItemTypeTask {
		fm.Status = input.Status
	}

	content, err := SerializeFrontmatter(fm.toUnquoted(), input.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
