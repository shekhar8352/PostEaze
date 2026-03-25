# Middleware

Gin `HandlerFunc` chains for cross-cutting behavior: CORS, request logging, JWT authentication, role checks, and Instagram channel analytics access.

## Order in the app

In `api/router.go`:

1. **CORS** (`gin-contrib/cors`) — Allowed origins, methods, headers; credentials as configured.
2. **`GinLoggingMiddleware`** — Structured request/response logging with correlation IDs (global).

Per-route middleware:

- **`AuthMiddleware`** — Requires `Authorization: Bearer <access_token>`; parses JWT with `utils.ParseToken(..., false)`; sets `user_id` and `role` on the Gin context.
- **`RequireRole(allowed...)`** — After `AuthMiddleware`; returns 403 if `role` is not in the allow list.
- **`RequireInstagramChannelAnalyticsAccess`** — After `AuthMiddleware` on `/channels/:channelId/analytics/*`. Parses `channelId`, checks `repositories.UserCanAccessChannel`, ensures Instagram provider, sets channel on context for handlers.

## Files

| File | Role |
|------|------|
| `log_middleware.go` | Request logging |
| `auth_middleware.go` | JWT + `RequireRole` |
| `channel_analytics_middleware.go` | Channel ownership / team access for analytics routes |

## Authentication behavior

- Missing or non-Bearer header → 401 JSON.
- Invalid/expired token → 401.
- `RequireRole` without matching role → 403.

## Related documentation

- [API router](../api/router.go)
- [Utils JWT](../utils/jwt.go) (token parsing)
- [Channel access queries](../entities/repositories/channel_access.go)
