# PostEaze Backend

The PostEaze backend is a Go-based REST API server that provides authentication, logging, and social media management functionality. Built with the Gin web framework, it follows a layered architecture pattern with clear separation of concerns between API routing, business logic, data access, and infrastructure components.

## Architecture Overview

The backend follows a clean architecture approach with the following layers:

```
┌─────────────────────────────────────────────────────────────┐
│                     API Layer (Gin Router)                  │
│                    /api/v1/* endpoints                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                  Middleware Layer                           │
│            Authentication, Logging, CORS                    │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                  Business Logic Layer                       │
│              Service implementations (v1)                   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                   Data Access Layer                         │
│            Entities, Repositories, Models                   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                  Infrastructure Layer                       │
│          Database, Config, Utils, Constants                 │
└─────────────────────────────────────────────────────────────┘
```

## Request Flow

1. **HTTP Request** → Gin Router (`api/router.go`)
2. **Middleware Processing** → Authentication, Logging (`middleware/`)
3. **Route Matching** → API v1 handlers (`api/v1/`)
4. **Business Logic** → Service layer (`business/v1/`)
5. **Data Access** → Repository pattern (`entities/repositories/`)
6. **Database** → PostgreSQL via `lib/pq` driver

## Key Components

### Core Application Files
- **`main.go`** - Application entry point and service initialization
- **`go.mod`** - Go module definition and dependencies
- **`Dockerfile`** - Container configuration for deployment

### Architecture Layers

#### API Layer (`api/`)
- **`router.go`** - Main router configuration and route grouping
- **`v1/`** - Version 1 API handlers and endpoints
- Handles HTTP request/response processing and routing

#### Business Logic (`business/`)
- **`v1/`** - Version 1 business logic implementations
- Contains service layer with core application logic
- Orchestrates data operations and business rules

#### Data Layer
- **`entities/`** - Domain entities and repository interfaces
- **`models/`** - Data models and structures (versioned)
- **`migrations/`** - Database schema migrations

#### Infrastructure
- **`middleware/`** - HTTP middleware components (auth, logging)
- **`utils/`** - Utility functions and helper modules
- **`constants/`** - Application constants and configuration keys
- **`resources/`** - Configuration files and static resources

## Service Initialization

The `main.go` file initializes services in the following order:

```go
func main() {
    ctx := context.Background()
    initEnv()        // Load environment variables
    initConfigs(ctx) // Initialize configuration management
    initDatabase(ctx)// Set up database connection pool
    initRouter(ctx)  // Configure HTTP router and start server
    initHttp(ctx)    // Initialize HTTP client utilities
}
```

### Environment & Configuration
- **Environment Variables**: Loaded via `utils/env` package
- **Configuration Management**: Multi-environment config support (dev/release)
- **AWS Integration**: Supports AWS AppConfig and Secrets Manager in release mode

### Database Connection
- **Driver**: PostgreSQL (`lib/pq`)
- **Connection Pooling**: Configurable max connections and timeouts
- **Configuration**: Database URL and connection parameters from config files

### HTTP Server
- **Framework**: Gin web framework
- **Middleware**: Logging and authentication middleware applied globally
- **Routing**: Grouped routes with versioning (`/api/v1/*`)

## Key Dependencies

### Core Framework
- **`gin-gonic/gin`** - HTTP web framework
- **`lib/pq`** - PostgreSQL driver

### Authentication & Security
- **`golang-jwt/jwt/v5`** - JWT token handling
- **`google/uuid`** - UUID generation

### Configuration & Environment
- **`joho/godotenv`** - Environment variable loading
- **`sinhashubham95/go-config-client`** - Advanced configuration management
- **AWS SDK v2** - Cloud configuration and secrets management

## API Endpoints

### Health Check
- **GET** `/api/health` - Service health status

### Authentication (`/api/v1/auth`)
- **POST** `/signup` - User registration
- **POST** `/login` - User authentication
- **POST** `/refresh` - Token refresh (requires auth)
- **POST** `/logout` - User logout (requires auth)

### Logging (`/api/v1/log`)
- **GET** `/byDate/:date` - Retrieve logs by date
- **GET** `/byId/:log_id` - Retrieve specific log entry

## Development Setup

### Prerequisites
- Go 1.23 or higher
- PostgreSQL database
- Environment variables configured (see `.env` file)

### Running the Application
```bash
# Install dependencies
go mod download

# Run database migrations
# (See migrations/ folder for SQL files)

# Start the server
go run main.go
```

### Configuration
The application supports two modes:
- **Development Mode**: Uses local config files from `resources/configs/`
- **Release Mode**: Uses AWS AppConfig and environment variables

## Security Features

- **JWT Authentication**: Stateless token-based authentication
- **Middleware Protection**: Routes protected by authentication middleware
- **Environment Variable Security**: Sensitive data loaded from environment
- **Database Connection Security**: Connection pooling with timeout controls

## Logging

- **Request Logging**: All HTTP requests logged via Gin middleware
- **Application Logging**: Structured logging with date-based log files
- **Log Storage**: Daily log files in `logs/` directory
- **Log API**: Endpoints to retrieve and query log data

## Related Documentation

- [`api/`](api/) - API layer and routing documentation
- [`business/`](business/) - Business logic layer documentation
- [`entities/`](entities/) - Data models and repository patterns
- [`middleware/`](middleware/) - HTTP middleware components
- [`utils/`](utils/) - Utility functions and helpers
- [`migrations/`](migrations/) - Database schema and migration procedures