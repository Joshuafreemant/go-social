package middleware

import (
	"github.com/Joshuafreemant/go-social/helpers" // Update the import path according to your project structure
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("jwt")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token missing"})
		}

		claims, err := helpers.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("user", claims)

		return c.Next()
	}
}
