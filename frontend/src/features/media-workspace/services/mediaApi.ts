import apiClient from "@/services/api/client";
import type {
  ApiSuccessEnvelope,
  MediaAsset,
  MediaAssetListResponse,
  MediaVersion,
  UpdateMediaAssetPayload,
} from "../types";

export const mediaApi = {
  async list(
    status?: string,
    limit = 20,
    offset = 0
  ): Promise<MediaAssetListResponse> {
    const params = new URLSearchParams();
    if (status) params.set("status", status);
    params.set("limit", String(limit));
    params.set("offset", String(offset));
    const { data } = await apiClient.get<
      ApiSuccessEnvelope<MediaAssetListResponse>
    >(`/v1/media-assets?${params}`);
    return data.data;
  },

  async get(id: number): Promise<MediaAsset> {
    const { data } = await apiClient.get<ApiSuccessEnvelope<MediaAsset>>(
      `/v1/media-assets/${id}`
    );
    return data.data;
  },

  async create(formData: FormData): Promise<MediaAsset> {
    const { data } = await apiClient.post<ApiSuccessEnvelope<MediaAsset>>(
      "/v1/media-assets",
      formData,
      { headers: { "Content-Type": "multipart/form-data" }, timeout: 120000 }
    );
    return data.data;
  },

  async update(id: number, body: UpdateMediaAssetPayload): Promise<void> {
    await apiClient.put(`/v1/media-assets/${id}`, body);
  },

  async remove(id: number): Promise<void> {
    await apiClient.delete(`/v1/media-assets/${id}`);
  },

  async addVersion(id: number, formData: FormData): Promise<MediaVersion> {
    const { data } = await apiClient.post<ApiSuccessEnvelope<MediaVersion>>(
      `/v1/media-assets/${id}/versions`,
      formData,
      { headers: { "Content-Type": "multipart/form-data" }, timeout: 120000 }
    );
    return data.data;
  },

  async deleteVersion(assetId: number, versionId: number): Promise<void> {
    await apiClient.delete(
      `/v1/media-assets/${assetId}/versions/${versionId}`
    );
  },

  async setCurrentVersion(
    assetId: number,
    versionId: number
  ): Promise<void> {
    await apiClient.put(`/v1/media-assets/${assetId}/current-version`, {
      version_id: versionId,
    });
  },
};
