# API Layer

The API layer serves as the HTTP interface for the PostEaze backend, handling REST API requests and routing them to appropriate business logic handlers. Built using the Gin web framework, this layer provides a clean separation between HTTP concerns and business logic.

## Architecture

The API layer follows a versioned REST architecture with the following structure:

```
/api
├── router.go          # Main router configuration and route registration
└── v1/                # Version 1 API endpoints
    ├── auth.go        # Authentication endpoints
    └── log.go         # Logging endpoints
```

## Key Components

- **router.go**: Central router configuration that initializes the Gin server, applies middleware, and registers all API routes
- **v1/**: Version 1 API handlers organized by feature domain
- **Route Groups**: Logical grouping of related endpoints (auth, logs, etc.)

## Router Configuration

The main router (`router.go`) handles:

### Server Initialization
- Creates Gin server instance with default configuration
- Applies global middleware (logging, CORS, etc.)
- Configures route groups for API versioning

### Route Structure
All API routes follow the pattern: `/api/v{version}/{feature}/{endpoint}`

Example routes:
- `/api/health` - Health check endpoint
- `/api/v1/auth/signup` - User registration
- `/api/v1/auth/login` - User authentication
- `/api/v1/log/byDate/:date` - Retrieve logs by date

### Middleware Integration
- **GinLoggingMiddleware**: Request/response logging
- **AuthMiddleware**: JWT token validation for protected routes
- Applied selectively based on endpoint security requirements

## Route Registration

Routes are organized into logical groups using Gin's RouterGroup:

```go
api := s.Group(constants.ApiRoute)           // /api
v1 := api.Group(constants.V1Route)           // /api/v1
authv1 := v1.Group(constants.AuthRoute)      // /api/v1/auth
logv1 := v1.Group(constants.LogRoute)        // /api/v1/log
```

### Authentication Routes
- `POST /signup` - User registration
- `POST /login` - User authentication  
- `POST /refresh` - Token refresh (protected)
- `POST /logout` - User logout (protected)

### Logging Routes
- `GET /byDate/:date` - Retrieve logs by date
- `GET /byId/:log_id` - Retrieve specific log entry

## Error Handling

The API layer uses consistent error handling patterns:
- HTTP status codes follow REST conventions
- Error responses include descriptive messages
- Logging for debugging and monitoring
- Graceful handling of validation errors

## Dependencies

- **Gin Framework**: HTTP router and middleware
- **Business Layer**: Core application logic
- **Models**: Request/response data structures
- **Middleware**: Authentication and logging
- **Constants**: Route definitions and configuration
- **Utils**: Response formatting and logging utilities

## Usage Example

```go
// Adding new route group
func addV1NewFeatureRoutes(v1 *gin.RouterGroup) {
    featurev1 := v1.Group(constants.NewFeatureRoute)
    featurev1.GET("/endpoint", handler.NewFeatureHandler)
}

// Register in Init() function
addV1NewFeatureRoutes(v1)
```

## Related Documentation

- [API v1 Documentation](./v1/README.md) - Detailed endpoint documentation
- [Business Layer](../business/README.md) - Business logic implementation
- [Middleware](../middleware/README.md) - Request processing middleware
- [Models](../models/README.md) - Data structures and validation