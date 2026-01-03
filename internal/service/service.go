package service

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/markelca/prioritty/internal/editor"
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

func (s Service) RemoveItem(item items.ItemInterface) error {
	switch v := item.(type) {
	case *items.Note:
		return s.removeNote(v.GetId())
	case *items.Task:
		return s.removeTask(v.GetId())
	default:
		return fmt.Errorf("Cannot remove item %v", v)
	}
}

func (s Service) UpdateItemFromEditorMsg(i items.ItemInterface, msg editor.EditorFinishedMsg) error {
	switch v := i.(type) {
	case *items.Task:
		v.Title = msg.Title
		v.Body = msg.Body
		// Update status if provided
		if msg.Status != "" {
			v.Status = items.ParseStatus(msg.Status)
		}
		if err := s.UpdateTask(*v); err != nil {
			log.Println("Error updating the task - ", err)
		}
		// Update tag if changed
		if err := s.updateTagFromEditor(v, msg.Tag); err != nil {
			log.Println("Error updating tag - ", err)
		}
	case *items.Note:
		v.Title = msg.Title
		v.Body = msg.Body
		if err := s.UpdateNote(*v); err != nil {
			log.Println("Error updating the note - ", err)
		}
		// Update tag if changed
		if err := s.updateTagFromEditor(v, msg.Tag); err != nil {
			log.Println("Error updating tag - ", err)
		}
	default:
		return fmt.Errorf("Can't update the item, no implementation: %v", v)
	}
	return nil
}

// updateTagFromEditor updates an item's tag based on the editor message.
func (s Service) updateTagFromEditor(i items.ItemInterface, tagName string) error {
	currentTag := i.GetTag()
	currentTagName := ""
	if currentTag != nil {
		currentTagName = currentTag.Name
	}

	// No change needed
	if currentTagName == tagName {
		return nil
	}

	// Remove tag if empty
	if tagName == "" {
		return s.UnsetTag(i)
	}

	// Set new tag
	return s.SetTag(i, tagName)
}

func (s Service) CreateTaskFromEditorMsg(msg editor.EditorFinishedMsg) error {
	task := items.Task{
		Item: items.Item{
			Title: msg.Title,
			Body:  msg.Body,
		},
		Status: items.ParseStatus(msg.Status),
	}
	if err := s.repository.CreateTask(task); err != nil {
		return err
	}
	// Set tag if provided
	if msg.Tag != "" {
		// Need to get the created task to set the tag
		// For now, we'll set it directly since the task ID will be the title
		return s.SetTag(&task, msg.Tag)
	}
	return nil
}

func (s Service) CreateNoteFromEditorMsg(msg editor.EditorFinishedMsg) error {
	note := items.Note{
		Item: items.Item{
			Title: msg.Title,
			Body:  msg.Body,
		},
	}
	if err := s.repository.CreateNote(note); err != nil {
		return err
	}
	// Set tag if provided
	if msg.Tag != "" {
		return s.SetTag(&note, msg.Tag)
	}
	return nil
}

func (s Service) SetTag(i items.ItemInterface, name string) error {
	var tag *items.Tag
	var err error
	tag, err = s.repository.GetTag(name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			tag, err = s.repository.CreateTag(name)
			if err != nil {
				log.Printf("Error creating tag: %v", err)
				return err
			}
		} else {
			log.Printf("Error getting tag: %v", err)
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

func (s Service) GetTags() ([]items.Tag, error) {
	return s.repository.GetTags()
}

func (s Service) RemoveTag(name string) error {
	itemsWithTag, err := s.repository.GetItemsWithTag(name)
	if err != nil {
		return fmt.Errorf("error checking items with tag: %v", err)
	}

	if len(itemsWithTag) > 0 {
		return fmt.Errorf("cannot remove tag '%s' because it is assigned to %d item(s)", name, len(itemsWithTag))
	}

	return s.repository.RemoveTag(name)
}

func (s Service) GetItemsWithTag(name string) ([]items.ItemInterface, error) {
	return s.repository.GetItemsWithTag(name)
}
