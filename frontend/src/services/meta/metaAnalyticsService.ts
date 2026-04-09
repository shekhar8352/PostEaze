import apiClient from '@/services/api/client';
import { type ApiResponse } from '@/services/api/types';

export interface SyncMetaAnalyticsResult {
  channel_id: number;
  provider: string;
  status: string;
  message?: string;
}

export interface SyncMetaAnalyticsResponse {
  results: SyncMetaAnalyticsResult[];
}

/**
 * Pull Instagram + Facebook insights for the current user’s channels (can take 30–90s).
 */
export async function syncMetaAnalytics(channelIds?: number[]): Promise<SyncMetaAnalyticsResponse> {
  const body =
    channelIds && channelIds.length > 0 ? { channel_ids: channelIds } : {};
  const response = await apiClient.post<ApiResponse<SyncMetaAnalyticsResponse>>(
    '/v1/meta/analytics/sync',
    body,
    { timeout: 120_000 }
  );
  return response.data.data;
}
