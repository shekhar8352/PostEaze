import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { youtubeChannelService } from './youtubeChannelService';
import type { GetChannelsParams } from '../types/youtube.types';

export const youtubeChannelKeys = {
  all: ['youtube', 'channels'] as const,
  lists: () => [...youtubeChannelKeys.all, 'list'] as const,
  list: (params?: GetChannelsParams) => [...youtubeChannelKeys.lists(), params] as const,
};

export function useYouTubeChannels(params?: GetChannelsParams) {
  return useQuery({
    queryKey: youtubeChannelKeys.list(params),
    queryFn: () => youtubeChannelService.getChannels(params),
  });
}

export function useCreateYouTubeChannel() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (authCode: string) => youtubeChannelService.createChannel(authCode),
    onSuccess: (data) => {
      void qc.invalidateQueries({ queryKey: youtubeChannelKeys.lists() });
      notifications.show({
        title: 'YouTube connected',
        message: data.channel_name || 'Channel ready for video publishing',
        color: 'green',
      });
    },
    onError: (e: Error) => {
      notifications.show({ title: 'Connection failed', message: e.message, color: 'red' });
    },
  });
}

export function useDeleteYouTubeChannel() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => youtubeChannelService.deleteChannel(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: youtubeChannelKeys.lists() });
    },
  });
}
