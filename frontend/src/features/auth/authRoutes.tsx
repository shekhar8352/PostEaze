import { lazy, Suspense } from 'react';
import type { RouteObject } from "react-router-dom";
import PublicRoute from '@/app/routes/PublicRoute';
import { LoadingFallback } from '@/app/routes/LoadingFallback';

const LoginPage = lazy(() => import('./pages/Login'));
const RegisterPage = lazy(() => import('./pages/Register'));
const ForgotPasswordPage = lazy(() => import('./pages/ForgotPassword'));
const EmailVerificationPage = lazy(() => import('./pages/EmailVerfication'));

const authRoutes: RouteObject[] = [
  {
    path: '/login',
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <PublicRoute element={<LoginPage />} />
      </Suspense>
    ),
  },
  {
    path: '/register',
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <PublicRoute element={<RegisterPage />} />
      </Suspense>
    ),
  },
  {
    path: '/forgot-password',
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <PublicRoute element={<ForgotPasswordPage />} />
      </Suspense>
    ),
  },
  {
    path: '/email-verify',
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <EmailVerificationPage />
      </Suspense>
    ),
  },
];

export default authRoutes;
