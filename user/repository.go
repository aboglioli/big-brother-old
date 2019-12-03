package user

import (
	"context"
	"time"

	"github.com/aboglioli/big-brother/infrastructure/db"
	"github.com/aboglioli/big-brother/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	FindByID(id string) (*User, errors.Error)
	FindByUsername(username string) (*User, errors.Error)
	FindByEmail(email string) (*User, errors.Error)

	Insert(u *User) errors.Error
	Update(u *User) errors.Error
	Delete(id string) errors.Error
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository() (Repository, errors.Error) {
	db, err := db.Get("User")

	if err != nil {
		return nil, err
	}

	return &repository{
		collection: db.Collection("user"),
	}, nil
}

func (r *repository) FindByID(id string) (*User, errors.Error) {
	path := "user/repository.FindByID"
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewInternal("OBJECTID_FROM_HEX").SetPath(path).SetMessage(err.Error())
	}

	filter := bson.M{
		"_id": objID,
	}

	res := r.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, errors.NewInternal("FIND_ONE").SetPath(path).SetMessage(err.Error())
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return nil, errors.NewInternal("DECODE").SetPath(path).SetMessage(err.Error())
	}

	return &user, nil
}

func (r *repository) FindByUsername(username string) (*User, errors.Error) {
	path := "user/repository.FindByUsername"
	ctx := context.Background()

	filter := bson.M{
		"username": username,
	}

	res := r.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, errors.NewInternal("FIND_ONE").SetPath(path).SetMessage(res.Err().Error())
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return nil, errors.NewInternal("DECODE").SetPath(path).SetMessage(err.Error())
	}

	return &user, nil
}

func (r *repository) FindByEmail(email string) (*User, errors.Error) {
	path := "user/repository.FindByEmail"
	ctx := context.Background()

	filter := bson.M{
		"email": email,
	}

	res := r.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, errors.NewInternal("FIND_ONE").SetPath(path).SetMessage(res.Err().Error())
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return nil, errors.NewInternal("DECODE").SetPath(path).SetMessage(err.Error())
	}

	return &user, nil
}

func (r *repository) Insert(u *User) errors.Error {
	ctx := context.Background()

	_, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		return errors.NewInternal("INSERT_ONE").SetPath("user/repository.Insert").SetMessage(err.Error())
	}

	return nil
}

func (r *repository) Update(u *User) errors.Error {
	path := "user/repository.Update"
	ctx := context.Background()

	if u.ID.IsZero() {
		return errors.NewInternal("INVALID_OBJECTID").SetPath(path)
	}

	u.UpdatedAt = time.Now()

	filter := bson.M{
		"_id": u.ID,
	}

	update := bson.D{
		{"$set", u},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewInternal("UPDATE_ONE").SetPath(path).SetMessage(err.Error())
	}

	return nil
}

func (r *repository) Delete(id string) errors.Error {
	path := "composition/repository.Delete"
	ctx := context.Background()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewInternal("OBJECTID_FROM_HEX").SetPath(path).SetMessage(err.Error())
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
		return errors.NewInternal("UPDATE_ONE").SetPath(path).SetMessage(err.Error())
	}

	return nil
}
