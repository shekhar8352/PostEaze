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
	// Valid metrics: impressions, reach, profile_views, follower_count, email_contacts, phone_call_clicks, text_message_clicks, get_directions_clicks, website_clicks
	metrics := []string{"impressions", "reach", "profile_views", "follower_count"}

	// Calculate time range (last 7 days)
	now := time.Now()
	since := now.AddDate(0, 0, -7).Unix()
	until := now.Unix()

	insightsResp, err := provider.GetProfileInsights(accessToken, igUserID, metrics, since, until)
	if err != nil {
		return fmt.Errorf("failed to fetch profile insights: %w", err)
	}

	// Process insights data
	for _, insight := range insightsResp.Data {
		for _, value := range insight.Values {
			// Parse end_time
			var date time.Time
			if value.EndTime != "" {
				date, err = time.Parse(time.RFC3339, value.EndTime)
				if err != nil {
					utils.Logger.Error(ctx, fmt.Sprintf("Failed to parse end_time: %v", err))
					continue
				}
			} else {
				date = time.Now()
			}

			// Get or create analytics record
			analytics := &entities.InstagramProfileAnalytics{
				ChannelID: channelID,
				Date:      date.Truncate(24 * time.Hour),
			}

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
			case "impressions":
				analytics.Impressions = &intValue
			case "reach":
				analytics.Reach = &intValue
			case "profile_views":
				analytics.ProfileViews = &intValue
			case "follower_count":
				analytics.FollowerCount = &intValue
			}

			// Marshal raw data
			rawData, _ := json.Marshal(insight)
			analytics.Raw = rawData

			// Upsert analytics
			if err := repositories.UpsertInstagramProfileAnalytics(ctx, analytics); err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to upsert profile analytics: %v", err))
			}
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

		if isStory {
			insightsResp, err = provider.GetStoryInsights(accessToken, post.ProviderPostID)
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
			insightsResp, err = provider.GetMediaInsights(accessToken, post.ProviderPostID, metrics)
		}

		if err != nil {
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to fetch insights for post %s: %v", post.ProviderPostID, err))
			continue
		}

		// Process insights
		if isStory {
			if err := processStoryInsights(ctx, post, insightsResp); err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to process story insights: %v", err))
			}
		} else {
			if err := processPostInsights(ctx, post, insightsResp); err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to process post insights: %v", err))
			}
		}
	}

	return nil
}

func processPostInsights(ctx context.Context, post entities.Post, insightsResp *instagram.InsightsResponse) error {
	analytics := &entities.InstagramPostAnalytics{
		ChannelID: post.ChannelID,
		PostID:    post.ID,
		Date:      time.Now().Truncate(24 * time.Hour),
	}

	// Extract metrics
	for _, insight := range insightsResp.Data {
		if len(insight.Values) == 0 {
			continue
		}

		// Get the value (lifetime metrics have single value)
		var intValue int
		switch v := insight.Values[0].Value.(type) {
		case float64:
			intValue = int(v)
		case int:
			intValue = v
		default:
			continue
		}

		// Map to analytics fields
		switch insight.Name {
		case "impressions":
			analytics.Impressions = &intValue
		case "reach":
			analytics.Reach = &intValue
		case "likes":
			analytics.Likes = &intValue
		case "comments":
			analytics.Comments = &intValue
		case "saved", "saves":
			analytics.Saves = &intValue
		case "shares":
			analytics.Shares = &intValue
		case "plays", "video_views":
			analytics.VideoViews = &intValue
		case "total_interactions":
			// Total interactions is a sum of all engagement, we can skip or store separately
			continue
		}
	}

	// Marshal raw data
	rawData, _ := json.Marshal(insightsResp)
	analytics.Raw = rawData

	// Upsert analytics
	return repositories.UpsertInstagramPostAnalytics(ctx, analytics)
}

func processStoryInsights(ctx context.Context, post entities.Post, insightsResp *instagram.InsightsResponse) error {
	analytics := &entities.InstagramStoryAnalytics{
		ChannelID: post.ChannelID,
		PostID:    post.ID,
		Date:      time.Now().Truncate(24 * time.Hour),
	}

	// Extract metrics
	for _, insight := range insightsResp.Data {
		if len(insight.Values) == 0 {
			continue
		}

		var intValue int
		switch v := insight.Values[0].Value.(type) {
		case float64:
			intValue = int(v)
		case int:
			intValue = v
		default:
			continue
		}

		// Map to analytics fields
		switch insight.Name {
		case "impressions":
			analytics.Impressions = &intValue
		case "reach":
			analytics.Reach = &intValue
		case "exits":
			analytics.Exits = &intValue
		case "replies":
			analytics.Replies = &intValue
		case "taps_forward":
			analytics.TapsForward = &intValue
		case "taps_back":
			analytics.TapsBackward = &intValue
		}
	}

	// Marshal raw data
	rawData, _ := json.Marshal(insightsResp)
	analytics.Raw = rawData

	// Upsert analytics
	return repositories.UpsertInstagramStoryAnalytics(ctx, analytics)
}
