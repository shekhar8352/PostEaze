# Business logic v1

Versioned domain services used by `api/v1` handlers. Each file groups related operations.

## Modules

| File | Role |
|------|------|
| `auth.go` | `AuthenticateWithFirebase` — verify Firebase token, create/update user, platforms; `RefreshToken`, `Logout`; helpers for new/existing users and token issuance |
| `user.go` | `GetUserById`, `UpdateUser` |
| `team.go` | Team create, list, get by ID/owner, update, status |
| `channel.go` | Instagram channel orchestration (delegates to services/repositories) |
| `log.go` | `ReadLogsByLogID`, `ReadLogsByDate` — scan JSON log files under `logs/` |
| `scheduled_post.go` | Create, get, list (calendar range), cancel scheduled posts for the authenticated owner |
| `media_asset.go` | Blob upload, media asset + version lifecycle, set current version, publish (scheduled post from asset) |

## Authentication flow (Firebase)

1. Client sends Firebase ID token + platform to `AuthenticateWithFirebase`.
2. `utils` Firebase Admin verifies the token.
3. If no row matches `firebase_id`, a user row is created; else platforms may be updated.
4. Access and refresh JWTs are returned for subsequent API calls.

Refresh/logout use refresh tokens stored in `refresh_tokens` via repository helpers.

## Transactions

Multi-step writes (e.g. user + token rows) use `database.GetTx` and commit/rollback patterns consistent with the rest of the codebase.

## Related documentation

- [Business overview](../README.md)
- [API v1](../../api/v1/README.md)
- [Models v1](../../models/v1/README.md)
- [Repositories](../../entities/repositories/README.md)
