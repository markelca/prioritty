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
	"github.com/spf13/viper"
)

type TaskEditorFinishedMsg struct {
	Id    int
	Title string
	Body  string
	Err   error
}

func EditTask(id int, title, body string) (tea.Cmd, error) {
	tempFile, err := os.CreateTemp(os.TempDir(), "task_*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	content := title + "\n\n" + body
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
			return TaskEditorFinishedMsg{Err: fmt.Errorf("failed to read modified file: %w", err)}
		}

		msg := parseTaskContent(string(modifiedContent))
		msg.Id = id

		return msg
	}), nil
}

func parseTaskContent(content string) TaskEditorFinishedMsg {
	// Check if content is completely empty or only whitespace
	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" {
		return TaskEditorFinishedMsg{Err: fmt.Errorf("operation cancelled - no content provided")}
	}

	lines := strings.Split(trimmedContent, "\n")
	var body, title string

	if len(lines) == 0 {
		return TaskEditorFinishedMsg{Err: fmt.Errorf("operation cancelled - no content provided")}
	}

	title = strings.TrimSpace(lines[0])
	if title == "" {
		return TaskEditorFinishedMsg{Err: fmt.Errorf("operation cancelled - no content provided")}
	}

	// Rest is the content (skip empty line if present)
	var contentLines []string
	startIndex := 1

	// Skip the first empty line if it exists
	if len(lines) > 1 && strings.TrimSpace(lines[1]) == "" {
		startIndex = 2
	}

	if len(lines) > startIndex {
		contentLines = lines[startIndex:]
		body = strings.Join(contentLines, "\n")
	} else {
		body = ""
	}

	return TaskEditorFinishedMsg{Title: title, Body: body, Err: nil}
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
