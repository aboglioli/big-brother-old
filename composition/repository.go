package composition

import (
	"context"
	"time"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/infrastructure/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FindAll() ([]*Composition, errors.Error)
	FindByID(id string) (*Composition, errors.Error)
	FindUses(id string) ([]*Composition, errors.Error)

	Insert(*Composition) errors.Error
	InsertMany([]*Composition) errors.Error
	Update(*Composition) errors.Error
	Delete(id string) errors.Error
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository() (Repository, errors.Error) {
	db, err := db.Get("Composition")

	if err != nil {
		return nil, err
	}

	return &repository{
		collection: db.Collection("composition"),
	}, nil
}

func (r *repository) FindAll() ([]*Composition, errors.Error) {
	errGen := errors.InternalFromPath("composition/repository.FindAll")
	ctx := context.Background()

	cur, err := r.collection.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	if err != nil {
		return nil, errGen("FIND", err.Error())
	}

	var comps []*Composition
	for cur.Next(ctx) {
		var comp Composition

		if err := cur.Decode(&comp); err != nil {
			return nil, errGen("CUR_DECODE", err.Error())
		}

		comps = append(comps, &comp)
	}

	if err := cur.Err(); err != nil {
		return nil, errGen("CUR", err.Error())
	}

	return comps, nil
}

func (r *repository) FindByID(id string) (*Composition, errors.Error) {
	errGen := errors.InternalFromPath("composition/repository.FindByID")

	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errGen("OBJECTID_FROM_HEX", err.Error())
	}

	filter := bson.M{
		"_id": objID,
	}

	res := r.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, errGen("FIND_ONE", res.Err().Error())
	}

	var comp Composition
	if err := res.Decode(&comp); err != nil {
		return nil, errGen("DECODE", err.Error())
	}

	return &comp, nil
}

func (r *repository) FindUses(id string) ([]*Composition, errors.Error) {
	errGen := errors.InternalFromPath("composition/repository.FindUses")

	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errGen("OBJECTID_FROM_HEX", err.Error())
	}
	filter := bson.M{
		"dependencies.of": objID,
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, errGen("FIND", err.Error())
	}

	var comps []*Composition
	for cur.Next(ctx) {
		var comp Composition

		if err := cur.Decode(&comp); err != nil {
			return nil, errGen("DECODE", err.Error())
		}

		comps = append(comps, &comp)
	}

	if err := cur.Err(); err != nil {
		return nil, errGen("CUR", err.Error())
	}

	return comps, nil
}

func (r *repository) Insert(c *Composition) errors.Error {
	ctx := context.Background()

	_, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		return errors.NewInternal("composition/repository.Insert", "INSERT_ONE", err.Error())
	}

	return nil
}

func (r *repository) InsertMany(comps []*Composition) errors.Error {
	ctx := context.Background()

	rawComps := make([]interface{}, len(comps))
	for i, c := range comps {
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		rawComps[i] = c
	}

	_, err := r.collection.InsertMany(ctx, rawComps)
	if err != nil {
		return errors.NewInternal("composition/repository.InsertMany", "INSERT_MANY", err.Error())
	}

	return nil
}

func (r *repository) Update(c *Composition) errors.Error {
	errGen := errors.InternalFromPath("composition/repository.Update")

	ctx := context.Background()

	if c.ID.IsZero() {
		return errGen("INVALID_OBJECTID", "")
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
		return errGen("UPDATE_ONE", err.Error())
	}

	return nil
}

func (r *repository) Delete(id string) errors.Error {
	errGen := errors.InternalFromPath("composition/repository.Delete")

	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errGen("OBJECTID_FROM_HEX", err.Error())
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
		return errGen("UPDATE_ONE", err.Error())
	}

	return nil
}
