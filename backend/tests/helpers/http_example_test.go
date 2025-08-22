package helpers

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// Example test showing how to use the HTTP helper for API testing
func TestHTTPHelper_APIExample(t *testing.T) {
	// Create a custom router with actual API endpoints
	router := gin.New()
	
	// Mock signup endpoint
	router.POST("/api/v1/auth/signup", func(c *gin.Context) {
		var req SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"status": "error", "message": "Invalid request"})
			return
		}
		
		if req.Email == "" {
			c.JSON(400, gin.H{"status": "error", "message": "Email is required"})
			return
		}
		
		c.JSON(201, gin.H{
			"status": "success", 
			"message": "User created successfully",
			"data": gin.H{"id": "user-123", "email": req.Email},
		})
	})
	
	// Mock login endpoint
	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"status": "error", "message": "Invalid request"})
			return
		}
		
		if req.Email == "test@example.com" && req.Password == "testpassword123" {
			c.JSON(200, gin.H{
				"status": "success",
				"message": "Login successful",
				"data": gin.H{
					"access_token": "mock-access-token",
					"refresh_token": "mock-refresh-token",
				},
			})
		} else {
			c.JSON(401, gin.H{"status": "error", "message": "Invalid credentials"})
		}
	})
	
	// Mock protected endpoint
	router.POST("/api/v1/auth/logout", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"status": "error", "message": "Authorization header required"})
			return
		}
		
		c.JSON(200, gin.H{"status": "success", "message": "Logged out successfully"})
	})
	
	helper := NewHTTPTestWithRouter(router)
	
	t.Run("signup with valid data", func(t *testing.T) {
		signupReq := NewSignupRequest(func(r *SignupRequest) {
			r.Email = "newuser@example.com"
			r.Name = "New User"
		})
		
		resp := helper.Request("POST", "/api/v1/auth/signup", signupReq)
		
		// Assert status code
		if err := AssertStatus(resp, 201); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
		
		// Assert JSON response
		if err := AssertJSON(resp); err != nil {
			t.Errorf("JSON assertion failed: %v", err)
		}
		
		// Assert response contains success
		if err := AssertContains(resp, "success"); err != nil {
			t.Errorf("Contains assertion failed: %v", err)
		}
		
		// Parse and verify response data
		var result map[string]interface{}
		if err := GetResponseJSON(resp, &result); err != nil {
			t.Errorf("Failed to parse JSON: %v", err)
		}
		
		if result["status"] != "success" {
			t.Errorf("Expected status 'success', got %v", result["status"])
		}
	})
	
	t.Run("signup with missing email", func(t *testing.T) {
		signupReq := NewSignupRequest(func(r *SignupRequest) {
			r.Email = "" // Missing email
		})
		
		resp := helper.Request("POST", "/api/v1/auth/signup", signupReq)
		
		// Should return 400 Bad Request
		if err := AssertStatus(resp, 400); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
		
		if err := AssertContains(resp, "Email is required"); err != nil {
			t.Errorf("Error message assertion failed: %v", err)
		}
	})
	
	t.Run("login with valid credentials", func(t *testing.T) {
		loginReq := NewLoginRequest() // Uses default test credentials
		
		resp := helper.Request("POST", "/api/v1/auth/login", loginReq)
		
		if err := AssertStatus(resp, 200); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
		
		if err := AssertContains(resp, "access_token"); err != nil {
			t.Errorf("Token assertion failed: %v", err)
		}
	})
	
	t.Run("login with invalid credentials", func(t *testing.T) {
		loginReq := NewLoginRequest(func(r *LoginRequest) {
			r.Password = "wrongpassword"
		})
		
		resp := helper.Request("POST", "/api/v1/auth/login", loginReq)
		
		if err := AssertStatus(resp, 401); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
		
		if err := AssertContains(resp, "Invalid credentials"); err != nil {
			t.Errorf("Error message assertion failed: %v", err)
		}
	})
	
	t.Run("logout with authentication", func(t *testing.T) {
		resp := helper.AuthRequest("POST", "/api/v1/auth/logout", "user-123", nil)
		
		if err := AssertStatus(resp, 200); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
		
		if err := AssertContains(resp, "Logged out successfully"); err != nil {
			t.Errorf("Success message assertion failed: %v", err)
		}
	})
	
	t.Run("logout without authentication", func(t *testing.T) {
		resp := helper.Request("POST", "/api/v1/auth/logout", nil)
		
		if err := AssertStatus(resp, 401); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
		
		if err := AssertContains(resp, "Authorization header required"); err != nil {
			t.Errorf("Error message assertion failed: %v", err)
		}
	})
	
	t.Run("custom headers example", func(t *testing.T) {
		headers := map[string]string{
			"X-Request-ID": "test-123",
			"User-Agent":   "test-client/1.0",
		}
		
		// Use an existing endpoint from our custom router
		loginReq := NewLoginRequest()
		resp := helper.RequestWithHeaders("POST", "/api/v1/auth/login", headers, loginReq)
		
		if err := AssertStatus(resp, 200); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
	})
}

// Example showing how to test with different user roles
func TestHTTPHelper_RoleBasedExample(t *testing.T) {
	router := gin.New()
	
	// Mock admin-only endpoint
	router.GET("/api/v1/admin/users", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"status": "error", "message": "Unauthorized"})
			return
		}
		
		// In a real app, you'd parse the JWT and check the role
		// For this example, we'll assume the token is valid
		c.JSON(200, gin.H{
			"status": "success",
			"data": []gin.H{
				{"id": "user-1", "name": "User 1"},
				{"id": "user-2", "name": "User 2"},
			},
		})
	})
	
	helper := NewHTTPTestWithRouter(router)
	
	t.Run("admin endpoint with admin role", func(t *testing.T) {
		resp := helper.AuthRequestWithRole("GET", "/api/v1/admin/users", "admin-123", "admin", nil)
		
		if err := AssertStatus(resp, 200); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
		
		var result map[string]interface{}
		if err := GetResponseJSON(resp, &result); err != nil {
			t.Errorf("Failed to parse JSON: %v", err)
		}
		
		data, ok := result["data"].([]interface{})
		if !ok {
			t.Error("Expected data to be an array")
		}
		
		if len(data) != 2 {
			t.Errorf("Expected 2 users, got %d", len(data))
		}
	})
	
	t.Run("admin endpoint without auth", func(t *testing.T) {
		resp := helper.Request("GET", "/api/v1/admin/users", nil)
		
		if err := AssertStatus(resp, 401); err != nil {
			t.Errorf("Status assertion failed: %v", err)
		}
	})
}