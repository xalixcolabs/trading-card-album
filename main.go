package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"com.xalixcolabs.trading-card-album/context/album"
	"com.xalixcolabs.trading-card-album/context/album_participant"
	"com.xalixcolabs.trading-card-album/context/auth"
	"com.xalixcolabs.trading-card-album/context/contact"
	"com.xalixcolabs.trading-card-album/context/user"
	"com.xalixcolabs.trading-card-album/database"
	_ "com.xalixcolabs.trading-card-album/docs"
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
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
	app.Get("/swagger/*", swaggo.HandlerDefault)
	registerApiResources(app)
	registerFrontend(app, embedFrontend)
	log.Fatal(app.Listen(":" + port))
}

func registerApiResources(app *fiber.App) {
	auth_resource.RegisterAuthResource(app)
	user_resource.RegisterUserResource(app)
	album_resource.RegisterAlbumResource(app)
	album_participant_resource.RegisterAlbumParticipantResource(app)
	contact_resource.RegisterContactResource(app)
}

// registerFrontend sirve el frontend embebido (SPA de Nuxt). Si el path no
// corresponde a un archivo estático ni a una ruta de API/swagger, responde con
// el fallback SPA (200.html / index.html) para que el cliente de Nuxt resuelva
// la ruta (deep links como /album/:id).
func registerFrontend(app *fiber.App, frontend embed.FS) {
	frontendFS, err := fs.Sub(frontend, "webui/.output/public")
	if err != nil {
		log.Fatal(err)
	}

	app.Use("/*", func(c fiber.Ctx) error {
		// Solo se sirven assets en GET/HEAD.
		if c.Method() != fiber.MethodGet && c.Method() != fiber.MethodHead {
			return c.Next()
		}

		requestPath := c.Path()
		// Las rutas de API y swagger las resuelven sus propios handlers.
		if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") ||
			requestPath == "/swagger" || strings.HasPrefix(requestPath, "/swagger/") {
			return c.Next()
		}

		// Servir el archivo estático si existe (index, rutas prerenderizadas, assets).
		relative := strings.Trim(strings.TrimPrefix(requestPath, "/"), "/")
		if relative == "" {
			relative = "index.html"
		}
		if data, err := fs.ReadFile(frontendFS, relative); err == nil {
			return sendFile(c, relative, data)
		}
		if data, err := fs.ReadFile(frontendFS, path.Join(relative, "index.html")); err == nil {
			return sendFile(c, "index.html", data)
		}

		// Fallback SPA: Nuxt (cliente) resuelve la ruta.
		spaFile := "index.html"
		if _, err := fs.Stat(frontendFS, "200.html"); err == nil {
			spaFile = "200.html"
		}
		data, err := fs.ReadFile(frontendFS, spaFile)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
		c.Set(fiber.HeaderCacheControl, "no-cache")
		return c.Send(data)
	})
}

func sendFile(c fiber.Ctx, name string, data []byte) error {
	c.Type(path.Ext(name))
	return c.Send(data)
}