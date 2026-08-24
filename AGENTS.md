# AGENTS.md

Guía para agentes de IA que trabajan en **Trading Card Album (TCA)**.

## Resumen del proyecto

MVP de intercambio de tarjetas coleccionables (estilo álbum) para el DevFest 2026.
- Un **admin** crea un álbum con su baraja de tarjetas (cada tarjeta apunta a una imagen webp servida por un contenedor nginx).
- Cada **participante** se une a un álbum y recibe una **tarjeta asignada aleatoriamente** desde el `card_pool`.
- Puede compartir su tarjeta vía **QR**; al escanear el QR de otra persona, se registra un contacto y se desbloquea la tarjeta asignada del otro.
- Cada usuario gestiona su perfil y colección.

Monorepo con backend en **Go** y frontend en **Nuxt (Vue)** dentro de `webui/`. El frontend compilado se **embebe** en el binario de Go y se sirve en producción por el propio backend.

## Stack

**Backend**
- Go 1.25 (módulo `com.xalixcolabs.trading-card-album`)
- Framework HTTP: **Fiber v3** (`github.com/gofiber/fiber/v3`)
- DB: **SQLite** (driver `modernc.org/sqlite`, sin CGO)
- Migraciones: **dbmate** (`github.com/amacneil/dbmate/v2`) — ver sección Base de datos
- Queries tipadas: **sqlc** (genera `database/sqlc`)
- Auth: **Google OAuth2** + **JWT** (cookie `jwt`) + Swagger docs con **swaggo**
- Utilidades: `gonanoid` (IDs), `go-qrcode` (generación de QR)

**Frontend (`webui/`)**
- **Nuxt 4** (`ssr: false`), Vue 3
- **Tailwind CSS 4** + **daisyUI**
- Cliente API generado con **orval** a partir de `docs/swagger.yaml`
- `nuxt-toast`, `nuxt-qrcode`, `pouchdb`

## Arquitectura

Backend organizado por **contextos de dominio** (clean-ish architecture) en `context/`:

```
context/<contexto>/
  <contexto>_resource.go          # Router + handlers HTTP (Fiber)
  application/<caso_de_uso>.go    # Lógica de negocio / casos de uso
  model/<entidad>.go              # Entidad + constructor desde sqlc
  model/dto/<...>.go              # DTOs de request/response (JSON tags)
```

- **`*_resource.go`**: define el grupo de rutas, aplica middlewares y expone handlers. Contiene los comentarios de anotación swagger (`@Router`, `@Param`, `@Success`...).
- **`application/`**: casos de uso. Reciben el usuario de sesión (`c.Locals("session").(user_model.User)`), usan `sqlc.New(database.GetDatabase())` y devuelven entidades de `model/`.
- **`model/`**: entidades puras + `model/dto/` para request/response.
- Los handlers obtienen la sesión desde `c.Locals("session")` (inyectada por `CheckJwtCoockie`).

Contextos existentes: `user`, `auth`, `album`, `album_participant`, `card`, `card_pool`, `contact`.

### Rutas API (todas bajo `/api/v1`)

- **Auth** (`/auth`): `GET /`, `GET /me` (JWT), `GET /google/callback`
- **Album** (`/album`): `GET /`, `GET /:id` (solo participantes; devuelve solo las tarjetas recolectadas), `POST /` (admin), `GET /:id/card`, `GET /:id/assigned_card`, `POST /new_card`, `GET /:id/share_assigned_card?qr=`
- **Album participant** (`/album_participant`): `POST /`
- **Card** (`/card`): `POST /`
- **Contact** (`/contact`): `GET /`
- **User** (`/user`): `PUT /:id`

Swagger en `/swagger/*`.

## Base de datos

- Motor: **SQLite**. El archivo vive en `database/trading-card-album.sqlite3` (definido por `DATABASE_URL`).
- Las migraciones se gestionan con **dbmate**:
  - Archivos en `database/migrations/*.sql`, formato dbmate (`-- migrate:up` ... `-- migrate:down`).
  - En tiempo de ejecución se aplican automáticamente: `main.go` llama a `database.RunMigrations()`, que **embebe** las migraciones (`//go:embed migrations`) y las ejecuta con la librería `dbmate` (no el binario CLI).
  - El CLI de dbmate también está disponible (`make db-migrate`) para desarrollo.
- `database/schema.sql` es el **schema dump** consolidado (referenciado por `DBMATE_SCHEMA_FILE`). `sqlc` lo usa como esquema de entrada.
- **Flujo de cambio de esquema**: crear nueva migración dbmate en `database/migrations/` → regenerar `schema.sql` → `make gen-sql`.

### Convenciones de datos

- IDs: **nanoid** (texto, `gonanoid.New()`).
- Timestamps: **unix** como `INTEGER` (p. ej. `created_at INTEGER NOT NULL`).
- Flags booleanas: `INTEGER` (0/1), p. ej. `is_admin`, `is_drawn`.
- Tablas: `user`, `album`, `card`, `album_participant`, `contact`, `card_pool`, `user_card_collection`.
- Los `card_pool` se reponen cuando se agota (ver `AssignCard`): si no hay tarjetas disponibles, se resetea el pool.

