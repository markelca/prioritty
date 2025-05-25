package tasks

type Service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return Service{repository: r}
}

func (s Service) FindAll() ([]Task, error) {
	return s.repository.FindAll()
}

func (s Service) UpdateStatus(t *Task, status Status) error {
	err := s.repository.UpdateStatus(*t, status)
	if err != nil {
		return err
	}
	t.SetStatus(status)
	return nil
}

func (s Service) AddTask(title string) error {
	t := Task{Title: title}
	return s.repository.CreateTask(t)
}

func (s Service) RemoveTask(id int) error {
	return s.repository.RemoveTask(id)
}
