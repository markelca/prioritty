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
