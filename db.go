package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	Client       *mongo.Client
	Database     *mongo.Database
	databaseName string
}

var Mongo = MongoDb{
	databaseName: "gin-basic-api",
}

func (db *MongoDb) connectDatabase() {

	log.Println("Connecting to database")

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatalln("MONGODB_URI not found")
	}

	options := options.Client().ApplyURI(uri).SetAppName("Gin Basic Api")

	client, err := mongo.Connect(context.Background(), options)

	if err != nil {
		log.Fatalln(err)
	}

	Mongo.Client = client
	Mongo.Database = client.Database(Mongo.databaseName)
	log.Println("Database connected...")
}
