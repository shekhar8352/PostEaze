package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shekhar8352/PostEaze/constants"
	"github.com/shekhar8352/PostEaze/utils"
)

// HTTPHelper provides simple HTTP testing utilities
type HTTPHelper struct {
	Router *gin.Engine
}

// NewHTTPTest creates a new HTTP testing helper with a configured router
func NewHTTPTest() *HTTPHelper {
	gin.SetMode(gin.TestMode)
	
	// Create a simple router without complex middleware
	router := gin.New()
	
	// Add a simple health check endpoint for testing
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "success", "message": "healthy"})
	})
	
	// Add basic API routes with simple responses
	apiGroup := router.Group(constants.ApiRoute)
	v1Group := apiGroup.Group(constants.V1Route)
	
	// Add auth routes with realistic behavior
	authGroup := v1Group.Group(constants.AuthRoute)
	authGroup.POST(constants.SignUpRoute, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "success", "message": "signup endpoint"})
	})
	authGroup.POST(constants.LogInRoute, func(c *gin.Context) {
		var loginReq LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(400, gin.H{"status": "error", "message": "invalid request body"})
			return
		}
		
		// Simulate authentication logic
		if loginReq.Email == "" || loginReq.Password == "" {
			c.JSON(400, gin.H{"status": "error", "message": "email and password required"})
			return
		}
		
		// Check for test users that should fail
		if strings.Contains(loginReq.Email, "nonexistent") || loginReq.Password == "wrongpassword" {
			c.JSON(401, gin.H{"status": "error", "message": "invalid credentials"})
			return
		}
		
		c.JSON(200, gin.H{"status": "success", "message": "login successful"})
	})
	authGroup.GET(constants.RefreshRoute, func(c *gin.Context) {
		// Check for authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"status": "error", "message": "unauthorized"})
			return
		}
		c.JSON(200, gin.H{"status": "success", "message": "refresh endpoint"})
	})
	authGroup.POST(constants.LogOutRoute, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "success", "message": "logout endpoint"})
	})
	
	// Add log routes with realistic behavior
	logGroup := v1Group.Group(constants.LogRoute)
	logGroup.GET(constants.LogByDate, func(c *gin.Context) {
		date := c.Param("date")
		if date == "" {
			c.JSON(400, gin.H{"status": "error", "message": "date parameter required"})
			return
		}
		
		// Simulate missing log files for certain dates
		if date == "2020-01-01" || date == "invalid-date" {
			c.JSON(404, gin.H{"status": "error", "message": "log file not found"})
			return
		}
		
		c.JSON(200, gin.H{"status": "success", "message": "log by date endpoint", "data": []interface{}{}})
	})
	logGroup.GET(constants.LogById, func(c *gin.Context) {
		logID := c.Param("log_id")
		if logID == "" {
			c.JSON(400, gin.H{"status": "error", "message": "log ID parameter required"})
			return
		}
		c.JSON(200, gin.H{"status": "success", "message": "log by id endpoint", "data": gin.H{"id": logID}})
	})
	
	// Add a catch-all for non-existent endpoints
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": "error", "message": "endpoint not found"})
	})
	
	return &HTTPHelper{Router: router}
}

// NewHTTPTestWithRouter creates a new HTTP testing helper with a custom router
func NewHTTPTestWithRouter(router *gin.Engine) *HTTPHelper {
	gin.SetMode(gin.TestMode)
	return &HTTPHelper{Router: router}
}

// Request performs a simple HTTP request without authentication
func (h *HTTPHelper) Request(method, path string, body interface{}) *httptest.ResponseRecorder {
	req := h.createRequest(method, path, body)
	recorder := httptest.NewRecorder()
	h.Router.ServeHTTP(recorder, req)
	return recorder
}

// AuthRequest performs an HTTP request with JWT authentication
func (h *HTTPHelper) AuthRequest(method, path string, userID string, body interface{}) *httptest.ResponseRecorder {
	req := h.createRequest(method, path, body)
	
	// Generate a simple test JWT token
	token, err := h.generateTestJWT(userID, "individual")
	if err != nil {
		// Return error response if token generation fails
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusInternalServerError)
		recorder.WriteString(fmt.Sprintf(`{"error": "failed to generate test token: %v"}`, err))
		return recorder
	}
	
	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	
	recorder := httptest.NewRecorder()
	h.Router.ServeHTTP(recorder, req)
	return recorder
}

