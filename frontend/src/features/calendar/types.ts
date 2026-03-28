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
}

export interface CreateScheduledPostRequest {
  channel_ids: number[];
  platforms: string[];
  scheduled_at: string;
  post_type: PostType;
  caption: string;
  media: { items: ScheduledMediaItem[] };
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
  overall_status: "scheduled" | "partial_failure" | "failed";
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
