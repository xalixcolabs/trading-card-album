package auth_middleware

import (
	"log"
	"os"

	"com.xalixcolabs.trading-card-album/context/user/application"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func CheckJwtCoockie(c fiber.Ctx) error {
	cookieJwt := c.Cookies("jwt")
	if cookieJwt == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sesion invalida",
		})
	}
	token, err := jwt.Parse(cookieJwt, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Printf("[CheckJwtCoockie] %s", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sesion invalida",
		})
	}
	if !token.Valid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sesion invalida",
		})
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("[CheckJwtCoockie GetSubject] %s", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sesion invalida",
		})
	}
	user, err := user_application.GetUserById(id)
	if err != nil {
		log.Printf("[CheckJwtCoockie GetUserById] %s", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sesion invalida",
		})
	}
	c.Locals("session", user)
	return c.Next()
}
