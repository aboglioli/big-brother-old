package quantity

import (
	"errors"

	"github.com/aboglioli/big-brother/unit"
)

type Service interface {
	IsValid(*Quantity) bool
	Add(*Quantity, *Quantity) (*Quantity, error)
	Substract(*Quantity, *Quantity) (*Quantity, error)
}

type service struct {
	unitRepository unit.Repository
}

func NewService() Service {
	return &service{
		unitRepository: unit.NewRepository(),
	}
}

func (s *service) IsValid(q *Quantity) bool {
	if s.unitRepository.Exists(q.Unit) && q.Quantity >= 0 {
		return true
	}
	return false
}

func (s *service) Add(q1 *Quantity, q2 *Quantity) (*Quantity, error) {
	u1 := s.unitRepository.FindByName(q1.Unit)
	if u1 == nil {
		return nil, errors.New("Unit does not exist")
	}

	u2 := s.unitRepository.FindByName(q2.Unit)
	if u2 == nil {
		return nil, errors.New("Unit does not exist")
	}

	if u1.Type != u2.Type {
		return nil, errors.New("Incompatible units")
	}

	nQ1 := normalizeQuantity(q1, u1)
	nQ2 := normalizeQuantity(q2, u2)

	total := nQ1 + nQ2
	total = total / u1.Modifier

	return &Quantity{
		Unit:     q1.Unit,
		Quantity: total,
	}, nil
}

func (s *service) Substract(q1 *Quantity, q2 *Quantity) (*Quantity, error) {
	u1 := s.unitRepository.FindByName(q1.Unit)
	if u1 == nil {
		return nil, errors.New("Unit does not exist")
	}

	u2 := s.unitRepository.FindByName(q2.Unit)
	if u2 == nil {
		return nil, errors.New("Unit does not exist")
	}

	if u1.Type != u2.Type {
		return nil, errors.New("Incompatible units")
	}

	nQ1 := normalizeQuantity(q1, u1)
	nQ2 := normalizeQuantity(q2, u2)

	total := nQ1 - nQ2
	total = total / u1.Modifier

	return &Quantity{
		Unit:     q1.Unit,
		Quantity: total,
	}, nil
}

func normalizeQuantity(q *Quantity, u *unit.Unit) float64 {
	return q.Quantity * u.Modifier
}
