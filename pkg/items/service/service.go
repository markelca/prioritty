package service

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/markelca/prioritty/pkg/editor"
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
		return allItems[i].After(allItems[j])
	})

	return allItems, nil
}

func (s Service) UpdateItemFromEditorMsg(i items.ItemInterface, msg editor.TaskEditorFinishedMsg) error {
	switch v := i.(type) {
	case *items.Task:
		v.Title = msg.Title
		v.Body = msg.Body
		if err := s.UpdateTask(*v); err != nil {
			log.Println("Error updating the task - ", err)
		}
	case *items.Note:
		v.Title = msg.Title
		v.Body = msg.Body
		if err := s.UpdateNote(*v); err != nil {
			log.Println("Error updating the task - ", err)
		}
	default:
		return fmt.Errorf("Can't update the item, no implementation: %v", v)
	}
	return nil
}

func (s Service) SetTag(i items.ItemInterface, name string) error {
	var tag *items.Tag
	var err error
	tag, err = s.repository.GetTag(name)
	if err != nil {
		if err == sql.ErrNoRows {
			tag, err = s.repository.CreateTag(name)
			if err != nil {
				log.Printf("Error creating tag: %v", err)
				return err
			}
		} else {
			log.Printf("Error geting tag: %v", err)
			return err
		}
	}
	switch v := i.(type) {
	case *items.Task:
		return s.repository.SetTaskTag(*v, *tag)
	case *items.Note:
		return s.repository.SetNoteTag(*v, *tag)
	default:
		return fmt.Errorf("Can't update the item, no implementation: %v", v)
	}

}

func (s Service) UnsetTag(i items.ItemInterface) error {
	switch v := i.(type) {
	case *items.Task:
		return s.repository.UnsetTaskTag(*v)
	case *items.Note:
		return s.repository.UnsetNoteTag(*v)
	default:
		return fmt.Errorf("Can't unset tag for item, no implementation: %v", v)
	}
}
