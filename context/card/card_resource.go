package card_resource

import (
	"com.xalixcolabs.trading-card-album/context/auth/middleware"
	"com.xalixcolabs.trading-card-album/context/card/application"
	"com.xalixcolabs.trading-card-album/context/card/model/dto"
	"com.xalixcolabs.trading-card-album/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RegisterCardResource(app *fiber.App) {
	apiv1 := app.Group(
		"/api/v1/card",
		logger.New(),
		auth_middleware.CheckJwtCoockie,
	)

	apiv1.Post("", CreateCard)
}

// @Description	Create card
// @Tags		Card
// @Accept		json
// @Produce		json
// @Param		request body card_dto.CreateCardRequest true "payload"
// @Success		200 {object} card_model.Card
// @Router /api/v1/card [post]
func CreateCard(c fiber.Ctx) error {
	request := new(card_dto.CreateCardRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}
	card, err := card_application.CreateCard(database.DefaultQuerier(), *request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al crear tarjeta",
		})
	}
	return c.JSON(card)
}
