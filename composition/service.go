package composition

import (
	"fmt"
	"math"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/events"
	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	GetByID(id string) (*Composition, errors.Error)
	Create(req *CreateRequest) (*Composition, errors.Error)
	Update(compID string, req *UpdateRequest) (*Composition, errors.Error)
	Delete(id string) errors.Error

	UpdateUses(c *Composition) (int, errors.Error)
}

type service struct {
	repository Repository
	eventMgr   events.Manager
}

func NewService(r Repository, e events.Manager) Service {
	return &service{
		repository: r,
		eventMgr:   e,
	}
}

func (s *service) GetByID(id string) (*Composition, errors.Error) {
	comp, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.NewValidation("composition/service.GetByID", "COMPOSITION_NOT_FOUND", err.Error())
	}
	return comp, nil
}

type CreateRequest struct {
	ID           *string            `json:"id"`
	Name         string             `json:"name"`
	Cost         float64            `json:"cost"`
	Unit         quantity.Quantity  `json:"unit" binding:"required"`
	Stock        *quantity.Quantity `json:"stock"`
	Dependencies []Dependency       `json:"dependencies"`

	AutoupdateCost *bool `json:"autoupdateCost"`
}

func (s *service) Create(req *CreateRequest) (*Composition, errors.Error) {
	errGen := errors.ValidationFromPath("composition/service.Create")
	c := NewComposition()

	if req.ID != nil {
		id, err := primitive.ObjectIDFromHex(*req.ID)
		if err != nil {
			return nil, errGen("INVALID_ID", err.Error())
		}
		if existingComp, err := s.repository.FindByID(*req.ID); existingComp != nil || err == nil {
			return nil, errGen("COMPOSITION_ALREADY_EXISTS", fmt.Sprintf("Composition with ID %s exists", *req.ID))
		}
		c.ID = id
	}

	c.Name = req.Name
	c.Cost = req.Cost
	c.Unit = req.Unit
	if req.Stock != nil {
		c.Stock = *req.Stock
	} else {
		c.Stock = quantity.Quantity{0, c.Unit.Unit}
	}
	if req.AutoupdateCost != nil {
		c.AutoupdateCost = *req.AutoupdateCost
	}

	c.SetDependencies(req.Dependencies)

	if err := s.validateSchema(c); err != nil {
		return nil, err
	}

	deps, err := s.calculateDependenciesSubvalues(c.Dependencies)
	if err != nil {
		return nil, err
	}

	c.SetDependencies(deps)

	if err := s.repository.Insert(c); err != nil {
		return nil, errGen("INSERT", err.Error())
	}

	// Publish event: composition.created
	evt := NewEvent("CompositionCreated", c)
	body, err := evt.ToBytes()
	if err != nil {
		return nil, err
	}
	if err := s.eventMgr.Publish("composition", "topic", "composition.created", body); err != nil {
		return nil, err
	}

	return c, nil
}

type UpdateRequest struct {
	ID           *string            `json:"id"`
	Name         *string            `json:"name"`
	Cost         *float64           `json:"cost"`
	Unit         *quantity.Quantity `json:"unit"`
	Stock        *quantity.Quantity `json:"stock"`
	Dependencies []Dependency       `json:"dependencies"`

	AutoupdateCost *bool `json:"autoupdateCost"`
}

func (s *service) Update(id string, req *UpdateRequest) (*Composition, errors.Error) {
	errGen := errors.ValidationFromPath("composition/service.Update")

	if req.ID != nil && *req.ID != id {
		return nil, errGen("ID_DOES_NOT_MATCH", fmt.Sprintf("%s != %s", *req.ID, id))
	}

	c, rawErr := s.repository.FindByID(id)
	if rawErr != nil {
		return nil, errGen("COMPOSITION_DOES_NOT_EXIST", rawErr.Error())
	}

	if !c.Enabled {
		return nil, errGen("COMPOSITION_IS_DELETED", "")
	}

	savedUnit := c.Unit

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Cost != nil {
		c.Cost = *req.Cost
	}
	if req.Unit != nil {
		c.Unit = *req.Unit
	}
	if req.Stock != nil {
		c.Stock = *req.Stock
	}
	if req.AutoupdateCost != nil {
		c.AutoupdateCost = *req.AutoupdateCost
	}

	if err := s.validateSchema(c); err != nil {
		return nil, err
	}

	if !savedUnit.Compatible(c.Unit) {
		return nil, errGen("CANNOT_CHANGE_UNIT_TYPE", fmt.Sprintf("%s != %s", c.Unit.Unit, req.Unit.Unit))
	}

	removed, _, added := c.CompareDependencies(req.Dependencies)

	for _, dep := range removed {
		c.RemoveDependency(dep.Of.Hex())
	}

	for _, dep := range added {
		depComp, err := s.repository.FindByID(dep.Of.Hex())
		if err != nil {
			return nil, errGen("DEPENDENCY_DOES_NOT_EXIST", err.Error())
		}

		if !quantity.IsValid(dep.Quantity) {
			return nil, errGen("INVALID_DEPENDENCY_QUANTITY", "")
		}

		if !dep.Quantity.Compatible(depComp.Unit) {
			return nil, errGen("INCOMPATIBLE_DEPENDENCY_QUANTITY", "")
		}

		subvalue := depComp.CostFromQuantity(dep.Quantity)
		dep.Subvalue = math.Round(subvalue*100) / 100

		c.UpsertDependency(dep)
	}

	if err := s.repository.Update(c); err != nil {
		return nil, errGen("UPDATE", err.Error())
	}

	// Publish event: composition.updated
	evt := NewEvent("CompositionUpdatedManually", c)
	body, err := evt.ToBytes()
	if err != nil {
		return nil, err
	}
	if err := s.eventMgr.Publish("composition", "topic", "composition.updated", body); err != nil {
		return nil, err
	}

	// Uses are updated asynchronously in another service reacting to the above sent event
	// Update uses
	// if err := s.UpdateUses(c); err != nil {
	// 	return nil, errGen("UPDATE_USES", err.Error())
	// }

	return c, nil
}

