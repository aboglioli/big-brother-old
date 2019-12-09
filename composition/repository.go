package composition

import (
	"context"
	"time"

	"github.com/aboglioli/big-brother/impl/db"
	"github.com/aboglioli/big-brother/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FindAll() ([]*Composition, error)
	FindByID(id string) (*Composition, error)
	FindUses(id string) ([]*Composition, error)
	FindByUpdateUses(upateUses bool) ([]*Composition, error)

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
	path := "composition/repository.FindAll"
	ctx := context.Background()

	cur, err := r.collection.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	if err != nil {
		return nil, errors.NewInternal("FIND_ALL").SetPath(path).SetRef(err)
	}

	var comps []*Composition
	for cur.Next(ctx) {
		var comp Composition

		if err := cur.Decode(&comp); err != nil {
			return nil, errors.NewInternal("CUR_DECODE").SetPath(path).SetRef(err)
		}

		comps = append(comps, &comp)
	}

	if err := cur.Err(); err != nil {
		return nil, errors.NewInternal("CUR").SetPath(path).SetRef(err)
	}

	return comps, nil
}

func (r *repository) FindByID(id string) (*Composition, error) {
	path := "composition/repository.FindByID"
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewInternal("OBJECTID").SetPath(path).SetRef(err)
	}

	filter := bson.M{
		"_id": objID,
	}

	res := r.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, errors.NewInternal("FIND_ONE").SetPath(path).SetRef(res.Err())
	}

	var comp Composition
	if err := res.Decode(&comp); err != nil {
		return nil, errors.NewInternal("DECODE").SetPath(path).SetRef(err)
	}

	return &comp, nil
}

func (r *repository) FindUses(id string) ([]*Composition, error) {
	path := "composition/repository.FindUses"
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewInternal("OBJECTID").SetPath(path).SetRef(err)
	}
	filter := bson.M{
		"dependencies.of": objID,
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewInternal("FIND").SetPath(path).SetRef(err)
	}

	var comps []*Composition
	for cur.Next(ctx) {
		var comp Composition

		if err := cur.Decode(&comp); err != nil {
			return nil, errors.NewInternal("DECODE").SetPath(path).SetRef(err)
		}

		comps = append(comps, &comp)
	}

	if err := cur.Err(); err != nil {
		return nil, errors.NewInternal("CUR").SetPath(path).SetRef(err)
	}

	return comps, nil
}

func (r *repository) FindByUpdateUses(updateUses bool) ([]*Composition, error) {
	path := "composition/repository.FindByUpdateUses"
	ctx := context.Background()

	filter := bson.M{
		"updateUses": updateUses,
	}

	cur, err := r.collection.Find(ctx, filter)
	defer cur.Close(ctx)
	if err != nil {
		return nil, errors.NewInternal("FIND").SetPath(path).SetRef(err)
	}

	var comps []*Composition
	for cur.Next(ctx) {
		var comp Composition

		if err := cur.Decode(&comp); err != nil {
			return nil, errors.NewInternal("CUR_DECODE").SetPath(path).SetRef(err)
		}

		comps = append(comps, &comp)
	}

	if err := cur.Err(); err != nil {
		return nil, errors.NewInternal("CUR").SetPath(path).SetRef(err)
	}

	return comps, nil
}

func (r *repository) Insert(c *Composition) error {
	ctx := context.Background()

	_, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		return errors.NewInternal("INSERT_ONE").SetPath("composition/repository.Insert").SetRef(err)
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
		return errors.NewInternal("INSERT_MANY").SetPath("composition/repository.InsertMany").SetRef(err)
	}

	return nil
}

func (r *repository) Update(c *Composition) error {
	path := "composition/repository.Update"
	ctx := context.Background()

	if c.ID.IsZero() {
		return errors.NewInternal("OBJECTID").SetPath(path)
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
		return errors.NewInternal("UPDATE_ONE").SetPath(path).SetRef(err)
	}

	return nil
}

func (r *repository) Delete(id string) error {
	path := "composition/repository.Delete"
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewInternal("OBJECTID").SetPath(path).SetRef(err)
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
		return errors.NewInternal("UPDATE_ONE").SetPath(path).SetRef(err)
	}

	return nil
}
