// YouTube Channel Types - Extends base channel types
// This is a template showing how to add a new channel provider

import type {
  BaseChannel,
  BaseChannelDisplay,
  BaseChannelMetadata,
  BaseCreateChannelRequest,
  BaseCreateChannelPayload,
  BaseUpdateChannelRequest,
  BaseChannelStats,
} from './base.types';

import {
  transformChannelToDisplay as baseTransformChannelToDisplay,
  transformChannelsToDisplay as baseTransformChannelsToDisplay,
} from './base.types';

/**
 * YouTube-specific metadata extending base metadata
 * Add YouTube-specific fields here
 */
export interface YouTubeChannelMetadata extends BaseChannelMetadata {
  // YouTube-specific fields
  channel_url?: string;
  subscriber_count?: number;
  video_count?: number;
  view_count?: number;
  // Add more YouTube-specific fields as needed
}

/**
 * YouTube channel interface - extends base channel with YouTube provider type
 */
export interface YouTubeChannel extends BaseChannel<'youtube', YouTubeChannelMetadata> {}

/**
 * YouTube channel display interface - frontend-friendly version
 */
export interface YouTubeChannelDisplay extends BaseChannelDisplay<'youtube', YouTubeChannelMetadata> {
  // Add any YouTube-specific display properties if needed
  channelUrl?: string;
  subscriberCount?: number;
  videoCount?: number;
  viewCount?: number;
}

/**
 * Frontend form data for creating YouTube channel
 */
export interface CreateYouTubeChannelRequest extends BaseCreateChannelRequest {}

/**
 * Backend API payload for creating YouTube channel
 */
export interface CreateYouTubeChannelPayload extends BaseCreateChannelPayload {}

/**
 * Update YouTube channel request
 */
export interface UpdateYouTubeChannelRequest extends BaseUpdateChannelRequest {}

/**
 * YouTube channel form data
 */
export interface YouTubeChannelFormData {
  channelName: string;
  email: string;
  website: string;
}

/**
 * YouTube OAuth configuration
 */
export interface YouTubeOAuthConfig {
  clientId: string;
  redirectUri: string;
  scope: string;
  responseType: string;
}

/**
 * YouTube channel statistics - extends base stats
 */
export interface YouTubeChannelStats extends BaseChannelStats {
  // Add YouTube-specific stats if needed
  subscriberCount?: number;
  videoCount?: number;
  viewCount?: number;
  averageViewDuration?: number;
}

// Re-export common types and utilities from base for convenience
export type { GetChannelsParams, GetChannelsResponse } from './base.types';

/**
 * Transform YouTube channel to display format
 * Extends base transformer with YouTube-specific fields
 */
export function transformChannelToDisplay(channel: YouTubeChannel): YouTubeChannelDisplay {
  const baseDisplay = baseTransformChannelToDisplay<YouTubeChannel, YouTubeChannelDisplay>(channel);
  
  // Add YouTube-specific transformations
  return {
    ...baseDisplay,
    channelUrl: channel.metadata.channel_url,
    subscriberCount: channel.metadata.subscriber_count,
    videoCount: channel.metadata.video_count,
    viewCount: channel.metadata.view_count,
  };
}

/**
 * Transform array of YouTube channels to display format
 */
export function transformChannelsToDisplay(channels: YouTubeChannel[]): YouTubeChannelDisplay[] {
  return baseTransformChannelsToDisplay<YouTubeChannel, YouTubeChannelDisplay>(
    channels,
    transformChannelToDisplay
  );
}
