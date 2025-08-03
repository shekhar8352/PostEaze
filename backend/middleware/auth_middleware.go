package middleware

import (
	"net/http"
	"strings"

	"github.com/shekhar8352/PostEaze/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized", "msg": "Unauthorized"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenStr, false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized", "msg": "Invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIface, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No role in token"})
			return
		}
		userRole := roleIface.(string)
		for _, role := range allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	}
}
