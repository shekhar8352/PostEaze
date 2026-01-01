package entities

import (
	"encoding/json"
	"time"
)

// AnalyticsSnapshot represents a generic analytics snapshot for any entity (channel, post)
type AnalyticsSnapshot struct {
	ID         int64           `db:"id"`
	EntityType string          `db:"entity_type"` // 'channel', 'post'
	EntityID   int64           `db:"entity_id"`
	PeriodType string          `db:"period_type"` // 'daily', 'weekly', 'monthly'
	StartDate  time.Time       `db:"start_date"`
	EndDate    time.Time       `db:"end_date"`
	Metrics    json.RawMessage `db:"metrics"` // Flexible JSON storage
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}
