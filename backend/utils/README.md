# Backend Utilities

This directory contains utility functions and helper modules that provide common functionality across the PostEaze backend application. These utilities handle cross-cutting concerns like authentication, HTTP requests, database operations, configuration management, and general helper functions.

## Architecture

The utilities are organized into specialized modules, each handling a specific domain of functionality:

- **Core Utilities**: General-purpose helper functions
- **Authentication**: JWT token handling and password hashing
- **HTTP Client**: Configurable HTTP request utilities
- **Database**: Database connection management and query utilities
- **Configuration**: Application configuration loading and management
- **Environment**: Environment variable handling
- **Command-line Flags**: Application startup flag processing

## Key Files

### Core Utility Files

- **`utils.go`**: Environment variable placeholder replacement utilities
- **`general_utils.go`**: Common HTTP response helpers for Gin framework
- **`auth_utils.go`**: Authentication utilities including JWT and password hashing
- **`jwt.go`**: Additional JWT token processing utilities
- **`logger.go`**: Logging utilities and configuration
- **`logapi_utils.go`**: API-specific logging helpers

### Utility Modules

- **`database/`**: Database connection management, DAO patterns, and transaction handling
- **`configs/`**: Configuration loading for development and production environments
- **`env/`**: Environment variable processing and template substitution
- **`flags/`**: Command-line flag parsing and application startup parameters
- **`http/`**: HTTP client utilities with configurable request handling

## Usage Examples

### General Response Utilities
```go
import "github.com/shekhar8352/PostEaze/utils"

// Send error response
utils.SendError(c, http.StatusBadRequest, "Invalid request")

// Send success response
utils.SendSuccess(c, userData, "User retrieved successfully")
```

### Authentication Utilities
```go
import "github.com/shekhar8352/PostEaze/utils"

// Hash password
hashedPassword, err := utils.HashPassword("plaintext")

// Generate JWT token
token, err := utils.GenerateJWT(userID, email)

// Parse JWT token
claims, err := utils.ParseJWT(tokenString)
```

### Environment Variable Processing
```go
import "github.com/shekhar8352/PostEaze/utils"

// Replace environment placeholders
processedString := utils.ReplacePlaceHoldersWithEnv("Database URL: ${DATABASE_URL}")
```

## Dependencies

The utilities depend on several external packages:
- **Gin**: Web framework for HTTP response utilities
- **JWT-Go**: JWT token processing
- **Bcrypt**: Password hashing
- **Godotenv**: Environment variable loading

## Related Documentation

- [Database Utilities](./database/README.md) - Database connection and query management
- [Configuration Utilities](./configs/README.md) - Application configuration handling
- [HTTP Utilities](./http/README.md) - HTTP client and request utilities
- [Environment Utilities](./env/README.md) - Environment variable processing
- [Flags Utilities](./flags/README.md) - Command-line flag handling