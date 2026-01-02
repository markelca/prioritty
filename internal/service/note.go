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

func (s NoteService) UpdateNote(n items.Note) error {
	return s.repository.UpdateNote(n)
}

func (s NoteService) AddNote(title string) error {
	t := items.Note{}
	t.Title = title
	return s.repository.CreateNote(t)
}

func (s NoteService) removeNote(id string) error {
	return s.repository.RemoveNote(id)
}
