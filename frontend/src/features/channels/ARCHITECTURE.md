# Channel Management - Modular Architecture Guide

This guide explains the modular architecture for managing social media channels (Instagram, Facebook, YouTube, etc.) in PostEaze.

## 📁 Architecture Overview

The channel management system is built with a **generic base layer** that can be extended for any social media provider with minimal code duplication.

```
features/channels/
├── types/
│   ├── base.types.ts           # Generic base types for all providers
│   ├── instagram.types.ts      # Instagram-specific types (extends base)
│   ├── facebook.types.ts       # Facebook-specific types (extends base)
│   └── youtube.types.ts        # YouTube-specific types (extends base)
├── services/
│   ├── BaseChannelService.ts   # Generic service class
│   └── instagramChannelService.ts  # Instagram service (extends base)
└── components/
    └── [provider-specific components]
```

## 🎯 Key Principles

1. **DRY (Don't Repeat Yourself)**: Common functionality is in base classes
2. **Type Safety**: Full TypeScript support with generics
3. **Extensibility**: Easy to add new providers
4. **Consistency**: All providers follow the same patterns

## 🔧 Adding a New Channel Provider

### Step 1: Create Type Definitions

Create `types/[provider].types.ts`:

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

// 1. Define provider-specific metadata
export interface FacebookChannelMetadata extends BaseChannelMetadata {
  // Add provider-specific fields
  page_id?: string;
  page_likes?: number;
}

// 2. Define channel type
export interface FacebookChannel extends BaseChannel<'facebook', FacebookChannelMetadata> {}

// 3. Define display type
export interface FacebookChannelDisplay extends BaseChannelDisplay<'facebook', FacebookChannelMetadata> {
  // Add computed properties for easier access
  pageId?: string;
  pageLikes?: number;
}

// 4. Extend request types
export interface CreateFacebookChannelRequest extends BaseCreateChannelRequest {}
export interface UpdateFacebookChannelRequest extends BaseUpdateChannelRequest {}

// 5. Define stats type
export interface FacebookChannelStats extends BaseChannelStats {
  pageLikes?: number;
}

// 6. Create transformer functions
export function transformChannelToDisplay(channel: FacebookChannel): FacebookChannelDisplay {
  const baseDisplay = baseTransformChannelToDisplay<FacebookChannel, FacebookChannelDisplay>(channel);
  
  return {
    ...baseDisplay,
    pageId: channel.metadata.page_id,
    pageLikes: channel.metadata.page_likes,
  };
}

export function transformChannelsToDisplay(channels: FacebookChannel[]): FacebookChannelDisplay[] {
  return baseTransformChannelsToDisplay<FacebookChannel, FacebookChannelDisplay>(
    channels,
    transformChannelToDisplay
  );
}
```

### Step 2: Create Service Class

Create `services/[provider]ChannelService.ts`:

```typescript
import { BaseChannelService } from "./BaseChannelService";
import type {
  FacebookChannel,
  FacebookChannelDisplay,
  CreateFacebookChannelRequest,
  CreateFacebookChannelPayload,
  UpdateFacebookChannelRequest,
  FacebookChannelStats,
} from "../types/facebook.types";
import { transformChannelToDisplay } from "../types/facebook.types";

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

  // Override methods only if provider has specific requirements
  // Otherwise, all methods are inherited from BaseChannelService
  
  // Example: Custom create payload
  protected prepareCreatePayload(data: CreateFacebookChannelRequest): CreateFacebookChannelPayload {
    return {
      code: data.authCode,
      channel_name: data.channelName,
      metadata: {
        email: data.email,
        website: data.website,
        // Add Facebook-specific fields
      }
    };
  }
}

export const facebookChannelService = new FacebookChannelService();
```

### Step 3: That's It! 🎉

Your new provider now has all these methods automatically:
- ✅ `createChannel()`
- ✅ `getChannels()`
- ✅ `getChannel()`
- ✅ `updateChannel()`
- ✅ `deleteChannel()`
- ✅ `reconnectChannel()`
- ✅ `getChannelStats()`
- ✅ `syncChannel()`

## 📊 What's Included in Base Types

### BaseChannel
```typescript
{
  channel_id: number;
  channel_name: string;
  provider: 'instagram' | 'facebook' | 'youtube';
  provider_channel_id: string;
  is_active: boolean;
  metadata: BaseChannelMetadata;
  created_at: string;
  updated_at?: string;
}
```

### BaseChannelMetadata
```typescript
{
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
  [key: string]: any; // Extensible for provider-specific fields
}
```

### BaseChannelDisplay
Frontend-friendly version with camelCase properties:
```typescript
{
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
  // ... plus all BaseChannel properties
}
```

## 🔄 Data Flow

```
Backend API (snake_case)
    ↓
BaseChannel<Provider, Metadata>
    ↓
transformChannelToDisplay()
    ↓
BaseChannelDisplay (camelCase)
    ↓
React Components
```

## 💡 Benefits

### Before (Without Modular Architecture)
- ❌ ~110 lines per provider service
- ❌ Duplicate code across providers
- ❌ Inconsistent implementations
- ❌ Hard to maintain

### After (With Modular Architecture)
- ✅ ~30 lines per provider service
- ✅ Shared base implementation
- ✅ Consistent patterns
- ✅ Easy to maintain and extend

## 🎨 Example: Instagram Implementation

The Instagram implementation demonstrates the pattern:

**Types** (`instagram.types.ts`): ~90 lines
- Extends base types
- Adds Instagram-specific fields (if any)
- Provides transformer functions

**Service** (`instagramChannelService.ts`): ~45 lines
- Extends BaseChannelService
- Overrides only what's needed
- Inherits all common methods

**Total**: ~135 lines vs ~250+ lines without modularity

## 🚀 Future Providers

Adding Facebook or YouTube now requires:
1. Copy template files
2. Update provider name
3. Add provider-specific fields
4. Done! (~30 minutes of work)

## 📝 Best Practices

1. **Always extend base types** - Don't create standalone types
2. **Use transformers** - Keep components clean with camelCase
3. **Override sparingly** - Only override when provider needs custom logic
4. **Document differences** - Comment why a provider is different
5. **Keep metadata flexible** - Use `[key: string]: any` for extensibility

## 🔍 Type Safety

The architecture provides full type safety:
- ✅ Generic constraints ensure type compatibility
- ✅ TypeScript catches errors at compile time
- ✅ IntelliSense works perfectly
- ✅ Refactoring is safe and easy

## 📚 Related Files

- `types/base.types.ts` - Base type definitions
- `services/BaseChannelService.ts` - Base service implementation
- `types/instagram.types.ts` - Instagram example
- `types/facebook.types.ts` - Facebook template
- `types/youtube.types.ts` - YouTube template
