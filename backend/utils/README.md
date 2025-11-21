# Backend Utilities

This directory contains a collection of utility functions and helper modules that provide common functionality across the PostEaze backend application.

## 📂 Sub-packages

Each sub-package has its own detailed documentation:

-   **[`configs`](./configs/README.md)**: Configuration loading and management (Dev/Prod).
-   **[`database`](./database/README.md)**: Database connection, transactions, and DAO patterns.
-   **[`env`](./env/README.md)**: Environment variable handling and substitution.
-   **[`flags`](./flags/README.md)**: Command-line flag parsing.
-   **[`http`](./http/README.md)**: HTTP client wrapper with retries and logging.
-   **[`redis`](./redis/redis.go)**: Redis client initialization and access.

---

## 🛠 Core Utilities (Root)

The root `utils/` directory contains essential helpers for authentication, logging, validation, and API responses.

### 🔐 Authentication & JWT

**Files:** `auth_utils.go`, `jwt.go`

Utilities for handling JSON Web Tokens (Access & Refresh) and password hashing.

| Function | Description |
| :--- | :--- |
| `GenerateJWT(userID int, email string)` | Generates a legacy JWT token (72h expiry). |
| `ParseJWT(tokenString string)` | Parses and validates a legacy JWT token. |
| `GenerateAccessToken(userID, role)` | Generates a short-lived (15m) access token. |
| `GenerateRefreshToken(userID)` | Generates a long-lived (7d) refresh token. |
| `ParseToken(tokenStr, isRefresh)` | Parses an access or refresh token. |
| `GetUserIDFromToken(tokenStr)` | Extracts the User ID from a token string. |

### 📝 Logging

**Files:** `logger.go`, `logapi_utils.go`

A structured, context-aware logging system with file rotation and search capabilities.

#### Structured Logger
```go
// Usage
utils.Logger.Info(ctx, "User logged in", "user_id", userID)
utils.Logger.Error(ctx, "Database connection failed", err)
```

| Function | Description |
| :--- | :--- |
| `NewLogger(config)` | Creates a new logger instance. |
| `AddLogID(ctx)` | Adds a unique Request ID to the context. |
| `GetLogID(ctx)` | Retrieves the Request ID from the context. |
| `Logger.Info/Debug/Warn/Error` | Logs messages with context and metadata. |

#### Log Search API
Utilities to query logs for the admin dashboard.

| Function | Description |
| :--- | :--- |
| `SearchLogs(query, date, level)` | Searches logs for keywords, optionally filtering by date/level. |
| `ReadLogsByDate(date, level)` | Reads all logs for a specific date. |
| `GetAvailableLogFiles()` | Returns a list of available log files. |

### ✅ Validation

**Files:** `validation.go`

Common validation logic for API inputs.

| Function | Description |
| :--- | :--- |
| `ValidateLogID(logID)` | Validates format and length of Log IDs. |
| `ValidateDate(dateStr)` | Validates date format (YYYY-MM-DD) and range. |
| `GetValidationErrorMessage(err)` | Returns user-friendly error messages. |

### 📡 API Responses (General)

**Files:** `general_utils.go`

Standardized JSON response helpers for Gin handlers.

| Function | Description |
| :--- | :--- |
| `SendSuccess(c, data, msg)` | Sends a 200 OK JSON response with data. |
| `SendError(c, code, msg)` | Sends a JSON error response. |
| `SendLogAPISuccess(c, data, msg)` | Sends a structured success response for Log APIs. |
| `SendLogAPIError(c, code, msg, type)` | Sends a structured error response for Log APIs. |
| `ContainsAny(slice, targets...)` | Checks if a slice contains any of the target strings. |

### 🔥 External Services

**Files:** `firebase.go`, `redis/redis.go`

Wrappers for external service clients.

| Function | Description |
| :--- | :--- |
| `InitializeFirebase(ctx)` | Initializes Firebase Admin SDK using env vars. |
| `GetFirebaseService()` | Returns the Firebase service instance. |
| `ValidateToken(ctx, idToken)` | Verifies a Firebase ID token. |
| `redis.Init(ctx)` | Initializes the Redis client. |
| `redis.GetClient()` | Returns the global Redis client. |