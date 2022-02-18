package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/user/financial-literacy/database"
	"github.com/user/financial-literacy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var transactionCollection *mongo.Collection = database.OpenCollection(database.Client, "transaction")
var validate = validator.New()

func GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := transactionCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var allTransactions []bson.M

		if err := result.All(ctx, &allTransactions); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allTransactions)
	}
}

func GetTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var transaction models.Transaction
		transactionId := c.Param("transaction_id")

		err := transactionCollection.FindOne(ctx, bson.M{"transaction_id": transactionId}).Decode(&transaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching transaction"})
			return
		}
		c.JSON(http.StatusOK, transaction)
	}
}

func CreateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var transaction models.Transaction
		var category models.Category

		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(transaction)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := categoriesCollection.FindOne(ctx, bson.M{"title": transaction.Category.Title}).Decode(&category)
		defer cancel()
		if err != nil {
			log.Println("category was not found")
			transaction.Category.ID = primitive.NewObjectID()
			transaction.Category.CategoryId = transaction.Category.ID.Hex()
			_, err := categoriesCollection.InsertOne(ctx, transaction.Category)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "new category was not created"})
				return
			}
			log.Println("new category was created")
		}

		transaction.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		transaction.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		transaction.ID = primitive.NewObjectID()
		transaction.TransactionId = transaction.ID.Hex()

		result, err := transactionCollection.InsertOne(ctx, transaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction was not created"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var transaction models.Transaction

		transactionId := c.Param("transaction_id")

		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var updateObj primitive.D

		if transaction.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: transaction.Name})
		}

		if transaction.Price != nil {
			updateObj = append(updateObj, bson.E{"price", transaction.Price})
		}

		if transaction.Data != nil {
			updateObj = append(updateObj, bson.E{"data", transaction.Data})
		}

		if transaction.TypeOfTransaction != nil {
			updateObj = append(updateObj, bson.E{"type_of_transaction", transaction.TypeOfTransaction})
		}

		if transaction.Comment != nil {
			updateObj = append(updateObj, bson.E{"comment", transaction.Comment})
		}

		updateObj = append(updateObj, bson.E{"category", transaction.Category})


		transaction.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", transaction.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"transaction_id": transactionId}

		result, err := transactionCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction was not updated"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func DeleteTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		transactionId := c.Param("transaction_id")
		filter := bson.M{"transaction_id": transactionId}

		res, err := transactionCollection.DeleteOne(ctx, filter)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction was not deleted"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
