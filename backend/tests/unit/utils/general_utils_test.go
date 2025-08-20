package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GeneralUtilsTestSuite defines the test suite for general utilities
type GeneralUtilsTestSuite struct {
	suite.Suite
	router *gin.Engine
}

// SetupSuite runs before all tests in the suite
func (suite *GeneralUtilsTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
}

// TestSendError tests the SendError utility function
func (suite *GeneralUtilsTestSuite) TestSendError_StandardError() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	code := http.StatusBadRequest
	message := "Invalid input provided"
	
	// Execute
	utils.SendError(c, code, message)
	
	// Assert
	assert.Equal(suite.T(), code, w.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "error", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
}

func (suite *GeneralUtilsTestSuite) TestSendError_InternalServerError() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	code := http.StatusInternalServerError
	message := "Internal server error occurred"
	
	// Execute
	utils.SendError(c, code, message)
	
	// Assert
	assert.Equal(suite.T(), code, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "error", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
}

func (suite *GeneralUtilsTestSuite) TestSendError_UnauthorizedError() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	code := http.StatusUnauthorized
	message := "Authentication required"
	
	// Execute
	utils.SendError(c, code, message)
	
	// Assert
	assert.Equal(suite.T(), code, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "error", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
}

func (suite *GeneralUtilsTestSuite) TestSendError_EmptyMessage() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	code := http.StatusBadRequest
	message := ""
	
	// Execute
	utils.SendError(c, code, message)
	
	// Assert
	assert.Equal(suite.T(), code, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "error", response["status"])
	assert.Equal(suite.T(), "", response["msg"])
}

func (suite *GeneralUtilsTestSuite) TestSendError_LongMessage() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	code := http.StatusBadRequest
	message := "This is a very long error message that contains multiple sentences and should still be handled correctly by the SendError function without any issues or truncation problems."
	
	// Execute
	utils.SendError(c, code, message)
	
	// Assert
	assert.Equal(suite.T(), code, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "error", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
}

// TestSendSuccess tests the SendSuccess utility function
func (suite *GeneralUtilsTestSuite) TestSendSuccess_WithStringData() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	data := "Operation completed successfully"
	message := "Success"
	
	// Execute
	utils.SendSuccess(c, data, message)
	
	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "success", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
	assert.Equal(suite.T(), data, response["data"])
}

func (suite *GeneralUtilsTestSuite) TestSendSuccess_WithMapData() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	data := map[string]interface{}{
		"user_id": 123,
		"email":   "test@example.com",
		"role":    "admin",
	}
	message := "User retrieved successfully"
	
	// Execute
	utils.SendSuccess(c, data, message)
	
	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "success", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
	
	responseData := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), float64(123), responseData["user_id"])
	assert.Equal(suite.T(), "test@example.com", responseData["email"])
	assert.Equal(suite.T(), "admin", responseData["role"])
}

func (suite *GeneralUtilsTestSuite) TestSendSuccess_WithSliceData() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	data := []string{"item1", "item2", "item3"}
	message := "Items retrieved successfully"
	
	// Execute
	utils.SendSuccess(c, data, message)
	
	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "success", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
	
	responseData := response["data"].([]interface{})
	assert.Len(suite.T(), responseData, 3)
	assert.Equal(suite.T(), "item1", responseData[0])
	assert.Equal(suite.T(), "item2", responseData[1])
	assert.Equal(suite.T(), "item3", responseData[2])
}

func (suite *GeneralUtilsTestSuite) TestSendSuccess_WithNilData() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	var data interface{} = nil
	message := "Operation completed"
	
	// Execute
	utils.SendSuccess(c, data, message)
	
	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "success", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
	assert.Nil(suite.T(), response["data"])
}

func (suite *GeneralUtilsTestSuite) TestSendSuccess_WithEmptyMessage() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	data := "test data"
	message := ""
	
	// Execute
	utils.SendSuccess(c, data, message)
	
	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "success", response["status"])
	assert.Equal(suite.T(), "", response["msg"])
	assert.Equal(suite.T(), "test data", response["data"])
}

func (suite *GeneralUtilsTestSuite) TestSendSuccess_WithComplexData() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	data := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "John", "active": true},
			{"id": 2, "name": "Jane", "active": false},
		},
		"total": 2,
		"metadata": map[string]interface{}{
			"page":     1,
			"per_page": 10,
		},
	}
	message := "Users retrieved successfully"
	
	// Execute
	utils.SendSuccess(c, data, message)
	
	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "success", response["status"])
	assert.Equal(suite.T(), message, response["msg"])
	
	responseData := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), float64(2), responseData["total"])
	
	users := responseData["users"].([]interface{})
	assert.Len(suite.T(), users, 2)
	
	firstUser := users[0].(map[string]interface{})
	assert.Equal(suite.T(), float64(1), firstUser["id"])
	assert.Equal(suite.T(), "John", firstUser["name"])
	assert.Equal(suite.T(), true, firstUser["active"])
}

// TestResponseFormat tests the consistent response format
func (suite *GeneralUtilsTestSuite) TestResponseFormat_ErrorVsSuccess() {
	// Test error response format
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	utils.SendError(c1, http.StatusBadRequest, "Error message")
	
	var errorResponse map[string]interface{}
	err := json.Unmarshal(w1.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	
	// Test success response format
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	utils.SendSuccess(c2, "data", "Success message")
	
	var successResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &successResponse)
	assert.NoError(suite.T(), err)
	
	// Both should have status and msg fields
	assert.Contains(suite.T(), errorResponse, "status")
	assert.Contains(suite.T(), errorResponse, "msg")
	assert.Contains(suite.T(), successResponse, "status")
	assert.Contains(suite.T(), successResponse, "msg")
	
	// Error response should not have data field
	assert.NotContains(suite.T(), errorResponse, "data")
	
	// Success response should have data field
	assert.Contains(suite.T(), successResponse, "data")
	
	// Status values should be different
	assert.Equal(suite.T(), "error", errorResponse["status"])
	assert.Equal(suite.T(), "success", successResponse["status"])
}

// Run the test suite
func TestGeneralUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(GeneralUtilsTestSuite))
}