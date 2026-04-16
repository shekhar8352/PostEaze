package entities

import "time"

// SocialComment is a comment on provider content (e.g. IG media, FB post), keyed for idempotent upserts.
type SocialComment struct {
	ID                      int64      `db:"id"`
	ChannelID               int64      `db:"channel_id"`
	Provider                string     `db:"provider"`
	PostID                  *int64     `db:"post_id"`
	ProviderContentID       string     `db:"provider_content_id"`
	ProviderCommentID       string     `db:"provider_comment_id"`
	ParentProviderCommentID *string    `db:"parent_provider_comment_id"`
	Text                    *string    `db:"text"`
	AuthorProviderUserID    *string    `db:"author_provider_user_id"`
	AuthorUsername          *string    `db:"author_username"`
	RawPayload              []byte     `db:"raw_payload"`
	CommentCreatedAt        *time.Time `db:"comment_created_at"`
	CreatedAt               time.Time  `db:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at"`
}
