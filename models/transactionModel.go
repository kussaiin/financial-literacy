package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID                primitive.ObjectID `bson:"_id"`
	Name              *string            `json:"name" bson:"name" validate:"required"`
	Price             *string            `json:"price" bson:"price" validate:"required"`
	Data              *string            `json:"data" bson:"data" validate:"required"`
	TypeOfTransaction *string            `json:"type_of_transaction" bson:"type_of_transaction" validate:"required"`
	Comment           *string            `json:"comment" bson:"comment"`
	Category          Category           `json:"category" bson:"category" `
	TransactionId     string             `json:"transaction_id" bson:"transaction_id"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}
