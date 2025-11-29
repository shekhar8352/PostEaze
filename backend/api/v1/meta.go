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
