package tasks

type Status int

const (
	Todo Status = iota
	InProgress
	Blocked
	Done
)

type Task struct {
	Title  string
	Status Status
}
