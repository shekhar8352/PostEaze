package apiv1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

func ConnectGoogleDriveHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	var req modelsv1.ConnectGoogleDriveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.ConnectGoogleDrive(c.Request.Context(), userIDStr.(string), &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Google Drive connected")
}

func GetGoogleDriveStatusHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	resp, code, err := businessv1.GetGoogleDriveStatus(c.Request.Context(), userIDStr.(string))
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "OK")
}

func DisconnectGoogleDriveHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	code, err := businessv1.DisconnectGoogleDrive(c.Request.Context(), userIDStr.(string))
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, nil, "Google Drive disconnected")
}

func ListGoogleDriveFilesHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	resp, code, err := businessv1.ListGoogleDriveFiles(
		c.Request.Context(),
		userIDStr.(string),
		c.Query("folderId"),
		c.Query("pageToken"),
		c.Query("q"),
	)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "OK")
}

func ListGoogleDriveRevisionsHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	fileID := c.Param("fileId")
	resp, code, err := businessv1.ListGoogleDriveRevisions(c.Request.Context(), userIDStr.(string), fileID)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "OK")
}

func ImportGoogleDriveHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	var req modelsv1.ImportGoogleDriveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.ImportFromGoogleDrive(c.Request.Context(), userIDStr.(string), &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Imported from Google Drive")
}

func ImportDriveRevisionHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	assetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req modelsv1.ImportDriveRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.ImportDriveRevision(c.Request.Context(), userIDStr.(string), assetID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Drive revision imported")
}
