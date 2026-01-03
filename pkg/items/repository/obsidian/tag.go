package obsidian

import (
	"log"
	"os"
	"sort"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
	"github.com/markelca/prioritty/pkg/markdown"
)

// GetTag returns a tag by name if it's used by any item.
// Returns ErrNotFound if the tag doesn't exist (not used by any item).
func (r *ObsidianRepository) GetTag(name string) (*items.Tag, error) {
	files, err := scanMarkdownFiles(r.vaultPath)
	if err != nil {
		return nil, err
	}

	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var fm markdown.Frontmatter
		if _, err := markdown.Parse(string(content), &fm); err != nil {
			continue
		}

		if fm.Tag == name {
			return &items.Tag{
				Id:   name,
				Name: name,
			}, nil
		}
	}

	return nil, repository.ErrNotFound
}

// GetTags returns all unique tags used by items in the vault.
func (r *ObsidianRepository) GetTags() ([]items.Tag, error) {
	files, err := scanMarkdownFiles(r.vaultPath)
	if err != nil {
		return nil, err
	}

	tagSet := make(map[string]struct{})

	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: failed to read file %s: %v", filePath, err)
			continue
		}

		var fm markdown.Frontmatter
		if _, err := markdown.Parse(string(content), &fm); err != nil {
			log.Printf("Warning: failed to parse frontmatter in %s: %v", filePath, err)
			continue
		}

		if fm.Tag != "" {
			tagSet[fm.Tag] = struct{}{}
		}
	}

	// Convert to slice and sort by name
	var tags []items.Tag
	for name := range tagSet {
		tags = append(tags, items.Tag{
			Id:   name,
			Name: name,
		})
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	return tags, nil
}

// CreateTag is a no-op for Obsidian since tags are created implicitly when assigned.
// It returns a Tag struct with the given name.
func (r *ObsidianRepository) CreateTag(name string) (*items.Tag, error) {
	return &items.Tag{
		Id:   name,
		Name: name,
	}, nil
}

// RemoveTag is a no-op for Obsidian since tags only exist when used.
// The service layer checks if items use the tag before calling this.
func (r *ObsidianRepository) RemoveTag(name string) error {
	// Check if tag exists (is used by any item)
	_, err := r.GetTag(name)
	if err != nil {
		return err
	}
	// If we get here, the tag exists (is used by items)
	// The service layer should prevent this call if items use the tag
	// This is a no-op since removing unused tags is implicit
	return nil
}

// GetItemsWithTag returns all items (tasks and notes) with the given tag.
func (r *ObsidianRepository) GetItemsWithTag(tagName string) ([]items.ItemInterface, error) {
	files, err := scanMarkdownFiles(r.vaultPath)
	if err != nil {
		return nil, err
	}

	var result []items.ItemInterface

	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: failed to read file %s: %v", filePath, err)
			continue
		}

		var fm markdown.Frontmatter
		body, err := markdown.Parse(string(content), &fm)
		if err != nil {
			log.Printf("Warning: failed to parse frontmatter in %s: %v", filePath, err)
			continue
		}

		if fm.Tag != tagName {
			continue
		}

		id := relativeID(r.vaultPath, filePath)

		switch fm.Type {
		case string(items.ItemTypeTask):
			task := taskFromFrontmatter(fm, body, id)
			result = append(result, &task)
		case string(items.ItemTypeNote):
			note := noteFromFrontmatter(fm, body, id)
			result = append(result, &note)
		}
	}

	return result, nil
}
