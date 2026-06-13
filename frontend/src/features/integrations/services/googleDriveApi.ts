import apiClient from '@/services/api/client';
import type { ApiSuccessEnvelope } from '@/features/media-workspace/types';
import { getGoogleRedirectUri } from '../utils/googleOAuth.utils';

export interface GoogleDriveStatus {
  connected: boolean;
  email?: string;
  scopes?: string;
}

export interface DriveFileItem {
  id: string;
  name: string;
  mime_type: string;
  size: number;
  is_folder: boolean;
  modified_time?: string;
  thumbnail_link?: string;
}

export interface DriveFileListResponse {
  files: DriveFileItem[];
  next_page_token?: string;
}

export interface DriveRevisionItem {
  id: string;
  modified_time: string;
  keep_forever: boolean;
  original_filename?: string;
  size: number;
  mime_type?: string;
}

export interface DriveRevisionListResponse {
  revisions: DriveRevisionItem[];
}

export const googleDriveApi = {
  async connect(code: string): Promise<GoogleDriveStatus> {
    const { data } = await apiClient.post<ApiSuccessEnvelope<GoogleDriveStatus>>(
      '/v1/integrations/google-drive/connect',
      { code, redirect_uri: getGoogleRedirectUri('drive') }
    );
    return data.data;
  },

  async status(): Promise<GoogleDriveStatus> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<GoogleDriveStatus>>(
      '/v1/integrations/google-drive'
    );
    return data.data;
  },

  async disconnect(): Promise<void> {
    await apiClient.delete('/v1/integrations/google-drive');
  },

  async listFiles(params: {
    folderId?: string;
    pageToken?: string;
    q?: string;
  }): Promise<DriveFileListResponse> {
    const search = new URLSearchParams();
    if (params.folderId) search.set('folderId', params.folderId);
    if (params.pageToken) search.set('pageToken', params.pageToken);
    if (params.q) search.set('q', params.q);
    const qs = search.toString();
    const { data } = await apiClient.get<ApiSuccessEnvelope<DriveFileListResponse>>(
      `/v1/integrations/google-drive/files${qs ? `?${qs}` : ''}`
    );
    return data.data;
  },

  async listRevisions(fileId: string): Promise<DriveRevisionListResponse> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<DriveRevisionListResponse>>(
      `/v1/integrations/google-drive/files/${fileId}/revisions`
    );
    return data.data;
  },

  async importFile(body: {
    file_id: string;
    revision_id?: string;
    title?: string;
    label?: string;
  }) {
    const { data } = await apiClient.post<ApiSuccessEnvelope<unknown>>(
      '/v1/media-assets/import/google-drive',
      body,
      { timeout: 300000 }
    );
    return data.data;
  },

  async importRevision(
    assetId: number,
    body: { revision_id: string; label?: string }
  ) {
    const { data } = await apiClient.post<ApiSuccessEnvelope<unknown>>(
      `/v1/media-assets/${assetId}/versions/import-drive-revision`,
      body,
      { timeout: 300000 }
    );
    return data.data;
  },
};
