package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

// UpdatePostFromInstagramSync refreshes fields from the Graph API when the post row already exists.
func UpdatePostFromInstagramSync(ctx context.Context, postID int64, postType *string, caption *string, media []byte, publishedAt *time.Time) error {
	db := database.GetDB()
	query := `
		UPDATE posts
		SET post_type = $2,
		    caption = $3,
		    media = $4,
		    published_at = $5,
		    updated_at = NOW()
		WHERE id = $1
	`
	res, err := db.ExecContext(ctx, query, postID, postType, caption, media, publishedAt)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("no post row updated for id %d", postID)
	}
	return nil
}

// GetPostByProviderID retrieves a post by channel ID and provider post ID (Instagram).
func GetPostByProviderID(ctx context.Context, channelID int64, providerPostID string) (*entities.Post, error) {
	return GetPostByProviderKey(ctx, channelID, "instagram", providerPostID)
}

// GetPostByProviderKey retrieves a post where provider_post_ids->>provider matches providerPostID.
func GetPostByProviderKey(ctx context.Context, channelID int64, provider, providerPostID string) (*entities.Post, error) {
	db := database.GetDB()
	query := `
		SELECT id, channel_ids, owner_id, providers, provider_post_ids, 
		       source, post_type, caption, media, published_at, created_at, updated_at
		FROM posts
		WHERE $1 = ANY(channel_ids) 
		  AND provider_post_ids->>$3 = $2
	`

	var post entities.Post
	err := db.QueryRowContext(ctx, query, channelID, providerPostID, provider).Scan(
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

// GetFirstChannelIDForInstagramMedia returns the first channel_id from a post row that owns this Instagram media id.
// Used when webhook payloads omit entry id (legacy task format). ok is false if no match.
func GetFirstChannelIDForInstagramMedia(ctx context.Context, mediaID string) (channelID int64, ok bool, err error) {
	db := database.GetDB()
	query := `
		SELECT (channel_ids)[1]
		FROM posts
		WHERE provider_post_ids->>'instagram' = $1
		  AND cardinality(channel_ids) > 0
		LIMIT 1
	`
	var cid sql.NullInt64
	err = db.QueryRowContext(ctx, query, mediaID).Scan(&cid)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	if !cid.Valid {
		return 0, false, nil
	}
	return cid.Int64, true, nil
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
	return GetProviderPostID(post, "instagram")
}

// GetProviderPostID returns the external post id for a provider key in provider_post_ids JSONB.
func GetProviderPostID(post *entities.Post, provider string) string {
	if post.ProviderPostIDs == nil {
		return ""
	}
	var ids map[string]string
	if err := json.Unmarshal(post.ProviderPostIDs, &ids); err != nil {
		return ""
	}
	return ids[provider]
}
