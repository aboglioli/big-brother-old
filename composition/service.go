package composition

type Service interface {
	Create(*Composition) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(c *Composition) error {
	return s.repository.Insert(c)
}
