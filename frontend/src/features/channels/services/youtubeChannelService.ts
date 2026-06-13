import apiClient from '@/services/api/client';
import type { ApiResponse } from '@/services/api/types';
import type {
  YouTubeChannel,
  YouTubeChannelDisplay,
  GetChannelsParams,
} from '../types/youtube.types';
import { transformChannelToDisplay } from '../types/youtube.types';
import { getGoogleRedirectUri } from '@/features/integrations/utils/googleOAuth.utils';

export interface CreateYouTubeChannelResponse {
  channel_id: number;
  channel_name: string;
}

class YouTubeChannelService {
  async createChannel(authCode: string): Promise<CreateYouTubeChannelResponse> {
    const response = await apiClient.post<ApiResponse<CreateYouTubeChannelResponse>>(
      '/v1/channels/youtube/create',
      {
        code: authCode,
        redirect_uri: getGoogleRedirectUri('youtube'),
      }
    );
    return response.data.data;
  }

  async getChannels(params?: GetChannelsParams): Promise<YouTubeChannelDisplay[]> {
    const response = await apiClient.get<ApiResponse<{ channels: YouTubeChannel[] }>>(
      '/v1/channels',
      { params: { provider: 'youtube', ...params } }
    );
    return (response.data.data.channels ?? []).map(transformChannelToDisplay);
  }

  async deleteChannel(channelId: string): Promise<void> {
    await apiClient.delete(`/v1/channels/youtube/${channelId}`);
  }
}

export const youtubeChannelService = new YouTubeChannelService();
