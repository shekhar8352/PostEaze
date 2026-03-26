package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/instagram"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

const postInsightsWorkerCount = 5

// Profile insights: fetch up to this many days (Instagram often caps a single request; we chunk).
const (
	profileInsightsLookbackDays = 90
	profileInsightsChunkDays    = 30
)

// Split metrics into batches to stay within Graph API limits and reduce all-or-nothing failures.
var profileInsightsMetricBatches = [][]string{
	{"reach", "profile_views", "follower_count", "website_clicks"},
	{"accounts_engaged", "total_interactions", "views"},
}

// Lifetime audience breakdowns (Graph API period=lifetime). Batched so one failing metric does not block others.
var audienceInsightMetricBatches = [][]string{
	{"audience_city", "audience_country"},
	{"audience_gender_age", "audience_locale"},
}

// HandleSyncInstagramAnalyticsTask syncs Instagram analytics for all active channels
func HandleSyncInstagramAnalyticsTask(ctx context.Context, t *asynq.Task) error {
	utils.Logger.Info(ctx, "Starting Instagram analytics sync job")

	// Fetch all active Instagram channels
	channels, err := repositories.GetAllActiveInstagramChannels(ctx)
	if err != nil {
		utils.Logger.Error(ctx, "Failed to fetch Instagram channels: %v", err)
		return err
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Found %d active Instagram channels to sync analytics", len(channels)))

	// Process each channel
	for _, channel := range channels {
		if err := syncChannelAnalytics(ctx, channel); err != nil {
			// Log error and continue with next channel
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to sync analytics for channel %d: %v", channel.ID, err))
			continue
		}
	}

	utils.Logger.Info(ctx, "Instagram analytics sync job completed")
	return nil
}

func syncChannelAnalytics(ctx context.Context, channel entities.Channel) error {
	// Get the latest access token
	token, err := repositories.GetLatestTokenByChannelID(ctx, channel.ID)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Decrypt the access token
	decryptedToken, err := encryption.Decrypt(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt access token: %w", err)
	}

	// Get Instagram user ID from channel metadata
	var metadata map[string]interface{}
	if len(channel.Metadata) > 0 {
		if err := json.Unmarshal(channel.Metadata, &metadata); err != nil {
			return fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	igUserID, ok := metadata["id"].(string)
	if !ok || igUserID == "" {
		return fmt.Errorf("instagram user ID not found in metadata")
	}

	provider := instagram.NewInstagramProvider()

	// Sync profile analytics
	if err := syncProfileAnalytics(ctx, provider, channel.ID, igUserID, decryptedToken); err != nil {
		utils.Logger.Error(ctx, fmt.Sprintf("Failed to sync profile analytics for channel %d: %v", channel.ID, err))
	}

	if err := syncAudienceInsights(ctx, provider, channel.ID, igUserID, decryptedToken); err != nil {
		utils.Logger.Error(ctx, fmt.Sprintf("Failed to sync audience insights for channel %d: %v", channel.ID, err))
	}

	// Sync post analytics
	if err := syncPostsAnalytics(ctx, provider, channel.ID, decryptedToken); err != nil {
		utils.Logger.Error(ctx, fmt.Sprintf("Failed to sync posts analytics for channel %d: %v", channel.ID, err))
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Successfully synced analytics for channel %d", channel.ID))
	return nil
}

func syncProfileAnalytics(ctx context.Context, provider instagram.InstagramProvider, channelID int64, igUserID string, accessToken string) error {
	// Daily profile insights: chunk last N days (API window limits) and merge metric batches.
	metricsByDate := make(map[string]*entities.InstagramProfileAnalytics)

	now := time.Now().UTC()
	remaining := profileInsightsLookbackDays
	cursor := now
	for remaining > 0 {
		span := profileInsightsChunkDays
		if span > remaining {
			span = remaining
		}
		until := cursor.Unix()
		since := cursor.AddDate(0, 0, -span).Unix()
		cursor = cursor.AddDate(0, 0, -span)
		remaining -= span

		for _, batch := range profileInsightsMetricBatches {
			resp, err := provider.GetProfileInsights(accessToken, igUserID, batch, since, until)
			if err != nil {
				utils.Logger.Warn(ctx, fmt.Sprintf("Profile insights batch failed for channel %d (since=%d until=%d metrics=%v): %v", channelID, since, until, batch, err))
				continue
			}
			mergeProfileInsightsIntoByDate(ctx, metricsByDate, channelID, resp)
		}
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Upserting profile analytics for channel %d: %d day rows", channelID, len(metricsByDate)))
	for _, analytics := range metricsByDate {
		analytics.Raw = nil
		if err := repositories.UpsertInstagramProfileAnalytics(ctx, analytics); err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to upsert profile analytics: %v", err))
		}
	}

	return nil
}

func syncAudienceInsights(ctx context.Context, provider instagram.InstagramProvider, channelID int64, igUserID, accessToken string) error {
	var merged []instagram.InsightData
	for _, batch := range audienceInsightMetricBatches {
		resp, err := provider.GetAudienceInsights(accessToken, igUserID, batch)
		if err != nil {
			var gerr *instagram.GraphAPIError
			if errors.As(err, &gerr) {
				utils.Logger.Warn(ctx, fmt.Sprintf("Audience insights batch skipped for channel %d metrics=%v: %s", channelID, batch, gerr.Message))
			} else {
				utils.Logger.Warn(ctx, fmt.Sprintf("Audience insights batch failed for channel %d metrics=%v: %v", channelID, batch, err))
			}
			continue
		}
		merged = append(merged, resp.Data...)
	}
	if len(merged) == 0 {
		utils.Logger.Info(ctx, fmt.Sprintf("No audience demographic data stored for channel %d (ineligible or unavailable from Instagram)", channelID))
		return nil
	}
	wrapped := struct {
		Data []instagram.InsightData `json:"data"`
	}{Data: merged}
	raw, err := json.Marshal(wrapped)
	if err != nil {
		return fmt.Errorf("marshal audience insights: %w", err)
	}
	snap := &entities.InstagramAudienceSnapshot{
		ChannelID:    channelID,
		SnapshotDate: utils.ToUTCDate(time.Now().UTC()),
		Raw:          raw,
	}
	if err := repositories.UpsertInstagramAudienceSnapshot(ctx, snap); err != nil {
		return fmt.Errorf("upsert audience snapshot: %w", err)
	}
	utils.Logger.Info(ctx, fmt.Sprintf("Stored audience snapshot for channel %d (%d insight series)", channelID, len(merged)))
	return nil
}

func mergeProfileInsightsIntoByDate(ctx context.Context, metricsByDate map[string]*entities.InstagramProfileAnalytics, channelID int64, insightsResp *instagram.InsightsResponse) {
	for _, insight := range insightsResp.Data {
		for _, value := range insight.Values {
			var date time.Time
			var err error
			if value.EndTime != "" {
				date, err = time.Parse("2006-01-02T15:04:05-0700", value.EndTime)
				if err != nil {
					utils.Logger.Error(ctx, fmt.Sprintf("Failed to parse end_time: %v", err))
					continue
				}
			} else {
				date = time.Now().UTC()
			}

			dateKey := date.UTC().Truncate(24 * time.Hour).Format("2006-01-02")
			if metricsByDate[dateKey] == nil {
				metricsByDate[dateKey] = &entities.InstagramProfileAnalytics{
					ChannelID: channelID,
					Date:      date.UTC().Truncate(24 * time.Hour),
				}
			}
			analytics := metricsByDate[dateKey]

			v64, ok := instagram.ParseScalarInsightValue(value.Value)
			if !ok {
				continue
			}
			intValue := int(v64)

			switch insight.Name {
			case "reach":
				analytics.Reach = &intValue
			case "profile_views":
				analytics.ProfileViews = &intValue
			case "follower_count":
				analytics.FollowerCount = &intValue
			case "website_clicks":
				analytics.WebsiteClicks = &intValue
			case "accounts_engaged":
				analytics.AccountsEngaged = &intValue
			case "total_interactions":
				analytics.TotalInteractions = &intValue
			case "views":
				analytics.Views = &intValue
			}
		}
	}
}

func syncPostsAnalytics(ctx context.Context, provider instagram.InstagramProvider, channelID int64, accessToken string) error {
	// Get all posts for this channel (limit to recent 100)
	posts, err := repositories.GetPostsByChannel(ctx, channelID, 100)
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Syncing analytics for %d posts with %d workers", len(posts), postInsightsWorkerCount))

	jobs := make(chan entities.Post, len(posts))
	for _, p := range posts {
		jobs <- p
	}
	close(jobs)

	var wg sync.WaitGroup
	for w := 0; w < postInsightsWorkerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for post := range jobs {
				syncSinglePostInsights(ctx, provider, channelID, accessToken, post)
			}
		}()
	}
	wg.Wait()

	return nil
}

// mediaInsightMetricFallbacks returns ordered metric lists. Carousels and some media types reject
// "impressions" (and occasionally "plays") in a single batch; later slices omit problematic metrics.
func mediaInsightMetricFallbacks(postType *string) [][]string {
	pt := ""
	if postType != nil {
		pt = *postType
	}
	noImpression := []string{"reach", "likes", "comments", "saved", "shares", "total_interactions"}
	full := []string{"impressions", "reach", "likes", "comments", "saved", "shares", "total_interactions"}
	fullVideo := []string{"impressions", "reach", "likes", "comments", "saved", "shares", "plays", "total_interactions"}
	minimal := []string{"reach", "likes", "comments", "saved", "shares"}

	switch pt {
	case "carousel":
		return [][]string{
			noImpression,
			{"likes", "comments", "saved", "shares", "total_interactions"},
			{"likes", "comments", "saved"},
		}
	case "video", "reel":
		return [][]string{
			fullVideo,
			noImpression,
			{"reach", "likes", "comments", "saved", "shares", "plays"},
			minimal,
		}
	default:
		return [][]string{
			full,
			noImpression,
			minimal,
			{"likes", "comments"},
		}
	}
}

func fetchMediaInsightsWithFallback(provider instagram.InstagramProvider, accessToken, mediaID string, postType *string) (*instagram.InsightsResponse, error) {
	var lastErr error
	for _, metrics := range mediaInsightMetricFallbacks(postType) {
		resp, err := provider.GetMediaInsights(accessToken, mediaID, metrics)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		var gerr *instagram.GraphAPIError
		if errors.As(err, &gerr) {
			if gerr.IsInsightsUnavailableForever() {
				return nil, err
			}
			continue
		}
		return nil, err
	}
	return nil, lastErr
}

func syncSinglePostInsights(ctx context.Context, provider instagram.InstagramProvider, channelID int64, accessToken string, post entities.Post) {
	isStory := post.PostType != nil && *post.PostType == "story"

	if isStory && post.PublishedAt != nil && time.Since(*post.PublishedAt) > 24*time.Hour {
		return
	}

	var insightsResp *instagram.InsightsResponse
	var err error

	instagramPostID := repositories.GetInstagramPostID(&post)

	if isStory {
		if instagramPostID == "" {
			utils.Logger.Error(ctx, fmt.Sprintf("Skipping story analytics for post %d: instagram_post_id is nil", post.ID))
			return
		}
		insightsResp, err = provider.GetStoryInsights(accessToken, instagramPostID)
	} else {
		if instagramPostID == "" {
			utils.Logger.Error(ctx, fmt.Sprintf("Skipping post analytics for post %d: instagram_post_id is nil", post.ID))
			return
		}
		insightsResp, err = fetchMediaInsightsWithFallback(provider, accessToken, instagramPostID, post.PostType)
	}

	if err != nil {
		var gerr *instagram.GraphAPIError
		if errors.As(err, &gerr) && gerr.IsInsightsUnavailableForever() {
			utils.Logger.Warn(ctx, fmt.Sprintf("Skipping insights for post %s (not available from Instagram): %s", instagramPostID, gerr.Message))
			return
		}
		utils.Logger.Warn(ctx, fmt.Sprintf("Failed to fetch insights for post %s after metric fallbacks: %v", instagramPostID, err))
		return
	}

	if isStory {
		if err := processStoryInsights(ctx, channelID, post, insightsResp); err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to process story insights: %v", err))
		}
	} else {
		if err := processPostInsights(ctx, channelID, post, insightsResp); err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to process post insights: %v", err))
		}
	}
}

