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

// maxPostsPerChannelCommentsSync limits Graph calls per cron run (each post may paginate comments).
const maxPostsPerChannelCommentsSync = 100

// HandleSyncInstagramCommentsTask fetches comments from Instagram Graph for synced posts and upserts social_comments.
func HandleSyncInstagramCommentsTask(ctx context.Context, t *asynq.Task) error {
	utils.Logger.Info(ctx, "Starting Instagram comments sync job")

	channels, err := repositories.GetAllActiveInstagramChannels(ctx)
	if err != nil {
		utils.Logger.Error(ctx, "Failed to fetch Instagram channels: %v", err)
		return err
	}

	api := instagram.NewInstagramProvider()
	totalUpserted := 0

	for _, channel := range channels {
		n, err := syncChannelInstagramComments(ctx, api, channel)
		if err != nil {
			utils.Logger.Error(ctx, "Comments sync failed for channel %d: %v", channel.ID, err)
			continue
		}
		totalUpserted += n
	}

	utils.Logger.Info(ctx, "Instagram comments sync job completed (upserted rows this run: %d)", totalUpserted)
	return nil
}

func syncChannelInstagramComments(ctx context.Context, api instagram.InstagramProvider, channel entities.Channel) (int, error) {
	token, err := repositories.GetLatestTokenByChannelID(ctx, channel.ID)
	if err != nil {
		return 0, fmt.Errorf("access token: %w", err)
	}

	decryptedToken, err := encryption.Decrypt(token.AccessToken)
	if err != nil {
		return 0, fmt.Errorf("decrypt token: %w", err)
	}

	posts, err := repositories.GetPostsByChannel(ctx, channel.ID, maxPostsPerChannelCommentsSync)
	if err != nil {
		return 0, fmt.Errorf("list posts: %w", err)
	}

	upserted := 0
	for i := range posts {
		post := &posts[i]
		mediaID := repositories.GetInstagramPostID(post)
		if mediaID == "" {
			continue
		}

		after := ""
		for {
			resp, err := api.GetMediaComments(decryptedToken, mediaID, after)
			if err != nil {
				utils.Logger.Warn(ctx, "GetMediaComments channel=%d media=%s: %v", channel.ID, mediaID, err)
				break
			}
			for i := range resp.Data {
				c := &resp.Data[i]
				if err := upsertGraphComment(ctx, channel.ID, post.ID, mediaID, c); err != nil {
					utils.Logger.Warn(ctx, "upsert comment %s: %v", c.ID, err)
					continue
				}
				upserted++
			}
			if resp.Paging == nil || resp.Paging.Cursors == nil || resp.Paging.Cursors.After == "" {
				break
			}
			after = resp.Paging.Cursors.After
		}
	}

	utils.Logger.Info(ctx, "Instagram comments sync channel=%d upserts=%d posts_scanned=%d", channel.ID, upserted, len(posts))
	return upserted, nil
}

func upsertGraphComment(ctx context.Context, channelID, postID int64, mediaID string, c *instagram.MediaCommentItem) error {
	raw, _ := json.Marshal(map[string]interface{}{
		"source": "instagram_graph_sync",
		"item":   c,
	})

	text := c.Text
	var textPtr *string
	if text != "" {
		textPtr = &text
	}

	var authorID *string
	var authorUsername *string
	if c.From != nil {
		if c.From.ID != "" {
			authorID = &c.From.ID
		}
		if c.From.Username != "" {
			authorUsername = &c.From.Username
		}
	}
	if authorUsername == nil && c.Username != "" {
		authorUsername = &c.Username
	}

	var parent *string
	if c.ParentID != nil && *c.ParentID != "" {
		parent = c.ParentID
	}

	var commentAt *time.Time
	if c.Timestamp != "" {
		commentAt = parseInstagramTimestamp(c.Timestamp)
	}

	sc := &entities.SocialComment{
		ChannelID:               channelID,
		Provider:                "instagram",
		PostID:                  &postID,
		ProviderContentID:       mediaID,
		ProviderCommentID:       c.ID,
		ParentProviderCommentID: parent,
		Text:                    textPtr,
		AuthorProviderUserID:    authorID,
		AuthorUsername:          authorUsername,
		RawPayload:              raw,
		CommentCreatedAt:        commentAt,
	}

	return repositories.UpsertSocialComment(ctx, sc)
}

func parseInstagramTimestamp(s string) *time.Time {
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05Z07:00",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			u := t.UTC()
			return &u
		}
	}
	return nil
}
