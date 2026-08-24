package auth_resource

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"

	"com.xalixcolabs.trading-card-album/context/auth/application"
	"com.xalixcolabs.trading-card-album/context/auth/middleware"
	"com.xalixcolabs.trading-card-album/context/user/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterAuthResource(app *fiber.App) {
	apiV1 := app.Group("/api/v1/auth", logger.New())

	apiV1.Get("", auth)
	apiV1.Get("/me", auth_middleware.CheckJwtCoockie, me)
	apiV1.Get("/google/callback", callback)
}

// @Description	Get current session user
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Success		200  {object}   user_model.User
// @Router /api/v1/auth/me [get]
func me(c fiber.Ctx) {
	c.JSON(c.Locals("session").(user_model.User))
}

func auth(c fiber.Ctx) error {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	state := hex.EncodeToString(b)
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
	url := auth_application.ConfigGoogleAuth().AuthCodeURL(state)
	return c.Redirect().To(url)
}

func callback(c fiber.Ctx) error {
	cookieState := c.Cookies("oauth_state")
	queryState := c.Query("state")
	if cookieState == "" || cookieState != queryState {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "invalid OAuth state"})
	}
	token, err := auth_application.ConfigGoogleAuth().Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	response, err := auth_application.GetGoogleResponse(token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	admins := os.Getenv("TCA_ADMINS")
	isAdmin := strings.Contains(admins, response.Email)
	user, err := user_application.CreateUser(response.Email, isAdmin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: user.ID,
	})
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
		Domain:   "localhost",
	})
	return c.Redirect().To("/")
}
