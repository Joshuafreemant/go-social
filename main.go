package main

import (
	"log"
	"os"

	"github.com/Joshuafreemant/go-social/database"
	"github.com/Joshuafreemant/go-social/helpers"
	"github.com/Joshuafreemant/go-social/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading env")
	}
	// app.Use(cors.New())
	file, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	app := fiber.New()
	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	app.Static("/images", "./images")
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		Output: file, // Write logs to the file
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Replace "*" with your allowed origins
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	router.RoutesSetup(app, db)
	// MONGO_URL := os.Getenv("MONGO_URL")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	app.Use(func(c *fiber.Ctx) error {
		return helpers.ResponseMsg(c, 404, "NotFound", nil)
	})
	log.Fatal(app.Listen(":" + PORT))
}
