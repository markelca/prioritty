package testutils

import (
	"slices"

	"github.com/markelca/prioritty/pkg/items"
	"github.com/markelca/prioritty/pkg/items/repository"
)

// Ensure MockRepository implements Repository interface
var _ repository.Repository = (*MockRepository)(nil)

// MockRepository is a test double for Repository that stores data in memory
type MockRepository struct {
	// Storage
	tasks map[string]items.Task
	notes map[string]items.Note
	tags  map[string]items.Tag

	// Call tracking for assertions
	Calls []string

	// Error injection for testing error paths
	GetTasksError         error
	GetNotesError         error
	CreateTaskError       error
	CreateNoteError       error
	UpdateTaskError       error
	UpdateNoteError       error
	RemoveTaskError       error
	RemoveNoteError       error
	UpdateTaskStatusError error
	SetTaskTagError       error
	UnsetTaskTagError     error
	SetNoteTagError       error
	UnsetNoteTagError     error
	GetTagError           error
	GetTagsError          error
	CreateTagError        error
	RemoveTagError        error
	GetItemsWithTagError  error
	ResetError            error
}

// NewMockRepository creates a new MockRepository with initialized storage
func NewMockRepository() *MockRepository {
	return &MockRepository{
		tasks: make(map[string]items.Task),
		notes: make(map[string]items.Note),
		tags:  make(map[string]items.Tag),
		Calls: []string{},
	}
}

// recordCall adds a method call to the call tracking list
func (m *MockRepository) recordCall(method string) {
	m.Calls = append(m.Calls, method)
}

// HasCall checks if a method was called
func (m *MockRepository) HasCall(method string) bool {
	return slices.Contains(m.Calls, method)
}

// CallCount returns the number of times a method was called
func (m *MockRepository) CallCount(method string) int {
	count := 0
	for _, call := range m.Calls {
		if call == method {
			count++
		}
	}
	return count
}

// ResetCalls clears the call tracking list
func (m *MockRepository) ResetCalls() {
	m.Calls = []string{}
}

// --- Task Repository Methods ---

func (m *MockRepository) GetTasks() ([]items.Task, error) {
	m.recordCall("GetTasks")
	if m.GetTasksError != nil {
		return nil, m.GetTasksError
	}

	result := make([]items.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		result = append(result, task)
	}
	return result, nil
}

func (m *MockRepository) UpdateTask(task items.Task) error {
	m.recordCall("UpdateTask")
	if m.UpdateTaskError != nil {
		return m.UpdateTaskError
	}

	if _, exists := m.tasks[task.Id]; !exists {
		return repository.ErrNotFound
	}
	m.tasks[task.Id] = task
	return nil
}

func (m *MockRepository) CreateTask(task *items.Task) error {
	m.recordCall("CreateTask")
	if m.CreateTaskError != nil {
		return m.CreateTaskError
	}

	m.tasks[task.Id] = *task
	return nil
}

func (m *MockRepository) RemoveTask(id string) error {
	m.recordCall("RemoveTask")
	if m.RemoveTaskError != nil {
		return m.RemoveTaskError
	}

	if _, exists := m.tasks[id]; !exists {
		return repository.ErrNotFound
	}
	delete(m.tasks, id)
	return nil
}

func (m *MockRepository) UpdateTaskStatus(task items.Task, status items.Status) error {
	m.recordCall("UpdateTaskStatus")
	if m.UpdateTaskStatusError != nil {
		return m.UpdateTaskStatusError
	}

	if existing, exists := m.tasks[task.Id]; exists {
		existing.Status = status
		m.tasks[task.Id] = existing
		return nil
	}
	return repository.ErrNotFound
}

func (m *MockRepository) SetTaskTag(task items.Task, tag items.Tag) error {
	m.recordCall("SetTaskTag")
	if m.SetTaskTagError != nil {
		return m.SetTaskTagError
	}

	if existing, exists := m.tasks[task.Id]; exists {
		existing.Tag = &tag
		m.tasks[task.Id] = existing
		return nil
	}
	return repository.ErrNotFound
}

func (m *MockRepository) UnsetTaskTag(task items.Task) error {
	m.recordCall("UnsetTaskTag")
	if m.UnsetTaskTagError != nil {
		return m.UnsetTaskTagError
	}

	if existing, exists := m.tasks[task.Id]; exists {
		existing.Tag = nil
		m.tasks[task.Id] = existing
		return nil
	}
	return repository.ErrNotFound
}

// --- Note Repository Methods ---

func (m *MockRepository) GetNotes() ([]items.Note, error) {
	m.recordCall("GetNotes")
	if m.GetNotesError != nil {
		return nil, m.GetNotesError
	}

	result := make([]items.Note, 0, len(m.notes))
	for _, note := range m.notes {
		result = append(result, note)
	}
	return result, nil
}

func (m *MockRepository) UpdateNote(note items.Note) error {
	m.recordCall("UpdateNote")
	if m.UpdateNoteError != nil {
		return m.UpdateNoteError
	}

	if _, exists := m.notes[note.Id]; !exists {
		return repository.ErrNotFound
	}
	m.notes[note.Id] = note
	return nil
}

