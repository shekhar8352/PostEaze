import { BaseService } from "@/services/base/BaseService";
import apiClient from "@/services/api/client";
import { type ApiResponse } from "@/services/api/types";
import type {
  InstagramChannel,
  CreateInstagramChannelRequest,
  CreateInstagramChannelPayload,
  UpdateInstagramChannelRequest,
  InstagramChannelStats,
  GetChannelsParams,
} from "../types/instagram.types";

class InstagramChannelService extends BaseService {
  constructor() {
    super("/v1/channels/instagram");
  }

  // Create Instagram Channel
  async createChannel(
    data: CreateInstagramChannelRequest
  ): Promise<InstagramChannel> {
    // Transform frontend data to backend payload format
    const payload: CreateInstagramChannelPayload = {
      code: data.authCode,
      channel_name: data.channelName,
      metadata: {
        email: data.email,
      }
    };
    
    const response = await apiClient.post<ApiResponse<InstagramChannel>>(
      `${this.endpoint}/create`,
      payload
    );
    return response.data.data;
  }

  // Get All Instagram Channels with optional filters
  async getChannels(params?: GetChannelsParams): Promise<InstagramChannel[]> {
    const response = await apiClient.get<ApiResponse<{ channels: InstagramChannel[], total: number }>>(
      '/v1/channels',
      {
        params: {
          provider: 'instagram',
          ...params,
        }
      }
    );
    // Backend returns { channels: [], total: 0 }, so we need to extract the channels array
    return response.data.data.channels;
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
