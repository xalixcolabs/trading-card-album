package admin_resource

import (
	"com.xalixcolabs.trading-card-album/context/admin/application"
	"com.xalixcolabs.trading-card-album/context/admin/model/dto"
	"com.xalixcolabs.trading-card-album/context/auth/middleware"
	"com.xalixcolabs.trading-card-album/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func RegisterAdminResource(app *fiber.App) {
	apiV1 := app.Group(
		"/api/v1/admin",
		logger.New(),
		auth_middleware.CheckJwtCoockie,
		auth_middleware.CheckIsAdmin,
	)

	apiV1.Get("/overview", getOverview)
	apiV1.Get("/albums", getAlbums)
	apiV1.Delete("/albums/:id", deleteAlbum)
	apiV1.Get("/users", getUsers)
	apiV1.Put("/users/:id/role", updateUserRole)
	apiV1.Get("/cards", getCards)
	apiV1.Delete("/cards/:id", deleteCard)
}

// @Description	Get platform overview metrics
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Success		200  {object}   admin_dto.Overview
// @Router /api/v1/admin/overview [get]
func getOverview(c fiber.Ctx) error {
	overview, err := admin_application.GetOverview(database.DefaultQuerier())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al obtener métricas",
		})
	}
	return c.JSON(overview)
}

// @Description	List all albums with stats
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Success		200  {array}   admin_dto.Album
// @Router /api/v1/admin/albums [get]
func getAlbums(c fiber.Ctx) error {
	albums, err := admin_application.ListAlbums(database.DefaultQuerier())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al listar álbumes",
		})
	}
	return c.JSON(albums)
}

// @Description	Delete an album
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Param		id   path  string  true  "Album ID"
// @Success		200  {object}   admin_dto.Message
// @Router /api/v1/admin/albums/{id} [delete]
func deleteAlbum(c fiber.Ctx) error {
	id := c.Params("id")
	if err := admin_application.DeleteAlbum(database.DefaultQuerier(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al eliminar álbum",
		})
	}
	return c.JSON(fiber.Map{"message": "Álbum eliminado"})
}

// @Description	List all users
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Success		200  {array}   admin_dto.User
// @Router /api/v1/admin/users [get]
func getUsers(c fiber.Ctx) error {
	users, err := admin_application.ListUsers(database.DefaultQuerier())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al listar usuarios",
		})
	}
	return c.JSON(users)
}

// @Description	Update user admin role
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Param		id   path  string  true  "User ID"
// @Param		request body admin_dto.UpdateUserRoleRequest true "payload"
// @Success		200  {object}   admin_dto.User
// @Router /api/v1/admin/users/{id}/role [put]
func updateUserRole(c fiber.Ctx) error {
	id := c.Params("id")
	request := new(admin_dto.UpdateUserRoleRequest)
	if err := c.Bind().JSON(request); err != nil {
		return err
	}
	user, err := admin_application.UpdateUserRole(database.DefaultQuerier(), id, request.IsAdmin)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al actualizar rol",
		})
	}
	return c.JSON(user)
}

// @Description	List all cards
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Success		200  {array}   card_model.Card
// @Router /api/v1/admin/cards [get]
func getCards(c fiber.Ctx) error {
	cards, err := admin_application.ListCards(database.DefaultQuerier())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al listar tarjetas",
		})
	}
	return c.JSON(cards)
}

// @Description	Delete a card
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Param		id   path  string  true  "Card ID"
// @Success		200  {object}   admin_dto.Message
// @Router /api/v1/admin/cards/{id} [delete]
func deleteCard(c fiber.Ctx) error {
	id := c.Params("id")
	if err := admin_application.DeleteCard(database.DefaultQuerier(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al eliminar tarjeta",
		})
	}
	return c.JSON(fiber.Map{"message": "Tarjeta eliminada"})
}