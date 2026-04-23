/**
 * Studio feature types.
 *
 * These mirror the response contracts defined in
 * backend/models/v1/studio.go so the frontend can consume the
 * /v1/studios, /v1/phases, and /v1/pieces endpoints type-safely.
 */

export interface ApiSuccessEnvelope<T> {
  status: string;
  msg: string;
  data: T;
}

// ---------------------------------------------------------------------------
// Core domain
// ---------------------------------------------------------------------------

export type PhaseKind =
  | "idea"
  | "script"
  | "shoot"
  | "edit"
  | "review"
  | "scheduled"
  | "published"
  | "custom";

export type PieceContentType =
  | "post"
  | "reel"
  | "story"
  | "video"
  | "carousel"
  | "other";

export type PieceStatus = "active" | "archived";

export type PieceAssetRole =
  | "script"
  | "raw"
  | "edit"
  | "thumbnail"
  | "final"
  | "attachment";

export type PieceScheduledPostRole = "primary" | "cross_post" | "repost";

// ---------------------------------------------------------------------------
// Entities
// ---------------------------------------------------------------------------

export interface Studio {
  id: number;
  team_id: string;
  name: string;
  piece_label: string;
  created_at: string;
  updated_at: string;
}

export interface Phase {
  id: number;
  studio_id: number;
  name: string;
  slug: string;
  order_index: number;
  kind: PhaseKind;
  wip_limit: number | null;
  is_default: boolean;
  color: string;
  created_at: string;
  updated_at: string;
}

export interface Piece {
  id: number;
  studio_id: number;
  phase_id: number;
  title: string;
  description: string;
  content_type: PieceContentType;
  status: PieceStatus;
  assignee_user_id: string | null;
  due_at: string | null;
  position: string;
  created_by: string;
  created_at: string;
  updated_at: string;
  archived_at: string | null;
}

export interface PieceAssetLink {
  media_asset_id: number;
  role: PieceAssetRole;
  added_at: string;
}

export interface PieceScheduledPostLink {
  scheduled_post_id: number;
  role: PieceScheduledPostRole;
  created_at: string;
}

export interface PieceActivity {
  id: number;
  piece_id: number;
  actor_user_id: string | null;
  kind: string;
  payload: unknown;
  created_at: string;
}

export interface PieceComment {
  id: number;
  piece_id: number;
  user_id: string;
  body: string;
  parent_id: number | null;
  created_at: string;
  updated_at: string;
}

// ---------------------------------------------------------------------------
// Aggregated responses
// ---------------------------------------------------------------------------

export interface PieceDetail {
  piece: Piece;
  assets: PieceAssetLink[];
  scheduled_posts: PieceScheduledPostLink[];
  activities: PieceActivity[];
  comments: PieceComment[];
}

export interface StudioBoard {
  studio: Studio;
  phases: Phase[];
  pieces: Piece[];
}

// ---------------------------------------------------------------------------
// Request payloads
// ---------------------------------------------------------------------------

export interface UpdateStudioPayload {
  name?: string;
  piece_label?: string;
}

export interface CreatePhasePayload {
  name: string;
  slug?: string;
  kind?: PhaseKind;
  color?: string;
  wip_limit?: number | null;
}

export interface UpdatePhasePayload {
  name?: string;
  slug?: string;
  kind?: PhaseKind;
  color?: string;
  wip_limit?: number | null;
}

export interface CreatePiecePayload {
  phase_id: number;
  title: string;
  description?: string;
  content_type?: PieceContentType;
  assignee_user_id?: string | null;
  /** RFC3339 timestamp */
  due_at?: string | null;
}

export interface UpdatePiecePayload {
  title?: string;
  description?: string | null;
  content_type?: PieceContentType;
  assignee_user_id?: string | null;
  /** RFC3339 timestamp; empty string clears */
  due_at?: string | null;
}

export interface MovePiecePayload {
  phase_id: number;
  before_piece_id?: number | null;
  after_piece_id?: number | null;
}

export interface LinkAssetPayload {
  media_asset_id: number;
  role?: PieceAssetRole;
}

export interface LinkScheduledPostPayload {
  scheduled_post_id: number;
  role?: PieceScheduledPostRole;
}

export interface CreateCommentPayload {
  body: string;
  parent_id?: number | null;
}
