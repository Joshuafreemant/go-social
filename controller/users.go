package controller

import (
	"context"
	"log"

	"github.com/Joshuafreemant/go-social/database"
	"github.com/Joshuafreemant/go-social/helpers"
	"github.com/Joshuafreemant/go-social/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Index(c *fiber.Ctx) error {
	return helpers.ResponseMsg(c, 200, "Api is running", nil)
}

func GetUsers(c *fiber.Ctx) error {
	ctx := context.Background()
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	var users []model.Users
	cursor, err := db.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		return helpers.ResponseMsg(c, 500, "Failed to fetch users", err.Error())
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return helpers.ResponseMsg(c, 500, "Failed to parse users", err.Error())
	}
	var usersResponse []model.Users
	for _, user := range users {
		userResponse := model.Users{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		usersResponse = append(usersResponse, userResponse)
	}

	return helpers.ResponseMsg(c, 200, "Fetched users successfully", usersResponse)
}
