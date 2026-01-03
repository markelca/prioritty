package obsidian

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/markelca/prioritty/pkg/frontmatter"
	"github.com/markelca/prioritty/pkg/items/repository/obsidian"
	"github.com/spf13/viper"
)

//go:embed all_items.base.yaml
var allItemsBaseContent string

//go:embed demo_data.json
var demoDataJSON []byte

// TypesJSON represents the Obsidian types.json structure
type TypesJSON struct {
	Types map[string]string `json:"types"`
}

// DemoData represents the structure of demo_data.json
type DemoData struct {
	Tasks []DemoTask `json:"tasks"`
	Notes []DemoNote `json:"notes"`
}

type DemoTask struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status string `json:"status"`
	Tag    string `json:"tag"`
}

type DemoNote struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Tag   string `json:"tag"`
}

// defaultTypes returns the default property types for Prioritty
func defaultTypes() TypesJSON {
	return TypesJSON{
		Types: map[string]string{
			"title":      "text",
			"type":       "text",
			"status":     "text",
			"tag":        "text",
			"created_at": "datetime",
		},
	}
}

// NewObsidianRepository creates and initializes an Obsidian repository.
// It ensures the vault directory exists and initializes .obsidian/types.json.
func NewObsidianRepository(vaultPath string) (*obsidian.ObsidianRepository, error) {
	// Ensure vault directory exists
	if err := os.MkdirAll(vaultPath, 0755); err != nil {
		return nil, err
	}

	// Initialize .obsidian folder and types.json
	obsidianDir := filepath.Join(vaultPath, ".obsidian")
	if err := os.MkdirAll(obsidianDir, 0755); err != nil {
		return nil, err
	}

	typesPath := filepath.Join(obsidianDir, "types.json")
	if err := initTypesJSON(typesPath); err != nil {
		return nil, err
	}

	// Initialize All Items.base file
	basePath := filepath.Join(vaultPath, "All Items.base")
	if err := initBaseFile(basePath); err != nil {
		return nil, err
	}

	repo := obsidian.NewObsidianRepository(vaultPath)

	// Seed demo data if demo mode
	if viper.GetBool("demo") {
		if err := seedDemoData(repo); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// initBaseFile creates the All Items.base file if it doesn't exist.
func initBaseFile(basePath string) error {
	// Only create if file doesn't exist
	if _, err := os.Stat(basePath); err == nil {
		return nil
	}
	return os.WriteFile(basePath, []byte(allItemsBaseContent), 0644)
}

// initTypesJSON creates or updates the types.json file with Prioritty's property types.
func initTypesJSON(typesPath string) error {
	var types TypesJSON

	// Check if file exists and read it
	if data, err := os.ReadFile(typesPath); err == nil {
		if err := json.Unmarshal(data, &types); err != nil {
			// If invalid JSON, start fresh
			types = defaultTypes()
		} else {
			// Merge our types with existing (our types take precedence)
			defaults := defaultTypes()
			if types.Types == nil {
				types.Types = make(map[string]string)
			}
			for k, v := range defaults.Types {
				types.Types[k] = v
			}
		}
	} else {
		// File doesn't exist, create with defaults
		types = defaultTypes()
	}

	// Write back
	data, err := json.MarshalIndent(types, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(typesPath, data, 0644)
}

// seedDemoData creates sample tasks and notes for demo mode.
func seedDemoData(repo *obsidian.ObsidianRepository) error {
	// Check if vault already has files (skip seeding)
	files, err := os.ReadDir(repo.VaultPath())
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".md" {
			// Already has markdown files, skip seeding
			return nil
		}
	}

	// Load demo data from embedded JSON
	var demoData DemoData
	if err := json.Unmarshal(demoDataJSON, &demoData); err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)

	// Create demo tasks
	for _, t := range demoData.Tasks {
		fm := obsidian.Frontmatter{
			Title:     t.Title,
			Type:      "task",
			Status:    t.Status,
			Tag:       t.Tag,
			CreatedAt: now,
		}
		content, err := frontmatter.Serialize(fm, t.Body)
		if err != nil {
			return err
		}

		filename := obsidian.FilenameFromTitle(t.Title)
		filePath := filepath.Join(repo.VaultPath(), filename)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return err
		}
	}

	// Create demo notes
	for _, n := range demoData.Notes {
		fm := obsidian.Frontmatter{
			Title:     n.Title,
			Type:      "note",
			Tag:       n.Tag,
			CreatedAt: now,
		}
		content, err := frontmatter.Serialize(fm, n.Body)
		if err != nil {
			return err
		}

		filename := obsidian.FilenameFromTitle(n.Title)
		filePath := filepath.Join(repo.VaultPath(), filename)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return err
		}
	}

	return nil
}
