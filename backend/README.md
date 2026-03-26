# PostEaze Backend

Go REST API for PostEaze: Firebase-based authentication, teams, Instagram channels, webhooks, background jobs (Asynq/Redis), and analytics. HTTP layer uses Gin; data access uses PostgreSQL with a raw-query entity pattern (`lib/pq`).

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
│  Business — business/v1 (auth, user, team, channel, log)    │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Data — entities, repositories, models/v1                   │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Infrastructure — configs, Redis, encryption, Firebase      │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│  Providers — Meta, Instagram                                  │
│  Tasks — Asynq client (API) + worker (cmd/worker)             │
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
| `business/v1/` | Domain logic (Firebase auth, users, teams, channels, logs) |
| `entities/` | `RawEntity` SQL patterns; `repositories/` data access |
| `models/v1/` | Request/response and shared structs |
| `migrations/` | Numbered `*.up.sql` / `*.down.sql` |
| `middleware/` | Logging, JWT auth, roles, Instagram analytics access |
| `provider/` | Meta Graph API, Instagram OAuth |
| `services/` | Email, Redis, Meta, Instagram orchestration |
| `tasks/` | Asynq task definitions, handlers, scheduler |
| `cmd/worker/` | Standalone worker process |
| `utils/` | Config, DB, env, flags, HTTP, JWT, Firebase, Redis, encryption |
| `resources/configs/` | Per-environment YAML (`dev/`, `cug/`, `prod/`) |

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
| Channels | `GET /channels`, `GET /channels/details`, `POST /channels/instagram/create`, `POST /channels/instagram/subscribe-webhooks` | Most require JWT |
| Webhooks | `GET`, `POST /webhooks/instagram` | Meta verification + events |
| Posts | `GET /posts` | JWT |
| Analytics | `GET /channels/:channelId/analytics/...` | JWT + `RequireInstagramChannelAnalyticsAccess` (profile, posts, overview, dashboard, etc.) |
| Dev | `POST /dev/generate-token` | Test JWT helpers when `ENV=development` / `dev` |
| Cron (dev-oriented) | `POST /cron/trigger-instagram-sync`, `.../trigger-instagram-posts`, `.../trigger-instagram-analytics` | Guarded by `ENV` in handlers |

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
# Apply SQL in migrations/ in order (see migrations/README.md)
go run .
```

Flags: `-mode` (`dev` | `release`), `-port`, `-base-config-path` (default `resources/configs/dev`).

## Swagger

Regenerate after changing handler comments:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
cd backend && swag init
```

Swagger host comes from `API_HOST` (see `api/router.go`). Open `http://<API_HOST>/api/swagger/index.html` (scheme/path must match your deployment).

## Environment variables (common)

| Variable | Purpose |
|----------|---------|
| `ENV` | `development` / `dev` vs production behavior (e.g. dev-only routes) |
| `API_HOST` | Swagger `Host` field |
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
