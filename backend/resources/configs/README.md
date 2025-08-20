# Configuration Files

Environment-specific configuration files for the PostEaze backend application. This directory contains YAML configuration files organized by deployment environment, providing a clean separation of settings for development, staging, and production deployments.

## Directory Structure

```
configs/
├── dev/              # Development environment
│   ├── application.yml
│   ├── database.yml
│   └── api.yml
├── prod/             # Production environment
│   ├── application.yml
│   ├── database.yml
│   └── api.yml
└── cug/              # Customer User Group (staging)
    ├── application.yml
    ├── database.yml
    └── api.yml
```

## Configuration Files

### application.yml
Core application settings and environment identification:
```yaml
env: "dev"  # Environment identifier used throughout the application
```

**Purpose**: Defines the basic application environment and can include application-wide settings like feature flags, application name, or global configuration parameters.

### database.yml
Database connection and pool configuration:
```yaml
env: "dev"

driverName: "postgres"
url: "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_URL}/${POSTGRES_DB}?sslmode=disable"
maxOpenConnections: 20
maxIdleConnections: 10
maxConnectionLifetimeInSeconds: 300
maxConnectionIdleTimeInSeconds: 60
```

**Configuration Parameters**:
- `driverName`: Database driver (postgres, mysql, etc.)
- `url`: Connection string with environment variable substitution
- `maxOpenConnections`: Maximum number of open database connections
- `maxIdleConnections`: Maximum number of idle connections in the pool
- `maxConnectionLifetimeInSeconds`: Maximum lifetime of a connection
- `maxConnectionIdleTimeInSeconds`: Maximum idle time before connection closure

### api.yml
External API service configurations with retry policies:
```yaml
env: "dev"

getCatsFact:
  method: "GET"
  url: "http://localhost:8082/account/details"
  timeoutInMillis: 1000
  retryCount: 3
  backoffPolicy:
    constantBackoff:
      intervalInMillis: 2
      maxJitterIntervalInMillis: 5
```

**Configuration Structure**:
- Service name (`getCatsFact`): Identifier for the external API
- `method`: HTTP method for the API call
- `url`: Endpoint URL for the external service
- `timeoutInMillis`: Request timeout in milliseconds
- `retryCount`: Number of retry attempts on failure
- `backoffPolicy`: Retry backoff strategy with timing parameters

## Environment-Specific Configurations

### Development Environment (`dev/`)
**Characteristics**:
- Local database connections
- Debug-friendly timeouts and connection limits
- Local or mock external API endpoints
- Verbose logging and debugging features enabled

**Typical Settings**:
- Lower connection pool limits for resource conservation
- Shorter timeouts for faster development feedback
- Local service URLs (localhost, docker containers)

### Production Environment (`prod/`)
**Characteristics**:
- Production database connections with optimized pools
- Production API endpoints
- Performance-optimized settings
- Security-hardened configurations

**Typical Settings**:
- Higher connection pool limits for performance
- Longer timeouts for network reliability
- Production service URLs with proper SSL/TLS
- Optimized retry policies for resilience

### Customer User Group Environment (`cug/`)
**Characteristics**:
- Staging environment for pre-production testing
- Production-like settings with test data
- Integration testing configurations
- User acceptance testing support

**Typical Settings**:
- Production-similar connection pools
- Test environment service URLs
- Balanced timeout and retry settings
- Feature flag configurations for testing

## Configuration Loading Process

1. **Environment Detection**: Application determines target environment via flags or environment variables
2. **Path Resolution**: Configuration directory resolved using `--base-config-path` flag
3. **File Discovery**: Configuration files loaded based on constants from `constants/configs.go`
4. **Variable Substitution**: Environment variables resolved using `${VARIABLE_NAME}` syntax
5. **Validation**: Configuration values validated against expected schemas
6. **Application**: Loaded configuration applied to application components

## Environment Variable Integration

Configuration files support environment variable substitution for sensitive or environment-specific values:

### Database Configuration
```yaml
url: "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_URL}/${POSTGRES_DB}?sslmode=disable"
```

**Required Environment Variables**:
- `POSTGRES_USER`: Database username
- `POSTGRES_PASSWORD`: Database password  
- `POSTGRES_URL`: Database host and port
- `POSTGRES_DB`: Database name

### API Configuration
```yaml
url: "${API_BASE_URL}/account/details"
```

**Environment-Specific Variables**:
- Development: `API_BASE_URL=http://localhost:8082`
- Production: `API_BASE_URL=https://api.posteaze.com`

## Configuration Management Best Practices

### Security
- **No Secrets in Files**: Sensitive data stored in environment variables, not configuration files
- **Environment Isolation**: Each environment has isolated configuration with appropriate access controls
- **Variable Validation**: Environment variables validated at application startup

### Maintainability
- **Consistent Structure**: All environments use identical file structure and naming
- **Documentation**: Configuration parameters documented with purpose and valid values
- **Version Control**: Configuration files tracked in version control (excluding sensitive data)

### Deployment
- **Environment Selection**: Configuration environment selected via command-line flags
- **Validation**: Configuration validated before application startup
- **Fallback Values**: Default values provided where appropriate

## Usage Examples

### Loading Configuration
```go
// Application startup with specific environment
./posteaze --base-config-path=resources/configs/prod --mode=release

// Development with default settings
./posteaze --mode=dev
```

### Accessing Configuration Values
```go
// Database configuration access
dbConfig := config.GetDatabaseConfig()
maxConnections := dbConfig.MaxOpenConnections

// API configuration access  
apiConfig := config.GetAPIConfig()
endpoint := apiConfig.GetCatsFact.URL
```

## Related Documentation

- [Resources Overview](../README.md) - Overall resource management strategy
- [Configuration Constants](../../constants/README.md) - Constants used for configuration loading
- [Configuration Utilities](../../utils/configs/README.md) - Configuration loading and parsing utilities
- [Environment Utilities](../../utils/env/README.md) - Environment variable handling
- [Database Utilities](../../utils/database/README.md) - Database connection management using these configurations