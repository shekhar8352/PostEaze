# Database migrations

Numbered SQL migrations for PostgreSQL. Each change has a matching `.up.sql` and `.down.sql` file.

## Files (current)

| Migration | Purpose |
|-----------|---------|
| `001_initial_schema.up/down.sql` | Core schema: users (Firebase + `platforms`), `refresh_tokens`, teams, team members, channels, posts, Instagram analytics tables, etc. Enables `pgcrypto` and `pg_trgm`. |
| `002_analytics_extra_columns.up/down.sql` | Extra columns for analytics / reporting (see file for details). |
| `003_user_firebase_auth_columns.up/down.sql` | User columns aligned with Firebase auth (see file for details). |
| `004_posts_multi_provider.up/down.sql` | Post storage updates for multiple providers (see file). |
| `005_instagram_audience_snapshots.up/down.sql` | Instagram audience snapshot storage (see file). |
| `006_scheduled_posts.up/down.sql` | Scheduled posts tables for the calendar API (see file). |
| `007_media_workspace.up/down.sql` | Media assets and versions (blob-backed workspace). |

Apply in numeric order on empty or known-state databases. For a greenfield dev DB, run `001` through the latest migration (or as required by your branch).

## golang-migrate: “Dirty database version N”

If a migration fails mid-run, `schema_migrations` is left with `dirty = true` and the same version, and every later `migrate up` errors until you fix it.

**After a failed `007` (e.g. bad FK types),** reset the version to the last known-good migration, clear dirty, then apply again:

```bash
cd backend
export DATABASE_URL='postgres://USER:PASS@HOST:5432/DBNAME?sslmode=disable'

make migrate-dirty-redo-seven
```

That runs `migrate force 6` then `migrate up`, so `007` runs again with the fixed SQL.

**Manual equivalent:**

```bash
migrate -path ./migrations -database "$DATABASE_URL" force 6
migrate -path ./migrations -database "$DATABASE_URL" up
```

**If migration N actually completed** (you verified objects exist) but the tool still reports dirty, you can clear only the flag with `migrate force N` (same N), then avoid re-running destructive SQL.

**SQL-only** (if you cannot use the CLI): golang-migrate keeps a single row in `schema_migrations` (`version`, `dirty`). Setting `version` to the last good migration and `dirty` to `false` matches what `migrate force` does—prefer the CLI when possible.

## Naming

```
{NNN}_{description}.up.sql
{NNN}_{description}.down.sql
```

## Applying manually

```bash
psql "$DATABASE_URL" -f backend/migrations/001_initial_schema.up.sql
# then 002 … 007 in order as needed
```

Use your real connection string or Docker `psql` invocation. Down migrations reverse the corresponding up migration.

## Schema source of truth

The canonical picture of tables and indexes is the latest `001` + subsequent migrations. Notable areas:

- **Auth** — Users identified by `firebase_id`; refresh tokens in `refresh_tokens`.
- **Teams** — `teams`, `team_members` with roles and status.
- **Social** — Channels, encrypted tokens, posts, Instagram post/story analytics (see `001` and analytics migrations).
- **Media workspace** — `media_assets` and `media_versions` (blob-backed assets with versioning); see `007_media_workspace`.

## golang-migrate (recommended)

Install the CLI (Postgres driver):

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

From `backend/` with `DATABASE_URL` set:

| Make target | Effect |
|-------------|--------|
| `make migrate-up` | Apply all pending migrations (`migrate up`). |
| `make migrate-force VERSION=N` | Set `schema_migrations` to version `N` and clear `dirty` without running SQL (recovery only). |
| `make migrate-dirty-redo-seven` | `migrate force 6` then `migrate up` — use after a failed `007` when you need to re-run it (see above). |

Alternatively run `migrate -path ./migrations -database "$DATABASE_URL" up` directly.

You can still apply `.up.sql` files with `psql` when you do not use golang-migrate; then there is no `schema_migrations` row unless you add tooling yourself.

## Related documentation

- [Database utilities](../utils/database/README.md)
- [Entities](../entities/README.md)
- [Repository layer](../entities/repositories/README.md)
- [Docker / DB init](../../init-db/README.md) (repository root)
