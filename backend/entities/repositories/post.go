package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/shekhar8352/PostEaze/entities"
	"github.com/shekhar8352/PostEaze/utils/database"
)

// CreatePost creates a new post in the database
func CreatePost(ctx context.Context, post *entities.Post) error {
	db := database.GetDB()
	query := `
		INSERT INTO posts (
			channel_ids, owner_id, providers, provider_post_ids,
			source, post_type, caption, media, published_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRowContext(ctx, query,
		post.ChannelIDs,
		post.OwnerID,
		post.Providers,
		post.ProviderPostIDs,
		post.Source,
		post.PostType,
		post.Caption,
		post.Media,
		post.PublishedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

// GetPostByProviderID retrieves a post by channel ID and provider post ID
// Uses JSONB containment to check if provider_post_ids contains the given provider:id pair
func GetPostByProviderID(ctx context.Context, channelID int64, providerPostID string) (*entities.Post, error) {
	db := database.GetDB()
	// Check if channel_ids contains channelID AND provider_post_ids contains instagram:providerPostID
	query := `
		SELECT id, channel_ids, owner_id, providers, provider_post_ids, 
		       source, post_type, caption, media, published_at, created_at, updated_at
		FROM posts
		WHERE $1 = ANY(channel_ids) 
		  AND provider_post_ids->>'instagram' = $2
	`

	var post entities.Post
	err := db.QueryRowContext(ctx, query, channelID, providerPostID).Scan(
		&post.ID,
		&post.ChannelIDs,
		&post.OwnerID,
		&post.Providers,
		&post.ProviderPostIDs,
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
	return GetPosts(ctx, PostFilters{
		ChannelID: &channelID,
		Limit:     limit,
	})
}

// PostFilters defines filters for GetPosts
type PostFilters struct {
	ChannelID *int64
	Provider  *string
	Limit     int
	Offset    int
}

// GetPosts returns posts matching filters
func GetPosts(ctx context.Context, filters PostFilters) ([]entities.Post, error) {
	db := database.GetDB()

	query := `
		SELECT id, channel_ids, owner_id, providers, provider_post_ids,
		       source, post_type, caption, media, published_at, created_at, updated_at
		FROM posts
		WHERE 1=1
	`
	var args []interface{}
	argCount := 1

	if filters.ChannelID != nil {
		query += fmt.Sprintf(" AND $%d = ANY(channel_ids)", argCount)
		args = append(args, *filters.ChannelID)
		argCount++
	}

	if filters.Provider != nil {
		query += fmt.Sprintf(" AND $%d = ANY(providers)", argCount)
		args = append(args, *filters.Provider)
		argCount++
	}

	query += " ORDER BY published_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
		argCount++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
		argCount++
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entities.Post
	for rows.Next() {
		var post entities.Post
		err := rows.Scan(
			&post.ID,
			&post.ChannelIDs,
			&post.OwnerID,
			&post.Providers,
			&post.ProviderPostIDs,
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

// Helper function to get Instagram post ID from ProviderPostIDs JSONB
func GetInstagramPostID(post *entities.Post) string {
	if post.ProviderPostIDs == nil {
		return ""
	}
	var ids map[string]string
	if err := json.Unmarshal(post.ProviderPostIDs, &ids); err != nil {
		return ""
	}
	return ids["instagram"]
}
