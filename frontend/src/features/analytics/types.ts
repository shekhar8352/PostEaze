/** Backend API envelope for Gin SendSuccess responses */
export interface ApiSuccessEnvelope<T> {
  status: string;
  msg: string;
  data: T;
}

export interface DateRangeMeta {
  start_date: string;
  end_date: string;
}

export interface PaginationMeta {
  limit: number;
  offset: number;
  total: number;
}

export interface AggregatedProfileOverview {
  total_reach: number;
  total_impressions: number;
  total_profile_views: number;
  total_website_clicks: number;
  total_views: number;
  total_accounts_engaged: number;
  total_interactions: number;
  average_reach: number;
  follower_growth: number;
  start_followers: number;
  end_followers: number;
}

export interface TopPostItem {
  post_id: number;
  post_type: string;
  caption: string;
  published_at?: string;
  thumbnail_url?: string;
  permalink?: string;
  impressions: number;
  reach: number;
  likes: number;
  comments: number;
  saves: number;
  shares: number;
  plays: number;
  engagement: number;
}

export interface PostsOverview {
  new_posts: number;
  total_posts: number;
  new_likes: number;
  new_comments: number;
  new_saves: number;
  new_shares: number;
  total_likes: number;
  total_comments: number;
  total_saves: number;
  total_shares: number;
  total_reach: number;
  total_impressions: number;
}

export interface AudienceBucket {
  label: string;
  value: number;
  pct: number;
}

export interface AudienceMetricBlock {
  key: string;
  title: string;
  items: AudienceBucket[];
}

/** Latest lifetime audience snapshot from Meta (dashboard payload). */
export interface AudienceDashboard {
  snapshot_date?: string;
  created_at?: string;
  metrics: AudienceMetricBlock[];
}

export interface DailyEngagementPoint {
  date: string;
  likes: number;
  comments: number;
  shares: number;
  saves: number;
  total: number;
}

export interface DailyEngagementSeriesResponse {
  meta: DateRangeMeta;
  series: DailyEngagementPoint[];
  note: string;
}

export interface DashboardResponse {
  meta: DateRangeMeta;
  overview: AggregatedProfileOverview;
  top_posts: TopPostItem[];
  posts_overview: PostsOverview;
  audience?: AudienceDashboard;
}

export interface PeriodBounds {
  start: string;
  end: string;
  prev_start: string;
  prev_end: string;
}

export interface ComparisonResponse {
  meta: DateRangeMeta;
  current: AggregatedProfileOverview;
  previous?: AggregatedProfileOverview;
  period: PeriodBounds;
}

export interface ProfileAnalyticsItem {
  id: number;
  date: string;
  follower_count?: number;
  impressions?: number;
  profile_views?: number;
  reach?: number;
  website_clicks?: number;
  email_clicks?: number;
  views?: number;
  accounts_engaged?: number;
  total_interactions?: number;
}

export interface ProfileAnalyticsResponse {
  meta: DateRangeMeta;
  pagination: PaginationMeta;
  analytics: ProfileAnalyticsItem[];
}

export interface StoryAnalyticsItem {
  id: number;
  post_id: number;
  date: string;
  impressions?: number;
  reach?: number;
  exits?: number;
  replies?: number;
  taps_forward?: number;
  taps_backward?: number;
  taps_exit?: number;
  caption: string;
  post_type: string;
  published_at?: string;
}

export interface StoryAnalyticsResponse {
  meta: DateRangeMeta;
  pagination: PaginationMeta;
  stories: StoryAnalyticsItem[];
}

export type TopPostsSort = "engagement" | "reach" | "impressions" | "plays";

export interface TopPostsResponse {
  meta: DateRangeMeta;
  sort: string;
  limit: number;
  top_posts: TopPostItem[];
}
