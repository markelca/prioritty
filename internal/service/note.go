package service

import (
	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
)

type NoteService struct {
	repository repository.Repository
}

func (s NoteService) GetNotes() ([]items.Note, error) {
	return s.repository.GetNotes()
}
