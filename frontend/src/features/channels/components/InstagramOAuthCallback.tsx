import { useEffect } from 'react';
import { Box, Text, Loader, Stack } from '@mantine/core';
import { extractAuthCodeFromUrl, extractErrorFromUrl } from '../utils/instagramOAuth.utils';

/**
 * Instagram OAuth Callback Component
 * Handles the OAuth redirect from Instagram and sends the code back to the parent window
 */
export const InstagramOAuthCallback = () => {
    useEffect(() => {
        const code = extractAuthCodeFromUrl();
        const error = extractErrorFromUrl();

        if (window.opener) {
            if (code) {
                // Send authorization code to parent window
                window.opener.postMessage(
                    {
                        type: 'INSTAGRAM_AUTH_CODE',
                        code,
                    },
                    window.location.origin
                );
            } else if (error) {
                // Send error to parent window
                window.opener.postMessage(
                    {
                        type: 'INSTAGRAM_AUTH_CODE',
                        error: error || 'Authorization failed',
                    },
                    window.location.origin
                );
            } else {
                // No code or error found
                window.opener.postMessage(
                    {
                        type: 'INSTAGRAM_AUTH_CODE',
                        error: 'No authorization code received',
                    },
                    window.location.origin
                );
            }

            // Close popup after sending message
            setTimeout(() => {
                window.close();
            }, 500);
        }
    }, []);

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
                <Loader size="lg" color="grape" />
                <Text size="lg" fw={500}>
                    Completing Instagram authorization...
                </Text>
                <Text size="sm" c="dimmed">
                    This window will close automatically
                </Text>
            </Stack>
        </Box>
    );
};

export default InstagramOAuthCallback;
