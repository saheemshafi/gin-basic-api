package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Address string
}

func (server *Server) Initialize() {

	app := gin.New()

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Server is ready")
	})

	log.Fatal(app.Run(server.Address))
}
