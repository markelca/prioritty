package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNote_ImplementsItemInterface(t *testing.T) {
	var _ ItemInterface = (*Note)(nil)
}

func TestNote_Render(t *testing.T) {
	note := Note{
		Item: Item{
			Id:    "note-1",
			Title: "Test Note",
			Body:  "Note body",
		},
	}

	mockRenderer := &mockRenderer{output: "rendered note"}
	result := note.Render(mockRenderer)

	assert.Equal(t, "rendered note", result)
	assert.True(t, mockRenderer.called)
}

func TestTask_Render(t *testing.T) {
	task := Task{
		Item: Item{
			Id:    "task-1",
			Title: "Test Task",
			Body:  "Task body",
		},
		Status: Done,
	}

	mockRenderer := &mockRenderer{output: "rendered task"}
	result := task.Render(mockRenderer)

	assert.Equal(t, "rendered task", result)
	assert.True(t, mockRenderer.called)
}

// mockRenderer is a test double for Renderer
type mockRenderer struct {
	output string
	called bool
}

func (m *mockRenderer) Render(r Renderable) string {
	m.called = true
	return m.output
}
