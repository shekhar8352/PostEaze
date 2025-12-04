// Instagram Channel Types - Extends base channel types

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
 * Instagram-specific metadata extending base metadata
 * Add Instagram-specific fields here if needed in the future
 */
export interface InstagramChannelMetadata extends BaseChannelMetadata {
  // Instagram-specific fields can be added here
  // For now, it uses all base fields
}

/**
 * Instagram channel interface - extends base channel with Instagram provider type
 */
export interface InstagramChannel extends BaseChannel<'instagram', InstagramChannelMetadata> {}

/**
 * Instagram channel display interface - frontend-friendly version
 */
export interface InstagramChannelDisplay extends BaseChannelDisplay<'instagram', InstagramChannelMetadata> {}

/**
 * Frontend form data for creating Instagram channel
 */
export interface CreateInstagramChannelRequest extends BaseCreateChannelRequest {}

/**
 * Backend API payload for creating Instagram channel
 */
export interface CreateInstagramChannelPayload extends BaseCreateChannelPayload {}

/**
 * Update Instagram channel request
 */
export interface UpdateInstagramChannelRequest extends BaseUpdateChannelRequest {}

/**
 * Instagram channel form data
 */
export interface InstagramChannelFormData {
  channelName: string;
  email: string;
  website: string;
}

/**
 * Instagram OAuth configuration
 */
export interface InstagramOAuthConfig {
  clientId: string;
  redirectUri: string;
  scope: string;
  responseType: string;
}

/**
 * Instagram channel statistics - extends base stats
 */
export interface InstagramChannelStats extends BaseChannelStats {}

// Re-export common types and utilities from base for convenience
export type { GetChannelsParams, GetChannelsResponse } from './base.types';

/**
 * Transform Instagram channel to display format
 * Uses the base transformer
 */
export function transformChannelToDisplay(channel: InstagramChannel): InstagramChannelDisplay {
  return baseTransformChannelToDisplay<InstagramChannel, InstagramChannelDisplay>(channel);
}

/**
 * Transform array of Instagram channels to display format
 */
export function transformChannelsToDisplay(channels: InstagramChannel[]): InstagramChannelDisplay[] {
  return baseTransformChannelsToDisplay<InstagramChannel, InstagramChannelDisplay>(channels);
}

