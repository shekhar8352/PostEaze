# PostEaze Backend

Go REST API for PostEaze: Firebase-based authentication, teams, Instagram and Facebook channels, Meta OAuth and analytics sync, webhooks, scheduled posts, **Studio** (per-team phases and pieces with media and scheduled-post links), media workspace (versioned assets and publish-to-scheduled-post), background jobs (Asynq/Redis), and channel analytics. HTTP layer uses Gin; data access uses PostgreSQL with a raw-query entity pattern (`lib/pq`).

## Architecture overview

```
┌─────────────────────────────────────────────────────────────┐
│                     API (Gin) — /api/*                       │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Middleware — CORS, request logging, JWT, channel analytics   │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Business — business/v1 (auth, user, team, channel, log,       │
│            scheduled posts, media assets, analytics, Meta,      │
│            studio / phases / pieces, posts)                      │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Data — entities, repositories, models/v1                   │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Infrastructure — configs, Redis, encryption, Firebase, blob    │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Helpers — services/studio (fractional index, phase templates)  │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Providers — Meta / Instagram Graph API, insights               │
│  Tasks — Asynq client (API) + worker (cmd/worker)               │
└───────────────────────────────────────────────────────────────┘
```

## Request flow

1. **HTTP** → `api/router.go` (`api.Init`)
2. **Middleware** → CORS, `GinLoggingMiddleware`, route-specific `AuthMiddleware` / `RequireInstagramChannelAnalyticsAccess`
3. **Handlers** → `api/v1/*.go`
4. **Business** → `business/v1/`
5. **Repositories** → `entities/repositories/` → PostgreSQL

## Key directories

| Path | Role |
|------|------|
| `main.go` | Startup: env, configs, DB, Redis, encryption, Firebase, Asynq client, router, HTTP client |
| `api/` | Router, Swagger, `v1` handlers, Instagram webhooks |
| `business/v1/` | Domain logic (Firebase auth, users, teams, channels, logs, scheduled posts, media, analytics, Meta, **studio/pieces**, posts) |
| `entities/` | `RawEntity` SQL patterns; `repositories/` data access |
| `models/v1/` | Request/response and shared structs |
| `migrations/` | Numbered `*.up.sql` / `*.down.sql` |
| `middleware/` | Logging, JWT auth, roles, Instagram analytics access |
| `provider/` | Meta Graph API, Instagram OAuth, insights |
| `services/` | Email, Redis, Meta, Instagram orchestration |
| `services/studio/` | Fractional indexing + default phase templates for the Studio pipeline |
| `tasks/` | Asynq task definitions, handlers, scheduler |
| `cmd/worker/` | Standalone worker process |
| `utils/` | Config, DB, env, flags, HTTP, JWT, Firebase, Redis, encryption |
| `resources/configs/` | Per-environment YAML (`dev/`, `cug/`, `prod/`) |
| `docs/` | Generated Swagger + human guides (`firebase-authentication.md`, `studio-pipeline.md`, …) |

## Service initialization (`main.go`)

Order: `initEnv` → `initConfigs` → `initDatabase` → `initRedis` → `initEncryption` → `initFirebase` → `initAsynq` → `initRouter` → `initHttp`.

- **Configs**: `dev` uses `resources/configs/dev` (override with `-base-config-path`); `release` uses AWS App Config.
- **Database**: PostgreSQL URL from config + `utils/env` substitution.
- **Asynq**: Client only in the API process; run `cmd/worker` for consumers.

## API surface (summary)

Base path: `/api/v1` unless noted.

| Area | Methods | Notes |
|------|---------|--------|
| Health | `GET /api/health` (liveness), `GET /api/health/postgres`, `GET /api/health/redis`, `GET /api/health/ready` | No version prefix; `ready` returns 503 if PG or Redis unavailable |
| Swagger | `GET /api/swagger/*` | UI at `/api/swagger/index.html` |
| Auth | `POST /auth/authenticate`, `POST /auth/refresh`, `POST /auth/logout`, `GET /auth/me` | Firebase ID token + platform; JWT for API |
| Logs | `GET /log/byDate/:date`, `GET /log/byId/:log_id` | |
| User | `GET /user/:user_id`, `PUT /user/:user_id` | |
| Team | `POST /team/create`, `GET /team/all`, `GET /team/:id`, `GET /team/owner/:id`, `PUT /team/update`, `PUT /team/update-status` | |
| Meta | `POST /meta/callback` | OAuth callback |
| Meta analytics | `POST /meta/analytics/sync` | JWT — sync analytics from Meta for connected accounts |
| Channels | `GET /channels`, `GET /channels/details`, `POST /channels/instagram/create`, `POST /channels/instagram/subscribe-webhooks`, `POST /channels/facebook/create` | Most require JWT |
| Webhooks | `GET`, `POST /webhooks/instagram` | Meta verification + events |
| Posts | `GET /posts` | JWT |
| Scheduled posts | `GET /scheduled-posts`, `POST /scheduled-posts`, `GET /scheduled-posts/:id`, `DELETE /scheduled-posts/:id` | JWT — list supports calendar range query params (see handlers) |
| Media workspace | `POST /media/upload`; `GET|POST /media-assets`, `GET|PUT|DELETE /media-assets/:id`, `POST /media-assets/:id/versions`, `DELETE /media-assets/:id/versions/:vid`, `PUT /media-assets/:id/current-version`, `POST /media-assets/:id/publish` | JWT — upload and versioning; publish creates a scheduled post from the current version |
| Studio | `/studios/*`, `/phases/*`, `/pieces/*` (board, phases, pieces, move, assets, scheduled posts, comments, activities) | JWT — see [`docs/studio-pipeline.md`](docs/studio-pipeline.md) |
| Analytics | See below | JWT + `RequireInstagramChannelAnalyticsAccess` (Instagram channel) |
| Dev | `POST /dev/generate-token` | Test JWT helpers when `ENV=development` / `dev` |
| Cron (dev-oriented) | `POST /cron/trigger-instagram-sync`, `.../trigger-instagram-posts`, `.../trigger-instagram-analytics` | Guarded by `ENV` in handlers |

