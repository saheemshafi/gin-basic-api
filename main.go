package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	Mongo.connectDatabase()
	defer func() {
		log.Fatalln(Mongo.Client.Disconnect(context.Background()))
	}()

	server := Server{
		Address: ":5000",
	}

	server.Initialize()
}
