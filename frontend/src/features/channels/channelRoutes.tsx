import { lazy, Suspense } from 'react';
import type { RouteObject } from 'react-router-dom';
import { LoadingFallback } from '@/app/routes/LoadingFallback';

const InstagramChannelPage = lazy(() => import('./pages/InstagramChannelPage'));
const InstagramOAuthCallback = lazy(() => import('./components/InstagramOAuthCallback'));
const FacebookChannel = lazy(() => import('./pages/FacebookChannel'));
const FacebookOAuthCallback = lazy(() => import('./components/FacebookOAuthCallback'));
const YouTubeChannelPage = lazy(() => import('./pages/YouTubeChannelPage'));
const GoogleDriveOAuthCallback = lazy(() =>
    import('@/features/integrations/components/GoogleOAuthCallback').then((m) => ({
        default: () => <m.GoogleOAuthCallback purpose="drive" label="Google Drive" />,
    }))
);
const GoogleYouTubeOAuthCallback = lazy(() =>
    import('@/features/integrations/components/GoogleOAuthCallback').then((m) => ({
        default: () => <m.GoogleOAuthCallback purpose="youtube" label="YouTube" />,
    }))
);

const channelRoutes: RouteObject[] = [
    {
        path: '/channels/instagram',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <InstagramChannelPage />
            </Suspense>
        ),
    },
    {
        path: '/auth/instagram/callback',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <InstagramOAuthCallback />
            </Suspense>
        ),
    },
    {
        path: '/channels/facebook',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <FacebookChannel />
            </Suspense>
        ),
    },
    {
        path: '/auth/facebook/callback',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <FacebookOAuthCallback />
            </Suspense>
        ),
    },
    {
        path: '/channels/youtube',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <YouTubeChannelPage />
            </Suspense>
        ),
    },
    {
        path: '/oauth/google/drive/callback',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <GoogleDriveOAuthCallback />
            </Suspense>
        ),
    },
    {
        path: '/oauth/google/youtube/callback',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <GoogleYouTubeOAuthCallback />
            </Suspense>
        ),
    },
];

export default channelRoutes;
