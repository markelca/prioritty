package frontmatter

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const Delimiter = "---"

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

// Serialize creates markdown content from frontmatter and body.
func Serialize[T any](fm T, body string) ([]byte, error) {
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
