// Base Channel Types - Generic types for all social media channels

/**
 * Supported social media providers
 */
export type ChannelProvider = 'instagram' | 'facebook' | 'youtube';

/**
 * Base metadata structure that all channel providers should extend
 * This ensures common fields are available across all channels
 */
export interface BaseChannelMetadata {
  email: string;
  name: string;
  username: string;
  profile_picture_url: string;
  followers_count: number;
  follows_count: number;
  media_count: number;
  last_synced_at: string;
  website: string;
  biography: string;
  [key: string]: any; // Allow provider-specific additional fields
}

/**
 * Base channel interface matching backend API response structure
 * All provider-specific channels should extend this
 */
export interface BaseChannel<TProvider extends ChannelProvider = ChannelProvider, TMetadata extends BaseChannelMetadata = BaseChannelMetadata> {
  channel_id: number;
  channel_name: string;
  provider: TProvider;
  provider_channel_id: string;
  is_active: boolean;
  metadata: TMetadata;
  created_at: string;
  updated_at?: string;
}

/**
 * Frontend-friendly display version with camelCase properties
 * All provider-specific display types should extend this
 */
export interface BaseChannelDisplay<TProvider extends ChannelProvider = ChannelProvider, TMetadata extends BaseChannelMetadata = BaseChannelMetadata> extends BaseChannel<TProvider, TMetadata> {
  // Computed properties for easier component access
  id: string;
  channelName: string;
  isConnected: boolean;
  email: string;
  username: string;
  profilePicture: string;
  followersCount: number;
  followingCount: number;
  mediaCount: number;
  lastSyncedAt: string;
  createdAt: string;
}

/**
 * Base form data for creating a channel
 */
export interface BaseCreateChannelRequest {
  channelName: string;
  email: string;
  website?: string;
  authCode: string;
}

/**
 * Base API payload for creating a channel
 */
export interface BaseCreateChannelPayload {
  code: string;
  channel_name: string;
  metadata: Record<string, any>;
}

/**
 * Base form data for updating a channel
 */
export interface BaseUpdateChannelRequest {
  channelName?: string;
  email?: string;
  website?: string;
}

/**
 * Base channel statistics
 */
export interface BaseChannelStats {
  followersCount: number;
  followingCount: number;
  mediaCount: number;
  engagementRate: number;
  recentPosts: number;
}

/**
 * Query parameters for fetching channels
 */
export interface GetChannelsParams {
  provider?: ChannelProvider;
  isConnected?: boolean;
  limit?: number;
  offset?: number;
  sortBy?: 'createdAt' | 'updatedAt' | 'channelName';
  sortOrder?: 'asc' | 'desc';
  search?: string;
  [key: string]: string | number | boolean | string[] | undefined;
}

/**
 * API Response wrapper for channel lists
 */
export interface GetChannelsResponse<TChannel extends BaseChannel = BaseChannel> {
  channels: TChannel[];
  total: number;
}

/**
 * Generic transformer function type
 */
export type ChannelTransformer<TChannel extends BaseChannel, TDisplay extends BaseChannelDisplay> = (channel: TChannel) => TDisplay;

/**
 * Base transformer function to convert backend channel to frontend display format
 * Can be used as-is or extended for provider-specific transformations
 */
export function transformChannelToDisplay<TChannel extends BaseChannel, TDisplay extends BaseChannelDisplay>(
  channel: TChannel
): TDisplay {
  const transformed = {
    ...channel,
    id: channel.channel_id.toString(),
    channelName: channel.channel_name,
    isConnected: channel.is_active,
    email: channel.metadata.email,
    username: channel.metadata.username,
    profilePicture: channel.metadata.profile_picture_url,
    followersCount: channel.metadata.followers_count,
    followingCount: channel.metadata.follows_count,
    mediaCount: channel.metadata.media_count,
    lastSyncedAt: channel.metadata.last_synced_at,
    createdAt: channel.created_at,
  };
  
  return transformed as unknown as TDisplay;
}

/**
 * Transform array of channels
 */
export function transformChannelsToDisplay<TChannel extends BaseChannel, TDisplay extends BaseChannelDisplay>(
  channels: TChannel[],
  transformer?: ChannelTransformer<TChannel, TDisplay>
): TDisplay[] {
  const transformFn = transformer || transformChannelToDisplay;
  return channels.map(transformFn);
}