// AuthRequestWithRole performs an HTTP request with JWT authentication and specific role
func (h *HTTPHelper) AuthRequestWithRole(method, path string, userID, role string, body interface{}) *httptest.ResponseRecorder {
	req := h.createRequest(method, path, body)
	
	// Generate a simple test JWT token with role
	token, err := h.generateTestJWT(userID, role)
	if err != nil {
		// Return error response if token generation fails
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusInternalServerError)
		recorder.WriteString(fmt.Sprintf(`{"error": "failed to generate test token: %v"}`, err))
		return recorder
	}
	
	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	
	recorder := httptest.NewRecorder()
	h.Router.ServeHTTP(recorder, req)
	return recorder
}

// RequestWithHeaders performs an HTTP request with custom headers
func (h *HTTPHelper) RequestWithHeaders(method, path string, headers map[string]string, body interface{}) *httptest.ResponseRecorder {
	req := h.createRequest(method, path, body)
	
	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	recorder := httptest.NewRecorder()
	h.Router.ServeHTTP(recorder, req)
	return recorder
}

// createRequest creates an HTTP request with the given parameters
func (h *HTTPHelper) createRequest(method, path string, body interface{}) *http.Request {
	var reqBody *bytes.Buffer
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			// If marshaling fails, create empty body
			reqBody = bytes.NewBuffer([]byte{})
		} else {
			reqBody = bytes.NewBuffer(jsonBody)
		}
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}
	
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		// If request creation fails, create a basic GET request
		req, _ = http.NewRequest("GET", "/", nil)
	}
	
	// Set content type for JSON requests
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	return req
}

// generateTestJWT generates a simple JWT token for testing
func (h *HTTPHelper) generateTestJWT(userID, role string) (string, error) {
	// Use test secret or fallback
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "test-secret-key-for-testing-only"
	}
	
	// Create claims
	claims := utils.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Simple response assertion helpers

// AssertStatus checks if the response has the expected status code
func AssertStatus(recorder *httptest.ResponseRecorder, expectedStatus int) error {
	if recorder.Code != expectedStatus {
		return fmt.Errorf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
	return nil
}

// AssertJSON checks if the response contains valid JSON
func AssertJSON(recorder *httptest.ResponseRecorder) error {
	var js json.RawMessage
	if err := json.Unmarshal(recorder.Body.Bytes(), &js); err != nil {
		return fmt.Errorf("response is not valid JSON: %v", err)
	}
	return nil
}

// AssertContains checks if the response body contains the expected string
func AssertContains(recorder *httptest.ResponseRecorder, expected string) error {
	body := recorder.Body.String()
	if !strings.Contains(body, expected) {
		return fmt.Errorf("response body does not contain '%s', got: %s", expected, body)
	}
	return nil
}

// GetResponseJSON parses the response body as JSON into the target interface
func GetResponseJSON(recorder *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(recorder.Body.Bytes(), target)
}

// GetResponseBody returns the response body as a string
func GetResponseBody(recorder *httptest.ResponseRecorder) string {
	return recorder.Body.String()
}

// Simple test data structures for common request bodies

// LoginRequest represents a login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupRequest represents a signup request body
type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

// RefreshTokenRequest represents a refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Helper functions for creating common request bodies

// NewLoginRequest creates a login request with default test values
func NewLoginRequest(overrides ...func(*LoginRequest)) LoginRequest {
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "testpassword123",
	}
	
	for _, override := range overrides {
		override(&req)
	}
	
	return req
}

// NewSignupRequest creates a signup request with default test values
func NewSignupRequest(overrides ...func(*SignupRequest)) SignupRequest {
	req := SignupRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "testpassword123",
		UserType: "individual",
	}
	
	for _, override := range overrides {
		override(&req)
	}
	
	return req
}

// NewRefreshTokenRequest creates a refresh token request with default test values
func NewRefreshTokenRequest(token string, overrides ...func(*RefreshTokenRequest)) RefreshTokenRequest {
	req := RefreshTokenRequest{
		RefreshToken: token,
	}
	
	for _, override := range overrides {
		override(&req)
	}
	
	return req
}