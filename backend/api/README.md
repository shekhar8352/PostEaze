# API layer

HTTP interface for the PostEaze backend: Gin router, versioned REST handlers under `/api/v1`, Swagger, and Instagram webhook routes.

## Layout

```
api/
├── router.go       # Gin engine, CORS, middleware, route groups, Swagger
├── health.go       # Liveness + Postgres/Redis checks + readiness aggregate
├── v1/             # Version 1 handlers (auth, user, team, channel, log, meta, posts, analytics, dev, cron)
└── webhooks/       # Instagram webhook verify + POST handler
```

## Router (`router.go`)

1. **CORS** — Allowed origins include local Vite and configured dev hosts; methods `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`; `Authorization` allowed.
2. **`GinLoggingMiddleware`** — Request/response logging (global).
3. **Groups** — Health: `GET /api/health` (liveness), `GET /api/health/postgres`, `GET /api/health/redis`, `GET /api/health/ready` (503 if PG or Redis down). `/api/v1` registers auth, logs, user, team, meta, channels (+ webhooks path), dev, cron, posts, analytics.
4. **Swagger** — `GET /api/swagger/*` via `gin-swagger`; `docs.SwaggerInfo` uses `API_HOST` and base path `/api/v1`.

## Route registration (v1)

| Registration | Base under `/api/v1` |
|--------------|----------------------|
| `addV1UserAuthRoutes` | `/auth` — authenticate, refresh, logout, `/me` |
| `addV1LogRoutes` | `/log` |
| `addV1UserRoutes` | `/user` |
| `addV1TeamRoutes` | `/team` |
| `addV1MetaRoutes` | `/meta` |
| `addV1ChannelRoutes` | `/channels`, `/webhooks` |
| `addV1DevRoutes` | `/dev` |
| `addV1CronRoutes` | `/cron` |
| `addV1PostRoutes` | `/posts` |
| `addV1AnalyticsRoutes` | `/channels/:channelId/analytics` |

Protected routes use `middleware.AuthMiddleware()`; analytics also uses `middleware.RequireInstagramChannelAnalyticsAccess()`.

## Error handling

Handlers typically use `utils.SendError` / `utils.SendSuccess` with appropriate HTTP status codes. Validation errors return 400; auth failures 401; channel analytics access 403/404 as applicable.

## Dependencies

- **Gin** — Routing and middleware
- **`business/v1`** — Business logic
- **`models/v1`** — Request/response shapes
- **`middleware`** — Auth, logging, analytics access
- **`constants`** — Path segments

## Related documentation

- [API v1 handlers](./v1/README.md)
- [Webhooks](./webhooks/README.md)
- [Middleware](../middleware/README.md)
- [Business layer](../business/README.md)
