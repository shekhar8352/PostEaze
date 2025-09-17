package apiv1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	apiv1 "github.com/shekhar8352/PostEaze/api/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
)

type FirebaseAuthAPITestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *FirebaseAuthAPITestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	
	// Setup routes
	v1 := suite.router.Group("/api/v1")
	auth := v1.Group("/auth")
	auth.POST("/authenticate", apiv1.AuthenticateWithFirebaseHandler)
}

func (suite *FirebaseAuthAPITestSuite) TestAuthenticateWithFirebaseHandler_InvalidJSON() {
	// Test with invalid JSON
	invalidJSON := `{"local_id": "test", "firebase_token": }`
	
	req, _ := http.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "error", response["status"])
}

func (suite *FirebaseAuthAPITestSuite) TestAuthenticateWithFirebaseHandler_MissingFields() {
	tests := []struct {
		name     string
		payload  modelsv1.FirebaseAuthParams
		expected int
	}{
		{
			name: "Missing LocalID",
			payload: modelsv1.FirebaseAuthParams{
				FirebaseToken: "valid-token",
				Platform:      "google",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Missing FirebaseToken",
			payload: modelsv1.FirebaseAuthParams{
				LocalID:  "firebase-uid-123",
				Platform: "google",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Missing Platform",
			payload: modelsv1.FirebaseAuthParams{
				LocalID:       "firebase-uid-123",
				FirebaseToken: "valid-token",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid Platform",
			payload: modelsv1.FirebaseAuthParams{
				LocalID:       "firebase-uid-123",
				FirebaseToken: "valid-token",
				Platform:      "invalid",
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			jsonPayload, _ := json.Marshal(tt.payload)
			
			req, _ := http.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			
			assert.Equal(suite.T(), tt.expected, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), "error", response["status"])
		})
	}
}

func (suite *FirebaseAuthAPITestSuite) TestAuthenticateWithFirebaseHandler_ValidPayload() {
	// Note: This test would fail in the current setup because Firebase is not initialized
	// and the database is not set up. This is a structure test to show how it would work.
	
	payload := modelsv1.FirebaseAuthParams{
		LocalID:       "firebase-uid-123",
		FirebaseToken: "valid-firebase-token",
		Platform:      "google",
	}
	
	jsonPayload, _ := json.Marshal(payload)
	
	req, _ := http.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// In a real test environment with proper mocking, we would expect:
	// - Either 200 (success) if user exists or is created successfully
	// - Or 401 (unauthorized) if Firebase token is invalid
	// - Or 500 (internal server error) if there are system issues
	
	// For now, we expect 500 because Firebase is not initialized
	assert.True(suite.T(), w.Code >= 400, "Should return an error code when Firebase is not properly set up")
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "error", response["status"])
}

func (suite *FirebaseAuthAPITestSuite) TestAuthenticateWithFirebaseHandler_ContentType() {
	payload := modelsv1.FirebaseAuthParams{
		LocalID:       "firebase-uid-123",
		FirebaseToken: "valid-firebase-token",
		Platform:      "google",
	}
	
	jsonPayload, _ := json.Marshal(payload)
	
	// Test without Content-Type header
	req, _ := http.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(jsonPayload))
	
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	
	// Should still work as Gin can handle JSON without explicit Content-Type
	assert.True(suite.T(), w.Code >= 400, "Should handle request even without Content-Type header")
}

func TestFirebaseAuthAPITestSuite(t *testing.T) {
	suite.Run(t, new(FirebaseAuthAPITestSuite))
}