/**
 * Facebook Login (Graph) — used to list Pages and obtain a user authorization code.
 * Configure VITE_META_APP_ID to match the Meta app used by the PostEaze API (META_APP_ID).
 */

const getRedirectUri = () =>
  import.meta.env.VITE_META_REDIRECT_URI || `${window.location.origin}/auth/facebook/callback`;

/** Scopes needed for page list + insights sync via stored page token */
const FB_SCOPES = [
  'pages_show_list',
  'pages_read_engagement',
  'read_insights',
  'business_management',
].join(',');

export const getFacebookAuthUrl = (): string => {
  const appId = import.meta.env.VITE_META_APP_ID || '';
  const params = new URLSearchParams({
    client_id: appId,
    redirect_uri: getRedirectUri(),
    scope: FB_SCOPES,
    response_type: 'code',
    state:
      typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function'
        ? `fb_${crypto.randomUUID()}`
        : `fb_${Date.now()}`,
  });
  return `https://www.facebook.com/v18.0/dialog/oauth?${params.toString()}`;
};

export const validateFacebookOAuthConfig = (): boolean => {
  if (!import.meta.env.VITE_META_APP_ID) {
    console.error('VITE_META_APP_ID is not configured');
    return false;
  }
  return true;
};

export const extractAuthCodeFromUrl = (): string | null => {
  const params = new URLSearchParams(window.location.search);
  return params.get('code');
};

export const extractErrorFromUrl = (): string | null => {
  const params = new URLSearchParams(window.location.search);
  return params.get('error_description') || params.get('error');
};

export const getFacebookRedirectUri = (): string => getRedirectUri();
