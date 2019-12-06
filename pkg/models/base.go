package models

import "time"

import "go.mongodb.org/mongo-driver/bson/primitive"

type Base struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Enabled   bool               `json:"-" bson:"enabled"`
	Active    bool               `json:"active" bson:"active"`
	Validated bool               `json:"-" bson:"validated"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func NewBase() Base {
	return Base{
		ID:        primitive.NewObjectID(),
		Enabled:   true,
		Active:    true,
		Validated: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
