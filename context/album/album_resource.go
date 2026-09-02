package album_resource

import (
	"encoding/json"
	"log"

	"com.xalixcolabs.trading-card-album/context/album/application"
	"com.xalixcolabs.trading-card-album/context/album/model/dto"
	"com.xalixcolabs.trading-card-album/context/album_participant/application"
	"com.xalixcolabs.trading-card-album/context/auth/middleware"
	"com.xalixcolabs.trading-card-album/context/card/application"
	"com.xalixcolabs.trading-card-album/context/card_pool/application"
	"com.xalixcolabs.trading-card-album/context/events"
	"com.xalixcolabs.trading-card-album/context/user/model"
	"com.xalixcolabs.trading-card-album/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/sse"
	"github.com/skip2/go-qrcode"
)

func RegisterAlbumResource(app *fiber.App) {
	apiv1 := app.Group(
		"/api/v1/album",
		logger.New(),
		auth_middleware.CheckJwtCoockie,
	)

	apiv1.Get("/", getAlbumsByUser)
	apiv1.Get("/:id", getAlbumById)
	apiv1.Get("/:id/join_qr", joinQr)
	apiv1.Post("", auth_middleware.CheckIsAdmin, createAlbum)
	apiv1.Get("/:id/card", getCardPoolByAlbumId)
	apiv1.Get("/:id/assigned_card", assignedCard)
	apiv1.Post("/new_card", registerNewCard)
	apiv1.Get("/:id/share_assigned_card", shareAssignedCard)
	apiv1.Get("/:id/qr_events", qrEvents)
}

// @Description	SSE: avisa cuando el QR compartido de este álbum fue escaneado
// @Tags		Album
// @Accept		json
// @Produce		text/event-stream
// @Param		id   path  string  true  "Album ID"
// @Success		200  {string}   string
// @Router /api/v1/album/{id}/qr_events [get]
func qrEvents(c fiber.Ctx) error {
	sse.New(sse.Config{
		Handler: func(c fiber.Ctx, stream *sse.Stream) error {
			session := c.Locals("session").(user_model.User)
			sub := events.Subscribe(session.ID)
			defer events.Unsubscribe(session.ID, sub)
			for {
				select {
				case <-stream.Done():
					return nil
				case <-sub.Closed:
					return nil
				case data := <-sub.Data:
					if err := stream.Event(sse.Event{Data: data}); err != nil {
						return err
					}
				}
			}
		},
	})(c)
	return nil
}

// @Description	Get albums by user
// @Tags		Album
// @Accept		json
// @Produce		json
// @Success		200 {array} sqlc.Album
// @Router /api/v1/album [get]
func getAlbumsByUser(c fiber.Ctx) error {
	session := c.Locals("session").(user_model.User)
	album, err := album_application.GetAlbumsByUser(database.DefaultQuerier(), session)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al listar albums",
		})
	}
	return c.JSON(album)
}

// @Description	Get album by id
// @Tags		Album
// @Accept		json
// @Produce		json
// @Param		id   path  string  true  "Album ID"
// @Success		200  {object}   album_model.Album
// @Router /api/v1/album/{id} [get]
func getAlbumById(c fiber.Ctx) error {
	id := c.Params("id")
	session := c.Locals("session").(user_model.User)
	album, err := album_application.GetAlbumById(database.DefaultQuerier(), session, id)
	if err != nil {
		log.Default().Printf("[Error getAlbumById] %s", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Álbum no encontrado",
		})
	}
	return c.JSON(album)
}

// @Description	Get join QR of an album (encode del código de invitación)
// @Tags		Album
// @Accept		json
// @Produce		image/png
// @Param		id   path  string  true  "Album ID"
// @Success		200  {file}   binary
// @Router /api/v1/album/{id}/join_qr [get]
func joinQr(c fiber.Ctx) error {
	id := c.Params("id")
	png, err := qrcode.Encode(id, qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al generar el QR",
		})
	}
	c.Set("Content-Type", "image/png")
	c.Set("Cache-Control", "no-store")
	return c.Send(png)
}

// @Description	Create album
// @Tags		Album
// @Accept		json
// @Produce		json
// @Param		request body album_dto.CreateAlbumRequest true "payload"
// @Success		200  {object}  album_model.Album
// @Router /api/v1/album [post]
func createAlbum(c fiber.Ctx) error {
	request := new(album_dto.CreateAlbumRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}
	album, err := album_application.CreateAlbum(database.DefaultQuerier(), *request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al crear album",
		})
	}
	return c.JSON(album)
}

// @Description	Share assigned card
// @Tags		Album
// @Accept		json
// @Produce		json,image/png
// @Param		id   path  string  true  "Album ID"
// @Param		qr   query  string  false  "qr view"
// @Success		200  {object}   album_participant_dto.ShareAssignedCardResponse
// @Router /api/v1/album/{id}/share_assigned_card [get]
func shareAssignedCard(c fiber.Ctx) error {
	isQr := c.Query("qr")
	id := c.Params("id")
	session := c.Locals("session").(user_model.User)
	response, err := album_participant_application.ShareAssignedCard(database.DefaultQuerier(), session, id)
	if err != nil {
		log.Default().Printf("[Error shareAssignedCard] %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al compartir tarjeta asignada",
		})
	}
	if isQr != "" {
		jsonPayload, err := json.Marshal(response)
		if err != nil {
			return c.JSON(response)
		}
		png, err := qrcode.Encode(string(jsonPayload), qrcode.Medium, 256)
		if err != nil {
			return c.JSON(response)
		}
		c.Set("Content-Type", "image/png")
		c.Set("Cache-Control", "no-store")
		return c.Send(png)
	}
	return c.JSON(response)
}

// @Description	Register new card
// @Tags		Album
// @Accept		json
// @Produce		json
// @Param		request body album_dto.RegisterCardRequest true "payload"
// @Success		200 {object} card_model.Card
// @Router /api/v1/album/new_card [post]
func registerNewCard(c fiber.Ctx) error {
	session := c.Locals("session").(user_model.User)
	request := new(album_dto.RegisterCardRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}
	response, err := album_application.RegisterCard(database.DefaultQuerier(), session, *request)
	if err != nil {
		log.Default().Printf("[Error registerNewCard] %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(response)
}

// @Description	Assigned card
// @Tags		Album
// @Accept		json
// @Produce		json
// @Param		id   path  string  true  "Album ID"
// @Success		200  {object}   card_model.Card
// @Router /api/v1/album/{id}/assigned_card [get]
func assignedCard(c fiber.Ctx) error {
	id := c.Params("id")
	session := c.Locals("session").(user_model.User)
	response, err := card_application.GetCardByAlbumId(database.DefaultQuerier(), session, id)
	if err != nil {
		log.Default().Printf("[Error shareAssignedCard] %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al compartir tarjeta asignada",
		})
	}
	return c.JSON(response)
}

// @Description	Card Pool
// @Tags		Album
// @Accept		json
// @Produce		json
// @Param		id   path  string  true  "Album ID"
// @Success		200  {array}   card_model.Card
// @Router /api/v1/album/{id}/card [get]
func getCardPoolByAlbumId(c fiber.Ctx) error {
	id := c.Params("id")
	session := c.Locals("session").(user_model.User)
	response, err := card_pool_application.GetMyCardsByAlbumId(database.DefaultQuerier(), session, id)
	if err != nil {
		log.Default().Printf("[Error shareAssignedCard] %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al compartir tarjeta asignada",
		})
	}
	return c.JSON(response)
}
