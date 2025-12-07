# Quick Start: Adding a New Channel Provider

## 🚀 5-Minute Guide

Follow these steps to add a new social media provider (e.g., Facebook, YouTube, Twitter, TikTok, etc.)

### 1️⃣ Create Type File

**File**: `types/[provider].types.ts`

```typescript
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

// Replace 'Provider' with your provider name (e.g., Facebook, YouTube)
// Replace 'provider' with lowercase version

export interface ProviderChannelMetadata extends BaseChannelMetadata {
  // Add provider-specific fields here
}

export interface ProviderChannel extends BaseChannel<'provider', ProviderChannelMetadata> {}

export interface ProviderChannelDisplay extends BaseChannelDisplay<'provider', ProviderChannelMetadata> {}

export interface CreateProviderChannelRequest extends BaseCreateChannelRequest {}
export interface CreateProviderChannelPayload extends BaseCreateChannelPayload {}
export interface UpdateProviderChannelRequest extends BaseUpdateChannelRequest {}

export interface ProviderChannelStats extends BaseChannelStats {}

export function transformChannelToDisplay(channel: ProviderChannel): ProviderChannelDisplay {
  return baseTransformChannelToDisplay<ProviderChannel, ProviderChannelDisplay>(channel);
}

export function transformChannelsToDisplay(channels: ProviderChannel[]): ProviderChannelDisplay[] {
  return baseTransformChannelsToDisplay<ProviderChannel, ProviderChannelDisplay>(channels);
}
```

### 2️⃣ Create Service File

**File**: `services/[provider]ChannelService.ts`

```typescript
import { BaseChannelService } from "./BaseChannelService";
import type {
  ProviderChannel,
  ProviderChannelDisplay,
  CreateProviderChannelRequest,
  CreateProviderChannelPayload,
  UpdateProviderChannelRequest,
  ProviderChannelStats,
} from "../types/provider.types";
import { transformChannelToDisplay } from "../types/provider.types";

class ProviderChannelService extends BaseChannelService<
  ProviderChannel,
  ProviderChannelDisplay,
  CreateProviderChannelRequest,
  CreateProviderChannelPayload,
  UpdateProviderChannelRequest,
  ProviderChannelStats
> {
  constructor() {
    super('provider', transformChannelToDisplay);
  }
}

export const providerChannelService = new ProviderChannelService();
```

### 3️⃣ Update Base Types (if needed)

**File**: `types/base.types.ts`

Add your provider to the `ChannelProvider` type:

```typescript
export type ChannelProvider = 'instagram' | 'facebook' | 'youtube' | 'provider';
```

### 4️⃣ Done! ✅

Your new provider now has all these methods:
- `createChannel(data)`
- `getChannels(params?)`
- `getChannel(id)`
- `updateChannel(id, data)`
- `deleteChannel(id)`
- `reconnectChannel(id, authCode)`
- `getChannelStats(id)`
- `syncChannel(id)`

## 📋 Checklist

- [ ] Created `types/[provider].types.ts`
- [ ] Created `services/[provider]ChannelService.ts`
- [ ] Added provider to `ChannelProvider` type
- [ ] Tested basic CRUD operations
- [ ] Created provider-specific components (optional)

## 🎯 Real Example: Facebook

### types/facebook.types.ts
```typescript
export interface FacebookChannelMetadata extends BaseChannelMetadata {
  page_id?: string;
  page_likes?: number;
}

export interface FacebookChannel extends BaseChannel<'facebook', FacebookChannelMetadata> {}
export interface FacebookChannelDisplay extends BaseChannelDisplay<'facebook', FacebookChannelMetadata> {
  pageId?: string;
  pageLikes?: number;
}

// ... rest of the types
```

### services/facebookChannelService.ts
```typescript
class FacebookChannelService extends BaseChannelService<
  FacebookChannel,
  FacebookChannelDisplay,
  CreateFacebookChannelRequest,
  CreateFacebookChannelPayload,
  UpdateFacebookChannelRequest,
  FacebookChannelStats
> {
  constructor() {
    super('facebook', transformChannelToDisplay);
  }
}

export const facebookChannelService = new FacebookChannelService();
```

## 💡 Tips

1. **Start simple**: Use the base types as-is, add provider-specific fields later
2. **Copy from Instagram**: Use `instagram.types.ts` as a reference
3. **Test incrementally**: Test each method as you add it
4. **Use TypeScript**: Let IntelliSense guide you

## 🔧 Common Customizations

### Custom Metadata Fields
```typescript
export interface TwitterChannelMetadata extends BaseChannelMetadata {
  tweet_count?: number;
  verified?: boolean;
}
```

### Custom Display Properties
```typescript
export interface TwitterChannelDisplay extends BaseChannelDisplay<'twitter', TwitterChannelMetadata> {
  tweetCount?: number;
  isVerified?: boolean;
}
```

### Custom Transformer
```typescript
export function transformChannelToDisplay(channel: TwitterChannel): TwitterChannelDisplay {
  const baseDisplay = baseTransformChannelToDisplay<TwitterChannel, TwitterChannelDisplay>(channel);
  
  return {
    ...baseDisplay,
    tweetCount: channel.metadata.tweet_count,
    isVerified: channel.metadata.verified,
  };
}
```

### Custom Service Method
```typescript
class TwitterChannelService extends BaseChannelService<...> {
  constructor() {
    super('twitter', transformChannelToDisplay);
  }

  // Add custom method
  async getTweets(channelId: string, limit: number = 10) {
    const response = await apiClient.get(`${this.endpoint}/${channelId}/tweets`, {
      params: { limit }
    });
    return response.data.data;
  }
}
```

## 📚 Next Steps

1. Create provider-specific components (Card, List, Form)
2. Add provider-specific OAuth configuration
3. Create provider-specific validation schemas
4. Add provider-specific routes and pages

## 🆘 Need Help?

- Check `ARCHITECTURE.md` for detailed documentation
- Look at `instagram.types.ts` for a complete example
- Review `BaseChannelService.ts` for available methods
