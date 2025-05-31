package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return Service{repository: r}
}

func (s Service) FindAll() ([]Task, error) {
	return s.repository.FindAll()
}

func (s Service) DestroyDemo() error {
	return s.repository.DropSchema()
}

func (s Service) UpdateTask(t Task) error {
	return s.repository.UpdateTask(t)
}

func (s Service) UpdateStatus(t *Task, status Status) error {
	if t.Status == status {
		status = Todo
	}
	err := s.repository.UpdateStatus(*t, status)
	if err != nil {
		return err
	}
	t.SetStatus(status)
	return nil
}

func (s Service) AddTask(title string) error {
	t := Task{Title: title}
	return s.repository.CreateTask(t)
}

func (s Service) RemoveTask(id int) error {
	return s.repository.RemoveTask(id)
}

func (s Service) EditWithEditor(t *Task) error {
	// Create temporary file with current task content
	tempFile, err := ioutil.TempFile("", "task_*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write current task content to temp file
	content := t.Title + "\n\n" + t.Body
	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	tempFile.Close()

	// Get the editor command
	editor := getEditor()

	// Open the editor
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run editor: %w", err)
	}

	// Read the modified content
	modifiedContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read modified file: %w", err)
	}

	// Parse the content back into title and body
	if err := t.parseContent(string(modifiedContent)); err != nil {
		return fmt.Errorf("failed to parse content: %w", err)
	}

	if err := s.UpdateTask(*t); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// parseContent parses the editor content into title and body
func (t *Task) parseContent(content string) error {
	lines := strings.Split(strings.TrimSpace(content), "\n")

	if len(lines) == 0 {
		return fmt.Errorf("content cannot be empty")
	}

	// First line is the title
	t.Title = strings.TrimSpace(lines[0])
	if t.Title == "" {
		return fmt.Errorf("title cannot be empty")
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
		t.Body = strings.Join(contentLines, "\n")
	} else {
		t.Body = ""
	}

	return nil
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

// CreateTaskFromEditor creates a new task using the editor
func CreateTaskFromEditor() (*Task, error) {
	task := &Task{}

	// Create temp file with template
	tempFile, err := ioutil.TempFile("", "new_task_*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write template to temp file
	template := "# Enter task title here\n\n# Enter task description here (optional)\n# Lines starting with # are comments and will be ignored"
	if _, err := tempFile.WriteString(template); err != nil {
		tempFile.Close()
		return nil, fmt.Errorf("failed to write template: %w", err)
	}
	tempFile.Close()

	// Open editor
	editor := getEditor()
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run editor: %w", err)
	}

	// Read and parse content
	content, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse content, ignoring comment lines
	if err := task.parseContentIgnoringComments(string(content)); err != nil {
		return nil, err
	}

	return task, nil
}

// parseContentIgnoringComments parses content while ignoring comment lines
func (t *Task) parseContentIgnoringComments(content string) error {
	lines := strings.Split(content, "\n")
	var validLines []string

	// Filter out comment lines and empty lines
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			validLines = append(validLines, line)
		}
	}

	if len(validLines) == 0 {
		return fmt.Errorf("no valid content found")
	}

	// First non-comment line is title
	t.Title = strings.TrimSpace(validLines[0])
	if t.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	// Rest is content
	if len(validLines) > 1 {
		t.Body = strings.Join(validLines[1:], "\n")
	}

	return nil
}
