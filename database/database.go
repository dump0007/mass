package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() (*mongo.Client) {
	MongoDb := "mongodb+srv://dump0007:Dlh592$eL@cluster0.6jugd.mongodb.net/"
	// MongoDb := "mongodb://localhost:27017"
	// fmt.Print(MongoDb)
	clientOptions := options.Client().ApplyURI(MongoDb)
	client, err := mongo.Connect(context.Background(), clientOptions)
	// client, err := mongo.Connect(options.Client().ApplyURI(MongoDb))
	if err != nil {
		return nil
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil
	}

	green := "\033[32m"
	reset := "\033[0m"
	fmt.Println(green + "Mongodb Connected" + reset)
	fmt.Println(MongoDb)
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)

	return collection
}