package tasks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/markelca/prioritty/pkg/editor"
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

func (s Service) EditWithEditor(t *Task) (tea.Cmd, error) {
	return editor.EditTask(t.Id, t.Title, t.Body)
}
