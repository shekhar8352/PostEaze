import { useEffect } from 'react';
import { Box, Text, Loader, Stack } from '@mantine/core';
import { extractAuthCodeFromUrl, extractErrorFromUrl } from '../utils/facebookOAuth.utils';

/**
 * OAuth redirect target for Facebook Login. Posts the code to the opener and closes.
 */
export const FacebookOAuthCallback = () => {
  useEffect(() => {
    const code = extractAuthCodeFromUrl();
    const error = extractErrorFromUrl();

    if (window.opener) {
      if (code) {
        window.opener.postMessage(
          { type: 'FACEBOOK_AUTH_CODE', code },
          window.location.origin
        );
      } else if (error) {
        window.opener.postMessage(
          { type: 'FACEBOOK_AUTH_CODE', error: error || 'Authorization failed' },
          window.location.origin
        );
      } else {
        window.opener.postMessage(
          { type: 'FACEBOOK_AUTH_CODE', error: 'No authorization code received' },
          window.location.origin
        );
      }
      setTimeout(() => window.close(), 400);
    }
  }, []);

  return (
    <Box
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(160deg, #0f172a 0%, #1e3a5f 45%, #0c4a6e 100%)',
      }}
    >
      <Stack align="center" gap="lg">
        <Loader size="lg" color="cyan" />
        <Text size="lg" fw={600} c="white">
          Completing Facebook authorization…
        </Text>
        <Text size="sm" c="rgba(255,255,255,0.7)">
          This window will close automatically.
        </Text>
      </Stack>
    </Box>
  );
};

export default FacebookOAuthCallback;
