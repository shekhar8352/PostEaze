// Instagram Channel Types

export interface InstagramChannel {
  id: string;
  channelName: string;
  email: string;
  website?: string;
  instagramUsername?: string;
  instagramUserId?: string;
  profilePicture?: string;
  followersCount?: number;
  followingCount?: number;
  mediaCount?: number;
  isConnected: boolean;
  lastSyncedAt?: string;
  createdAt: string;
  updatedAt: string;
}

// Frontend form data
export interface CreateInstagramChannelRequest {
  channelName: string;
  email: string;
  website?: string;
  authCode: string;
}

// Backend API payload
export interface CreateInstagramChannelPayload {
  code: string;
  channel_name: string;
}

export interface UpdateInstagramChannelRequest {
  channelName?: string;
  email?: string;
  website?: string;
}

export interface InstagramChannelFormData {
  channelName: string;
  email: string;
  website: string;
}

export interface InstagramOAuthConfig {
  clientId: string;
  redirectUri: string;
  scope: string;
  responseType: string;
}

export interface InstagramChannelStats {
  followersCount: number;
  followingCount: number;
  mediaCount: number;
  engagementRate: number;
  recentPosts: number;
}
