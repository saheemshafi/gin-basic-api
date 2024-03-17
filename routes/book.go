package routes

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/gin-gonic/gin"
	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/models"
	"github.com/saheemshafi/gin-basic-api/utils"
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
			"message": "You can't add page to this book",
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

	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book id",
		})
		return
	}

	pageId, err := primitive.ObjectIDFromHex(ctx.Param("pageId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page id",
		})
		return
	}

	var pageInfo struct {
		Title   string
		Content string
	}

	if err := ctx.ShouldBindJSON(&pageInfo); err != nil {
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
			"message": "You can't update page from this book",
		})
		return
	}

	updates := bson.M{
		"$set": bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	updateMap := map[string]string{
		"title":   pageInfo.Title,
		"content": pageInfo.Content,
	}

	for key, value := range updateMap {
		if strings.TrimSpace(value) != "" {
			updates["$set"].(bson.M)[key] = value
		}
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := db.Db.Collection("pages").FindOneAndUpdate(
		context.Background(),
		bson.M{
			"_id": pageId,
		},
		updates,
		options,
	)

	if err := result.Err(); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Page not found",
			})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	var page models.Page
	result.Decode(&page)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Updated page",
		"data":    page,
	})
}

func DeletePage(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book id",
		})
		return
	}

	pageId, err := primitive.ObjectIDFromHex(ctx.Param("pageId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page id",
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
			"message": "You can't delete page from this book",
		})
		return
	}

	page := db.Db.Collection("pages").FindOneAndDelete(context.Background(), bson.M{
		"_id": pageId,
	})

	if err := page.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Page not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete",
		})
		return
	}

	result := db.Db.Collection("books").FindOneAndUpdate(
		context.Background(),
		bson.M{
			"_id": bookId,
		},
		bson.M{
			"$pull": bson.M{
				"pages": pageId,
			},
		},
	)

	if err := result.Err(); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to remove page",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Removed page from book",
	})
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

	var updates = bson.M{
		"$set": bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	updateMap := map[string]string{
		"title":       bookInfo.Title,
		"description": bookInfo.Description,
	}

	for key, value := range updateMap {
		if strings.TrimSpace(value) != "" {
			updates["$set"].(bson.M)[key] = value
		}
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := db.Db.Collection("books").FindOneAndUpdate(
		context.Background(),
		bson.M{
			"_id": bookId,
		},
		updates,
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

		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Book not found",
			})
			return
		}

		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
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

func ChangeBookCover(ctx *gin.Context) {
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
			"message": "You can't change this book's cover",
		})
		return
	}

	formFile, err := ctx.FormFile("cover")
	log.Println(formFile.Filename, formFile.Size)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	tempCover, err := formFile.Open()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to open file",
		})
		return
	}

	defer tempCover.Close()
	uploadResult, err := utils.UploadFile(tempCover)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Failed to upload file",
		})
		return
	}

	updateResult, err := db.Db.Collection("books").UpdateByID(context.Background(), bookId, bson.M{
		"$set": bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
			"cover":     uploadResult.PublicID,
		},
	})

	if err != nil || updateResult.ModifiedCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to change cover",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Changed book cover",
		"data":    uploadResult.PublicID,
	})

	if book.Cover != "" {
		result, err := utils.DeleteFile(book.Cover, api.Image)

		if result.Error.Message != "" || err != nil {
			log.Println(result.Error.Message)
		}
	}

}

func ChangePageCover(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid book id",
		})
		return
	}

	pageId, err := primitive.ObjectIDFromHex(ctx.Param("pageId"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page id",
		})
		return
	}

	existingBook := db.Db.Collection("books").FindOne(context.Background(), bson.M{
		"_id": bookId,
	})

	if err := existingBook.Err(); err != nil {
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
			"message": "You can't change this page's cover",
		})
		return
	}

	formFile, err := ctx.FormFile("cover")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	tempCover, err := formFile.Open()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to open file",
		})
		return
	}

	defer tempCover.Close()
	uploadResult, err := utils.UploadFile(tempCover)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Failed to upload file",
		})
		return
	}

	updateResult := db.Db.Collection("pages").FindOneAndUpdate(
		context.Background(),
		bson.M{
			"_id": pageId,
		},
		bson.M{
			"$set": bson.M{
				"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
				"cover":     uploadResult.PublicID,
			},
		})

	if err := updateResult.Err(); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Page not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to change cover",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Changed page cover",
		"data":    uploadResult.PublicID,
	})

	var page models.Page
	updateResult.Decode(&page)

	if page.Cover != "" {
		result, err := utils.DeleteFile(page.Cover, api.Image)

		log.Println("Deleting", page.Cover)
		if result.Error.Message != "" || err != nil {
			log.Println(result.Error.Message)
		}
	}

}
