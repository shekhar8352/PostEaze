import { BaseService } from "@/services/base/BaseService";
import apiClient from "@/services/api/client";
import { type ApiResponse } from "@/services/api/types";
import type {
  BaseChannel,
  BaseChannelDisplay,
  BaseCreateChannelRequest,
  BaseCreateChannelPayload,
  BaseUpdateChannelRequest,
  BaseChannelStats,
  GetChannelsParams,
  GetChannelsResponse,
  ChannelProvider,
} from "../types/base.types";

/**
 * ============================================================================
 * BASE CHANNEL SERVICE
 * ============================================================================
 * 
 * This is a reusable service class that handles common operations for ALL 
 * social media channels (Instagram, Facebook, YouTube, etc.)
 * 
 * WHY USE THIS?
 * Instead of writing the same code for each channel type, we write it once here
 * and let each channel (Instagram, Facebook) extend this class.
 * 
 * WHAT ARE THOSE <T...> THINGS?
 * Those are TypeScript "generics" - think of them as placeholders that get 
 * replaced with actual types when you use this class.
 * 
 * EXAMPLE:
 * When InstagramChannelService extends this class, it says:
 * - BackendData = InstagramChannel (what the API returns)
 * - FrontendData = InstagramChannelDisplay (what the UI shows)
 * - CreateData = InstagramCreateRequest (form data to create a channel)
 * - etc.
 * 
 * ============================================================================
 */
export abstract class BaseChannelService<
  // What the backend API sends us (raw data from server)
  BackendData extends BaseChannel,
  
  // What the frontend UI displays (transformed, user-friendly data)
  FrontendData extends BaseChannelDisplay,
  
  // Data we collect from user to create a new channel
  CreateData extends BaseCreateChannelRequest,
  
  // Data we send to backend API to create a channel
  CreatePayload extends BaseCreateChannelPayload,
  
  // Data we collect from user to update a channel
  UpdateData extends BaseUpdateChannelRequest,
  
  // Statistics data for a channel (followers, posts, etc.)
  StatsData extends BaseChannelStats
