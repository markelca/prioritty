package service

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/internal/editor"
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
)

type TaskService struct {
	repository repository.Repository
}

func (s TaskService) GetTasks() ([]items.Task, error) {
	return s.repository.GetTasks()
}

func (s TaskService) DestroyDemo() error {
	return s.repository.Reset()
}

func (s TaskService) UpdateTask(t items.Task) error {
	return s.repository.UpdateTask(t)
}

func (s TaskService) UpdateStatus(t *items.Task, status items.Status) error {
	if t.Status == status {
		status = items.Todo
	}
	err := s.repository.UpdateTaskStatus(*t, status)
	if err != nil {
		return err
	}
	t.SetStatus(status)
	return nil
}

func (s TaskService) SetTag(title string) error {
	return nil
}

func (s TaskService) TagItem(items items.ItemInterface, tag string) error {
	return nil
}

func (s TaskService) AddTask(title string) error {
	t := items.Task{}
	t.Title = title
	return s.repository.CreateTask(t)
}

func (s TaskService) removeTask(id string) error {
	return s.repository.RemoveTask(id)
}

func (s TaskService) EditWithEditor(t items.ItemInterface) (tea.Cmd, error) {
	input := editor.EditorInput{
		Id:    t.GetId(),
		Title: t.GetTitle(),
		Body:  t.GetBody(),
	}

	// Set item type and status based on the concrete type
	switch task := t.(type) {
	case *items.Task:
		input.ItemType = items.ItemTypeTask
		input.Status = items.StatusToString(task.Status)
	case *items.Note:
		input.ItemType = items.ItemTypeNote
	}

	// Set tag if present
	if tag := t.GetTag(); tag != nil {
		input.Tag = tag.Name
	}

	return editor.EditItem(input)
}

func (s TaskService) CreateWithEditor(itemType items.ItemType) (tea.Cmd, error) {
	input := editor.EditorInput{
		ItemType: itemType,
	}
	if itemType == items.ItemTypeTask {
		input.Status = "todo"
	}
	return editor.EditItem(input)
}
