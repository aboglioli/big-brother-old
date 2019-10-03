package db

import (
	"context"
	"log"

	"github.com/aboglioli/big-brother/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var databases map[string]*mongo.Database

func init() {
	databases = make(map[string]*mongo.Database)
}

func connect() (*mongo.Client, error) {
	c := config.Get()

	clientOptions := options.Client().ApplyURI(c.MongoURL).SetAuth(
		options.Credential{
			AuthSource: "admin",
			Username:   "admin",
			Password:   "admin",
		})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client, nil
}

func Get(database string) (*mongo.Database, error) {
	if client == nil {
		c, err := connect()
		if err != nil {
			return nil, err
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
