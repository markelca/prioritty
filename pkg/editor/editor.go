package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

	editor := getEditor()

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
	lines := strings.Split(strings.TrimSpace(content), "\n")
	var body, title string

	if len(lines) == 0 {
		return TaskEditorFinishedMsg{Err: fmt.Errorf("content cannot be empty")}
	}

	title = strings.TrimSpace(lines[0])
	if title == "" {
		return TaskEditorFinishedMsg{Err: fmt.Errorf("title cannot be empty")}
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

// getEditor returns the preferred editor command
func getEditor() string {
	// Check environment variables in order of preference
	editors := []string{"VISUAL", "EDITOR"}

	for _, env := range editors {
		if editor := os.Getenv(env); editor != "" {
			return editor
		}
	}

	// Fall back to system defaults
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	default:
		// Try common editors in order of preference
		commonEditors := []string{"vim", "vi", "nano", "emacs"}
		for _, editor := range commonEditors {
			if _, err := exec.LookPath(editor); err == nil {
				return editor
			}
		}
		// Ultimate fallback
		return "vi"
	}
}
