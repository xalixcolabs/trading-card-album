# AGENTS.md

Guía para agentes de IA que trabajan en **Trading Card Album (TCA)**. Monorepo: backend Go + frontend Nuxt embebido en el binario Go.

## Resumen

MVP de intercambio de tarjetas coleccionables (DevFest 2026).
- Un **admin** crea un álbum con su baraja de tarjetas (imágenes webp servidas por nginx desde `cards/`).
- Cada **participante** se une a un álbum y recibe una tarjeta aleatoria del `card_pool`.
- Comparte su tarjeta vía **QR**; al escanear el QR de otra persona se registra un contacto y se desbloquea su tarjeta.
- El frontend de Nuxt se compila y se **embebe** (`//go:embed webui/.output/public/*`); el binario lo sirve con fallback SPA.

## Stack

- Go 1.25 (módulo `com.xalixcolabs.trading-card-album`), Fiber v3, SQLite.
- Queries con **sqlc**, migraciones con **dbmate** (librería embebida, no CLI), Swagger con **swaggo**.
- Frontend en `webui/`: Nuxt 4 (`ssr:false`) + Tailwind CSS 4 + `@phosphor-icons/vue`; cliente API generado por **orval** desde `docs/swagger.yaml`.

## Backend — arquitectura

`context/<contexto>/`:
```
<contexto>_resource.go        # grupo de rutas Fiber + handlers + anotaciones swagger
application/<caso_de_uso>.go  # lógica de negocio
model/<entidad>.go            # entidad + New<Entidad>FromSqlc<Entidad>
model/dto/<...>.go            # DTOs JSON
```

Contextos: `auth`, `user`, `album`, `album_participant`, `card`, `card_pool`, `contact`, `admin`, `events`.
- `card` y `card_pool` **no** exponen resource (se usan internamente desde `album`).
- `events` es un hub pub/sub en memoria para SSE (`context/events/hub.go`); `main.go` lo limpia cada 24h.

**Inyección de dependencias (clave):** los casos de uso en `application/` reciben `database.Querier` como **primer parámetro** y **no** abren la DB. Los resources pasan `database.DefaultQuerier()`; los tests pasan `database/queriermock.Querier`.
- Añadir/cambiar una query obliga a tocar: interfaz `database/querier.go` + mock `database/queriermock/queriermock.go` (mantenido a mano) + tests.

Los handlers obtienen la sesión con `c.Locals("session").(user_model.User)` (inyectada por `CheckJwtCoockie`).

### Rutas API (bajo `/api/v1`)

- **Auth** `/auth`: `GET ""` (redirige a Google), `GET /me` (JWT), `GET /google/callback`.
- **User** `/user`: `PUT /:id` (solo el propio usuario).
- **Album** `/album`: `GET /`, `GET /:id`, `GET /:id/join_qr`, `POST ""` (admin), `GET /:id/card` (tarjetas recolectadas por el usuario), `GET /:id/assigned_card`, `POST /new_card`, `GET /:id/share_assigned_card?qr=`, `GET /:id/qr_events` (SSE).
- **Album participant** `/album_participant`: `POST ""`.
- **Contact** `/contact`: `GET /`.
- **Admin** `/admin` (JWT + `CheckIsAdmin`): `GET /overview`, `GET|POST /albums`, `PUT|DELETE /albums/:id`, `GET /albums/:id/cards`, `GET /users`, `GET /users/:id`, `PUT /users/:id/role`, `GET|POST /cards`, `PUT|DELETE /cards/:id`.

Swagger en `/swagger/*`.

## Base de datos

- Migraciones dbmate en `database/migrations/*.sql` (`-- migrate:up` / `-- migrate:down`). Se **embeben** (`//go:embed migrations`) y se aplican al arrancar (`database.RunMigrations()`), que exige `DATABASE_URL`.
- `database/schema.sql` es el dump consolidado que lee sqlc; incluye la lista de versiones en `schema_migrations`.
- Flujo de cambio de esquema: nueva migración → actualizar `database/schema.sql` (+ versión) → `make gen-sql`.
- IDs: nanoid (TEXT). Timestamps: unix `INTEGER`. Booleanos: `INTEGER` 0/1.
- El `card_pool` se repone al agotarse: `AssignCard` reinicia el pool si `GetRandomAvailableCard` no devuelve filas.

**Gotcha:** `database.GetDatabase()` usa la ruta fija `database/trading-card-album.sqlite3` e **ignora `DATABASE_URL`**; solo las migraciones leen `DATABASE_URL`. Cambiar la ruta de la DB requiere editar `database/database.go`.

## Tests

- `make test` → `go test ./...`. No hay tests de frontend.
- Un solo paquete/test: `go test ./context/album/application/ -run TestCreateAlbum`.
- Los casos de uso se testean con `database/queriermock.Querier` (campos `<Metodo>Fn`; sin configurar devuelven valor cero y `nil`). Añade éxito + errores para casos nuevos.

## Código generado — ¡NO editar a mano!

