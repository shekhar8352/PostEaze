package repositories

import (
	"context"
	"database/sql"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// CreatePost creates a new post in the database
func CreatePost(ctx context.Context, post *entities.Post) error {
	db := database.GetDB()
	query := `
		INSERT INTO posts (channel_id, provider, provider_post_id, source, post_type, caption, media, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query,
		post.ChannelID,
		post.Provider,
		post.ProviderPostID,
		post.Source,
		post.PostType,
		post.Caption,
		post.Media,
		post.PublishedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

// GetPostByProviderID retrieves a post by channel ID and provider post ID
func GetPostByProviderID(ctx context.Context, channelID int64, providerPostID string) (*entities.Post, error) {
	db := database.GetDB()
	query := `
		SELECT id, channel_id, provider, provider_post_id, source, post_type, caption, media, published_at, created_at, updated_at
		FROM posts
		WHERE channel_id = $1 AND provider_post_id = $2
	`

	var post entities.Post
	err := db.QueryRowContext(ctx, query, channelID, providerPostID).Scan(
		&post.ID,
		&post.ChannelID,
		&post.Provider,
		&post.ProviderPostID,
		&post.Source,
		&post.PostType,
		&post.Caption,
		&post.Media,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// GetPostsByChannel retrieves all posts for a channel
func GetPostsByChannel(ctx context.Context, channelID int64, limit int) ([]entities.Post, error) {
	db := database.GetDB()
	query := `
		SELECT id, channel_id, provider, provider_post_id, source, post_type, caption, media, published_at, created_at, updated_at
		FROM posts
		WHERE channel_id = $1
		ORDER BY published_at DESC
		LIMIT $2
	`

	rows, err := db.QueryContext(ctx, query, channelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entities.Post
	for rows.Next() {
		var post entities.Post
		err := rows.Scan(
			&post.ID,
			&post.ChannelID,
			&post.Provider,
			&post.ProviderPostID,
			&post.Source,
			&post.PostType,
			&post.Caption,
			&post.Media,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}
