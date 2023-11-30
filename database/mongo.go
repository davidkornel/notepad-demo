package database

import (
	"context"
	"github.com/go-logr/logr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	uri    string
	logger logr.Logger
}

func NewMongoDB(logger logr.Logger, uri string) *MongoDB {
	return &MongoDB{
		uri:    uri,
		logger: logger,
	}
}

func (db *MongoDB) Connect2MongoDB() error {
	log := db.logger.WithName("Connect2MongoDB")
	var err error
	db.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.uri))
	if err != nil {
		return err
	}
	log.Info("MongoDB uri looks well formatted")
	return nil
}

func (db *MongoDB) CloseMongoDBConnection() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (db *MongoDB) GetClient() *mongo.Client {
	return db.client
}
