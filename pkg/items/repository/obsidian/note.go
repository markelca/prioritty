package obsidian

import (
	"log"
	"os"
	"time"

	"github.com/markelca/prioritty/pkg/items"
)

// GetNotes returns all notes from the vault.
func (r *ObsidianRepository) GetNotes() ([]items.Note, error) {
	files, err := scanMarkdownFiles(r.vaultPath)
	if err != nil {
		return nil, err
	}

	var notes []items.Note
	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: failed to read file %s: %v", filePath, err)
			continue
		}

		fm, body, err := parseFrontmatter(content)
		if err != nil {
			log.Printf("Warning: failed to parse frontmatter in %s: %v", filePath, err)
			continue
		}

		if fm.Type != string(items.ItemTypeNote) {
			continue
		}

		id := relativeID(r.vaultPath, filePath)
		note := noteFromFrontmatter(fm, body, id)
		notes = append(notes, note)
	}

	return notes, nil
}

// CreateNote creates a new note file in the vault.
func (r *ObsidianRepository) CreateNote(n items.Note) error {
	// Set created_at if not set
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}

	// Generate unique filename
	filePath := uniqueFilename(r.vaultPath, n.Title)

	// Create frontmatter
	fm := frontmatterFromNote(n)

	// Serialize to markdown
	content, err := serializeFrontmatter(fm, n.Body)
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filePath, content, 0644)
}

// UpdateNote updates an existing note file.
func (r *ObsidianRepository) UpdateNote(n items.Note) error {
	oldPath := fullPathFromID(r.vaultPath, n.Id)

	// Read existing file to preserve created_at if not set
	existingContent, err := os.ReadFile(oldPath)
	if err != nil {
		return err
	}

	existingFm, _, err := parseFrontmatter(existingContent)
	if err != nil {
		return err
	}

	// Preserve created_at from existing file if not set on update
	if n.CreatedAt.IsZero() {
		n.CreatedAt = parseCreatedAt(existingFm.CreatedAt)
	}

	// Create frontmatter
	fm := frontmatterFromNote(n)

	// Serialize to markdown
	content, err := serializeFrontmatter(fm, n.Body)
	if err != nil {
		return err
	}

	// Check if title changed (requires rename)
	newFilename := filenameFromTitle(n.Title)
	oldFilename := relativeID(r.vaultPath, oldPath)

	if newFilename != oldFilename {
		// Generate unique filename for new title
		newPath := uniqueFilename(r.vaultPath, n.Title)

		// Write new file
		if err := os.WriteFile(newPath, content, 0644); err != nil {
			return err
		}

		// Remove old file
		return os.Remove(oldPath)
	}

	// Title unchanged, write in place
	return os.WriteFile(oldPath, content, 0644)
}

// RemoveNote removes a note file from the vault.
func (r *ObsidianRepository) RemoveNote(id string) error {
	filePath := fullPathFromID(r.vaultPath, id)
	return os.Remove(filePath)
}

// SetNoteTag sets the tag on a note.
func (r *ObsidianRepository) SetNoteTag(n items.Note, tag items.Tag) error {
	filePath := fullPathFromID(r.vaultPath, n.Id)

	// Read existing file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fm, body, err := parseFrontmatter(content)
	if err != nil {
		return err
	}

	// Update tag
	fm.Tag = tag.Name

	// Serialize and write back
	newContent, err := serializeFrontmatter(fm, body)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, newContent, 0644)
}

// UnsetNoteTag removes the tag from a note.
func (r *ObsidianRepository) UnsetNoteTag(n items.Note) error {
	filePath := fullPathFromID(r.vaultPath, n.Id)

	// Read existing file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fm, body, err := parseFrontmatter(content)
	if err != nil {
		return err
	}

	// Remove tag
	fm.Tag = ""

	// Serialize and write back
	newContent, err := serializeFrontmatter(fm, body)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, newContent, 0644)
}
