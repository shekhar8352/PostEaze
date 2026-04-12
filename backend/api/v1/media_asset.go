package apiv1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// UploadMediaHandler godoc
// @Summary      Upload a media file to blob storage
// @Tags         Media
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Param        file formData file true "Media file (max 50MB)"
// @Success      200 {object} map[string]interface{}
// @Router       /media/upload [post]
func UploadMediaHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	resp, code, err := businessv1.UploadMedia(
		c.Request.Context(),
		userIDStr.(string),
		header.Filename,
		contentType,
		header.Size,
		file,
	)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "File uploaded")
}

// CreateMediaAssetHandler godoc
// @Summary      Create a media asset with initial version
// @Tags         Media
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Param        file  formData file   true  "Media file"
// @Param        title formData string true  "Asset title"
// @Param        asset_type formData string true "photo or video"
// @Param        label formData string false "Version label (default: raw)"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets [post]
func CreateMediaAssetHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	title := c.PostForm("title")
	assetType := c.PostForm("asset_type")
	label := c.PostForm("label")

	if title == "" {
		title = header.Filename
	}
	if assetType == "" {
		ct := header.Header.Get("Content-Type")
		if strings.HasPrefix(ct, "video/") {
			assetType = "video"
		} else {
			assetType = "photo"
		}
	}
	if assetType != "photo" && assetType != "video" {
		utils.SendError(c, http.StatusBadRequest, "asset_type must be photo or video")
		return
	}

	contentType := header.Header.Get("Content-Type")
	upload, code, err := businessv1.UploadMedia(
		c.Request.Context(), userIDStr.(string),
		header.Filename, contentType, header.Size, file,
	)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}

	req := &modelsv1.CreateMediaAssetRequest{
		Title:     title,
		AssetType: assetType,
		Label:     label,
	}
	resp, code, err := businessv1.CreateMediaAsset(c.Request.Context(), userIDStr.(string), req, upload)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Media asset created")
}

// ListMediaAssetsHandler godoc
// @Summary      List media assets
// @Tags         Media
// @Security     BearerAuth
// @Param        status query string false "Filter by status (draft, ready, published)"
// @Param        limit  query int    false "Limit (default 20)"
// @Param        offset query int    false "Offset (default 0)"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets [get]
func ListMediaAssetsHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	var q modelsv1.ListMediaAssetsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, code, err := businessv1.ListMediaAssets(c.Request.Context(), userIDStr.(string), q)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Media assets retrieved")
}

// GetMediaAssetHandler godoc
// @Summary      Get a media asset with all versions
// @Tags         Media
// @Security     BearerAuth
// @Param        id path int true "Asset ID"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets/{id} [get]
func GetMediaAssetHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}

	resp, code, err := businessv1.GetMediaAsset(c.Request.Context(), userIDStr.(string), id)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "OK")
}

// UpdateMediaAssetHandler godoc
// @Summary      Update media asset title or status
// @Tags         Media
// @Security     BearerAuth
// @Param        id   path int true "Asset ID"
// @Param        body body modelsv1.UpdateMediaAssetRequest true "Fields to update"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets/{id} [put]
func UpdateMediaAssetHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req modelsv1.UpdateMediaAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	code, err := businessv1.UpdateMediaAsset(c.Request.Context(), userIDStr.(string), id, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"id": id}, "Updated")
}

// DeleteMediaAssetHandler godoc
// @Summary      Delete a media asset and all its versions
// @Tags         Media
// @Security     BearerAuth
// @Param        id path int true "Asset ID"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets/{id} [delete]
func DeleteMediaAssetHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid id")
		return
	}
	code, err := businessv1.DeleteMediaAsset(c.Request.Context(), userIDStr.(string), id)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"id": id}, "Deleted")
}

// AddVersionHandler godoc
// @Summary      Add a new version to a media asset
// @Tags         Media
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Param        id    path     int    true  "Asset ID"
// @Param        file  formData file   true  "New version file"
// @Param        label formData string false "Version label"
// @Param        notes formData string false "Notes about this version"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets/{id}/versions [post]
func AddVersionHandler(c *gin.Context) {
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

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	upload, code, err := businessv1.UploadMedia(
		c.Request.Context(), userIDStr.(string),
		header.Filename, contentType, header.Size, file,
	)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}

	req := &modelsv1.AddVersionRequest{
		Label: c.PostForm("label"),
		Notes: c.PostForm("notes"),
	}
	resp, code, err := businessv1.AddVersion(c.Request.Context(), userIDStr.(string), assetID, req, upload)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Version added")
}

// DeleteVersionHandler godoc
// @Summary      Delete a specific version
// @Tags         Media
// @Security     BearerAuth
// @Param        id  path int true "Asset ID"
// @Param        vid path int true "Version ID"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets/{id}/versions/{vid} [delete]
func DeleteVersionHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	assetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid asset id")
		return
	}
	versionID, err := strconv.ParseInt(c.Param("vid"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid version id")
		return
	}
	code, err := businessv1.DeleteVersion(c.Request.Context(), userIDStr.(string), assetID, versionID)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"id": versionID}, "Version deleted")
}

// SetCurrentVersionHandler godoc
// @Summary      Set the active version for a media asset
// @Tags         Media
// @Security     BearerAuth
// @Param        id   path int true "Asset ID"
// @Param        body body modelsv1.SetCurrentVersionRequest true "Version ID"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets/{id}/current-version [put]
func SetCurrentVersionHandler(c *gin.Context) {
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
	var req modelsv1.SetCurrentVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	code, err := businessv1.SetCurrentVersion(c.Request.Context(), userIDStr.(string), assetID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"asset_id": assetID, "version_id": req.VersionID}, "Current version updated")
}

// PublishMediaAssetHandler godoc
// @Summary      Publish current version to selected channels
// @Tags         Media
// @Security     BearerAuth
// @Param        id   path int true "Asset ID"
// @Param        body body modelsv1.PublishMediaAssetRequest true "Channels and caption"
// @Success      200 {object} map[string]interface{}
// @Router       /media-assets/{id}/publish [post]
func PublishMediaAssetHandler(c *gin.Context) {
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
	var req modelsv1.PublishMediaAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.PublishMediaAsset(c.Request.Context(), userIDStr.(string), assetID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Published")
}
