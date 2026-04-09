package businessv1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	modelsv1 "github.com/shekhar8352/PostEaze/models/v1"
	"github.com/shekhar8352/PostEaze/tasks"
)

// SyncMetaAnalyticsForUser runs Instagram and/or Facebook analytics sync for channels owned by the user.
func SyncMetaAnalyticsForUser(ctx context.Context, userID string, channelIDs []int64) (*modelsv1.SyncMetaAnalyticsResponse, error) {
	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	var candidates []int64
	if len(channelIDs) == 0 {
		channels, err := repositories.GetChannelsByUserID(ctx, userID, "")
		if err != nil {
			return nil, err
		}
		for _, ch := range channels {
			if !ch.IsActive {
				continue
			}
			if ch.Provider != "instagram" && ch.Provider != "facebook" {
				continue
			}
			candidates = append(candidates, ch.ID)
		}
	} else {
		candidates = append(candidates, channelIDs...)
	}

	var results []modelsv1.SyncMetaAnalyticsResult
	for _, cid := range candidates {
		ch, err := repositories.GetChannelByID(ctx, cid)
		if err != nil {
			msg := "channel lookup failed"
			if errors.Is(err, sql.ErrNoRows) {
				msg = "channel not found"
			}
			results = append(results, modelsv1.SyncMetaAnalyticsResult{
				ChannelID: cid,
				Status:    "error",
				Message:   msg,
			})
			continue
		}
		if ch.OwnerUserID != ownerID {
			results = append(results, modelsv1.SyncMetaAnalyticsResult{
				ChannelID: cid,
				Provider:  ch.Provider,
				Status:    "error",
				Message:   "forbidden",
			})
			continue
		}
		if ch.Provider != "instagram" && ch.Provider != "facebook" {
			results = append(results, modelsv1.SyncMetaAnalyticsResult{
				ChannelID: cid,
				Provider:  ch.Provider,
				Status:    "error",
				Message:   "unsupported provider",
			})
			continue
		}

		var syncErr error
		switch ch.Provider {
		case "instagram":
			syncErr = tasks.SyncInstagramChannelAnalytics(ctx, *ch)
		case "facebook":
			syncErr = tasks.SyncFacebookChannelAnalytics(ctx, *ch)
		}

		if syncErr != nil {
			results = append(results, modelsv1.SyncMetaAnalyticsResult{
				ChannelID: cid,
				Provider:  ch.Provider,
				Status:    "error",
				Message:   syncErr.Error(),
			})
			continue
		}
		results = append(results, modelsv1.SyncMetaAnalyticsResult{
			ChannelID: cid,
			Provider:  ch.Provider,
			Status:    "ok",
		})
	}

	return &modelsv1.SyncMetaAnalyticsResponse{Results: results}, nil
}
