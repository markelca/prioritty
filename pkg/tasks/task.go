package tasks

type Status int

const (
	Todo Status = iota
	InProgress
	Blocked
	Done
)

type Task struct {
	Id     int
	Title  string
	Status Status
}
