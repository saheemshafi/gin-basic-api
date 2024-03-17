package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateBook(ctx *gin.Context) {

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	var book models.Book

	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	book.Author = user.Id

	_, err := book.Insert()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book created",
		"data":    book,
	})
}

func AddPage(ctx *gin.Context) {

}

func UpdatePage(ctx *gin.Context) {

}

func UpdateBook(ctx *gin.Context) {

}

func DeleteBook(ctx *gin.Context) {

}

func GetBooks(ctx *gin.Context) {
	cursor, err := db.Db.Collection("books").Find(context.Background(), bson.M{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	var books []models.Book

	err = cursor.All(ctx.Request.Context(), &books)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve books",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Books retrieved",
		"data":    books,
	})
}

func GetBook(ctx *gin.Context) {

	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book id",
		})
		return
	}

	result := db.Db.Collection("books").FindOne(context.Background(), bson.M{"_id": bookId})

	if err := result.Err(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	var book models.Book

	if err := result.Decode(&book); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book retrieved",
		"data":    book,
	})
}