func (s *service) Delete(id string) errors.Error {
	errGen := errors.ValidationFromPath("composition/service.Delete")

	c, err := s.repository.FindByID(id)
	if err != nil {
		return errGen("DELETE", err.Error())
	}

	uses, _ := s.repository.FindUses(id)
	if len(uses) > 0 {
		return errGen("COMPOSITION_USED_AS_DEPENDENCY", "")
	}

	if err := s.repository.Delete(id); err != nil {
		return errGen("NOT_FOUND", err.Error())
	}

	// Publish event
	evt := NewEvent("CompositionDeleted", c)
	body, err := evt.ToBytes()
	if err != nil {
		return err
	}
	if err := s.eventMgr.Publish("composition", "topic", "composition.deleted", body); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateUses(c *Composition) (int, errors.Error) {
	uses, _ := s.repository.FindUses(c.ID.Hex())
	count := 0

	for _, u := range uses {
		dep := u.FindDependencyByID(c.ID.Hex())

		subvalue := c.CostFromQuantity(dep.Quantity)
		dep.Subvalue = math.Round(subvalue*1000) / 1000

		u.UpsertDependency(*dep)

		if err := s.repository.Update(u); err != nil {
			return count, errors.NewValidation("composition/service.updateUses", "UPDATE_USES", err.Error())
		}

		count++

		// Publish event
		evt := NewEvent("CompositionUpdatedAutomatically", u)
		body, err := evt.ToBytes()
		if err != nil {
			return count, err
		}
		if err := s.eventMgr.Publish("composition", "topic", "composition.updated", body); err != nil {
			return count, err
		}

		// Update uses
		subcount, err := s.UpdateUses(u)
		if err != nil {
			return count + subcount, err
		}

		count = count + subcount
	}

	return count, nil
}

func (s *service) calculateDependenciesSubvalues(dependencies []Dependency) ([]Dependency, errors.Error) {
	errGen := errors.ValidationFromPath("composition/service.calculateDependenciesSubvalue")

	newDependencies := make([]Dependency, len(dependencies))
	for i, dep := range dependencies {
		comp, err := s.repository.FindByID(dep.Of.Hex())
		if err != nil || comp == nil {
			return nil, errGen("DEPENDENCY_DOES_NOT_EXIST", err.Error())
		}

		if !dep.Quantity.Compatible(comp.Unit) {
			return nil, errGen("INCOMPATIBLE_DEPENDENCY_QUANTITY", "")
		}

		subvalue := comp.CostFromQuantity(dep.Quantity)
		dep.Subvalue = math.Round(subvalue*1000) / 1000
		newDependencies[i] = dep
	}

	return newDependencies, nil
}

func (s *service) validateSchema(c *Composition) errors.Error {
	errGen := errors.ValidationFromPath("composition/service.validateSchema")
	if c.Cost < 0 {
		return errGen("NEGATIVE_COST", fmt.Sprintf("%v", c.Cost))
	}
	if !quantity.IsValid(c.Unit) {
		return errGen("INVALID_UNIT", fmt.Sprintf("%v", c.Unit))
	}
	if !quantity.IsValid(c.Stock) {
		return errGen("INVALID_STOCK", fmt.Sprintf("%v", c.Stock))
	}

	if !c.Stock.Compatible(c.Unit) {
		return errGen("INCOMPATIBLE_STOCK_AND_UNIT", fmt.Sprintf("%s != %s", c.Stock.Unit, c.Unit.Unit))
	}

	for i, d := range c.Dependencies {
		_, err := s.repository.FindByID(d.Of.Hex())
		if err != nil {
			return errGen("DEPENDENCY_DOES_NOT_EXIST", err.Error())
		}

		if !quantity.IsValid(d.Quantity) {
			return errGen("INVALID_DEPENDENCY_QUANTITY", fmt.Sprintf("%d> %v", i, d.Quantity))
		}
	}

	return nil
}
