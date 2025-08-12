# Middleware

This folder contains middleware components that handle cross-cutting concerns in the PostEaze application's request processing pipeline. Middleware functions execute before route handlers and provide functionality like authentication, logging, and request processing.

## Architecture

The middleware follows the Gin framework's middleware pattern, where each middleware is a `gin.HandlerFunc` that can:
- Process incoming requests before they reach route handlers
- Modify request/response data
- Terminate request processing early (e.g., for authentication failures)
- Pass control to the next middleware or handler in the chain

## Request Processing Pipeline

The middleware is applied in the following order:
1. **GinLoggingMiddleware** - Logs all incoming requests and responses
2. **AuthMiddleware** - Validates JWT tokens for protected routes
3. **RequireRole** - Enforces role-based access control
4. Route handlers execute after middleware validation

## Key Components

### Authentication Middleware (`auth_middleware.go`)

Provides JWT-based authentication and role-based authorization:

- **AuthMiddleware()**: Validates Bearer tokens in Authorization headers
  - Extracts JWT token from `Authorization: Bearer <token>` header
  - Validates token using `utils.ParseToken()`
  - Sets `user_id` and `role` in Gin context for downstream handlers
  - Returns 401 Unauthorized for invalid/missing tokens

- **RequireRole(allowedRoles...)**: Enforces role-based access control
  - Checks if user's role matches any of the allowed roles
  - Must be used after AuthMiddleware to access role information
  - Returns 403 Forbidden for insufficient permissions

#### Usage Example
```go
// Protect route with authentication
authv1.POST("/refresh", middleware.AuthMiddleware(), apiv1.RefreshTokenHandler)

// Protect route with role-based access
adminRoutes.GET("/users", middleware.AuthMiddleware(), middleware.RequireRole("admin"), handler)
```

### Logging Middleware (`log_middleware.go`)

Provides comprehensive request/response logging with correlation IDs:

- **GinLoggingMiddleware()**: Logs HTTP requests and responses
  - Generates unique log ID for request correlation
  - Logs request start with method, path, IP, and User-Agent
  - Logs request completion with status code and duration
  - Extracts real client IP from headers (X-Forwarded-For, X-Real-IP)
  - Integrates with the application's structured logging system

#### Features
- **Request Correlation**: Each request gets a unique log ID for tracing
- **Client IP Detection**: Handles proxied requests correctly
- **Performance Monitoring**: Tracks request duration
- **Structured Logging**: Uses the application's logger with context

#### Usage Example
```go
// Applied globally to all routes
s := gin.Default()
s.Use(middleware.GinLoggingMiddleware())
```

## Security Patterns

### JWT Token Validation
- Tokens must be provided in `Authorization: Bearer <token>` format
- Token parsing and validation handled by `utils.ParseToken()`
- Invalid tokens result in immediate request termination
- User context (ID and role) extracted and made available to handlers

### Role-Based Access Control (RBAC)
- Roles are embedded in JWT tokens
- `RequireRole` middleware enforces access control
- Supports multiple allowed roles per route
- Fails securely with 403 Forbidden for unauthorized access

### Request Logging Security
- Sensitive headers are not logged
- Client IP detection handles proxy scenarios
- Log correlation IDs help with security incident investigation

## Integration with Application

### Router Integration
Middleware is registered in `api/router.go`:
```go
// Global logging middleware
s.Use(middleware.GinLoggingMiddleware())

// Route-specific authentication
authv1.POST("/refresh", middleware.AuthMiddleware(), apiv1.RefreshTokenHandler)
```

### Context Data
Middleware sets the following data in Gin context:
- `user_id`: Authenticated user's ID (from AuthMiddleware)
- `role`: User's role for authorization (from AuthMiddleware)
- Request context includes log ID for correlation

## Dependencies

- **Gin Framework**: Web framework providing middleware support
- **utils Package**: Token parsing and logging utilities
- **JWT**: Token-based authentication system

## Error Handling

### Authentication Errors
- Missing Authorization header: 401 Unauthorized
- Invalid token format: 401 Unauthorized  
- Expired/invalid token: 401 Unauthorized
- Missing role in token: 403 Forbidden
- Insufficient role permissions: 403 Forbidden

### Logging Errors
- Logging middleware is designed to be non-blocking
- Logging failures don't interrupt request processing
- Uses structured logging for error tracking

## Best Practices

1. **Middleware Order**: Apply logging before authentication for complete request tracking
2. **Error Responses**: Use consistent JSON error format across middleware
3. **Context Usage**: Store user data in Gin context for handler access
4. **Security**: Always validate tokens before trusting user data
5. **Performance**: Middleware should be lightweight to avoid request latency

## Related Documentation

- [API Layer](../api/README.md) - How middleware integrates with routing
- [Utils](../utils/README.md) - Token parsing and logging utilities
- [Business Logic](../business/README.md) - How handlers use middleware context