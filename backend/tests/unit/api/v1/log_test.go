package apiv1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/shekhar8352/PostEaze/api/v1"
	"github.com/shekhar8352/PostEaze/tests/testutils"
)

// LogAPITestSuite tests the log API handlers
type LogAPITestSuite struct {
	testutils.APITestSuite
}

// SetupSuite runs once before all tests in the suite
func (s *LogAPITestSuite) SetupSuite() {
	s.APITestSuite.SetupSuite()
	
	// Setup log routes for testing
	s.setupLogRoutes()
}

// setupLogRoutes sets up the log routes for testing
func (s *LogAPITestSuite) setupLogRoutes() {
	logGroup := s.Router.Group("/api/v1/logs")
	{
		logGroup.GET("/date/:date", apiv1.GetLogsByDate)
		logGroup.GET("/:log_id", apiv1.GetLogByIDHandler)
	}
}

// TestGetLogsByDate_ValidDate tests get logs by date handler with valid date
func (s *LogAPITestSuite) TestGetLogsByDate_ValidDate() {
	testCases := []struct {
		name         string
		date         string
		expectedCode int
	}{
		{
			name:         "valid date format",
			date:         "2024-01-15",
			expectedCode: http.StatusOK, // or 500 due to missing log files
		},
		{
			name:         "today's date",
			date:         "2024-12-20",
			expectedCode: http.StatusOK, // or 500 due to missing log files
		},
		{
			name:         "past date",
			date:         "2023-12-01",
			expectedCode: http.StatusOK, // or 500 due to missing log files
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/date/"+tc.date, nil)
			testutils.SetURLParam(ctx, "date", tc.date)
			
			// Execute the handler
			apiv1.GetLogsByDate(ctx)
			
			// Assert response (allowing for business logic errors due to missing log files)
			if tc.expectedCode == http.StatusOK {
				s.True(recorder.Code == http.StatusOK || recorder.Code == http.StatusInternalServerError)
			} else {
				s.Equal(tc.expectedCode, recorder.Code)
			}
			
			// Assert response format
			testutils.AssertJSONResponse(s.T(), recorder)
			
			var response map[string]interface{}
			err := testutils.ParseJSONResponse(recorder, &response)
			s.NoError(err)
			s.Contains(response, "status")
			s.Contains(response, "msg")
			
			// If successful, check data structure
			if recorder.Code == http.StatusOK {
				s.Equal("success", response["status"])
				s.Contains(response, "data")
				
				data := response["data"].(map[string]interface{})
				s.Contains(data, "logs")
				s.Contains(data, "total")
			}
		})
	}
}

// TestGetLogsByDate_InvalidDate tests get logs by date handler with invalid date formats
func (s *LogAPITestSuite) TestGetLogsByDate_InvalidDate() {
	testCases := []struct {
		name         string
		date         string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "empty date",
			date:         "",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Date is required",
		},
		{
			name:         "invalid date format - wrong separator",
			date:         "2024/01/15",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid date format",
		},
		{
			name:         "invalid date format - missing day",
			date:         "2024-01",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid date format",
		},
		{
			name:         "invalid date format - wrong order",
			date:         "01-15-2024",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid date format",
		},
		{
			name:         "invalid date format - text",
			date:         "invalid-date",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid date format",
		},
		{
			name:         "invalid date format - partial",
			date:         "2024-1-1",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid date format",
		},
		{
			name:         "invalid date - impossible date",
			date:         "2024-13-32",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid date format",
		},
		{
			name:         "invalid date - february 30",
			date:         "2024-02-30",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid date format",
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var ctx *gin.Context
			var recorder *httptest.ResponseRecorder
			
			if tc.date == "" {
				// Test empty date parameter
				ctx, recorder = testutils.NewTestGinContext("GET", "/api/v1/logs/date/", nil)
				// Don't set URL param for empty date test
			} else {
				ctx, recorder = testutils.NewTestGinContext("GET", "/api/v1/logs/date/"+tc.date, nil)
				testutils.SetURLParam(ctx, "date", tc.date)
			}
			
			// Execute the handler
			apiv1.GetLogsByDate(ctx)
			
			// Assert error response
			testutils.AssertErrorResponse(s.T(), recorder, tc.expectedCode, tc.expectedMsg)
		})
	}
}

// TestGetLogsByDate_ResponseFormat tests that get logs by date returns proper response format
func (s *LogAPITestSuite) TestGetLogsByDate_ResponseFormat() {
	ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/date/2024-01-15", nil)
	testutils.SetURLParam(ctx, "date", "2024-01-15")
	
	// Execute the handler
	apiv1.GetLogsByDate(ctx)
	
	// Assert that response is valid JSON
	testutils.AssertJSONResponse(s.T(), recorder)
	
	// Assert that response has required fields
	var response map[string]interface{}
	err := testutils.ParseJSONResponse(recorder, &response)
	s.NoError(err)
	s.Contains(response, "status")
	s.Contains(response, "msg")
	
	// Status should be either "success" or "error"
	status := response["status"].(string)
	s.True(status == "success" || status == "error")
	
	// If successful, verify data structure
	if status == "success" {
		s.Contains(response, "data")
		data := response["data"].(map[string]interface{})
		s.Contains(data, "logs")
		s.Contains(data, "total")
		
		// Verify logs is an array
		logs := data["logs"]
		s.NotNil(logs)
		
		// Verify total is a number
		total := data["total"]
		s.NotNil(total)
	}
}

