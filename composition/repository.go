package composition

import (
	"context"
	"errors"
	"time"

	"github.com/aboglioli/big-brother/infrastructure/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FindAll() ([]*Composition, error)
	FindByID(id string) (*Composition, error)
	FindUses(id string) ([]*Composition, error)

	Insert(*Composition) error
	InsertMany([]*Composition) error
	Update(*Composition) error
	Delete(id string) error
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository() (Repository, error) {
	db, err := db.Get("Composition")

	if err != nil {
		return nil, err
	}

	return &repository{
		collection: db.Collection("composition"),
	}, nil
}

func (r *repository) FindAll() ([]*Composition, error) {
	ctx := context.Background()

	cur, err := r.collection.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	var comps []*Composition
	for cur.Next(ctx) {
		var comp Composition

		if err := cur.Decode(&comp); err != nil {
			return nil, err
		}

		comps = append(comps, &comp)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return comps, nil
}

func (r *repository) FindByID(id string) (*Composition, error) {
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objID,
	}

	res := r.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var comp Composition
	if err := res.Decode(&comp); err != nil {
		return nil, err
	}

	return &comp, nil
}

func (r *repository) FindUses(id string) ([]*Composition, error) {
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"dependencies.of": objID,
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var comps []*Composition
	for cur.Next(ctx) {
		var comp Composition

		if err := cur.Decode(&comp); err != nil {
			return nil, err
		}

		comps = append(comps, &comp)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return comps, nil
}

func (r *repository) Insert(c *Composition) error {
	ctx := context.Background()

	_, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) InsertMany(comps []*Composition) error {
	ctx := context.Background()

	rawComps := make([]interface{}, len(comps))
	for i, c := range comps {
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		rawComps[i] = c
	}

	_, err := r.collection.InsertMany(ctx, rawComps)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(c *Composition) error {
	ctx := context.Background()

	if c.ID.IsZero() {
		return errors.New("Invalid ObjectID")
	}

	c.UpdatedAt = time.Now()

	filter := bson.M{
		"_id": c.ID,
	}

	update := bson.D{
		{"$set", c},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(id string) error {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}

	update := bson.D{
		{"$set", bson.D{
			{"updatedAt", time.Now()},
			{"enabled", false},
		}},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
