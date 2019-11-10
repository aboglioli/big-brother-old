package composition

import (
	"math"
	"time"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Composition struct {
	ID           primitive.ObjectID `bson:"_id" validate:"required"`
	Name         string             `bson:"name"`
	Cost         float64            `bson:"cost" validate:"required"`
	Unit         quantity.Quantity  `bson:"unit" validate:"required"`
	Stock        quantity.Quantity  `bson:"stock" validate:"required"`
	Dependencies []Dependency       `bson:"dependencies" validate:"required"`

	AutoupdateCost bool      `bson:"autoupdate_cost"`
	Enabled        bool      `bson:"enabled" `
	Validated      bool      `bson:"validated"`
	CreatedAt      time.Time `bson:"createdAt"`
	UpdatedAt      time.Time `bson:"updatedAt"`
}

func NewComposition() *Composition {
	return &Composition{
		ID:             primitive.NewObjectID(),
		AutoupdateCost: true,
		Enabled:        true,
		Validated:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (c *Composition) CostFromQuantity(q quantity.Quantity) float64 {
	nQuantity := q.Normalize()
	nUnit := c.Unit.Normalize()

	if nUnit == 0 {
		return 0
	}

	return nQuantity * c.Cost / nUnit
}

func (c *Composition) SetDependencies(deps []Dependency) {
	c.Dependencies = deps
	c.calculateCostFromDependencies()
}

func (c *Composition) FindDependencyByID(id string) *Dependency {
	for _, d := range c.Dependencies {
		if d.Of.Hex() == id {
			return &d
		}
	}
	return nil
}

func (c *Composition) UpsertDependency(d Dependency) {
	updated := false
	for i, dep := range c.Dependencies {
		if dep.Of.Hex() == d.Of.Hex() {
			c.Dependencies[i] = d
			updated = true
			break
		}
	}

	if !updated {
		c.Dependencies = append(c.Dependencies, d)
	}

	c.calculateCostFromDependencies()
}

func (c *Composition) RemoveDependency(depID string) errors.Error {
	removed := false
	dependencies := make([]Dependency, 0, len(c.Dependencies))
	for _, dep := range c.Dependencies {
		if dep.Of.Hex() != depID {
			dependencies = append(dependencies, dep)
			continue
		}
		removed = true
	}
	c.Dependencies = dependencies

	if !removed {
		return errors.New("composition/composition.RemoveDependency", "DEPENDENCY_DOES_NOT_EXIST", "")
	}

	c.calculateCostFromDependencies()

	return nil
}

func (c1 *Composition) CompareDependencies(deps []Dependency) (left []Dependency, common []Dependency, right []Dependency) {
	for _, dep := range c1.Dependencies {
		if !isDependencyInArray(dep, deps) {
			left = append(left, dep)
		} else {
			common = append(common, dep)
		}
	}

	for _, dep := range deps {
		if !isDependencyInArray(dep, c1.Dependencies) {
			right = append(right, dep)
		}
	}

	return
}

func (c *Composition) calculateCostFromDependencies() {
	if c.AutoupdateCost && len(c.Dependencies) > 0 {
		var cost float64
		for _, d := range c.Dependencies {
			cost += d.Subvalue
		}
		c.Cost = math.Round(cost*1000) / 1000
	}
}

func isDependencyInArray(d Dependency, dependencies []Dependency) bool {
	for _, dep := range dependencies {
		if d.Equals(dep) {
			return true
		}
	}
	return false
}
