package repositories

import (
	"context"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// UpsertSocialComment inserts or updates a comment; unique on (channel_id, provider, provider_comment_id).
func UpsertSocialComment(ctx context.Context, c *entities.SocialComment) error {
	db := database.GetDB()
	query := `
		INSERT INTO social_comments (
			channel_id, provider, post_id, provider_content_id, provider_comment_id,
			parent_provider_comment_id, text, author_provider_user_id, author_username,
			raw_payload, comment_created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (channel_id, provider, provider_comment_id)
		DO UPDATE SET
			post_id = COALESCE(EXCLUDED.post_id, social_comments.post_id),
			provider_content_id = EXCLUDED.provider_content_id,
			text = COALESCE(EXCLUDED.text, social_comments.text),
			parent_provider_comment_id = COALESCE(EXCLUDED.parent_provider_comment_id, social_comments.parent_provider_comment_id),
			author_provider_user_id = COALESCE(EXCLUDED.author_provider_user_id, social_comments.author_provider_user_id),
			author_username = COALESCE(EXCLUDED.author_username, social_comments.author_username),
			raw_payload = COALESCE(EXCLUDED.raw_payload, social_comments.raw_payload),
			comment_created_at = COALESCE(EXCLUDED.comment_created_at, social_comments.comment_created_at),
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query,
		c.ChannelID,
		c.Provider,
		c.PostID,
		c.ProviderContentID,
		c.ProviderCommentID,
		c.ParentProviderCommentID,
		c.Text,
		c.AuthorProviderUserID,
		c.AuthorUsername,
		jsonbOrEmpty(c.RawPayload),
		c.CommentCreatedAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}
