# HTTP Testing Utilities Guide

This guide explains how to use the HTTP testing utilities provided in the `testutils` package for testing Gin HTTP handlers and API endpoints.

## Overview

The HTTP testing utilities provide a comprehensive set of functions to:
- Create test HTTP contexts and requests
- Assert API responses and status codes
- Parse and validate JSON responses
- Test authenticated endpoints
- Set headers, parameters, and query strings
- Perform integration testing with routers

## Core Functions

### Creating Test Contexts and Requests

#### `NewTestGinContext(method, url string, body interface{}) (*gin.Context, *httptest.ResponseRecorder)`
Creates a new Gin context for testing with the specified HTTP method, URL, and request body.

```go
// GET request without body
ctx, recorder := NewTestGinContext("GET", "/api/users", nil)

// POST request with JSON body
requestBody := map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
}
ctx, recorder := NewTestGinContext("POST", "/api/users", requestBody)
```

#### `CreateTestRequest(method, url string, body interface{}) *http.Request`
Creates an HTTP request for testing purposes.

```go
req := CreateTestRequest("POST", "/api/login", loginData)
```

### Response Parsing and Validation

#### `ParseJSONResponse(recorder *httptest.ResponseRecorder, target interface{}) error`
Parses the JSON response from a ResponseRecorder into the target interface.

```go
var response APIResponse
err := ParseJSONResponse(recorder, &response)
```

#### `GetResponseJSON(recorder *httptest.ResponseRecorder) (map[string]interface{}, error)`
Parses the response body as JSON and returns it as a map.

```go
responseData, err := GetResponseJSON(recorder)
if err == nil {
    userID := responseData["data"].(map[string]interface{})["id"]
}
```

### Response Assertions

#### `AssertSuccessResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedData interface{})`
Asserts that the response is a successful API response with expected data.

```go
expectedData := map[string]interface{}{
    "id": "123",
    "name": "John Doe",
}
AssertSuccessResponse(t, recorder, expectedData)
```

#### `AssertErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedMessage string)`
Asserts that the response is an error response with expected status and message.

```go
AssertErrorResponse(t, recorder, http.StatusBadRequest, "invalid input")
```

#### `AssertStatusCode(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int)`
Asserts that the response has the expected status code.

```go
AssertStatusCode(t, recorder, http.StatusCreated)
```

### Setting Request Parameters

#### `SetAuthorizationHeader(ctx *gin.Context, token string)`
Sets the Authorization header with Bearer token.

```go
SetAuthorizationHeader(ctx, "your-jwt-token")
```

#### `SetURLParam(ctx *gin.Context, key, value string)`
Sets a URL parameter in the Gin context.

```go
SetURLParam(ctx, "id", "123")
SetURLParam(ctx, "type", "user")
```

#### `SetQueryParam(ctx *gin.Context, key, value string)`
Sets a query parameter in the Gin context.

```go
SetQueryParam(ctx, "page", "1")
SetQueryParam(ctx, "limit", "10")
```

#### `SetRequestHeader(ctx *gin.Context, key, value string)`
Sets a custom header on the request.

```go
SetRequestHeader(ctx, "X-Client-Version", "1.0.0")
```

### Router Testing

#### `CreateTestRouter() *gin.Engine`
Creates a new Gin router for testing.

```go
router := CreateTestRouter()
router.POST("/api/users", UserHandler)
```

#### `PerformRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder`
Performs an HTTP request on the given router and returns the response recorder.

```go
recorder := PerformRequest(router, "POST", "/api/users", userData)
```

## Usage Examples

### Testing a Simple API Handler

```go
func TestUserHandler(t *testing.T) {
    t.Run("successful user creation", func(t *testing.T) {
        // Arrange
        requestBody := map[string]interface{}{
            "name":  "John Doe",
            "email": "john@example.com",
        }
        
        ctx, recorder := NewTestGinContext("POST", "/api/users", requestBody)
        
        // Act
        UserHandler(ctx)
        
        // Assert
        AssertStatusCode(t, recorder, http.StatusCreated)
        AssertJSONResponse(t, recorder)
        
        expectedData := map[string]interface{}{
            "id":   "123",
            "name": "John Doe",
        }
        AssertSuccessResponse(t, recorder, expectedData)
    })
    
    t.Run("validation error", func(t *testing.T) {
        // Arrange
        requestBody := map[string]interface{}{
            "email": "john@example.com", // missing name
        }
        
        ctx, recorder := NewTestGinContext("POST", "/api/users", requestBody)
        
        // Act
        UserHandler(ctx)
        
        // Assert
        AssertErrorResponse(t, recorder, http.StatusBadRequest, "name is required")
    })
}
```

### Testing Authenticated Endpoints

```go
func TestAuthenticatedHandler(t *testing.T) {
    t.Run("successful authenticated request", func(t *testing.T) {
        // Arrange
        ctx, recorder := NewTestGinContext("GET", "/api/profile", nil)
        SetAuthorizationHeader(ctx, "valid-jwt-token")
        
        // Act
        ProfileHandler(ctx)
        
        // Assert
        AssertStatusCode(t, recorder, http.StatusOK)
        AssertSuccessResponse(t, recorder, nil)
    })
    
    t.Run("missing authorization", func(t *testing.T) {
        // Arrange
        ctx, recorder := NewTestGinContext("GET", "/api/profile", nil)
        
        // Act
        ProfileHandler(ctx)
        
        // Assert
        AssertErrorResponse(t, recorder, http.StatusUnauthorized, "authorization required")
    })
}
```

