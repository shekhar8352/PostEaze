import { lazy, Suspense } from 'react';
import type { RouteObject } from 'react-router-dom';
import { LoadingFallback } from '@/app/routes/LoadingFallback';

const InstagramChannelPage = lazy(() => import('./pages/InstagramChannelPage'));
const InstagramOAuthCallback = lazy(() => import('./components/InstagramOAuthCallback'));

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
];

export default channelRoutes;
