package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/models"
	"github.com/saheemshafi/gin-basic-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Authorize(ctx *gin.Context) {
	cookie, err := ctx.Cookie("token")

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You are not logged in",
		})
		ctx.Abort()
		return
	}

	token, err := utils.DecodeJWT(cookie)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	userId, _ := primitive.ObjectIDFromHex(token["jti"].(string))

	options := options.FindOne().SetProjection(bson.M{"password": 0})
	result := db.Db.Collection("users").FindOne(context.Background(), bson.M{"_id": userId}, options)

	log.Println(result)

	if err := result.Err(); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized",
		})
		ctx.Abort()
		return
	}

	var user models.User
	result.Decode(&user)

	ctx.Set("user", user)
	ctx.Next()
}
