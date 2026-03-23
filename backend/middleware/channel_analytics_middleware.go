package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/utils"
)

const ContextKeyAnalyticsChannel = "analytics_channel"

// RequireInstagramChannelAnalyticsAccess ensures the user is authenticated (run after AuthMiddleware),
// owns the channel or is an active team member, and the channel provider is Instagram.
func RequireInstagramChannelAnalyticsAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, ok := c.Get("user_id")
		if !ok {
			utils.SendError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		userID, _ := userIDVal.(string)

		channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
		if err != nil {
			utils.SendError(c, http.StatusBadRequest, "Invalid channel ID")
			c.Abort()
			return
		}

		allowed, ch, err := repositories.UserCanAccessChannel(c.Request.Context(), channelID, userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.SendError(c, http.StatusNotFound, "Channel not found")
				c.Abort()
				return
			}
			utils.SendError(c, http.StatusInternalServerError, "Failed to verify channel access")
			c.Abort()
			return
		}
		if ch == nil || !allowed {
			utils.SendError(c, http.StatusForbidden, "You do not have access to this channel")
			c.Abort()
			return
		}

		if ch.Provider != "instagram" {
			utils.SendError(c, http.StatusBadRequest, "Analytics is only available for Instagram channels")
			c.Abort()
			return
		}

		c.Set(ContextKeyAnalyticsChannel, ch)
		c.Next()
	}
}
