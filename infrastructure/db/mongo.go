package db

import (
	"context"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var databases = make(map[string]*mongo.Database)

func connect() (*mongo.Client, error) {
	conf := config.Get()
	ctx := context.Background()

	options := options.Client().ApplyURI(conf.MongoURL).SetAuth(
		options.Credential{
			AuthSource: conf.MongoAuthSource,
			Username:   conf.MongoUsername,
			Password:   conf.MongoPassword,
		})

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func Get(database string) (*mongo.Database, errors.Error) {
	if client == nil {
		c, err := connect()
		if err != nil {
			return nil, errors.NewInternal("FAILED_TO_CONNECT_TO_DB").SetPath("infrastructure/db/database.Get").SetMessage(err.Error())
		}
		client = c
	}

	d, ok := databases[database]

	if !ok {
		d = client.Database(database)
		databases[database] = d
	}

	return d, nil
}
