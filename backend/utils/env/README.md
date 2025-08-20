# Environment Utilities

This package provides environment variable handling and template substitution utilities for the PostEaze backend application. It supports loading environment variables from `.env` files and performing template-based string substitution with environment values.

## Architecture

The environment utilities implement a simple but powerful environment management system:

- **Environment Loading**: Automatic loading of `.env` files using godotenv
- **Variable Caching**: In-memory caching of environment variables for performance
- **Template Substitution**: String template processing with environment variable replacement
- **Flexible Parsing**: Support for complex environment variable formats

## Key Components

### Core Files

- **`env.go`**: Main environment variable processing and template substitution

### Core Functions

- **`InitEnv()`**: Initialize environment variable loading and caching
- **`ApplyEnvironmentToString(value string) string`**: Apply environment variable substitution to strings

## Usage Examples

### Environment Initialization

```go
import "github.com/shekhar8352/PostEaze/utils/env"

// Initialize environment variables
// This loads .env file and caches all environment variables
env.InitEnv()
```

### Template String Processing

```go
import "github.com/shekhar8352/PostEaze/utils/env"

// Initialize environment first
env.InitEnv()

// Apply environment variable substitution
template := "Database URL: ${DATABASE_URL}, Port: ${PORT}"
processed := env.ApplyEnvironmentToString(template)
// Result: "Database URL: postgres://localhost:5432/postease, Port: 8080"
```

### Configuration Template Processing

```go
// Process configuration strings with environment variables
configTemplate := `
server:
  host: ${SERVER_HOST}
  port: ${SERVER_PORT}
database:
  url: ${DATABASE_URL}
  max_connections: ${DB_MAX_CONNECTIONS}
`

processedConfig := env.ApplyEnvironmentToString(configTemplate)
```

## Environment Variable Loading

### .env File Support

The package automatically loads environment variables from `.env` files in the current directory:

```env
# .env file example
DATABASE_URL=postgres://localhost:5432/postease
SERVER_HOST=localhost
SERVER_PORT=8080
JWT_SECRET=your-secret-key
REDIS_URL=redis://localhost:6379
```

### Environment Variable Parsing

The environment loading process:

1. **Load .env File**: Uses `godotenv.Load()` to load variables from `.env` file
2. **Parse System Environment**: Reads all system environment variables
3. **Cache Variables**: Stores variables in an in-memory map for fast access
4. **Handle Complex Values**: Properly handles environment variables with `=` signs in values

## Template Substitution

### Substitution Pattern

The package supports `${VARIABLE_NAME}` pattern for environment variable substitution:

```go
// Template patterns supported
"${DATABASE_URL}"           // Simple variable
"${SERVER_HOST}:${PORT}"    // Multiple variables
"prefix-${ENV}-suffix"      // Variables with prefix/suffix
```

### Substitution Process

1. **Pattern Matching**: Finds all `${VARIABLE_NAME}` patterns in the input string
2. **Variable Lookup**: Looks up each variable in the cached environment map
3. **String Replacement**: Replaces patterns with actual environment values
4. **Fallback Handling**: Leaves unmatched patterns unchanged

## Environment Variable Handling

### Variable Caching

Environment variables are cached in memory for performance:

```go
var envObj map[string]string  // In-memory cache of environment variables
```

### Complex Value Support

The package properly handles environment variables with complex values:

```env
# Environment variables with equals signs
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
JWT_SECRET=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ
```

### Environment Variable Parsing Logic

```go
// Parsing logic handles complex values
for _, e := range env {
    s := strings.Split(e, "=")
    if len(s) >= 2 {
        key := strings.TrimSpace(s[0])
        value := strings.TrimSpace(strings.Join(s[1:], "="))  // Rejoin for complex values
        envObj[key] = value
    }
}
```

## Integration Patterns

### Configuration Processing

```go
// Use with configuration utilities
configString := utils.ReplacePlaceHoldersWithEnv("${DATABASE_URL}")
processedConfig := env.ApplyEnvironmentToString(configString)
```

### Application Startup

```go
func main() {
    // Initialize environment variables early in application startup
    env.InitEnv()
    
    // Process configuration templates
    dbURL := env.ApplyEnvironmentToString("${DATABASE_URL}")
    
    // Continue with application initialization
}
```

### Dynamic Configuration

```go
// Process configuration templates at runtime
func getProcessedConfig(template string) string {
    return env.ApplyEnvironmentToString(template)
}
```

## Error Handling

The package handles various scenarios gracefully:

- **Missing .env File**: Continues with system environment variables only
- **Invalid Environment Format**: Skips malformed environment variable entries
- **Missing Variables**: Leaves unmatched template patterns unchanged
- **Empty Values**: Properly handles empty environment variable values

## Best Practices

### Environment File Management

1. **Use .env for Development**: Keep development settings in `.env` files
2. **Don't Commit Secrets**: Add `.env` to `.gitignore` for security
3. **Provide .env.example**: Include example environment file in repository
4. **Document Variables**: Comment environment variables in example files

### Template Usage

1. **Consistent Patterns**: Use `${VARIABLE_NAME}` pattern consistently
2. **Descriptive Names**: Use clear, descriptive environment variable names
3. **Default Values**: Consider providing default values in application logic
4. **Validation**: Validate critical environment variables after loading

### Security Considerations

1. **Sensitive Data**: Never commit sensitive environment variables to version control
2. **Production Secrets**: Use secure secret management in production
3. **Access Control**: Limit access to environment files and variables
4. **Rotation**: Regularly rotate sensitive environment variables

## Dependencies

- **godotenv**: For loading `.env` files
- **strings**: For string manipulation and parsing
- **os**: For accessing system environment variables

## Related Documentation

- [Backend Utilities Overview](../README.md) - Main utilities documentation
- [Configuration Utilities](../configs/README.md) - Configuration management
- [Flags Utilities](../flags/README.md) - Command-line flag processing