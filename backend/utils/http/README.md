# HTTP Utilities

This package provides HTTP client utilities and request handling for the PostEaze backend application. It offers a configurable HTTP client system with support for request configuration, execution, and response handling patterns.

## Architecture

The HTTP utilities implement a flexible HTTP client system:

- **Configurable Clients**: Pre-configured HTTP clients with specific settings
- **Request Builder Pattern**: Fluent API for building HTTP requests
- **Response Handling**: Multiple response handling patterns (execute, get response, bind to struct)
- **Configuration Management**: Centralized configuration for different HTTP endpoints
- **Thread-Safe Operations**: Concurrent-safe request configuration storage

## Key Components

### Core Files

- **`http.go`**: Main HTTP client initialization and request building
- **`request.go`**: Request configuration and setup utilities
- **`execute.go`**: HTTP request execution and response handling

### Core Interfaces

- **`client`**: Main HTTP client with configuration storage
- **`Request`**: Request builder with fluent API
- **`RequestConfig`**: Configuration structure for HTTP requests

## HTTP Client Management

### Client Initialization

```go
import httpclient "github.com/shekhar8352/PostEaze/utils/http"

// Create request configurations
config1 := httpclient.NewRequestConfig("api-service", map[string]interface{}{
    "method": "GET",
    "url": "https://api.example.com/users",
    "timeoutinmillis": 5000,
    "headers": map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer ${API_TOKEN}",
    },
})

config2 := httpclient.NewRequestConfig("auth-service", map[string]interface{}{
    "method": "POST",
    "url": "https://auth.example.com/login",
    "timeoutinmillis": 3000,
})

// Initialize HTTP client with configurations
httpclient.InitHttp(config1, config2)
```

### Getting Client Instance

```go
// Get the HTTP client instance
client := httpclient.GetClient()
```

## Request Building

### Fluent Request API

```go
// Get pre-configured request
request := client.getRequestDetails("api-service")

// Customize request using fluent API
request.SetContext(ctx).
    SetMethod("POST").
    SetURL("https://api.example.com/users/123").
    SetQueryParam("include", "profile").
    SetQueryParams(map[string]string{
        "page": "1",
        "limit": "10",
    }).
    SetHeaderParam("X-Request-ID", "12345").
    SetHeaderParams(map[string]string{
        "X-Client-Version": "1.0.0",
        "X-Platform": "backend",
    }).
    SetBody(strings.NewReader(`{"name": "John Doe"}`))
```

### Request Configuration

```go
// Create request configuration from map
configMap := map[string]interface{}{
    "method": "POST",
    "url": "https://api.example.com/endpoint",
    "timeoutinmillis": 10000,
    "retrycount": 3,
    "headers": map[string]string{
        "Content-Type": "application/json",
        "User-Agent": "PostEaze-Backend/1.0",
    },
}

config := httpclient.NewRequestConfig("my-service", configMap)
```

## Request Execution

### Basic Execution

```go
import httpclient "github.com/shekhar8352/PostEaze/utils/http"

// Create HTTP request
httpRequest, err := client.createRequest(ctx, request)
if err != nil {
    return err
}

// Execute request (fire and forget)
err = httpclient.Call(httpRequest, httpClient)
if err != nil {
    return err
}
```

### Get Response Body

```go
// Execute request and get response body
responseBody, err := httpclient.CallAndGetResponse(httpRequest, httpClient)
if err != nil {
    return err
}

// Process response body
fmt.Println("Response:", string(responseBody))
```

### Bind Response to Struct

```go
type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var user UserResponse

// Execute request and bind response to struct
err = httpclient.CallAndBind(httpRequest, &user, httpClient)
if err != nil {
    return err
}

fmt.Printf("User: %+v\n", user)
```

## Configuration Options

### Request Configuration Parameters

- **`method`**: HTTP method (GET, POST, PUT, DELETE, etc.)
- **`url`**: Target URL for the request
- **`timeoutinmillis`**: Request timeout in milliseconds
- **`retrycount`**: Number of retry attempts (future enhancement)
- **`headers`**: Default headers for the request

### Timeout Configuration

