# Flags Utilities

This package provides command-line flag handling and application startup parameter management for the PostEaze backend application. It defines and parses command-line flags that control application behavior, configuration paths, and runtime settings.

## Architecture

The flags utilities implement a centralized command-line argument processing system:

- **Flag Definition**: Pre-defined flags for common application parameters
- **Automatic Parsing**: Flags are parsed automatically during package initialization
- **Default Values**: Sensible default values for all flags
- **Environment Integration**: Combines flag-based and environment-based configuration
- **Constants Integration**: Uses application constants for consistency

## Key Components

### Core Files

- **`flags.go`**: Main flag definitions, parsing, and accessor functions

### Supported Flags

The package defines several key application flags:

- **`mode`**: Application running mode (development, production, etc.)
- **`port`**: Server port number
- **`baseConfigPath`**: Base path for configuration files

## Flag Definitions

### Application Mode Flag

```go
var mode = flag.String(constants.ModeKey, constants.DefaultMode, constants.ModeUsage)
```

- **Flag Name**: Defined by `constants.ModeKey`
- **Default Value**: `constants.DefaultMode`
- **Usage Description**: `constants.ModeUsage`

### Port Flag

```go
var port = flag.Int(constants.PortKey, constants.DefaultPort, constants.PortUsage)
```

- **Flag Name**: Defined by `constants.PortKey`
- **Default Value**: `constants.DefaultPort`
- **Usage Description**: `constants.PortUsage`

### Base Configuration Path Flag

```go
var baseConfigPath = flag.String(constants.BaseConfigPathKey, constants.DefaultBaseConfigPath, constants.BaseConfigPathUsage)
```

- **Flag Name**: Defined by `constants.BaseConfigPathKey`
- **Default Value**: `constants.DefaultBaseConfigPath`
- **Usage Description**: `constants.BaseConfigPathUsage`

## Usage Examples

### Command-Line Usage

```bash
# Run with default values
./postease

# Specify custom port
./postease -port=9090

# Specify custom mode
./postease -mode=production

# Specify custom config path
./postease -baseConfigPath=/etc/postease/configs

# Combine multiple flags
./postease -mode=production -port=8080 -baseConfigPath=/opt/configs
```

### Programmatic Access

```go
import "github.com/shekhar8352/PostEaze/utils/flags"

func main() {
    // Get application mode
    appMode := flags.Mode()
    fmt.Printf("Running in %s mode\n", appMode)
    
    // Get server port
    serverPort := flags.Port()
    fmt.Printf("Server will run on port %d\n", serverPort)
    
    // Get base configuration path
    configPath := flags.BaseConfigPath()
    fmt.Printf("Loading configs from %s\n", configPath)
    
    // Get environment variables
    env := flags.Env()
    region := flags.AWSRegion()
    fmt.Printf("Environment: %s, AWS Region: %s\n", env, region)
}
```

### Application Initialization

```go
func initializeApplication() {
    // Flags are automatically parsed during package init
    
    // Use flag values for application setup
    mode := flags.Mode()
    port := flags.Port()
    configPath := flags.BaseConfigPath()
    
    // Configure application based on flags
    if mode == "development" {
        setupDevelopmentMode(configPath)
    } else {
        setupProductionMode(flags.Env(), flags.AWSRegion())
    }
    
    // Start server on specified port
    startServer(port)
}
```

## Flag Accessor Functions

### Mode()

Returns the application running mode:

```go
func Mode() string
```

**Usage:**
```go
mode := flags.Mode()
switch mode {
case "development":
    // Development-specific setup
case "production":
    // Production-specific setup
default:
    // Default setup
}
```

### Port()

Returns the server port number:

```go
func Port() int
```

**Usage:**
```go
port := flags.Port()
server := &http.Server{
    Addr: fmt.Sprintf(":%d", port),
    Handler: router,
}
```

### BaseConfigPath()

Returns the base configuration path for file-based configurations:

```go
func BaseConfigPath() string
```

