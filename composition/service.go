package composition

import (
	"errors"

	"github.com/aboglioli/big-brother/quantity"
)

type Service interface {
	Create(*Composition) error
}

type service struct {
	repository      Repository
	quantityService quantity.Service
}

func NewService(r Repository, qServ quantity.Service) Service {
	return &service{
		repository:      r,
		quantityService: qServ,
	}
}

func (s *service) Create(c *Composition) error {
	if err := s.validateSchema(c); err != nil {
		return err
	}
	return s.repository.Insert(c)
}

func (s *service) validateSchema(c *Composition) error {
	if c.Cost < 0 {
		return errors.New("Negative cost")
	}
	if !s.quantityService.IsValid(&c.Unit) {
		return errors.New("Invalid Unit")
	}
	if !s.quantityService.IsValid(&c.Stock) {
		return errors.New("Invalid Stock")
	}
	return nil
}
