package database

import (
	"context"

	"github.com/nadirbasalamah/go-vrent/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database represents mongoDB database
var Database *mongo.Database
var client *mongo.Client
var err error

// Connect func returns error if connection failed
func Connect() error {
	client, err = mongo.NewClient(options.Client().ApplyURI(config.Config("MONGO_URI")))
	if err != nil {
		return err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}

	Database = client.Database(config.Config("DATABASE"))
	return nil
}

// Disconnect func to disconnect from DB
func Disconnect() {
	client.Disconnect(context.TODO())
}
