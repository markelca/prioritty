package editor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/config"
	"github.com/markelca/prioritty/pkg/frontmatter"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// unquotedString is a string type that marshals to YAML without quotes.
type unquotedString string

func (s unquotedString) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: string(s),
	}, nil
}

// TaskEditorFrontmatter represents the YAML frontmatter for tasks in the editor.
type TaskEditorFrontmatter struct {
	Title  unquotedString `yaml:"title"`
	Status unquotedString `yaml:"status"`
	Tag    unquotedString `yaml:"tag"`
}

// NoteEditorFrontmatter represents the YAML frontmatter for notes in the editor.
type NoteEditorFrontmatter struct {
	Title unquotedString `yaml:"title"`
	Tag   unquotedString `yaml:"tag"`
}

// EditorInput contains the data to populate the editor temp file.
type EditorInput struct {
	Id       string
	ItemType items.ItemType
	Title    string
	Body     string
	Status   string
	Tag      string
}

// EditorFinishedMsg contains the parsed result from the editor.
type EditorFinishedMsg struct {
	Id     string
	Title  string
	Body   string
	Status string
	Tag    string
	Err    error
}

func EditItem(input EditorInput) (tea.Cmd, error) {
	tempFile, err := os.CreateTemp(os.TempDir(), "item_*.md")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	content, err := serializeEditorContent(input)
	if err != nil {
		tempFile.Close()
		return nil, fmt.Errorf("failed to serialize content: %w", err)
	}

	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}
	tempFile.Close()

	editor, err := getEditor()
	if err != nil {
		return nil, err
	}

	return tea.ExecProcess(exec.Command(editor, tempFile.Name()), func(err error) tea.Msg {
		defer os.Remove(tempFile.Name())
		modifiedContent, err := os.ReadFile(tempFile.Name())
		if err != nil {
			return EditorFinishedMsg{Err: fmt.Errorf("failed to read modified file: %w", err)}
		}

		msg := parseEditorContent(string(modifiedContent), input.ItemType)
		msg.Id = input.Id

		return msg
	}), nil
}

// serializeEditorContent creates the content for the editor temp file with frontmatter.
func serializeEditorContent(input EditorInput) (string, error) {
	var content []byte
	var err error

	if input.ItemType == items.ItemTypeTask {
		fm := TaskEditorFrontmatter{
			Title:  unquotedString(input.Title),
			Status: unquotedString(input.Status),
			Tag:    unquotedString(input.Tag),
		}
		content, err = frontmatter.Serialize(fm, input.Body)
	} else {
		fm := NoteEditorFrontmatter{
			Title: unquotedString(input.Title),
			Tag:   unquotedString(input.Tag),
		}
		content, err = frontmatter.Serialize(fm, input.Body)
	}

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// parsedFrontmatter is used for parsing (uses regular strings)
type parsedFrontmatter struct {
	Title  string `yaml:"title"`
	Status string `yaml:"status"`
	Tag    string `yaml:"tag"`
}

// parseEditorContent parses the editor content including frontmatter.
func parseEditorContent(content string, itemType items.ItemType) EditorFinishedMsg {
	// Check if content is completely empty or only whitespace
	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return EditorFinishedMsg{Err: fmt.Errorf("operation cancelled - no content provided")}
	}

	// Parse frontmatter
	var fm parsedFrontmatter
	body, err := frontmatter.Parse(content, &fm)
	if err != nil {
		log.Printf("Error parsing frontmatter: %v", err)
		return EditorFinishedMsg{Err: fmt.Errorf("invalid frontmatter: %w", err)}
	}

	// Title is required
	title := strings.TrimSpace(fm.Title)
	if title == "" {
		return EditorFinishedMsg{Err: fmt.Errorf("operation cancelled - no title provided")}
	}

	return EditorFinishedMsg{
		Title:  title,
		Body:   strings.TrimSpace(body),
		Status: fm.Status,
		Tag:    fm.Tag,
	}
}

func getEditor() (string, error) {
	if editor := viper.GetString(config.CONF_EDITOR); editor != "" {
		if _, err := exec.LookPath(editor); err == nil {
			return editor, nil
		} else {
			log.Printf("Warning - Couldn't open the editor from your prioritty.yaml file (%s), make sure is installed and accesible globally", editor)
		}
	}

	// Check environment variables in order of preference
	editors := []string{"VISUAL", "EDITOR"}

	for _, env := range editors {
		if editor := os.Getenv(env); editor != "" {
			if _, err := exec.LookPath(editor); err == nil {
				return editor, nil
			} else {
				log.Printf("Warning - Couldn't open the editor from your environment vars (%s), make sure is installed and accesible globally", editor)
			}
		}
	}

	switch runtime.GOOS {
	case "windows":
		return "notepad", nil
	default:
		commonEditors := []string{"nano", "vim", "vi", "emacs"}
		for _, editor := range commonEditors {
			if _, err := exec.LookPath(editor); err == nil {
				return editor, nil
			}
		}
	}
	return "", fmt.Errorf("Error - No available editor could be found")
}
