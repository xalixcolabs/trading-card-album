package main

import (
	"embed"
	"io/fs"
	"log"
	"os"

	"com.xalixcolabs.trading-card-album/context/album"
	"com.xalixcolabs.trading-card-album/context/album_participant"
	"com.xalixcolabs.trading-card-album/context/auth"
	"com.xalixcolabs.trading-card-album/context/user"
	"com.xalixcolabs.trading-card-album/database"
	_ "com.xalixcolabs.trading-card-album/docs"
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
)

//go:embed webui/.output/public/*
var embedFrontend embed.FS

// @title          Trading Card Album API
// @version        1.0
// @description    MVP para intercambio de tarjetas coleccionables.
func main() {
	godotenv.Load()

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	database.RunMigrations()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))
	app.Get("/swagger/*", swaggo.HandlerDefault)
	registerApiResources(app)
	frontendFS, err := fs.Sub(embedFrontend, "webui/.output/public")
	if err != nil {
		log.Fatal(err)
	}
	app.Use("/*", static.New("", static.Config{
		FS: frontendFS,
	}))
	log.Fatal(app.Listen(":" + port))
}

func registerApiResources(app *fiber.App) {
	auth_resource.RegisterAuthResource(app)
	user_resource.RegisterUserResource(app)
	album_resource.RegisterAlbumResource(app)
	album_participant_resource.RegisterAlbumParticipantResource(app)
}
