package service

import (
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
)

type Service struct {
	TaskService
	NoteService
	repository repository.Repository
}

func NewService(r repository.Repository) Service {
	return Service{
		TaskService: TaskService{repository: r},
		NoteService: NoteService{repository: r},
		repository:  r,
	}
}

func (s Service) GetAll() ([]items.Item, error) {
	var allItems []items.Item

	// Get all notes
	notes, err := s.GetNotes()
	if err != nil {
		return nil, err
	}

	// Convert notes to items and append
	for _, note := range notes {
		allItems = append(allItems, note.Item)
	}

	// Get all tasks
	tasks, err := s.GetTasks()
	if err != nil {
		return nil, err
	}

	// Convert tasks to items and append
	for _, task := range tasks {
		allItems = append(allItems, task.Item)
	}

	return allItems, nil
}
