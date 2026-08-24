# Trading Card Album (TCA)

MVP de intercambio de tarjetas coleccionables (estilo álbum) para el **DevFest 2026**.

- Un **admin** crea un álbum con su baraja de tarjetas.
- Cada **participante** se une a un álbum y recibe una **tarjeta asignada aleatoriamente** desde el `card_pool`.
- Comparte su tarjeta vía **QR**; al escanear el QR de otro participante se registra un **contacto** y se desbloquea su tarjeta en tu colección.
- Cada usuario gestiona su perfil, su colección y sus contactos.

El frontend (Nuxt) se **embebe en el binario de Go** y se sirve junto con la API desde un único ejecutable.

---

## Stack

**Backend** (`/`)
- **Go 1.25** (módulo `com.xalixcolabs.trading-card-album`)
- **Fiber v3** (`github.com/gofiber/fiber/v3`)
- **SQLite** vía `modernc.org/sqlite` (sin CGO)
- Migraciones con **dbmate**, queries tipadas con **sqlc**
- Auth: **Google OAuth2** + **JWT** en cookie
- QR con `go-qrcode`, Swagger con **swaggo**

**Frontend** (`webui/`)
- **Nuxt 4** (`ssr: false`), Vue 3
- **Tailwind CSS 4** + tema propio (tokens en `app/assets/css/main.css`)
- Iconos **@phosphor-icons/vue**, fuentes auto-hospedadas (Space Grotesk + JetBrains Mono)
- Cliente API generado con **orval** desde `docs/swagger.yaml`
- `nuxt-toast`, `nuxt-qrcode`

---

## Arquitectura

Backend organizado por **contextos de dominio** (clean-ish architecture) en `context/`:

```
context/<contexto>/
  <contexto>_resource.go        # Router + handlers HTTP (Fiber) + anotaciones swagger
  application/<caso_de_uso>.go  # Lógica de negocio / casos de uso
  model/<entidad>.go            # Entidad + constructor desde sqlc
  model/dto/<...>.go            # DTOs request/response
```

- Los casos de uso reciben una `database.Querier` (interfaz implementada por `*sqlc.Queries`) inyectada por parámetro, lo que permite testear con mocks (`database/queriermock`).
- La sesión se obtiene desde `c.Locals("session")` (inyectada por `CheckJwtCoockie`).

Contextos: `user`, `auth`, `album`, `album_participant`, `card`, `card_pool`, `contact`, `admin`.

### Estructura del proyecto

```
context/                 # Contextos de dominio
database/
  migrations/            # Migraciones dbmate (-- migrate:up / -- migrate:down)
  query/                 # Queries fuente de sqlc
  sqlc/                  # Código generado por sqlc (¡no editar!)
  schema.sql             # Dump consolidado del esquema (entrada de sqlc)
  querier.go             # Interfaz Querier + DefaultQuerier()
  queriermock/           # Mock de Querier para tests
docs/                    # Swagger generado (swag)
webui/                   # Frontend Nuxt
main.go                  # Entrada: API + frontend embebido con fallback SPA
Makefile                 # Comandos
```

---

## Cómo empezar

Requisitos: **Go 1.25+**, **Node 20+**, **sqlc**, **swag**, y opcionalmente **dbmate** (CLI).

### 1. Variables de entorno (`.env`)

Copia el ejemplo y complétalo:

```bash
cp .env.example .env   # si existe; si no, crea .env con las variables de abajo
```

| Variable | Descripción |
| --- | --- |
| `APP_PORT` | Puerto del backend (por defecto `8080`) |
| `DATABASE_URL` | `sqlite:database/trading-card-album.sqlite3` |
| `DBMATE_MIGRATIONS_DIR` | `database/migrations` |
| `DBMATE_SCHEMA_FILE` | `database/schema.sql` |
| `TCA_ADMINS` | Emails (separados por coma) que serán admin |
| `GOOGLE_CLIENT_ID` | OAuth de Google |
| `GOOGLE_CLIENT_SECRET` | OAuth de Google |
| `GOOGLE_REDIRECT_URL` | Origen de la app + `/api/v1/auth/google/callback` |
| `JWT_SECRET` | Secreto para firmar el JWT |

> `GOOGLE_REDIRECT_URL` apunta al **origen de la app**: en desarrollo `http://localhost:3000/...` (Nuxt proxya `/api`), en producción el host real del binario.

> `.env` contiene secretos reales. No se versiona.

### 2. Desarrollo

```bash
make dev   # genera sqlc+swagger, inicia backend (:8080) y frontend (:3000)
```

