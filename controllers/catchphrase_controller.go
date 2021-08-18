package controllers

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/Aceix/movie-catchprase/config"
	"github.com/Aceix/movie-catchprase/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCatchphrasesBy(ctx *fiber.Ctx) error {
	catchphraseColl := config.MI.DB.Collection("catchphrases")
	ctx1, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var catchphrases []models.Catchphrase

	filter := bson.M{}
	findOptions := options.Find()

	if s := ctx.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"movieName": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"catchphrase": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limitVal, _ := strconv.Atoi(ctx.Query("limit", "10"))
	limit := int64(limitVal)

	total, _ := catchphraseColl.CountDocuments(ctx1, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := catchphraseColl.Find(ctx1, filter, findOptions)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrases not found",
			"error":   err,
		})
	}
	defer cursor.Close(ctx1)

	for cursor.Next(ctx1) {
		var catchphrase models.Catchphrase
		err = cursor.Decode(&catchphrase)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "An error occured whiles retrieving results",
				"error":   err,
			})
		}
		catchphrases = append(catchphrases, catchphrase)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"data":     catchphrases,
		"total":    total,
		"page":     page,
		"lastPage": last,
		"limit":    limit,
	})

}

func GetCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var catchphrase models.Catchphrase
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := catchphraseCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrase Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&catchphrase)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrase Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    catchphrase,
		"success": true,
	})
}

func AddCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	catchphrase := new(models.Catchphrase)

	if err := c.BodyParser(catchphrase); err != nil {
		log.Println(err)
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := catchphraseCollection.InsertOne(ctx, catchphrase)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrase failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Catchphrase inserted successfully",
	})

}

func UpdateCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	catchphrase := new(models.Catchphrase)

	if err := c.BodyParser(catchphrase); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrase not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": catchphrase,
	}
	_, err = catchphraseCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrase failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Catchphrase updated successfully",
	})
}

func DeleteCatchphrase(c *fiber.Ctx) error {
	catchphraseCollection := config.MI.DB.Collection("catchphrases")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrase not found",
			"error":   err,
		})
	}
	_, err = catchphraseCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Catchphrase failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Catchphrase deleted successfully",
	})
}
