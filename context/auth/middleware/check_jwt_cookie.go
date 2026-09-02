package auth_middleware

import (
	"log"
	"os"
	"strings"

	"com.xalixcolabs.trading-card-album/context/user/application"
	"com.xalixcolabs.trading-card-album/database"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func CheckJwtCoockie(c fiber.Ctx) error {
	// Aceptar el JWT desde el header Authorization (Bearer) o desde la cookie.
	tokenString := ""
	if auth := c.Get(fiber.HeaderAuthorization); strings.HasPrefix(auth, "Bearer ") {
		tokenString = strings.TrimPrefix(auth, "Bearer ")
	}
	if tokenString == "" {
		tokenString = c.Cookies("jwt")
	}
	if tokenString == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sesion invalida",
		})
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
	user, err := user_application.GetUserById(database.DefaultQuerier(), id)
	if err != nil {
		log.Printf("[CheckJwtCoockie GetUserById] %s", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sesion invalida",
		})
	}
	c.Locals("session", user)
	return c.Next()
}
