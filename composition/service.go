package composition

import (
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/quantity"
)

type Service interface {
	Create(*Composition) errors.Error
	CalculateCostFromDependencies(*Composition) errors.Error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
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

// Calculate cost of composition from dependencies
func (s *service) CalculateCostFromDependencies(c *Composition) errors.Error {
	if c.AutoupdateCost && len(c.Dependencies) > 0 {
		for _, dep := range c.Dependencies {
			comp, err := s.repository.FindByID(dep.Of.String())
			if err != nil {
				return errors.New("composition/service.calculateCost", "DEPENDENCY_DOES_NOT_EXIST", err.Error())
			}

			if !dep.Quantity.Compatible(comp.Unit) {
				return errors.New("compatible/service.calculateCost", "INCOMPATBLE_QUANTITIES", "")
			}

			nDepQ := dep.Quantity.Normalize()
			nCompQ := comp.Unit.Normalize()

			dep.Subvalue = nDepQ * comp.Cost / nCompQ
		}
		c.CalculateCost()
	}
	return nil
}

func (s *service) validateSchema(c *Composition) errors.Error {
	errGen := errors.FromPath("composition/service.validateSchema")
	if c.Cost < 0 {
		return errGen("NEGATIVE_COST", "")
	}
	if !quantity.IsValid(c.Unit) {
		return errGen("INVALID_UNIT", "")
	}
	if !quantity.IsValid(c.Stock) {
		return errGen("INVALID_STOCK", "")
	}

	for _, d := range c.Dependencies {
		if !quantity.IsValid(d.Quantity) {
			return errGen("INVALID_DEPENDENCY_QUANTITY", "")
		}
	}

	return nil
}
