package tasks

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/entities/repositories"
	"github.com/shekhar8352/PostEaze/socialcomments"
	"github.com/shekhar8352/PostEaze/utils"
)

// processInstagramCommentPayload parses webhook task JSON, resolves channel and post, upserts social_comments.
func processInstagramCommentPayload(ctx context.Context, raw []byte) error {
	wrapper, err := socialcomments.DecodeJSONMap(raw)
	if err != nil {
		return fmt.Errorf("decode task json: %v: %w", err, asynq.SkipRetry)
	}

	var change map[string]interface{}
	if nested, ok := wrapper["change"].(map[string]interface{}); ok {
		change = nested
	} else {
		change = wrapper
	}

	field, _ := change["field"].(string)
	if field != "" && field != "comments" {
		utils.Logger.Info(ctx, "Skipping non-comments field: %s", field)
		return nil
	}

	value, ok := change["value"].(map[string]interface{})
	if !ok || value == nil {
		return fmt.Errorf("missing value map: %w", asynq.SkipRetry)
	}

	norm, err := socialcomments.ParseInstagramCommentValue(value)
	if err != nil {
		utils.Logger.Warn(ctx, "Instagram comment parse failed: %v (value keys: check Meta payload shape)", err)
		return fmt.Errorf("parse instagram comment: %v: %w", err, asynq.SkipRetry)
	}

	entryID := socialcomments.StringFromAny(wrapper["entry_id"])
	var ch *entities.Channel
	if entryID != "" {
		var err error
		ch, err = repositories.GetChannelByProviderChannelID(ctx, "instagram", entryID)
		if err != nil {
			return fmt.Errorf("get channel by provider id: %w", err)
		}
	}
	if ch == nil {
		cid, ok, err := repositories.GetFirstChannelIDForInstagramMedia(ctx, norm.MediaID)
		if err != nil {
			return fmt.Errorf("fallback channel from media: %w", err)
		}
		if !ok {
			utils.Logger.Warn(ctx, "No channel for Instagram comment (entry_id=%q media_id=%s comment_id=%s)", entryID, norm.MediaID, norm.CommentID)
			return nil
		}
		ch, err = repositories.GetChannelByID(ctx, cid)
		if err != nil {
			return fmt.Errorf("get channel by id: %w", err)
		}
		if ch == nil {
			return nil
		}
	}

	var postID *int64
	post, err := repositories.GetPostByProviderID(ctx, ch.ID, norm.MediaID)
	if err != nil {
		return fmt.Errorf("get post by provider: %w", err)
	}
	if post != nil {
		postID = &post.ID
	}

	text := norm.Text
	var textPtr *string
	if text != "" {
		textPtr = &text
	}

	rawPayload := socialcomments.MarshalRawPayload(wrapper)

	sc := &entities.SocialComment{
		ChannelID:               ch.ID,
		Provider:                "instagram",
		PostID:                  postID,
		ProviderContentID:       norm.MediaID,
		ProviderCommentID:       norm.CommentID,
		ParentProviderCommentID: norm.ParentCommentID,
		Text:                    textPtr,
		AuthorProviderUserID:    norm.AuthorUserID,
		AuthorUsername:          norm.AuthorUsername,
		RawPayload:              rawPayload,
		CommentCreatedAt:        nil,
	}

	if err := repositories.UpsertSocialComment(ctx, sc); err != nil {
		return fmt.Errorf("upsert social comment: %w", err)
	}

	utils.Logger.Info(ctx, "Stored Instagram comment id=%s channel=%d post=%v", norm.CommentID, ch.ID, postID)
	return nil
}
