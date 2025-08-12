# Resources

Resource files and configuration management for the PostEaze backend. This directory contains all external configuration files, templates, and static resources needed by the application across different environments.

## Contents

The resources directory is organized to support multiple deployment environments:

```
resources/
└── configs/           # Configuration files organized by environment
    ├── dev/          # Development environment configurations
    ├── prod/         # Production environment configurations
    └── cug/          # Customer User Group (staging) environment configurations
```

## Resource Organization

### Configuration Files
The primary resource type is configuration files stored in YAML format. These files define environment-specific settings for:

- **Application Settings**: Basic application configuration and environment identification
- **Database Configuration**: Database connection parameters, connection pooling, and driver settings
- **API Configuration**: External API endpoints, timeout settings, and retry policies
- **Logger Configuration**: Logging levels and output parameters

### Environment Structure
Each environment directory (`dev/`, `prod/`, `cug/`) contains the same set of configuration files:

- `application.yml` - Core application settings
- `database.yml` - Database connection and pool configuration
- `api.yml` - External API service configurations

## Configuration Loading

The application uses a hierarchical configuration loading system:

1. **Environment Detection**: The application determines the target environment through command-line flags or environment variables
2. **Path Resolution**: Configuration path is resolved using the `base-config-path` flag (default: `resources/configs/dev`)
3. **File Loading**: Configuration files are loaded based on the constants defined in the `constants` package
4. **Environment Variables**: Configuration values can reference environment variables using `${VARIABLE_NAME}` syntax

## Environment-Specific Configurations

### Development (`dev/`)
- Local development settings
- Database connections to local PostgreSQL instances
- Debug-level logging
- Local API endpoints for testing

### Production (`prod/`)
- Production-ready settings
- Optimized database connection pools
- Production logging levels
- Live API endpoints

### Customer User Group (`cug/`)
- Staging environment settings
- Testing configurations
- Pre-production validation settings

## Configuration File Format

All configuration files use YAML format with consistent structure:

```yaml
env: "environment_name"

# Environment-specific settings
driverName: "postgres"
url: "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_URL}/${POSTGRES_DB}?sslmode=disable"
maxOpenConnections: 20
maxIdleConnections: 10

# API configurations with retry policies
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

## Usage Patterns

### Configuration Access
Configuration values are accessed through the configuration utilities:
```go
// Load database configuration
dbConfig := config.GetDatabaseConfig()
connectionString := dbConfig.URL
```

### Environment Variables
Configuration files support environment variable substitution:
```yaml
url: "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_URL}/${POSTGRES_DB}?sslmode=disable"
```

### Multi-Environment Deployment
The same application binary can be deployed to different environments by changing the configuration path:
```bash
# Development
./app --base-config-path=resources/configs/dev

# Production  
./app --base-config-path=resources/configs/prod
```

## Best Practices

1. **Environment Isolation**: Each environment has its own configuration directory
2. **Consistent Structure**: All environments use the same file names and structure
3. **Environment Variables**: Sensitive data (passwords, keys) are referenced via environment variables
4. **Version Control**: Configuration files are version controlled, but sensitive values are externalized
5. **Validation**: Configuration loading includes validation to catch errors early

## Related Documentation

- [Configuration Constants](../constants/README.md) - Constants used for configuration loading
- [Configuration Utilities](../utils/configs/README.md) - Utilities for loading and parsing configurations
- [Environment Utilities](../utils/env/README.md) - Environment variable handling
- [Flags Utilities](../utils/flags/README.md) - Command-line flag processing for configuration paths