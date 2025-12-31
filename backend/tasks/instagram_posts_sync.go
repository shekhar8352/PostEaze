package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/provider/instagram"
	"github.com/shekhar8352/PostEaze/utils"
	"github.com/shekhar8352/PostEaze/utils/encryption"
)

// HandleSyncInstagramPostsTask syncs Instagram posts for all active channels
func HandleSyncInstagramPostsTask(ctx context.Context, t *asynq.Task) error {
	utils.Logger.Info(ctx, "Starting Instagram posts sync job")

	// Fetch all active Instagram channels
	channels, err := repositories.GetAllActiveInstagramChannels(ctx)
	if err != nil {
		utils.Logger.Error(ctx, "Failed to fetch Instagram channels: ", err)
		return err
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Found %d active Instagram channels to sync posts", len(channels)))

	// Process each channel
	for _, channel := range channels {
		if err := syncChannelPosts(ctx, channel); err != nil {
			// Log error and continue with next channel
			utils.Logger.Error(ctx, fmt.Sprintf("Failed to sync posts for channel %d: %v", channel.ID, err))
			continue
		}
	}

	utils.Logger.Info(ctx, "Instagram posts sync job completed")
	return nil
}

func syncChannelPosts(ctx context.Context, channel entities.Channel) error {
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

	// Fetch media from Instagram
	provider := instagram.NewInstagramProvider()
	after := ""
	totalSynced := 0

	for {
		mediaResp, err := provider.GetMedia(decryptedToken, igUserID, after)
		if err != nil {
			return fmt.Errorf("failed to fetch media: %w", err)
		}

		// Process each media item
		for _, media := range mediaResp.Data {
			// Check if post already exists
			existingPost, err := repositories.GetPostByProviderID(ctx, channel.ID, media.ID)
			if err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Error checking existing post: %v", err))
				continue
			}

			if existingPost != nil {
				// Post already exists, skip
				continue
			}

			// Parse timestamp (Instagram format: "2020-08-09T04:12:37+0000")
			publishedAt, err := time.Parse("2006-01-02T15:04:05-0700", media.Timestamp)
			if err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to parse timestamp for media %s: %v", media.ID, err))
				publishedAt = time.Now()
			}

			// Determine post type
			postType := getPostType(media.MediaType)

			// Create media JSON
			mediaJSON, err := json.Marshal([]map[string]string{
				{
					"url":  media.MediaURL,
					"type": media.MediaType,
				},
			})
			if err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to marshal media JSON: %v", err))
				mediaJSON = []byte("[]")
			}

			// Create provider_post_ids JSONB
			providerPostIDs, _ := json.Marshal(map[string]string{
				"instagram": media.ID,
			})

			// Get owner ID as pointer
			ownerIDStr := channel.OwnerUserID.String()

			// Create new post
			post := &entities.Post{
				ChannelIDs:      pq.Int64Array{channel.ID},
				OwnerID:         &ownerIDStr,
				Providers:       pq.StringArray{"instagram"},
				ProviderPostIDs: providerPostIDs,

				Source:      "native",
				PostType:    &postType,
				Caption:     &media.Caption,
				Media:       mediaJSON,
				PublishedAt: &publishedAt,
			}

			if err := repositories.CreatePost(ctx, post); err != nil {
				utils.Logger.Error(ctx, fmt.Sprintf("Failed to create post %s: %v", media.ID, err))
				continue
			}

			totalSynced++
		}

		// Check if there are more pages
		if mediaResp.Paging == nil || mediaResp.Paging.Next == "" {
			break
		}

		if mediaResp.Paging.Cursors != nil {
			after = mediaResp.Paging.Cursors.After
		} else {
			break
		}
	}

	utils.Logger.Info(ctx, fmt.Sprintf("Synced %d new posts for channel %d", totalSynced, channel.ID))
	return nil
}

func getPostType(mediaType string) string {
	switch mediaType {
	case "IMAGE":
		return "image"
	case "VIDEO":
		return "video"
	case "CAROUSEL_ALBUM":
		return "carousel"
	default:
		return "post"
	}
}
