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

// RegisterHandlers registers task handlers to the mux
func RegisterHandlers(mux *asynq.ServeMux) {
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	mux.HandleFunc(TypeLogMessage, HandleLogMessageTask)
	mux.HandleFunc(TypeInstagramComment, HandleInstagramCommentTask)
	mux.HandleFunc(TypeInstagramMention, HandleInstagramMentionTask)
	mux.HandleFunc(TypeInstagramStoryInsight, HandleInstagramStoryInsightTask)
	mux.HandleFunc(TypeSyncInstagramProfiles, HandleSyncInstagramProfilesTask)
	mux.HandleFunc(TypeSyncInstagramPosts, HandleSyncInstagramPostsTask)
	mux.HandleFunc(TypeSyncInstagramAnalytics, HandleSyncInstagramAnalyticsTask)
}

// HandleEmailDeliveryTask handles email delivery tasks
func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Sending Email to User: ", p.UserID, ", Subject: ", p.Subject)
	// Logic to send email would go here
	return nil
}

// HandleLogMessageTask handles log message tasks
func HandleLogMessageTask(ctx context.Context, t *asynq.Task) error {
	var p LogMessagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Log Task: ", p.Message)
	return nil
}

// HandleInstagramCommentTask handles Instagram comment events
func HandleInstagramCommentTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Processing Instagram Comment: ", change)
	// TODO: Implement comment processing logic (e.g., save to DB, notify user)
	return nil
}

// HandleInstagramMentionTask handles Instagram mention events
func HandleInstagramMentionTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Processing Instagram Mention: ", change)
	// TODO: Implement mention processing logic
	return nil
}

// HandleInstagramStoryInsightTask handles Instagram story insight events
func HandleInstagramStoryInsightTask(ctx context.Context, t *asynq.Task) error {
	var change map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &change); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	utils.Logger.Info(ctx, "Processing Instagram Story Insight: ", change)
	// TODO: Implement story insight processing logic
	return nil
}

// HandleSyncInstagramProfilesTask syncs Instagram profile data for all active channels
func HandleSyncInstagramProfilesTask(ctx context.Context, t *asynq.Task) error {
	utils.Logger.Info(ctx, "Starting Instagram profile sync job")

	// Fetch all active Instagram channels
	channels, err := repositories.GetAllActiveInstagramChannels(ctx)
	if err != nil {
		utils.Logger.Error(ctx, "Failed to fetch Instagram channels: ", err)
		return err
	}

	utils.Logger.Info(ctx, "Found ", len(channels), " active Instagram channels to sync")

	// Process each channel
	for _, channel := range channels {
		if err := syncChannelProfile(ctx, channel); err != nil {
			// Log error and continue with next channel
			utils.Logger.Error(ctx, "Failed to sync channel ", channel.ID, ": ", err)
			continue
		}
	}

	utils.Logger.Info(ctx, "Instagram profile sync job completed")
	return nil
}

func syncChannelProfile(ctx context.Context, channel entities.Channel) error {
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

	// Fetch Instagram profile details
	provider := instagram.NewInstagramProvider()
	pageDetails, err := provider.GetPageDetails(decryptedToken)
	if err != nil {
		return fmt.Errorf("failed to fetch page details from Instagram: %w", err)
	}

	// Parse existing metadata
	var existingMetadata map[string]interface{}
	if len(channel.Metadata) > 0 {
		if err := json.Unmarshal(channel.Metadata, &existingMetadata); err != nil {
			existingMetadata = make(map[string]interface{})
		}
	} else {
		existingMetadata = make(map[string]interface{})
	}

	// Merge new data with existing metadata (preserve email)
	existingMetadata["username"] = pageDetails.Username
	existingMetadata["name"] = pageDetails.Name
	existingMetadata["biography"] = pageDetails.Biography
	existingMetadata["followers_count"] = pageDetails.FollowersCount
	existingMetadata["follows_count"] = pageDetails.FollowsCount
	existingMetadata["media_count"] = pageDetails.MediaCount
	existingMetadata["profile_picture_url"] = pageDetails.ProfilePictureURL
	existingMetadata["website"] = pageDetails.Website
	existingMetadata["last_synced_at"] = time.Now().Format(time.RFC3339)

	// Marshal updated metadata
	updatedMetadata, err := json.Marshal(existingMetadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Update channel metadata
	if err := repositories.UpdateChannelMetadata(ctx, channel.ID, updatedMetadata); err != nil {
		return fmt.Errorf("failed to update channel metadata: %w", err)
	}

	utils.Logger.Info(ctx, "Successfully synced profile for channel ", channel.ID)
	return nil
}
