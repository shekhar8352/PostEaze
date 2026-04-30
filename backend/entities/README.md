# Entities

The entities package contains the core data models and database interaction patterns for PostEaze. This package implements a custom query pattern where entities define their own SQL queries and parameter binding logic, providing a lightweight alternative to traditional ORM solutions.

## Architecture

The entities package follows a pattern where each entity:
- Defines its data structure with JSON tags for API serialization
- Implements the `RawEntity` interface for database operations
- Contains query constants and SQL query definitions
- Handles parameter binding and result scanning

### RawEntity Interface

All entities implement the `database.RawEntity` interface which provides:
- `GetQuery(code int)` - Returns SQL query for a given operation code
- `GetQueryValues(code int)` - Returns parameter values for the query
- `BindRawRow(code int, row Scanner)` - Binds database row results to entity fields
- `GetNextRaw()` - Returns a new instance for result scanning

## Key files

- **user.go** — User row + refresh-token helpers (Firebase identity, `platforms` array)
- **team.go** — Team and membership queries
- **channel.go** (and related) — Social channel + token storage
- **media_asset.go**, **media_version.go** — Media workspace rows and query codes
- **studio.go**, **phase.go**, **piece.go** — Studio pipeline (team board, columns, work items, joins)
- **repositories/** — Repository functions calling `QueryRaw` / transactions

## Entity structure

### User entity

Users are keyed by internal UUID and `firebase_id`; there is no password column.

```go
type User struct {
    ID           string         `json:"id"`
    FirebaseID   string         `json:"firebase_id"`
    Name         string         `json:"name"`
    Email        string         `json:"email"`
    Platforms    pq.StringArray `json:"platforms"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    RefreshToken string         `json:"refresh_token"` // used when binding token flows
    ExpiresAt    time.Time      `json:"expire_at"`
}
```

**Operation codes (see `user.go`):** `CreateUserWithFirebase`, `InsertRefreshToken`, `GetUserByEmail`, `GetUserByToken`, `GetUserByID`, `GetUserByFirebaseID`, `UpdateUserPlatforms`, `UpdateUser`, `RevokeTokens`

### Team Entity

The Team entity manages team creation and member relationships:

```go
type Team struct {
    ID        string       `json:"id"`
    Name      string       `json:"name"`
    OwnerID   string       `json:"owner_id"`
    Members   []TeamMember `json:"members"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
}
```

**Supported Operations:**
- `CreateTeam` - Create new team record
- `AddUsersToTeam` - Add multiple users to team with roles

## Query Pattern Implementation

Each entity uses integer constants to identify different database operations:

Operation constants are defined per entity (e.g. `CreateUserWithFirebase = iota` in `user.go`). `GetQuery` / `GetQueryValues` / `BindRawRow` implement `database.RawEntity`.

## Usage examples

### Creating a user (Firebase)

```go
user := entities.User{
    FirebaseID: "firebase-uid",
    Name:       "Jane",
    Email:      "jane@example.com",
    Platforms:  pq.StringArray{"web"},
}
err := db.QueryRaw(ctx, &user, entities.CreateUserWithFirebase)
```

### Retrieving a User

```go
user := entities.User{Email: "john@example.com"}
err := db.QueryRaw(ctx, &user, entities.GetUserByEmail)
// user struct is populated with database values
```

## Database Integration

Entities integrate with the database layer through:
- **database.RawEntity interface** - Standardized query execution pattern
- **Parameter binding** - Type-safe parameter passing to SQL queries
- **Result scanning** - Automatic mapping of database rows to struct fields
- **Transaction support** - Compatible with database transaction handling

## Related Documentation

- [Repositories](./repositories/README.md) - Repository pattern implementations
- [Models v1](../models/v1/README.md) - API request/response models
- [Database Utils](../utils/database/README.md) - Database connection and utilities
- [Migrations](../migrations/README.md) - Database schema management