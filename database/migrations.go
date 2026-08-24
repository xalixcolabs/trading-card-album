package database

import (
	"embed"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/sqlite"
)

//go:embed migrations
var embedMigrations embed.FS

func RunMigrations() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL env variable is required")
	}

	if strings.HasPrefix(databaseURL, "sqlite:") {
		path := strings.TrimPrefix(databaseURL, "sqlite:")
		dir := filepath.Dir(path)
		if dir != "." && dir != "/" {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatalf("Error al crear el directorio de la base de datos: %v", err)
			}
		}
	}

	tempDir, err := os.MkdirTemp("", "dbmate-migrations-*")
	if err != nil {
		log.Fatalf("Error al crear directorio temporal: %v", err)
	}
	defer os.RemoveAll(tempDir)

	entries, err := fs.ReadDir(embedMigrations, "migrations")
	if err != nil {
		log.Fatalf("Error al leer migraciones embebidas: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := embedMigrations.ReadFile(filepath.Join("migrations", entry.Name()))
		if err != nil {
			log.Fatalf("Error al leer archivo %s: %v", entry.Name(), err)
		}
		err = os.WriteFile(filepath.Join(tempDir, entry.Name()), content, 0644)
		if err != nil {
			log.Fatalf("Error al escribir archivo temporal %s: %v", entry.Name(), err)
		}
	}

	u, err := url.Parse(databaseURL)
	if err != nil {
		log.Fatalf("Error al parsear URL de la base de datos: %v", err)
	}

	dbm := dbmate.New(u)
	dbm.AutoDumpSchema = false
	dbm.MigrationsDir = []string{tempDir}

	err = dbm.Migrate()
	if err != nil {
		log.Fatalf("Error ejecutando migraciones con dbmate: %v", err)
	}

	log.Println("✅ Migraciones de dbmate aplicadas con éxito usando la librería nativa")
}
