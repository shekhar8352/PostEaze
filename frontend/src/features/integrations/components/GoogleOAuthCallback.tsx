import { useEffect } from 'react';
import { Box, Loader, Stack, Text } from '@mantine/core';
import {
  extractAuthCodeFromUrl,
  extractErrorFromUrl,
  getGoogleMessageType,
  type GoogleOAuthPurpose,
} from '../utils/googleOAuth.utils';

interface GoogleOAuthCallbackProps {
  purpose: GoogleOAuthPurpose;
  label: string;
}

export function GoogleOAuthCallback({ purpose, label }: GoogleOAuthCallbackProps) {
  const messageType = getGoogleMessageType(purpose);

  useEffect(() => {
    const code = extractAuthCodeFromUrl();
    const error = extractErrorFromUrl();

    if (window.opener) {
      if (code) {
        window.opener.postMessage({ type: messageType, code }, window.location.origin);
      } else if (error) {
        window.opener.postMessage(
          { type: messageType, error: error || 'Authorization failed' },
          window.location.origin
        );
      } else {
        window.opener.postMessage(
          { type: messageType, error: 'No authorization code received' },
          window.location.origin
        );
      }
      setTimeout(() => window.close(), 500);
    }
  }, [messageType]);

  return (
    <Box
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
      }}
    >
      <Stack align="center" gap="lg">
        <Loader size="lg" color="blue" />
        <Text size="lg" fw={500}>
          Completing {label} authorization…
        </Text>
        <Text size="sm" c="dimmed">
          This window will close automatically
        </Text>
      </Stack>
    </Box>
  );
}

export default GoogleOAuthCallback;
