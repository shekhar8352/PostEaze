// Facebook Channel Types - Extends base channel types
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
 * Facebook-specific metadata extending base metadata
 * Add Facebook-specific fields here
 */
export interface FacebookChannelMetadata extends BaseChannelMetadata {
  // Facebook-specific fields
  page_id?: string;
  page_category?: string;
  page_likes?: number;
  // Add more Facebook-specific fields as needed
}

/**
 * Facebook channel interface - extends base channel with Facebook provider type
 */
export interface FacebookChannel extends BaseChannel<'facebook', FacebookChannelMetadata> {}

/**
 * Facebook channel display interface - frontend-friendly version
 */
export interface FacebookChannelDisplay extends BaseChannelDisplay<'facebook', FacebookChannelMetadata> {
  // Add any Facebook-specific display properties if needed
  pageId?: string;
  pageLikes?: number;
}

/**
 * Frontend form data for creating Facebook channel
 */
export interface CreateFacebookChannelRequest extends BaseCreateChannelRequest {}

/**
 * Backend API payload for creating Facebook channel
 */
export interface CreateFacebookChannelPayload extends BaseCreateChannelPayload {}

/**
 * Update Facebook channel request
 */
export interface UpdateFacebookChannelRequest extends BaseUpdateChannelRequest {}

/**
 * Facebook channel form data
 */
export interface FacebookChannelFormData {
  channelName: string;
  email: string;
  website: string;
}

/**
 * Facebook OAuth configuration
 */
export interface FacebookOAuthConfig {
  clientId: string;
  redirectUri: string;
  scope: string;
  responseType: string;
}

/**
 * Facebook channel statistics - extends base stats
 */
export interface FacebookChannelStats extends BaseChannelStats {
  // Add Facebook-specific stats if needed
  pageLikes?: number;
  pageViews?: number;
}

// Re-export common types and utilities from base for convenience
export type { GetChannelsParams, GetChannelsResponse } from './base.types';

/**
 * Transform Facebook channel to display format
 * Extends base transformer with Facebook-specific fields
 */
export function transformChannelToDisplay(channel: FacebookChannel): FacebookChannelDisplay {
  const baseDisplay = baseTransformChannelToDisplay<FacebookChannel, FacebookChannelDisplay>(channel);
  
  // Add Facebook-specific transformations
  return {
    ...baseDisplay,
    pageId: channel.metadata.page_id,
    pageLikes: channel.metadata.page_likes,
  };
}

/**
 * Transform array of Facebook channels to display format
 */
export function transformChannelsToDisplay(channels: FacebookChannel[]): FacebookChannelDisplay[] {
  return baseTransformChannelsToDisplay<FacebookChannel, FacebookChannelDisplay>(
    channels,
    transformChannelToDisplay
  );
}
