package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func Connect(connectionCh chan<- string) {

	connectionCh <- "Connecting to database..."

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatalln("MONGODB_URI not found")
	}

	options := options.Client().ApplyURI(uri).SetAppName("Gin Basic Api")

	client, _ := mongo.Connect(context.Background(), options)

	connectionCh <- "Pinging database instance..."
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalln(err)
	}

	Db = client.Database("gin-basic-api")
	connectionCh <- "Database connected..."
}

func FindOne(
	context context.Context,
	collection string,
	filter any,
	options ...*options.FindOneOptions,
) *mongo.SingleResult {
	return Db.Collection(collection).FindOne(context, filter, options...)
}

func UpdateOne(
	context context.Context,
	collection string,
	filter any,
	update any,
	options ...*options.FindOneAndUpdateOptions,
) *mongo.SingleResult {
	return Db.Collection(collection).FindOneAndUpdate(context, filter, update, options...)
}

func DeleteOne(
	context context.Context,
	collection string,
	filter any,
	options ...*options.FindOneAndDeleteOptions,
) *mongo.SingleResult {
	return Db.Collection(collection).FindOneAndDelete(context, filter, options...)
}

func InsertOne(
	context context.Context,
	collection string,
	document any,
	options ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {

	return Db.Collection(collection).InsertOne(context, document, options...)
}

func Find(
	context context.Context,
	collection string,
	filter any,
	options ...*options.FindOptions,
) (cur *mongo.Cursor, err error) {
	return Db.Collection(collection).Find(context, filter, options...)
}
