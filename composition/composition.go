package composition

import (
	"time"

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
	var cost float64
	for _, d := range c.Dependencies {
		cost += d.Subvalue
	}
	c.Cost = cost
}
