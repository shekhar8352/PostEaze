# API v1 handlers

All routes below are prefixed with `/api/v1` unless stated. Request/response bodies are JSON unless noted.

## Files

| File | Responsibility |
|------|----------------|
| `auth.go` | Firebase authentication, refresh, logout, current user |
| `user.go` | Get/update user by ID |
| `team.go` | Team CRUD-style operations |
| `channel.go` | Instagram and Facebook channels, page details, webhook subscription |
| `meta.go` | Meta OAuth callback |
| `meta_analytics.go` | Sync Meta analytics (authenticated) |
| `log.go` | Read application logs by date or ID |
| `posts.go` | List posts for authenticated user |
| `scheduled_post.go` | Scheduled posts CRUD + calendar listing |
| `media_asset.go` | Media upload, assets, versions, current version, publish to scheduled post |
| `studio.go` | Studios, phases, pieces, moves, asset/post links, comments, activities |
| `analytics.go` | Instagram analytics under `/channels/:channelId/analytics` |
| `dev.go` | Development-only test JWT |
| `cron.go` | Dev-oriented triggers for background sync jobs |

## Authentication (`/auth`)

Base: `/api/v1/auth`

| Method | Path | Auth | Description |
|--------|------|------|---------------|
| POST | `/authenticate` | No | Body: `FirebaseAuthParams` — Firebase ID token, platform, optional local ID. Creates or updates user, returns JWTs. |
| POST | `/refresh` | No | Body: `refresh_token`. Returns new tokens. |
| POST | `/logout` | No | Body: refresh token; revokes session server-side. |
| GET | `/me` | Bearer JWT | Returns current user profile. |

## Users (`/user`)

| Method | Path | Auth (router) |
|--------|------|----------------|
| GET | `/user/:user_id` | Not enforced by `AuthMiddleware` |
| PUT | `/user/:user_id` | Not enforced by `AuthMiddleware` |

Swagger may reference `/users`; the live paths are `/user/:user_id` per `constants.UserRoute`.

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

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/meta/callback` | No | Exchange OAuth code for pages/tokens (Meta flow). |
| POST | `/meta/analytics/sync` | JWT | Trigger sync of Meta analytics for the authenticated user’s connected assets. |

## Channels (`/channels`)

| Method | Path | Auth |
|--------|------|------|
| GET | `/channels` | JWT |
| GET | `/channels/details` | JWT — query `channel_id` |
| POST | `/channels/instagram/create` | JWT |
| POST | `/channels/instagram/subscribe-webhooks` | JWT |
| POST | `/channels/facebook/create` | JWT |

## Webhooks

Mounted in router under `/api/v1/webhooks` (see [webhooks README](../webhooks/README.md)):

- `GET /webhooks/instagram` — Meta verification
- `POST /webhooks/instagram` — Events (signature verified)

## Posts

| Method | Path | Auth |
|--------|------|------|
| GET | `/posts` | JWT |

## Scheduled posts (`/scheduled-posts`)

All routes require JWT (`AuthMiddleware`).

| Method | Path | Description |
|--------|------|-------------|
| GET | `/scheduled-posts` | List posts in a time range (calendar); query params from `ListScheduledPostsQuery` |
| POST | `/scheduled-posts` | Create a scheduled post |
| GET | `/scheduled-posts/:id` | Get one scheduled post |
| DELETE | `/scheduled-posts/:id` | Cancel / remove a scheduled post |

## Media workspace

All routes require JWT (`AuthMiddleware`).

### Upload (`/media`)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/media/upload` | Multipart file upload to blob storage; returns URL/metadata for use when creating assets or versions |

### Assets (`/media-assets`)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/media-assets` | List assets (query params from `ListMediaAssetsQuery`) |
| POST | `/media-assets` | Create asset + initial version (multipart: metadata + file) |
| GET | `/media-assets/:id` | Get one asset with all versions |
| PUT | `/media-assets/:id` | Update title/status |
| DELETE | `/media-assets/:id` | Delete asset and versions |
| POST | `/media-assets/:id/versions` | Add a version (multipart) |
| DELETE | `/media-assets/:id/versions/:vid` | Remove a version |
| PUT | `/media-assets/:id/current-version` | Set active version |
| POST | `/media-assets/:id/publish` | Create a scheduled post from the current version (body: channels, caption, schedule) |

## Studio (`/studios`, `/phases`, `/pieces`)

JWT required (`AuthMiddleware`). Full route map and semantics: **[`docs/studio-pipeline.md`](../../docs/studio-pipeline.md)** (`EnsureStudio`, board, phases CRUD/reorder, pieces CRUD/move/status, asset and scheduled-post links, comments, activities).

## Analytics (`/channels/:channelId/analytics`)

All routes require JWT and `RequireInstagramChannelAnalyticsAccess` (user must own channel or have team access; Instagram provider).

| Method | Path | Description |
|--------|------|-------------|
| GET | `/channels/:channelId/analytics/profile` | Time-series profile metrics |
| GET | `/channels/:channelId/analytics/posts` | Post-level analytics list |
| GET | `/channels/:channelId/analytics/overview` | Aggregated overview |
| GET | `/channels/:channelId/analytics/top-posts` | Top posts by engagement |
| GET | `/channels/:channelId/analytics/posts-overview` | Totals (likes, comments, etc.) |
| GET | `/channels/:channelId/analytics/posts/:postId` | Single post insights |
| GET | `/channels/:channelId/analytics/dashboard` | Dashboard bundle |
| GET | `/channels/:channelId/analytics/comparison` | Period-over-period comparison |
| GET | `/channels/:channelId/analytics/stories` | Story analytics |
| GET | `/channels/:channelId/analytics/audience` | Audience snapshots |

## Logs (`/log`)

| Method | Path |
|--------|------|
| GET | `/log/byDate/:date` |
| GET | `/log/byId/:log_id` |

## Development (`/dev`)

| Method | Path | Notes |
|--------|------|-------|
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
- [Studio pipeline](../../docs/studio-pipeline.md)