### Testing with URL Parameters

```go
func TestGetUserByID(t *testing.T) {
    t.Run("user found", func(t *testing.T) {
        // Arrange
        ctx, recorder := NewTestGinContext("GET", "/api/users/123", nil)
        SetURLParam(ctx, "id", "123")
        
        // Act
        GetUserByIDHandler(ctx)
        
        // Assert
        AssertStatusCode(t, recorder, http.StatusOK)
        AssertResponseContains(t, recorder, "123")
    })
}
```

### Testing with Query Parameters

```go
func TestGetUsers(t *testing.T) {
    t.Run("with pagination", func(t *testing.T) {
        // Arrange
        ctx, recorder := NewTestGinContext("GET", "/api/users", nil)
        SetQueryParam(ctx, "page", "2")
        SetQueryParam(ctx, "limit", "10")
        
        // Act
        GetUsersHandler(ctx)
        
        // Assert
        AssertStatusCode(t, recorder, http.StatusOK)
        
        // Verify pagination in response
        responseData, err := GetResponseJSON(recorder)
        assert.NoError(t, err)
        
        data := responseData["data"].(map[string]interface{})
        assert.Equal(t, "2", data["page"])
        assert.Equal(t, "10", data["limit"])
    })
}
```

### Integration Testing with Router

```go
func TestUserAPIIntegration(t *testing.T) {
    // Setup
    router := CreateTestRouter()
    router.POST("/api/users", CreateUserHandler)
    router.GET("/api/users/:id", GetUserHandler)
    
    t.Run("create and retrieve user", func(t *testing.T) {
        // Create user
        userData := map[string]interface{}{
            "name":  "Jane Doe",
            "email": "jane@example.com",
        }
        
        recorder := PerformRequest(router, "POST", "/api/users", userData)
        AssertStatusCode(t, recorder, http.StatusCreated)
        
        // Extract user ID from response
        responseData, err := GetResponseJSON(recorder)
        require.NoError(t, err)
        userID := responseData["data"].(map[string]interface{})["id"].(string)
        
        // Retrieve user
        recorder = PerformRequest(router, "GET", "/api/users/"+userID, nil)
        AssertStatusCode(t, recorder, http.StatusOK)
        AssertResponseContains(t, recorder, "Jane Doe")
    })
}
```

### Testing Custom Headers

```go
func TestCustomHeaders(t *testing.T) {
    t.Run("with client version header", func(t *testing.T) {
        // Arrange
        ctx, recorder := NewTestGinContext("GET", "/api/version", nil)
        SetRequestHeader(ctx, "X-Client-Version", "2.1.0")
        
        // Act
        VersionHandler(ctx)
        
        // Assert
        AssertStatusCode(t, recorder, http.StatusOK)
        AssertResponseHeaders(t, recorder, map[string]string{
            "X-API-Version": "1.0",
        })
    })
}
```

## Best Practices

1. **Use descriptive test names**: Use clear, descriptive names for test functions and sub-tests.

2. **Follow AAA pattern**: Structure tests with Arrange, Act, Assert sections.

3. **Test both success and error cases**: Always test both happy path and error scenarios.

4. **Use appropriate assertions**: Choose the most specific assertion function for your needs.

5. **Clean up resources**: Use proper setup and teardown for database connections and other resources.

6. **Test edge cases**: Include tests for boundary conditions, empty inputs, and invalid data.

7. **Use table-driven tests**: For testing multiple similar scenarios, use table-driven tests.

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name           string
        input          map[string]interface{}
        expectedStatus int
        expectedError  string
    }{
        {
            name:           "valid input",
            input:          map[string]interface{}{"name": "John", "email": "john@example.com"},
            expectedStatus: http.StatusOK,
        },
        {
            name:           "missing name",
            input:          map[string]interface{}{"email": "john@example.com"},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "name is required",
        },
        {
            name:           "invalid email",
            input:          map[string]interface{}{"name": "John", "email": "invalid"},
            expectedStatus: http.StatusBadRequest,
            expectedError:  "invalid email",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx, recorder := NewTestGinContext("POST", "/api/users", tt.input)
            
            UserHandler(ctx)
            
            AssertStatusCode(t, recorder, tt.expectedStatus)
            if tt.expectedError != "" {
                AssertErrorResponse(t, recorder, tt.expectedStatus, tt.expectedError)
            }
        })
    }
}
```

## Available Assertion Functions

- `AssertSuccessResponse()` - Assert successful API response
- `AssertErrorResponse()` - Assert error API response
- `AssertStatusCode()` - Assert HTTP status code
- `AssertJSONResponse()` - Assert valid JSON response
- `AssertContentType()` - Assert content type header
- `AssertResponseHeaders()` - Assert custom headers
- `AssertResponseContains()` - Assert response contains text
- `AssertResponseNotContains()` - Assert response doesn't contain text
- `AssertEmptyResponse()` - Assert empty response body
- `AssertNonEmptyResponse()` - Assert non-empty response body

This comprehensive set of utilities makes it easy to write thorough tests for your API endpoints while maintaining consistency and readability across your test suite.