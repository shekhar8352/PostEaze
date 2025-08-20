# Business Logic v1

This folder contains the version 1 implementation of PostEaze's core business logic services. These services implement the domain-specific operations and business rules for the application's primary features.

## Contents

### Authentication Services (`auth.go`)
Handles all user authentication and authorization operations:

- **Signup**: User registration with team creation support
- **Login**: User authentication with token generation
- **RefreshToken**: JWT token refresh mechanism
- **Logout**: User session termination and token revocation

### Logging Services (`log.go`)
Manages system log retrieval and analysis:

- **ReadLogsByLogID**: Retrieves logs by correlation ID across multiple days
- **ReadLogsByDate**: Fetches all logs for a specific date

## Service Responsibilities

### Authentication Service Functions

#### `Signup(ctx context.Context, params modelsv1.SignupParams)`
**Purpose**: Registers new users and creates associated resources

**Business Logic**:
- Hashes user passwords securely
- Creates user records in the database
- For team users: creates team and assigns admin role
- Generates access and refresh tokens
- Manages database transactions for data consistency

**Returns**: User details with authentication tokens

#### `Login(ctx context.Context, params modelsv1.LoginParams)`
**Purpose**: Authenticates existing users

**Business Logic**:
- Validates user credentials against stored data
- Verifies password using secure hash comparison
- Generates new authentication tokens
- Logs authentication attempts and outcomes

**Returns**: User details with fresh authentication tokens

#### `RefreshToken(ctx context.Context, token string)`
**Purpose**: Refreshes expired access tokens

**Business Logic**:
- Validates refresh token authenticity
- Retrieves user information from token
- Generates new access token
- Maintains session continuity

**Returns**: New access token

#### `Logout(ctx context.Context, refreshToken string)`
**Purpose**: Terminates user sessions

**Business Logic**:
- Validates refresh token
- Revokes all tokens for the user
- Cleans up session data

### Logging Service Functions

#### `ReadLogsByLogID(ctx context.Context, logID string)`
**Purpose**: Retrieves correlated logs across multiple days

**Business Logic**:
- Searches logs from the last 3 days
- Filters logs by correlation ID (logID)
- Handles missing log files gracefully
- Sorts results chronologically

**Returns**: Ordered list of log entries

#### `ReadLogsByDate(date string)`
**Purpose**: Retrieves all logs for a specific date

**Business Logic**:
- Constructs log file path based on date
- Validates file existence
- Reads and parses log entries
- Returns total count and entries

**Returns**: Log entries and total count

## Implementation Patterns

### Transaction Management
Complex operations use database transactions:

```go
tx, err := database.GetTx(ctx, nil)
if err != nil {
    return nil, err
}
// ... perform multiple operations
err = database.CommitTx(tx)
```

### Error Handling and Logging
All operations include comprehensive error handling:

```go
utils.Logger.Info(ctx, "Attempting to login user with email: %s", params.Email)
if err != nil {
    utils.Logger.Error(ctx, "Error validating password for user with email: %s", params.Email)
    return nil, errors.New("invalid credentials")
}
```

### Security Practices
- Password hashing using secure algorithms
- JWT token generation with proper expiration
- Input validation through model parameters
- Context-aware logging for security auditing

## Data Flow

### Authentication Flow
1. API receives request with user credentials
2. Business service validates input parameters
3. Repository layer queries user data
4. Password verification using utils
5. Token generation and storage
6. Response with user data and tokens

### Logging Flow
1. API receives log query request
2. Business service determines search parameters
3. File system operations to read log files
4. Log parsing and filtering
5. Result aggregation and sorting
6. Response with formatted log data

## Dependencies

### Internal Dependencies
- `modelsv1`: Data structures and validation
- `repositories`: Data access operations
- `utils`: Utility functions (hashing, tokens, logging)
- `database`: Transaction management

### External Dependencies
- `context`: Request context propagation
- `errors`: Error handling
- Standard library packages for file operations

## Error Scenarios

### Authentication Errors
- Invalid credentials
- Token generation failures
- Database transaction failures
- User not found scenarios

### Logging Errors
- Missing log files
- File read permissions
- Log parsing failures
- Invalid date formats

## Usage Examples

### From API Handlers
```go
// User signup
user, err := businessv1.Signup(c.Request.Context(), signupParams)
if err != nil {
    utils.SendError(c, http.StatusInternalServerError, err.Error())
    return
}

// Log retrieval
logs, err := businessv1.ReadLogsByLogID(c.Request.Context(), logID)
if err != nil {
    utils.Logger.Info(c.Request.Context(), "Error reading log by ID: ", err)
    utils.SendError(c, http.StatusInternalServerError, err.Error())
    return
}
```

## Related Documentation

- [Business Layer Overview](../README.md) - Architecture and patterns
- [API v1](../../api/v1/README.md) - HTTP handlers that consume these services
- [Models v1](../../models/v1/README.md) - Data structures used by these services
- [Repositories](../../entities/repositories/README.md) - Data access layer