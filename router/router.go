package router

import (
	"github.com/Joshuafreemant/go-social/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RoutesSetup(app *fiber.App) {
	api := app.Group("/api", logger.New())
	// routes
	api.Get("/", controller.Index)
	auth := api.Group("/auth")
	users := api.Group("/users")
	auth.Post("/create-account", controller.CreateUsers)
	auth.Post("/login", controller.Login)
	auth.Post("/forgot-password", controller.ForgotPassword)
	auth.Post("/verify-otp", controller.VerifyOTP)
	auth.Post("/reset-password", controller.ResetPassword)
	users.Get("/", controller.GetUsers)

}