## Código generado — ¡NO editar a mano!

- `database/sqlc/*.go` → generado por **sqlc** (fuente: `database/query/*.sql` + `schema.sql`).
- `docs/*` (`docs.go`, `swagger.json`, `swagger.yaml`) → generado por **swag** (a partir de los comentarios de anotación en los resources).
- `webui/app/services/*` y `webui/app/models/*` → generados por **orval** a partir de `docs/swagger.yaml`.
- `webui/.nuxt/`, `webui/.output/`, `webui/dist/` → artefactos de build.
- Modifica las **fuentes** y vuelve a generar (ver comandos abajo).

## Comandos

Todo vía `Makefile`:

```bash
make dev             # gen (sqlc+swagger) + backend (go run) + frontend (nuxt dev) en paralelo
make dev-backend     # go run main.go  (backend en :8080)
make dev-frontend    # cd webui && npm run dev  (frontend en :3000)
make build           # build-frontend + build-backend
make build-backend   # go build -o build/trading-card-album .
make build-frontend  # cd webui && npm run generate
make db-migrate      # dbmate migrate (CLI)
make gen-sql         # sqlc generate
make gen-swagger     # swag init
make gen             # gen-sql + gen-swagger
```

Flujo típico de cambios:
1. Modificar queries en `database/query/*.sql` → `make gen-sql`.
2. Cambiar anotaciones swagger / handlers → `make gen-swagger`.
3. Para el frontend, `npm run dev` ejecuta `orval` automáticamente (regenera servicios/modelos desde el swagger).

## Variables de entorno (`.env`)

- `APP_PORT` (por defecto `8080`)
- `DATABASE_URL` (`sqlite:database/trading-card-album.sqlite3`)
- `DBMATE_MIGRATIONS_DIR`, `DBMATE_SCHEMA_FILE`
- `TCA_ADMINS` — lista de emails separada por comas que serán admin
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL`
- `JWT_SECRET`

> `GOOGLE_REDIRECT_URL` apunta al **origen de la app**: `http://localhost:3000/api/v1/auth/google/callback` en desarrollo (Nuxt proxya `/api` hacia el backend) y al host real del binario en producción.

> `.env` contiene **secretos reales** (OAuth y JWT). No debe versionarse ni exponerse en logs ni en código.

## Flujo de autenticación

1. `GET /api/v1/auth` → redirige a Google OAuth con `state` en cookie `oauth_state`.
2. `GET /api/v1/auth/google/callback` → valida el state, intercambia el código, obtiene email, crea/actualiza usuario (admin si está en `TCA_ADMINS`) y firma un **JWT** que se guarda en cookie `jwt` (no HTTPOnly).
3. Middleware `CheckJwtCoockie` (en `context/auth/middleware/check_jwt_cookie.go`) valida el JWT en cada request protegido y deja el usuario en `c.Locals("session")`.
4. `CheckIsAdmin` protege rutas de admin (ej. crear álbum).

## Convenciones de código

- Paquetes de `context/*` usan sufijos `_resource`, `_application`, `_model`, `_dto` (ej. `album_resource`, `album_application`, `album_model`, `album_dto`).
- Handlers Fiber: devuelven `error`; en errores usan `c.Status(...).JSON(fiber.Map{"message": ...})`.
- Constructores de entidad desde sqlc: `New<Entidad>FromSqlc<Entidad>` (ver `model/card.go`, `model/User.go`).
- Campos sensibles no se exponen en JSON: `Secret json:"-"` (en `AlbumParticipant`) y `IsAdmin json:"-"` (en `User`).
- Los mensajes de UI/comentarios están en **español**.
- No añadir comentarios al código salvo que se pida.

## Frontend — notas

- Servicios generados por orval usan el mutator `customFetch` (`webui/app/services/CustomFetch.ts`) que añade `credentials: 'include'` (para las cookies) y `baseURL` desde `runtimeConfig.public.apiBase`.
- `runtimeConfig.public.apiBase`: siempre `''` (mismo origen). En dev Nuxt (puerto `3000`) proxya `/api` hacia el backend (`vite.server.proxy` → `http://localhost:8080`); en producción Fiber sirve la SPA y el API juntos.
- **`app.buildAssetsDir: 'assets'`**: los assets salen en `/assets` (no `/_nuxt`) porque `go:embed` excluye directorios que empiezan con `_` — sin esto el frontend embebido queda roto.
- Fiber sirve el frontend embebido con **fallback SPA** (`main.go` → `registerFrontend`): si el path no es un archivo estático ni `/api/*` ni `/swagger/*`, responde `200.html`/`index.html` para que el cliente de Nuxt resuelva rutas como `/album/:id`.
- Composables: `useApiData` (envuelve `useAsyncData`), `useProfile` (estado global del perfil).
- Páginas principales: `index.vue`, `album/[id].vue`, `contactos.vue`, `profile.vue`, `login.vue`. Componentes: `Card`, `AppSheet`, `TabBar`, `MyQrModal`, `ScanQrModal`, `GoogleMark`.
- Middleware global `auth.global.ts`: redirige a `/login` si no existe la cookie `jwt`.