```go
// Configure timeout in request config
configMap := map[string]interface{}{
    "timeoutinmillis": 5000,  // 5 second timeout
}
```

### Header Configuration

```go
// Configure default headers
configMap := map[string]interface{}{
    "headers": map[string]string{
        "Content-Type": "application/json",
        "Accept": "application/json",
        "User-Agent": "PostEaze-Backend",
    },
}
```

## Advanced Usage Patterns

### Service-Specific Clients

```go
// Configure different services with different settings
apiConfig := httpclient.NewRequestConfig("api", map[string]interface{}{
    "url": "https://api.postease.com",
    "timeoutinmillis": 10000,
    "headers": map[string]string{
        "Authorization": "Bearer ${API_TOKEN}",
    },
})

authConfig := httpclient.NewRequestConfig("auth", map[string]interface{}{
    "url": "https://auth.postease.com",
    "timeoutinmillis": 5000,
    "headers": map[string]string{
        "Content-Type": "application/json",
    },
})

httpclient.InitHttp(apiConfig, authConfig)
```

### Dynamic Request Building

```go
func makeAPICall(endpoint string, params map[string]string) ([]byte, error) {
    client := httpclient.GetClient()
    request := client.getRequestDetails("api")
    
    // Build request dynamically
    request.SetContext(context.Background()).
        SetURL(fmt.Sprintf("https://api.postease.com/%s", endpoint))
    
    // Add query parameters
    for key, value := range params {
        request.SetQueryParam(key, value)
    }
    
    // Create and execute request
    httpRequest, err := client.createRequest(context.Background(), request)
    if err != nil {
        return nil, err
    }
    
    return httpclient.CallAndGetResponse(httpRequest, http.Client{})
}
```

## Error Handling

### HTTP Status Code Handling

The execution functions automatically handle HTTP status codes:

- **Success (2xx)**: Request processed successfully
- **Client/Server Errors (non-2xx)**: Returns error with status information

```go
// Error handling example
err := httpclient.Call(request, client)
if err != nil {
    // Handle HTTP errors (non-2xx status codes)
    log.Printf("HTTP request failed: %v", err)
    return err
}
```

### Response Processing Errors

```go
// Handle JSON unmarshaling errors
var response MyStruct
err := httpclient.CallAndBind(request, &response, client)
if err != nil {
    // Could be HTTP error or JSON unmarshaling error
    log.Printf("Request or parsing failed: %v", err)
    return err
}
```

## Thread Safety

The HTTP client utilities are designed to be thread-safe:

- **Configuration Storage**: Uses `sync.RWMutex` for concurrent access to request configurations
- **Request Building**: Each request instance is independent
- **Client Reuse**: HTTP clients can be safely reused across goroutines

## Best Practices

### Configuration Management

1. **Pre-configure Services**: Set up request configurations for different services during initialization
2. **Use Environment Variables**: Configure URLs and tokens using environment variables
3. **Set Appropriate Timeouts**: Configure reasonable timeouts for different service types
4. **Default Headers**: Set common headers in configuration rather than per-request

### Request Building

1. **Reuse Configurations**: Use pre-configured request templates and customize as needed
2. **Context Usage**: Always provide context for request cancellation and timeouts
3. **Error Handling**: Handle both HTTP errors and parsing errors appropriately
4. **Resource Cleanup**: HTTP clients handle resource cleanup automatically

### Performance Optimization

1. **Client Reuse**: Reuse HTTP clients rather than creating new ones for each request
2. **Connection Pooling**: HTTP clients automatically manage connection pooling
3. **Timeout Configuration**: Set appropriate timeouts to avoid hanging requests
4. **Concurrent Requests**: The utilities support concurrent request execution

## Dependencies

- **net/http**: Standard Go HTTP client
- **context**: For request cancellation and timeouts
- **encoding/json**: For JSON response parsing
- **io**: For request/response body handling
- **sync**: For thread-safe operations
- **spf13/cast**: For configuration value type conversion

## Related Documentation

- [Backend Utilities Overview](../README.md) - Main utilities documentation
- [Configuration Utilities](../configs/README.md) - Configuration management
- [Environment Utilities](../env/README.md) - Environment variable handling