En desarrollo Nuxt (`:3000`) proxya `/api/*` hacia el backend (`:8080`) mediante `vite.server.proxy`; el navegador solo habla con un único origen, así que `apiBase` es siempre `''`.

También puedes lanzar cada parte por separado:

```bash
make dev-backend     # go run main.go  (:8080)
make dev-frontend    # cd webui && npm run dev  (:3000)
```

### 3. Build (ejecutable unificado)

```bash
make build   # genera el frontend y lo embebe en el binario
./build/trading-card-album
```

El binario sirve la API **y** el frontend embebido (con fallback SPA para rutas como `/album/:id`). El frontend se compila con `nuxt generate` y se embebe vía `//go:embed webui/.output/public/*`; por eso los assets se emiten en `/assets` (no `/_nuxt`), ya que `go:embed` excluye directorios que empiezan con `_`.

### Otros comandos

```bash
make test            # go test ./...
make gen-sql         # sqlc generate
make gen-swagger     # swag init
make gen             # gen-sql + gen-swagger
make db-migrate      # dbmate migrate (CLI)
```

---

## API (resumen)

Todo bajo `/api/v1`. La documentación completa está en `/swagger` (el frontend la consume vía orval).

- **Auth** `/auth`: `GET /` (inicia OAuth), `GET /me`, `GET /google/callback`
- **Album** `/album`: `GET /` (tus álbumes), `GET /:id` (solo participantes; devuelve solo las tarjetas recolectadas), `POST /` (admin), `GET /:id/join_qr` (QR de invitación), `GET /:id/card`, `GET /:id/assigned_card`, `POST /new_card`, `GET /:id/share_assigned_card?qr=`
- **Album participant** `/album_participant`: `POST /` (unirse con código)
- **Contact** `/contact`: `GET /`
- **User** `/user`: `PUT /:id`
- **Admin** `/admin` (requiere rol admin): `GET /overview`, `GET|POST /albums`, `PUT|DELETE /albums/:id`, `GET /albums/:id/cards`, `GET /users?email=`, `GET /users/:id`, `PUT /users/:id/role`, `GET|POST /cards`, `PUT|DELETE /cards/:id`

---

## Flujo de autenticación

1. `GET /api/v1/auth` → redirige a Google OAuth con `state` en cookie.
2. `GET /api/v1/auth/google/callback` → valida el state, intercambia el código, obtiene email y foto, crea/actualiza el usuario (admin si está en `TCA_ADMINS`) y firma un **JWT** en cookie `jwt`.
3. El middleware `CheckJwtCoockie` valida el JWT y deja el usuario en `c.Locals("session")`.
4. `CheckIsAdmin` protege las rutas de administración.

---

## Base de datos

- Motor: **SQLite**. Archivo en `database/trading-card-album.sqlite3` (`DATABASE_URL`).
- Migraciones con **dbmate**: se aplican automáticamente al arrancar (`database.RunMigrations()` usa las migraciones embebidas) o con `make db-migrate`.
- `database/schema.sql` es el dump que sqlc usa como esquema.

**Flujo de cambio de esquema**: crear migración en `database/migrations/` → actualizar `database/schema.sql` → `make gen-sql`.

### Convenciones de datos

- IDs: **nanoid** (texto).
- Timestamps: **unix** (`INTEGER`).
- Flags: `INTEGER` (0/1) (`is_admin`, `is_drawn`).
- El `card_pool` se repone automáticamente cuando se agota.
- SQLite no activa FKs por defecto: el borrado de álbumes en admin limpia explícitamente las filas relacionadas.

---

## Tests

Cobertura de unit tests en la capa `application` de cada contexto usando el mock `database/queriermock`:

```bash
make test   # go test ./...
```

## Frontend — notas

- Servicios/modelos generados por **orval** (`webui/app/services`, `webui/app/models`) desde `docs/swagger.yaml`; `npm run dev` y `npx orval` los regeneran.
- `customFetch` (`webui/app/services/CustomFetch.ts`) añade `credentials: 'include'` y usa `apiBase` (siempre `''`).
- App **mobile-only**: shell de teléfono (max 480px), barra de pestañas inferior, bottom sheets (`AppSheet`), safe areas.
- Páginas: `login`, `index`, `album/[id]`, `contactos`, `profile`, `admin/*`. Componentes: `Card`, `UserAvatar`, `TabBar`, `AppSheet`, `MyQrModal`, `ScanQrModal`, `JoinAlbumSheet`, `InviteSheet`, `AdminShell`, etc.