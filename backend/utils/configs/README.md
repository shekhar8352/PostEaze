# Configuration Utilities

This package provides configuration management utilities for the PostEaze backend application. It supports both development and production environments with different configuration sources, including file-based configurations for development and AWS App Config for production deployments.

## Architecture

The configuration utilities implement a flexible configuration loading system:

- **Multi-Environment Support**: Different initialization methods for development and production
- **Multiple Providers**: File-based configurations and AWS App Config
- **YAML Configuration**: Standardized YAML format for all configuration files
- **Centralized Access**: Single client interface for accessing configurations

## Key Components

### Core Files

- **`configs.go`**: Main configuration client initialization and management

### Configuration Providers

The package supports two primary configuration providers:

1. **File-Based Provider** (`configs.FileBased`): For development environments
2. **AWS App Config Provider** (`configs.AWSAppConfig`): For production environments

## Usage Examples

### Development Environment Initialization

```go
import "github.com/shekhar8352/PostEaze/utils/configs"

// Initialize for development mode
err := configs.InitDev("/path/to/config/directory", "app", "database", "redis")
if err != nil {
    log.Fatal("Failed to initialize dev configs:", err)
}

// Get configuration client
client := configs.Get()
```

### Production Environment Initialization

```go
import "github.com/shekhar8352/PostEaze/utils/configs"

// Initialize for production mode with AWS App Config
err := configs.InitRelease("production", "us-west-2", "app", "database", "redis")
if err != nil {
    log.Fatal("Failed to initialize release configs:", err)
}

// Get configuration client
client := configs.Get()
```

### Configuration Access

```go
// Get configuration client
client := configs.Get()

// Access configuration values (example usage)
// Note: Specific methods depend on the go-config-client library implementation
```

## Configuration Parameters

### Development Mode Parameters

- **`constants.ConfigDirectoryKey`**: Directory path containing configuration files
- **`constants.ConfigNamesKey`**: Array of configuration file names (without extension)
- **`constants.ConfigTypeKey`**: Configuration file type (YAML)

### Production Mode Parameters

- **`constants.ConfigIDKey`**: Application identifier
- **`constants.ConfigRegionKey`**: AWS region for App Config
- **`constants.ConfigEnvKey`**: Environment name (e.g., "production", "staging")
- **`constants.ConfigAppKey`**: Application name
- **`constants.ConfigNamesKey`**: Array of configuration names
- **`constants.ConfigTypeKey`**: Configuration type (YAML)
- **`constants.ConfigCredentialsModeKey`**: AWS credentials mode

## Environment-Specific Configuration

### Development Configuration

For development environments, configurations are loaded from local YAML files:

```yaml
# app.yaml
server:
  port: 8080
  host: localhost

database:
  host: localhost
  port: 5432
  name: postease_dev
```

### Production Configuration

For production environments, configurations are retrieved from AWS App Config:
- Configurations are stored in AWS App Config service
- Supports environment-specific configurations
- Automatic configuration updates and rollbacks
- Integrated with AWS IAM for security

## Configuration Loading Patterns

### Multi-File Configuration

The system supports loading multiple configuration files:

```go
// Load multiple configuration files
err := configs.InitDev("/config", "app", "database", "redis", "auth")
```

Each configuration file can contain related settings:
- **app.yaml**: Application-specific settings
- **database.yaml**: Database connection settings
- **redis.yaml**: Redis cache settings
- **auth.yaml**: Authentication settings

### Environment-Specific Settings

Configuration values can be environment-specific:

```yaml
# Development
database:
  host: localhost
  port: 5432

# Production (via AWS App Config)
database:
  host: prod-db.amazonaws.com
  port: 5432
```

## Error Handling

The configuration utilities handle various error scenarios:
- **File Not Found**: Missing configuration files in development
- **AWS Connection Errors**: Network or authentication issues with AWS App Config
- **Invalid YAML**: Malformed configuration files
- **Missing Parameters**: Required configuration parameters not provided

## Dependencies

The package depends on:
- **go-config-client**: Third-party configuration client library
- **Constants Package**: Application constants for configuration keys
- **AWS SDK**: For AWS App Config integration (in production mode)

## Constants Integration

The package uses constants from the application constants package:

```go
// Example constants used
constants.ConfigDirectoryKey    // "directory"
constants.ConfigNamesKey       // "names"
constants.ConfigTypeKey        // "type"
constants.ConfigYAML          // "yaml"
constants.ApplicationName     // "PostEaze"
```

## Best Practices

### Development Setup

1. **Organize Configuration Files**: Keep configuration files in a dedicated directory
2. **Use Descriptive Names**: Name configuration files based on their purpose
3. **Environment Variables**: Use environment variables for sensitive data
4. **Version Control**: Include sample configuration files in version control

### Production Setup

1. **AWS App Config**: Use AWS App Config for centralized configuration management
2. **Environment Separation**: Maintain separate configurations for different environments
3. **Security**: Use AWS IAM roles for secure access to configurations
4. **Monitoring**: Monitor configuration changes and deployments

## Related Documentation

- [Backend Utilities Overview](../README.md) - Main utilities documentation
- [Environment Utilities](../env/README.md) - Environment variable handling
- [Flags Utilities](../flags/README.md) - Command-line flag processing
- [Constants](../../constants/README.md) - Application constants