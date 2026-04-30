# Business logic layer

Orchestrates domain rules between HTTP handlers and repositories. Stateless functions take `context.Context` first; database work uses `utils/database` transactions where needed.

## Responsibilities

- **Firebase authentication** — Validate Firebase ID tokens, upsert users by `firebase_id`, issue JWT access/refresh tokens, refresh and logout.
- **Users & teams** — User lookup/update; team create, list, update, status.
- **Channels** — Instagram channel creation and related flows (via services/repositories).
- **Media workspace** — Upload to blob storage, CRUD on `media_assets` / `media_versions`, set current version, publish into scheduled posts.
- **Studio** — Team studios, phases, pieces, moves, asset/scheduled-post links, comments, activities (`studio.go` in `v1/`).
- **Logs** — Read structured log files for admin/debug (`ReadLogsByDate`, `ReadLogsByLogID`).

## Patterns

- **Functions** — e.g. `AuthenticateWithFirebase`, `GetUserById`, `CreateTeam` (see `v1/`).
- **Context** — Passed through for cancellation and logging.
- **Transactions** — `database.GetTx` / `CommitTx` / `RollbackTx` when multiple writes must stay consistent.
- **Errors** — Returned to API layer for HTTP mapping; security-sensitive paths return generic messages.

## Dependencies

- `models/v1` — Request/response types
- `entities/repositories` — SQL access
- `utils` — JWT, Firebase client, logging, hashing where applicable
- `utils/database` — Pool and transactions

## Layout

- **`v1/`** — Current implementation (`auth.go`, `user.go`, `team.go`, `channel.go`, `scheduled_post.go`, `media_asset.go`, `studio.go`, `log.go`, …)

## Related documentation

- [API v1](../api/v1/README.md)
- [Business v1](./v1/README.md)
- [Entities](../entities/README.md)
- [Studio pipeline](../docs/studio-pipeline.md)
- [Repositories](../entities/repositories/README.md)
