import { lazy, Suspense } from 'react';
import type { RouteObject } from 'react-router-dom';
import { LoadingFallback } from '@/app/routes/LoadingFallback';

const DashboardPage = lazy(() => import('./pages/Dashboard'));

const dashboardRoutes: RouteObject[] = [
    {
        path: '/dashboard',
        element: (
            <Suspense fallback={<LoadingFallback />}>
                <DashboardPage />
            </Suspense>
        ),
    },
];

export default dashboardRoutes;
