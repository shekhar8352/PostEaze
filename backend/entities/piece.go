package entities

import (
	"time"

	"github.com/google/uuid"
)

// PieceContentType narrows the shape of a Piece for UI and channel-fit hints.
type PieceContentType string

const (
	PieceContentTypePost     PieceContentType = "post"
	PieceContentTypeReel     PieceContentType = "reel"
	PieceContentTypeStory    PieceContentType = "story"
	PieceContentTypeVideo    PieceContentType = "video"
	PieceContentTypeCarousel PieceContentType = "carousel"
	PieceContentTypeOther    PieceContentType = "other"
)

// PieceStatus is independent of Phase; archived pieces are hidden from the board by default.
type PieceStatus string

const (
	PieceStatusActive   PieceStatus = "active"
	PieceStatusArchived PieceStatus = "archived"
)

// Piece is the central content item ("Piece" is the canonical DB term; display
// label per-studio is "Piece" or "Drop" or custom). Pieces flow through Phases
// and can be linked to media_assets and scheduled_posts.
type Piece struct {
	ID             int64      `db:"id"`
	StudioID       int64      `db:"studio_id"`
	PhaseID        int64      `db:"phase_id"`
	Title          string     `db:"title"`
	Description    string     `db:"description"`
	ContentType    string     `db:"content_type"`
	Status         string     `db:"status"`
	AssigneeUserID *uuid.UUID `db:"assignee_user_id"`
	DueAt          *time.Time `db:"due_at"`
	Position       string     `db:"position"` // fractional-index string for cheap reorders
	CreatedBy      uuid.UUID  `db:"created_by"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	ArchivedAt     *time.Time `db:"archived_at"`
}

// PieceAssetRole tags the purpose of a linked media_asset on a Piece.
type PieceAssetRole string

const (
	PieceAssetRoleScript     PieceAssetRole = "script"
	PieceAssetRoleRaw        PieceAssetRole = "raw"
	PieceAssetRoleEdit       PieceAssetRole = "edit"
	PieceAssetRoleThumbnail  PieceAssetRole = "thumbnail"
	PieceAssetRoleFinal      PieceAssetRole = "final"
	PieceAssetRoleAttachment PieceAssetRole = "attachment"
)

// PieceAsset is the join between a Piece and a media_asset.
type PieceAsset struct {
	PieceID      int64     `db:"piece_id"`
	MediaAssetID int64     `db:"media_asset_id"`
	Role         string    `db:"role"`
	AddedAt      time.Time `db:"added_at"`
}

// PieceScheduledPostRole distinguishes primary vs. cross-post vs. repost publications.
type PieceScheduledPostRole string

const (
	PieceScheduledPostRolePrimary   PieceScheduledPostRole = "primary"
	PieceScheduledPostRoleCrossPost PieceScheduledPostRole = "cross_post"
	PieceScheduledPostRoleRepost    PieceScheduledPostRole = "repost"
)

// PieceScheduledPost is the join between a Piece and a scheduled_post.
type PieceScheduledPost struct {
	PieceID          int64     `db:"piece_id"`
	ScheduledPostID  int64     `db:"scheduled_post_id"`
	Role             string    `db:"role"`
	CreatedAt        time.Time `db:"created_at"`
}

// PieceActivityKind enumerates audit-log event types for a Piece.
type PieceActivityKind string

const (
	PieceActivityCreated        PieceActivityKind = "created"
	PieceActivityUpdated        PieceActivityKind = "updated"
	PieceActivityMoved          PieceActivityKind = "moved"
	PieceActivityAssigned       PieceActivityKind = "assigned"
	PieceActivityAssetLinked    PieceActivityKind = "asset_linked"
	PieceActivityAssetUnlinked  PieceActivityKind = "asset_unlinked"
	PieceActivityScheduled      PieceActivityKind = "scheduled"
	PieceActivityPublished      PieceActivityKind = "published"
	PieceActivityPublishFailed  PieceActivityKind = "publish_failed"
	PieceActivityCommented      PieceActivityKind = "commented"
	PieceActivityArchived       PieceActivityKind = "archived"
	PieceActivityRestored       PieceActivityKind = "restored"
)

// PieceActivity is an append-only audit log entry on a Piece.
type PieceActivity struct {
	ID          int64      `db:"id"`
	PieceID     int64      `db:"piece_id"`
	ActorUserID *uuid.UUID `db:"actor_user_id"`
	Kind        string     `db:"kind"`
	Payload     []byte     `db:"payload"` // JSONB
	CreatedAt   time.Time  `db:"created_at"`
}

// PieceComment is a threaded discussion entry on a Piece.
type PieceComment struct {
	ID        int64     `db:"id"`
	PieceID   int64     `db:"piece_id"`
	UserID    uuid.UUID `db:"user_id"`
	Body      string    `db:"body"`
	ParentID  *int64    `db:"parent_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
