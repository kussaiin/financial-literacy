package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Category struct {
	ID         primitive.ObjectID `bson:"_id"`
	Title      *string            `json:"title" bson:"title" validate:"required"`
	CategoryId string             `json:"category_id" bson:"category_id"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}
