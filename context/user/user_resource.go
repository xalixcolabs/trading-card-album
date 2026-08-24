package user_resource

import (
	"com.xalixcolabs.trading-card-album/context/auth/middleware"
	"com.xalixcolabs.trading-card-album/context/user/application"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/context/user/model/dto"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RegisterUserResource(app *fiber.App) {
	apiV1 := app.Group(
		"/api/v1/user",
		logger.New(),
		auth_middleware.CheckJwtCoockie,
	)

	apiV1.Put("/:id", updateUserProfile)
}

// @Description	Update user by ID
// @Tags		User
// @Accept		json
// @Produce		json
// @Param       id path string true "User ID"
// @Param		request body user_dto.UpdateUserRequest true "payload"
// @Success		200  {object}   user_model.User
// @Router /api/v1/user/{id} [put]
func updateUserProfile(c fiber.Ctx) error {
	id := c.Params("id")
	session := c.Locals("session").(user_model.User)
	if session.ID != id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sin privilegios para actualizar usuario",
		})
	}
	request := new(user_dto.UpdateUserRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}
	response, err := user_application.UpdateUser(id, *request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al actualizar usuario",
		})
	}
	return c.JSON(response)
}
