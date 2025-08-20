# Business Logic Layer

The business logic layer contains the core application logic and business rules for PostEaze. This layer acts as an intermediary between the API handlers and the data access layer, implementing the business processes and domain-specific operations.

## Architecture

The business layer follows a service-oriented architecture pattern where:

- **Service Functions**: Stateless functions that implement specific business operations
- **Transaction Management**: Handles database transactions for complex operations
- **Business Rule Enforcement**: Validates business constraints and rules
- **Domain Logic**: Implements core application functionality independent of external interfaces

## Key Components

- **Authentication Services**: User registration, login, token management, and logout operations
- **Logging Services**: Log retrieval and filtering operations for system monitoring
- **Transaction Coordination**: Manages database transactions across multiple repository operations

## Service Layer Patterns

### Function-Based Services
The business layer uses function-based services rather than class-based services:

```go
func Signup(ctx context.Context, params modelsv1.SignupParams) (map[string]interface{}, error)
func Login(ctx context.Context, params modelsv1.LoginParams) (map[string]interface{}, error)
```

### Context Propagation
All business functions accept a `context.Context` as the first parameter for:
- Request tracing and logging correlation
- Timeout and cancellation handling
- Database transaction context

### Transaction Management
Complex operations use database transactions to ensure data consistency:

```go
tx, err := database.GetTx(ctx, nil)
// ... perform multiple repository operations
err = database.CommitTx(tx)
```

### Error Handling
Business functions return errors that are:
- Logged with appropriate context
- Propagated to the API layer for proper HTTP response handling
- Wrapped with additional business context when needed

## Business Rules Implementation

### User Registration
- Password hashing before storage
- Team creation for team-type users
- Automatic admin role assignment for team creators
- Token generation for immediate authentication

### Authentication Flow
- Email-based user lookup
- Password verification using secure hashing
- JWT token generation (access and refresh tokens)
- Token storage and management

### Logging Operations
- Multi-day log file reading
- Log filtering by correlation ID
- Chronological sorting of log entries

## Dependencies

The business layer depends on:
- **Models**: Data structures and validation (`models/v1`)
- **Repositories**: Data access layer (`entities/repositories`)
- **Utilities**: Common functions for hashing, tokens, logging (`utils`)
- **Database**: Transaction management (`utils/database`)

## Usage Patterns

Business functions are called from API handlers:

```go
// In API handler
user, err := businessv1.Signup(c.Request.Context(), signupParams)
if err != nil {
    // Handle error and return appropriate HTTP response
}
```

## Version Organization

The business layer is organized by API version:
- **v1/**: Version 1 business logic implementation

## Related Documentation

- [API Layer](../api/README.md) - HTTP handlers that call business services
- [Entities](../entities/README.md) - Data models and repository interfaces
- [Models](../models/README.md) - Request/response data structures
- [Utils](../utils/README.md) - Shared utility functions