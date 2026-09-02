package album_participant_resource

import (
	"log"

	"com.xalixcolabs.trading-card-album/context/album_participant/application"
	"com.xalixcolabs.trading-card-album/context/album_participant/model/dto"
	"com.xalixcolabs.trading-card-album/context/auth/middleware"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RegisterAlbumParticipantResource(app *fiber.App) {
	apiV1 := app.Group(
		"/api/v1/album_participant",
		logger.New(),
		auth_middleware.CheckJwtCoockie,
	)

	apiV1.Post("", createAlbumParticipant)
}

// @Description	Create Album participant
// @Tags		Album Participant
// @Accept		json
// @Produce		json
// @Param		request body album_participant_dto.CreateAlbumParticipantRequest true "payload"
// @Success		200  {object}   album_participant_model.AlbumParticipant
// @Router /api/v1/album_participant [post]
// @Security BearerAuth
func createAlbumParticipant(c fiber.Ctx) error {
	session := c.Locals("session").(user_model.User)
	request := new(album_participant_dto.CreateAlbumParticipantRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}
	response, err := album_participant_application.CreateAlbumParticipant(database.DefaultQuerier(), session, *request)
	if err != nil {
		log.Default().Printf("[Error createAlbumParticipant] %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al crear album participant",
		})
	}
	return c.JSON(response)
}
