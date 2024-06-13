package helpers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type resMessage struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseMsg(c *fiber.Ctx, code int, msg string, data interface{}) error {
	resPonse := &resMessage{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	return c.Status(code).JSON(resPonse)
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
