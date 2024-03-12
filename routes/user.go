package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateAccount(ctx *gin.Context) {

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var existingUser models.User
	db.Db.Collection("users").
		FindOne(context.TODO(), bson.M{"email": user.Email}).
		Decode(&existingUser)

	if existingUser.Email == user.Email {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User with email already exists",
		})
		return
	}

	_, err := user.Insert()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account created",
		"data":    user,
	})
}

func Login(ctx *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := db.Db.Collection("users").
		FindOne(context.TODO(), bson.M{"email": credentials.Email})

	if err := result.Err(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	var user models.User
	if err := result.Decode(&user); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	if user.Email != credentials.Email && user.Password != credentials.Password {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
	})
}

func UpdateUser(ctx *gin.Context) {

}
