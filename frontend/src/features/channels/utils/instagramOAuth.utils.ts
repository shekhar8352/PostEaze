import type { InstagramOAuthConfig } from '../types/instagram.types';

// Instagram OAuth Configuration
const INSTAGRAM_CONFIG: InstagramOAuthConfig = {
  clientId: import.meta.env.VITE_INSTAGRAM_CLIENT_ID || '',
  redirectUri:
    import.meta.env.VITE_INSTAGRAM_REDIRECT_URI ||
    `${window.location.origin}/auth/instagram/callback`,
  scope: 'instagram_business_basic, instagram_business_manage_messages, instagram_business_manage_comments, instagram_business_content_publish, instagram_business_manage_insights',
  responseType: 'code',
};

/**
 * Generate Instagram OAuth authorization URL
 */
export const getInstagramAuthUrl = (): string => {
  const params = new URLSearchParams({
    client_id: INSTAGRAM_CONFIG.clientId,
    redirect_uri: INSTAGRAM_CONFIG.redirectUri,
    scope: INSTAGRAM_CONFIG.scope,
    response_type: INSTAGRAM_CONFIG.responseType,
  });

  return `https://api.instagram.com/oauth/authorize?${params.toString()}`;
};

/**
 * Extract authorization code from URL query parameters
 */
export const extractAuthCodeFromUrl = (): string | null => {
  const params = new URLSearchParams(window.location.search);
  return params.get('code');
};

/**
 * Extract error from URL query parameters
 */
export const extractErrorFromUrl = (): string | null => {
  const params = new URLSearchParams(window.location.search);
  return params.get('error') || params.get('error_description');
};

/**
 * Validate Instagram OAuth configuration
 */
export const validateOAuthConfig = (): boolean => {
  if (!INSTAGRAM_CONFIG.clientId) {
    console.error('Instagram Client ID is not configured');
    return false;
  }
  return true;
};

/**
 * Get OAuth configuration (for debugging)
 */
export const getOAuthConfig = (): InstagramOAuthConfig => {
  return { ...INSTAGRAM_CONFIG };
};
