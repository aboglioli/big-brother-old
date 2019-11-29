package composition

import (
	"fmt"
	"math"
	"time"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Composition struct {
	ID           primitive.ObjectID `json:"id" bson:"_id""`
	Name         string             `json:"name" bson:"name"`
	Cost         float64            `json:"cost" bson:"cost"`
	Unit         quantity.Quantity  `json:"unit" bson:"unit"`
	Stock        quantity.Quantity  `json:"stock" bson:"stock"`
	Dependencies []Dependency       `json:"dependencies" bson:"dependencies"`

	AutoupdateCost             bool      `json:"autoupdateCost" bson:"autoupdateCost"`
	Enabled                    bool      `json:"enabled" bson:"enabled" `
	Validated                  bool      `json:"validated" bson:"validated"`
	UsesUpdatedSinceLastChange bool      `json:"usesUpdatedSinceLastChange" bson:"usesUpdatedSinceLastChange"`
	CreatedAt                  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt" bson:"updatedAt"`
}

func NewComposition() *Composition {
	return &Composition{
		ID:                         primitive.NewObjectID(),
		AutoupdateCost:             true,
		Enabled:                    true,
		Validated:                  false, // TODO: should be validated asynchronously
		UsesUpdatedSinceLastChange: true,
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
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
		if d.On.Hex() == id {
			return &d
		}
	}
	return nil
}

func (c *Composition) UpsertDependency(d Dependency) {
	updated := false
	for i, dep := range c.Dependencies {
		if dep.On.Hex() == d.On.Hex() {
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
		if dep.On.Hex() != depID {
			dependencies = append(dependencies, dep)
			continue
		}
		removed = true
	}
	c.Dependencies = dependencies

	if !removed {
		return errors.NewValidation().SetPath("composition/composition.RemoveDependency").SetCode("DEPENDENCY_DOES_NOT_EXIST")
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

func (c *Composition) ValidateSchema() errors.Error {
	errGen := errors.NewValidation().SetPath("composition/composition.ValidateSchema")

	if c.Cost < 0 {
		return errGen.SetCode("NEGATIVE_COST").SetMessage(fmt.Sprintf("%v", c.Cost))
	}
	if !c.Unit.IsValid() {
		return errGen.SetCode("INVALID_UNIT").SetMessage(fmt.Sprintf("%v", c.Unit))
	}
	if !c.Stock.IsValid() {
		return errGen.SetCode("INVALID_STOCK").SetMessage(fmt.Sprintf("%v", c.Stock))
	}

	if !c.Stock.Compatible(c.Unit) {
		return errGen.SetCode("INCOMPATIBLE_STOCK_AND_UNIT").SetMessage(fmt.Sprintf("%s != %s", c.Stock.Unit, c.Unit.Unit))
	}

	for i, d := range c.Dependencies {
		if !d.Quantity.IsValid() {
			return errGen.SetCode("INVALID_DEPENDENCY_QUANTITY").SetMessage(fmt.Sprintf("Dependency %d: %v", i, d.Quantity))
		}
	}

	return nil
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
