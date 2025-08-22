package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"
)

func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		code    int
		message string
	}{
		{
			name:    "bad request error",
			code:    http.StatusBadRequest,
			message: "Invalid input provided",
		},
		{
			name:    "internal server error",
			code:    http.StatusInternalServerError,
			message: "Internal server error occurred",
		},
		{
			name:    "unauthorized error",
			code:    http.StatusUnauthorized,
			message: "Authentication required",
		},
		{
			name:    "empty message",
			code:    http.StatusBadRequest,
			message: "",
		},
		{
			name:    "long message",
			code:    http.StatusBadRequest,
			message: "This is a very long error message that contains multiple sentences and should still be handled correctly by the SendError function without any issues or truncation problems.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			utils.SendError(c, tt.code, tt.message)

			if w.Code != tt.code {
				t.Errorf("SendError() status code = %v, want %v", w.Code, tt.code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json; charset=utf-8" {
				t.Errorf("SendError() content type = %v, want application/json; charset=utf-8", contentType)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("SendError() response not valid JSON: %v", err)
				return
			}

			if response["status"] != "error" {
				t.Errorf("SendError() status = %v, want error", response["status"])
			}
			if response["msg"] != tt.message {
				t.Errorf("SendError() message = %v, want %v", response["msg"], tt.message)
			}
			if _, exists := response["data"]; exists {
				t.Error("SendError() should not include data field")
			}
		})
	}
}

func TestSendSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		data    interface{}
		message string
	}{
		{
			name:    "string data",
			data:    "Operation completed successfully",
			message: "Success",
		},
		{
			name: "map data",
			data: map[string]interface{}{
				"user_id": 123,
				"email":   "test@example.com",
				"role":    "admin",
			},
			message: "User retrieved successfully",
		},
		{
			name:    "slice data",
			data:    []string{"item1", "item2", "item3"},
			message: "Items retrieved successfully",
		},
		{
			name:    "nil data",
			data:    nil,
			message: "Operation completed",
		},
		{
			name:    "empty message",
			data:    "test data",
			message: "",
		},
		{
			name: "complex data",
			data: map[string]interface{}{
				"users": []map[string]interface{}{
					{"id": 1, "name": "John", "active": true},
					{"id": 2, "name": "Jane", "active": false},
				},
				"total": 2,
				"metadata": map[string]interface{}{
					"page":     1,
					"per_page": 10,
				},
			},
			message: "Users retrieved successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			utils.SendSuccess(c, tt.data, tt.message)

			if w.Code != http.StatusOK {
				t.Errorf("SendSuccess() status code = %v, want %v", w.Code, http.StatusOK)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json; charset=utf-8" {
				t.Errorf("SendSuccess() content type = %v, want application/json; charset=utf-8", contentType)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("SendSuccess() response not valid JSON: %v", err)
				return
			}

			if response["status"] != "success" {
				t.Errorf("SendSuccess() status = %v, want success", response["status"])
			}
			if response["msg"] != tt.message {
				t.Errorf("SendSuccess() message = %v, want %v", response["msg"], tt.message)
			}
			if _, exists := response["data"]; !exists {
				t.Error("SendSuccess() should include data field")
			}

			// Verify data content for specific test cases
			switch tt.name {
			case "map data":
				responseData := response["data"].(map[string]interface{})
				if responseData["user_id"] != float64(123) {
					t.Errorf("SendSuccess() user_id = %v, want 123", responseData["user_id"])
				}
				if responseData["email"] != "test@example.com" {
					t.Errorf("SendSuccess() email = %v, want test@example.com", responseData["email"])
				}
			case "slice data":
				responseData := response["data"].([]interface{})
				if len(responseData) != 3 {
					t.Errorf("SendSuccess() slice length = %v, want 3", len(responseData))
				}
				if responseData[0] != "item1" {
					t.Errorf("SendSuccess() first item = %v, want item1", responseData[0])
				}
			case "complex data":
				responseData := response["data"].(map[string]interface{})
				if responseData["total"] != float64(2) {
					t.Errorf("SendSuccess() total = %v, want 2", responseData["total"])
				}
				users := responseData["users"].([]interface{})
				if len(users) != 2 {
					t.Errorf("SendSuccess() users length = %v, want 2", len(users))
				}
			}
		})
	}
}

func TestResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("error vs success format consistency", func(t *testing.T) {
		// Test error response format
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		utils.SendError(c1, http.StatusBadRequest, "Error message")

		var errorResponse map[string]interface{}
		err := json.Unmarshal(w1.Body.Bytes(), &errorResponse)
		if err != nil {
			t.Fatalf("Error response not valid JSON: %v", err)
		}

		// Test success response format
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		utils.SendSuccess(c2, "data", "Success message")

		var successResponse map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &successResponse)
		if err != nil {
			t.Fatalf("Success response not valid JSON: %v", err)
		}

		// Both should have status and msg fields
		if _, exists := errorResponse["status"]; !exists {
			t.Error("Error response should have status field")
		}
		if _, exists := errorResponse["msg"]; !exists {
			t.Error("Error response should have msg field")
		}
		if _, exists := successResponse["status"]; !exists {
			t.Error("Success response should have status field")
		}
		if _, exists := successResponse["msg"]; !exists {
			t.Error("Success response should have msg field")
		}

		// Error response should not have data field
		if _, exists := errorResponse["data"]; exists {
			t.Error("Error response should not have data field")
		}

		// Success response should have data field
		if _, exists := successResponse["data"]; !exists {
			t.Error("Success response should have data field")
		}

		// Status values should be different
		if errorResponse["status"] != "error" {
			t.Errorf("Error response status = %v, want error", errorResponse["status"])
		}
		if successResponse["status"] != "success" {
			t.Errorf("Success response status = %v, want success", successResponse["status"])
		}
	})
}