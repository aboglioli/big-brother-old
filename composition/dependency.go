package composition

import (
	"github.com/aboglioli/big-brother/quantity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dependency struct {
	Of       primitive.ObjectID `bson:"of" validate:"required"`
	Quantity quantity.Quantity  `bson:"quantity" validate:"required"`
	Subvalue float64            `bson:"subvalue"`
}
