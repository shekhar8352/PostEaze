package modelsv1

// --- Studio requests ---

type CreateStudioRequest struct {
	Name       string `json:"name"`
	PieceLabel string `json:"piece_label" binding:"omitempty,max=32"`
}

type UpdateStudioRequest struct {
	Name       string `json:"name"`
	PieceLabel string `json:"piece_label" binding:"omitempty,max=32"`
}

// --- Phase requests ---

type CreatePhaseRequest struct {
	Name     string `json:"name" binding:"required"`
	Slug     string `json:"slug"`
	Kind     string `json:"kind" binding:"omitempty,oneof=idea script shoot edit review scheduled published custom"`
	Color    string `json:"color"`
	WIPLimit *int   `json:"wip_limit"`
}

type UpdatePhaseRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Kind     string `json:"kind" binding:"omitempty,oneof=idea script shoot edit review scheduled published custom"`
	Color    string `json:"color"`
	WIPLimit *int   `json:"wip_limit"`
}

type ReorderPhasesRequest struct {
	PhaseIDs []int64 `json:"phase_ids" binding:"required,min=1"`
}

// --- Piece requests ---

type CreatePieceRequest struct {
	PhaseID        int64   `json:"phase_id" binding:"required"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	ContentType    string  `json:"content_type" binding:"omitempty,oneof=post reel story video carousel other"`
	AssigneeUserID *string `json:"assignee_user_id"`
	DueAt          *string `json:"due_at"` // RFC3339
}

type UpdatePieceRequest struct {
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	ContentType    string  `json:"content_type" binding:"omitempty,oneof=post reel story video carousel other"`
	AssigneeUserID *string `json:"assignee_user_id"`
	DueAt          *string `json:"due_at"` // RFC3339 ("" clears)
}

type MovePieceRequest struct {
	PhaseID       int64  `json:"phase_id" binding:"required"`
	BeforePieceID *int64 `json:"before_piece_id"`
	AfterPieceID  *int64 `json:"after_piece_id"`
}

type SetPieceStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active archived"`
}

type LinkAssetRequest struct {
	MediaAssetID int64  `json:"media_asset_id" binding:"required"`
	Role         string `json:"role" binding:"omitempty,oneof=script raw edit thumbnail final attachment"`
}

type LinkScheduledPostRequest struct {
	ScheduledPostID int64  `json:"scheduled_post_id" binding:"required"`
	Role            string `json:"role" binding:"omitempty,oneof=primary cross_post repost"`
}

type CreatePieceCommentRequest struct {
	Body     string `json:"body" binding:"required"`
	ParentID *int64 `json:"parent_id"`
}

// --- Responses ---

type StudioResponse struct {
	ID         int64  `json:"id"`
	TeamID     string `json:"team_id"`
	Name       string `json:"name"`
	PieceLabel string `json:"piece_label"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PhaseResponse struct {
	ID         int64  `json:"id"`
	StudioID   int64  `json:"studio_id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	OrderIndex int    `json:"order_index"`
	Kind       string `json:"kind"`
	WIPLimit   *int   `json:"wip_limit"`
	IsDefault  bool   `json:"is_default"`
	Color      string `json:"color"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PieceResponse struct {
	ID             int64   `json:"id"`
	StudioID       int64   `json:"studio_id"`
	PhaseID        int64   `json:"phase_id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	ContentType    string  `json:"content_type"`
	Status         string  `json:"status"`
	AssigneeUserID *string `json:"assignee_user_id"`
	DueAt          *string `json:"due_at"`
	Position       string  `json:"position"`
	CreatedBy      string  `json:"created_by"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	ArchivedAt     *string `json:"archived_at"`
}

type PieceAssetResponse struct {
	MediaAssetID int64  `json:"media_asset_id"`
	Role         string `json:"role"`
	AddedAt      string `json:"added_at"`
}

type PieceScheduledPostResponse struct {
	ScheduledPostID int64  `json:"scheduled_post_id"`
	Role            string `json:"role"`
	CreatedAt       string `json:"created_at"`
}

type PieceActivityResponse struct {
	ID          int64   `json:"id"`
	PieceID     int64   `json:"piece_id"`
	ActorUserID *string `json:"actor_user_id"`
	Kind        string  `json:"kind"`
	Payload     any     `json:"payload"`
	CreatedAt   string  `json:"created_at"`
}

type PieceCommentResponse struct {
	ID        int64   `json:"id"`
	PieceID   int64   `json:"piece_id"`
	UserID    string  `json:"user_id"`
	Body      string  `json:"body"`
	ParentID  *int64  `json:"parent_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// PieceDetailResponse bundles a piece with its linked assets, scheduled posts,
// and recent activities/comments for the detail drawer.
type PieceDetailResponse struct {
	Piece          PieceResponse                `json:"piece"`
	Assets         []PieceAssetResponse         `json:"assets"`
	ScheduledPosts []PieceScheduledPostResponse `json:"scheduled_posts"`
	Activities     []PieceActivityResponse      `json:"activities"`
	Comments       []PieceCommentResponse       `json:"comments"`
}

// StudioBoardResponse is the single payload used by the Kanban board UI.
type StudioBoardResponse struct {
	Studio StudioResponse  `json:"studio"`
	Phases []PhaseResponse `json:"phases"`
	Pieces []PieceResponse `json:"pieces"`
}

// ListPiecesQuery filters for GET /v1/pieces.
type ListPiecesQuery struct {
	StudioID       int64  `form:"studio_id" binding:"required"`
	PhaseID        int64  `form:"phase_id"`
	Status         string `form:"status"`
	AssigneeUserID string `form:"assignee_user_id"`
}
