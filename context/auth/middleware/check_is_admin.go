package auth_middleware

import (
	"log"

	"com.xalixcolabs.trading-card-album/context/user/model"
	"github.com/gofiber/fiber/v3"
)

func CheckIsAdmin(c fiber.Ctx) error {
	user := c.Locals("session").(user_model.User)
	if user.IsAdmin == 0 {
		log.Printf("[CheckIsAdmin] no admin")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "No eres administrador",
		})
	}
	return c.Next()
}
