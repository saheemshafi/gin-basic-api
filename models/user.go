package models

import (
	"context"
	"time"

	"github.com/saheemshafi/gin-basic-api/db"
	"github.com/saheemshafi/gin-basic-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required,email"`
	Password  string             `json:"-" bson:"password" binding:"required"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}

func (user *User) Insert() (*mongo.InsertOneResult, error) {

	hash, err := utils.HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	user.Id = primitive.NewObjectID()
	user.Password = hash
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return db.Db.Collection("users").InsertOne(context.TODO(), user)
}
