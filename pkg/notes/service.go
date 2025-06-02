package notes

type Service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return Service{repository: r}
}

func (s Service) FindAll() ([]Note, error) {
	return s.repository.FindAll()
}
