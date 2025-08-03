package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"status": "error", "msg": msg})
}

func SendSuccess(c *gin.Context, data any, msg string) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "msg": msg, "data": data})
}