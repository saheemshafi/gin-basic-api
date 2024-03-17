package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func Connect() {

	log.Println("Connecting to database")

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatalln("MONGODB_URI not found")
	}

	options := options.Client().ApplyURI(uri).SetAppName("Gin Basic Api")

	client, _ := mongo.Connect(context.Background(), options)

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalln(err)
	}

	Db = client.Database("gin-basic-api")
	log.Println("Database connected...")
}
