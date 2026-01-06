package items

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Status
	}{
		{"todo lowercase", "todo", Todo},
		{"todo uppercase", "TODO", Todo},
		{"todo mixed case", "ToDo", Todo},
		{"in-progress with hyphen", "in-progress", InProgress},
		{"inprogress without hyphen", "inprogress", InProgress},
		{"InProgress mixed case", "InProgress", InProgress},
		{"done lowercase", "done", Done},
		{"done uppercase", "DONE", Done},
		{"cancelled UK spelling", "cancelled", Cancelled},
		{"canceled US spelling", "canceled", Cancelled},
		{"invalid defaults to todo", "invalid", Todo},
		{"empty defaults to todo", "", Todo},
		{"random string defaults to todo", "something", Todo},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseStatus(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTask_SetStatus(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  Status
		newStatus      Status
		expectedStatus Status
	}{
		{"todo to done", Todo, Done, Done},
		{"todo to in-progress", Todo, InProgress, InProgress},
		{"todo to cancelled", Todo, Cancelled, Cancelled},
		{"done to todo", Done, Todo, Todo},
		{"in-progress to done", InProgress, Done, Done},
		{"toggle todo to todo resets to todo", Todo, Todo, Todo},
		{"toggle done to done resets to todo", Done, Done, Todo},
		{"toggle in-progress to in-progress resets to todo", InProgress, InProgress, Todo},
		{"toggle cancelled to cancelled resets to todo", Cancelled, Cancelled, Todo},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{Status: tc.initialStatus}
			task.SetStatus(tc.newStatus)
			assert.Equal(t, tc.expectedStatus, task.Status)
		})
	}
}

func TestTask_SetStatus_ToggleBehavior(t *testing.T) {
	task := &Task{Status: Todo}

	// Set to Done
	task.SetStatus(Done)
	assert.Equal(t, Done, task.Status, "Should change from Todo to Done")

	// Toggle Done (should reset to Todo)
	task.SetStatus(Done)
	assert.Equal(t, Todo, task.Status, "Toggling same status should reset to Todo")

	// Set to InProgress
	task.SetStatus(InProgress)
	assert.Equal(t, InProgress, task.Status, "Should change to InProgress")

	// Change from InProgress to Done (not toggle)
	task.SetStatus(Done)
	assert.Equal(t, Done, task.Status, "Should change from InProgress to Done")
}

func TestStatusConstants(t *testing.T) {
	assert.Equal(t, Status("todo"), Todo)
	assert.Equal(t, Status("in-progress"), InProgress)
	assert.Equal(t, Status("done"), Done)
	assert.Equal(t, Status("cancelled"), Cancelled)
	assert.Equal(t, Status("note"), NoteType)
}

func TestTask_ImplementsItemInterface(t *testing.T) {
	var _ ItemInterface = (*Task)(nil)
}
