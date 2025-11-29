import { useState, useCallback } from 'react';
import { getInstagramAuthUrl, validateOAuthConfig } from '../utils/instagramOAuth.utils';

interface UseInstagramOAuthReturn {
  openOAuthPopup: () => Promise<string>;
  isLoading: boolean;
  error: string | null;
}

/**
 * Hook for managing Instagram OAuth popup flow
 */
export const useInstagramOAuth = (): UseInstagramOAuthReturn => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const openOAuthPopup = useCallback((): Promise<string> => {
    return new Promise((resolve, reject) => {
      // Validate configuration
      if (!validateOAuthConfig()) {
        const configError = 'Instagram OAuth is not properly configured. Please check your environment variables.';
        setError(configError);
        reject(new Error(configError));
        return;
      }

      setIsLoading(true);
      setError(null);

      const authUrl = getInstagramAuthUrl();
      const width = 600;
      const height = 700;
      const left = window.screen.width / 2 - width / 2;
      const top = window.screen.height / 2 - height / 2;

      const popup = window.open(
        authUrl,
        'Instagram Authorization',
        `width=${width},height=${height},left=${left},top=${top},toolbar=no,menubar=no,scrollbars=yes,resizable=yes`
      );

      if (!popup) {
        const popupError = 'Popup blocked. Please allow popups for this site and try again.';
        setError(popupError);
        setIsLoading(false);
        reject(new Error(popupError));
        return;
      }

      // Listen for message from popup
      const messageHandler = (event: MessageEvent) => {
        // Verify origin for security
        if (event.origin !== window.location.origin) {
          return;
        }

        if (event.data.type === 'INSTAGRAM_AUTH_CODE') {
          window.removeEventListener('message', messageHandler);
          clearInterval(checkPopupClosed);
          setIsLoading(false);

          if (event.data.code) {
            resolve(event.data.code);
          } else if (event.data.error) {
            const authError = event.data.error || 'Authorization failed';
            setError(authError);
            reject(new Error(authError));
          } else {
            const noCodeError = 'No authorization code received';
            setError(noCodeError);
            reject(new Error(noCodeError));
          }
        }
      };

      window.addEventListener('message', messageHandler);

      // Check if popup was closed manually
      const checkPopupClosed = setInterval(() => {
        if (popup.closed) {
          clearInterval(checkPopupClosed);
          window.removeEventListener('message', messageHandler);
          setIsLoading(false);
          
          const cancelledError = 'Authorization cancelled by user';
          setError(cancelledError);
          reject(new Error(cancelledError));
        }
      }, 500);
    });
  }, []);

  return {
    openOAuthPopup,
    isLoading,
    error,
  };
};
