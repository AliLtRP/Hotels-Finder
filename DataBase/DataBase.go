package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Schema struct {
	Name     string   `bson: "name"`
	Rate     float64  `bson: "rate"`
	Price    int      `bson: "price"`
	Location string   `bson: "location"`
	Images   []string `bson: "images"`
	Comments []string `bson: "comments"`
}

type User struct {
	User     string `bson: "user"`
	password string `bson: "password"`
}

func Connection() (*mongo.Client, context.Context) {

	// mad
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:admin@hotels.ncxhxs1.mongodb.net/?retryWrites=true&w=majority"))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("DataBase is connected")
	return client, ctx
}

func Collection(client *mongo.Client, ctx context.Context, collectionName string) *mongo.Collection {
	collection := client.Database("GOLang-Hotels").Collection(collectionName)
	return collection
}

func InsertOne(ctx context.Context, collection *mongo.Collection, data bson.D) *mongo.InsertOneResult {

	res, err := collection.InsertOne(ctx, data)

	if err != nil {
		fmt.Println(err)
	}

	return res
}
