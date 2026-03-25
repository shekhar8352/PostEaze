# Repositories

This folder contains the repository layer implementation that provides data access patterns and abstracts database operations for the PostEaze application. The repository pattern separates business logic from data access logic, making the codebase more maintainable and testable.

## Architecture

The repository layer follows a functional approach where each repository file contains functions that operate on specific entities. The repositories use the custom database utility layer to execute raw SQL queries through entity-defined query patterns.

### Key Components

- **Repository Functions**: Stateless functions that handle specific data operations
- **Entity Integration**: Direct integration with entity query patterns and constants
- **Transaction Support**: Support for both transactional and non-transactional operations
- **Context Handling**: Proper context propagation for request lifecycle management

## Repository Pattern Implementation

### Data Access Flow

```
Repository Function → Entity Query Pattern → Database Utility → PostgreSQL
```

1. **Repository Function**: Handles business-specific data operations
2. **Entity Query Pattern**: Provides SQL queries and parameter binding
3. **Database Utility**: Manages connection pooling and query execution
4. **PostgreSQL**: Underlying database storage

### Entity Integration

Each repository function works with corresponding entity structs that implement the `RawEntity` interface:

```go
// Entity defines query patterns and parameter binding
func (o *User) GetQuery(code int) string { ... }
func (o *User) GetQueryValues(code int) []any { ... }
func (o *User) BindRawRow(code int, row Scanner) error { ... }
```

## Repository Files

### user.go
User and refresh-token access aligned with Firebase auth (see `entities.User` operations): create user with Firebase ID, lookup by email/Firebase ID, refresh token insert/revoke, etc.

### team.go
Contains team-related data access functions:

- **SaveTeam**: Creates new team records with owner assignment
- **AddListOfUsersToTeam**: Bulk addition of users to teams with role assignment

### channel.go
Channel CRUD, tokens, and Instagram-specific helpers (see source for full list).

### post.go
Post storage and listing (`posts` table), filters by channel, provider post IDs.

### channel_access.go / post_access.go
Authorization helpers (`UserCanAccessChannel`, post/channel membership checks) used by middleware and handlers.

### analytics_queries.go / instagram_analytics.go
Instagram post and story analytics reads/writes for dashboards and reports.

### teamMembers.go
Team membership operations alongside `team.go`.

## Usage Patterns

### Basic Repository Function

```go
func ExampleCreate(ctx context.Context, tx database.Database) error {
    data := entities.User{}
    // Populate FirebaseID, Name, Email, Platforms, then:
    return tx.QueryRaw(ctx, &data, entities.CreateUserWithFirebase)
}
```

### Transaction Usage

```go
// With transaction
tx, err := database.GetTx(ctx, nil)
if err != nil {
    return err
}
defer database.RollbackTx(tx)

user, err := CreateUser(ctx, tx, userData)
if err != nil {
    return err
}

return database.CommitTx(tx)
```

### Non-transactional Usage

```go
// Direct database access
user, err := GetUserByEmail(ctx, "user@example.com")
if err != nil {
    return err
}
```

## Data Access Patterns

### Query Execution Pattern

1. **Entity Preparation**: Create entity struct with required data
2. **Query Code**: Use entity-defined constants for query identification
3. **Database Execution**: Call `QueryRaw` or `QueryMultiRaw` through database interface
4. **Result Binding**: Entity automatically binds results through `BindRawRow`

### Error Handling

- **Database Errors**: Propagated from database layer
- **No Records**: `database.ErrNoRecords` for empty result sets
- **Context Cancellation**: Proper context handling for request timeouts

## Dependencies

### Internal Dependencies
- **entities**: Entity definitions and query patterns
- **models/v1**: API model definitions for data transformation
- **utils/database**: Database connection and query execution utilities

### External Dependencies
- **context**: Request context management
- **time**: Time handling for token expiration

## Best Practices

### Function Design
- **Single Responsibility**: Each function handles one specific data operation
- **Context First**: Always accept context as the first parameter
- **Transaction Support**: Accept database interface for transaction flexibility
- **Error Propagation**: Return errors for proper error handling upstream

### Data Transformation
- **Model Conversion**: Transform between API models and entity structs
- **Field Mapping**: Explicit field mapping for data integrity
- **Type Safety**: Proper type conversion and validation

### Performance Considerations
- **Connection Reuse**: Leverage database connection pooling
- **Query Optimization**: Use entity-defined optimized queries
- **Batch Operations**: Support bulk operations where applicable

## Related Documentation

- [../README.md](../README.md) - Entity layer overview
- [../../utils/database/README.md](../../utils/database/README.md) - Database utilities
- [../../models/README.md](../../models/README.md) - Data models
- [../../business/README.md](../../business/README.md) - Business logic layer