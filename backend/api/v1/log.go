package apiv1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

func GetLogByIDHandler(c *gin.Context) {
	logID := c.Param("log_id")
	if logID == "" {
		utils.Logger.Info(c.Request.Context(), "Log ID is required")
		utils.SendError(c, http.StatusBadRequest, "Log ID is required")
		return
	}

	log, err := businessv1.ReadLogsByLogID(c.Request.Context(), logID)
	if err != nil {
		utils.Logger.Info(c.Request.Context(), "Error reading log by ID: ", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Logger.Info(c.Request.Context(), "Read log by ID successfully: ", log)

	utils.SendSuccess(c, log)
}

func GetLogsByDate(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		utils.SendError(c, http.StatusBadRequest, "Date is required")
		return
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid date format")
		return
	}

	logs, total, err := businessv1.ReadLogsByDate(date)
	if err != nil {
		utils.Logger.Error(c.Request.Context(), "Failed to read logs: %v", err)
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, gin.H{"logs": logs, "total": total})
}