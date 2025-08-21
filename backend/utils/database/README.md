# Database Utilities

This package provides database connection management, query utilities, and transaction handling for the PostEaze backend. It implements a clean abstraction layer over SQL database operations with support for connection pooling, transactions, and raw SQL queries.

## Architecture

The database utilities follow a layered architecture:

- **Connection Management**: Database initialization, connection pooling, and lifecycle management
- **Query Interface**: Abstract interface for database operations (Database interface)
- **Transaction Support**: Transactional operations with automatic rollback
- **Entity Pattern**: Raw entity interface for flexible query handling
- **Error Handling**: Standardized error types for common database scenarios

## Key Components

### Core Files

- **`database.go`**: Main database client initialization and connection management
- **`dao.go`**: Data Access Object implementation for regular database operations
- **`daotx.go`**: Transactional Data Access Object for transaction-based operations
- **`entity.go`**: Entity interfaces and structures for raw SQL operations

### Database Interface

The `Database` interface provides a consistent API for all database operations:

```go
type Database interface {
    QueryRaw(ctx context.Context, entity RawEntity, code int) error
    QueryMultiRaw(ctx context.Context, entity RawEntity, code int) ([]RawEntity, error)
    ExecRaws(ctx context.Context, source string, execs ...RawExec) error
    ExecRawsConsistent(ctx context.Context, source string, execs ...RawExec) error
}
```

## Connection Management

### Database Configuration

```go
type Config struct {
    DriverName            string        `json:"driverName"`
    URL                   string        `json:"url"`
    MaxOpenConnections    int           `json:"maxOpenConnections"`
    MaxIdleConnections    int           `json:"maxIdleConnections"`
    ConnectionMaxLifetime time.Duration `json:"connectionMaxLifetime"`
    ConnectionMaxIdleTime time.Duration `json:"connectionMaxIdleTime"`
}
```

### Connection Pooling

The database utilities implement connection pooling with configurable parameters:
- **Max Open Connections**: Maximum number of open connections to the database
- **Max Idle Connections**: Maximum number of idle connections in the pool
- **Connection Max Lifetime**: Maximum amount of time a connection may be reused
- **Connection Max Idle Time**: Maximum amount of time a connection may be idle

## Usage Examples

### Database Initialization

```go
import "github.com/shekhar8352/PostEaze/utils/database"

config := database.Config{
    DriverName:            "postgres",
    URL:                   "postgres://user:password@localhost/dbname",
    MaxOpenConnections:    25,
    MaxIdleConnections:    5,
    ConnectionMaxLifetime: time.Hour,
    ConnectionMaxIdleTime: time.Minute * 30,
}

err := database.Init(ctx, config)
if err != nil {
    log.Fatal("Failed to initialize database:", err)
}
defer database.Close()
```

### Basic Database Operations

```go
// Get database instance
db := database.Get()

// Single record query
err := db.QueryRaw(ctx, entity, queryCode)
if err != nil {
    if errors.Is(err, database.ErrNoRecords) {
        // Handle no records found
    }
    return err
}

// Multiple records query
entities, err := db.QueryMultiRaw(ctx, entity, queryCode)
if err != nil {
    return err
}
```

### Transaction Management

```go
// Begin transaction
tx, err := database.GetTx(ctx, nil)
if err != nil {
    return err
}

// Execute operations within transaction
err = tx.ExecRaws(ctx, "source", exec1, exec2, exec3)
if err != nil {
    database.RollbackTx(tx)
    return err
}

// Commit transaction
err = database.CommitTx(tx)
if err != nil {
    return err
}
```

### Consistent Execution

```go
// Execute with row count validation
err := db.ExecRawsConsistent(ctx, "source", execs...)
if err != nil {
    if errors.Is(err, database.ErrNoRowsAffected) {
        // Handle case where no rows were affected
    }
    return err
}
```

## Entity Pattern

The `RawEntity` interface allows flexible query handling:

```go
type RawEntity interface {
    GetQuery(code int) string
    GetQueryValues(code int) []any
    GetMultiQuery(code int) string
    GetMultiQueryValues(code int) []any
    GetNextRaw() RawEntity
    BindRawRow(code int, row Scanner) error
    GetExec(code int) string
    GetExecValues(code int, source string) []any
}
```

This pattern enables:
- **Dynamic Queries**: Different queries based on operation codes
- **Flexible Binding**: Custom row binding logic for each entity
- **Batch Operations**: Multiple operations in a single transaction

## Error Handling

The package defines standard error types:

- **`ErrNoRecords`**: No records found for the query
- **`ErrNoRowsAffected`**: No rows were affected by the operation

These errors can be checked using `errors.Is()` for proper error handling.

## Transaction Safety

- **Automatic Rollback**: Transactions are automatically rolled back on errors
- **Deferred Cleanup**: Proper resource cleanup using defer statements
- **Context Support**: All operations support context for cancellation and timeouts

## Dependencies

- **database/sql**: Standard Go database interface
- **Context**: For operation cancellation and timeouts

## Related Documentation

- [Backend Utilities Overview](../README.md) - Main utilities documentation
- [Entities](../../entities/README.md) - Data models and entity implementations
- [Models](../../models/README.md) - Data model definitions