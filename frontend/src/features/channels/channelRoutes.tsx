import { lazy } from 'react';

const InstagramChannel = lazy(() => import('./pages/InstagramChannel'));
const FacebookChannel = lazy(() => import('./pages/FacebookChannel'));
const YouTubeChannel = lazy(() => import('./pages/YouTubeChannel'));

export const channelRoutes = [
    {
        path: '/channels/instagram',
        element: <InstagramChannel />,
    },
    {
        path: '/channels/facebook',
        element: <FacebookChannel />,
    },
    {
        path: '/channels/youtube',
        element: <YouTubeChannel />,
    },
];
