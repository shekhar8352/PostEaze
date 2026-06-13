export type GoogleOAuthPurpose = 'drive' | 'youtube';

const CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID || '';

const REDIRECT_BY_PURPOSE: Record<GoogleOAuthPurpose, string> = {
  drive:
    import.meta.env.VITE_GOOGLE_DRIVE_REDIRECT_URI ||
    `${window.location.origin}/oauth/google/drive/callback`,
  youtube:
    import.meta.env.VITE_GOOGLE_YOUTUBE_REDIRECT_URI ||
    `${window.location.origin}/oauth/google/youtube/callback`,
};

const SCOPES_BY_PURPOSE: Record<GoogleOAuthPurpose, string> = {
  drive: [
    'https://www.googleapis.com/auth/drive.readonly',
    'https://www.googleapis.com/auth/userinfo.email',
  ].join(' '),
  youtube: [
    'https://www.googleapis.com/auth/youtube.upload',
    'https://www.googleapis.com/auth/youtube.readonly',
    'https://www.googleapis.com/auth/userinfo.email',
  ].join(' '),
};

const MESSAGE_TYPE_BY_PURPOSE: Record<GoogleOAuthPurpose, string> = {
  drive: 'GOOGLE_DRIVE_AUTH_CODE',
  youtube: 'GOOGLE_YOUTUBE_AUTH_CODE',
};

export function getGoogleRedirectUri(purpose: GoogleOAuthPurpose): string {
  return REDIRECT_BY_PURPOSE[purpose];
}

export function getGoogleAuthUrl(purpose: GoogleOAuthPurpose): string {
  const params = new URLSearchParams({
    client_id: CLIENT_ID,
    redirect_uri: REDIRECT_BY_PURPOSE[purpose],
    response_type: 'code',
    scope: SCOPES_BY_PURPOSE[purpose],
    access_type: 'offline',
    prompt: 'consent',
  });
  return `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`;
}

export function getGoogleMessageType(purpose: GoogleOAuthPurpose): string {
  return MESSAGE_TYPE_BY_PURPOSE[purpose];
}

export function extractAuthCodeFromUrl(): string | null {
  const params = new URLSearchParams(window.location.search);
  return params.get('code');
}

export function extractErrorFromUrl(): string | null {
  const params = new URLSearchParams(window.location.search);
  return params.get('error') || params.get('error_description');
}

export function validateGoogleOAuthConfig(): boolean {
  if (!CLIENT_ID) {
    console.error('VITE_GOOGLE_CLIENT_ID is not configured');
    return false;
  }
  return true;
}
