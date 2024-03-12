package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/routes"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db.Connect()
	defer db.Db.Client().Disconnect(context.TODO())

	app := gin.Default()

	routes.Register(app)

	log.Fatal(app.Run(":5000"))
}
