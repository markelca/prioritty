package repository

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/spf13/viper"
)

var ErrNotFound = errors.New("not found")

const (
	RepoTypeObsidian = "obsidian"
	RepoTypeSQLite   = "sqlite"
)

type TaskRepository interface {
	GetTasks() ([]items.Task, error)
	UpdateTask(items.Task) error
	CreateTask(items.Task) error
	RemoveTask(string) error
	UpdateTaskStatus(items.Task, items.Status) error
	SetTaskTag(items.Task, items.Tag) error
	UnsetTaskTag(items.Task) error
}

type NoteRepository interface {
	GetNotes() ([]items.Note, error)
	UpdateNote(items.Note) error
	CreateNote(items.Note) error
	RemoveNote(string) error
	SetNoteTag(items.Note, items.Tag) error
	UnsetNoteTag(items.Note) error
}

type Repository interface {
	TaskRepository
	NoteRepository
	GetTag(string) (*items.Tag, error)
	GetTags() ([]items.Tag, error)
	CreateTag(string) (*items.Tag, error)
	RemoveTag(string) error
	GetItemsWithTag(string) ([]items.ItemInterface, error)
	Reset() error
}

func expandTilde(p string) string {
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		return path.Join(home, p[2:])
	}
	return p
}

func GetDatabasePath(repoType string, isDemo bool) (dbPath string, err error) {
	switch repoType {
	case RepoTypeObsidian:
		if isDemo {
			dbPath = path.Join(os.TempDir(), "prioritty_demo_vault")
		} else {
			dbPath = expandTilde(viper.GetString(config.CONF_DATABASE_PATH))
		}
	case RepoTypeSQLite:
		if isDemo {
			dbPath = path.Join(os.TempDir(), "prioritty_demo.db")
		} else {
			dbPath = expandTilde(viper.GetString(config.CONF_DATABASE_PATH))
		}
	default:
		return "", fmt.Errorf("repository type not supported (%s)", repoType)
	}
	return dbPath, nil
}
