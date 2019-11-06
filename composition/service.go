package composition

import (
	"fmt"
	"math"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/quantity"
)

type Service interface {
	GetByID(id string) (*Composition, errors.Error)
	Create(*Composition) errors.Error
	Update(*Composition) errors.Error
	Delete(id string) errors.Error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) GetByID(id string) (*Composition, errors.Error) {
	comp, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.New("composition/service.GetByID", "COMPOSITION_NOT_FOUND", err.Error())
	}
	return comp, nil
}

func (s *service) Create(c *Composition) errors.Error {
	if err := s.validateSchema(c); err != nil {
		return err
	}

	if err := s.calculateDependenciesSubvalue(c.Dependencies); err != nil {
		return err
	}
	c.CalculateCost()

	if err := s.repository.Insert(c); err != nil {
		return errors.New("composition/service.Create", "INSERT", err.Error())
	}

	return nil
}

func (s *service) Update(c *Composition) errors.Error {
	errGen := errors.FromPath("composition/service.Update")

	if err := s.validateSchema(c); err != nil {
		return err
	}

	saved, err := s.repository.FindByID(c.ID.String())

	if err != nil {
		return errGen("COMPOSITION_DOES_NOT_EXIST", err.Error())
	}

	new, _, old := c.CompareDependencies(saved)

	if saved.Cost != c.Cost || len(new) > 0 || len(old) > 0 {
		if err := s.calculateDependenciesSubvalue(c.Dependencies); err != nil {
			return err
		}
		c.CalculateCost()
	}

	if err := s.repository.Update(c); err != nil {
		return errGen("UPDATE", err.Error())
	}

	if err := s.updateUses(c); err != nil {
		return errGen("UPDATE_USES", err.Error())
	}

	return nil
}

func (s *service) Delete(id string) errors.Error {
	errGen := errors.FromPath("composition/service.Delete")

	uses, _ := s.repository.FindUses(id)
	if len(uses) > 0 {
		return errGen("COMPOSITION_USED_AS_DEPENDENCY", "")
	}

	if err := s.repository.Delete(id); err != nil {
		return errGen("NOT_FOUND", err.Error())
	}

	return nil
}

func (s *service) updateUses(c *Composition) errors.Error {
	uses, _ := s.repository.FindUses(c.ID.String())

	for _, c := range uses {
		if err := s.Update(c); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) calculateDependenciesSubvalue(dependencies []*Dependency) errors.Error {
	errGen := errors.FromPath("composition/service.calculateDependenciesSubvalue")
	for _, dep := range dependencies {
		comp, err := s.repository.FindByID(dep.Of.String())
		if err != nil {
			return errGen("DEPENDENCY_DOES_NOT_EXIST", err.Error())
		}

		if !dep.Quantity.Compatible(comp.Unit) {
			return errGen("INCOMPATIBLE_DEPENDENCY_QUANTITY", "")
		}

		subvalue := comp.CostFromQuantity(dep.Quantity)

		dep.Subvalue = math.Round(subvalue*100) / 100
	}

	return nil
}

func (s *service) validateSchema(c *Composition) errors.Error {
	errGen := errors.FromPath("composition/service.validateSchema")
	if c.Cost < 0 {
		return errGen("NEGATIVE_COST", fmt.Sprintf("%v", c.Cost))
	}
	if !quantity.IsValid(c.Unit) {
		return errGen("INVALID_UNIT", fmt.Sprintf("%v", c.Unit))
	}
	if !quantity.IsValid(c.Stock) {
		return errGen("INVALID_STOCK", fmt.Sprintf("%v", c.Stock))
	}

	for i, d := range c.Dependencies {
		_, err := s.repository.FindByID(d.Of.String())
		if err != nil {
			return errGen("DEPENDENCY_DOES_NOT_EXIST", err.Error())
		}

		if !quantity.IsValid(d.Quantity) {
			return errGen("INVALID_DEPENDENCY_QUANTITY", fmt.Sprintf("%d> %v", i, d.Quantity))
		}
	}

	return nil
}
