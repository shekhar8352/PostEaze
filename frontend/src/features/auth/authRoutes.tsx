import { lazy } from 'react'
import type { RouteObject } from "react-router-dom";


const LoginPage = lazy(()=> import('./pages/Login'))
const RegisterPage = lazy(()=> import('./pages/Register'))
const ForgotPasswordPage = lazy(() => import('./pages/ForgotPassword'));
const EmailVerificationPage = lazy(() => import('./pages/EmailVerfication'));

const authRoutes: RouteObject[] = [
  {
    path: '/login',
    element: <LoginPage/>
  },
  {
    path: '/register',
    element: <RegisterPage/>
  },
  {
    path: '/forgot-password',
    element: <ForgotPasswordPage/>
  },
  {
    path: '/email-verify',
    element: <EmailVerificationPage/>
  }
];

export default authRoutes;
