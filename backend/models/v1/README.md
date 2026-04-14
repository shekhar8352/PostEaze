# Models v1

HTTP and shared structs for API version 1: JSON binding/validation tags for Gin, plus types reused in business logic.

## Layout

| File | Contents |
|------|----------|
| `user.go` | `FirebaseAuthParams`, user DTOs, `UpdateUserParams` |
| `tokens.go` | Refresh/access token request payloads |
| `log.go` | Log query responses |
| `team.go` | Team and membership models |
| `channel.go` | Instagram channel creation, page details |
| `analytics.go` | Analytics API request/response types |
| `scheduled_post.go` | Scheduled post create/list/get/cancel DTOs and query types |
| `media_asset.go` | Media upload response, asset/version DTOs, list query, publish payload |

## Authentication

Clients call `POST /api/v1/auth/authenticate` with a **Firebase ID token** and **platform** (and optional identifiers). The handler binds to `FirebaseAuthParams` (see `models/v1` source for exact fields and validation tags).

Refresh and logout use `RefreshTokenParams` / body fields expected by `auth.go` handlers.

## Validation

Handlers use `ShouldBindJSON` with `binding:"..."` tags on structs. Prefer adding new fields in `models/v1` rather than anonymous maps.

## Related documentation

- [API v1](../../api/v1/README.md)
- [Business v1](../../business/v1/README.md)
- [Models package](../README.md)
