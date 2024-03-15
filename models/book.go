package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	Id          primitive.ObjectID   `json:"_id" bson:"_id"`
	Title       string               `json:"title" bson:"title" binding:"required"`
	Author      primitive.ObjectID   `json:"author" bson:"author"`
	Description string               `json:"description" bson:"description" binding:"required"`
	Cover       string               `json:"cover" bson:"cover" binding:"required"`
	Pages       []primitive.ObjectID `json:"pages" bson:"pages"`
}

type Page struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id"`
	Title   string             `json:"title" bson:"title" binding:"required"`
	Cover   string             `json:"cover" bson:"cover"`
	Content string             `json:"content" bson:"content"`
}
