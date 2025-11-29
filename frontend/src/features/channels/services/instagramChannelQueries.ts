import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { instagramChannelService } from "./instagramChannelService";
import type {
  CreateInstagramChannelRequest,
  UpdateInstagramChannelRequest,
} from "../types/instagram.types";

// Query Keys
export const instagramChannelKeys = {
  all: ['instagram', 'channels'] as const,
  lists: () => [...instagramChannelKeys.all, 'list'] as const,
  list: (filters: string) => [...instagramChannelKeys.lists(), { filters }] as const,
  details: () => [...instagramChannelKeys.all, 'detail'] as const,
  detail: (id: string) => [...instagramChannelKeys.details(), id] as const,
  stats: (id: string) => [...instagramChannelKeys.all, 'stats', id] as const,
};

// Get All Instagram Channels
export const useInstagramChannels = () => {
  return useQuery({
    queryKey: instagramChannelKeys.lists(),
    queryFn: () => instagramChannelService.getChannels(),
  });
};

// Get Single Instagram Channel
export const useInstagramChannel = (id: string) => {
  return useQuery({
    queryKey: instagramChannelKeys.detail(id),
    queryFn: () => instagramChannelService.getChannel(id),
    enabled: !!id,
  });
};

// Get Instagram Channel Stats
export const useInstagramChannelStats = (id: string) => {
  return useQuery({
    queryKey: instagramChannelKeys.stats(id),
    queryFn: () => instagramChannelService.getChannelStats(id),
    enabled: !!id,
  });
};

// Create Instagram Channel
export const useCreateInstagramChannel = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateInstagramChannelRequest) =>
      instagramChannelService.createChannel(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.lists() });
    },
    onError: (error: any) => {
      console.error("Create Instagram channel failed:", error.message);
    },
  });
};

// Update Instagram Channel
export const useUpdateInstagramChannel = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateInstagramChannelRequest }) =>
      instagramChannelService.updateChannel(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.lists() });
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.detail(variables.id) });
    },
    onError: (error: any) => {
      console.error("Update Instagram channel failed:", error.message);
    },
  });
};

// Delete Instagram Channel
export const useDeleteInstagramChannel = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => instagramChannelService.deleteChannel(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.lists() });
    },
    onError: (error: any) => {
      console.error("Delete Instagram channel failed:", error.message);
    },
  });
};

// Reconnect Instagram Channel
export const useReconnectInstagramChannel = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, authCode }: { id: string; authCode: string }) =>
      instagramChannelService.reconnectChannel(id, authCode),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.lists() });
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.detail(variables.id) });
    },
    onError: (error: any) => {
      console.error("Reconnect Instagram channel failed:", error.message);
    },
  });
};

// Sync Instagram Channel
export const useSyncInstagramChannel = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => instagramChannelService.syncChannel(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.lists() });
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: instagramChannelKeys.stats(id) });
    },
    onError: (error: any) => {
      console.error("Sync Instagram channel failed:", error.message);
    },
  });
};
