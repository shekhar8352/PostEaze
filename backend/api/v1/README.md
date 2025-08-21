# API Version 1

This directory contains the Version 1 REST API endpoints for PostEaze. All endpoints are prefixed with `/api/v1` and organized by feature domain. The API follows RESTful conventions and uses JSON for request/response payloads.

## Architecture

The v1 API is organized into feature-based modules:

```
/v1
├── auth.go            # Authentication and user management endpoints
└── log.go             # Application logging and monitoring endpoints
```

## Authentication Endpoints (`auth.go`)

Base path: `/api/v1/auth`

### POST /signup
**Purpose**: Register a new user account

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "User Name"
}
```

**Response**: User object with authentication tokens
**Status Codes**: 200 (success), 400 (invalid data), 500 (server error)

### POST /login  
**Purpose**: Authenticate existing user and generate tokens

**Request Body**:
```json
{
  "email": "user@example.com", 
  "password": "securepassword"
}
```

**Response**: User object with JWT access and refresh tokens
**Status Codes**: 200 (success), 400 (missing credentials), 401 (invalid credentials)

### POST /refresh
**Purpose**: Generate new access token using refresh token
**Authentication**: Required (Bearer token)

**Request Body**:
```json
{
  "refreshToken": "jwt_refresh_token_here"
}
```

**Response**: New user object with updated tokens
**Status Codes**: 200 (success), 400 (missing token), 401 (invalid token)

### POST /logout
**Purpose**: Invalidate user session and tokens
**Authentication**: Required (Bearer token)

**Request Body**: None (uses Authorization header)
**Response**: Success confirmation
**Status Codes**: 200 (success), 500 (server error)

## Logging Endpoints (`log.go`)

Base path: `/api/v1/log`

### GET /byId/:log_id
**Purpose**: Retrieve a specific log entry by its unique identifier

**Parameters**:
- `log_id` (path parameter): Unique identifier for the log entry

**Response**: Log entry object with details
**Status Codes**: 200 (success), 400 (missing ID), 500 (server error)

**Example**: `GET /api/v1/log/byId/12345`

### GET /byDate/:date
**Purpose**: Retrieve all log entries for a specific date

**Parameters**:
- `date` (path parameter): Date in YYYY-MM-DD format

**Response**: 
```json
{
  "logs": [...],
  "total": 25
}
```

**Status Codes**: 200 (success), 400 (invalid date format), 500 (server error)

**Example**: `GET /api/v1/log/byDate/2024-01-15`

## Request/Response Patterns

### Standard Response Format
All endpoints return responses in a consistent format:

```json
{
  "data": {...},
  "message": "Operation completed successfully",
  "status": "success"
}
```

### Error Response Format
```json
{
  "error": "Error description",
  "status": "error"
}
```

### Authentication
Protected endpoints require JWT token in Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Handler Implementation

### Request Processing Flow
1. **Request Binding**: Parse and validate JSON request body
2. **Business Logic**: Delegate to appropriate business layer service
3. **Response Formatting**: Use utility functions for consistent responses
4. **Error Handling**: Log errors and return appropriate HTTP status codes
5. **Success Response**: Return formatted success response with data

### Validation
- Request body validation using Gin's `ShouldBindJSON`
- Required field validation (email, password, etc.)
- Date format validation for log endpoints
- Token validation for protected endpoints

### Logging
- Request/response logging for debugging
- Error logging with context information
- Success operation logging for monitoring

## Dependencies

- **Gin Context**: HTTP request/response handling
- **Business Layer v1**: Core application logic (`businessv1` package)
- **Models v1**: Request/response data structures (`modelsv1` package)
- **Utils**: Response formatting and logging utilities
- **Middleware**: Authentication validation

## Error Handling

### Common Error Scenarios
- **400 Bad Request**: Invalid JSON, missing required fields, invalid date format
- **401 Unauthorized**: Invalid credentials, expired tokens, missing authentication
- **500 Internal Server Error**: Database errors, business logic failures

### Logging Strategy
- Info level: Successful operations and request details
- Error level: Server errors and business logic failures
- Context preservation: Request context passed through all operations

## Usage Examples

### Adding New Endpoint
```go
func NewEndpointHandler(c *gin.Context) {
    var body modelsv1.NewEndpointParams
    if err := c.ShouldBindJSON(&body); err != nil {
        utils.SendError(c, http.StatusBadRequest, "Invalid request data")
        return
    }

    result, err := businessv1.NewEndpointLogic(c.Request.Context(), body)
    if err != nil {
        utils.SendError(c, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SendSuccess(c, result, "Operation completed successfully")
}
```

## Related Documentation

- [API Layer Overview](../README.md) - Router configuration and architecture
- [Business Layer v1](../../business/v1/README.md) - Business logic implementation
- [Models v1](../../models/v1/README.md) - Data structures and validation
- [Middleware](../../middleware/README.md) - Authentication and request processing