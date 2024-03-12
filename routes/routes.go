package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine) {

	v1 := app.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	users.POST("/create-account", CreateAccount)
	users.POST("/login", Login)
	users.PUT("/{id}", UpdateUser)
}
