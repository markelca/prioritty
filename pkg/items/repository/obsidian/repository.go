package obsidian

import (
	"os"
	"path/filepath"
)

// ObsidianRepository implements repository.Repository using Obsidian markdown files.
type ObsidianRepository struct {
	vaultPath string
}

// NewObsidianRepository creates a new ObsidianRepository for the given vault path.
func NewObsidianRepository(vaultPath string) *ObsidianRepository {
	return &ObsidianRepository{
		vaultPath: vaultPath,
	}
}

// VaultPath returns the path to the Obsidian vault.
func (r *ObsidianRepository) VaultPath() string {
	return r.vaultPath
}

// Reset removes all markdown files from the vault (used for demo cleanup).
// It preserves the .obsidian folder.
func (r *ObsidianRepository) Reset() error {
	entries, err := os.ReadDir(r.vaultPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		// Skip .obsidian folder
		if name == ".obsidian" {
			continue
		}
		fullPath := filepath.Join(r.vaultPath, name)
		if entry.IsDir() {
			continue
		}
		// Remove markdown files
		if filepath.Ext(name) == ".md" {
			if err := os.Remove(fullPath); err != nil {
				return err
			}
		}
	}
	return nil
}
