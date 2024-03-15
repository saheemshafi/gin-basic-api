package models

import (
	"context"
	"time"

	"github.com/saheemshafi/gin-basic-api/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Page struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     string             `json:"title" bson:"title" binding:"required"`
	Cover     string             `json:"cover" bson:"cover"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func (page *Page) Insert() (*mongo.InsertOneResult, error) {

	page.Id = primitive.NewObjectID()
	page.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	page.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return db.Db.Collection("pages").InsertOne(context.TODO(), page)
}
