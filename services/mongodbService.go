package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDB(ctx context.Context) (*mongo.Collection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	log.Println("Connected to MongoDB")

	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")

	return collection, nil
}
