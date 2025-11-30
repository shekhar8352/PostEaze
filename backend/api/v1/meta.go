package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shekhar8352/PostEaze/services/meta_service"
)

type MetaCallbackRequest struct {
	Code        string `json:"code" binding:"required"`
	State       string `json:"state"`
	RedirectURI string `json:"redirect_uri" binding:"required"`
}

// HandleMetaCallback godoc
// @Summary      Handle Meta Callback
// @Description  Exchanges authorization code for access token and fetches pages
// @Tags         Meta
// @Accept       json
// @Produce      json
// @Param        request body MetaCallbackRequest true "Meta Callback Request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /meta/callback [post]
func HandleMetaCallback(c *gin.Context) {
	var req MetaCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := meta_service.NewMetaService()
	pages, err := service.GetPagesFromCode(req.Code, req.RedirectURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages})
}
