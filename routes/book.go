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
		utils.WriteResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	book.Author = user.Id
	book.Pages = []primitive.ObjectID{}

	_, err := book.Insert()

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteResponse(ctx, http.StatusCreated, "Book created", book)
}

func AddPage(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	var page models.Page

	if err := ctx.ShouldBindJSON(&page); err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	existingBook := db.FindOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		})

	if err := existingBook.Err(); err != nil {
		log.Println(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "You can't add page to this book")
		return
	}

	insertResult, err := page.Insert()

	if err != nil {
		log.Println(err)
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to add page")
		return
	}

	result := db.UpdateOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		},
		bson.M{
			"$push": bson.M{
				"pages": insertResult.InsertedID,
			},
			"$set": bson.M{
				"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)

	if err := result.Err(); err != nil {
		log.Println(err)
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to add page")
		return
	}

	utils.WriteResponse(ctx, http.StatusCreated, "Added page to book", page)
}

func UpdatePage(ctx *gin.Context) {

	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	pageId, err := primitive.ObjectIDFromHex(ctx.Param("pageId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid page id")
		return
	}

	var pageInfo struct {
		Title   string
		Content string
	}

	if err := ctx.ShouldBindJSON(&pageInfo); err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	existingBook := db.FindOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		})

	if err := existingBook.Err(); err != nil {
		log.Println(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "You can't update page from this book")
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
	result := db.UpdateOne(
		context.Background(),
		models.PageCollection,
		bson.M{
			"_id": pageId,
		},
		updates,
		options,
	)

	if err := result.Err(); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Page not found")
			return
		}

		log.Print(err)
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var page models.Page
	result.Decode(&page)

	utils.WriteResponse(ctx, http.StatusOK, "Updated page", page)
}

func DeletePage(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	pageId, err := primitive.ObjectIDFromHex(ctx.Param("pageId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid page id")
		return
	}

	existingBook := db.FindOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		})

	if err := existingBook.Err(); err != nil {
		log.Println(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "You can't delete page from this book")
		return
	}

	page := db.DeleteOne(
		context.Background(),
		models.PageCollection,
		bson.M{
			"_id": pageId,
		})

	if err := page.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Page not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to delete")
		return
	}

	result := db.UpdateOne(
		context.Background(),
		models.BookCollection,
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
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to remove page")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Removed page from book")
}

func UpdateBook(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	var bookInfo struct {
		Title       string
		Description string
	}

	if err := ctx.ShouldBindJSON(&bookInfo); err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	existingBook := db.FindOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		})

	if err := existingBook.Err(); err != nil {
		log.Println(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "You can't update this book")
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
	result := db.UpdateOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		},
		updates,
		options,
	)

	if err := result.Err(); err != nil {
		log.Println(err)

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Updated book")
}

func DeleteBook(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	existingBook := db.FindOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		})

	if err := existingBook.Err(); err != nil {
		log.Println(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "You can't delete this book")
		return
	}

	result := db.DeleteOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		})

	if err := result.Err(); err != nil {
		log.Print(err)
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Book deleted")
}

func GetBooks(ctx *gin.Context) {

	cursor, err := db.Find(
		context.Background(),
		models.BookCollection,
		bson.M{})

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var books []models.Book

	err = cursor.All(context.Background(), &books)

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Books retrieved", books)
}

func GetBook(ctx *gin.Context) {

	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	result := db.FindOne(
		context.Background(),
		models.BookCollection,
		bson.M{"_id": bookId})

	if err := result.Err(); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		log.Println(err)
		utils.WriteResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var book models.Book

	if err := result.Decode(&book); err != nil {
		log.Println(err.Error())
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Book retrieved", book)
}

func ChangeBookCover(ctx *gin.Context) {
	bookId, err := primitive.ObjectIDFromHex(ctx.Param("bookId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	existingBook := db.FindOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		})

	if err := existingBook.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "You can't change this book's cover")
		return
	}

	formFile, err := ctx.FormFile("cover")
	log.Println(formFile.Filename, formFile.Size)

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	tempCover, err := formFile.Open()

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to open file")
		return
	}

	defer tempCover.Close()
	uploadResult, err := utils.UploadFile(tempCover)

	if err != nil {
		log.Println(err)
		utils.WriteResponse(ctx, http.StatusOK, "Failed to upload file")
		return
	}

	updateResult := db.UpdateOne(
		context.Background(),
		models.BookCollection,
		bson.M{
			"_id": bookId,
		},
		bson.M{
			"$set": bson.M{
				"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
				"cover":     uploadResult.PublicID,
			},
		})

	if updateResult.Err(); err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to change cover")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Changed book cover", uploadResult.PublicID)

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
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid book id")
		return
	}

	pageId, err := primitive.ObjectIDFromHex(ctx.Param("pageId"))

	if err != nil {
		utils.WriteResponse(ctx, http.StatusBadRequest, "Invalid page id")
		return
	}

	existingBook := db.Db.Collection("books").FindOne(context.Background(), bson.M{
		"_id": bookId,
	})

	if err := existingBook.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteResponse(ctx, http.StatusNotFound, "Book not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var book models.Book
	existingBook.Decode(&book)

	userFromCtx, _ := ctx.Get("user")
	user := userFromCtx.(models.User)

	if book.Author != user.Id {
		utils.WriteResponse(ctx, http.StatusUnauthorized, "You can't change this page's cover")
		return
	}

	formFile, err := ctx.FormFile("cover")

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	tempCover, err := formFile.Open()

	if err != nil {
		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to open file")
		return
	}

	defer tempCover.Close()
	uploadResult, err := utils.UploadFile(tempCover)

	if err != nil {
		log.Println(err)
		utils.WriteResponse(ctx, http.StatusOK, "Failed to upload file")
		return
	}

	updateResult := db.UpdateOne(
		context.Background(),
		models.PageCollection,
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
			utils.WriteResponse(ctx, http.StatusNotFound, "Page not found")
			return
		}

		utils.WriteResponse(ctx, http.StatusInternalServerError, "Failed to change cover")
		return
	}

	utils.WriteResponse(ctx, http.StatusOK, "Changed page cover", uploadResult.PublicID)

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
