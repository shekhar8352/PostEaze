package apiv1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

// ---------------------------------------------------------------------------
// Studio
// ---------------------------------------------------------------------------

// EnsureStudioHandler godoc
// @Summary      Get or create the user's team Studio
// @Tags         Studio
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Router       /studios/default [get]
func EnsureStudioHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	resp, code, err := businessv1.EnsureStudioForUser(c.Request.Context(), userIDStr.(string))
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Studio ready")
}

// UpdateStudioHandler godoc
// @Summary      Rename a studio or change its piece label
// @Tags         Studio
// @Security     BearerAuth
// @Param        id path int true "Studio ID"
// @Param        body body modelsv1.UpdateStudioRequest true "Update"
// @Success      200 {object} map[string]interface{}
// @Router       /studios/{id} [put]
func UpdateStudioHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid studio id")
		return
	}
	var req modelsv1.UpdateStudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.UpdateStudio(c.Request.Context(), userIDStr.(string), id, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Studio updated")
}

// GetStudioBoardHandler godoc
// @Summary      Get the full board (studio + phases + pieces)
// @Tags         Studio
// @Security     BearerAuth
// @Param        id path int true "Studio ID"
// @Success      200 {object} map[string]interface{}
// @Router       /studios/{id}/board [get]
func GetStudioBoardHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid studio id")
		return
	}
	resp, code, err := businessv1.GetStudioBoard(c.Request.Context(), userIDStr.(string), id)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Board")
}

// ---------------------------------------------------------------------------
// Phases
// ---------------------------------------------------------------------------

// ListPhasesHandler godoc
// @Summary      List phases for a studio
// @Tags         Studio
// @Security     BearerAuth
// @Param        id path int true "Studio ID"
// @Success      200 {object} map[string]interface{}
// @Router       /studios/{id}/phases [get]
func ListPhasesHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	studioID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid studio id")
		return
	}
	resp, code, err := businessv1.ListPhases(c.Request.Context(), userIDStr.(string), studioID)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Phases")
}

// CreatePhaseHandler godoc
// @Summary      Add a new phase to a studio
// @Tags         Studio
// @Security     BearerAuth
// @Param        id path int true "Studio ID"
// @Param        body body modelsv1.CreatePhaseRequest true "Phase"
// @Success      200 {object} map[string]interface{}
// @Router       /studios/{id}/phases [post]
func CreatePhaseHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	studioID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid studio id")
		return
	}
	var req modelsv1.CreatePhaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.CreatePhase(c.Request.Context(), userIDStr.(string), studioID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Phase created")
}

// ReorderPhasesHandler godoc
// @Summary      Reorder phases inside a studio
// @Tags         Studio
// @Security     BearerAuth
// @Param        id path int true "Studio ID"
// @Param        body body modelsv1.ReorderPhasesRequest true "Ordered phase IDs"
// @Success      200 {object} map[string]interface{}
// @Router       /studios/{id}/phases/reorder [post]
func ReorderPhasesHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	studioID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid studio id")
		return
	}
	var req modelsv1.ReorderPhasesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if code, err := businessv1.ReorderPhases(c.Request.Context(), userIDStr.(string), studioID, &req); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Phases reordered")
}

// UpdatePhaseHandler godoc
// @Summary      Update a phase's metadata
// @Tags         Studio
// @Security     BearerAuth
// @Param        phaseId path int true "Phase ID"
// @Param        body body modelsv1.UpdatePhaseRequest true "Update"
// @Success      200 {object} map[string]interface{}
// @Router       /phases/{phaseId} [put]
func UpdatePhaseHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	phaseID, err := strconv.ParseInt(c.Param("phaseId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid phase id")
		return
	}
	var req modelsv1.UpdatePhaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.UpdatePhase(c.Request.Context(), userIDStr.(string), phaseID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Phase updated")
}

// DeletePhaseHandler godoc
// @Summary      Delete a phase (fails if it has active pieces)
// @Tags         Studio
// @Security     BearerAuth
// @Param        phaseId path int true "Phase ID"
// @Success      200 {object} map[string]interface{}
// @Router       /phases/{phaseId} [delete]
func DeletePhaseHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	phaseID, err := strconv.ParseInt(c.Param("phaseId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid phase id")
		return
	}
	if code, err := businessv1.DeletePhase(c.Request.Context(), userIDStr.(string), phaseID); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Phase deleted")
}

// ---------------------------------------------------------------------------
// Pieces
// ---------------------------------------------------------------------------

// CreatePieceHandler godoc
// @Summary      Create a piece in a studio
// @Tags         Studio
// @Security     BearerAuth
// @Param        id path int true "Studio ID"
// @Param        body body modelsv1.CreatePieceRequest true "Piece"
// @Success      200 {object} map[string]interface{}
// @Router       /studios/{id}/pieces [post]
func CreatePieceHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	studioID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid studio id")
		return
	}
	var req modelsv1.CreatePieceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.CreatePiece(c.Request.Context(), userIDStr.(string), studioID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Piece created")
}

// GetPieceHandler godoc
// @Summary      Get a piece with linked assets, scheduled posts, activities, and comments
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId} [get]
func GetPieceHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	resp, code, err := businessv1.GetPieceDetail(c.Request.Context(), userIDStr.(string), pieceID)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Piece")
}

