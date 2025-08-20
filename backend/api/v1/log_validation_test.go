package apiv1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"
)

func TestGetLogByIDHandler_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		logID          string
		expectedStatus int
		expectedError  string
		expectedType   string
	}{
		{
			name:           "whitespace log ID",
			logID:          "   ",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Log ID is required",
			expectedType:   utils.ErrorTypeInvalidInput,
		},
		{
			name:           "invalid log ID with special characters",
			logID:          "log@id#123",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid log ID format. Only alphanumeric characters, hyphens, underscores, and dots are allowed",
			expectedType:   utils.ErrorTypeInvalidInput,
		},
		{
			name:           "log ID with spaces",
			logID:          "log id 123",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid log ID format. Only alphanumeric characters, hyphens, underscores, and dots are allowed",
			expectedType:   utils.ErrorTypeInvalidInput,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/api/v1/log/byId/:log_id", GetLogByIDHandler)

			req, _ := http.NewRequest("GET", "/api/v1/log/byId/"+tt.logID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if success, ok := response["success"].(bool); !ok || success {
				t.Errorf("Expected success to be false")
			}

			if errorObj, ok := response["error"].(map[string]interface{}); ok {
				if message, ok := errorObj["message"].(string); !ok || message != tt.expectedError {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedError, message)
				}
				if errorType, ok := errorObj["type"].(string); !ok || errorType != tt.expectedType {
					t.Errorf("Expected error type '%s', got '%s'", tt.expectedType, errorType)
				}
			} else {
				t.Errorf("Expected error object in response")
			}
		})
	}
}

func TestGetLogsByDate_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		date           string
		expectedStatus int
		expectedError  string
		expectedType   string
	}{
		{
			name:           "whitespace date",
			date:           "   ",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Date is required",
			expectedType:   utils.ErrorTypeInvalidInput,
		},
		{
			name:           "invalid date format - wrong order",
			date:           "01-15-2024",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date format. Expected YYYY-MM-DD",
			expectedType:   utils.ErrorTypeInvalidInput,
		},
		{
			name:           "invalid date format - missing day",
			date:           "2024-01",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date format. Expected YYYY-MM-DD",
			expectedType:   utils.ErrorTypeInvalidInput,
		},
		{
			name:           "invalid date - February 30th",
			date:           "2024-02-30",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid date format. Expected YYYY-MM-DD",
			expectedType:   utils.ErrorTypeInvalidInput,
		},
		{
			name:           "date too far in past",
			date:           "2010-01-01",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Date is out of valid range. Please provide a date within the last 5 years",
			expectedType:   utils.ErrorTypeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/api/v1/log/byDate/:date", GetLogsByDate)

			req, _ := http.NewRequest("GET", "/api/v1/log/byDate/"+tt.date, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if success, ok := response["success"].(bool); !ok || success {
				t.Errorf("Expected success to be false")
			}

			if errorObj, ok := response["error"].(map[string]interface{}); ok {
				if message, ok := errorObj["message"].(string); !ok || message != tt.expectedError {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedError, message)
				}
				if errorType, ok := errorObj["type"].(string); !ok || errorType != tt.expectedType {
					t.Errorf("Expected error type '%s', got '%s'", tt.expectedType, errorType)
				}
			} else {
				t.Errorf("Expected error object in response")
			}
		})
	}
}