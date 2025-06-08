package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

func SendSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}