import apiClient from "@/services/api/client";
import type {
  ApiSuccessEnvelope,
  ComparisonResponse,
  DashboardResponse,
  ProfileAnalyticsResponse,
  StoryAnalyticsResponse,
  TopPostsResponse,
  TopPostsSort,
} from "../types";

function qs(params: Record<string, string | number | undefined>): string {
  const sp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== "") sp.set(k, String(v));
  });
  const s = sp.toString();
  return s ? `?${s}` : "";
}

export interface DateRangeParams {
  startDate: string;
  endDate: string;
}

async function unwrap<T>(promise: Promise<{ data: ApiSuccessEnvelope<T> }>): Promise<T> {
  const { data } = await promise;
  return data.data;
}

export const analyticsService = {
  getDashboard(channelId: number, range: DateRangeParams) {
    return unwrap(
      apiClient.get<ApiSuccessEnvelope<DashboardResponse>>(
        `/v1/channels/${channelId}/analytics/dashboard${qs({
          start_date: range.startDate,
          end_date: range.endDate,
        })}`
      )
    );
  },

  getComparison(channelId: number, range: DateRangeParams) {
    return unwrap(
      apiClient.get<ApiSuccessEnvelope<ComparisonResponse>>(
        `/v1/channels/${channelId}/analytics/comparison${qs({
          start_date: range.startDate,
          end_date: range.endDate,
        })}`
      )
    );
  },

  getProfileSeries(channelId: number, range: DateRangeParams) {
    return unwrap(
      apiClient.get<ApiSuccessEnvelope<ProfileAnalyticsResponse>>(
        `/v1/channels/${channelId}/analytics/profile${qs({
          start_date: range.startDate,
          end_date: range.endDate,
          limit: 0,
        })}`
      )
    );
  },

  getStories(channelId: number, range: DateRangeParams) {
    return unwrap(
      apiClient.get<ApiSuccessEnvelope<StoryAnalyticsResponse>>(
        `/v1/channels/${channelId}/analytics/stories${qs({
          start_date: range.startDate,
          end_date: range.endDate,
          limit: 100,
          offset: 0,
        })}`
      )
    );
  },

  getTopPosts(
    channelId: number,
    range: DateRangeParams,
    options: { sort?: TopPostsSort; limit?: number; postType?: string }
  ) {
    return unwrap(
      apiClient.get<ApiSuccessEnvelope<TopPostsResponse>>(
        `/v1/channels/${channelId}/analytics/top-posts${qs({
          start_date: range.startDate,
          end_date: range.endDate,
          sort: options.sort,
          limit: options.limit ?? 10,
          post_type: options.postType,
        })}`
      )
    );
  },
};
