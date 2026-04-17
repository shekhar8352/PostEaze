import { useQuery } from "@tanstack/react-query";
import { analyticsService, type DateRangeParams } from "./analyticsService";
import type { TopPostsSort } from "../types";

export const analyticsKeys = {
  all: ["analytics"] as const,
  dashboard: (channelId: number, range: DateRangeParams) =>
    [...analyticsKeys.all, "dashboard", channelId, range] as const,
  comparison: (channelId: number, range: DateRangeParams) =>
    [...analyticsKeys.all, "comparison", channelId, range] as const,
  profile: (channelId: number, range: DateRangeParams) =>
    [...analyticsKeys.all, "profile", channelId, range] as const,
  stories: (channelId: number, range: DateRangeParams) =>
    [...analyticsKeys.all, "stories", channelId, range] as const,
  topPosts: (
    channelId: number,
    range: DateRangeParams,
    sort: TopPostsSort,
    postType?: string
  ) => [...analyticsKeys.all, "topPosts", channelId, range, sort, postType ?? ""] as const,
  dailyEngagement: (channelId: number, range: DateRangeParams) =>
    [...analyticsKeys.all, "dailyEngagement", channelId, range] as const,
};

export function useAnalyticsDashboard(channelId: number | null, range: DateRangeParams | null) {
  return useQuery({
    queryKey: analyticsKeys.dashboard(channelId ?? 0, range ?? { startDate: "", endDate: "" }),
    queryFn: () => analyticsService.getDashboard(channelId!, range!),
    enabled: Boolean(channelId && range?.startDate && range?.endDate),
  });
}

export function useAnalyticsComparison(channelId: number | null, range: DateRangeParams | null) {
  return useQuery({
    queryKey: analyticsKeys.comparison(channelId ?? 0, range ?? { startDate: "", endDate: "" }),
    queryFn: () => analyticsService.getComparison(channelId!, range!),
    enabled: Boolean(channelId && range?.startDate && range?.endDate),
  });
}

export function useProfileTimeSeries(channelId: number | null, range: DateRangeParams | null) {
  return useQuery({
    queryKey: analyticsKeys.profile(channelId ?? 0, range ?? { startDate: "", endDate: "" }),
    queryFn: () => analyticsService.getProfileSeries(channelId!, range!),
    enabled: Boolean(channelId && range?.startDate && range?.endDate),
    select: (data) => ({
      ...data,
      analytics: [...data.analytics].sort(
        (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
      ),
    }),
  });
}

export function useStoriesAnalytics(channelId: number | null, range: DateRangeParams | null) {
  return useQuery({
    queryKey: analyticsKeys.stories(channelId ?? 0, range ?? { startDate: "", endDate: "" }),
    queryFn: () => analyticsService.getStories(channelId!, range!),
    enabled: Boolean(channelId && range?.startDate && range?.endDate),
  });
}

export function useTopPosts(
  channelId: number | null,
  range: DateRangeParams | null,
  sort: TopPostsSort,
  postType?: string
) {
  return useQuery({
    queryKey: analyticsKeys.topPosts(
      channelId ?? 0,
      range ?? { startDate: "", endDate: "" },
      sort,
      postType
    ),
    queryFn: () =>
      analyticsService.getTopPosts(channelId!, range!, { sort, limit: 15, postType }),
    enabled: Boolean(channelId && range?.startDate && range?.endDate),
  });
}

export function useDailyEngagement(channelId: number | null, range: DateRangeParams | null) {
  return useQuery({
    queryKey: analyticsKeys.dailyEngagement(channelId ?? 0, range ?? { startDate: "", endDate: "" }),
    queryFn: () => analyticsService.getDailyEngagement(channelId!, range!),
    enabled: Boolean(channelId && range?.startDate && range?.endDate),
  });
}
