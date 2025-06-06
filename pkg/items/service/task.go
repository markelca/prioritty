package service

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/pkg/editor"
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
	return s.repository.DropSchema()
}

func (s TaskService) UpdateTask(t items.Task) error {
	return s.repository.UpdateTask(t)
}

func (s TaskService) UpdateTaskFromEditorMsg(t *items.Task, msg editor.TaskEditorFinishedMsg) {
	t.Title = msg.Title
	t.Body = msg.Body
	if err := s.UpdateTask(*t); err != nil {
		fmt.Println("Error updating the task - ", err)
	}
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

func (s TaskService) AddTask(title string) error {
	t := items.Task{}
	t.Title = title
	return s.repository.CreateTask(t)
}

func (s TaskService) AddNote(title string) error {
	t := items.Note{}
	t.Title = title
	return s.repository.CreateNote(t)
}

func (s TaskService) RemoveTask(id int) error {
	return s.repository.RemoveTask(id)
}

func (s TaskService) EditWithEditor(t *items.Task) (tea.Cmd, error) {
	return editor.EditTask(t.Id, t.Title, t.Body)
}