> extends BaseService {
  
  // ========================================================================
  // PROPERTIES
  // ========================================================================
  
  /**
   * Which social media platform this service is for
   * Examples: 'instagram', 'facebook', 'youtube'
   */
  protected provider: ChannelProvider;
  
  /**
   * A function that converts backend data to frontend data
   * Think of it as a translator between what the API gives us 
   * and what the UI needs to display
   */
  protected dataTransformer: (backendChannel: BackendData) => FrontendData;

  // ========================================================================
  // CONSTRUCTOR
  // ========================================================================
  
  /**
   * Sets up the service when it's created
   * 
   * @param provider - Which platform (e.g., 'instagram')
   * @param dataTransformer - Function to convert backend data to frontend format
   * 
   * EXAMPLE:
   * new InstagramChannelService('instagram', transformInstagramChannel)
   */
  constructor(
    provider: ChannelProvider,
    dataTransformer: (backendChannel: BackendData) => FrontendData
  ) {
    // Set up the base API endpoint (e.g., /v1/channels/instagram)
    super(`/v1/channels/${provider}`);
    
    this.provider = provider;
    this.dataTransformer = dataTransformer;
  }

  // ========================================================================
  // HELPER METHODS (used internally by this class)
  // ========================================================================
  
  /**
   * Converts an array of backend channels to frontend format
   * 
   * WHY? Backend sends data in one format, but our UI needs it differently
   * 
   * @param backendChannels - Array of channels from the API
   * @returns Array of channels ready for the UI to display
   */
  protected convertToFrontendFormat(backendChannels: BackendData[]): FrontendData[] {
    // Apply the transformer function to each channel
    return backendChannels.map(this.dataTransformer);
  }

  /**
   * Prepares the data to send to the backend when creating a channel
   * 
   * WHY? The form collects data in one format, but the API expects it differently
   * 
   * @param formData - Data from the create channel form
   * @returns Payload formatted for the backend API
   * 
   * NOTE: Child classes (like InstagramChannelService) can override this
   * if they need custom payload formatting
   */
  protected buildCreatePayload(formData: CreateData): CreatePayload {
    // Build the standard payload structure
    const payload = {
      code: formData.authCode,              // OAuth authorization code
      channel_name: formData.channelName,   // User-provided channel name
      metadata: {
        email: formData.email,              // Optional: user's email
        website: formData.website,          // Optional: user's website
      }
    };
    
    // TypeScript type casting (we know this matches CreatePayload structure)
    return payload as unknown as CreatePayload;
  }

  // ========================================================================
  // PUBLIC API METHODS (used by components and hooks)
  // ========================================================================
  
  /**
   * CREATE: Add a new channel
   * 
   * FLOW:
   * 1. User fills out form → formData
   * 2. We convert formData to API payload
   * 3. Send POST request to backend
   * 4. Return the created channel
   * 
   * @param formData - Data from the create channel form
   * @returns The newly created channel (backend format)
   * 
   * EXAMPLE USAGE:
   * const newChannel = await instagramService.createChannel({
   *   authCode: 'abc123',
   *   channelName: 'My Instagram',
   *   email: 'user@example.com'
   * });
   */
  async createChannel(formData: CreateData): Promise<BackendData> {
    // Step 1: Convert form data to API payload format
    const apiPayload = this.buildCreatePayload(formData);
    
    // Step 2: Send POST request to create the channel
    const response = await apiClient.post<ApiResponse<BackendData>>(
      `${this.endpoint}/create`,
      apiPayload
    );
    
    // Step 3: Extract and return the channel data
    return response.data.data;
  }

  /**
   * READ: Get all channels for this provider
   * 
   * FLOW:
   * 1. Send GET request with optional filters
   * 2. Receive array of channels from backend
   * 3. Convert to frontend format for UI display
   * 
   * @param filters - Optional filters (status, search, pagination)
   * @returns Array of channels ready for UI display
   * 
   * EXAMPLE USAGE:
   * const channels = await instagramService.getChannels({ 
   *   status: 'active',
   *   page: 1 
   * });
   */
  async getChannels(filters?: GetChannelsParams): Promise<FrontendData[]> {
    // Step 1: Send GET request with provider and filters
    const response = await apiClient.get<ApiResponse<GetChannelsResponse<BackendData>>>(
      '/v1/channels',
      {
        params: {
          provider: this.provider,  // e.g., 'instagram'
          ...filters,               // spread any additional filters
        }
      }
    );
    
    // Step 2: Extract channels from response
    const backendChannels = response.data.data.channels;
    
    // Step 3: Convert to frontend format and return
    return this.convertToFrontendFormat(backendChannels);
  }

  /**
   * READ: Get a single channel by its ID
   * 
   * @param channelId - Unique identifier for the channel
   * @returns The channel data (backend format)
   * 
   * EXAMPLE USAGE:
   * const channel = await instagramService.getChannel('channel_123');
   */
  async getChannel(channelId: string): Promise<BackendData> {
    const response = await apiClient.get<ApiResponse<BackendData>>(
      `${this.endpoint}/${channelId}`
    );
    return response.data.data;
  }

  /**
   * UPDATE: Modify an existing channel
   * 
   * @param channelId - ID of the channel to update
   * @param updateData - New data to apply
   * @returns The updated channel
   * 
   * EXAMPLE USAGE:
   * const updated = await instagramService.updateChannel('channel_123', {
   *   channelName: 'New Name'
   * });
   */
  async updateChannel(channelId: string, updateData: UpdateData): Promise<BackendData> {
    const response = await apiClient.put<ApiResponse<BackendData>>(
      `${this.endpoint}/${channelId}`,
      updateData
    );
    return response.data.data;
  }

  /**
   * DELETE: Remove a channel
   * 
   * @param channelId - ID of the channel to delete
   * 
   * EXAMPLE USAGE:
   * await instagramService.deleteChannel('channel_123');
   */
  async deleteChannel(channelId: string): Promise<void> {
    await apiClient.delete(`${this.endpoint}/${channelId}`);
  }

  /**
   * RECONNECT: Refresh the OAuth connection for a channel
   * 
   * WHY? OAuth tokens expire, so users need to reconnect periodically
   * 
   * @param channelId - ID of the channel to reconnect
   * @param newAuthCode - New OAuth authorization code
   * @returns The updated channel with fresh token
   * 
   * EXAMPLE USAGE:
   * const reconnected = await instagramService.reconnectChannel(
   *   'channel_123', 
   *   'new_auth_code_xyz'
   * );
   */
  async reconnectChannel(channelId: string, newAuthCode: string): Promise<BackendData> {
    const response = await apiClient.post<ApiResponse<BackendData>>(
      `${this.endpoint}/${channelId}/reconnect`,
      { authCode: newAuthCode }
    );
    return response.data.data;
  }

  /**
   * STATS: Get statistics for a channel
   * 
   * @param channelId - ID of the channel
   * @returns Statistics (followers, posts, engagement, etc.)
   * 
   * EXAMPLE USAGE:
   * const stats = await instagramService.getChannelStats('channel_123');
   * console.log(stats.followers, stats.posts);
   */
  async getChannelStats(channelId: string): Promise<StatsData> {
    const response = await apiClient.get<ApiResponse<StatsData>>(
      `${this.endpoint}/${channelId}/stats`
    );
    return response.data.data;
  }

  /**
   * SYNC: Fetch latest data from the social media platform
   * 
   * WHY? Channel data (followers, profile pic, etc.) changes on the platform
   * This refreshes our local copy with the latest data
   * 
   * @param channelId - ID of the channel to sync
   * @returns The updated channel with fresh data
   * 
   * EXAMPLE USAGE:
   * const synced = await instagramService.syncChannel('channel_123');
   */
  async syncChannel(channelId: string): Promise<BackendData> {
    const response = await apiClient.post<ApiResponse<BackendData>>(
      `${this.endpoint}/${channelId}/sync`
    );
    return response.data.data;
  }
}
