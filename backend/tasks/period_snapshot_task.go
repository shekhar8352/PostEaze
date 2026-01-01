package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/utils"
)

// HandlePeriodSnapshotTask generates analytics snapshots for all active channels
func HandlePeriodSnapshotTask(ctx context.Context, t *asynq.Task) error {
	utils.Logger.Info(ctx, "Starting period snapshot generation task")

	// Get all active channels
	// Note: We should probably get ALL channels, not just Instagram, if we support others later.
	// For now, we reuse GetAllActiveInstagramChannels or generic GetAllChannels if available.
	// But GetAllActiveInstagramChannels is available in repository.
	channels, err := repositories.GetAllActiveInstagramChannels(ctx)
	if err != nil {
		utils.Logger.Error(ctx, "Failed to fetch channels: "+err.Error())
		return err
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	// Dates for Daily Snapshot (Yesterday)
	dailyStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	dailyEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, time.UTC)

	for _, channel := range channels {
		// 1. Generate Daily Snapshot
		if err := generateSnapshot(ctx, "daily", channel.ID, dailyStart, dailyEnd); err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to generate daily snapshot for channel %d: %v", channel.ID, err))
		}

		// 2. Generate Weekly Snapshot (if today is Monday)
		if now.Weekday() == time.Monday {
			// Last week (Monday to Sunday)
			lastWeekStart := dailyStart.AddDate(0, 0, -6) // Yesterday (Sunday) - 6 days = Previous Monday
			// Verify: If today is Monday, Yesterday is Sunday. Sunday - 6 days = Monday. correct.
			lastWeekEnd := dailyEnd // Yesterday (Sunday) end.

			if err := generateSnapshot(ctx, "weekly", channel.ID, lastWeekStart, lastWeekEnd); err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to generate weekly snapshot for channel %d: %v", channel.ID, err))
			}
		}

		// 3. Generate Monthly Snapshot (if today is 1st of month)
		if now.Day() == 1 {
			// Last month
			// Yesterday was last day of last month.
			lastMonthEnd := dailyEnd
			lastMonthStart := time.Date(yesterday.Year(), yesterday.Month(), 1, 0, 0, 0, 0, time.UTC)

			if err := generateSnapshot(ctx, "monthly", channel.ID, lastMonthStart, lastMonthEnd); err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to generate monthly snapshot for channel %d: %v", channel.ID, err))
			}
		}
	}

	utils.Logger.Info(ctx, "Period snapshot generation task completed")
	return nil
}

func generateSnapshot(ctx context.Context, periodType string, channelID int64, start, end time.Time) error {
	// 1. Fetch Aggregated Metrics
	agg, err := repositories.GetAggregatedProfileAnalytics(ctx, channelID, start, end)
	if err != nil {
		return fmt.Errorf("fetching analytics: %w", err)
	}

	// 2. Marshal metrics to JSON
	metricsJSON, err := json.Marshal(agg)
	if err != nil {
		return fmt.Errorf("marshaling metrics: %w", err)
	}

	// 3. Create/Upsert Snapshot
	snapshot := &entities.AnalyticsSnapshot{
		EntityType: "channel",
		EntityID:   channelID,
		PeriodType: periodType,
		StartDate:  start,
		EndDate:    end,
		Metrics:    metricsJSON,
	}

	return repositories.UpsertAnalyticsSnapshot(ctx, snapshot)
}
