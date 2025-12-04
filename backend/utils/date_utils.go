package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

// ParseDateRange parses start_date and end_date query parameters
// Returns start and end dates, defaulting to last 7 days if not provided
func ParseDateRange(c *gin.Context) (time.Time, time.Time) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7) // Default: last 7 days

	if startStr := c.Query("start_date"); startStr != "" {
		if parsed, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = parsed
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		if parsed, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = parsed
		}
	}

	return startDate, endDate
}
