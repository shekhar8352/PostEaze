import { useQuery } from "@tanstack/react-query";
import apiClient from "@/services/api/client";
import { type ApiResponse } from "@/services/api/types";
import { 
  type BaseChannel, 
  type GetChannelsResponse,
  transformChannelsToDisplay
} from "../types/base.types";

export const channelKeys = {
  all: ['channels'] as const,
  lists: () => [...channelKeys.all, 'list'] as const,
};

/**
 * Hook to fetch all channels for the authenticated user
 */
export const useChannels = () => {
  return useQuery({
    queryKey: channelKeys.lists(),
    queryFn: async () => {
      const response = await apiClient.get<ApiResponse<GetChannelsResponse<BaseChannel>>>(
        '/v1/channels'
      );
      
      const backendChannels = response.data.data.channels;
      return transformChannelsToDisplay(backendChannels);
    },
  });
};
