package composition

import (
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/quantity"
)

type Service interface {
	Create(*Composition) errors.Error
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

func (s *service) Create(c *Composition) errors.Error {
	if err := s.validateSchema(c); err != nil {
		return err
	}

	for _, d := range c.Dependencies {
		_, err := s.repository.FindByID(d.Of.String())
		if err != nil {
			return errors.New("composition/service.Create", "DEPENDENCY_DOES_NOT_EXIST", err.Error())
		}
	}

	err := s.repository.Insert(c)
	if err != nil {
		return errors.New("composition/service.Create", "INSERT", err.Error())
	}
	return nil
}

func (s *service) validateSchema(c *Composition) errors.Error {
	errGen := errors.FromPath("composition/service.validateSchema")
	if c.Cost < 0 {
		return errGen("NEGATIVE_COST", "")
	}
	if !s.quantityService.IsValid(&c.Unit) {
		return errGen("INVALID_UNIT", "")
	}
	if !s.quantityService.IsValid(&c.Stock) {
		return errGen("INVALID_STOCK", "")
	}

	for _, d := range c.Dependencies {
		if !s.quantityService.IsValid(&d.Quantity) {
			return errGen("INVALID_DEPENDENCY_QUANTITY", "")
		}
	}

	return nil
}
