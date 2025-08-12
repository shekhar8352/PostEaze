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

## Key Files

- **user.go**: User entity with authentication and profile management queries
- **team.go**: Team entity with team creation and member management queries
- **repositories/**: Repository pattern implementations for data access

## Entity Structure

### User Entity

The User entity handles user authentication, profile management, and refresh token operations:

```go
type User struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    Password     string    `json:"-"`          // Hidden from JSON serialization
    UserType     string    `json:"user_type"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expire_at"`
}
```

**Supported Operations:**
- `CreateUser` - Insert new user record
- `InsertRefreshToken` - Store refresh token for authentication
- `GetUserByEmail` - Retrieve user by email address
- `GetUserByToken` - Validate and retrieve user by refresh token
- `GetUserByID` - Retrieve user by ID
- `RevokeTokens` - Invalidate all refresh tokens for a user

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

```go
const (
    CreateUser = iota
    InsertRefreshToken
    GetUserByEmail
    // ... more operations
)
```

The `GetQuery()` method returns the appropriate SQL query based on the operation code:

```go
func (o *User) GetQuery(code int) string {
    switch code {
    case CreateUser:
        return `INSERT INTO users (name, email, password, user_type) 
                VALUES ($1, $2, $3, $4) 
                RETURNING id, created_at, updated_at;`
    // ... more cases
    }
}
```

## Usage Examples

### Creating a User

```go
user := entities.User{
    Name:     "John Doe",
    Email:    "john@example.com",
    Password: "hashedpassword",
    UserType: "individual",
}

err := db.QueryRaw(ctx, &user, entities.CreateUser)
// user.ID, user.CreatedAt, user.UpdatedAt are populated from RETURNING clause
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