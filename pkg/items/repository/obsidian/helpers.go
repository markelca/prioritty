package obsidian

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// toKebabCase converts a title string to kebab-case for use as filename.
// It handles unicode normalization, removes special characters, and converts spaces to hyphens.
func toKebabCase(title string) string {
	// Normalize unicode (NFD) and remove diacritics
	normalized := norm.NFD.String(title)
	var builder strings.Builder
	for _, r := range normalized {
		if unicode.Is(unicode.Mn, r) {
			// Skip combining marks (diacritics)
			continue
		}
		builder.WriteRune(r)
	}
	result := builder.String()

	// Convert to lowercase
	result = strings.ToLower(result)

	// Replace spaces and underscores with hyphens
	result = strings.ReplaceAll(result, " ", "-")
	result = strings.ReplaceAll(result, "_", "-")

	// Remove invalid filesystem characters: / \ : * ? " < > |
	invalidChars := regexp.MustCompile(`[/\\:*?"<>|]`)
	result = invalidChars.ReplaceAllString(result, "")

	// Remove any character that's not alphanumeric or hyphen
	validChars := regexp.MustCompile(`[^a-z0-9-]`)
	result = validChars.ReplaceAllString(result, "")

	// Collapse multiple hyphens into one
	multiHyphen := regexp.MustCompile(`-+`)
	result = multiHyphen.ReplaceAllString(result, "-")

	// Trim leading/trailing hyphens
	result = strings.Trim(result, "-")

	// If result is empty, use a default
	if result == "" {
		result = "untitled"
	}

	return result
}

// filenameFromTitle generates a .md filename from a title.
func filenameFromTitle(title string) string {
	return toKebabCase(title) + ".md"
}

// FilenameFromTitle is the exported version of filenameFromTitle.
func FilenameFromTitle(title string) string {
	return filenameFromTitle(title)
}

// scanMarkdownFiles returns all .md files in the vault directory (non-recursive).
// It excludes the .obsidian folder.
func scanMarkdownFiles(vaultPath string) ([]string, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".md") {
			files = append(files, filepath.Join(vaultPath, name))
		}
	}
	return files, nil
}

// uniqueFilename generates a unique filename by appending a counter if needed.
// Returns the full path to the file.
func uniqueFilename(vaultPath, title string) string {
	base := toKebabCase(title)
	filename := base + ".md"
	fullPath := filepath.Join(vaultPath, filename)

	// Check if file exists, if so append counter
	counter := 2
	for {
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fullPath
		}
		filename = base + "-" + strconv.Itoa(counter) + ".md"
		fullPath = filepath.Join(vaultPath, filename)
		counter++
	}
}

// relativeID returns the filename (relative to vault) from a full path.
// This is used as the item ID.
func relativeID(vaultPath, fullPath string) string {
	rel, err := filepath.Rel(vaultPath, fullPath)
	if err != nil {
		return filepath.Base(fullPath)
	}
	return rel
}

// fullPath converts an ID (relative filename) to a full path.
func fullPathFromID(vaultPath, id string) string {
	return filepath.Join(vaultPath, id)
}
