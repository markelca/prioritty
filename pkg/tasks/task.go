package tasks

type Status int

const (
	Todo Status = iota
	InProgress
	Done
	Cancelled
)

type Task struct {
	Id     int
	Title  string
	Status Status
}
