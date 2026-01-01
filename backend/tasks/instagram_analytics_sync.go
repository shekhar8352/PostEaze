package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/instagram"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

// HandleSyncInstagramAnalyticsTask syncs Instagram analytics for all active channels
func HandleSyncInstagramAnalyticsTask(ctx context.Context, t *asynq.Task) error {
	utils.Logger.Info(ctx, "Starting Instagram analytics sync job")

	// Fetch all active Instagram channels
	channels, err := repositories.GetAllActiveInstagramChannels(ctx)
	if err != nil {
		utils.Logger.Error(ctx, "Failed to fetch Instagram channels: ", err)
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

	// Sync post analytics
	if err := syncPostsAnalytics(ctx, provider, channel.ID, decryptedToken); err != nil {
		utils.Logger.Error(ctx, fmt.Sprintf("Failed to sync posts analytics for channel %d: %v", channel.ID, err))
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Successfully synced analytics for channel %d", channel.ID))
	return nil
}

func syncProfileAnalytics(ctx context.Context, provider instagram.InstagramProvider, channelID int64, igUserID string, accessToken string) error {
	// Profile insights metrics (daily period)
	// Valid metrics from API: reach, follower_count, website_clicks, profile_views, online_followers, accounts_engaged, total_interactions
	// Note: impressions is NOT available for profile insights
	metrics := []string{"reach", "profile_views", "follower_count", "website_clicks"}

	// Calculate time range (last 7 days)
	now := time.Now()
	since := now.AddDate(0, 0, -7).Unix()
	until := now.Unix()

	insightsResp, err := provider.GetProfileInsights(accessToken, igUserID, metrics, since, until)
	if err != nil {
		return fmt.Errorf("failed to fetch profile insights: %w", err)
	}

	// Group metrics by date
	metricsByDate := make(map[string]*entities.InstagramProfileAnalytics)

	// Process insights data
	for _, insight := range insightsResp.Data {
		for _, value := range insight.Values {
			// Parse end_time (Instagram format: "2025-11-28T08:00:00+0000")
			var date time.Time
			if value.EndTime != "" {
				date, err = time.Parse("2006-01-02T15:04:05-0700", value.EndTime)
				if err != nil {
					utils.Logger.Error(ctx, fmt.Sprintf("Failed to parse end_time: %v", err))
					continue
				}
			} else {
				date = time.Now()
			}

			dateKey := date.Truncate(24 * time.Hour).Format("2006-01-02")

			// Get or create analytics record for this date
			if metricsByDate[dateKey] == nil {
				metricsByDate[dateKey] = &entities.InstagramProfileAnalytics{
					ChannelID: channelID,
					Date:      date.Truncate(24 * time.Hour),
				}
			}

			analytics := metricsByDate[dateKey]

			// Extract value
			var intValue int
			switch v := value.Value.(type) {
			case float64:
				intValue = int(v)
			case int:
				intValue = v
			default:
				continue
			}

			// Set the appropriate field based on metric name
			switch insight.Name {
			case "reach":
				analytics.Reach = &intValue
				utils.Logger.Info(ctx, fmt.Sprintf("Setting reach=%d for date %s", intValue, dateKey))
			case "profile_views":
				analytics.ProfileViews = &intValue
				utils.Logger.Info(ctx, fmt.Sprintf("Setting profile_views=%d for date %s", intValue, dateKey))
			case "follower_count":
				analytics.FollowerCount = &intValue
				utils.Logger.Info(ctx, fmt.Sprintf("Setting follower_count=%d for date %s", intValue, dateKey))
			case "website_clicks":
				analytics.WebsiteClicks = &intValue
				utils.Logger.Info(ctx, fmt.Sprintf("Setting website_clicks=%d for date %s", intValue, dateKey))
			}
		}
	}

	// Marshal raw data and upsert all analytics
	rawData, _ := json.Marshal(insightsResp)
	utils.Logger.Info(ctx, fmt.Sprintf("Upserting profile analytics for %d dates", len(metricsByDate)))
	for _, analytics := range metricsByDate {
		analytics.Raw = rawData
		utils.Logger.Info(ctx, fmt.Sprintf("Upserting analytics for date %s: reach=%v, profile_views=%v, follower_count=%v, website_clicks=%v",
			analytics.Date.Format("2006-01-02"),
			analytics.Reach,
			analytics.ProfileViews,
			analytics.FollowerCount,
			analytics.WebsiteClicks))
		if err := repositories.UpsertInstagramProfileAnalytics(ctx, analytics); err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to upsert profile analytics: %v", err))
		} else {
			utils.Logger.Info(ctx, fmt.Sprintf("Upserted profile analytics for date %s", analytics.Date.Format("2006-01-02")))
		}
	}

	return nil
}

func syncPostsAnalytics(ctx context.Context, provider instagram.InstagramProvider, channelID int64, accessToken string) error {
	// Get all posts for this channel (limit to recent 100)
	posts, err := repositories.GetPostsByChannel(ctx, channelID, 100)
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Syncing analytics for %d posts", len(posts)))

	for _, post := range posts {
		// Determine if it's a story (stories expire after 24 hours)
		isStory := post.PostType != nil && *post.PostType == "story"

		// Skip stories older than 24 hours
		if isStory && post.PublishedAt != nil && time.Since(*post.PublishedAt) > 24*time.Hour {
			continue
		}

		// Fetch insights based on post type
		var insightsResp *instagram.InsightsResponse
		var err error

		// Get Instagram post ID from JSONB
		instagramPostID := repositories.GetInstagramPostID(&post)

		if isStory {
			if instagramPostID == "" {
				utils.Logger.Error(ctx, fmt.Sprintf("Skipping story analytics for post %d: instagram_post_id is nil", post.ID))
				continue
			}
			insightsResp, err = provider.GetStoryInsights(accessToken, instagramPostID)
		} else {
			// Determine metrics based on post type
			// Valid metrics from API: impressions, reach, likes, comments, saved, shares, plays, total_interactions
			var metrics []string
			if post.PostType != nil && (*post.PostType == "video" || *post.PostType == "reel") {
				// Video/Reel metrics
				metrics = []string{"impressions", "reach", "likes", "comments", "saved", "shares", "plays", "total_interactions"}
			} else {
				// Image/Carousel metrics
				metrics = []string{"impressions", "reach", "likes", "comments", "saved", "shares", "total_interactions"}
			}
			if instagramPostID == "" {
				utils.Logger.Error(ctx, fmt.Sprintf("Skipping post analytics for post %d: instagram_post_id is nil", post.ID))
				continue
			}
			insightsResp, err = provider.GetMediaInsights(accessToken, instagramPostID, metrics)
		}

		if err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to fetch insights for post %s: %v", instagramPostID, err))
			continue
		}

		// Process insights
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

	return nil
}

// getMetricValue is a helper function to extract a specific metric value from InsightsResponse.
func getMetricValue(insightsResp *instagram.InsightsResponse, metricName string) *int {
	for _, insight := range insightsResp.Data {
		if insight.Name == metricName {
			if len(insight.Values) > 0 {
				var intValue int
				switch v := insight.Values[0].Value.(type) {
				case float64:
					intValue = int(v)
				case int:
					intValue = v
				default:
					return nil
				}
				return &intValue
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
		Date:          time.Now(),
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
		Date:         time.Now(),
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
