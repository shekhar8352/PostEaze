// Barrel exports for channels feature
export { default as channelRoutes } from './channelRoutes';

// Components
export { InstagramChannelCard } from './components/InstagramChannelCard';
export { InstagramChannelList } from './components/InstagramChannelList';
export { InstagramChannelForm } from './components/InstagramChannelForm';
export { InstagramChannelModal } from './components/InstagramChannelModal';
export { default as InstagramOAuthCallback } from './components/InstagramOAuthCallback';

// Pages
export { default as InstagramChannelPage } from './pages/InstagramChannelPage';

// Services
export { instagramChannelService } from './services/instagramChannelService';
export * from './services/instagramChannelQueries';
export * from './services/channelQueries';

// Store
export * from './store/instagramChannelSlice';

// Hooks
export { useInstagramOAuth } from './hooks/useInstagramOAuth';

// Types
export type * from './types/instagram.types';

// Validation
export * from './validation/instagramChannel.schema';
