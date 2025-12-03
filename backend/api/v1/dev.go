package apiv1

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/utils"
)

// GenerateTestTokenRequest represents the request to generate a test token
type GenerateTestTokenRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// GenerateTestTokenResponse represents the response with the generated token
type GenerateTestTokenResponse struct {
	Token string `json:"token"`
}

// GenerateTestTokenHandler godoc
// @Summary      Generate Test Token (Development Only)
// @Description  Generates a JWT token for testing purposes. Only available in development mode.
// @Tags         Development
// @Accept       json
// @Produce      json
// @Param        request body GenerateTestTokenRequest true "User ID"
// @Success      200  {object}  GenerateTestTokenResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /dev/generate-token [post]
func GenerateTestTokenHandler(c *gin.Context) {
	// Only allow in development mode
	env := os.Getenv("ENV")
	if env != "development" && env != "dev" && env != "" {
		utils.SendError(c, http.StatusForbidden, "This endpoint is only available in development mode")
		return
	}

	var req GenerateTestTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Generate token with long validity (30 days for testing)
	token, err := utils.GenerateDevToken(req.UserID, "user", 30)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to generate token: "+err.Error())
		return
	}

	utils.SendSuccess(c, GenerateTestTokenResponse{Token: token}, "Test token generated successfully (valid for 30 days)")
}
