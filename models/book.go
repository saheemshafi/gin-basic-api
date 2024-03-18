package models

import (
	"context"
	"time"

	"github.com/saheemshafi/gin-basic-api/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const BookCollection = "books"

type Book struct {
	Id          primitive.ObjectID   `json:"_id" bson:"_id"`
	Title       string               `json:"title" bson:"title" binding:"required"`
	Author      primitive.ObjectID   `json:"author" bson:"author"`
	Description string               `json:"description" bson:"description" binding:"required"`
	Cover       string               `json:"cover" bson:"cover"`
	Pages       []primitive.ObjectID `json:"pages" bson:"pages"`
	CreatedAt   primitive.DateTime   `json:"createdAt" bson:"createdAt"`
	UpdatedAt   primitive.DateTime   `json:"updatedAt" bson:"updatedAt"`
}

func (book *Book) Insert() (*mongo.InsertOneResult, error) {

	book.Id = primitive.NewObjectID()
	book.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	book.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return db.InsertOne(context.Background(), BookCollection, book)
}
