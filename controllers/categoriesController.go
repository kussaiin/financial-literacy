package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
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

var categoriesCollection *mongo.Collection = database.OpenCollection(database.Client, "categories")

func GetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := categoriesCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var allCategories []bson.M

		if err := result.All(ctx, &allCategories); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allCategories)
	}
}

func GetCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var category models.Category
		categoryId := c.Param("category_id")

		err := categoriesCollection.FindOne(ctx, bson.M{"category_id": categoryId}).Decode(&category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching category"})
			return
		}
		c.JSON(http.StatusOK, category)
	}
}


func CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var category models.Category

		if err := c.BindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(category)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		category.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		category.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		category.ID = primitive.NewObjectID()
		category.CategoryId = category.ID.Hex()

		result, err := categoriesCollection.InsertOne(ctx, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "category was not created"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var category models.Category

		categoryId := c.Param("category_id")

		if err := c.BindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var updateObj primitive.D

		if category.Title != nil {
			updateObj = append(updateObj, bson.E{Key: "title", Value: category.Title})
		}

		category.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", category.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"category_id": categoryId}

		result, err := categoriesCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "category was not updated"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func DeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		categoryId := c.Param("category_id")
		filter := bson.M{"category_id": categoryId}

		res, err := categoriesCollection.DeleteOne(ctx, filter)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "category was not deleted"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
