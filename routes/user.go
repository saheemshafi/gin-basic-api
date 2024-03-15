package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/models"
	"github.com/saheemshafi/gin-basic-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
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

	if !utils.ComparePasswordHashes(credentials.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials",
		})
		return
	}

	sessionTime := time.Now().Add(24 * time.Hour)
	token, err := utils.EncodeJWT(user.Id.Hex(), sessionTime)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Failed to login",
		})
		return
	}

	ctx.SetCookie(
		"token",
		token,
		int(time.Until(sessionTime).Seconds()),
		"/",
		"localhost",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
		"token":   token,
	})
}

func UpdateUser(ctx *gin.Context) {

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	var updates struct {
		Name string `binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, err := db.Db.Collection("users").UpdateByID(
		context.Background(),
		user.Id,
		bson.M{
			"$set": bson.M{
				"name":      updates.Name,
				"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)

	if result.ModifiedCount == 0 || err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Updated user details",
	})
}
