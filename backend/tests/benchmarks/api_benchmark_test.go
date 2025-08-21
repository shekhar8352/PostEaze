package benchmarks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

// setupBenchmarkRouter creates a router for API benchmarking
func setupBenchmarkRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add simple mock auth endpoints for benchmarking
	authGroup := router.Group("/api/v1/auth")
	authGroup.POST("/signup", func(c *gin.Context) {
		var body modelsv1.SignupParams
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signup data"})
			return
		}
		// Mock successful response
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"user": gin.H{
					"id":    "mock-user-id",
					"name":  body.Name,
					"email": body.Email,
				},
				"access_token":  "mock-access-token",
				"refresh_token": "mock-refresh-token",
			},
		})
	})
	
	authGroup.POST("/login", func(c *gin.Context) {
		var body modelsv1.LoginParams
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
			return
		}
		// Mock successful response
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"user": gin.H{
					"id":    "mock-user-id",
					"email": body.Email,
				},
				"access_token":  "mock-access-token",
				"refresh_token": "mock-refresh-token",
			},
		})
	})
	
	authGroup.POST("/refresh", func(c *gin.Context) {
		var body modelsv1.RefreshTokenParams
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
			return
		}
		// Mock successful response
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"access_token": "mock-new-access-token",
			},
		})
	})
	
	authGroup.POST("/logout", func(c *gin.Context) {
		// Mock successful response
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Logged out successfully",
		})
	})
	
	return router
}

// BenchmarkSignupEndpoint benchmarks the signup API endpoint
func BenchmarkSignupEndpoint(b *testing.B) {
	router := setupBenchmarkRouter()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create unique request data for each iteration
		signupData := modelsv1.SignupParams{
			Name:     "Test User",
			Email:    "test" + time.Now().Format("20060102150405.000000") + "@example.com",
			Password: "testpassword123",
			UserType: modelsv1.UserTypeIndividual,
		}
		
		jsonData, _ := json.Marshal(signupData)
		req := httptest.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		b.StartTimer()

		router.ServeHTTP(recorder, req)
		
		// Verify successful response
		if recorder.Code != http.StatusOK && recorder.Code != http.StatusCreated {
			b.Logf("Signup endpoint returned status: %d", recorder.Code)
		}
	}
}

// BenchmarkLoginEndpoint benchmarks the login API endpoint
func BenchmarkLoginEndpoint(b *testing.B) {
	router := setupBenchmarkRouter()
	
	loginData := modelsv1.LoginParams{
		Email:    "benchmark@example.com",
		Password: "testpassword123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)
		
		// Verify successful response
		if recorder.Code != http.StatusOK {
			b.Logf("Login endpoint returned status: %d", recorder.Code)
		}
	}
}

// BenchmarkRefreshTokenEndpoint benchmarks the refresh token API endpoint
func BenchmarkRefreshTokenEndpoint(b *testing.B) {
	router := setupBenchmarkRouter()
	
	refreshData := modelsv1.RefreshTokenParams{
		RefreshToken: "mock-refresh-token",
	}
	
	jsonData, _ := json.Marshal(refreshData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)
		
		// Verify successful response
		if recorder.Code != http.StatusOK {
			b.Logf("Refresh endpoint returned status: %d", recorder.Code)
		}
	}
}

// BenchmarkLogoutEndpoint benchmarks the logout API endpoint
func BenchmarkLogoutEndpoint(b *testing.B) {
	router := setupBenchmarkRouter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer mock-refresh-token")
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)
		
		// Verify successful response
		if recorder.Code != http.StatusOK {
			b.Logf("Logout endpoint returned status: %d", recorder.Code)
		}
	}
}

// BenchmarkConcurrentAPIRequests benchmarks concurrent API requests
func BenchmarkConcurrentAPIRequests(b *testing.B) {
	router := setupBenchmarkRouter()
	
	// Prepare request data
	signupData := modelsv1.SignupParams{
		Name:     "Concurrent User",
		Email:    "concurrent@example.com",
		Password: "testpassword123",
		UserType: modelsv1.UserTypeIndividual,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Create unique email for each request
			uniqueData := signupData
			uniqueData.Email = "concurrent" + time.Now().Format("20060102150405.000000") + "@example.com"
			uniqueJSON, _ := json.Marshal(uniqueData)
			
			req := httptest.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(uniqueJSON))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)
		}
	})
}

// BenchmarkJSONMarshaling benchmarks JSON marshaling for API responses
func BenchmarkJSONMarshaling(b *testing.B) {
	user := modelsv1.User{
		ID:       "user-123",
		Name:     "Test User",
		Email:    "test@example.com",
		UserType: modelsv1.UserTypeIndividual,
	}

	response := map[string]interface{}{
		"user":          user,
		"access_token":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(response)
		if err != nil {
			b.Fatalf("JSON marshaling failed: %v", err)
		}
	}
}

// BenchmarkJSONUnmarshaling benchmarks JSON unmarshaling for API requests
func BenchmarkJSONUnmarshaling(b *testing.B) {
	jsonData := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "testpassword123",
		"user_type": "individual"
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var params modelsv1.SignupParams
		err := json.Unmarshal([]byte(jsonData), &params)
		if err != nil {
			b.Fatalf("JSON unmarshaling failed: %v", err)
		}
	}
}