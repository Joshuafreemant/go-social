package router

import (
	"github.com/Joshuafreemant/go-social/controller"
	"github.com/Joshuafreemant/go-social/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

func RoutesSetup(app *fiber.App, db *mongo.Database) {
	api := app.Group("/api", logger.New())
	// routes
	api.Get("/", controller.Index)
	auth := api.Group("/auth")
	users := api.Group("/users")
	posts := api.Group("/posts")
	auth.Post("/create-account", func(c *fiber.Ctx) error {
		return controller.CreateUsers(c, db)
	})
	auth.Post("/login", func(c *fiber.Ctx) error {
		return controller.Login(c, db)
	})
	auth.Post("/forgot-password", controller.ForgotPassword)
	auth.Post("/verify-otp", controller.VerifyOTP)
	auth.Post("/reset-password", controller.ResetPassword)

	users.Get("/", controller.GetUsers)

	posts.Use(middleware.JWTMiddleware())
	posts.Post("/add-post", func(c *fiber.Ctx) error {
		return controller.CreatePost(c, db)
	})
	posts.Get("/", func(c *fiber.Ctx) error {
		return controller.GetAllPosts(c, db)
	})
	posts.Delete("/delete/:id", func(c *fiber.Ctx) error {
		return controller.DeletePost(c, db)
	})

}
