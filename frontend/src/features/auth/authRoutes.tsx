import type { RouteObject } from 'react-router-dom';
import SignupPage from './pages/SignupPage';

const authRoutes: RouteObject[] = [
  { path: '/signup', element: <SignupPage /> },
];

export default authRoutes;