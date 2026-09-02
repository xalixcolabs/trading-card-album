# Contribuyendo a Trading Card Album

Gracias por querer contribuir. Estas guías buscan que los cambios sean pequeños, revisables y consistentes con el resto del proyecto.

## Configuración del entorno

1. Clona el repositorio y crea tu `.env` (ver `README.md` → "Cómo empezar").
2. Instala dependencias:

   ```bash
   go mod download          # backend
   cd webui && npm install   # frontend
   ```

3. Requisitos de herramientas: **sqlc**, **swag**, y opcionalmente **dbmate** (CLI).
4. Corre `make dev` y verifica que todo arranca.

## Ejecutar con Docker

Hay dos formas de obtener la imagen:

- **`make docker`**: compila en local (`make build`) y **empaqueta el binario de `build/`** en una imagen **distroless** (`Dockerfile`). Solo necesita el binario ya compilado.
- **`make docker-full`**: **compila todo dentro de Docker** (frontend + backend) con el multi-stage `Dockerfile.full`. No requiere herramientas locales.

En ambos casos las variables de entorno se pasan en tiempo de ejecución (el `.env` no se copia a la imagen).

```bash
make docker       # local: make build + empaquetado distroless
make docker-full  # todo dentro de Docker
```

Ejemplo de `docker-compose.yaml` (el backend expone la API + el frontend embebido en `:8080`.

```yaml
services:
  app:
    image: trading-card-album:latest
    container_name: tca_app
    ports:
      - "8080:8080"
    environment:
      APP_PORT: "8080"
      DATABASE_URL: "sqlite:data/trading-card-album.sqlite3"
      TCA_ADMINS: "admin@example.com"
      JWT_SECRET: "cambia-este-secreto"
      GOOGLE_CLIENT_ID: "xxxx.apps.googleusercontent.com"
      GOOGLE_CLIENT_SECRET: "xxxx"
      # El origen de la app: apunta al host/puerto donde se accede al contenedor.
      GOOGLE_REDIRECT_URL: "http://localhost:8080/api/v1/auth/google/callback"
    volumes:
      # Persistencia de la base de datos SQLite.
      - tca-data:/app/data
    restart: unless-stopped

volumes:
  tca-data:
```
- `GOOGLE_REDIRECT_URL` debe coincidir con el origen desde el que se abre la app (en este ejemplo `http://localhost:8080`).

## Flujo de trabajo git

- Trabaja en una **rama** descriptiva (`feat/...`, `fix/...`, `refactor/...`).
- Usa **commits convencionales** (`feat`, `fix`, `refactor`, `chore`, `docs`, `test`, `build`) **sin emojis**.
- Haz **commits pequeños y por cambio lógico**; no juntes varias cosas en un commit gigante.
- Mensaje en **inglés** (o en español si el proyecto lo usa para el contexto del cambio), conciso y descriptivo:

  ```
  feat: add admin endpoint to list album cards
  fix: keep user name on re-login
  ```

- Antes de commitear, revisa `git status` y `git diff` para no subir archivos generados que no correspondan ni secretos (`.env` está gitignoreado).

## Código generado — ¡no editar a mano!

Estos archivos se generan y **no deben modificarse manualmente**:

- `database/sqlc/*.go` → **sqlc** (fuente: `database/query/*.sql` + `database/schema.sql`)
- `docs/*` → **swag** (a partir de las anotaciones en los `*_resource.go`)
- `webui/app/services/*` y `webui/app/models/*` → **orval** (a partir de `docs/swagger.yaml`)
- `webui/.nuxt`, `webui/.output`, `webui/dist` → artefactos de build

Modifica las **fuentes** y regenera (ver secciones siguientes).

## Backend

### Estructura de un contexto

Sigue el patrón de `context/<contexto>/`:

- `<contexto>_resource.go`: grupo de rutas, middlewares, handlers y anotaciones swagger (`@Router`, `@Param`, `@Success`).
- `application/`: casos de uso. Reciben `database.Querier` (inyectada) y el usuario de sesión (`user_model.User`).
- `model/`: entidades + constructores `New<Entidad>FromSqlc<Entidad>`.
- `model/dto/`: DTOs de request/response.

### Convenciones

- Paquetes con sufijos `_resource`, `_application`, `_model`, `_dto`.
- Handlers devuelven `error`; errores con `c.Status(...).JSON(fiber.Map{"message": ...})`.
- No expongas campos sensibles en JSON (`Secret json:"-"`, etc.).
- Comentarios/mensajes de UI en **español**; no añadas comentarios salvo que aporten.
- Para exponer un caso de uso al frontend, regístralo en el resource **y** pasa `database.DefaultQuerier()`.

### Cambiar una query

1. Edita `database/query/<archivo>.sql`.
2. Si la query la usa la capa `application`, añádela a la interfaz `database.Querier` y al mock `database/queriermock`.
3. `make gen-sql`.
4. Añade/actualiza los **unit tests** de la capa `application` (`*_test.go`).

### Cambiar el esquema

1. Crea una migración dbmate en `database/migrations/<timestamp>_<nombre>.sql`:

   ```sql
   -- migrate:up
   ALTER TABLE ...;

   -- migrate:down
   ALTER TABLE ...;
   ```

2. Actualiza `database/schema.sql` (dump que usa sqlc) e incluye la versión en `schema_migrations`.
3. `make gen-sql`.
4. Aplica la migración a tu DB local (se aplica sola al arrancar) y verifica.

### Endpoints y swagger

- Añade el handler y sus anotaciones swagger en el `*_resource.go`, registra la ruta y `make gen-swagger`.
- Para el frontend, regenera el cliente: `cd webui && npx orval` (o `npm run dev`, que lo ejecuta).
- Swag no resuelve tipos no importados: si el `@Success` referencia un tipo, impórtalo y úsalo.

## Frontend

- Páginas en `webui/app/pages/`, componentes en `webui/app/components/`, composables en `webui/app/composables/`.
- Diseño **mobile-only**: shell de teléfono, tokens de `app/assets/css/main.css`, bottom sheets con `AppSheet`, safe areas.
- Servicios/modelos los genera **orval**: si cambias una respuesta, regenera el cliente antes de usarlo.
- Mantén los estados: loading (skeleton), empty y error en cada pantalla.
- Usa iconos de `@phosphor-icons/vue` (no SVGs a mano).

## Tests

- Ejecuta: `make test` (`go test ./...`).
- Los casos de uso se testean con `database/queriermock` (cada método tiene un `<X>Fn` configurable).
- Añade tests para cada caso de uso nuevo o cambiado (éxito + errores).

## Antes de abrir un PR

1. `make gen` para regenerar sqlc + swagger (si aplica).
2. `make test` y `cd webui && npm run generate` en verde.
3. Revisa que no queden dependencias sin usar (`go mod tidy`, `npm prune`).
4. Verifica el diff: solo los archivos intencionados.

## Reportar bugs y sugerencias

- Abre un issue describiendo: qué esperabas, qué pasó, y cómo reproducirlo (pasos + entorno).
- Para preguntas de diseño, menciona el contexto (pantalla/flujo afectado).