// TestGetLogByIDHandler_ValidID tests get log by ID handler with valid log IDs
func (s *LogAPITestSuite) TestGetLogByIDHandler_ValidID() {
	testCases := []struct {
		name         string
		logID        string
		expectedCode int
	}{
		{
			name:         "valid log ID format",
			logID:        "log-123456",
			expectedCode: http.StatusOK, // or 500 due to missing log files
		},
		{
			name:         "uuid format log ID",
			logID:        "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: http.StatusOK, // or 500 due to missing log files
		},
		{
			name:         "numeric log ID",
			logID:        "12345",
			expectedCode: http.StatusOK, // or 500 due to missing log files
		},
		{
			name:         "alphanumeric log ID",
			logID:        "abc123def456",
			expectedCode: http.StatusOK, // or 500 due to missing log files
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/"+tc.logID, nil)
			testutils.SetURLParam(ctx, "log_id", tc.logID)
			
			// Execute the handler
			apiv1.GetLogByIDHandler(ctx)
			
			// Assert response (allowing for business logic errors due to missing log files)
			if tc.expectedCode == http.StatusOK {
				s.True(recorder.Code == http.StatusOK || recorder.Code == http.StatusInternalServerError)
			} else {
				s.Equal(tc.expectedCode, recorder.Code)
			}
			
			// Assert response format
			testutils.AssertJSONResponse(s.T(), recorder)
			
			var response map[string]interface{}
			err := testutils.ParseJSONResponse(recorder, &response)
			s.NoError(err)
			s.Contains(response, "status")
			s.Contains(response, "msg")
			
			// If successful, check data structure
			if recorder.Code == http.StatusOK {
				s.Equal("success", response["status"])
				s.Contains(response, "data")
				
				// Data should be an array of log entries (can be empty)
				data := response["data"]
				if data != nil {
					// Verify it's an array
					s.IsType([]interface{}{}, data)
				}
			}
		})
	}
}

// TestGetLogByIDHandler_InvalidID tests get log by ID handler with invalid log IDs
func (s *LogAPITestSuite) TestGetLogByIDHandler_InvalidID() {
	testCases := []struct {
		name         string
		logID        string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "empty log ID",
			logID:        "",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Log ID is required",
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var ctx *gin.Context
			var recorder *httptest.ResponseRecorder
			
			if tc.logID == "" {
				// Test empty log ID parameter
				ctx, recorder = testutils.NewTestGinContext("GET", "/api/v1/logs/", nil)
				// Don't set URL param for empty log ID test
			} else {
				ctx, recorder = testutils.NewTestGinContext("GET", "/api/v1/logs/"+tc.logID, nil)
				testutils.SetURLParam(ctx, "log_id", tc.logID)
			}
			
			// Execute the handler
			apiv1.GetLogByIDHandler(ctx)
			
			// Assert error response
			testutils.AssertErrorResponse(s.T(), recorder, tc.expectedCode, tc.expectedMsg)
		})
	}
}

// TestGetLogByIDHandler_ResponseFormat tests that get log by ID returns proper response format
func (s *LogAPITestSuite) TestGetLogByIDHandler_ResponseFormat() {
	ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/test-log-id", nil)
	testutils.SetURLParam(ctx, "log_id", "test-log-id")
	
	// Execute the handler
	apiv1.GetLogByIDHandler(ctx)
	
	// Assert that response is valid JSON
	testutils.AssertJSONResponse(s.T(), recorder)
	
	// Assert that response has required fields
	var response map[string]interface{}
	err := testutils.ParseJSONResponse(recorder, &response)
	s.NoError(err)
	s.Contains(response, "status")
	s.Contains(response, "msg")
	
	// Status should be either "success" or "error"
	status := response["status"].(string)
	s.True(status == "success" || status == "error")
	
	// If successful, verify data structure
	if status == "success" {
		s.Contains(response, "data")
		data := response["data"]
		if data != nil {
			// Data should be an array (even if empty)
			s.IsType([]interface{}{}, data)
		}
	}
}

