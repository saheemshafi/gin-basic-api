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
		utils.WriteResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser models.User

	db.FindOne(
		context.Background(),
		models.UserCollection,
		bson.M{"email": user.Email},
	).Decode(&existingUser)

	if existingUser.Email == user.Email {
		utils.WriteResponse(ctx, http.StatusBadRequest, "User with email already exists")
		return
	}

	_, err := user.Insert()

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = ""
	utils.WriteResponse(ctx, http.StatusOK, "Account created", user)
}

func Login(ctx *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result := db.FindOne(
		context.Background(),
		models.UserCollection,
		bson.M{"email": credentials.Email},
	)

	if err := result.Err(); err != nil {
		utils.WriteResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	var user models.User

	if err := result.Decode(&user); err != nil {
		log.Println(err.Error())
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if !utils.ComparePasswordHashes(credentials.Password, user.Password) {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	sessionTime := time.Now().Add(24 * time.Hour)
	token, err := utils.EncodeJWT(user.Id.Hex(), sessionTime)

	if err != nil {
		utils.WriteResponse(ctx, http.StatusOK, "Failed to login")
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

	utils.WriteResponse(ctx, http.StatusOK, "Logged in", token)

}

func UpdateUser(ctx *gin.Context) {

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	var updates struct {
		Name string `binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&updates); err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result := db.UpdateOne(
		context.Background(),
		models.UserCollection,
		bson.M{
			"_id": user.Id,
		},
		bson.M{
			"$set": bson.M{
				"name":      updates.Name,
				"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)

	if err := result.Err(); err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Updated user details")
}
