import apiClient from '@/services/api/client';
import { type ApiResponse } from '@/services/api/types';

export interface MetaPageRow {
  id: string;
  name: string;
  access_token: string;
  category?: string;
  tasks?: string[];
}

/**
 * Exchange OAuth code for long-lived user token and list Facebook Pages (code is consumed).
 */
export async function fetchMetaPagesFromCode(
  code: string,
  redirectUri: string
): Promise<MetaPageRow[]> {
  const { data } = await apiClient.post<{ pages: MetaPageRow[] }>('/v1/meta/callback', {
    code,
    redirect_uri: redirectUri,
  });
  return data.pages ?? [];
}

export interface CreateFacebookChannelPayload {
  page_id: string;
  page_access_token: string;
  channel_name?: string;
}

export interface CreateFacebookChannelResponse {
  channel_id: number;
  channel_name: string;
}

/**
 * Persist a Facebook Page channel using tokens returned from fetchMetaPagesFromCode.
 */
export async function createFacebookChannelFromPageToken(
  payload: CreateFacebookChannelPayload
): Promise<CreateFacebookChannelResponse> {
  const response = await apiClient.post<ApiResponse<CreateFacebookChannelResponse>>(
    '/v1/channels/facebook/create',
    {
      page_id: payload.page_id,
      page_access_token: payload.page_access_token,
      channel_name: payload.channel_name,
    }
  );
  return response.data.data;
}
