package obsidian

import (
	"log"
	"os"
	"time"

	"github.com/markelca/prioritty/pkg/items"
)

// GetTasks returns all tasks from the vault.
func (r *ObsidianRepository) GetTasks() ([]items.Task, error) {
	files, err := scanMarkdownFiles(r.vaultPath)
	if err != nil {
		return nil, err
	}

	var tasks []items.Task
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

		if fm.Type != typeTask {
			continue
		}

		id := relativeID(r.vaultPath, filePath)
		task := taskFromFrontmatter(fm, body, id)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// CreateTask creates a new task file in the vault.
func (r *ObsidianRepository) CreateTask(t items.Task) error {
	// Set created_at if not set
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}

	// Generate unique filename
	filePath := uniqueFilename(r.vaultPath, t.Title)

	// Create frontmatter
	fm := frontmatterFromTask(t)

	// Serialize to markdown
	content, err := serializeFrontmatter(fm, t.Body)
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filePath, content, 0644)
}

// UpdateTask updates an existing task file.
func (r *ObsidianRepository) UpdateTask(t items.Task) error {
	oldPath := fullPathFromID(r.vaultPath, t.Id)

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
	if t.CreatedAt.IsZero() {
		t.CreatedAt = parseCreatedAt(existingFm.CreatedAt)
	}

	// Create frontmatter
	fm := frontmatterFromTask(t)

	// Serialize to markdown
	content, err := serializeFrontmatter(fm, t.Body)
	if err != nil {
		return err
	}

	// Check if title changed (requires rename)
	newFilename := filenameFromTitle(t.Title)
	oldFilename := relativeID(r.vaultPath, oldPath)

	if newFilename != oldFilename {
		// Generate unique filename for new title
		newPath := uniqueFilename(r.vaultPath, t.Title)

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

// RemoveTask removes a task file from the vault.
func (r *ObsidianRepository) RemoveTask(id string) error {
	filePath := fullPathFromID(r.vaultPath, id)
	return os.Remove(filePath)
}

// UpdateTaskStatus updates only the status of a task.
func (r *ObsidianRepository) UpdateTaskStatus(t items.Task, status items.Status) error {
	filePath := fullPathFromID(r.vaultPath, t.Id)

	// Read existing file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fm, body, err := parseFrontmatter(content)
	if err != nil {
		return err
	}

	// Update status
	fm.Status = statusToString(status)

	// Serialize and write back
	newContent, err := serializeFrontmatter(fm, body)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, newContent, 0644)
}

// SetTaskTag sets the tag on a task.
func (r *ObsidianRepository) SetTaskTag(t items.Task, tag items.Tag) error {
	filePath := fullPathFromID(r.vaultPath, t.Id)

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

// UnsetTaskTag removes the tag from a task.
func (r *ObsidianRepository) UnsetTaskTag(t items.Task) error {
	filePath := fullPathFromID(r.vaultPath, t.Id)

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