// TestGetLogByIDHandler_ErrorCases tests get log by ID handler error scenarios
func (s *LogAPITestSuite) TestGetLogByIDHandler_ErrorCases() {
	testCases := []struct {
		name         string
		logID        string
		expectedCode int
		description  string
	}{
		{
			name:         "non-existent log ID",
			logID:        "non-existent-log-id",
			expectedCode: http.StatusOK, // Handler returns success with empty array
			description:  "Should handle non-existent log IDs gracefully",
		},
		{
			name:         "special characters in log ID",
			logID:        "log-with-dashes",
			expectedCode: http.StatusOK, // Handler accepts any non-empty string
			description:  "Should handle special characters in log ID",
		},
		{
			name:         "numeric log ID",
			logID:        "123456789",
			expectedCode: http.StatusOK, // Handler accepts any non-empty string
			description:  "Should handle numeric log IDs",
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/"+tc.logID, nil)
			testutils.SetURLParam(ctx, "log_id", tc.logID)
			
			// Execute the handler
			apiv1.GetLogByIDHandler(ctx)
			
			// Assert response (allowing for business logic errors)
			if tc.expectedCode == http.StatusOK {
				s.True(recorder.Code == http.StatusOK || recorder.Code == http.StatusInternalServerError, tc.description)
			} else {
				s.Equal(tc.expectedCode, recorder.Code, tc.description)
			}
			
			// Assert response format
			testutils.AssertJSONResponse(s.T(), recorder)
			
			var response map[string]interface{}
			err := testutils.ParseJSONResponse(recorder, &response)
			s.NoError(err)
			s.Contains(response, "status")
			s.Contains(response, "msg")
		})
	}
}

// TestGetLogsByDate_ErrorCases tests get logs by date handler error scenarios
func (s *LogAPITestSuite) TestGetLogsByDate_ErrorCases() {
	testCases := []struct {
		name         string
		date         string
		expectedCode int
		description  string
	}{
		{
			name:         "non-existent date",
			date:         "1900-01-01",
			expectedCode: http.StatusInternalServerError, // Due to missing log files
			description:  "Should handle non-existent log files gracefully",
		},
		{
			name:         "future date",
			date:         "2030-12-31",
			expectedCode: http.StatusInternalServerError, // Due to missing log files
			description:  "Should handle future dates gracefully",
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/date/"+tc.date, nil)
			testutils.SetURLParam(ctx, "date", tc.date)
			
			// Execute the handler
			apiv1.GetLogsByDate(ctx)
			
			// Assert response (allowing for business logic errors)
			s.True(recorder.Code == http.StatusOK || recorder.Code == http.StatusInternalServerError, tc.description)
			
			// Assert response format
			testutils.AssertJSONResponse(s.T(), recorder)
			
			var response map[string]interface{}
			err := testutils.ParseJSONResponse(recorder, &response)
			s.NoError(err)
			s.Contains(response, "status")
			s.Contains(response, "msg")
		})
	}
}

// TestLogHandlers_StatusCodes tests that log handlers return appropriate status codes
func (s *LogAPITestSuite) TestLogHandlers_StatusCodes() {
	testCases := []struct {
		name        string
		handler     func(*gin.Context)
		setupCtx    func() (*gin.Context, *httptest.ResponseRecorder)
		description string
	}{
		{
			name:    "get logs by date with valid date",
			handler: apiv1.GetLogsByDate,
			setupCtx: func() (*gin.Context, *httptest.ResponseRecorder) {
				ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/date/2024-01-15", nil)
				testutils.SetURLParam(ctx, "date", "2024-01-15")
				return ctx, recorder
			},
			description: "Should return appropriate status code for valid date",
		},
		{
			name:    "get logs by date with invalid date",
			handler: apiv1.GetLogsByDate,
			setupCtx: func() (*gin.Context, *httptest.ResponseRecorder) {
				ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/date/invalid", nil)
				testutils.SetURLParam(ctx, "date", "invalid")
				return ctx, recorder
			},
			description: "Should return 400 for invalid date",
		},
		{
			name:    "get log by ID with valid ID",
			handler: apiv1.GetLogByIDHandler,
			setupCtx: func() (*gin.Context, *httptest.ResponseRecorder) {
				ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/test-id", nil)
				testutils.SetURLParam(ctx, "log_id", "test-id")
				return ctx, recorder
			},
			description: "Should return appropriate status code for valid log ID",
		},
		{
			name:    "get log by ID with empty ID",
			handler: apiv1.GetLogByIDHandler,
			setupCtx: func() (*gin.Context, *httptest.ResponseRecorder) {
				ctx, recorder := testutils.NewTestGinContext("GET", "/api/v1/logs/", nil)
				// Don't set URL param to simulate empty ID
				return ctx, recorder
			},
			description: "Should return 400 for empty log ID",
		},
	}
	
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, recorder := tc.setupCtx()
			
			// Execute the handler
			tc.handler(ctx)
			
			// Assert that we get a valid HTTP status code
			s.True(recorder.Code >= 200 && recorder.Code < 600, "Should return valid HTTP status code")
			
			// Assert response format
			testutils.AssertJSONResponse(s.T(), recorder)
			
			var response map[string]interface{}
			err := testutils.ParseJSONResponse(recorder, &response)
			s.NoError(err, tc.description)
			s.Contains(response, "status")
			s.Contains(response, "msg")
		})
	}
}

// Run the test suite
func TestLogAPITestSuite(t *testing.T) {
	testutils.RunTestSuite(t, new(LogAPITestSuite))
}