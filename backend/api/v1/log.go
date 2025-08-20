package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	businessv1 "github.com/shekhar8352/PostEaze/business/v1"
	"github.com/shekhar8352/PostEaze/utils"
)

func GetLogByIDHandler(c *gin.Context) {
	logID := c.Param("log_id")
	
	// Enhanced log ID validation
	if err := utils.ValidateLogID(logID); err != nil {
		errorMsg := utils.GetValidationErrorMessage(err)
		utils.Logger.Info(c.Request.Context(), "Log ID validation failed: %s", errorMsg)
		utils.SendLogAPIError(c, http.StatusBadRequest, errorMsg, utils.ErrorTypeInvalidInput)
		return
	}

	logs, err := businessv1.ReadLogsByLogID(c.Request.Context(), logID)
	if err != nil {
		// Check if it's a log file read error (directory issues, etc.)
		if utils.IsLogFileReadError(err) {
			utils.Logger.Error(c.Request.Context(), "Failed to read logs by ID '%s': %v", logID, err)
			utils.SendLogAPIError(c, http.StatusInternalServerError, "Failed to read log files", utils.ErrorTypeInternal)
			return
		}
		
		// Check if it's a log file not found error
		if utils.IsLogFileNotFoundError(err) {
			utils.Logger.Info(c.Request.Context(), "No logs found for ID '%s': %v", logID, err)
			utils.SendLogAPIError(c, http.StatusNotFound, "No logs found for the specified log ID", utils.ErrorTypeNotFound)
			return
		}
		
		// Fallback for other errors
		utils.Logger.Error(c.Request.Context(), "Unexpected error reading logs by ID '%s': %v", logID, err)
		utils.SendLogAPIError(c, http.StatusInternalServerError, "Internal server error", utils.ErrorTypeInternal)
		return
	}

	utils.Logger.Info(c.Request.Context(), "Successfully retrieved logs for ID '%s', found %d entries", logID, len(logs))

	utils.SendLogAPISuccess(c, logs, "Logs retrieved successfully")
}

func GetLogsByDate(c *gin.Context) {
	date := c.Param("date")
	
	// Enhanced date validation
	if err := utils.ValidateDate(date); err != nil {
		errorMsg := utils.GetValidationErrorMessage(err)
		utils.Logger.Info(c.Request.Context(), "Date validation failed: %s", errorMsg)
		utils.SendLogAPIError(c, http.StatusBadRequest, errorMsg, utils.ErrorTypeInvalidInput)
		return
	}

	logs, total, err := businessv1.ReadLogsByDate(date)
	if err != nil {
		// Check if it's a log file not found error
		if utils.IsLogFileNotFoundError(err) {
			utils.Logger.Info(c.Request.Context(), "No logs found for date '%s': %v", date, err)
			utils.SendLogAPIError(c, http.StatusNotFound, "No logs found for the specified date", utils.ErrorTypeNotFound)
			return
		}
		
		// Check if it's a log file read error
		if utils.IsLogFileReadError(err) {
			utils.Logger.Error(c.Request.Context(), "Failed to read logs for date '%s': %v", date, err)
			utils.SendLogAPIError(c, http.StatusInternalServerError, "Failed to read log file", utils.ErrorTypeInternal)
			return
		}
		
		// Fallback for other errors
		utils.Logger.Error(c.Request.Context(), "Unexpected error reading logs for date '%s': %v", date, err)
		utils.SendLogAPIError(c, http.StatusInternalServerError, "Internal server error", utils.ErrorTypeInternal)
		return
	}

	utils.Logger.Info(c.Request.Context(), "Successfully retrieved logs for date '%s', found %d entries", date, total)

	utils.SendLogAPISuccess(c, gin.H{"logs": logs, "total": total}, "Logs retrieved successfully")
}