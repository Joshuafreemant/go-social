package database

import (
	"context"
	"fmt"

	"github.com/Joshuafreemant/go-social/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(config.Config("MONGO_URL"))
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	var ctx = context.Background()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database(config.Config("DB_NAME")), nil
}
