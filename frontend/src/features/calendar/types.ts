/** Matches backend Gin success envelope */
export interface ApiSuccessEnvelope<T> {
  status: string;
  msg: string;
  data: T;
}

export type PostType = "image" | "video" | "carousel";

export interface ScheduledMediaItem {
  url: string;
  kind: "image" | "video";
  /** When set, backend marks this workspace asset `published` after a successful Instagram post */
  media_asset_id?: number;
}

export interface CreateScheduledPostRequest {
  channel_ids: number[];
  platforms: string[];
  /** Omit or null when publish_now is true */
  scheduled_at?: string | null;
  publish_now?: boolean;
  post_type: PostType;
  caption: string;
  media: { items: ScheduledMediaItem[] };
  /**
   * Optional Studio Piece to auto-link after the scheduled post is created.
   * Backend silently ignores link failures so schedule calls stay resilient.
   */
  piece_id?: number;
}

export interface ChannelScheduleResult {
  channel_id: number;
  success: boolean;
  creation_id?: string;
  published_media_id?: string;
  error_message?: string;
}

export interface CreateScheduledPostResponse {
  scheduled_post_id: number;
  overall_status: "published" | "scheduled" | "partial_failure" | "failed";
  publish_now: boolean;
  results: ChannelScheduleResult[];
}

export interface ScheduledPostListItem {
  id: number;
  channel_ids: number[];
  platforms: string[];
  scheduled_at: string;
  status: string;
  post_type: string;
  caption?: string;
  media: unknown;
  provider_state?: unknown;
  created_at: string;
  updated_at: string;
}

export interface ListScheduledPostsResponse {
  posts: ScheduledPostListItem[];
}

export type CalendarViewMode = "month" | "week" | "day";

/** react-big-calendar event shape for scheduled posts */
export interface CalendarScheduledEvent {
  id: number;
  title: string;
  start: Date;
  end: Date;
  resource: ScheduledPostListItem;
}