// getMetricValue is a helper function to extract a specific metric value from InsightsResponse.
func getMetricValue(insightsResp *instagram.InsightsResponse, metricName string) *int {
	for _, insight := range insightsResp.Data {
		if insight.Name == metricName {
			if len(insight.Values) > 0 {
				if v64, ok := instagram.ParseScalarInsightValue(insight.Values[0].Value); ok {
					i := int(v64)
					return &i
				}
			}
		}
	}
	return nil
}

func processPostInsights(ctx context.Context, channelID int64, post entities.Post, insightsResp *instagram.InsightsResponse) error {
	rawJSON, _ := json.Marshal(insightsResp)

	analytics := &entities.InstagramPostAnalytics{
		ChannelID:     channelID,
		PostID:        post.ID,
		Date:          utils.ToUTCDate(time.Now()),
		Impressions:   getMetricValue(insightsResp, "impressions"),
		Reach:         getMetricValue(insightsResp, "reach"),
		Likes:         getMetricValue(insightsResp, "likes"),
		Comments:      getMetricValue(insightsResp, "comments"),
		Saves:         getMetricValue(insightsResp, "saved"),
		Shares:        getMetricValue(insightsResp, "shares"),
		VideoViews:    getMetricValue(insightsResp, "video_views"),
		ProfileVisits: getMetricValue(insightsResp, "profile_visits"),
		Follows:       getMetricValue(insightsResp, "follows"),

		Views:             getMetricValue(insightsResp, "views"),
		Plays:             getMetricValue(insightsResp, "plays"),
		TotalInteractions: getMetricValue(insightsResp, "total_interactions"),

		Raw: rawJSON,
	}

	// Calculate engagement rate
	reach := 0
	if analytics.Reach != nil {
		reach = *analytics.Reach
	} else if analytics.Impressions != nil {
		reach = *analytics.Impressions // fallback
	}

	totalInteractions := 0
	if analytics.TotalInteractions != nil {
		totalInteractions = *analytics.TotalInteractions
	} else {
		// Fallback calculation
		likes := 0
		comments := 0
		saves := 0
		shares := 0
		if analytics.Likes != nil {
			likes = *analytics.Likes
		}
		if analytics.Comments != nil {
			comments = *analytics.Comments
		}
		if analytics.Saves != nil {
			saves = *analytics.Saves
		}
		if analytics.Shares != nil {
			shares = *analytics.Shares
		}
		totalInteractions = likes + comments + saves + shares
		if totalInteractions > 0 {
			analytics.TotalInteractions = &totalInteractions
		}
	}

	if reach > 0 && totalInteractions > 0 {
		rate := (float64(totalInteractions) / float64(reach)) * 100
		analytics.EngagementRate = &rate
	}

	return repositories.UpsertInstagramPostAnalytics(ctx, analytics)
}

func processStoryInsights(ctx context.Context, channelID int64, post entities.Post, insightsResp *instagram.InsightsResponse) error {
	rawJSON, _ := json.Marshal(insightsResp)

	analytics := &entities.InstagramStoryAnalytics{
		ChannelID:    channelID,
		PostID:       post.ID,
		Date:         utils.ToUTCDate(time.Now()),
		Impressions:  getMetricValue(insightsResp, "impressions"),
		Reach:        getMetricValue(insightsResp, "reach"),
		Exits:        getMetricValue(insightsResp, "exits"),
		Replies:      getMetricValue(insightsResp, "replies"),
		TapsForward:  getMetricValue(insightsResp, "taps_forward"),
		TapsBackward: getMetricValue(insightsResp, "taps_back"),
		Raw:          rawJSON,
	}

	return repositories.UpsertInstagramStoryAnalytics(ctx, analytics)
}
