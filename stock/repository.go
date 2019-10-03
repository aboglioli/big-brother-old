package stock

import (
	"context"

	"github.com/aboglioli/big-brother/infrastructure/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Insert(stock *Stock) (*Stock, error)
	FindByID(id string) (*Stock, error)
	FindByProductID(productID string) (*Stock, error)
}

type repository struct {
	collection *mongo.Collection
}

func newRepository() (Repository, error) {
	db, err := db.Get("stock")

	if err != nil {
		return nil, err
	}

	return &repository{
		collection: db.Collection("stock"),
	}, nil
}

func (r *repository) Insert(stock *Stock) (*Stock, error) {
	_, err := r.collection.InsertOne(context.Background(), stock)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (r repository) FindByID(id string) (*Stock, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	var stock Stock
	filter := bson.D{
		{"_id", objID},
	}

	if err := r.collection.FindOne(context.Background(), filter).Decode(&stock); err != nil {
		return nil, err
	}

	return &stock, nil
}

func (r repository) FindByProductID(productID string) (*Stock, error) {
	objID, _ := primitive.ObjectIDFromHex(productID)
	var stock Stock
	filter := bson.D{
		{"productId", objID},
	}

	if err := r.collection.FindOne(context.Background(), filter).Decode(&stock); err != nil {
		return nil, err
	}

	return &stock, nil
}
