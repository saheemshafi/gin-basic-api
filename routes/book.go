package routes

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	book.Pages = []primitive.ObjectID{}

	_, err := book.Insert()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Book created",
		"data":    book,
	})
}

func AddPage(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book id",
		})
		return
	}

	var page models.Page

	if err := ctx.ShouldBindJSON(&page); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	insertResult, err := page.Insert()

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add page",
		})
		return
	}

	result, err := db.Db.Collection("books").UpdateByID(
		context.Background(),
		bookId,
		bson.M{
			"$push": bson.M{
				"pages": insertResult.InsertedID,
			},
			"$set": bson.M{
				"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)

	if err != nil || result.ModifiedCount == 0 {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add page",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Added page to book",
		"data":    page,
	})
}

func UpdatePage(ctx *gin.Context) {

}

func UpdateBook(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book id",
		})
		return
	}

	var bookInfo struct {
		Title       string
		Description string
	}

	if err := ctx.ShouldBindJSON(&bookInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	existingBook := db.Db.Collection("books").FindOne(context.Background(), bson.M{
		"_id": bookId,
	})

	if err := existingBook.Err(); err != nil {
		log.Println(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Book not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You can't update this book",
		})
		return
	}

	var updates = bson.M{}

	if strings.TrimSpace(bookInfo.Title) != "" {
		updates["title"] = bookInfo.Title
	}

	if strings.TrimSpace(bookInfo.Description) != "" {
		updates["description"] = bookInfo.Description
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := db.Db.Collection("books").FindOneAndUpdate(
		context.Background(),
		bson.M{
			"_id": bookId,
		},
		bson.M{
			"$set": updates,
		},
		options,
	)

	if err := result.Err(); err != nil {
		log.Println(err)

		ctx.JSON(http.StatusInternalServerError, bson.M{
			"message": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Updated book",
	})
}

func DeleteBook(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book id",
		})
		return
	}

	existingBook := db.Db.Collection("books").FindOne(context.Background(), bson.M{
		"_id": bookId,
	})

	if err := existingBook.Err(); err != nil {
		log.Println(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Book not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You can't delete this book",
		})
		return
	}

	result := db.Db.Collection("books").FindOneAndDelete(context.Background(), bson.M{
		"_id": bookId,
	})

	if err := result.Err(); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book deleted",
	})
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

	err = cursor.All(context.Background(), &books)

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
