package helpers

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewHTTPTest(t *testing.T) {
	helper := NewHTTPTest()
	
	if helper == nil {
		t.Fatal("NewHTTPTest() returned nil")
	}
	
	if helper.Router == nil {
		t.Fatal("HTTPHelper.Router is nil")
	}
}

func TestHTTPHelper_Request(t *testing.T) {
	helper := NewHTTPTest()
	
	// Test simple GET request to health endpoint
	resp := helper.Request("GET", "/health", nil)
	
	if resp == nil {
		t.Fatal("Request() returned nil response")
	}
	
	// Should get 200 response
	if resp.Code != 200 {
		t.Errorf("Expected status code 200, got %d", resp.Code)
	}
}

func TestHTTPHelper_AuthRequest(t *testing.T) {
	helper := NewHTTPTest()
	
	// Test authenticated request
	resp := helper.AuthRequest("POST", "/api/v1/auth/refresh", "test-user-123", nil)
	
	if resp == nil {
		t.Fatal("AuthRequest() returned nil response")
	}
	
	// Should get some response
	if resp.Code == 0 {
		t.Error("Expected non-zero status code")
	}
}

func TestHTTPHelper_AuthRequestWithRole(t *testing.T) {
	helper := NewHTTPTest()
	
	// Test authenticated request with role
	resp := helper.AuthRequestWithRole("POST", "/api/v1/auth/logout", "test-user-123", "admin", nil)
	
	if resp == nil {
		t.Fatal("AuthRequestWithRole() returned nil response")
	}
	
	// Should get some response
	if resp.Code == 0 {
		t.Error("Expected non-zero status code")
	}
}

func TestHTTPHelper_RequestWithHeaders(t *testing.T) {
	helper := NewHTTPTest()
	
	headers := map[string]string{
		"X-Test-Header": "test-value",
		"Content-Type":  "application/json",
	}
	
	resp := helper.RequestWithHeaders("GET", "/health", headers, nil)
	
	if resp == nil {
		t.Fatal("RequestWithHeaders() returned nil response")
	}
	
	// Should get 200 response
	if resp.Code != 200 {
		t.Errorf("Expected status code 200, got %d", resp.Code)
	}
}

func TestHTTPHelper_WithCustomRouter(t *testing.T) {
	// Create custom router
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test endpoint"})
	})
	
	helper := NewHTTPTestWithRouter(router)
	resp := helper.Request("GET", "/test", nil)
	
	if resp.Code != 200 {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
	
	if err := AssertContains(resp, "test endpoint"); err != nil {
		t.Errorf("Response assertion failed: %v", err)
	}
}

func TestAssertStatus(t *testing.T) {
	helper := NewHTTPTest()
	resp := helper.Request("GET", "/api/v1/auth/signup", nil)
	
	// Test successful assertion
	if err := AssertStatus(resp, resp.Code); err != nil {
		t.Errorf("AssertStatus failed: %v", err)
	}
	
	// Test failed assertion
	if err := AssertStatus(resp, 999); err == nil {
		t.Error("AssertStatus should have failed for wrong status code")
	}
}

func TestAssertJSON(t *testing.T) {
	helper := NewHTTPTest()
	resp := helper.Request("GET", "/health", nil)
	
	// The health response should be JSON
	if err := AssertJSON(resp); err != nil {
		t.Errorf("AssertJSON failed: %v", err)
	}
}

func TestAssertContains(t *testing.T) {
	helper := NewHTTPTest()
	resp := helper.Request("GET", "/health", nil)
	
	body := GetResponseBody(resp)
	if body == "" {
		t.Skip("Empty response body, skipping contains test")
	}
	
	// Test with content that should exist
	if err := AssertContains(resp, "success"); err != nil {
		t.Errorf("AssertContains failed: %v", err)
	}
	
	// Test with content that shouldn't exist
	if err := AssertContains(resp, "nonexistent-content-xyz"); err == nil {
		t.Error("AssertContains should have failed for non-existent content")
	}
}

func TestNewLoginRequest(t *testing.T) {
	// Test default values
	req := NewLoginRequest()
	if req.Email == "" || req.Password == "" {
		t.Error("NewLoginRequest should have default values")
	}
	
	// Test with overrides
	req = NewLoginRequest(func(r *LoginRequest) {
		r.Email = "custom@example.com"
	})
	if req.Email != "custom@example.com" {
		t.Error("Override function should modify email")
	}
}

func TestNewSignupRequest(t *testing.T) {
	// Test default values
	req := NewSignupRequest()
	if req.Name == "" || req.Email == "" || req.Password == "" || req.UserType == "" {
		t.Error("NewSignupRequest should have default values")
	}
	
	// Test with overrides
	req = NewSignupRequest(func(r *SignupRequest) {
		r.UserType = "team"
	})
	if req.UserType != "team" {
		t.Error("Override function should modify user type")
	}
}

func TestNewRefreshTokenRequest(t *testing.T) {
	token := "test-refresh-token"
	req := NewRefreshTokenRequest(token)
	
	if req.RefreshToken != token {
		t.Error("NewRefreshTokenRequest should set the token")
	}
	
	// Test with overrides
	req = NewRefreshTokenRequest(token, func(r *RefreshTokenRequest) {
		r.RefreshToken = "modified-token"
	})
	if req.RefreshToken != "modified-token" {
		t.Error("Override function should modify token")
	}
}

func TestGetResponseJSON(t *testing.T) {
	// Create a custom router with JSON response
	router := gin.New()
	router.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{"test": "value", "number": 42})
	})
	
	helper := NewHTTPTestWithRouter(router)
	resp := helper.Request("GET", "/json", nil)
	
	var result map[string]interface{}
	if err := GetResponseJSON(resp, &result); err != nil {
		t.Errorf("GetResponseJSON failed: %v", err)
	}
	
	if result["test"] != "value" {
		t.Error("JSON parsing failed to get correct value")
	}
	
	if result["number"] != float64(42) { // JSON numbers are float64
		t.Error("JSON parsing failed to get correct number")
	}
}

func TestGetResponseBody(t *testing.T) {
	// Create a custom router with text response
	router := gin.New()
	router.GET("/text", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
	
	helper := NewHTTPTestWithRouter(router)
	resp := helper.Request("GET", "/text", nil)
	
	body := GetResponseBody(resp)
	if body != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", body)
	}
}