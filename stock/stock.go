package stock

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/aboglioli/big-brother/quantity"
)

type Stock struct {
	ID          primitive.ObjectID `bson:"_id"`
	ProductID   primitive.ObjectID `bson:"productId" validate:"required"`
	ProductName string             `bson:"productName" validate:"required"`
	Quantity    quantity.Quantity  `bson:"quantity" validate:"required"`
}

func NewStock(pID primitive.ObjectID, pName string, q quantity.Quantity) *Stock {
	return &Stock{
		ID:          primitive.NewObjectID(),
		ProductID:   pID,
		ProductName: pName,
		Quantity:    q,
	}
}
