package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"katseye/internal/infrastructure/persistence/mongodb"
)

type MongoResources struct {
	Client      *mongo.Client
	Database    *mongo.Database
	Collections MongoCollections
}

type MongoCollections struct {
	Products  *mongo.Collection
	Partners  *mongo.Collection
	Addresses *mongo.Collection
	Users     *mongo.Collection
	Consumers *mongo.Collection
}

func newMongoResources(cfg MongoConfig) (*MongoResources, error) {
	client, err := mongodb.NewMongoClient(cfg.URI)
	if err != nil {
		return nil, err
	}

	database := client.Database(cfg.Database)

	return &MongoResources{
		Client:   client,
		Database: database,
		Collections: MongoCollections{
			Products:  database.Collection("products"),
			Partners:  database.Collection("partners"),
			Addresses: database.Collection("addresses"),
			Users:     database.Collection("users"),
			Consumers: database.Collection("consumers"),
		},
	}, nil
}

func (r *MongoResources) Close(ctx context.Context) error {
	if r == nil || r.Client == nil {
		return nil
	}

	return r.Client.Disconnect(ctx)
}
