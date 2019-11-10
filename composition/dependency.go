package composition

import (
	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dependency struct {
	Of       primitive.ObjectID `bson:"of" json:"of" validate:"required" binding:"required"`
	Quantity quantity.Quantity  `bson:"quantity" json:"quantity" validate:"required" binding:"required"`
	Subvalue float64            `bson:"subvalue" json:"subvalue"`
}

func (d1 Dependency) Equals(d2 Dependency) bool {
	return d1.Of.Hex() == d2.Of.Hex() && d1.Quantity.Equals(d2.Quantity)
}
