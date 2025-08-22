package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogAPIError represents a structured error response for log API
type LogAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Error types for classification
const (
	ErrorTypeNotFound     = "not_found"
	ErrorTypeInvalidInput = "invalid_input"
	ErrorTypeInternal     = "internal_error"
)

func SendError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"status": "error", "msg": msg})
}

// SendLogAPIError sends a structured error response for log API endpoints
func SendLogAPIError(c *gin.Context, code int, message, errorType string) {
	c.JSON(code, gin.H{
		"success": false,
		"error": LogAPIError{
			Code:    code,
			Message: message,
			Type:    errorType,
		},
	})
}

func SendSuccess(c *gin.Context, data any, msg string) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "msg": msg, "data": data})
}

// SendLogAPISuccess sends a structured success response for log API endpoints
func SendLogAPISuccess(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"message": message,
	})
}