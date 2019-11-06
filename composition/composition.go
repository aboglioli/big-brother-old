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
	Cost         float64            `bson:"cost" validate:"required"`
	Unit         quantity.Quantity  `bson:"unit" validate:"required"`
	Stock        quantity.Quantity  `bson:"stock" validate:"required"`
	Dependencies []*Dependency      `bson:"dependencies" validate:"required"`

	AutoupdateCost bool      `bson:"autoupdate_cost"`
	Enabled        bool      `bson:"enabled" `
	CreatedAt      time.Time `bson:"createdAt"`
	UpdatedAt      time.Time `bson:"updatedAt"`
}

func NewComposition() *Composition {
	return &Composition{
		ID:             primitive.NewObjectID(),
		AutoupdateCost: true,
		Enabled:        true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (c *Composition) CalculateCost() {
	if c.AutoupdateCost && len(c.Dependencies) > 0 {
		var cost float64
		for _, d := range c.Dependencies {
			cost += d.Subvalue
		}
		c.Cost = math.Round(cost*1000) / 1000
	}
}

func (c *Composition) CostFromQuantity(q quantity.Quantity) float64 {
	nQuantity := q.Normalize()
	nUnit := c.Unit.Normalize()

	return nQuantity * c.Cost / nUnit
}

func (c *Composition) UpsertDependency(d *Dependency) errors.Error {
	if !c.dependencyExists(d.Of.String()) {
		c.Dependencies = append(c.Dependencies, d)
		return nil
	}

	for i, dep := range c.Dependencies {
		if dep.Of.String() == d.Of.String() {
			c.Dependencies[i] = d
			return nil
		}
	}

	return errors.New("composition.Composition.UpsertDependency", "UNKOWN", "")
}

func (c *Composition) RemoveDependency(depID string) errors.Error {
	if !c.dependencyExists(depID) {
		return errors.New("composition.Composition.RemoveDependency", "DEPENDENCY_DOES_NOT_EXIST", "")
	}

	dependencies := make([]*Dependency, 0, len(c.Dependencies))
	for _, dep := range c.Dependencies {
		if dep.Of.String() != depID {
			dependencies = append(dependencies, dep)
		}
	}
	c.Dependencies = dependencies

	return nil
}

func (c1 *Composition) CompareDependencies(c2 *Composition) (left []*Dependency, common []*Dependency, right []*Dependency) {
	for _, dep := range c1.Dependencies {
		if !isDependencyInArray(dep, c2.Dependencies) {
			left = append(left, dep)
		} else {
			common = append(common, dep)
		}
	}

	for _, dep := range c2.Dependencies {
		if !isDependencyInArray(dep, c1.Dependencies) {
			right = append(right, dep)
		}
	}

	return
}

func (c *Composition) dependencyExists(of string) bool {
	for _, d := range c.Dependencies {
		if d.Of.String() == of {
			return true
		}
	}
	return false
}

func isDependencyInArray(d *Dependency, dependencies []*Dependency) bool {
	for _, dep := range dependencies {
		if d.Of.String() == dep.Of.String() {
			return true
		}
	}
	return false
}