- `database/sqlc/*.go` ← sqlc (`database/query/*.sql` + `database/schema.sql`).
- `docs/*` (`docs.go`, `swagger.json`, `swagger.yaml`) ← swag (anotaciones en los `*_resource.go`).
- `webui/app/services/**` y `webui/app/models/**` ← orval (`docs/swagger.yaml`). Excepción: `webui/app/services/CustomFetch.ts` es el mutator y **sí** se edita a mano.
- `webui/.nuxt/`, `webui/.output/`, `webui/dist/` (symlink a `.output/public`): artefactos de build.

## Comandos

```bash
make dev            # make gen + backend (go run, :8080) + frontend (nuxt dev, :3000) en paralelo
make dev-backend    # go run main.go
make dev-frontend   # cd webui && npm run dev   (ejecuta orval y luego nuxt dev)
make build          # build-frontend (npm run generate) + build-backend
make build-backend  # go build -o build/trading-card-album .
make test           # go test ./...
make gen            # gen-sql (sqlc) + gen-swagger (swag)
make db-migrate     # dbmate migrate (CLI, opcional)
make docker         # empaqueta build/ local en imagen (requiere make build antes)
make docker-full    # compila frontend + backend dentro de Docker
```

Orden típico tras un cambio:
1. Queries: `database/query/*.sql` → actualizar `database/querier.go` + mock → `make gen-sql`.
2. Handlers/anotaciones swagger → `make gen-swagger`.
3. Frontend: `npm run dev` regenera orval; manual: `cd webui && npx orval`.

**Build gotcha (CGO):** el driver sqlite de dbmate usa `mattn/go-sqlite3`, así que compilar requiere CGO/compilador C. Las queries en runtime usan `modernc.org/sqlite` (sin CGO). `Dockerfile.full` compila con CGO.

**Build gotcha (embed):** `go build` requiere que exista `webui/.output/public`; hay que compilar el frontend antes. Por eso `make build` va frontend → backend.

## Auth

1. `GET /api/v1/auth` guarda `state` en cookie `oauth_state` y redirige a Google.
2. `/auth/google/callback` valida el state, crea/actualiza usuario (admin si el email aparece en `TCA_ADMINS`, comparado con `strings.Contains`), firma un JWT HS256 (subject = user ID) en la cookie `jwt` (`Secure`, **no** HTTPOnly) y redirige a `/`.
3. `CheckJwtCoockie` (así, con el typo) acepta el token desde el header `Authorization: Bearer ...` **o** la cookie `jwt`, carga el usuario y lo deja en `c.Locals("session")`.
4. `CheckIsAdmin` exige `session.IsAdmin != 0`.

Las cookies van con `Secure: true`; el desarrollo local asume `localhost`.

## Convenciones

- Paquetes de `context/*` con sufijos `_resource`, `_application`, `_model`, `_dto`; constructores `New<Entidad>FromSqlc<Entidad>`.
- Handlers Fiber devuelven `error`; errores con `c.Status(...).JSON(fiber.Map{"message": ...})`.
- Campos sensibles fuera del JSON: `Secret json:"-"` (AlbumParticipant). `IsAdmin` **sí** se expone.
- Comentarios y mensajes de UI en **español**; evita comentarios innecesarios.
- Commits convencionales (`feat:`, `fix:`, ...) **sin emojis**, mensajes en inglés; ramas `feat/|fix/|refactor/`. Antes de un PR: `make gen`, `make test`, `cd webui && npm run generate`.

## Frontend (`webui/`)

- Diseño **mobile-only** (shell de teléfono, tokens en `app/assets/css/main.css`, bottom sheets `AppSheet`, safe areas).
- Páginas: `index.vue`, `album/[id].vue`, `contactos.vue`, `profile.vue`, `login.vue`, `admin/*`. Middlewares: `auth.global.ts` (redirige a `/login` sin cookie `jwt`) y `admin.ts`.
- Composables: `useApiData` (envuelve `useAsyncData`), `useProfile`.
- Los servicios de orval usan el mutator `customFetch` (`credentials: 'include'`, `baseURL` desde `runtimeConfig.public.apiBase`).
- `apiBase` es `''` (mismo origen): en dev Nuxt proxya `/api` → `http://localhost:8080`; en prod Fiber sirve SPA + API juntos.
- **`app.buildAssetsDir: 'assets'`**: los assets deben salir en `/assets` porque `go:embed` ignora directorios que empiezan con `_` (`_nuxt` rompería el embed).
- Fiber sirve el frontend embebido con fallback SPA (`main.go` → `registerFrontend`): paths que no son archivo, ni `/api/*`, ni `/swagger/*` devuelven `200.html`/`index.html` para deep links como `/album/:id`.
- Swag no resuelve tipos no importados: si un `@Success` referencia un tipo, impórtalo y úsalo.

## Entorno (`.env`)

`APP_PORT` (8080), `DATABASE_URL`, `DBMATE_MIGRATIONS_DIR`, `DBMATE_SCHEMA_FILE`, `TCA_ADMINS` (emails separados por comas), `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`, `JWT_SECRET`.

`GOOGLE_REDIRECT_URL` apunta al origen de la app: `http://localhost:3000/api/v1/auth/google/callback` en dev, host real en prod. `.env` contiene secretos reales y está gitignoreado.
