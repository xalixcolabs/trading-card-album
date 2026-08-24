package contact_resource

import (
	"com.xalixcolabs.trading-card-album/context/auth/middleware"
	"com.xalixcolabs.trading-card-album/context/contact/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RegisterContactResource(app *fiber.App) {
	apiV1 := app.Group(
		"/api/v1/contact",
		logger.New(),
		auth_middleware.CheckJwtCoockie,
	)

	apiV1.Get("/", listContacts)
}

// @Description	List contacts of the current user
// @Tags		Contact
// @Accept		json
// @Produce		json
// @Success		200  {array}   contact_dto.Contact
// @Router /api/v1/contact [get]
func listContacts(c fiber.Ctx) error {
	session := c.Locals("session").(user_model.User)
	contacts, err := contact_application.ListContacts(database.DefaultQuerier(), session)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al listar contactos",
		})
	}
	return c.JSON(contacts)
}