**Usage:**
```go
configPath := flags.BaseConfigPath()
err := configs.InitDev(configPath, "app", "database")
```

### Env()

Returns the environment from environment variables:

```go
func Env() string
```

**Usage:**
```go
env := flags.Env()
if env == "production" {
    // Production-specific logic
}
```

### AWSRegion()

Returns the AWS region from environment variables:

```go
func AWSRegion() string
```

**Usage:**
```go
region := flags.AWSRegion()
err := configs.InitRelease(flags.Env(), region, configNames...)
```

## Integration Patterns

### Configuration Integration

```go
func setupConfiguration() error {
    mode := flags.Mode()
    
    if mode == "development" {
        // Use file-based configuration
        return configs.InitDev(flags.BaseConfigPath(), "app", "database")
    } else {
        // Use AWS App Config
        return configs.InitRelease(flags.Env(), flags.AWSRegion(), "app", "database")
    }
}
```

### Server Setup

```go
func startServer() {
    port := flags.Port()
    mode := flags.Mode()
    
    // Configure server based on mode
    if mode == "development" {
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // Start server
    router := setupRoutes()
    log.Printf("Starting server on port %d in %s mode", port, mode)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
```

### Environment-Aware Logic

```go
func setupLogging() {
    env := flags.Env()
    mode := flags.Mode()
    
    if mode == "development" || env == "development" {
        // Development logging
        log.SetLevel(log.DebugLevel)
    } else {
        // Production logging
        log.SetLevel(log.InfoLevel)
    }
}
```

## Flag Parsing

### Automatic Parsing

Flags are automatically parsed during package initialization:

```go
func init() {
    flag.Parse()
}
```

This ensures that all flags are parsed before any accessor functions are called.

### Custom Flag Parsing

If you need to add custom flags, follow the same pattern:

```go
var customFlag = flag.String("custom", "default", "Custom flag description")

func CustomFlag() string {
    return *customFlag
}
```

## Environment Variable Integration

The package integrates with environment variables for configuration:

- **`constants.EnvKey`**: Environment name (development, staging, production)
- **`constants.AWSRegionKey`**: AWS region for cloud services

```go
// Environment variables are accessed directly
env := os.Getenv(constants.EnvKey)
region := os.Getenv(constants.AWSRegionKey)
```

## Best Practices

### Flag Usage

1. **Consistent Naming**: Use constants for flag names to ensure consistency
2. **Sensible Defaults**: Provide reasonable default values for all flags
3. **Clear Descriptions**: Write clear usage descriptions for each flag
4. **Environment Fallback**: Use environment variables as fallback for sensitive values

### Application Setup

1. **Early Initialization**: Access flags early in application startup
2. **Validation**: Validate flag values before using them
3. **Documentation**: Document available flags in application help
4. **Testing**: Test application behavior with different flag combinations

### Configuration Strategy

1. **Mode-Based Logic**: Use mode flag to determine configuration strategy
2. **Path Configuration**: Use baseConfigPath for flexible configuration file location
3. **Environment Awareness**: Combine flags with environment variables for complete configuration
4. **Default Behavior**: Ensure application works with default flag values

## Error Handling

The flag package handles parsing errors automatically, but you should validate flag values:

```go
func validateFlags() error {
    port := flags.Port()
    if port < 1 || port > 65535 {
        return fmt.Errorf("invalid port number: %d", port)
    }
    
    mode := flags.Mode()
    validModes := []string{"development", "staging", "production"}
    if !contains(validModes, mode) {
        return fmt.Errorf("invalid mode: %s", mode)
    }
    
    return nil
}
```

## Dependencies

- **flag**: Standard Go flag parsing package
- **os**: For environment variable access
- **Constants Package**: Application constants for flag definitions

## Related Documentation

- [Backend Utilities Overview](../README.md) - Main utilities documentation
- [Configuration Utilities](../configs/README.md) - Configuration management
- [Environment Utilities](../env/README.md) - Environment variable handling
- [Constants](../../constants/README.md) - Application constants