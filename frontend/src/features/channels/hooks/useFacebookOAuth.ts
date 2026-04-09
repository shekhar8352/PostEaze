import { useState, useCallback } from 'react';
import { getFacebookAuthUrl, validateFacebookOAuthConfig } from '../utils/facebookOAuth.utils';

interface UseFacebookOAuthReturn {
  openOAuthPopup: () => Promise<string>;
  isLoading: boolean;
  error: string | null;
}

export const useFacebookOAuth = (): UseFacebookOAuthReturn => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const openOAuthPopup = useCallback((): Promise<string> => {
    return new Promise((resolve, reject) => {
      if (!validateFacebookOAuthConfig()) {
        const configError =
          'Facebook OAuth is not configured. Set VITE_META_APP_ID (and optional VITE_META_REDIRECT_URI).';
        setError(configError);
        reject(new Error(configError));
        return;
      }

      setIsLoading(true);
      setError(null);

      const authUrl = getFacebookAuthUrl();
      const width = 600;
      const height = 700;
      const left = window.screen.width / 2 - width / 2;
      const top = window.screen.height / 2 - height / 2;

      const popup = window.open(
        authUrl,
        'Facebook Authorization',
        `width=${width},height=${height},left=${left},top=${top},toolbar=no,menubar=no,scrollbars=yes,resizable=yes`
      );

      if (!popup) {
        const popupError = 'Popup blocked. Allow popups for this site and try again.';
        setError(popupError);
        setIsLoading(false);
        reject(new Error(popupError));
        return;
      }

      const messageHandler = (event: MessageEvent) => {
        if (event.origin !== window.location.origin) return;

        if (event.data.type === 'FACEBOOK_AUTH_CODE') {
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

      const checkPopupClosed = setInterval(() => {
        if (popup.closed) {
          clearInterval(checkPopupClosed);
          window.removeEventListener('message', messageHandler);
          setIsLoading(false);
          const cancelledError = 'Authorization cancelled';
          setError(cancelledError);
          reject(new Error(cancelledError));
        }
      }, 500);
    });
  }, []);

  return { openOAuthPopup, isLoading, error };
};
