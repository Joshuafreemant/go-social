package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Joshuafreemant/go-social/helpers"
	"github.com/Joshuafreemant/go-social/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePost(c *fiber.Ctx, db *mongo.Database) error {
	posts := new(model.Post)
	var ctx = context.Background()

	claims := c.Locals("user").(*helpers.Claims)
	userEmail := claims.Email

	if err := c.BodyParser(posts); err != nil {
		return helpers.ResponseMsg(c, 400, err.Error(), nil)
	}
	posts.UserId = userEmail
	imageDir := "./images"
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create image directory"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read form data"})
	}

	files := form.File["images"]
	for _, file := range files {
		filePath := filepath.Join(imageDir, file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image"})
		}

		log.Printf("Saved file: %s", filePath)

		posts.Images = append(posts.Images, fmt.Sprintf("/images/%s", file.Filename))
	}

	posts.CreatedAt = time.Now()
	posts.UpdatedAt = time.Now()

	result, err := db.Collection("posts").InsertOne(ctx, posts)
	if err != nil {
		return helpers.ResponseMsg(c, 500, "Failed to create post", err.Error())
	}

	return helpers.ResponseMsg(c, 200, "Post created", result.InsertedID)
}

func DeletePost(c *fiber.Ctx, db *mongo.Database) error {
	postID := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return helpers.ResponseMsg(c, 400, "Invalid post ID", err.Error())
	}

	filter := bson.M{"_id": objectID}

	var ctx = context.Background()

	result, err := db.Collection("posts").DeleteOne(ctx, filter)
	if err != nil {
		return helpers.ResponseMsg(c, 500, "Failed to delete post", err.Error())
	}

	if result.DeletedCount == 0 {
		return helpers.ResponseMsg(c, 404, "Post not found", nil)
	}

	return helpers.ResponseMsg(c, 200, "Post deleted", postID)
}

func GetAllPosts(c *fiber.Ctx, db *mongo.Database) error {
	var ctx = context.Background()

	var posts []model.Post

	result, err := db.Collection("posts").Find(ctx, bson.M{})
	if err != nil {
		return helpers.ResponseMsg(c, 500, "Failed to retrieve posts", err.Error())
	}
	defer result.Close(ctx)

	if err = result.All(ctx, &posts); err != nil {
		return helpers.ResponseMsg(c, 500, "Failed to decode posts", err.Error())
	}

	return helpers.ResponseMsg(c, 200, "Posts retrieved successfully", posts)
}
