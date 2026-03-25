# Constants

Application constants and configuration keys for the PostEaze backend. This package centralizes all constant values used throughout the application, including configuration keys, application settings, routing paths, and command-line flags.

## Contents

This package contains four main constant files:

- **constants.go**: Core application constants like application name and common values
- **configs.go**: Configuration keys and names used for loading application settings
- **flags.go**: Command-line flag definitions and default values for application startup
- **router.go**: API route path constants used throughout the routing system

## Constant Categories

### Application Constants (`constants.go`)
Basic application-wide constants:
```go
ApplicationName = "post-eaze"  // Used for application identification
Empty = ""                     // Common empty string constant
```

### Configuration Constants (`configs.go`)
Constants for configuration management:

**Config Initialization Keys:**
- `ConfigDirectoryKey`, `ConfigNamesKey`, `ConfigTypeKey` - Configuration loading parameters
- `ConfigIDKey`, `ConfigRegionKey`, `ConfigEnvKey` - Environment and deployment settings

**Config Names:**
- `ApplicationConfig`, `LoggerConfig`, `DatabaseConfig`, `APIConfig` - Configuration file names

**Config Property Keys:**
- Database: `DatabaseDriverNameConfigKey`, `DatabaseURLConfigKey`, connection pool settings
- Logger: `LoggerLevelConfigKey`, `LoggerParamsConfigKey`
- API: `APIGetCatsFactConfigKey` - External API configuration

### Command-Line Flags (`flags.go`)
Application startup flags and their defaults:
```go
ModeKey = "mode"                              // Application run mode (dev/release)
PortKey = "port"                             // Server port (default: 8080)
BaseConfigPathKey = "base-config-path"       // Config directory path
EnvKey = "ENV"                               // Environment variable key
AWSRegionKey = "AWS_REGION"                  // AWS region setting
```

### API Routes (`router.go`)
Route path constants for consistent API endpoint definitions (non-exhaustive):
```go
ApiRoute      = "/api"
V1Route       = "/v1"
AuthRoute     = "/auth"
Authenticate  = "/authenticate" // POST — Firebase sign-in
RefreshRoute  = "/refresh"
LogOutRoute     = "/logout"
LogRoute      = "/log"
UserRoute     = "/user"
TeamRoute     = "/team"
MetaRoute     = "/meta"
ChannelRoute  = "/channels"
// ... see router.go for log paths, webhooks, etc.
```

## Usage Patterns

### Configuration Loading
Constants from `configs.go` are used with the configuration utilities to load settings:
```go
// Example usage in configuration loading
configName := constants.DatabaseConfig
driverKey := constants.DatabaseDriverNameConfigKey
```

### Route definition
Router constants group paths consistently, e.g. `ApiRoute` + `V1Route` + `AuthRoute` for `/api/v1/auth`, then `Authenticate` for `POST /authenticate`.

### Flag Processing
Flag constants provide consistent command-line interface:
```go
// Example usage in flag parsing
mode := flag.String(constants.ModeKey, constants.DefaultMode, constants.ModeUsage)
port := flag.Int(constants.PortKey, constants.DefaultPort, constants.PortUsage)
```

## Best Practices

1. **Centralization**: All constants are defined here to avoid duplication across the codebase
2. **Naming Convention**: Constants use descriptive names with appropriate suffixes (Key, Route, Config)
3. **Grouping**: Related constants are grouped together with comments for clarity
4. **Immutability**: All values are constants to prevent accidental modification
5. **Documentation**: Each constant group includes usage examples and explanations

## Related Documentation

- [Configuration Utilities](../utils/configs/README.md) - How constants are used in config loading
- [API Documentation](../api/README.md) - How route constants are used in API setup
- [Flags Utilities](../utils/flags/README.md) - How flag constants are processed
- [Resources Configs](../resources/README.md) - Configuration files that use these constants