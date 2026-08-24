.PHONY: build dev dev-backend dev-frontend test db-migrate gen-sql gen-swagger gen

# ==============================================================================
# ENTORNO DE DESARROLLO
# ==============================================================================

dev:
	@echo "🚀 Iniciando entorno de desarrollo..."
	@make gen
	@make -j 2 dev-backend dev-frontend

dev-backend:
	@echo "📦 Iniciando Backend (Go)..."
	@go run main.go

dev-frontend:
	@echo "🎨 Iniciando Frontend (Nuxt)..."
	@cd webui && npm run dev -- --host


# ==============================================================================
# BUILD
# ==============================================================================

build:
	@echo "🚀 Construyendo proyecto..."
	@make build-frontend
	@make build-backend

build-backend:
	@go build -o build/trading-card-album .

build-frontend:
	@cd webui && npm run generate

# ==============================================================================
# TESTS
# ==============================================================================

test:
	@echo "🧪 Ejecutando tests..."
	@go test ./...

# ==============================================================================
# UTILIDADES Y GENERACIÓN DE CÓDIGO
# ==============================================================================

db-migrate:
	@echo "🗄️  Ejecutando migraciones con dbmate..."
	@dbmate migrate

gen-sql:
	@echo "🧬 Generando código de base de datos con sqlc..."
	@sqlc generate

gen-swagger:
	@echo "📝 Generando documentación de Swagger..."
	@swag init

gen: gen-sql gen-swagger
	@echo "✅ ¡Todo el código ha sido generado exitosamente!"