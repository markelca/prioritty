package service

import (
	"sort"

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

func (s Service) GetAll() ([]items.ItemInterface, error) {
	var allItems []items.ItemInterface

	notes, err := s.GetNotes()
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		allItems = append(allItems, &note)
	}

	tasks, err := s.GetTasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		allItems = append(allItems, &task)
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].GetCreatedAt().After(allItems[j].GetCreatedAt())
	})

	return allItems, nil
}
