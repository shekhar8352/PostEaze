package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

// ToUTCDate truncates t to UTC midnight (calendar date boundary).
func ToUTCDate(t time.Time) time.Time {
	t = t.UTC()
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// ParseDateRange parses start_date and end_date query parameters (YYYY-MM-DD in UTC).
// Defaults to the last 7 calendar days ending today if not provided.
func ParseDateRange(c *gin.Context) (time.Time, time.Time) {
	endDate := ToUTCDate(time.Now())
	startDate := endDate.AddDate(0, 0, -6) // inclusive 7-day window including endDate

	if startStr := c.Query("start_date"); startStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", startStr, time.UTC); err == nil {
			startDate = ToUTCDate(parsed)
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", endStr, time.UTC); err == nil {
			endDate = ToUTCDate(parsed)
		}
	}

	if endDate.Before(startDate) {
		startDate, endDate = endDate, startDate
	}

	return startDate, endDate
}

// PreviousPeriodInclusive returns the immediately preceding window with the same number of
// inclusive calendar days as [startDate, endDate] (both should be UTC date-truncated).
func PreviousPeriodInclusive(startDate, endDate time.Time) (prevStart, prevEnd time.Time) {
	startDate = ToUTCDate(startDate)
	endDate = ToUTCDate(endDate)
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days < 1 {
		days = 1
	}
	prevEnd = startDate.AddDate(0, 0, -1)
	prevStart = prevEnd.AddDate(0, 0, -(days - 1))
	return ToUTCDate(prevStart), ToUTCDate(prevEnd)
}