// UpdatePieceHandler godoc
// @Summary      Update a piece's editable fields
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        body body modelsv1.UpdatePieceRequest true "Update"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId} [put]
func UpdatePieceHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	var req modelsv1.UpdatePieceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.UpdatePiece(c.Request.Context(), userIDStr.(string), pieceID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Piece updated")
}

// MovePieceHandler godoc
// @Summary      Move a piece between/inside phases
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        body body modelsv1.MovePieceRequest true "Move"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/move [post]
func MovePieceHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	var req modelsv1.MovePieceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.MovePiece(c.Request.Context(), userIDStr.(string), pieceID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Piece moved")
}

// SetPieceStatusHandler godoc
// @Summary      Archive or restore a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        body body modelsv1.SetPieceStatusRequest true "Status"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/status [put]
func SetPieceStatusHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	var req modelsv1.SetPieceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if code, err := businessv1.SetPieceStatus(c.Request.Context(), userIDStr.(string), pieceID, req.Status); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Piece status updated")
}

// DeletePieceHandler godoc
// @Summary      Permanently delete a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId} [delete]
func DeletePieceHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	if code, err := businessv1.DeletePiece(c.Request.Context(), userIDStr.(string), pieceID); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Piece deleted")
}

// LinkPieceAssetHandler godoc
// @Summary      Link a media asset to a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        body body modelsv1.LinkAssetRequest true "Asset + role"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/assets [post]
func LinkPieceAssetHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	var req modelsv1.LinkAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if code, err := businessv1.LinkPieceAsset(c.Request.Context(), userIDStr.(string), pieceID, &req); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Asset linked")
}

// UnlinkPieceAssetHandler godoc
// @Summary      Unlink a media asset from a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        assetId path int true "Media Asset ID"
// @Param        role query string false "Role (defaults to attachment)"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/assets/{assetId} [delete]
func UnlinkPieceAssetHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	assetID, err := strconv.ParseInt(c.Param("assetId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid asset id")
		return
	}
	role := c.Query("role")
	if role == "" {
		role = "attachment"
	}
	if code, err := businessv1.UnlinkPieceAsset(c.Request.Context(), userIDStr.(string), pieceID, assetID, role); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Asset unlinked")
}

// LinkPieceScheduledPostHandler godoc
// @Summary      Link a scheduled post to a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        body body modelsv1.LinkScheduledPostRequest true "Scheduled post + role"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/scheduled-posts [post]
func LinkPieceScheduledPostHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	var req modelsv1.LinkScheduledPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	if code, err := businessv1.LinkPieceScheduledPost(c.Request.Context(), userIDStr.(string), pieceID, &req); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Scheduled post linked")
}

// UnlinkPieceScheduledPostHandler godoc
// @Summary      Unlink a scheduled post from a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        scheduledPostId path int true "Scheduled post ID"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/scheduled-posts/{scheduledPostId} [delete]
func UnlinkPieceScheduledPostHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	spID, err := strconv.ParseInt(c.Param("scheduledPostId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid scheduled post id")
		return
	}
	if code, err := businessv1.UnlinkPieceScheduledPost(c.Request.Context(), userIDStr.(string), pieceID, spID); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Scheduled post unlinked")
}

// ListPieceActivitiesHandler godoc
// @Summary      List recent activities for a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        limit query int false "Max entries (default 50)"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/activities [get]
func ListPieceActivitiesHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	limit := 50
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	resp, code, err := businessv1.ListPieceActivities(c.Request.Context(), userIDStr.(string), pieceID, limit)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Activities")
}

// CreatePieceCommentHandler godoc
// @Summary      Add a comment to a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        body body modelsv1.CreatePieceCommentRequest true "Comment"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/comments [post]
func CreatePieceCommentHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	var req modelsv1.CreatePieceCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	resp, code, err := businessv1.CreatePieceComment(c.Request.Context(), userIDStr.(string), pieceID, &req)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Comment added")
}

// ListPieceCommentsHandler godoc
// @Summary      List comments for a piece
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/comments [get]
func ListPieceCommentsHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	resp, code, err := businessv1.ListPieceComments(c.Request.Context(), userIDStr.(string), pieceID)
	if err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, resp, "Comments")
}

// DeletePieceCommentHandler godoc
// @Summary      Delete a comment
// @Tags         Studio
// @Security     BearerAuth
// @Param        pieceId path int true "Piece ID"
// @Param        commentId path int true "Comment ID"
// @Success      200 {object} map[string]interface{}
// @Router       /pieces/{pieceId}/comments/{commentId} [delete]
func DeletePieceCommentHandler(c *gin.Context) {
	userIDStr, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User ID not found in token")
		return
	}
	pieceID, err := strconv.ParseInt(c.Param("pieceId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid piece id")
		return
	}
	commentID, err := strconv.ParseInt(c.Param("commentId"), 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "invalid comment id")
		return
	}
	if code, err := businessv1.DeletePieceComment(c.Request.Context(), userIDStr.(string), pieceID, commentID); err != nil {
		utils.SendError(c, code, err.Error())
		return
	}
	utils.SendSuccess(c, gin.H{"ok": true}, "Comment deleted")
}
