package utils

import (
	"context"
	"errors"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE = "url_shortner"
)

var DatabaseClient *mongo.Client

// ConnectDB will establish connection with mongodb
func ConnectDB() (client *mongo.Client) {
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

	// Connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	handlerError(err)
	log.Println("Connected to Mongodb")
	return client
}

func GetCollection(client *mongo.Client, collectionName string) (*mongo.Collection, error) {
	db := client.Database(DATABASE)
	names, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	if !stringInSlice(collectionName, names) {
		return nil, errors.New("collection does not exist")
	}

	return db.Collection(collectionName), nil
}
