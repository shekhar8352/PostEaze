# API v1 handlers

All routes below are prefixed with `/api/v1` unless stated. Request/response bodies are JSON unless noted.

## Files

| File | Responsibility |
|------|----------------|
| `auth.go` | Firebase authentication, refresh, logout, current user |
| `user.go` | Get/update user by ID |
| `team.go` | Team CRUD-style operations |
| `channel.go` | Instagram channels, page details, webhook subscription |
| `meta.go` | Meta OAuth callback |
| `log.go` | Read application logs by date or ID |
| `posts.go` | List posts for authenticated user |
| `analytics.go` | Instagram analytics under `/channels/:channelId/analytics` |
| `dev.go` | Development-only test JWT |
| `cron.go` | Dev-oriented triggers for background sync jobs |

## Authentication (`/auth`)

Base: `/api/v1/auth`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/authenticate` | No | Body: `FirebaseAuthParams` — Firebase ID token, platform, optional local ID. Creates or updates user, returns JWTs. |
| POST | `/refresh` | No | Body: `refresh_token`. Returns new tokens. |
| POST | `/logout` | No | Body: refresh token; revokes session server-side. |
| GET | `/me` | Bearer JWT | Returns current user profile. |

## Users (`/user`)

| Method | Path | Auth (router) |
|--------|------|----------------|
| GET | `/user/:user_id` | Not enforced by `AuthMiddleware` |
| PUT | `/user/:user_id` | Not enforced by `AuthMiddleware` |

Swagger annotations may reference `/users`; the actual paths are `/user/:user_id` per `constants.UserRoute`.

## Teams (`/team`)

| Method | Path | Auth (router) |
|--------|------|----------------|
| POST | `/team/create` | Not enforced |
| GET | `/team/all` | Not enforced |
| GET | `/team/:id` | Not enforced |
| GET | `/team/owner/:id` | Not enforced |
| PUT | `/team/update` | Not enforced |
| PUT | `/team/update-status` | Not enforced |

## Meta (`/meta`)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/meta/callback` | Exchange OAuth code for pages/tokens (Meta flow). |

## Channels & Instagram (`/channels`)

| Method | Path | Auth |
|--------|------|------|
| GET | `/channels` | JWT |
| GET | `/channels/details` | JWT — query `channel_id` |
| POST | `/channels/instagram/create` | JWT |
| POST | `/channels/instagram/subscribe-webhooks` | JWT |

## Webhooks

Mounted in router under `/api/v1/webhooks` (see [webhooks README](../webhooks/README.md)):

- `GET /webhooks/instagram` — Meta verification
- `POST /webhooks/instagram` — Events (signature verified)

## Posts

| Method | Path | Auth |
|--------|------|------|
| GET | `/posts` | JWT |

## Analytics (`/channels/:channelId/analytics`)

All routes require JWT and `RequireInstagramChannelAnalyticsAccess` (user must own channel or have team access; Instagram provider).

Examples: `GET .../profile`, `.../posts`, `.../overview`, `.../top-posts`, `.../posts-overview`, `.../posts/:postId`, `.../dashboard`, `.../comparison`, `.../stories`.

## Logs (`/log`)

| Method | Path |
|--------|------|
| GET | `/log/byDate/:date` |
| GET | `/log/byId/:log_id` |

## Development (`/dev`)

| Method | Path | Notes |
|--------|------|--------|
| POST | `/dev/generate-token` | Only when `ENV` is `development`, `dev`, or empty (see handler). |

## Cron triggers (`/cron`)

Manual enqueue of background jobs (handlers restrict non-dev `ENV` — see `cron.go`):

- `POST /cron/trigger-instagram-sync`
- `POST /cron/trigger-instagram-posts`
- `POST /cron/trigger-instagram-analytics`

## Conventions

- **Success**: `utils.SendSuccess` with message + data payload.
- **Errors**: `utils.SendError` with HTTP status and message string.
- **Bearer**: `Authorization: Bearer <access_token>` for protected routes.

## Related documentation

- [API layer](../README.md)
- [Business v1](../../business/v1/README.md)
- [Models v1](../../models/v1/README.md)
- [Middleware](../../middleware/README.md)
