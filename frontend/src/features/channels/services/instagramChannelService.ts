import { BaseChannelService } from "./BaseChannelService";
import type {
  InstagramChannel,
  InstagramChannelDisplay,
  CreateInstagramChannelRequest,
  CreateInstagramChannelPayload,
  UpdateInstagramChannelRequest,
  InstagramChannelStats,
} from "../types/instagram.types";
import { transformChannelToDisplay } from "../types/instagram.types";

/**
 * Instagram Channel Service
 * Extends BaseChannelService with Instagram-specific functionality
 */
class InstagramChannelService extends BaseChannelService<
  InstagramChannel,
  InstagramChannelDisplay,
  CreateInstagramChannelRequest,
  CreateInstagramChannelPayload,
  UpdateInstagramChannelRequest,
  InstagramChannelStats
> {
  constructor() {
    super('instagram', transformChannelToDisplay);
  }

  /**
   * Override prepareCreatePayload if Instagram has specific requirements
   * Otherwise, the base implementation is used
   */
  protected prepareCreatePayload(data: CreateInstagramChannelRequest): CreateInstagramChannelPayload {
    return {
      code: data.authCode,
      channel_name: data.channelName,
      metadata: {
        email: data.email,
        website: data.website,
      }
    };
  }

  // Add Instagram-specific methods here if needed
  // All common methods (create, get, update, delete, reconnect, stats, sync)
  // are inherited from BaseChannelService
}

export const instagramChannelService = new InstagramChannelService();