### Channel analytics (`/channels/:channelId/analytics`)

All require JWT and Instagram channel analytics access.

| Method | Path |
|--------|------|
| GET | `/profile` |
| GET | `/posts` |
| GET | `/overview` |
| GET | `/top-posts` |
| GET | `/posts-overview` |
| GET | `/posts/:postId` |
| GET | `/dashboard` |
| GET | `/comparison` |
| GET | `/stories` |
| GET | `/audience` |

See [`api/v1/README.md`](api/v1/README.md) and [`api/webhooks/README.md`](api/webhooks/README.md) for detail.

## Dependencies (high level)

- **Web**: `gin-gonic/gin`, `gin-contrib/cors`, `swaggo` (Swagger UI)
- **DB**: `lib/pq`
- **Auth**: `golang-jwt/jwt/v5`, Firebase Admin SDK (`firebase.google.com/go/v4`)
- **Cache / queue**: `redis/go-redis`, `hibiken/asynq`
- **Config**: `sinhashubham95/go-config-client`, `joho/godotenv`

## Local development

**Prerequisites:** Go 1.23+, PostgreSQL, Redis (for Asynq), `.env` aligned with `resources/configs/dev`.

```bash
cd backend
go mod download
export DATABASE_URL='postgres://USER:PASS@HOST:5432/DBNAME?sslmode=disable'
make migrate-up   # requires golang-migrate CLI; see migrations/README.md
go run .
```

Flags: `-mode` (`dev` | `release`), `-port`, `-base-config-path` (default `resources/configs/dev`).

## Swagger

OpenAPI specs live in [`docs/`](docs/) (`swagger.json`, `swagger.yaml`). Regenerate after changing handler Swagger comments (`// @Summary`, `// @Router`, etc.):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
cd backend && swag init --generalInfo api/router.go --output docs --parseDependency --parseInternal
```

`docs.SwaggerInfo` in `api/router.go` sets **Host** from `API_HOST` and **BasePath** to `/api/v1`. Open `http://<API_HOST>/api/swagger/index.html` (scheme and host must match your deployment).

## Environment variables (common)

| Variable | Purpose |
|----------|---------|
| `ENV` | `development` / `dev` vs production behavior (e.g. dev-only routes) |
| `API_HOST` | Swagger `Host` field (e.g. `localhost:8000`) |
| Firebase / `GOOGLE_APPLICATION_CREDENTIALS` | Firebase Admin (see `utils/firebase.go`) |
| `ENCRYPTION_KEY` | Base64 32-byte key for channel tokens |
| `INSTAGRAM_*`, `META_*` | Instagram/Meta app and webhooks |
| `REDIS_ADDR` | Asynq broker (worker + API) |
| `SMTP_*` | Email (if used) |

Generate encryption key: `openssl rand -base64 32`

## Security

- JWT on protected routes via `Authorization: Bearer <access_token>`
- Firebase ID tokens verified on `/auth/authenticate`
- Instagram webhooks: `X-Hub-Signature-256` verification

## Logging

Structured logs under `logs/` (see [`logs/README.md`](logs/README.md)); HTTP logging middleware attaches request correlation IDs.

## Related docs

- [`api/README.md`](api/README.md), [`api/v1/README.md`](api/v1/README.md)
- [`business/README.md`](business/README.md), [`tasks/README.md`](tasks/README.md)
- [`migrations/README.md`](migrations/README.md), [`middleware/README.md`](middleware/README.md)
- [`docs/firebase-authentication.md`](docs/firebase-authentication.md)
- [`docs/studio-pipeline.md`](docs/studio-pipeline.md), [`services/studio/README.md`](services/studio/README.md)
