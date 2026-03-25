# Backend utilities

Shared helpers: configuration, database access, env substitution, CLI flags, HTTP client, Redis, encryption, Firebase Admin, JWT issuance/parsing, logging, and Gin response helpers.

## Sub-packages

| Package | Doc |
|---------|-----|
| `configs/` | [README](./configs/README.md) — dev file configs vs AWS App Config |
| `database/` | [README](./database/README.md) — pool, transactions, `RawEntity` |
| `env/` | [README](./env/README.md) — `.env` load + `${VAR}` substitution |
| `flags/` | [README](./flags/README.md) — `-mode`, `-port`, `-base-config-path` |
| `http/` | [README](./http/README.md) — configured HTTP client |
| `redis/` | See `redis/redis.go` — client init |

## Root-level modules

### JWT (`jwt.go`) and legacy helpers (`auth_utils.go`)

- **Current API auth** — `GenerateAccessToken`, `GenerateRefreshToken`, `ParseToken`, claims with user ID and role (used with Firebase-backed users).
- **Legacy** — `GenerateJWT` / `ParseJWT` (`auth_utils.go`) may remain for older flows; prefer the v5 JWT path in `jwt.go` for new code.

### Logging (`logger.go`, `logapi_utils.go`)

Structured logger with daily files under `logs/`, correlation IDs, and helpers to read/search log files for the log API.

### Validation (`validation.go`)

Shared checks (e.g. log ID and date strings) for log endpoints.

### API responses (`general_utils.go`)

`SendSuccess`, `SendError`, and log-specific response helpers used across handlers.

### Firebase (`firebase.go`)

`InitializeFirebase`, `GetFirebaseService`, `ValidateToken` — Firebase Admin verification for `/auth/authenticate`.

### Redis (`redis/`)

`Init`, `GetClient` — used by Asynq and caching.

### Encryption (`encryption/`)

Used for channel secrets; initialized from `main.go`.

## Related documentation

- [Middleware](../middleware/README.md) (uses `ParseToken`)
- [Business v1](../business/v1/README.md)
