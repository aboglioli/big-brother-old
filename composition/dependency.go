package composition

import (
	"github.com/aboglioli/big-brother/pkg/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dependency struct {
	On       primitive.ObjectID `bson:"on" json:"on" validate:"required" binding:"required"`
	Quantity quantity.Quantity  `bson:"quantity" json:"quantity" validate:"required" binding:"required"`
	Subvalue float64            `bson:"subvalue" json:"subvalue"`
}

func (d1 Dependency) Equals(d2 Dependency) bool {
	return d1.On.Hex() == d2.On.Hex() && d1.Quantity.Equals(d2.Quantity)
}
