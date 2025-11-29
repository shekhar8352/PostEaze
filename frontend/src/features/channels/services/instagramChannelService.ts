import { BaseService } from "@/services/base/BaseService";
import apiClient from "@/services/api/client";
import { type ApiResponse } from "@/services/api/types";
import type {
  InstagramChannel,
  CreateInstagramChannelRequest,
  UpdateInstagramChannelRequest,
  InstagramChannelStats,
} from "../types/instagram.types";

class InstagramChannelService extends BaseService {
  constructor() {
    super("/v1/channels/instagram");
  }

  // Create Instagram Channel
  async createChannel(
    data: CreateInstagramChannelRequest
  ): Promise<InstagramChannel> {
    const response = await apiClient.post<ApiResponse<InstagramChannel>>(
      this.endpoint,
      data
    );
    return response.data.data;
  }

  // Get All Instagram Channels
  async getChannels(): Promise<InstagramChannel[]> {
    const response = await apiClient.get<ApiResponse<InstagramChannel[]>>(
      this.endpoint
    );
    return response.data.data;
  }

  // Get Single Instagram Channel
  async getChannel(id: string): Promise<InstagramChannel> {
    const response = await apiClient.get<ApiResponse<InstagramChannel>>(
      `${this.endpoint}/${id}`
    );
    return response.data.data;
  }

  // Update Instagram Channel
  async updateChannel(
    id: string,
    data: UpdateInstagramChannelRequest
  ): Promise<InstagramChannel> {
    const response = await apiClient.put<ApiResponse<InstagramChannel>>(
      `${this.endpoint}/${id}`,
      data
    );
    return response.data.data;
  }

  // Delete Instagram Channel
  async deleteChannel(id: string): Promise<void> {
    await apiClient.delete(`${this.endpoint}/${id}`);
  }

  // Reconnect Instagram Channel (refresh OAuth token)
  async reconnectChannel(
    id: string,
    authCode: string
  ): Promise<InstagramChannel> {
    const response = await apiClient.post<ApiResponse<InstagramChannel>>(
      `${this.endpoint}/${id}/reconnect`,
      { authCode }
    );
    return response.data.data;
  }

  // Get Instagram Channel Stats
  async getChannelStats(id: string): Promise<InstagramChannelStats> {
    const response = await apiClient.get<ApiResponse<InstagramChannelStats>>(
      `${this.endpoint}/${id}/stats`
    );
    return response.data.data;
  }

  // Sync Instagram Channel Data
  async syncChannel(id: string): Promise<InstagramChannel> {
    const response = await apiClient.post<ApiResponse<InstagramChannel>>(
      `${this.endpoint}/${id}/sync`
    );
    return response.data.data;
  }
}

export const instagramChannelService = new InstagramChannelService();
