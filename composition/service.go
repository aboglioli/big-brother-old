package composition

import (
	"fmt"
	"math"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	GetByID(id string) (*Composition, error)
	Create(req *CreateRequest) (*Composition, error)
	Update(compID string, req *UpdateRequest) (*Composition, error)
	Delete(id string) error

	UpdateUses(c *Composition) ([]*Composition, error)
	Validate(id string) error
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

func (s *service) GetByID(id string) (*Composition, error) {
	comp, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.NewStatus("COMPOSITION_NOT_FOUND").SetPath("composition/service.GetByID").SetRef(err)
	}
	if !comp.Enabled {
		return nil, errors.NewStatus("COMPOSITION_IS_DELETED").SetPath("composition/service.GetByID")
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

// Create creates a new Composition.
// ID can be defined by the user or not.
/**
* @api {topic} composition.created composition.created
* @apiName CompositionCreated
* @apiGroup RabbitMQ
*
* @apiDescription Emits a new event when a new composition is created
*
* @apiSuccessExample {json} Body
* {
* 	"type": "CompositionCreated",
* 	"payload": composition data
* }
 */
func (s *service) Create(req *CreateRequest) (*Composition, error) {
	path := "composition/service.Create"
	c := NewComposition()

	if req.ID != nil {
		id, err := primitive.ObjectIDFromHex(*req.ID)
		if err != nil {
			return nil, errors.NewStatus("INVALID_ID").SetPath(path).SetRef(err)
		}
		if existingComp, err := s.repository.FindByID(*req.ID); existingComp != nil || err == nil {
			return nil, errors.NewStatus("COMPOSITION_ALREADY_EXISTS").SetPath(path).SetMessage(fmt.Sprintf("Composition with ID %s exists", *req.ID)).SetRef(err)
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

	if err := s.repository.Insert(c); err != nil {
		return nil, errors.NewStatus("INSERT").SetPath(path).SetRef(err)
	}

	// Publish event: composition.created
	event, opts := NewCompositionCreatedEvent(c)
	if err := s.eventMgr.Publish(event, opts); err != nil {
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

// Update updates an existing Composition.
/**
* @api {topic} composition.updated composition.updated
* @apiName CompositionUpdatedManually
* @apiGroup RabbitMQ
*
* @apiDescription Emits a new event when an existing composition is updated.
* This event can be of type "CompositionUpdatedManually" or
* "CompositionUpdatedAutomatically". The last one is published once a
* composition is updated due to a dependency change.
*
* @apiSuccessExample {json} Body
* {
* 	"type": "CompositionUpdatedManually",
* 	"payload": composition data
* }
 */
func (s *service) Update(id string, req *UpdateRequest) (*Composition, error) {
	path := "composition/service.Update"

	if req.ID != nil && *req.ID != id {
		return nil, errors.NewStatus("ID_DOES_NOT_MATCH").SetPath(path).SetMessage(fmt.Sprintf("%s != %s", *req.ID, id))
	}

	c, err := s.repository.FindByID(id)
	if err != nil {
		return nil, errors.NewStatus("COMPOSITION_DOES_NOT_EXIST").SetPath(path).SetRef(err)
	}

	if !c.Enabled {
		return nil, errors.NewStatus("COMPOSITION_IS_DELETED").SetPath(path)
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
		return nil, errors.NewStatus("CANNOT_CHANGE_UNIT_TYPE").SetPath(path).SetMessage(fmt.Sprintf("%v != %v", c.Unit, req.Unit))
	}

	removed, _, added := c.CompareDependencies(req.Dependencies)

	if len(removed) == 0 && len(added) == 0 {
		// If nothings changes, recalculate cost from dependencies
		if c.AutoupdateCost {
			c.calculateCostFromDependencies()
		}
	} else {
		for _, dep := range removed {
			c.RemoveDependency(dep.On.Hex())
		}

		for i, dep := range added {
			depComp, err := s.repository.FindByID(dep.On.Hex())
			if err != nil {
				return nil, errors.NewStatus("DEPENDENCY_DOES_NOT_EXIST").SetPath(path).SetRef(err)
			}

			if !dep.Quantity.IsValid() {
				return nil, errors.NewStatus("INVALID_DEPENDENCY_QUANTITY").SetPath(path).SetMessage("Dependency nro %d: %s", i, dep.On.Hex())
			}

			if !dep.Quantity.Compatible(depComp.Unit) {
				return nil, errors.NewStatus("INCOMPATIBLE_DEPENDENCY_QUANTITY").SetPath(path).SetMessage("Dependency nro %d (%s): %v != %v", i, dep.On.Hex(), dep.Quantity, depComp.Unit)
			}

			subvalue := depComp.CostFromQuantity(dep.Quantity)
			dep.Subvalue = math.Round(subvalue*1000) / 1000

			c.UpsertDependency(dep)
		}
	}

	c.UsesUpdatedSinceLastChange = false

	if err := s.repository.Update(c); err != nil {
		return nil, errors.NewStatus("UPDATE").SetRef(err)
	}

	// Publish event: composition.updated
	event, opts := NewCompositionUpdatedManuallyEvent(c)
	if err := s.eventMgr.Publish(event, opts); err != nil {
		return nil, errors.NewStatus("FAILED_TO_PUBLISH").SetRef(err)
	}

	return c, nil
}

// Delete deletes an existing Composition.
/**
* @api {topic} composition.deleted composition.deleted
* @apiName CompositionDeleted
* @apiGroup RabbitMQ
*
* @apiDescription Emits a new event when an existing composition is deleted.
*
* @apiSuccessExample {json} Body
* {
* 	"type": "CompositionDeleted",
* 	"payload": composition data
* }
 */
func (s *service) Delete(id string) error {
	path := "composition/service.Delete"

	c, err := s.repository.FindByID(id)
	if err != nil {
		return errors.NewStatus("DELETE").SetPath(path).SetRef(err)
	}

	uses, _ := s.repository.FindUses(id)
	if len(uses) > 0 {
		return errors.NewStatus("COMPOSITION_USED_AS_DEPENDENCY").SetPath(path).SetMessage("Composition used as dependecy in %d compositions", len(uses))
	}

	if err := s.repository.Delete(id); err != nil {
		return errors.NewStatus("NOT_FOUND").SetPath(path).SetRef(err)
	}

	// Publish event
	event, opts := NewCompositionDeletedEvent(c)
	if err := s.eventMgr.Publish(event, opts); err != nil {
		return errors.NewStatus("PUBLISH").SetPath(path).SetRef(err)
	}

	return nil
}

// Update updates uses of an updated composition.
/**
* @api {topic} composition.updated composition.updated
* @apiName CompositionUpdatedAutomatically
* @apiGroup RabbitMQ
*
* @apiDescription Emits a new event when an existing composition is updated.
* This event can be of type "CompositionUpdatedManually" or
* "CompositionUpdatedAutomatically". The last one is published once a
* composition is updated due to a dependency change.
*
* @apiSuccessExample {json} Body
* {
* 	"type": "CompositionsUpdatedAutomatically",
* 	"payload": list of compositions
* }
 */
func (s *service) UpdateUses(c *Composition) ([]*Composition, error) {
	path := "composition/service.UpdateUses"

	cache := make(map[string]*Composition)

	err := s.updateUses(c, cache)
	if err != nil {
		return nil, errors.NewStatus("UPDATE_USES").SetPath(path).SetRef(err)
	}

	comps := make([]*Composition, 0)
	for _, u := range cache {
		if err := s.repository.Update(u); err != nil {
			return nil, errors.NewStatus("UPDATE_USES").SetPath(path).SetRef(err)
		}
		comps = append(comps, u)
	}

	if len(comps) > 0 {
		event, opts := NewCompositionsUpdatedAutomaticallyEvent(comps)
		if err := s.eventMgr.Publish(event, opts); err != nil {
			return nil, err
		}
	}

	return comps, nil
}

func (s *service) Validate(compID string) error {
	comp, err := s.repository.FindByID(compID)
	if err != nil {
		return err
	}

	comp.Validated = true

	if err := s.repository.Update(comp); err != nil {
		return err
	}

	return nil
}

func (s *service) updateUses(c *Composition, cache map[string]*Composition) error {
	path := "composition/service.updateUses"

	uses, _ := s.repository.FindUses(c.ID.Hex())

	for _, u := range uses {
		cachedUse, ok := cache[u.ID.Hex()]
		if ok {
			u = cachedUse
		}

		dep := u.FindDependencyByID(c.ID.Hex())

		subvalue := c.CostFromQuantity(dep.Quantity)
		dep.Subvalue = math.Round(subvalue*1000) / 1000

		u.UpsertDependency(*dep)

		cache[u.ID.Hex()] = u

		// Update uses
		if err := s.updateUses(u, cache); err != nil {
			return errors.NewStatus("UPDATE_USES").SetPath(path).SetRef(err)
		}
	}

	return nil
}

func (s *service) validateSchema(c *Composition) error {
	path := "composition/service.validateSchema"

	if err := c.ValidateSchema(); err != nil {
		return err
	}

	newDependencies := make([]Dependency, len(c.Dependencies))
	for i, dep := range c.Dependencies {
		comp, err := s.repository.FindByID(dep.On.Hex())
		if err != nil {
			return errors.NewStatus("DEPENDENCY_DOES_NOT_EXIST").SetPath(path).SetRef(err)
		}

		if !dep.Quantity.Compatible(comp.Unit) {
			return errors.NewStatus("INCOMPATIBLE_DEPENDENCY_QUANTITY").SetPath(path).SetMessage("Dependency %d: %v != %v", i, dep.Quantity, comp.Unit)
		}

		subvalue := comp.CostFromQuantity(dep.Quantity)
		dep.Subvalue = math.Round(subvalue*1000) / 1000
		newDependencies[i] = dep
	}
	c.SetDependencies(newDependencies)

	return nil
}
