package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/routes"
	"github.com/saheemshafi/gin-basic-api/utils"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	/*
		The channel is only used to get familiar with it's implementation.
		Log statements within the functions can also be used.

		Also go routines can be fired so both db and cld start trying to connect at
		same time and then notify back or log.Fatal when failed
	*/
	connectionCh := make(chan string, 6)

	db.Connect(connectionCh)
	defer db.Db.Client().Disconnect(context.TODO())

	utils.InitializeCloudinary(connectionCh)
	/*
		Channel needs to be closed first else range will go into infinite loop.
		Buffered channel is used so it won't get into a deadlock after there is
		nothing to put into it
	*/
	close(connectionCh)

	for msg := range connectionCh {
		log.Println(msg)
	}

	app := gin.Default()

	routes.Register(app)

	var port string = os.Getenv("PORT")

	if port == "" {
		log.Print("Port not specified in env... Defaulting to 5000")
		port = "5000"
	}
	log.Fatal(app.Run(fmt.Sprintf(":%v", port)))
}
