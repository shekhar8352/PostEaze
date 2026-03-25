# Database migrations

Numbered SQL migrations for PostgreSQL. Each change has a matching `.up.sql` and `.down.sql` file.

## Files (current)

| Migration | Purpose |
|-----------|---------|
| `001_initial_schema.up/down.sql` | Core schema: users (Firebase + `platforms`), `refresh_tokens`, teams, team members, channels, posts, Instagram analytics tables, etc. Enables `pgcrypto` and `pg_trgm`. |
| `002_analytics_extra_columns.up/down.sql` | Extra columns for analytics / reporting (see file for details). |
| `003_user_firebase_auth_columns.up/down.sql` | User columns aligned with Firebase auth (see file for details). |

Apply in numeric order on empty or known-state databases. For a greenfield dev DB, run `001` then `002` then `003` (or as required by your branch).

## Naming

```
{NNN}_{description}.up.sql
{NNN}_{description}.down.sql
```

## Applying manually

```bash
psql "$DATABASE_URL" -f backend/migrations/001_initial_schema.up.sql
# then 002, 003 as needed
```

Use your real connection string or Docker `psql` invocation. Down migrations reverse the corresponding up migration.

## Schema source of truth

The canonical picture of tables and indexes is the latest `001` + subsequent migrations. Notable areas:

- **Auth** — Users identified by `firebase_id`; refresh tokens in `refresh_tokens`.
- **Teams** — `teams`, `team_members` with roles and status.
- **Social** — Channels, encrypted tokens, posts, Instagram post/story analytics (see `001` and analytics migrations).

## Automated runner

There is no built-in migration runner in the app binary yet; apply SQL explicitly or integrate a tool (e.g. `golang-migrate`, `goose`) if you need version tracking in a `schema_migrations` table.

## Related documentation

- [Database utilities](../utils/database/README.md)
- [Entities](../entities/README.md)
- [Repository layer](../entities/repositories/README.md)
- [Docker / DB init](../../init-db/README.md) (repository root)