func (m *MockRepository) CreateNote(note *items.Note) error {
	m.recordCall("CreateNote")
	if m.CreateNoteError != nil {
		return m.CreateNoteError
	}

	m.notes[note.Id] = *note
	return nil
}

func (m *MockRepository) RemoveNote(id string) error {
	m.recordCall("RemoveNote")
	if m.RemoveNoteError != nil {
		return m.RemoveNoteError
	}

	if _, exists := m.notes[id]; !exists {
		return repository.ErrNotFound
	}
	delete(m.notes, id)
	return nil
}

func (m *MockRepository) SetNoteTag(note items.Note, tag items.Tag) error {
	m.recordCall("SetNoteTag")
	if m.SetNoteTagError != nil {
		return m.SetNoteTagError
	}

	if existing, exists := m.notes[note.Id]; exists {
		existing.Tag = &tag
		m.notes[note.Id] = existing
		return nil
	}
	return repository.ErrNotFound
}

func (m *MockRepository) UnsetNoteTag(note items.Note) error {
	m.recordCall("UnsetNoteTag")
	if m.UnsetNoteTagError != nil {
		return m.UnsetNoteTagError
	}

	if existing, exists := m.notes[note.Id]; exists {
		existing.Tag = nil
		m.notes[note.Id] = existing
		return nil
	}
	return repository.ErrNotFound
}

// --- Tag Methods ---

func (m *MockRepository) GetTag(name string) (*items.Tag, error) {
	m.recordCall("GetTag")
	if m.GetTagError != nil {
		return nil, m.GetTagError
	}

	for _, tag := range m.tags {
		if tag.Name == name {
			return &tag, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *MockRepository) GetTags() ([]items.Tag, error) {
	m.recordCall("GetTags")
	if m.GetTagsError != nil {
		return nil, m.GetTagsError
	}

	result := make([]items.Tag, 0, len(m.tags))
	for _, tag := range m.tags {
		result = append(result, tag)
	}
	return result, nil
}

func (m *MockRepository) CreateTag(name string) (*items.Tag, error) {
	m.recordCall("CreateTag")
	if m.CreateTagError != nil {
		return nil, m.CreateTagError
	}

	tag := items.Tag{
		Id:   "mock-tag-" + name,
		Name: name,
	}
	m.tags[tag.Id] = tag
	return &tag, nil
}

func (m *MockRepository) RemoveTag(name string) error {
	m.recordCall("RemoveTag")
	if m.RemoveTagError != nil {
		return m.RemoveTagError
	}

	for id, tag := range m.tags {
		if tag.Name == name {
			delete(m.tags, id)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (m *MockRepository) GetItemsWithTag(tagName string) ([]items.ItemInterface, error) {
	m.recordCall("GetItemsWithTag")
	if m.GetItemsWithTagError != nil {
		return nil, m.GetItemsWithTagError
	}

	var result []items.ItemInterface

	for _, task := range m.tasks {
		if task.Tag != nil && task.Tag.Name == tagName {
			t := task // Create a copy to avoid pointer issues
			result = append(result, &t)
		}
	}

	for _, note := range m.notes {
		if note.Tag != nil && note.Tag.Name == tagName {
			n := note // Create a copy to avoid pointer issues
			result = append(result, &n)
		}
	}

	return result, nil
}

func (m *MockRepository) Reset() error {
	m.recordCall("Reset")
	if m.ResetError != nil {
		return m.ResetError
	}

	m.tasks = make(map[string]items.Task)
	m.notes = make(map[string]items.Note)
	m.tags = make(map[string]items.Tag)
	return nil
}

// --- Test Helper Methods ---

// AddTask adds a task directly to the mock storage (bypasses CreateTask)
func (m *MockRepository) AddTask(task items.Task) {
	m.tasks[task.Id] = task
}

// AddNote adds a note directly to the mock storage (bypasses CreateNote)
func (m *MockRepository) AddNote(note items.Note) {
	m.notes[note.Id] = note
}

// AddTag adds a tag directly to the mock storage (bypasses CreateTag)
func (m *MockRepository) AddTag(tag items.Tag) {
	m.tags[tag.Id] = tag
}

// GetTaskById retrieves a task by ID from mock storage
func (m *MockRepository) GetTaskById(id string) (*items.Task, bool) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, false
	}
	return &task, true
}

// GetNoteById retrieves a note by ID from mock storage
func (m *MockRepository) GetNoteById(id string) (*items.Note, bool) {
	note, exists := m.notes[id]
	if !exists {
		return nil, false
	}
	return &note, true
}

// TaskCount returns the number of tasks in storage
func (m *MockRepository) TaskCount() int {
	return len(m.tasks)
}

// NoteCount returns the number of notes in storage
func (m *MockRepository) NoteCount() int {
	return len(m.notes)
}

// TagCount returns the number of tags in storage
func (m *MockRepository) TagCount() int {
	return len(m.tags)
}
